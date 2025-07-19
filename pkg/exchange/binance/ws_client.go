package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
)

// WSStreamClient represents a high-level WebSocket stream client
type WSStreamClient struct {
	config    *Config
	clients   map[string]*WSClient
	mu        sync.RWMutex
	callbacks map[string]WebSocketCallback
}

// NewWSStreamClient creates a new WebSocket stream client
func NewWSStreamClient(config *Config) *WSStreamClient {
	return &WSStreamClient{
		config:    config,
		clients:   make(map[string]*WSClient),
		callbacks: make(map[string]WebSocketCallback),
	}
}

// SubscribeToKline subscribes to kline/candlestick data
func (c *WSStreamClient) SubscribeToKline(symbol string, interval string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s_%s", strings.ToLower(symbol), WSStreamKline, interval)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToTicker subscribes to 24hr ticker data
func (c *WSStreamClient) SubscribeToTicker(symbol string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamTicker)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToMiniTicker subscribes to mini ticker data
func (c *WSStreamClient) SubscribeToMiniTicker(symbol string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamMiniTicker)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToAllMiniTickers subscribes to all mini tickers
func (c *WSStreamClient) SubscribeToAllMiniTickers(callback WebSocketCallback) (func() error, error) {
	streamName := "!miniTicker@arr"
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToBookTicker subscribes to book ticker data
func (c *WSStreamClient) SubscribeToBookTicker(symbol string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamBookTicker)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToAllBookTickers subscribes to all book tickers
func (c *WSStreamClient) SubscribeToAllBookTickers(callback WebSocketCallback) (func() error, error) {
	streamName := "!bookTicker"
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToDepth subscribes to order book depth data
func (c *WSStreamClient) SubscribeToDepth(symbol string, levels string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s%s", strings.ToLower(symbol), WSStreamDepth, levels)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToTrade subscribes to trade data
func (c *WSStreamClient) SubscribeToTrade(symbol string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamTrade)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToAggTrade subscribes to aggregated trade data
func (c *WSStreamClient) SubscribeToAggTrade(symbol string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamAggTrade)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToCombinedStreams subscribes to multiple streams at once
func (c *WSStreamClient) SubscribeToCombinedStreams(streams []string, callback WebSocketCallback) (func() error, error) {
	streamName := strings.Join(streams, "/")
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToKlineWithCallback subscribes to kline/candlestick data with type-specific callback
func (c *WSStreamClient) SubscribeToKlineWithCallback(symbol string, interval string, callback KlineCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s_%s", strings.ToLower(symbol), WSStreamKline, interval)

	wsCallback := func(data []byte) error {
		klineData, err := ParseKlineData(data)
		if err != nil {
			return fmt.Errorf("failed to parse kline data: %w", err)
		}
		return callback(klineData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToTickerWithCallback subscribes to 24hr ticker data with type-specific callback
func (c *WSStreamClient) SubscribeToTickerWithCallback(symbol string, callback TickerCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamTicker)

	wsCallback := func(data []byte) error {
		tickerData, err := ParseTickerData(data)
		if err != nil {
			return fmt.Errorf("failed to parse ticker data: %w", err)
		}
		return callback(tickerData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToMiniTickerWithCallback subscribes to mini ticker data with type-specific callback
func (c *WSStreamClient) SubscribeToMiniTickerWithCallback(symbol string, callback MiniTickerCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamMiniTicker)

	wsCallback := func(data []byte) error {
		miniTickerData, err := ParseMiniTickerData(data)
		if err != nil {
			return fmt.Errorf("failed to parse mini ticker data: %w", err)
		}
		return callback(miniTickerData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToAllMiniTickersWithCallback subscribes to all mini tickers with type-specific callback
func (c *WSStreamClient) SubscribeToAllMiniTickersWithCallback(callback func([]*WSMiniTickerData) error) (func() error, error) {
	streamName := "!miniTicker@arr"

	wsCallback := func(data []byte) error {
		var miniTickers []*WSMiniTickerData
		err := json.Unmarshal(data, &miniTickers)
		if err != nil {
			return fmt.Errorf("failed to parse mini tickers array: %w", err)
		}
		return callback(miniTickers)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToBookTickerWithCallback subscribes to book ticker data with type-specific callback
func (c *WSStreamClient) SubscribeToBookTickerWithCallback(symbol string, callback BookTickerCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamBookTicker)

	wsCallback := func(data []byte) error {
		bookTickerData, err := ParseBookTickerData(data)
		if err != nil {
			return fmt.Errorf("failed to parse book ticker data: %w", err)
		}
		return callback(bookTickerData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToAllBookTickersWithCallback subscribes to all book tickers with type-specific callback
func (c *WSStreamClient) SubscribeToAllBookTickersWithCallback(callback BookTickerCallback) (func() error, error) {
	streamName := "!bookTicker"

	wsCallback := func(data []byte) error {
		bookTickerData, err := ParseBookTickerData(data)
		if err != nil {
			return fmt.Errorf("failed to parse book ticker data: %w", err)
		}
		return callback(bookTickerData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToDepthWithCallback subscribes to order book depth data with type-specific callback
func (c *WSStreamClient) SubscribeToDepthWithCallback(symbol string, levels string, callback DepthCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s%s", strings.ToLower(symbol), WSStreamDepth, levels)

	wsCallback := func(data []byte) error {
		depthData, err := ParseDepthData(data)
		if err != nil {
			return fmt.Errorf("failed to parse depth data: %w", err)
		}
		return callback(depthData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToTradeWithCallback subscribes to trade data with type-specific callback
func (c *WSStreamClient) SubscribeToTradeWithCallback(symbol string, callback TradeCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamTrade)

	wsCallback := func(data []byte) error {
		tradeData, err := ParseTradeData(data)
		if err != nil {
			return fmt.Errorf("failed to parse trade data: %w", err)
		}
		return callback(tradeData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToAggTradeWithCallback subscribes to aggregated trade data with type-specific callback
func (c *WSStreamClient) SubscribeToAggTradeWithCallback(symbol string, callback AggTradeCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamAggTrade)

	wsCallback := func(data []byte) error {
		aggTradeData, err := ParseAggTradeData(data)
		if err != nil {
			return fmt.Errorf("failed to parse aggregated trade data: %w", err)
		}
		return callback(aggTradeData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToPartialDepthWithCallback subscribes to partial book depth data with type-specific callback
func (c *WSStreamClient) SubscribeToPartialDepthWithCallback(symbol string, levels string, callback PartialDepthCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s%s", strings.ToLower(symbol), WSStreamDepth, levels)

	wsCallback := func(data []byte) error {
		partialDepthData, err := ParsePartialDepthData(data)
		if err != nil {
			return fmt.Errorf("failed to parse partial depth data: %w", err)
		}
		return callback(partialDepthData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToDiffDepthWithCallback subscribes to diff depth data with type-specific callback
func (c *WSStreamClient) SubscribeToDiffDepthWithCallback(symbol string, pushRate string, callback DiffDepthCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s%s", strings.ToLower(symbol), WSStreamDepth, pushRate)

	wsCallback := func(data []byte) error {
		diffDepthData, err := ParseDiffDepthData(data)
		if err != nil {
			return fmt.Errorf("failed to parse diff depth data: %w", err)
		}
		return callback(diffDepthData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToPartialDepth subscribes to partial book depth data (legacy method)
func (c *WSStreamClient) SubscribeToPartialDepth(symbol string, levels string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s%s", strings.ToLower(symbol), WSStreamDepth, levels)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToDiffDepth subscribes to diff depth data (legacy method)
func (c *WSStreamClient) SubscribeToDiffDepth(symbol string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamDepth)
	return c.subscribeToStream(streamName, callback)
}

// subscribeToStream is the internal method to subscribe to any stream
func (c *WSStreamClient) subscribeToStream(streamName string, callback WebSocketCallback) (func() error, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if already subscribed
	if _, exists := c.clients[streamName]; exists {
		return nil, fmt.Errorf("already subscribed to stream: %s", streamName)
	}

	// Create WebSocket client
	wsClient := NewWSClient(c.config,
		WithOnMessage(func(data []byte) {
			// Call the user's callback
			if callback != nil {
				if err := callback(data); err != nil {
					log.Printf("Error in stream callback for %s: %v", streamName, err)
				}
			}
		}),
		WithOnError(func(err error) {
			log.Printf("WebSocket error for stream %s: %v", streamName, err)
		}),
		WithOnConnect(func() {
			log.Printf("Connected to stream: %s", streamName)
		}),
		WithOnDisconnect(func() {
			log.Printf("Disconnected from stream: %s", streamName)
		}),
	)

	// Store the client and callback
	c.clients[streamName] = wsClient
	c.callbacks[streamName] = callback

	// Connect to the stream
	err := wsClient.SubscribeToStream(streamName)
	if err != nil {
		// Clean up on error
		delete(c.clients, streamName)
		delete(c.callbacks, streamName)
		return nil, fmt.Errorf("failed to subscribe to stream %s: %w", streamName, err)
	}

	// Return unsubscription function
	unsubscribe := func() error {
		return c.unsubscribeFromStream(streamName)
	}

	return unsubscribe, nil
}

// unsubscribeFromStream unsubscribes from a specific stream
func (c *WSStreamClient) unsubscribeFromStream(streamName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, exists := c.clients[streamName]
	if !exists {
		return fmt.Errorf("not subscribed to stream: %s", streamName)
	}

	// Disconnect the client
	err := client.Disconnect()
	if err != nil {
		log.Printf("Error disconnecting from stream %s: %v", streamName, err)
	}

	// Remove from maps
	delete(c.clients, streamName)
	delete(c.callbacks, streamName)

	return err
}

// UnsubscribeFromAllStreams unsubscribes from all active streams
func (c *WSStreamClient) UnsubscribeFromAllStreams() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errors []string
	for streamName, client := range c.clients {
		if err := client.Disconnect(); err != nil {
			errors = append(errors, fmt.Sprintf("stream %s: %v", streamName, err))
		}
	}

	// Clear maps
	c.clients = make(map[string]*WSClient)
	c.callbacks = make(map[string]WebSocketCallback)

	if len(errors) > 0 {
		return fmt.Errorf("errors during unsubscribe: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetActiveStreams returns a list of currently active stream names
func (c *WSStreamClient) GetActiveStreams() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	streams := make([]string, 0, len(c.clients))
	for streamName := range c.clients {
		streams = append(streams, streamName)
	}
	return streams
}

// IsStreamActive checks if a specific stream is currently active
func (c *WSStreamClient) IsStreamActive(streamName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	client, exists := c.clients[streamName]
	return exists && client.IsConnected()
}

// Close closes all WebSocket connections
func (c *WSStreamClient) Close() error {
	return c.UnsubscribeFromAllStreams()
}

// Helper functions for parsing WebSocket data

// ParseKlineData parses kline data from WebSocket message
func ParseKlineData(data []byte) (*WSKlineData, error) {
	var klineData WSKlineData
	err := json.Unmarshal(data, &klineData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kline data: %w", err)
	}
	return &klineData, nil
}

// ParseTickerData parses ticker data from WebSocket message
func ParseTickerData(data []byte) (*WSTickerData, error) {
	var tickerData WSTickerData
	err := json.Unmarshal(data, &tickerData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ticker data: %w", err)
	}
	return &tickerData, nil
}

// ParseMiniTickerData parses mini ticker data from WebSocket message
func ParseMiniTickerData(data []byte) (*WSMiniTickerData, error) {
	var miniTickerData WSMiniTickerData
	err := json.Unmarshal(data, &miniTickerData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mini ticker data: %w", err)
	}
	return &miniTickerData, nil
}

// ParseBookTickerData parses book ticker data from WebSocket message
func ParseBookTickerData(data []byte) (*WSBookTickerData, error) {
	var bookTickerData WSBookTickerData
	err := json.Unmarshal(data, &bookTickerData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse book ticker data: %w", err)
	}
	return &bookTickerData, nil
}

// ParseDepthData parses depth data from WebSocket message
func ParseDepthData(data []byte) (*WSDepthData, error) {
	var depthData WSDepthData
	err := json.Unmarshal(data, &depthData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse depth data: %w", err)
	}
	return &depthData, nil
}

// ParseTradeData parses trade data from WebSocket message
func ParseTradeData(data []byte) (*WSTradeData, error) {
	var tradeData WSTradeData
	err := json.Unmarshal(data, &tradeData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trade data: %w", err)
	}
	return &tradeData, nil
}

// ParseAggTradeData parses aggregated trade data from WebSocket message
func ParseAggTradeData(data []byte) (*WSAggTradeData, error) {
	var aggTradeData WSAggTradeData
	err := json.Unmarshal(data, &aggTradeData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse aggregated trade data: %w", err)
	}
	return &aggTradeData, nil
}

// ParsePartialDepthData parses partial depth data from WebSocket message
func ParsePartialDepthData(data []byte) (*WSPartialDepthData, error) {
	var partialDepthData WSPartialDepthData
	err := json.Unmarshal(data, &partialDepthData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse partial depth data: %w", err)
	}
	return &partialDepthData, nil
}

// ParseDiffDepthData parses diff depth data from WebSocket message
func ParseDiffDepthData(data []byte) (*WSDiffDepthData, error) {
	var diffDepthData WSDiffDepthData
	err := json.Unmarshal(data, &diffDepthData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse diff depth data: %w", err)
	}
	return &diffDepthData, nil
}
