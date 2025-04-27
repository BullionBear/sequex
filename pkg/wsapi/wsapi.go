package wsapi

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Common errors
var (
	ErrConnectionClosed = errors.New("connection closed")
	ErrClientClosed     = errors.New("client closed by user")
)

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

// NewBinanceWebsocketClient creates a new Binance WebSocket client
// It establishes a connection to the Binance WebSocket API
func NewBinanceWebsocketClient(opts ...BinanceWSClientOption) (BinanceWebsocketClient, error) {
	ctx, cancel := context.WithCancel(context.Background())

	client := &BinanceWSClient{
		endpoint:       "wss://ws-api.binance.com:443/ws-api/v3",
		ctx:            ctx,
		cancel:         cancel,
		dialer:         websocket.Dialer{},
		backoffBase:    1 * time.Second,
		backoffMax:     30 * time.Second,
		messageHandler: defaultMessageHandler,
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

		fmt.Println("Successfully reconnected to Binance WebSocket API")
		return nil
	}
}
