package binanceperp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// WSClient manages WebSocket connections for different Binance perpetual futures streams
type WSClient struct {
	subscriptions map[string]*WSSubscription
	mu            sync.RWMutex
	baseWsURL     string
	config        *WSConfig
}

// NewWSClient creates a new WebSocket client
func NewWSClient(config *WSConfig) *WSClient {
	// Use default config if not provided
	if config == nil {
		config = &WSConfig{
			BaseWSUrl:      MainnetWSBaseUrl,
			ReconnectDelay: reconnectDelay,
			PingInterval:   pingInterval,
			MaxReconnects:  -1,
		}
	}

	// Use default URL if not provided
	if config.BaseWSUrl == "" {
		config.BaseWSUrl = MainnetWSBaseUrl
	}

	return &WSClient{
		subscriptions: make(map[string]*WSSubscription),
		baseWsURL:     config.BaseWSUrl,
		config:        config,
	}
}

// SubscribeKline subscribes to kline/candlestick WebSocket stream
func (c *WSClient) SubscribeKline(symbol string, interval string, options *KlineSubscriptionOptions) (func(), error) {
	// Create stream name for kline subscription
	// Format: <symbol>@kline_<interval>
	streamName := fmt.Sprintf("%s@kline_%s", symbol, interval)
	subscriptionID := fmt.Sprintf("kline_%s_%s", symbol, interval)

	return c.subscribe(subscriptionID, streamName, options)
}

// SubscribeAggTrade subscribes to aggregate trade WebSocket stream
func (c *WSClient) SubscribeAggTrade(symbol string, options *AggTradeSubscriptionOptions) (func(), error) {
	// Create stream name for aggregate trade subscription
	// Format: <symbol>@aggTrade
	streamName := fmt.Sprintf("%s@aggTrade", symbol)
	subscriptionID := fmt.Sprintf("aggTrade_%s", symbol)

	return c.subscribe(subscriptionID, streamName, options)
}

// SubscribeTicker subscribes to 24hr ticker statistics WebSocket stream
func (c *WSClient) SubscribeTicker(symbol string, options *TickerSubscriptionOptions) (func(), error) {
	// Create stream name for ticker subscription
	// Format: <symbol>@ticker
	streamName := fmt.Sprintf("%s@ticker", symbol)
	subscriptionID := fmt.Sprintf("ticker_%s", symbol)

	return c.subscribe(subscriptionID, streamName, options)
}

// SubscribeLiquidation subscribes to liquidation order WebSocket stream
func (c *WSClient) SubscribeLiquidation(symbol string, options *LiquidationSubscriptionOptions) (func(), error) {
	// Create stream name for liquidation subscription
	// Format: <symbol>@forceOrder
	streamName := fmt.Sprintf("%s@forceOrder", symbol)
	subscriptionID := fmt.Sprintf("liquidation_%s", symbol)

	return c.subscribe(subscriptionID, streamName, options)
}

// SubscribeDepth subscribes to partial book depth WebSocket stream
func (c *WSClient) SubscribeDepth(symbol string, level DepthLevel, updateSpeed DepthUpdateSpeed, options *DepthSubscriptionOptions) (func(), error) {
	// Validate depth level
	switch level {
	case DepthLevel5, DepthLevel10, DepthLevel20:
		// Valid levels
	default:
		return nil, fmt.Errorf("invalid depth level: %d, must be 5, 10, or 20", level)
	}

	// Validate update speed
	switch updateSpeed {
	case DepthUpdate100ms, DepthUpdate250ms, DepthUpdate500ms:
		// Valid update speeds
	case "": // Empty string defaults to 250ms
		updateSpeed = DepthUpdate250ms
	default:
		return nil, fmt.Errorf("invalid update speed: %s, must be 100ms, 250ms, or 500ms", updateSpeed)
	}

	// Create stream name for depth subscription
	// Format: <symbol>@depth<levels> OR <symbol>@depth<levels>@<speed>
	var streamName string
	if updateSpeed == DepthUpdate250ms {
		// Default 250ms doesn't include speed in stream name
		streamName = fmt.Sprintf("%s@depth%d", symbol, level)
	} else {
		// Include speed for non-default update speeds
		streamName = fmt.Sprintf("%s@depth%d@%s", symbol, level, updateSpeed)
	}

	subscriptionID := fmt.Sprintf("depth_%s_%d_%s", symbol, level, updateSpeed)

	return c.subscribe(subscriptionID, streamName, options)
}

