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

// User Data Stream Constants
const (
	listenKeyRefreshInterval = 55 * time.Minute // Refresh every 55 minutes (5 minutes before expiry)
	listenKeyRetryDelay      = 30 * time.Second // Retry delay for listen key operations
)

// BinancePerpUserDataStream manages WebSocket connection to Binance perpetual futures user data streams
type BinancePerpUserDataStream struct {
	client       *Client // REST client for listen key management
	conn         *websocket.Conn
	mu           sync.RWMutex
	done         chan struct{}
	reconnect    chan struct{}
	logger       *log.Logger
	config       *WSConfig
	subscription *Subscription

	// Connection state
	connected       bool
	listenKey       string
	ctx             context.Context
	cancel          context.CancelFunc
	shouldReconnect bool
	reconnectCount  int

	// Listen key management
	refreshTicker *time.Ticker
	refreshDone   chan struct{}
}

// NewBinancePerpUserDataStream creates a new user data stream manager
func NewBinancePerpUserDataStream(client *Client, config *WSConfig, subscription *Subscription) *BinancePerpUserDataStream {
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

	return &BinancePerpUserDataStream{
		client:          client,
		config:          config,
		subscription:    subscription,
		ctx:             ctx,
		cancel:          cancel,
		done:            make(chan struct{}),
		reconnect:       make(chan struct{}),
		refreshDone:     make(chan struct{}),
		logger:          log.Default(),
		shouldReconnect: true,
	}
}

// Connect establishes user data stream connection to Binance
func (u *BinancePerpUserDataStream) Connect(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.connected {
		return nil // Already connected
	}

	// Step 1: Get listen key from REST API
	if err := u.ensureListenKey(ctx); err != nil {
		if u.subscription != nil && u.subscription.onError != nil {
			u.subscription.onError(err)
		}
		return fmt.Errorf("failed to get listen key: %w", err)
	}

	// Step 2: Connect WebSocket with listen key
	if err := u.connectWebSocket(ctx); err != nil {
		if u.subscription != nil && u.subscription.onError != nil {
			u.subscription.onError(err)
		}
		return err
	}

	u.connected = true
	u.reconnectCount = 0

	// Step 3: Start background routines
	go u.readLoop()
	go u.pingLoop()
	go u.reconnectLoop()
	go u.listenKeyRefreshLoop()

	// Call OnConnect callback
	if u.subscription != nil && u.subscription.onConnect != nil {
		u.subscription.onConnect()
	}

	u.logger.Printf("[BinancePerpUserData] Connected with listen key: %s", u.listenKey[:10]+"...")
	return nil
}

// Disconnect closes the user data stream connection gracefully
func (u *BinancePerpUserDataStream) Disconnect() error {
	u.mu.Lock()
	u.shouldReconnect = false
	conn := u.conn
	u.conn = nil
	u.connected = false
	u.mu.Unlock()

	// Cancel context and stop refresh timer
	u.cancel()
	if u.refreshTicker != nil {
		u.refreshTicker.Stop()
	}

	// Signal refresh loop to stop
	select {
	case <-u.refreshDone:
	default:
		close(u.refreshDone)
	}

	// Wait a moment for goroutines to stop
	time.Sleep(10 * time.Millisecond)

	// Close WebSocket connection
	if conn != nil {
		err := conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			u.logger.Printf("[BinancePerpUserData] Error sending close message: %v", err)
		}
		conn.Close()
	}

	// Close done channel
	select {
	case <-u.done:
	default:
		close(u.done)
	}

	// Close listen key via REST API
	if u.listenKey != "" {
		if err := u.closeListenKey(); err != nil {
			u.logger.Printf("[BinancePerpUserData] Error closing listen key: %v", err)
		}
	}

	// Call OnClose callback
	if u.subscription != nil && u.subscription.onClose != nil {
		u.subscription.onClose()
	}

	u.logger.Printf("[BinancePerpUserData] Disconnected")
	return nil
}

// IsConnected returns the current connection status
func (u *BinancePerpUserDataStream) IsConnected() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.connected
}

