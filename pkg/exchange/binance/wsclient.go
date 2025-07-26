package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// WSClient manages WebSocket connections for different Binance streams
type WSClient struct {
	subscriptions map[string]*Subscription
	mu            sync.RWMutex
	baseURL       string
}

// NewWSClient creates a new WebSocket client
func NewWSClient(config WSConfig) *WSClient {
	// Use default URL if not provided
	if config.BaseURL == "" {
		config.BaseURL = MainnetWSBaseUrl
	}

	return &WSClient{
		subscriptions: make(map[string]*Subscription),
		baseURL:       config.BaseURL,
	}
}

// SubscribeKline subscribes to kline/candlestick WebSocket stream
func (c *WSClient) SubscribeKline(symbol string, interval string, options KlineSubscriptionOptions) (func(), error) {
	// Create stream path for kline subscription
	// Format: /<symbol>@kline_<interval>
	streamPath := fmt.Sprintf("/%s@kline_%s", symbol, interval)
	subscriptionID := fmt.Sprintf("kline_%s_%s", symbol, interval)

	return c.subscribe(subscriptionID, streamPath, options)
}

// SubscribeAggTrade subscribes to aggregate trade WebSocket stream
func (c *WSClient) SubscribeAggTrade(symbol string, options AggTradeSubscriptionOptions) (func(), error) {
	// Create stream path for aggregate trade subscription
	// Format: /<symbol>@aggTrade
	streamPath := fmt.Sprintf("/%s@aggTrade", symbol)
	subscriptionID := fmt.Sprintf("aggTrade_%s", symbol)

	return c.subscribe(subscriptionID, streamPath, options)
}

// SubscribeTrade subscribes to raw trade WebSocket stream
func (c *WSClient) SubscribeTrade(symbol string, options TradeSubscriptionOptions) (func(), error) {
	// Create stream path for trade subscription
	// Format: /<symbol>@trade
	streamPath := fmt.Sprintf("/%s@trade", symbol)
	subscriptionID := fmt.Sprintf("trade_%s", symbol)

	return c.subscribe(subscriptionID, streamPath, options)
}

// SubscribeDepth subscribes to partial book depth WebSocket stream
func (c *WSClient) SubscribeDepth(symbol string, levels int, updateSpeed string, options DepthSubscriptionOptions) (func(), error) {
	// Validate levels
	if levels != 5 && levels != 10 && levels != 20 {
		return nil, fmt.Errorf("invalid levels: %d, must be 5, 10, or 20", levels)
	}

	// Create stream path for depth subscription
	// Format: /<symbol>@depth<levels> or /<symbol>@depth<levels>@100ms
	var streamPath, subscriptionID string
	if updateSpeed == "100ms" {
		streamPath = fmt.Sprintf("/%s@depth%d@100ms", symbol, levels)
		subscriptionID = fmt.Sprintf("depth_%s_%d_100ms", symbol, levels)
	} else {
		streamPath = fmt.Sprintf("/%s@depth%d", symbol, levels)
		subscriptionID = fmt.Sprintf("depth_%s_%d", symbol, levels)
	}

	return c.subscribe(subscriptionID, streamPath, options)
}

// SubscribeDepthUpdate subscribes to differential depth WebSocket stream
func (c *WSClient) SubscribeDepthUpdate(symbol string, updateSpeed string, options DepthUpdateSubscriptionOptions) (func(), error) {
	// Create stream path for differential depth subscription
	// Format: /<symbol>@depth or /<symbol>@depth@100ms
	var streamPath, subscriptionID string
	if updateSpeed == "100ms" {
		streamPath = fmt.Sprintf("/%s@depth@100ms", symbol)
		subscriptionID = fmt.Sprintf("depthUpdate_%s_100ms", symbol)
	} else {
		streamPath = fmt.Sprintf("/%s@depth", symbol)
		subscriptionID = fmt.Sprintf("depthUpdate_%s", symbol)
	}

	return c.subscribe(subscriptionID, streamPath, options)
}