// SubscribeDiffDepth subscribes to differential book depth WebSocket stream
func (c *WSClient) SubscribeDiffDepth(symbol string, updateSpeed DepthUpdateSpeed, options *DiffDepthSubscriptionOptions) (func(), error) {
	// Validate update speed
	switch updateSpeed {
	case DepthUpdate100ms, DepthUpdate250ms, DepthUpdate500ms:
		// Valid update speeds
	case "": // Empty string defaults to 250ms
		updateSpeed = DepthUpdate250ms
	default:
		return nil, fmt.Errorf("invalid update speed: %s, must be 100ms, 250ms, or 500ms", updateSpeed)
	}

	// Create stream name for differential depth subscription
	// Format: <symbol>@depth OR <symbol>@depth@<speed>
	var streamName string
	if updateSpeed == DepthUpdate250ms {
		// Default 250ms doesn't include speed in stream name
		streamName = fmt.Sprintf("%s@depth", symbol)
	} else {
		// Include speed for non-default update speeds
		streamName = fmt.Sprintf("%s@depth@%s", symbol, updateSpeed)
	}

	subscriptionID := fmt.Sprintf("diffdepth_%s_%s", symbol, updateSpeed)

	return c.subscribe(subscriptionID, streamName, options)
}

// subscribe is the common subscription logic for all stream types
func (c *WSClient) subscribe(subscriptionID, streamName string, options interface{}) (func(), error) {
	c.mu.Lock()
	// Check if already subscribed
	if _, exists := c.subscriptions[subscriptionID]; exists {
		c.mu.Unlock()
		return nil, fmt.Errorf("already subscribed to %s stream", subscriptionID)
	}

	// Create subscription with message handler
	subscription := c.createSubscription(subscriptionID, streamName, options)

	// Store subscription
	c.subscriptions[subscriptionID] = subscription
	c.mu.Unlock()

	// Connect to WebSocket
	ctx := context.Background()
	if err := subscription.conn.Connect(ctx, streamName); err != nil {
		c.mu.Lock()
		delete(c.subscriptions, subscriptionID)
		c.mu.Unlock()
		c.callOnError(options, err)
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Update state - OnConnect will be called by the low-level connection callback
	c.mu.Lock()
	subscription.state = StateConnected
	c.mu.Unlock()

	// Return unsubscribe function
	unsubscribeFunc := func() {
		c.unsubscribe(subscriptionID)
	}

	return unsubscribeFunc, nil
}

// createSubscription creates a new subscription with the appropriate message handler
func (c *WSClient) createSubscription(subscriptionID, streamName string, options interface{}) *WSSubscription {
	// Create the underlying subscription for the BinancePerpWSConn
	lowLevelSubscription := &Subscription{}
	lowLevelSubscription.
		WithConnect(func() {
			c.callOnConnect(options)
		}).
		WithReconnect(func() {
			c.callOnReconnect(options)
		}).
		WithError(func(err error) {
			c.callOnError(options, err)
		}).
		WithMessage(func(data []byte) {
			c.handleMessage(subscriptionID, data)
		}).
		WithClose(func() {
			c.callOnDisconnect(options)
		})

	// Create WebSocket connection
	conn := NewBinancePerpWSConn(c.config, lowLevelSubscription)

	// Create subscription
	subscription := &WSSubscription{
		id:      subscriptionID,
		conn:    conn,
		options: options,
		state:   StateConnecting,
	}

	return subscription
}

// handleMessage processes incoming WebSocket messages based on subscription type
func (c *WSClient) handleMessage(subscriptionID string, data []byte) {
	c.mu.RLock()
	subscription, exists := c.subscriptions[subscriptionID]
	c.mu.RUnlock()

	if !exists {
		log.Printf("[WSClient] Received message for unknown subscription: %s", subscriptionID)
		return
	}

	// Parse as a generic map to handle any JSON structure
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		log.Printf("[WSClient] Failed to parse JSON: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to parse JSON: %w", err))
		return
	}

	// Check if this message has an event type field
	eventTypeRaw, hasEventType := rawData["e"]
	if !hasEventType {
		log.Printf("[WSClient] Message missing event type 'e'")
		return
	}

	eventType, ok := eventTypeRaw.(string)
	if !ok {
		log.Printf("[WSClient] Event type 'e' is not a string: %T %v", eventTypeRaw, eventTypeRaw)
		return
	}

	// Route message based on event type and subscription type
	switch eventType {
	case "kline":
		c.handleKlineMessage(subscription, data)
	case "aggTrade":
		c.handleAggTradeMessage(subscription, data)
	case "24hrTicker":
		c.handleTickerMessage(subscription, data)
	case "forceOrder":
		c.handleLiquidationMessage(subscription, data)
	case "depthUpdate":
		c.handleDepthMessage(subscription, data)
	default:
		log.Printf("[WSClient] Unknown event type: %s for subscription: %s", eventType, subscriptionID)
	}
}

