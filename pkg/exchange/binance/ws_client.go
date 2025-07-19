package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

// WSClient represents a high-level WebSocket client for Binance
type WSClient struct {
	config *Config
	conn   *WSConnection
	logger *slog.Logger

	// Subscription management
	subscriptions map[string]bool
	subMu         sync.RWMutex

	// Event handlers
	klineHandlers      []KlineHandler
	tickerHandlers     []TickerHandler
	tradeHandlers      []TradeHandler
	depthHandlers      []DepthHandler
	bookTickerHandlers []BookTickerHandler
	aggTradeHandlers   []AggTradeHandler
	handlerMu          sync.RWMutex

	// Channels for stopping
	stopChan chan struct{}
	stopped  bool
	stopMu   sync.Mutex
}

// Event handler types
type KlineHandler func(event *WSKlineEvent)
type TickerHandler func(event *WSTickerEvent)
type TradeHandler func(event *WSTradeEvent)
type DepthHandler func(event *WSDepthEvent)
type BookTickerHandler func(event *WSBookTickerEvent)
type AggTradeHandler func(event *WSAggTradeEvent)

// NewWSClient creates a new high-level WebSocket client
func NewWSClient(config *Config) *WSClient {
	if config == nil {
		config = DefaultConfig()
	}

	logger := slog.Default().With("component", "binance-ws-client")

	client := &WSClient{
		config:        config,
		conn:          NewWSConnection(config),
		logger:        logger,
		subscriptions: make(map[string]bool),
		stopChan:      make(chan struct{}),
	}

	// Set up message and error handlers
	client.conn.SetMessageHandler(client.handleMessage)
	client.conn.SetErrorHandler(client.handleError)

	return client
}

// Connect establishes the WebSocket connection
func (ws *WSClient) Connect(ctx context.Context) error {
	ws.logger.Debug("connecting websocket client")
	return ws.conn.Connect(ctx)
}

// Disconnect closes the WebSocket connection
func (ws *WSClient) Disconnect() error {
	ws.stopMu.Lock()
	defer ws.stopMu.Unlock()

	if ws.stopped {
		return nil
	}

	ws.logger.Debug("disconnecting websocket client")
	ws.stopped = true

	// Signal stop
	select {
	case ws.stopChan <- struct{}{}:
	default:
	}

	return ws.conn.Disconnect()
}

// IsConnected returns whether the WebSocket is connected
func (ws *WSClient) IsConnected() bool {
	return ws.conn.IsConnected()
}

// SubscribeKline subscribes to kline/candlestick streams
func (ws *WSClient) SubscribeKline(symbols []string, interval string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols provided for kline subscription")
	}

	if !ValidateInterval(interval) {
		return fmt.Errorf("invalid interval: %s", interval)
	}

	streams := make([]string, len(symbols))
	for i, symbol := range symbols {
		streams[i] = BuildKlineStreamName(symbol, interval)
	}

	ws.logger.Debug("subscribing to kline streams", "symbols", symbols, "interval", interval)

	// Track subscriptions
	ws.subMu.Lock()
	for _, stream := range streams {
		ws.subscriptions[stream] = true
	}
	ws.subMu.Unlock()

	return ws.conn.Subscribe(streams)
}

// UnsubscribeKline unsubscribes from kline/candlestick streams
func (ws *WSClient) UnsubscribeKline(symbols []string, interval string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols provided for kline unsubscription")
	}

	streams := make([]string, len(symbols))
	for i, symbol := range symbols {
		streams[i] = BuildKlineStreamName(symbol, interval)
	}

	ws.logger.Debug("unsubscribing from kline streams", "symbols", symbols, "interval", interval)

	// Remove from tracked subscriptions
	ws.subMu.Lock()
	for _, stream := range streams {
		delete(ws.subscriptions, stream)
	}
	ws.subMu.Unlock()

	return ws.conn.Unsubscribe(streams)
}

