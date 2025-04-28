package wsapi

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"crypto/x509"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Common errors
var (
	ErrConnectionClosed = errors.New("connection closed")
	ErrClientClosed     = errors.New("client closed by user")
)

// LogonRequest represents a session.logon request to the Binance WebSocket API
type LogonRequest struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Params LogonParams `json:"params"`
}

// LogonParams contains the parameters for a session.logon request
type LogonParams struct {
	APIKey     string `json:"apiKey"`
	Signature  string `json:"signature"`
	Timestamp  int64  `json:"timestamp"`
	RecvWindow int64  `json:"recvWindow,omitempty"`
}

// SessionStatusRequest represents a session.status request to the Binance WebSocket API
type SessionStatusRequest struct {
	ID     string `json:"id"`
	Method string `json:"method"`
}

// BinanceWebsocketClient defines the interface for Binance WebSocket API client
type BinanceWebsocketClient interface {
	// AsyncDo sends a message to the WebSocket server asynchronously
	AsyncDo(msg []byte) error

	// Close closes the WebSocket connection and cleans up resources
	// It is idempotent - can be called multiple times safely
	Close() error
}

// BinanceWSClient implements the BinanceWebsocketClient interface
type BinanceWSClient struct {
	endpoint string
	conn     *websocket.Conn
	sendMu   sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	dialer   websocket.Dialer

	// Reconnection settings
	backoffBase time.Duration
	backoffMax  time.Duration

	// Authentication credentials
	apiKey         string
	privateKeyPath string // Path to the ED25519 private key file

	// Post-connection hooks
	postConnectHooks []func() error

	// Handler function for received messages
	messageHandler func(int, []byte)

	// Flag to track if client is closed
	closed   bool
	closedMu sync.RWMutex
}

// BinanceWSClientOption defines function type for client configuration options
type BinanceWSClientOption func(*BinanceWSClient)

// WithMessageHandler sets a custom message handler function
func WithMessageHandler(handler func(int, []byte)) BinanceWSClientOption {
	return func(c *BinanceWSClient) {
		c.messageHandler = handler
	}
}

// WithTestnet configures the client to use the testnet endpoint
func WithTestnet() BinanceWSClientOption {
	return func(c *BinanceWSClient) {
		c.endpoint = "wss://ws-api.testnet.binance.vision/ws-api/v3"
	}
}

// WithAlternativePort configures the client to use the alternative port 9443
func WithAlternativePort() BinanceWSClientOption {
	return func(c *BinanceWSClient) {
		if c.endpoint == "wss://ws-api.testnet.binance.vision/ws-api/v3" {
			c.endpoint = "wss://ws-api.testnet.binance.vision:9443/ws-api/v3"
		} else {
			c.endpoint = "wss://ws-api.binance.com:9443/ws-api/v3"
		}
	}
}

// WithBackoffSettings configures reconnection backoff settings
func WithBackoffSettings(base, max time.Duration) BinanceWSClientOption {
	return func(c *BinanceWSClient) {
		c.backoffBase = base
		c.backoffMax = max
	}
}

// WithLogon configures the client with API credentials and performs logon
func WithLogon(apiKey, privateKeyPath string) BinanceWSClientOption {
	return func(c *BinanceWSClient) {
		// Store credentials for potential reconnections
		c.apiKey = apiKey
		c.privateKeyPath = privateKeyPath

		// Add a post-connect hook to send logon request after connection is established
		c.postConnectHooks = append(c.postConnectHooks, func() error {
			return c.sendLogonRequest()
		})
	}
}

// WithCheckSessionStatus adds a session status check after logon
func WithCheckSessionStatus(delay time.Duration) BinanceWSClientOption {
	return func(c *BinanceWSClient) {
		if delay < 0 {
			delay = 1 * time.Second // Default delay if invalid value provided
		}

		// Add a post-connect hook that runs after logon (if delay > 0)
		if delay > 0 {
			c.postConnectHooks = append(c.postConnectHooks, func() error {
				// Wait a bit for logon to complete before checking status
				time.Sleep(delay)
				return c.SendSessionStatus()
			})
		}
	}
}

