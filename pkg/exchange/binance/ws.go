package binance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSConnection represents a low-level WebSocket connection to Binance
type WSConnection struct {
	config *Config
	logger *slog.Logger
	conn   *websocket.Conn
	url    string

	// Connection management
	mu             sync.RWMutex
	isConnected    bool
	isReconnecting bool
	shouldStop     bool

	// Message handling
	messageHandler MessageHandler
	errorHandler   ErrorHandler

	// Channels for communication
	writeChan chan []byte
	closeChan chan struct{}

	// Reconnection settings
	reconnectInterval    time.Duration
	maxReconnectAttempts int
	reconnectAttempts    int

	// Request tracking
	requestID int
	requestMu sync.Mutex
}

// MessageHandler is called when a message is received
type MessageHandler func(message []byte)

// ErrorHandler is called when an error occurs
type ErrorHandler func(err error)

// NewWSConnection creates a new WebSocket connection
func NewWSConnection(config *Config) *WSConnection {
	if config == nil {
		config = DefaultConfig()
	}

	wsURL := WSBaseURL
	if config.UseTestnet {
		wsURL = WSBaseURLTestnet
	}

	logger := slog.Default().With("component", "binance-websocket")

	return &WSConnection{
		config:               config,
		logger:               logger,
		url:                  wsURL + "/ws",
		writeChan:            make(chan []byte, 256),
		closeChan:            make(chan struct{}),
		reconnectInterval:    5 * time.Second,
		maxReconnectAttempts: 10,
	}
}

// SetMessageHandler sets the message handler function
func (ws *WSConnection) SetMessageHandler(handler MessageHandler) {
	ws.messageHandler = handler
}

// SetErrorHandler sets the error handler function
func (ws *WSConnection) SetErrorHandler(handler ErrorHandler) {
	ws.errorHandler = handler
}

// Connect establishes the WebSocket connection
func (ws *WSConnection) Connect(ctx context.Context) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.isConnected {
		return nil // Already connected
	}

	ws.logger.Debug("connecting to websocket", "url", ws.url)

	// Parse URL
	u, err := url.Parse(ws.url)
	if err != nil {
		return fmt.Errorf("failed to parse websocket URL: %w", err)
	}

	// Set up dialer with timeout
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = ws.config.Timeout

	// Connect
	conn, _, err := dialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	ws.conn = conn
	ws.isConnected = true
	ws.shouldStop = false
	ws.reconnectAttempts = 0

	ws.logger.Info("websocket connected successfully")

	// Start message handlers
	go ws.readLoop()
	go ws.writeLoop()

	return nil
}

// Disconnect closes the WebSocket connection
func (ws *WSConnection) Disconnect() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if !ws.isConnected {
		return nil // Already disconnected
	}

	ws.logger.Debug("disconnecting websocket")
	ws.shouldStop = true

	// Close the connection
	if ws.conn != nil {
		err := ws.conn.Close()
		ws.conn = nil
		ws.isConnected = false

		// Signal close
		select {
		case ws.closeChan <- struct{}{}:
		default:
		}

		ws.logger.Info("websocket disconnected")
		return err
	}

	return nil
}

// IsConnected returns whether the connection is active
func (ws *WSConnection) IsConnected() bool {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.isConnected
}

// SendMessage sends a message through the WebSocket connection
func (ws *WSConnection) SendMessage(message []byte) error {
	if !ws.IsConnected() {
		return errors.New("websocket not connected")
	}

	select {
	case ws.writeChan <- message:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("timeout sending message")
	}
}

// Subscribe sends a subscription request
func (ws *WSConnection) Subscribe(streams []string) error {
	if len(streams) == 0 {
		return errors.New("no streams to subscribe")
	}

	request := WSRequest{
		Method: WSMethodSubscribe,
		Params: streams,
		ID:     ws.getNextRequestID(),
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal subscribe request: %w", err)
	}

	ws.logger.Debug("subscribing to streams", "streams", streams, "requestId", request.ID)
	return ws.SendMessage(data)
}