// subscribe is the common subscription logic for all stream types
func (c *WSClient) subscribe(subscriptionID, streamPath string, options interface{}) (func(), error) {
	c.mu.Lock()
	// Check if already subscribed
	if _, exists := c.subscriptions[subscriptionID]; exists {
		c.mu.Unlock()
		return nil, fmt.Errorf("already subscribed to %s stream", subscriptionID)
	}

	// Create new WebSocket connection
	conn := NewBinanceWSConn(c.baseURL, streamPath)

	// Create subscription
	subscription := &Subscription{
		id:      subscriptionID,
		conn:    conn,
		options: options,
		state:   StateConnecting,
	}

	// Set up message handler
	conn.OnMessage = func(data []byte) {
		c.handleMessage(subscription, data)
	}

	// Store subscription
	c.subscriptions[subscriptionID] = subscription
	c.mu.Unlock()

	// Connect to WebSocket
	if err := conn.Connect(); err != nil {
		c.mu.Lock()
		delete(c.subscriptions, subscriptionID)
		c.mu.Unlock()
		c.callOnError(options, err)
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Update state and call OnConnect
	c.mu.Lock()
	subscription.state = StateConnected
	c.mu.Unlock()

	c.callOnConnect(options)

	// Return unsubscribe function
	unsubscribeFunc := func() {
		c.unsubscribe(subscriptionID)
	}

	return unsubscribeFunc, nil
}

// handleMessage processes incoming WebSocket messages based on event type or structure
func (c *WSClient) handleMessage(subscription *Subscription, data []byte) {
	// Parse as a generic map to handle any JSON structure
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		log.Printf("[WSClient] Failed to parse JSON: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to parse JSON: %w", err))
		return
	}

	// Check if this message has an event type field
	if eventTypeRaw, hasEventType := rawData["e"]; hasEventType {
		// Handle event-based streams
		eventType, ok := eventTypeRaw.(string)
		if !ok {
			log.Printf("[WSClient] Event type 'e' is not a string: %T %v", eventTypeRaw, eventTypeRaw)
			return
		}

		switch eventType {
		case "kline":
			c.handleKlineMessage(subscription, data)
		case "aggTrade":
			c.handleAggTradeMessage(subscription, data)
		case "trade":
			c.handleTradeMessage(subscription, data)
		case "depthUpdate":
			c.handleDepthUpdateMessage(subscription, data)
		default:
			log.Printf("[WSClient] Unknown event type: %s", eventType)
		}
		return
	}

	// Check if this is a partial depth stream (has lastUpdateId but no event type)
	if _, hasLastUpdateId := rawData["lastUpdateId"]; hasLastUpdateId {
		c.handleDepthMessage(subscription, data)
		return
	}

	// Unknown message format
	log.Printf("[WSClient] Unknown message format: no event type field and no lastUpdateId field")
}

// handleKlineMessage processes incoming kline WebSocket messages
func (c *WSClient) handleKlineMessage(subscription *Subscription, data []byte) {
	var event WSKlineEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal kline data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal kline data: %w", err))
		return
	}

	// Call the kline callback
	if klineOptions, ok := subscription.options.(KlineSubscriptionOptions); ok && klineOptions.OnKline != nil {
		klineOptions.OnKline(event.KlineData)
	}
}

// handleAggTradeMessage processes incoming aggregate trade WebSocket messages
func (c *WSClient) handleAggTradeMessage(subscription *Subscription, data []byte) {
	var event WSAggTradeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal aggregate trade data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal aggregate trade data: %w", err))
		return
	}

	// Call the aggregate trade callback
	if aggTradeOptions, ok := subscription.options.(AggTradeSubscriptionOptions); ok && aggTradeOptions.OnAggTrade != nil {
		aggTradeOptions.OnAggTrade(event)
	}
}

// handleTradeMessage processes incoming raw trade WebSocket messages
func (c *WSClient) handleTradeMessage(subscription *Subscription, data []byte) {
	var event WSTradeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal trade data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal trade data: %w", err))
		return
	}

	// Call the trade callback
	if tradeOptions, ok := subscription.options.(TradeSubscriptionOptions); ok && tradeOptions.OnTrade != nil {
		tradeOptions.OnTrade(event)
	}
}