// NewBinanceWebsocketClient creates a new Binance WebSocket client
// It establishes a connection to the Binance WebSocket API
func NewBinanceWebsocketClient(opts ...BinanceWSClientOption) (BinanceWebsocketClient, error) {
	ctx, cancel := context.WithCancel(context.Background())

	client := &BinanceWSClient{
		endpoint:         "wss://ws-api.binance.com:443/ws-api/v3",
		ctx:              ctx,
		cancel:           cancel,
		dialer:           websocket.Dialer{},
		backoffBase:      1 * time.Second,
		backoffMax:       30 * time.Second,
		messageHandler:   defaultMessageHandler,
		postConnectHooks: []func() error{},
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	// Establish initial connection
	if err := client.connect(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to establish connection: %w", err)
	}

	// Setup handlers
	client.setupPingHandler()
	client.setupPongHandler()

	// Run post-connect hooks
	for _, hook := range client.postConnectHooks {
		if err := hook(); err != nil {
			client.Close()
			return nil, fmt.Errorf("post-connect hook failed: %w", err)
		}
	}

	// Start read loop in a goroutine
	client.wg.Add(1)
	go client.readLoop()

	return client, nil
}

// defaultMessageHandler is the default implementation for processing received messages
func defaultMessageHandler(msgType int, data []byte) {
	fmt.Printf("Received message (type %d): %s\n", msgType, data)
}

// connect establishes a WebSocket connection to the Binance API endpoint
func (c *BinanceWSClient) connect() error {
	u, _ := url.Parse(c.endpoint)

	// Use a timeout for the dial attempt
	dialCtx, dialCancel := context.WithTimeout(c.ctx, 10*time.Second)
	defer dialCancel()

	conn, _, err := c.dialer.DialContext(dialCtx, u.String(), nil)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}

	c.conn = conn

	// Set the initial read deadline
	c.resetReadDeadline()

	return nil
}

// setupPingHandler configures the handler for ping frames
func (c *BinanceWSClient) setupPingHandler() {
	c.conn.SetPingHandler(func(appData string) error {
		// When we receive a ping, we must respond with a pong with the same payload
		// as stated in the Binance API documentation
		err := c.conn.WriteControl(
			websocket.PongMessage,
			[]byte(appData),
			time.Now().Add(5*time.Second),
		)
		if err != nil {
			fmt.Printf("Failed to send pong: %v\n", err)
		}

		// Extend the read deadline when we get a ping
		c.resetReadDeadline()
		return nil
	})
}

// setupPongHandler configures the handler for pong frames
func (c *BinanceWSClient) setupPongHandler() {
	c.conn.SetPongHandler(func(appData string) error {
		if appData != "" {
			fmt.Printf("Warning: unsolicited pong frame with payload: %s\n", appData)
		}

		// Extend the read deadline when we get a pong
		c.resetReadDeadline()
		return nil
	})
}

// resetReadDeadline extends the read deadline
// Binance sends a ping every 20 seconds and expects a pong within 60 seconds
func (c *BinanceWSClient) resetReadDeadline() {
	// We use 65 seconds to be safe (slightly longer than the 60s requirement)
	c.conn.SetReadDeadline(time.Now().Add(65 * time.Second))
}

// AsyncDo sends a message to the WebSocket server asynchronously
func (c *BinanceWSClient) AsyncDo(msg []byte) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	c.closedMu.RLock()
	closed := c.closed
	c.closedMu.RUnlock()

	if closed || c.conn == nil {
		return ErrClientClosed
	}

	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

// Close closes the WebSocket connection and cleans up resources
func (c *BinanceWSClient) Close() error {
	c.closedMu.Lock()
	if c.closed {
		c.closedMu.Unlock()
		return nil // Already closed, return without error (idempotent)
	}
	c.closed = true
	c.closedMu.Unlock()

	// Cancel context to signal all goroutines to stop
	c.cancel()

	// Wait for goroutines to finish
	c.wg.Wait()

	// Close the WebSocket connection if it exists
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// readLoop continuously reads messages from the WebSocket connection
func (c *BinanceWSClient) readLoop() {
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// Continue reading
		}

		// ReadMessage blocks until a message is received or an error occurs
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			c.closedMu.RLock()
			closed := c.closed
			c.closedMu.RUnlock()

			if closed {
				// Client is closed, exit gracefully
				return
			}

			fmt.Printf("Read error: %v - attempting reconnect\n", err)
			if reconnectErr := c.reconnectWithBackoff(); reconnectErr != nil {
				if errors.Is(reconnectErr, ErrClientClosed) {
					// Client was closed during reconnection attempt
					return
				}
				fmt.Printf("Failed to reconnect: %v - exiting read loop\n", reconnectErr)
				return
			}

			// Reset handlers after reconnection
			c.setupPingHandler()
			c.setupPongHandler()
			continue
		}

		// Process the received message through the handler
		if c.messageHandler != nil {
			go c.messageHandler(msgType, msg)
		}
	}
}

// sendLogonRequest sends a session.logon request with the API credentials
func (c *BinanceWSClient) sendLogonRequest() error {
	if c.apiKey == "" || c.privateKeyPath == "" {
		return errors.New("API key and private key path are required for logon")
	}

	// Generate a unique request ID
	requestID := uuid.New().String()

	// Get current timestamp in milliseconds
	timestamp := time.Now().UnixMilli()

	// Prepare parameters for signature
	params := map[string]string{
		"apiKey":    c.apiKey,
		"timestamp": fmt.Sprintf("%d", timestamp),
	}

	// Generate signature
	signature := c.generateSignature(params)

	// Create logon request
	logonRequest := LogonRequest{
		ID:     requestID,
		Method: "session.logon",
		Params: LogonParams{
			APIKey:    c.apiKey,
			Signature: signature,
			Timestamp: timestamp,
		},
	}

	// Marshal request to JSON
	requestJSON, err := json.Marshal(logonRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal logon request: %w", err)
	}

	// Send logon request
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	if c.conn == nil {
		return errors.New("connection is not established")
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, requestJSON); err != nil {
		return fmt.Errorf("failed to send logon request: %w", err)
	}

	fmt.Println("Sent session.logon request")
	return nil
}

