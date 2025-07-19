package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// UserDataStreamClient manages private WebSocket connections for user data streams
type UserDataStreamClient struct {
	config *Config
	client *Client
	conn   *WSConnection
	logger *slog.Logger

	// Listen key management
	listenKey       string
	listenKeyMu     sync.RWMutex
	keepAliveTicker *time.Ticker
	keepAliveStop   chan struct{}

	// Event handlers
	accountUpdateHandlers   []AccountUpdateHandler
	balanceUpdateHandlers   []BalanceUpdateHandler
	executionReportHandlers []ExecutionReportHandler
	listStatusHandlers      []ListStatusHandler
	handlerMu               sync.RWMutex

	// State management
	isConnected bool
	shouldStop  bool
	stateMu     sync.RWMutex
}

// Event handler types for user data stream
type AccountUpdateHandler func(event *WSAccountUpdate)
type BalanceUpdateHandler func(event *WSBalanceUpdate)
type ExecutionReportHandler func(event *WSExecutionReport)
type ListStatusHandler func(event *WSListStatus)

// NewUserDataStreamClient creates a new user data stream client
func NewUserDataStreamClient(config *Config) *UserDataStreamClient {
	if config == nil {
		config = DefaultConfig()
	}

	logger := slog.Default().With("component", "binance-user-data-stream")

	return &UserDataStreamClient{
		config:        config,
		client:        NewClient(config),
		logger:        logger,
		keepAliveStop: make(chan struct{}),
	}
}

// Connect establishes the user data stream connection
func (uds *UserDataStreamClient) Connect(ctx context.Context) error {
	uds.stateMu.Lock()
	defer uds.stateMu.Unlock()

	if uds.isConnected {
		return nil
	}

	uds.logger.Debug("connecting user data stream")

	// Create listen key
	streamResp, err := uds.client.CreateUserDataStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user data stream: %w", err)
	}

	uds.listenKeyMu.Lock()
	uds.listenKey = streamResp.ListenKey
	uds.listenKeyMu.Unlock()

	// Create WebSocket connection with listen key
	wsURL := uds.getWebSocketURL()
	uds.conn = &WSConnection{
		config:               uds.config,
		logger:               uds.logger,
		url:                  wsURL,
		writeChan:            make(chan []byte, 256),
		closeChan:            make(chan struct{}),
		reconnectInterval:    5 * time.Second,
		maxReconnectAttempts: 10,
	}

	// Set up message and error handlers
	uds.conn.SetMessageHandler(uds.handleMessage)
	uds.conn.SetErrorHandler(uds.handleError)

	// Connect WebSocket
	err = uds.conn.Connect(ctx)
	if err != nil {
		// Clean up listen key on connection failure
		uds.client.CloseUserDataStream(context.Background(), uds.listenKey)
		return fmt.Errorf("failed to connect websocket: %w", err)
	}

	uds.isConnected = true
	uds.shouldStop = false

	// Start keep-alive routine
	uds.startKeepAlive()

	uds.logger.Info("user data stream connected successfully")
	return nil
}

// Disconnect closes the user data stream connection
func (uds *UserDataStreamClient) Disconnect() error {
	uds.stateMu.Lock()
	defer uds.stateMu.Unlock()

	if !uds.isConnected {
		return nil
	}

	uds.logger.Debug("disconnecting user data stream")
	uds.shouldStop = true

	// Stop keep-alive
	uds.stopKeepAlive()

	// Close WebSocket connection
	var wsErr error
	if uds.conn != nil {
		wsErr = uds.conn.Disconnect()
	}

	// Close listen key
	uds.listenKeyMu.RLock()
	listenKey := uds.listenKey
	uds.listenKeyMu.RUnlock()

	if listenKey != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := uds.client.CloseUserDataStream(ctx, listenKey); err != nil {
			uds.logger.Warn("failed to close user data stream", "error", err)
		}
	}

	uds.isConnected = false
	uds.logger.Info("user data stream disconnected")

	return wsErr
}

// IsConnected returns whether the connection is active
func (uds *UserDataStreamClient) IsConnected() bool {
	uds.stateMu.RLock()
	defer uds.stateMu.RUnlock()
	return uds.isConnected
}

// OnAccountUpdate adds an account update event handler
func (uds *UserDataStreamClient) OnAccountUpdate(handler AccountUpdateHandler) {
	uds.handlerMu.Lock()
	defer uds.handlerMu.Unlock()
	uds.accountUpdateHandlers = append(uds.accountUpdateHandlers, handler)
}

// OnBalanceUpdate adds a balance update event handler
func (uds *UserDataStreamClient) OnBalanceUpdate(handler BalanceUpdateHandler) {
	uds.handlerMu.Lock()
	defer uds.handlerMu.Unlock()
	uds.balanceUpdateHandlers = append(uds.balanceUpdateHandlers, handler)
}