// ensureListenKey gets or creates a listen key using REST API
func (u *BinancePerpUserDataStream) ensureListenKey(ctx context.Context) error {
	if u.listenKey != "" {
		// Try to refresh existing listen key first
		if err := u.refreshListenKey(ctx); err != nil {
			u.logger.Printf("[BinancePerpUserData] Failed to refresh listen key, creating new one: %v", err)
			// If refresh fails, create a new one
			return u.createListenKey(ctx)
		}
		return nil
	}

	// No existing listen key, create a new one
	return u.createListenKey(ctx)
}

// createListenKey creates a new listen key using REST API
func (u *BinancePerpUserDataStream) createListenKey(ctx context.Context) error {
	resp, err := u.client.StartUserDataStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to start user data stream: %w", err)
	}

	if resp.Code != 0 {
		return fmt.Errorf("failed to start user data stream: code %d, message: %s", resp.Code, resp.Message)
	}

	if resp.Data == nil || resp.Data.ListenKey == "" {
		return fmt.Errorf("empty listen key received from API")
	}

	u.listenKey = resp.Data.ListenKey
	u.logger.Printf("[BinancePerpUserData] Created new listen key: %s", u.listenKey[:10]+"...")
	return nil
}

// refreshListenKey extends the validity of existing listen key
func (u *BinancePerpUserDataStream) refreshListenKey(ctx context.Context) error {
	resp, err := u.client.KeepaliveUserDataStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to keepalive user data stream: %w", err)
	}

	if resp.Code != 0 {
		// Check for specific error -1125 (listen key doesn't exist)
		if resp.Code == -1125 {
			return fmt.Errorf("listen key does not exist, need to recreate")
		}
		return fmt.Errorf("failed to keepalive user data stream: code %d, message: %s", resp.Code, resp.Message)
	}

	if resp.Data != nil && resp.Data.ListenKey != "" {
		u.listenKey = resp.Data.ListenKey
	}

	u.logger.Printf("[BinancePerpUserData] Refreshed listen key: %s", u.listenKey[:10]+"...")
	return nil
}

// closeListenKey closes the listen key using REST API
func (u *BinancePerpUserDataStream) closeListenKey() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := u.client.CloseUserDataStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to close user data stream: %w", err)
	}

	if resp.Code != 0 {
		return fmt.Errorf("failed to close user data stream: code %d, message: %s", resp.Code, resp.Message)
	}

	u.logger.Printf("[BinancePerpUserData] Closed listen key: %s", u.listenKey[:10]+"...")
	u.listenKey = ""
	return nil
}

// connectWebSocket establishes WebSocket connection using listen key
func (u *BinancePerpUserDataStream) connectWebSocket(ctx context.Context) error {
	if u.listenKey == "" {
		return fmt.Errorf("no listen key available")
	}

	// Construct WebSocket URL for user data stream
	url := u.config.BaseWSUrl + "/ws/" + u.listenKey

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to user data stream: %w", err)
	}

	u.conn = conn
	u.logger.Printf("[BinancePerpUserData] WebSocket connected to %s", url)
	return nil
}

// readLoop continuously reads messages from the user data stream
func (u *BinancePerpUserDataStream) readLoop() {
	for {
		select {
		case <-u.ctx.Done():
			return
		case <-u.done:
			return
		default:
		}

		u.mu.RLock()
		conn := u.conn
		u.mu.RUnlock()

		if conn == nil {
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			// Check if this is a graceful shutdown
			if u.ctx.Err() != nil {
				return
			}

			u.logger.Printf("[BinancePerpUserData] Read error: %v", err)
			if u.subscription != nil && u.subscription.onError != nil {
				u.subscription.onError(err)
			}

			u.handleDisconnect()
			return
		}

		// Call the message handler if set
		if u.subscription != nil && u.subscription.onMessage != nil {
			u.subscription.onMessage(message)
		}
	}
}