// handleKlineMessage processes incoming kline WebSocket messages
func (c *WSClient) handleKlineMessage(subscription *WSSubscription, data []byte) {
	var event WSKlineEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal kline data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal kline data: %w", err))
		return
	}

	// Call the kline callback
	if klineOptions, ok := subscription.options.(*KlineSubscriptionOptions); ok && klineOptions.onKline != nil {
		klineOptions.onKline(event.KlineData)
	}
}

// handleAggTradeMessage processes incoming aggregate trade WebSocket messages
func (c *WSClient) handleAggTradeMessage(subscription *WSSubscription, data []byte) {
	var event WSAggTradeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal aggregate trade data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal aggregate trade data: %w", err))
		return
	}

	// Call the aggregate trade callback
	if aggTradeOptions, ok := subscription.options.(*AggTradeSubscriptionOptions); ok && aggTradeOptions.onAggTrade != nil {
		aggTradeOptions.onAggTrade(event)
	}
}

// handleTickerMessage processes incoming 24hr ticker WebSocket messages
func (c *WSClient) handleTickerMessage(subscription *WSSubscription, data []byte) {
	var event WSTickerEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal ticker data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal ticker data: %w", err))
		return
	}

	// Call the ticker callback
	if tickerOptions, ok := subscription.options.(*TickerSubscriptionOptions); ok && tickerOptions.onTicker != nil {
		tickerOptions.onTicker(event)
	}
}

// handleLiquidationMessage processes incoming liquidation order WebSocket messages
func (c *WSClient) handleLiquidationMessage(subscription *WSSubscription, data []byte) {
	var event WSLiquidationEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal liquidation data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal liquidation data: %w", err))
		return
	}

	// Call the liquidation callback
	if liquidationOptions, ok := subscription.options.(*LiquidationSubscriptionOptions); ok && liquidationOptions.onLiquidation != nil {
		liquidationOptions.onLiquidation(event)
	}
}

// handleDepthMessage processes incoming depth WebSocket messages (both partial and differential)
func (c *WSClient) handleDepthMessage(subscription *WSSubscription, data []byte) {
	var event WSDepthEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[WSClient] Failed to unmarshal depth data: %v", err)
		c.callOnError(subscription.options, fmt.Errorf("failed to unmarshal depth data: %w", err))
		return
	}

	// Route to appropriate callback based on subscription type
	switch opts := subscription.options.(type) {
	case *DepthSubscriptionOptions:
		// Partial book depth (top N levels)
		if opts.onDepth != nil {
			opts.onDepth(event)
		}
	case *DiffDepthSubscriptionOptions:
		// Differential depth (order book changes)
		if opts.onDiffDepth != nil {
			opts.onDiffDepth(event)
		}
	default:
		log.Printf("[WSClient] Unknown depth subscription type for subscription: %s", subscription.id)
	}
}

