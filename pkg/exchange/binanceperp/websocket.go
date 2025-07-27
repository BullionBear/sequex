package binanceperp

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSConnection interface for all WebSocket connection types
type WSConnection interface {
	Connect(ctx context.Context, streamName string) error
	Disconnect() error
	IsConnected() bool
}

// Ping interval constants for Binance futures WebSocket
const (
	pingInterval      = 3 * time.Minute  // Binance sends ping every 3 minutes
	pongTimeout       = 10 * time.Minute // Disconnect if no pong received within 10 minutes
	reconnectDelay    = 5 * time.Second  // Delay between reconnection attempts
	connectionTimeout = 24 * time.Hour   // Connection valid for 24 hours
)

// WSConfig holds WebSocket client configuration
type WSConfig struct {
	BaseWSUrl      string
	ReconnectDelay time.Duration
	PingInterval   time.Duration
	MaxReconnects  int // -1 means no max reconnects
}

// Subscription provides a builder pattern for configuring WebSocket stream callbacks
type Subscription struct {
	onConnect   func()
	onReconnect func()
	onError     func(error)
	onMessage   func([]byte)
	onClose     func()
}

// WithConnect sets the OnConnect callback
func (s *Subscription) WithConnect(onConnect func()) *Subscription {
	s.onConnect = onConnect
	return s
}

// WithReconnect sets the OnReconnect callback
func (s *Subscription) WithReconnect(onReconnect func()) *Subscription {
	s.onReconnect = onReconnect
	return s
}

// WithError sets the OnError callback
func (s *Subscription) WithError(onError func(error)) *Subscription {
	s.onError = onError
	return s
}

// WithMessage sets the OnMessage callback
func (s *Subscription) WithMessage(onMessage func([]byte)) *Subscription {
	s.onMessage = onMessage
	return s
}

// WithClose sets the OnClose callback
func (s *Subscription) WithClose(onClose func()) *Subscription {
	s.onClose = onClose
	return s
}

// BinancePerpWSConn manages WebSocket connection to Binance perpetual futures streams
type BinancePerpWSConn struct {
	conn         *websocket.Conn
	mu           sync.RWMutex
	done         chan struct{}
	reconnect    chan struct{}
	logger       *log.Logger
	config       *WSConfig
	subscription *Subscription

	// Connection state
	connected       bool
	streamName      string
	ctx             context.Context
	cancel          context.CancelFunc
	shouldReconnect bool
	reconnectCount  int
}

// NewBinancePerpWSConn creates a new WebSocket client instance
func NewBinancePerpWSConn(config *WSConfig, subscription *Subscription) *BinancePerpWSConn {
	if config == nil {
		config = &WSConfig{
			BaseWSUrl:      MainnetWSBaseUrl,
			ReconnectDelay: reconnectDelay,
			PingInterval:   pingInterval,
			MaxReconnects:  -1, // No max reconnects by default
		}
	}

	// Set defaults if not provided
	if config.BaseWSUrl == "" {
		config.BaseWSUrl = MainnetWSBaseUrl
	}
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = reconnectDelay
	}
	if config.PingInterval == 0 {
		config.PingInterval = pingInterval
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &BinancePerpWSConn{
		config:          config,
		subscription:    subscription,
		ctx:             ctx,
		cancel:          cancel,
		done:            make(chan struct{}),
		reconnect:       make(chan struct{}),
		logger:          log.Default(),
		shouldReconnect: true,
	}
}

// Connect establishes WebSocket connection to Binance
func (c *BinancePerpWSConn) Connect(ctx context.Context, streamName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil // Already connected
	}

	c.streamName = streamName

	// Construct WebSocket URL for raw stream
	url := c.config.BaseWSUrl + "/ws/" + streamName

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		if c.subscription != nil && c.subscription.onError != nil {
			c.subscription.onError(err)
		}
		return err
	}

	c.conn = conn
	c.connected = true
	c.reconnectCount = 0

	// Start goroutines for message handling and ping/pong
	go c.readLoop()
	go c.pingLoop()
	go c.reconnectLoop()

	// Call OnConnect callback
	if c.subscription != nil && c.subscription.onConnect != nil {
		c.subscription.onConnect()
	}

	c.logger.Printf("[BinancePerpWS] Connected to %s", url)
	return nil
}

