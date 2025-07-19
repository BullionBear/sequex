package binancefuture

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// WSStreamClient represents a high-level WebSocket stream client
type WSStreamClient struct {
	config     *Config
	clients    map[string]*WSClient
	mu         sync.RWMutex
	callbacks  map[string]WebSocketCallback
	restClient *Client // For creating listen keys
}

// NewWSStreamClient creates a new WebSocket stream client
func NewWSStreamClient(config *Config) *WSStreamClient {
	return &WSStreamClient{
		config:     config,
		clients:    make(map[string]*WSClient),
		callbacks:  make(map[string]WebSocketCallback),
		restClient: NewClient(config),
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

// SubscribeToMarkPrice subscribes to mark price data
func (c *WSStreamClient) SubscribeToMarkPrice(symbol string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamMarkPrice)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToAllMarkPrices subscribes to all mark prices
func (c *WSStreamClient) SubscribeToAllMarkPrices(callback WebSocketCallback) (func() error, error) {
	streamName := "!markPrice@arr"
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToFundingRate subscribes to funding rate data
func (c *WSStreamClient) SubscribeToFundingRate(symbol string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamFundingRate)
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

// SubscribeToMarkPriceWithCallback subscribes to mark price data with type-specific callback
func (c *WSStreamClient) SubscribeToMarkPriceWithCallback(symbol string, callback MarkPriceCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamMarkPrice)

	wsCallback := func(data []byte) error {
		markPriceData, err := ParseMarkPriceData(data)
		if err != nil {
			return fmt.Errorf("failed to parse mark price data: %w", err)
		}
		return callback(markPriceData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToAllMarkPricesWithCallback subscribes to all mark prices with type-specific callback
func (c *WSStreamClient) SubscribeToAllMarkPricesWithCallback(callback func([]*WSMarkPriceData) error) (func() error, error) {
	streamName := "!markPrice@arr"

	wsCallback := func(data []byte) error {
		var markPrices []*WSMarkPriceData
		err := json.Unmarshal(data, &markPrices)
		if err != nil {
			return fmt.Errorf("failed to parse mark prices array: %w", err)
		}
		return callback(markPrices)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToFundingRateWithCallback subscribes to funding rate data with type-specific callback
func (c *WSStreamClient) SubscribeToFundingRateWithCallback(symbol string, callback FundingRateCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamFundingRate)

	wsCallback := func(data []byte) error {
		fundingRateData, err := ParseFundingRateData(data)
		if err != nil {
			return fmt.Errorf("failed to parse funding rate data: %w", err)
		}
		return callback(fundingRateData)
	}

	return c.subscribeToStream(streamName, wsCallback)
}

// SubscribeToPartialDepth subscribes to partial book depth data
func (c *WSStreamClient) SubscribeToPartialDepth(symbol string, levels string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s%s", strings.ToLower(symbol), WSStreamDepth, levels)
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToDiffDepth subscribes to diff depth data
func (c *WSStreamClient) SubscribeToDiffDepth(symbol string, callback WebSocketCallback) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s@100ms", strings.ToLower(symbol), WSStreamDepth)
	return c.subscribeToStream(streamName, callback)
}

// Note: User data stream methods are not yet implemented in the Binance Futures client.
// These will be added when the corresponding REST API methods are implemented.

// subscribeToStream subscribes to a WebSocket stream
func (c *WSStreamClient) subscribeToStream(streamName string, callback WebSocketCallback) (func() error, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we already have a client for this stream
	client, exists := c.clients[streamName]
	if !exists {
		// Create a new WebSocket client
		client = NewWSClient(c.config,
			WithOnMessage(func(data []byte) {
				// Route the message to the appropriate callback
				if callback != nil {
					if err := callback(data); err != nil {
						log.Printf("error in WebSocket callback: %v", err)
					}
				}
			}),
			WithOnError(func(err error) {
				log.Printf("WebSocket error: %v", err)
			}),
		)

		// Connect to the WebSocket
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := client.Connect(ctx)
		cancel()

		if err != nil {
			return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
		}

		// Subscribe to the stream
		err = client.SubscribeToStream(streamName)
		if err != nil {
			client.Disconnect()
			return nil, fmt.Errorf("failed to subscribe to stream %s: %w", streamName, err)
		}

		c.clients[streamName] = client
		c.callbacks[streamName] = callback
	}

	// Return unsubscribe function
	unsubscribe := func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		client, exists := c.clients[streamName]
		if !exists {
			return nil
		}

		// Unsubscribe from the stream
		err := client.UnsubscribeFromStream(streamName)
		if err != nil {
			return fmt.Errorf("failed to unsubscribe from stream %s: %w", streamName, err)
		}

		// Disconnect the client
		err = client.Disconnect()
		if err != nil {
			return fmt.Errorf("failed to disconnect client for stream %s: %w", streamName, err)
		}

		// Remove from maps
		delete(c.clients, streamName)
		delete(c.callbacks, streamName)

		return nil
	}

	return unsubscribe, nil
}

// unsubscribeFromStream unsubscribes from a WebSocket stream
func (c *WSStreamClient) unsubscribeFromStream(streamName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, exists := c.clients[streamName]
	if !exists {
		return nil
	}

	// Unsubscribe from the stream
	err := client.UnsubscribeFromStream(streamName)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from stream %s: %w", streamName, err)
	}

	// Disconnect the client
	err = client.Disconnect()
	if err != nil {
		return fmt.Errorf("failed to disconnect client for stream %s: %w", streamName, err)
	}

	// Remove from maps
	delete(c.clients, streamName)
	delete(c.callbacks, streamName)

	return nil
}

// UnsubscribeFromAllStreams unsubscribes from all WebSocket streams
func (c *WSStreamClient) UnsubscribeFromAllStreams() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errors []error

	for streamName, client := range c.clients {
		// Unsubscribe from the stream
		err := client.UnsubscribeFromStream(streamName)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to unsubscribe from stream %s: %w", streamName, err))
		}

		// Disconnect the client
		err = client.Disconnect()
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to disconnect client for stream %s: %w", streamName, err))
		}
	}

	// Clear maps
	c.clients = make(map[string]*WSClient)
	c.callbacks = make(map[string]WebSocketCallback)

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while unsubscribing: %v", errors)
	}

	return nil
}