// callOnConnect calls the OnConnect callback for any subscription type
func (c *WSClient) callOnConnect(options interface{}) {
	switch opts := options.(type) {
	case *KlineSubscriptionOptions:
		if opts.onConnect != nil {
			opts.onConnect()
		}
	case *AggTradeSubscriptionOptions:
		if opts.onConnect != nil {
			opts.onConnect()
		}
	case *TickerSubscriptionOptions:
		if opts.onConnect != nil {
			opts.onConnect()
		}
	case *LiquidationSubscriptionOptions:
		if opts.onConnect != nil {
			opts.onConnect()
		}
	case *DepthSubscriptionOptions:
		if opts.onConnect != nil {
			opts.onConnect()
		}
	case *DiffDepthSubscriptionOptions:
		if opts.onConnect != nil {
			opts.onConnect()
		}
	}
}

// callOnReconnect calls the OnReconnect callback for any subscription type
func (c *WSClient) callOnReconnect(options interface{}) {
	switch opts := options.(type) {
	case *KlineSubscriptionOptions:
		if opts.onReconnect != nil {
			opts.onReconnect()
		}
	case *AggTradeSubscriptionOptions:
		if opts.onReconnect != nil {
			opts.onReconnect()
		}
	case *TickerSubscriptionOptions:
		if opts.onReconnect != nil {
			opts.onReconnect()
		}
	case *LiquidationSubscriptionOptions:
		if opts.onReconnect != nil {
			opts.onReconnect()
		}
	case *DepthSubscriptionOptions:
		if opts.onReconnect != nil {
			opts.onReconnect()
		}
	case *DiffDepthSubscriptionOptions:
		if opts.onReconnect != nil {
			opts.onReconnect()
		}
	}
}

// callOnError calls the OnError callback for any subscription type
func (c *WSClient) callOnError(options interface{}, err error) {
	switch opts := options.(type) {
	case *KlineSubscriptionOptions:
		if opts.onError != nil {
			opts.onError(err)
		}
	case *AggTradeSubscriptionOptions:
		if opts.onError != nil {
			opts.onError(err)
		}
	case *TickerSubscriptionOptions:
		if opts.onError != nil {
			opts.onError(err)
		}
	case *LiquidationSubscriptionOptions:
		if opts.onError != nil {
			opts.onError(err)
		}
	case *DepthSubscriptionOptions:
		if opts.onError != nil {
			opts.onError(err)
		}
	case *DiffDepthSubscriptionOptions:
		if opts.onError != nil {
			opts.onError(err)
		}
	}
}

// callOnDisconnect calls the OnDisconnect callback for any subscription type
func (c *WSClient) callOnDisconnect(options interface{}) {
	switch opts := options.(type) {
	case *KlineSubscriptionOptions:
		if opts.onDisconnect != nil {
			opts.onDisconnect()
		}
	case *AggTradeSubscriptionOptions:
		if opts.onDisconnect != nil {
			opts.onDisconnect()
		}
	case *TickerSubscriptionOptions:
		if opts.onDisconnect != nil {
			opts.onDisconnect()
		}
	case *LiquidationSubscriptionOptions:
		if opts.onDisconnect != nil {
			opts.onDisconnect()
		}
	case *DepthSubscriptionOptions:
		if opts.onDisconnect != nil {
			opts.onDisconnect()
		}
	case *DiffDepthSubscriptionOptions:
		if opts.onDisconnect != nil {
			opts.onDisconnect()
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

	// OnDisconnect callback is called by the connection's OnClose callback
}

// Close closes all active subscriptions
func (c *WSClient) Close() {
	c.mu.Lock()
	subscriptions := make([]*WSSubscription, 0, len(c.subscriptions))
	for _, sub := range c.subscriptions {
		subscriptions = append(subscriptions, sub)
	}
	c.subscriptions = make(map[string]*WSSubscription)
	c.mu.Unlock()

	// Close all connections
	for _, sub := range subscriptions {
		if sub.conn != nil {
			sub.conn.Disconnect()
		}
	}
}

// GetSubscriptionCount returns the number of active subscriptions
func (c *WSClient) GetSubscriptionCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.subscriptions)
}

// IsSubscribed checks if a specific subscription is active
func (c *WSClient) IsSubscribed(subscriptionID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.subscriptions[subscriptionID]
	return exists
}