// SendSessionStatus sends a session.status request to check the authentication status
func (c *BinanceWSClient) SendSessionStatus() error {
	// Generate a unique request ID
	requestID := uuid.New().String()

	// Create session status request
	statusRequest := SessionStatusRequest{
		ID:     requestID,
		Method: "session.status",
	}

	// Marshal request to JSON
	requestJSON, err := json.Marshal(statusRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal session.status request: %w", err)
	}

	// Send the request
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	if c.conn == nil {
		return errors.New("connection is not established")
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, requestJSON); err != nil {
		return fmt.Errorf("failed to send session.status request: %w", err)
	}

	fmt.Println("Sent session.status request")
	return nil
}

// generateSignature creates a signature for Binance API authentication using ED25519
func (c *BinanceWSClient) generateSignature(params map[string]string) string {
	// Sort parameters by key
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string
	var sb strings.Builder
	for i, key := range keys {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(params[key])
	}
	queryString := sb.String()

	// Read ED25519 private key from file
	privateKeyBytes, err := os.ReadFile(c.privateKeyPath)
	if err != nil {
		fmt.Printf("Failed to read private key file: %v\n", err)
		return ""
	}

	// Try to parse the private key - handle different formats
	var privateKey ed25519.PrivateKey

	// First try direct parsing
	if len(privateKeyBytes) == ed25519.PrivateKeySize {
		privateKey = ed25519.PrivateKey(privateKeyBytes)
	} else {
		// The file might contain PEM encoded data
		// Try to decode PEM format
		block, _ := pem.Decode(privateKeyBytes)
		if block != nil {
			// If it's PEM encoded, use the decoded data
			privateKeyBytes = block.Bytes

			// Try to parse as PKCS8
			if key, err := x509.ParsePKCS8PrivateKey(privateKeyBytes); err == nil {
				if edKey, ok := key.(ed25519.PrivateKey); ok {
					privateKey = edKey
				}
			}

			// If that didn't work, try as raw ED25519 key
			if privateKey == nil && len(privateKeyBytes) == ed25519.PrivateKeySize {
				privateKey = ed25519.PrivateKey(privateKeyBytes)
			}
		}
	}

	// If we still don't have a valid key, try to handle hex encoded key
	if privateKey == nil {
		// Try to decode as hex string
		hexDecoded, err := hex.DecodeString(strings.TrimSpace(string(privateKeyBytes)))
		if err == nil && len(hexDecoded) == ed25519.PrivateKeySize {
			privateKey = ed25519.PrivateKey(hexDecoded)
		} else if err == nil && len(hexDecoded) == ed25519.SeedSize {
			// Maybe it's just the seed part (32 bytes)
			privateKey = ed25519.NewKeyFromSeed(hexDecoded)
		}
	}

	// If we still don't have a valid key, report an error
	if privateKey == nil {
		fmt.Printf("Invalid ED25519 private key format: length=%d bytes\n", len(privateKeyBytes))
		return ""
	}

	// For debugging - output key information
	fmt.Printf("Using ED25519 private key, length: %d bytes\n", len(privateKey))

	// Sign the data with ED25519
	// The message to be signed must be the raw query string bytes
	signature := ed25519.Sign(privateKey, []byte(queryString))

	// Return base64 encoded signature (instead of hex encoding)
	// This matches the Python example that uses base64.b64encode
	return base64.StdEncoding.EncodeToString(signature)
}

// reconnectWithBackoff attempts to reconnect to the WebSocket server with exponential backoff
func (c *BinanceWSClient) reconnectWithBackoff() error {
	backoff := c.backoffBase

	for {
		select {
		case <-c.ctx.Done():
			return ErrClientClosed
		default:
			// Continue reconnection attempts
		}

		fmt.Printf("Reconnecting in %v...\n", backoff)

		// Wait before retrying
		select {
		case <-c.ctx.Done():
			return ErrClientClosed
		case <-time.After(backoff):
			// Continue after backoff period
		}

		// Try to connect
		if err := c.connect(); err != nil {
			fmt.Printf("Reconnection attempt failed: %v\n", err)

			// Increase backoff for next attempt, capped at maximum
			backoff *= 2
			if backoff > c.backoffMax {
				backoff = c.backoffMax
			}
			continue
		}

		// Reset handlers after reconnection
		c.setupPingHandler()
		c.setupPongHandler()

		// Run post-connect hooks after reconnection
		for _, hook := range c.postConnectHooks {
			if err := hook(); err != nil {
				fmt.Printf("Post-connect hook failed after reconnection: %v\n", err)
				// Close the connection and try again
				if c.conn != nil {
					c.conn.Close()
				}
				continue
			}
		}

		fmt.Println("Successfully reconnected to Binance WebSocket API")
		return nil
	}
}