// Disconnect closes the WebSocket connection gracefully
func (c *BinancePerpWSConn) Disconnect() error {
	c.mu.Lock()
	c.shouldReconnect = false
	conn := c.conn
	c.conn = nil
	c.connected = false
	c.mu.Unlock()

	// Cancel context and close channels first to stop all goroutines
	c.cancel()

	// Wait a moment for goroutines to stop
	time.Sleep(10 * time.Millisecond)

	// Now safely send close message and close connection
	if conn != nil {
		// Send close frame before closing
		err := conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			c.logger.Printf("[BinancePerpWS] Error sending close message: %v", err)
		}

		conn.Close()
	}

	// Close done channel after connection is closed
	select {
	case <-c.done:
		// Already closed
	default:
		close(c.done)
	}

	// Call OnClose callback
	if c.subscription != nil && c.subscription.onClose != nil {
		c.subscription.onClose()
	}

	c.logger.Printf("[BinancePerpWS] Disconnected")
	return nil
}

// IsConnected returns the current connection status
func (c *BinancePerpWSConn) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// readLoop continuously reads messages from the WebSocket connection
func (c *BinancePerpWSConn) readLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.done:
			return
		default:
		}

		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			// Check if this is a graceful shutdown
			if c.ctx.Err() != nil {
				return
			}

			c.logger.Printf("[BinancePerpWS] Read error: %v", err)
			if c.subscription != nil && c.subscription.onError != nil {
				c.subscription.onError(err)
			}

			c.handleDisconnect()
			return
		}

		// Call the message handler if set
		if c.subscription != nil && c.subscription.onMessage != nil {
			c.subscription.onMessage(message)
		}
	}
}

// pingLoop handles ping/pong frames according to Binance requirements
func (c *BinancePerpWSConn) pingLoop() {
	ticker := time.NewTicker(c.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.done:
			return
		case <-ticker.C:
			c.mu.RLock()
			conn := c.conn
			connected := c.connected
			c.mu.RUnlock()

			if !connected || conn == nil {
				continue
			}

			// Send pong frame (unsolicited pong frames are allowed)
			// Use a write timeout to prevent blocking
			err := conn.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				c.logger.Printf("[BinancePerpWS] Pong error: %v", err)
				// Don't call error callback for pong errors during shutdown
				if c.ctx.Err() == nil && c.subscription != nil && c.subscription.onError != nil {
					c.subscription.onError(err)
				}
			}
		}
	}
}

// reconnectLoop handles automatic reconnection
func (c *BinancePerpWSConn) reconnectLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.done:
			return
		case <-c.reconnect:
			if !c.shouldReconnect {
				continue
			}

			// Check if we've exceeded max reconnects
			if c.config.MaxReconnects > 0 && c.reconnectCount >= c.config.MaxReconnects {
				c.logger.Printf("[BinancePerpWS] Max reconnects (%d) exceeded", c.config.MaxReconnects)
				if c.subscription != nil && c.subscription.onError != nil {
					c.subscription.onError(fmt.Errorf("max reconnects exceeded"))
				}
				return
			}

			c.logger.Printf("[BinancePerpWS] Reconnecting in %v... (attempt %d)",
				c.config.ReconnectDelay, c.reconnectCount+1)

			time.Sleep(c.config.ReconnectDelay)

			if err := c.Connect(c.ctx, c.streamName); err != nil {
				c.logger.Printf("[BinancePerpWS] Reconnect failed: %v", err)
				c.reconnectCount++
				// Trigger another reconnect attempt
				select {
				case c.reconnect <- struct{}{}:
				default:
				}
			} else {
				c.logger.Printf("[BinancePerpWS] Reconnected successfully")
				if c.subscription != nil && c.subscription.onReconnect != nil {
					c.subscription.onReconnect()
				}
			}
		}
	}
}

// handleDisconnect handles connection loss and triggers reconnection
func (c *BinancePerpWSConn) handleDisconnect() {
	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.connected = false
	shouldReconnect := c.shouldReconnect && c.ctx.Err() == nil
	c.mu.Unlock()

	if shouldReconnect {
		// Trigger reconnection
		select {
		case c.reconnect <- struct{}{}:
		default:
		}
	}
}