// handleDepthMessage processes incoming partial book depth WebSocket messages
func (c *WSClient) handleDepthMessage(subscription *Subscription, data []byte) {
	var event WSDepthEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal depth data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal depth data: %w", err))
		return
	}

	// Call the depth callback
	if depthOptions, ok := subscription.options.(DepthSubscriptionOptions); ok && depthOptions.OnDepth != nil {
		depthOptions.OnDepth(event)
	}
}

// handleDepthUpdateMessage processes incoming differential depth WebSocket messages
func (c *WSClient) handleDepthUpdateMessage(subscription *Subscription, data []byte) {
	var event WSDepthUpdateEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal depth update data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal depth update data: %w", err))
		return
	}

	// Call the depth update callback
	if depthUpdateOptions, ok := subscription.options.(DepthUpdateSubscriptionOptions); ok && depthUpdateOptions.OnDepthUpdate != nil {
		depthUpdateOptions.OnDepthUpdate(event)
	}
}

// callOnConnect calls the OnConnect callback for any subscription type
func (c *WSClient) callOnConnect(options interface{}) {
	switch opts := options.(type) {
	case KlineSubscriptionOptions:
		if opts.OnConnect != nil {
			opts.OnConnect()
		}
	case AggTradeSubscriptionOptions:
		if opts.OnConnect != nil {
			opts.OnConnect()
		}
	case TradeSubscriptionOptions:
		if opts.OnConnect != nil {
			opts.OnConnect()
		}
	case DepthSubscriptionOptions:
		if opts.OnConnect != nil {
			opts.OnConnect()
		}
	case DepthUpdateSubscriptionOptions:
		if opts.OnConnect != nil {
			opts.OnConnect()
		}
	}
}

// callOnError calls the OnError callback for any subscription type
func (c *WSClient) callOnError(options interface{}, err error) {
	switch opts := options.(type) {
	case KlineSubscriptionOptions:
		if opts.OnError != nil {
			opts.OnError(err)
		}
	case AggTradeSubscriptionOptions:
		if opts.OnError != nil {
			opts.OnError(err)
		}
	case TradeSubscriptionOptions:
		if opts.OnError != nil {
			opts.OnError(err)
		}
	case DepthSubscriptionOptions:
		if opts.OnError != nil {
			opts.OnError(err)
		}
	case DepthUpdateSubscriptionOptions:
		if opts.OnError != nil {
			opts.OnError(err)
		}
	}
}

// callOnDisconnect calls the OnDisconnect callback for any subscription type
func (c *WSClient) callOnDisconnect(options interface{}) {
	switch opts := options.(type) {
	case KlineSubscriptionOptions:
		if opts.OnDisconnect != nil {
			opts.OnDisconnect()
		}
	case AggTradeSubscriptionOptions:
		if opts.OnDisconnect != nil {
			opts.OnDisconnect()
		}
	case TradeSubscriptionOptions:
		if opts.OnDisconnect != nil {
			opts.OnDisconnect()
		}
	case DepthSubscriptionOptions:
		if opts.OnDisconnect != nil {
			opts.OnDisconnect()
		}
	case DepthUpdateSubscriptionOptions:
		if opts.OnDisconnect != nil {
			opts.OnDisconnect()
		}
	}
}

// unsubscribe removes and disconnects a subscription
func (c *WSClient) unsubscribe(subscriptionID string) {
	c.mu.Lock()
	subscription, exists := c.subscriptions[subscriptionID]
	if !exists {
		c.mu.Unlock()
		return
	}

	delete(c.subscriptions, subscriptionID)
	c.mu.Unlock()

	// Disconnect the WebSocket connection
	if subscription.conn != nil {
		subscription.conn.Disconnect()
	}

	// Call OnDisconnect callback
	c.callOnDisconnect(subscription.options)
}

// Close closes all active subscriptions
func (c *WSClient) Close() {
	c.mu.Lock()
	subscriptions := make([]*Subscription, 0, len(c.subscriptions))
	for _, sub := range c.subscriptions {
		subscriptions = append(subscriptions, sub)
	}
	c.subscriptions = make(map[string]*Subscription)
	c.mu.Unlock()

	// Close all connections
	for _, sub := range subscriptions {
		if sub.conn != nil {
			sub.conn.Disconnect()
		}
		c.callOnDisconnect(sub.options)
	}
}