// pingLoop handles ping/pong frames for user data stream
func (u *BinancePerpUserDataStream) pingLoop() {
	ticker := time.NewTicker(u.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-u.ctx.Done():
			return
		case <-u.done:
			return
		case <-ticker.C:
			u.mu.RLock()
			conn := u.conn
			connected := u.connected
			u.mu.RUnlock()

			if !connected || conn == nil {
				continue
			}

			// Send pong frame (unsolicited pong frames are allowed)
			err := conn.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				u.logger.Printf("[BinancePerpUserData] Pong error: %v", err)
				// Don't call error callback for pong errors during shutdown
				if u.ctx.Err() == nil && u.subscription != nil && u.subscription.onError != nil {
					u.subscription.onError(err)
				}
			}
		}
	}
}

// reconnectLoop handles automatic reconnection for user data stream
func (u *BinancePerpUserDataStream) reconnectLoop() {
	for {
		select {
		case <-u.ctx.Done():
			return
		case <-u.done:
			return
		case <-u.reconnect:
			if !u.shouldReconnect {
				continue
			}

			// Check if we've exceeded max reconnects
			if u.config.MaxReconnects > 0 && u.reconnectCount >= u.config.MaxReconnects {
				u.logger.Printf("[BinancePerpUserData] Max reconnects (%d) exceeded", u.config.MaxReconnects)
				if u.subscription != nil && u.subscription.onError != nil {
					u.subscription.onError(fmt.Errorf("max reconnects exceeded"))
				}
				return
			}

			u.logger.Printf("[BinancePerpUserData] Reconnecting in %v... (attempt %d)",
				u.config.ReconnectDelay, u.reconnectCount+1)

			time.Sleep(u.config.ReconnectDelay)

			if err := u.Connect(u.ctx); err != nil {
				u.logger.Printf("[BinancePerpUserData] Reconnect failed: %v", err)
				u.reconnectCount++
				// Trigger another reconnect attempt
				select {
				case u.reconnect <- struct{}{}:
				default:
				}
			} else {
				u.logger.Printf("[BinancePerpUserData] Reconnected successfully")
				if u.subscription != nil && u.subscription.onReconnect != nil {
					u.subscription.onReconnect()
				}
			}
		}
	}
}

// listenKeyRefreshLoop periodically refreshes the listen key to prevent expiry
func (u *BinancePerpUserDataStream) listenKeyRefreshLoop() {
	u.refreshTicker = time.NewTicker(listenKeyRefreshInterval)
	defer u.refreshTicker.Stop()

	for {
		select {
		case <-u.ctx.Done():
			return
		case <-u.done:
			return
		case <-u.refreshDone:
			return
		case <-u.refreshTicker.C:
			u.logger.Printf("[BinancePerpUserData] Refreshing listen key...")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := u.refreshListenKey(ctx)
			cancel()

			if err != nil {
				u.logger.Printf("[BinancePerpUserData] Listen key refresh failed: %v", err)

				// If refresh fails with -1125 error (listen key doesn't exist), try to recreate
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				if createErr := u.createListenKey(ctx); createErr != nil {
					u.logger.Printf("[BinancePerpUserData] Failed to recreate listen key: %v", createErr)
					if u.subscription != nil && u.subscription.onError != nil {
						u.subscription.onError(fmt.Errorf("listen key refresh and recreation failed: %w", createErr))
					}
					// Trigger reconnection to recover
					select {
					case u.reconnect <- struct{}{}:
					default:
					}
				} else {
					u.logger.Printf("[BinancePerpUserData] Successfully recreated listen key after refresh failure")
					// Trigger reconnection with new listen key
					select {
					case u.reconnect <- struct{}{}:
					default:
					}
				}
			} else {
				u.logger.Printf("[BinancePerpUserData] Listen key refreshed successfully")
			}
		}
	}
}

// handleDisconnect handles connection loss and triggers reconnection for user data stream
func (u *BinancePerpUserDataStream) handleDisconnect() {
	u.mu.Lock()
	if u.conn != nil {
		u.conn.Close()
		u.conn = nil
	}
	u.connected = false
	shouldReconnect := u.shouldReconnect && u.ctx.Err() == nil
	u.mu.Unlock()

	if shouldReconnect {
		// Trigger reconnection
		select {
		case u.reconnect <- struct{}{}:
		default:
		}
	}
}