// SubscribeTicker subscribes to 24hr ticker streams
func (ws *WSClient) SubscribeTicker(symbols []string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols provided for ticker subscription")
	}

	streams := make([]string, len(symbols))
	for i, symbol := range symbols {
		streams[i] = BuildTickerStreamName(symbol)
	}

	ws.logger.Debug("subscribing to ticker streams", "symbols", symbols)

	// Track subscriptions
	ws.subMu.Lock()
	for _, stream := range streams {
		ws.subscriptions[stream] = true
	}
	ws.subMu.Unlock()

	return ws.conn.Subscribe(streams)
}

// SubscribeTrade subscribes to trade streams
func (ws *WSClient) SubscribeTrade(symbols []string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols provided for trade subscription")
	}

	streams := make([]string, len(symbols))
	for i, symbol := range symbols {
		streams[i] = BuildTradeStreamName(symbol)
	}

	ws.logger.Debug("subscribing to trade streams", "symbols", symbols)

	// Track subscriptions
	ws.subMu.Lock()
	for _, stream := range streams {
		ws.subscriptions[stream] = true
	}
	ws.subMu.Unlock()

	return ws.conn.Subscribe(streams)
}

// OnKline adds a kline event handler
func (ws *WSClient) OnKline(handler KlineHandler) {
	ws.handlerMu.Lock()
	defer ws.handlerMu.Unlock()
	ws.klineHandlers = append(ws.klineHandlers, handler)
}

// OnTicker adds a ticker event handler
func (ws *WSClient) OnTicker(handler TickerHandler) {
	ws.handlerMu.Lock()
	defer ws.handlerMu.Unlock()
	ws.tickerHandlers = append(ws.tickerHandlers, handler)
}

// OnTrade adds a trade event handler
func (ws *WSClient) OnTrade(handler TradeHandler) {
	ws.handlerMu.Lock()
	defer ws.handlerMu.Unlock()
	ws.tradeHandlers = append(ws.tradeHandlers, handler)
}

// OnDepth adds a depth event handler
func (ws *WSClient) OnDepth(handler DepthHandler) {
	ws.handlerMu.Lock()
	defer ws.handlerMu.Unlock()
	ws.depthHandlers = append(ws.depthHandlers, handler)
}

// OnBookTicker adds a book ticker event handler
func (ws *WSClient) OnBookTicker(handler BookTickerHandler) {
	ws.handlerMu.Lock()
	defer ws.handlerMu.Unlock()
	ws.bookTickerHandlers = append(ws.bookTickerHandlers, handler)
}

// OnAggTrade adds an aggregate trade event handler
func (ws *WSClient) OnAggTrade(handler AggTradeHandler) {
	ws.handlerMu.Lock()
	defer ws.handlerMu.Unlock()
	ws.aggTradeHandlers = append(ws.aggTradeHandlers, handler)
}

// GetSubscriptions returns the current active subscriptions
func (ws *WSClient) GetSubscriptions() []string {
	ws.subMu.RLock()
	defer ws.subMu.RUnlock()

	subscriptions := make([]string, 0, len(ws.subscriptions))
	for stream := range ws.subscriptions {
		subscriptions = append(subscriptions, stream)
	}
	return subscriptions
}

// handleMessage processes incoming WebSocket messages
func (ws *WSClient) handleMessage(message []byte) {
	// Try to parse as a response first
	var response WSResponse
	if err := json.Unmarshal(message, &response); err == nil && response.ID != 0 {
		ws.handleResponse(&response)
		return
	}

	// Try to parse as a stream message
	var streamMsg WSStreamMessage
	if err := json.Unmarshal(message, &streamMsg); err == nil && streamMsg.Stream != "" {
		ws.handleStreamMessage(&streamMsg)
		return
	}

	// Try to parse as direct event data (for single stream connections)
	ws.handleDirectEvent(message)
}

// handleResponse processes subscription/unsubscription responses
func (ws *WSClient) handleResponse(response *WSResponse) {
	if response.Error != nil {
		ws.logger.Error("websocket response error", "error", response.Error, "requestId", response.ID)
		return
	}

	ws.logger.Debug("websocket response received", "requestId", response.ID)
}