// Unsubscribe sends an unsubscription request
func (ws *WSConnection) Unsubscribe(streams []string) error {
	if len(streams) == 0 {
		return errors.New("no streams to unsubscribe")
	}

	request := WSRequest{
		Method: WSMethodUnsubscribe,
		Params: streams,
		ID:     ws.getNextRequestID(),
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal unsubscribe request: %w", err)
	}

	ws.logger.Debug("unsubscribing from streams", "streams", streams, "requestId", request.ID)
	return ws.SendMessage(data)
}

// readLoop continuously reads messages from the WebSocket connection
func (ws *WSConnection) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			ws.logger.Error("panic in read loop", "panic", r)
		}
	}()

	for {
		if ws.shouldStop {
			break
		}

		if !ws.IsConnected() {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Set read deadline
		ws.conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// Read message
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			ws.logger.Error("failed to read websocket message", "error", err)
			ws.handleConnectionError(err)
			continue
		}

		// Handle message
		if ws.messageHandler != nil {
			go ws.messageHandler(message)
		}
	}
}

// writeLoop continuously writes messages to the WebSocket connection
func (ws *WSConnection) writeLoop() {
	defer func() {
		if r := recover(); r != nil {
			ws.logger.Error("panic in write loop", "panic", r)
		}
	}()

	ticker := time.NewTicker(54 * time.Second) // Ping every 54 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ws.closeChan:
			return

		case message := <-ws.writeChan:
			if !ws.IsConnected() {
				continue
			}

			ws.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := ws.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				ws.logger.Error("failed to write websocket message", "error", err)
				ws.handleConnectionError(err)
			}

		case <-ticker.C:
			if !ws.IsConnected() {
				continue
			}

			// Send ping
			ws.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				ws.logger.Error("failed to send ping", "error", err)
				ws.handleConnectionError(err)
			}
		}
	}
}

// handleConnectionError handles connection errors and triggers reconnection
func (ws *WSConnection) handleConnectionError(err error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.shouldStop || ws.isReconnecting {
		return
	}

	ws.logger.Warn("websocket connection error", "error", err)

	if ws.errorHandler != nil {
		go ws.errorHandler(err)
	}

	// Mark as disconnected
	ws.isConnected = false
	if ws.conn != nil {
		ws.conn.Close()
		ws.conn = nil
	}

	// Start reconnection
	go ws.reconnect()
}

// reconnect attempts to reconnect to the WebSocket
func (ws *WSConnection) reconnect() {
	ws.mu.Lock()
	if ws.isReconnecting || ws.shouldStop {
		ws.mu.Unlock()
		return
	}
	ws.isReconnecting = true
	ws.mu.Unlock()

	defer func() {
		ws.mu.Lock()
		ws.isReconnecting = false
		ws.mu.Unlock()
	}()

	for ws.reconnectAttempts < ws.maxReconnectAttempts && !ws.shouldStop {
		ws.reconnectAttempts++
		ws.logger.Info("attempting to reconnect", "attempt", ws.reconnectAttempts, "maxAttempts", ws.maxReconnectAttempts)

		// Wait before reconnecting
		time.Sleep(ws.reconnectInterval)

		if ws.shouldStop {
			return
		}

		// Try to reconnect
		ctx, cancel := context.WithTimeout(context.Background(), ws.config.Timeout)
		err := ws.Connect(ctx)
		cancel()

		if err == nil {
			ws.logger.Info("websocket reconnected successfully")
			return
		}

		ws.logger.Error("failed to reconnect", "error", err, "attempt", ws.reconnectAttempts)
	}

	ws.logger.Error("max reconnection attempts reached, giving up")
	if ws.errorHandler != nil {
		ws.errorHandler(errors.New("max reconnection attempts reached"))
	}
}

// getNextRequestID generates the next request ID
func (ws *WSConnection) getNextRequestID() int {
	ws.requestMu.Lock()
	defer ws.requestMu.Unlock()
	ws.requestID++
	return ws.requestID
}
