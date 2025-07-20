package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSClient represents a WebSocket client for Binance
type WSClient struct {
	conn           *websocket.Conn
	url            string
	config         *Config
	mu             sync.RWMutex
	isConnected    bool
	isReconnecting bool
	closeChan      chan struct{}
	reconnectChan  chan struct{}

	// Callbacks
	onConnect    func()
	onDisconnect func()
	onError      func(error)
	onMessage    func([]byte)

	// Reconnection settings
	maxReconnectAttempts int
	reconnectDelay       time.Duration
	pingInterval         time.Duration
	pongWait             time.Duration
	writeWait            time.Duration
}

// WSClientOption represents a configuration option for WSClient
type WSClientOption func(*WSClient)

// WithOnConnect sets the onConnect callback
func WithOnConnect(callback func()) WSClientOption {
	return func(c *WSClient) {
		c.onConnect = callback
	}
}

// WithOnDisconnect sets the onDisconnect callback
func WithOnDisconnect(callback func()) WSClientOption {
	return func(c *WSClient) {
		c.onDisconnect = callback
	}
}

// WithOnError sets the onError callback
func WithOnError(callback func(error)) WSClientOption {
	return func(c *WSClient) {
		c.onError = callback
	}
}

// WithOnMessage sets the onMessage callback
func WithOnMessage(callback func([]byte)) WSClientOption {
	return func(c *WSClient) {
		c.onMessage = callback
	}
}

// WithReconnectSettings sets reconnection settings
func WithReconnectSettings(maxAttempts int, delay time.Duration) WSClientOption {
	return func(c *WSClient) {
		c.maxReconnectAttempts = maxAttempts
		c.reconnectDelay = delay
	}
}

// NewWSClient creates a new WebSocket client
func NewWSClient(config *Config, options ...WSClientOption) *WSClient {
	client := &WSClient{
		config:               config,
		closeChan:            make(chan struct{}),
		reconnectChan:        make(chan struct{}, 1),
		maxReconnectAttempts: 5,
		reconnectDelay:       5 * time.Second,
		pingInterval:         30 * time.Second,
		pongWait:             60 * time.Second,
		writeWait:            10 * time.Second,
	}

	// Set WebSocket URL based on config
	if config.UseTestnet {
		client.url = WSBaseURLTestnet
	} else {
		client.url = WSBaseURL
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client
}

// Connect establishes a WebSocket connection
func (c *WSClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return fmt.Errorf("already connected")
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket %s: %w", c.url, err)
	}

	c.conn = conn
	c.isConnected = true

	// Set connection parameters
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
		return nil
	})

	// Start goroutines for handling connection
	go c.readPump()
	go c.pingPump()

	// Call onConnect callback
	if c.onConnect != nil {
		c.onConnect()
	}

	return nil
}

// Disconnect closes the WebSocket connection
func (c *WSClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return nil
	}

	// Signal close (use select to avoid closing already closed channel)
	select {
	case <-c.closeChan:
		// Channel already closed
	default:
		close(c.closeChan)
	}

	// Close connection
	if c.conn != nil {
		err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Printf("error sending close message: %v", err)
		}
		c.conn.Close()
	}

	c.isConnected = false
	c.isReconnecting = false

	// Call onDisconnect callback
	if c.onDisconnect != nil {
		c.onDisconnect()
	}

	return nil
}

// IsConnected returns whether the client is currently connected
func (c *WSClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// readPump handles reading messages from the WebSocket connection
func (c *WSClient) readPump() {
	defer func() {
		c.mu.Lock()
		c.isConnected = false
		c.mu.Unlock()
		c.triggerReconnect()
	}()

	for {
		select {
		case <-c.closeChan:
			return
		default:
			c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}

			// Call onMessage callback
			if c.onMessage != nil {
				c.onMessage(message)
			}
		}
	}
}

// pingPump sends periodic ping messages to keep the connection alive
func (c *WSClient) pingPump() {
	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeChan:
			return
		case <-ticker.C:
			c.mu.RLock()
			if !c.isConnected || c.conn == nil {
				c.mu.RUnlock()
				return
			}
			c.mu.RUnlock()

			c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, nil); err != nil {
				log.Printf("ping error: %v", err)
				return
			}
		}
	}
}

// triggerReconnect initiates reconnection with exponential backoff
func (c *WSClient) triggerReconnect() {
	c.mu.Lock()
	if c.isReconnecting {
		c.mu.Unlock()
		return
	}
	c.isReconnecting = true
	c.mu.Unlock()

	// Send reconnect signal
	select {
	case c.reconnectChan <- struct{}{}:
	default:
	}

	go c.reconnect()
}

// reconnect handles reconnection with exponential backoff
func (c *WSClient) reconnect() {
	attempt := 0
	for attempt < c.maxReconnectAttempts {
		attempt++

		// Calculate delay with exponential backoff
		delay := time.Duration(math.Pow(2, float64(attempt-1))) * c.reconnectDelay
		if delay > 60*time.Second {
			delay = 60 * time.Second
		}

		log.Printf("Attempting to reconnect in %v (attempt %d/%d)", delay, attempt, c.maxReconnectAttempts)
		time.Sleep(delay)

		// Try to connect
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := c.Connect(ctx)
		cancel()

		if err == nil {
			log.Printf("Successfully reconnected after %d attempts", attempt)
			c.mu.Lock()
			c.isReconnecting = false
			c.mu.Unlock()
			return
		}

		log.Printf("Reconnection attempt %d failed: %v", attempt, err)
	}

	log.Printf("Failed to reconnect after %d attempts", c.maxReconnectAttempts)
	c.mu.Lock()
	c.isReconnecting = false
	c.mu.Unlock()

	// Call onError callback
	if c.onError != nil {
		c.onError(fmt.Errorf("failed to reconnect after %d attempts", c.maxReconnectAttempts))
	}
}

// SendMessage sends a message through the WebSocket connection
func (c *WSClient) SendMessage(message interface{}) error {
	c.mu.RLock()
	if !c.isConnected || c.conn == nil {
		c.mu.RUnlock()
		return fmt.Errorf("not connected")
	}
	c.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// SubscribeToStream subscribes to a specific stream
func (c *WSClient) SubscribeToStream(streamName string) error {
	// For raw streams, we connect directly to the stream URL
	streamURL := fmt.Sprintf("%s/%s", c.url, streamName)

	// Close existing connection if any
	if c.isConnected {
		c.Disconnect()
	}

	// Update URL and connect
	c.url = streamURL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.Connect(ctx)
}

// UnsubscribeFromStream unsubscribes from a stream by disconnecting
func (c *WSClient) UnsubscribeFromStream() error {
	return c.Disconnect()
}