// handleStreamMessage processes stream data messages
func (ws *WSClient) handleStreamMessage(msg *WSStreamMessage) {
	ws.handleDirectEvent(msg.Data)
}

// handleDirectEvent processes event data directly
func (ws *WSClient) handleDirectEvent(data []byte) {
	// Try to determine event type by parsing the event type field
	var eventType struct {
		EventType string `json:"e"`
	}

	if err := json.Unmarshal(data, &eventType); err != nil {
		ws.logger.Warn("failed to parse event type", "error", err)
		return
	}

	switch eventType.EventType {
	case WSStreamKline:
		ws.handleKlineEvent(data)
	case WSStreamTicker:
		ws.handleTickerEvent(data)
	case WSStreamTrade:
		ws.handleTradeEvent(data)
	case WSStreamDepth:
		ws.handleDepthEvent(data)
	case WSStreamBookTicker:
		ws.handleBookTickerEvent(data)
	case WSStreamAggTrade:
		ws.handleAggTradeEvent(data)
	default:
		ws.logger.Debug("unknown event type", "eventType", eventType.EventType)
	}
}

// handleKlineEvent processes kline events
func (ws *WSClient) handleKlineEvent(data []byte) {
	var event WSKlineEvent
	if err := json.Unmarshal(data, &event); err != nil {
		ws.logger.Error("failed to parse kline event", "error", err)
		return
	}

	ws.handlerMu.RLock()
	handlers := make([]KlineHandler, len(ws.klineHandlers))
	copy(handlers, ws.klineHandlers)
	ws.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleTickerEvent processes ticker events
func (ws *WSClient) handleTickerEvent(data []byte) {
	var event WSTickerEvent
	if err := json.Unmarshal(data, &event); err != nil {
		ws.logger.Error("failed to parse ticker event", "error", err)
		return
	}

	ws.handlerMu.RLock()
	handlers := make([]TickerHandler, len(ws.tickerHandlers))
	copy(handlers, ws.tickerHandlers)
	ws.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleTradeEvent processes trade events
func (ws *WSClient) handleTradeEvent(data []byte) {
	var event WSTradeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		ws.logger.Error("failed to parse trade event", "error", err)
		return
	}

	ws.handlerMu.RLock()
	handlers := make([]TradeHandler, len(ws.tradeHandlers))
	copy(handlers, ws.tradeHandlers)
	ws.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleDepthEvent processes depth events
func (ws *WSClient) handleDepthEvent(data []byte) {
	var event WSDepthEvent
	if err := json.Unmarshal(data, &event); err != nil {
		ws.logger.Error("failed to parse depth event", "error", err)
		return
	}

	ws.handlerMu.RLock()
	handlers := make([]DepthHandler, len(ws.depthHandlers))
	copy(handlers, ws.depthHandlers)
	ws.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleBookTickerEvent processes book ticker events
func (ws *WSClient) handleBookTickerEvent(data []byte) {
	var event WSBookTickerEvent
	if err := json.Unmarshal(data, &event); err != nil {
		ws.logger.Error("failed to parse book ticker event", "error", err)
		return
	}

	ws.handlerMu.RLock()
	handlers := make([]BookTickerHandler, len(ws.bookTickerHandlers))
	copy(handlers, ws.bookTickerHandlers)
	ws.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleAggTradeEvent processes aggregate trade events
func (ws *WSClient) handleAggTradeEvent(data []byte) {
	var event WSAggTradeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		ws.logger.Error("failed to parse aggregate trade event", "error", err)
		return
	}

	ws.handlerMu.RLock()
	handlers := make([]AggTradeHandler, len(ws.aggTradeHandlers))
	copy(handlers, ws.aggTradeHandlers)
	ws.handlerMu.RUnlock()

	// Call handlers
	for _, handler := range handlers {
		go handler(&event)
	}
}

// handleError processes WebSocket errors
func (ws *WSClient) handleError(err error) {
	ws.logger.Error("websocket error", "error", err)
}
