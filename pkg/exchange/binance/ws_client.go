package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
)

// WSStreamClient represents a high-level WebSocket stream client
type WSStreamClient struct {
	config     *Config
	clients    map[string]*WSClient
	mu         sync.RWMutex
	callbacks  map[string]webSocketCallback
	restClient *Client // For creating listen keys
}

// NewWSStreamClient creates a new WebSocket stream client
func NewWSStreamClient(config *Config) *WSStreamClient {
	return &WSStreamClient{
		config:     config,
		clients:    make(map[string]*WSClient),
		callbacks:  make(map[string]webSocketCallback),
		restClient: NewClient(config),
	}
}

// SubscribeToKline subscribes to kline/candlestick data with subscription options
func (c *WSStreamClient) SubscribeToKline(symbol string, interval string, options *KlineSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s_%s", strings.ToLower(symbol), WSStreamKline, interval)

	// Create a wrapper callback that handles the subscription options
	wsCallback := func(data []byte) error {
		klineData, err := ParseKlineData(data)
		if err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse kline data: %w", err))
			}
			return fmt.Errorf("failed to parse kline data: %w", err)
		}

		if options.onKline != nil {
			if err := options.onKline(klineData); err != nil {
				return err
			}
		}
		return nil
	}

	// Create WebSocket client options
	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// SubscribeToTicker subscribes to 24hr ticker data with subscription options
func (c *WSStreamClient) SubscribeToTicker(symbol string, options *TickerSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamTicker)

	wsCallback := func(data []byte) error {
		tickerData, err := ParseTickerData(data)
		if err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse ticker data: %w", err))
			}
			return fmt.Errorf("failed to parse ticker data: %w", err)
		}

		if options.onTicker != nil {
			if err := options.onTicker(tickerData); err != nil {
				return err
			}
		}
		return nil
	}

	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// SubscribeToMiniTicker subscribes to mini ticker data with subscription options
func (c *WSStreamClient) SubscribeToMiniTicker(symbol string, options *MiniTickerSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamMiniTicker)

	wsCallback := func(data []byte) error {
		miniTickerData, err := ParseMiniTickerData(data)
		if err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse mini ticker data: %w", err))
			}
			return fmt.Errorf("failed to parse mini ticker data: %w", err)
		}

		if options.onMiniTicker != nil {
			return options.onMiniTicker(miniTickerData)
		}
		return nil
	}

	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// SubscribeToAllMiniTickers subscribes to all mini tickers with subscription options
func (c *WSStreamClient) SubscribeToAllMiniTickers(options *MiniTickerSubscriptionOptions) (func() error, error) {
	streamName := "!miniTicker@arr"

	wsCallback := func(data []byte) error {
		var miniTickers []*WSMiniTickerData
		err := json.Unmarshal(data, &miniTickers)
		if err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse mini tickers array: %w", err))
			}
			return fmt.Errorf("failed to parse mini tickers array: %w", err)
		}

		if options.onMiniTicker != nil {
			// For array data, we'll call the callback for each item
			for _, ticker := range miniTickers {
				if err := options.onMiniTicker(ticker); err != nil {
					return err
				}
			}
		}
		return nil
	}

	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// SubscribeToBookTicker subscribes to book ticker data with subscription options
func (c *WSStreamClient) SubscribeToBookTicker(symbol string, options *BookTickerSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamBookTicker)

	wsCallback := func(data []byte) error {
		bookTickerData, err := ParseBookTickerData(data)
		if err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse book ticker data: %w", err))
			}
			return fmt.Errorf("failed to parse book ticker data: %w", err)
		}

		if options.onBookTicker != nil {
			return options.onBookTicker(bookTickerData)
		}
		return nil
	}

	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// SubscribeToAllBookTickers subscribes to all book tickers with subscription options
func (c *WSStreamClient) SubscribeToAllBookTickers(options *BookTickerSubscriptionOptions) (func() error, error) {
	streamName := "!bookTicker"

	wsCallback := func(data []byte) error {
		bookTickerData, err := ParseBookTickerData(data)
		if err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse book ticker data: %w", err))
			}
			return fmt.Errorf("failed to parse book ticker data: %w", err)
		}

		if options.onBookTicker != nil {
			return options.onBookTicker(bookTickerData)
		}
		return nil
	}

	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// SubscribeToDepth subscribes to order book depth data with subscription options
func (c *WSStreamClient) SubscribeToDepth(symbol string, levels string, options *DepthSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s%s", strings.ToLower(symbol), WSStreamDepth, levels)

	wsCallback := func(data []byte) error {
		depthData, err := ParseDepthData(data)
		if err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse depth data: %w", err))
			}
			return fmt.Errorf("failed to parse depth data: %w", err)
		}

		if options.onDepth != nil {
			return options.onDepth(depthData)
		}
		return nil
	}

	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// SubscribeToTrade subscribes to trade data with subscription options
func (c *WSStreamClient) SubscribeToTrade(symbol string, options *TradeSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamTrade)

	wsCallback := func(data []byte) error {
		tradeData, err := ParseTradeData(data)
		if err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse trade data: %w", err))
			}
			return fmt.Errorf("failed to parse trade data: %w", err)
		}

		if options.onTrade != nil {
			return options.onTrade(tradeData)
		}
		return nil
	}

	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// SubscribeToAggTrade subscribes to aggregated trade data with subscription options
func (c *WSStreamClient) SubscribeToAggTrade(symbol string, options *AggTradeSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamAggTrade)

	wsCallback := func(data []byte) error {
		aggTradeData, err := ParseAggTradeData(data)
		if err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse aggregated trade data: %w", err))
			}
			return fmt.Errorf("failed to parse aggregated trade data: %w", err)
		}

		if options.onAggTrade != nil {
			return options.onAggTrade(aggTradeData)
		}
		return nil
	}

	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// SubscribeToUserDataStream subscribes to user data stream with subscription options