// OnExecutionReport adds an execution report event handler
func (uds *UserDataStreamClient) OnExecutionReport(handler ExecutionReportHandler) {
	uds.handlerMu.Lock()
	defer uds.handlerMu.Unlock()
	uds.executionReportHandlers = append(uds.executionReportHandlers, handler)
}

// OnListStatus adds a list status event handler
func (uds *UserDataStreamClient) OnListStatus(handler ListStatusHandler) {
	uds.handlerMu.Lock()
	defer uds.handlerMu.Unlock()
	uds.listStatusHandlers = append(uds.listStatusHandlers, handler)
}

// getWebSocketURL constructs the WebSocket URL with listen key
func (uds *UserDataStreamClient) getWebSocketURL() string {
	baseURL := WSBaseURL
	if uds.config.UseTestnet {
		baseURL = WSBaseURLTestnet
	}

	uds.listenKeyMu.RLock()
	listenKey := uds.listenKey
	uds.listenKeyMu.RUnlock()

	return baseURL + "/ws/" + listenKey
}

// startKeepAlive starts the listen key keep-alive routine
func (uds *UserDataStreamClient) startKeepAlive() {
	// Keep alive every 30 minutes (Binance recommends every 30-60 minutes)
	uds.keepAliveTicker = time.NewTicker(30 * time.Minute)

	go func() {
		defer uds.keepAliveTicker.Stop()

		for {
			select {
			case <-uds.keepAliveStop:
				return
			case <-uds.keepAliveTicker.C:
				if uds.shouldStop {
					return
				}

				uds.listenKeyMu.RLock()
				listenKey := uds.listenKey
				uds.listenKeyMu.RUnlock()

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				err := uds.client.KeepAliveUserDataStream(ctx, listenKey)
				cancel()

				if err != nil {
					uds.logger.Error("failed to keep alive user data stream", "error", err)
				} else {
					uds.logger.Debug("user data stream kept alive")
				}
			}
		}
	}()
}

// stopKeepAlive stops the keep-alive routine
func (uds *UserDataStreamClient) stopKeepAlive() {
	if uds.keepAliveTicker != nil {
		select {
		case uds.keepAliveStop <- struct{}{}:
		default:
		}
		uds.keepAliveTicker.Stop()
	}
}

// handleMessage processes incoming WebSocket messages
func (uds *UserDataStreamClient) handleMessage(message []byte) {
	// Try to determine event type
	var eventType struct {
		EventType string `json:"e"`
	}

	if err := json.Unmarshal(message, &eventType); err != nil {
		uds.logger.Warn("failed to parse event type", "error", err)
		return
	}

	switch eventType.EventType {
	case WSEventAccountUpdate:
		uds.handleAccountUpdate(message)
	case WSEventBalanceUpdate:
		uds.handleBalanceUpdate(message)
	case WSEventExecutionReport:
		uds.handleExecutionReport(message)
	case WSEventListStatus:
		uds.handleListStatus(message)
	default:
		uds.logger.Debug("unknown user data stream event type", "eventType", eventType.EventType)
	}
}

// handleAccountUpdate processes account update events
func (uds *UserDataStreamClient) handleAccountUpdate(data []byte) {
	var event WSAccountUpdate
	if err := json.Unmarshal(data, &event); err != nil {
		uds.logger.Error("failed to parse account update event", "error", err)
		return
	}

	uds.handlerMu.RLock()
	handlers := make([]AccountUpdateHandler, len(uds.accountUpdateHandlers))
	copy(handlers, uds.accountUpdateHandlers)
	uds.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleBalanceUpdate processes balance update events
func (uds *UserDataStreamClient) handleBalanceUpdate(data []byte) {
	var event WSBalanceUpdate
	if err := json.Unmarshal(data, &event); err != nil {
		uds.logger.Error("failed to parse balance update event", "error", err)
		return
	}

	uds.handlerMu.RLock()
	handlers := make([]BalanceUpdateHandler, len(uds.balanceUpdateHandlers))
	copy(handlers, uds.balanceUpdateHandlers)
	uds.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleExecutionReport processes execution report events
func (uds *UserDataStreamClient) handleExecutionReport(data []byte) {
	var event WSExecutionReport
	if err := json.Unmarshal(data, &event); err != nil {
		uds.logger.Error("failed to parse execution report event", "error", err)
		return
	}

	uds.handlerMu.RLock()
	handlers := make([]ExecutionReportHandler, len(uds.executionReportHandlers))
	copy(handlers, uds.executionReportHandlers)
	uds.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleListStatus processes list status events
func (uds *UserDataStreamClient) handleListStatus(data []byte) {
	var event WSListStatus
	if err := json.Unmarshal(data, &event); err != nil {
		uds.logger.Error("failed to parse list status event", "error", err)
		return
	}

	uds.handlerMu.RLock()
	handlers := make([]ListStatusHandler, len(uds.listStatusHandlers))
	copy(handlers, uds.listStatusHandlers)
	uds.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleError processes WebSocket errors
func (uds *UserDataStreamClient) handleError(err error) {
	uds.logger.Error("user data stream error", "error", err)
}