// GetActiveStreams returns a list of active stream names
func (c *WSStreamClient) GetActiveStreams() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	streams := make([]string, 0, len(c.clients))
	for streamName := range c.clients {
		streams = append(streams, streamName)
	}

	return streams
}

// IsStreamActive checks if a stream is currently active
func (c *WSStreamClient) IsStreamActive(streamName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.clients[streamName]
	return exists
}

// Close closes all WebSocket connections
func (c *WSStreamClient) Close() error {
	return c.UnsubscribeFromAllStreams()
}

// Parse functions for different WebSocket data types
func ParseKlineData(data []byte) (*WSKlineData, error) {
	var klineData WSKlineData
	err := json.Unmarshal(data, &klineData)
	return &klineData, err
}

func ParseTickerData(data []byte) (*WSTickerData, error) {
	var tickerData WSTickerData
	err := json.Unmarshal(data, &tickerData)
	return &tickerData, err
}

func ParseMiniTickerData(data []byte) (*WSMiniTickerData, error) {
	var miniTickerData WSMiniTickerData
	err := json.Unmarshal(data, &miniTickerData)
	return &miniTickerData, err
}

func ParseBookTickerData(data []byte) (*WSBookTickerData, error) {
	var bookTickerData WSBookTickerData
	err := json.Unmarshal(data, &bookTickerData)
	return &bookTickerData, err
}

func ParseDepthData(data []byte) (*WSDepthData, error) {
	var depthData WSDepthData
	err := json.Unmarshal(data, &depthData)
	return &depthData, err
}

func ParseTradeData(data []byte) (*WSTradeData, error) {
	var tradeData WSTradeData
	err := json.Unmarshal(data, &tradeData)
	return &tradeData, err
}

func ParseAggTradeData(data []byte) (*WSAggTradeData, error) {
	var aggTradeData WSAggTradeData
	err := json.Unmarshal(data, &aggTradeData)
	return &aggTradeData, err
}

func ParseMarkPriceData(data []byte) (*WSMarkPriceData, error) {
	var markPriceData WSMarkPriceData
	err := json.Unmarshal(data, &markPriceData)
	return &markPriceData, err
}

func ParseFundingRateData(data []byte) (*WSFundingRateData, error) {
	var fundingRateData WSFundingRateData
	err := json.Unmarshal(data, &fundingRateData)
	return &fundingRateData, err
}

func ParseOutboundAccountPosition(data []byte) (*WSOutboundAccountPosition, error) {
	var accountPosition WSOutboundAccountPosition
	err := json.Unmarshal(data, &accountPosition)
	return &accountPosition, err
}

func ParseBalanceUpdate(data []byte) (*WSBalanceUpdate, error) {
	var balanceUpdate WSBalanceUpdate
	err := json.Unmarshal(data, &balanceUpdate)
	return &balanceUpdate, err
}

func ParseExecutionReport(data []byte) (*WSExecutionReport, error) {
	var executionReport WSExecutionReport
	err := json.Unmarshal(data, &executionReport)
	return &executionReport, err
}