func (c *WSStreamClient) SubscribeToUserDataStream(listenKey string, options *UserDataSubscriptionOptions) (func() error, error) {
	streamName := listenKey

	wsCallback := func(data []byte) error {
		// Parse the message to determine the event type
		var message map[string]interface{}
		if err := json.Unmarshal(data, &message); err != nil {
			if options.onError != nil {
				options.onError(fmt.Errorf("failed to parse user data message: %w", err))
			}
			return fmt.Errorf("failed to parse user data message: %w", err)
		}

		eventType, ok := message["e"].(string)
		if !ok {
			if options.onError != nil {
				options.onError(fmt.Errorf("invalid event type in user data message"))
			}
			return fmt.Errorf("invalid event type in user data message")
		}

		switch eventType {
		case "executionReport":
			if options.onExecutionReport != nil {
				executionReport, err := ParseExecutionReport(data)
				if err != nil {
					if options.onError != nil {
						options.onError(fmt.Errorf("failed to parse execution report: %w", err))
					}
					return fmt.Errorf("failed to parse execution report: %w", err)
				}
				return options.onExecutionReport(executionReport)
			}
		case "outboundAccountPosition":
			if options.onAccountUpdate != nil {
				accountUpdate, err := ParseOutboundAccountPosition(data)
				if err != nil {
					if options.onError != nil {
						options.onError(fmt.Errorf("failed to parse account update: %w", err))
					}
					return fmt.Errorf("failed to parse account update: %w", err)
				}
				return options.onAccountUpdate(accountUpdate)
			}
		case "balanceUpdate":
			if options.onBalanceUpdate != nil {
				balanceUpdate, err := ParseBalanceUpdate(data)
				if err != nil {
					if options.onError != nil {
						options.onError(fmt.Errorf("failed to parse balance update: %w", err))
					}
					return fmt.Errorf("failed to parse balance update: %w", err)
				}
				return options.onBalanceUpdate(balanceUpdate)
			}
		default:
			if options.onError != nil {
				options.onError(fmt.Errorf("unknown event type: %s", eventType))
			}
			return fmt.Errorf("unknown event type: %s", eventType)
		}

		return nil
	}

	wsOptions := []WSClientOption{
		WithOnConnect(options.onConnect),
		WithOnDisconnect(options.onDisconnect),
		WithOnError(options.onError),
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, wsOptions)
}

// handleUserDataStreamReconnect handles reconnection logic for user data stream
func (c *WSStreamClient) handleUserDataStreamReconnect(listenKey *string, reconnectChan chan struct{}) error {
	// Close old stream
	if err := c.restClient.CloseUserDataStream(context.Background(), *listenKey); err != nil {
		log.Printf("Failed to close old user data stream: %v", err)
	}

	// Create new listen key
	userDataStream, err := c.restClient.CreateUserDataStream(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create new user data stream: %w", err)
	}

	*listenKey = userDataStream.ListenKey
	log.Printf("Reconnected user data stream with new listen key: %s...", (*listenKey)[:8])

	// Signal reconnect to update WebSocket connection
	select {
	case reconnectChan <- struct{}{}:
	default:
	}

	return nil
}

// subscribeToStreamWithOptions is a helper method that creates a WebSocket client with the given options
func (c *WSStreamClient) subscribeToStreamWithOptions(streamName string, callback webSocketCallback, wsOptions []WSClientOption) (func() error, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if stream is already active
	if _, exists := c.clients[streamName]; exists {
		return nil, fmt.Errorf("stream %s is already active", streamName)
	}

	// Create WebSocket client with options
	wsClient := NewWSClient(c.config, wsOptions...)

	// Set the message callback
	wsClient.onMessage = func(data []byte) {
		if callback != nil {
			if err := callback(data); err != nil {
				log.Printf("Error in stream callback for %s: %v", streamName, err)
			}
		}
	}

	// Construct the full stream URL and connect directly to it
	streamURL := fmt.Sprintf("%s/%s", wsClient.url, streamName)
	wsClient.url = streamURL

	// Connect to WebSocket
	ctx := context.Background()
	if err := wsClient.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to stream %s: %w", streamName, err)
	}

	// Store the client
	c.clients[streamName] = wsClient
	c.callbacks[streamName] = callback

	// Return unsubscribe function
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
	c.callbacks = make(map[string]webSocketCallback)

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

// ParseOutboundAccountPosition parses outbound account position data from WebSocket message
func ParseOutboundAccountPosition(data []byte) (*WSOutboundAccountPosition, error) {
	var accountPosition WSOutboundAccountPosition
	err := json.Unmarshal(data, &accountPosition)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outbound account position data: %w", err)
	}
	return &accountPosition, nil
}

// ParseBalanceUpdate parses balance update data from WebSocket message
func ParseBalanceUpdate(data []byte) (*WSBalanceUpdate, error) {
	var balanceUpdate WSBalanceUpdate
	err := json.Unmarshal(data, &balanceUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse balance update data: %w", err)
	}
	return &balanceUpdate, nil
}

// ParseExecutionReport parses execution report data from WebSocket message
func ParseExecutionReport(data []byte) (*WSExecutionReport, error) {
	var executionReport WSExecutionReport
	err := json.Unmarshal(data, &executionReport)
	if err != nil {
		return nil, fmt.Errorf("failed to parse execution report data: %w", err)
	}
	return &executionReport, nil
}
