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

// SubscribeToKline subscribes to kline/candlestick data with subscription options
func (c *WSStreamClient) SubscribeToKline(symbol string, interval string, options *KlineSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s_%s", strings.ToLower(symbol), WSStreamKline, interval)

	wsCallback := func(data []byte) error {
		klineData, err := ParseKlineData(data)
		if err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse kline data: %w", err))
			}
			return fmt.Errorf("failed to parse kline data: %w", err)
		}

		if options.klineCallback != nil {
			return options.klineCallback(klineData)
		}
		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
}

// SubscribeToTicker subscribes to 24hr ticker data with subscription options
func (c *WSStreamClient) SubscribeToTicker(symbol string, options *TickerSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamTicker)

	wsCallback := func(data []byte) error {
		tickerData, err := ParseTickerData(data)
		if err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse ticker data: %w", err))
			}
			return fmt.Errorf("failed to parse ticker data: %w", err)
		}

		if options.tickerCallback != nil {
			return options.tickerCallback(tickerData)
		}
		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
}

// SubscribeToMiniTicker subscribes to mini ticker data with subscription options
func (c *WSStreamClient) SubscribeToMiniTicker(symbol string, options *MiniTickerSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamMiniTicker)

	wsCallback := func(data []byte) error {
		miniTickerData, err := ParseMiniTickerData(data)
		if err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse mini ticker data: %w", err))
			}
			return fmt.Errorf("failed to parse mini ticker data: %w", err)
		}

		if options.miniTickerCallback != nil {
			return options.miniTickerCallback(miniTickerData)
		}
		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
}

// SubscribeToAllMiniTickers subscribes to all mini tickers
func (c *WSStreamClient) SubscribeToAllMiniTickers(callback WebSocketCallback) (func() error, error) {
	streamName := "!miniTicker@arr"
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToBookTicker subscribes to book ticker data with subscription options
func (c *WSStreamClient) SubscribeToBookTicker(symbol string, options *BookTickerSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamBookTicker)

	wsCallback := func(data []byte) error {
		bookTickerData, err := ParseBookTickerData(data)
		if err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse book ticker data: %w", err))
			}
			return fmt.Errorf("failed to parse book ticker data: %w", err)
		}

		if options.bookTickerCallback != nil {
			return options.bookTickerCallback(bookTickerData)
		}
		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
}

// SubscribeToAllBookTickers subscribes to all book tickers
func (c *WSStreamClient) SubscribeToAllBookTickers(callback WebSocketCallback) (func() error, error) {
	streamName := "!bookTicker"
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToDepth subscribes to order book depth data with subscription options
func (c *WSStreamClient) SubscribeToDepth(symbol string, levels string, options *DepthSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s%s", strings.ToLower(symbol), WSStreamDepth, levels)

	wsCallback := func(data []byte) error {
		depthData, err := ParseDepthData(data)
		if err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse depth data: %w", err))
			}
			return fmt.Errorf("failed to parse depth data: %w", err)
		}

		if options.depthCallback != nil {
			return options.depthCallback(depthData)
		}
		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
}

// SubscribeToTrade subscribes to trade data with subscription options
func (c *WSStreamClient) SubscribeToTrade(symbol string, options *TradeSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamTrade)

	wsCallback := func(data []byte) error {
		tradeData, err := ParseTradeData(data)
		if err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse trade data: %w", err))
			}
			return fmt.Errorf("failed to parse trade data: %w", err)
		}

		if options.tradeCallback != nil {
			return options.tradeCallback(tradeData)
		}
		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
}

// SubscribeToAggTrade subscribes to aggregated trade data with subscription options
func (c *WSStreamClient) SubscribeToAggTrade(symbol string, options *AggTradeSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamAggTrade)

	wsCallback := func(data []byte) error {
		aggTradeData, err := ParseAggTradeData(data)
		if err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse aggregated trade data: %w", err))
			}
			return fmt.Errorf("failed to parse aggregated trade data: %w", err)
		}

		if options.aggTradeCallback != nil {
			return options.aggTradeCallback(aggTradeData)
		}
		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
}

// SubscribeToMarkPrice subscribes to mark price data with subscription options
func (c *WSStreamClient) SubscribeToMarkPrice(symbol string, options *MarkPriceSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamMarkPrice)

	wsCallback := func(data []byte) error {
		markPriceData, err := ParseMarkPriceData(data)
		if err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse mark price data: %w", err))
			}
			return fmt.Errorf("failed to parse mark price data: %w", err)
		}

		if options.markPriceCallback != nil {
			return options.markPriceCallback(markPriceData)
		}
		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
}

// SubscribeToAllMarkPrices subscribes to all mark prices
func (c *WSStreamClient) SubscribeToAllMarkPrices(callback WebSocketCallback) (func() error, error) {
	streamName := "!markPrice@arr"
	return c.subscribeToStream(streamName, callback)
}

// SubscribeToFundingRate subscribes to funding rate data with subscription options
func (c *WSStreamClient) SubscribeToFundingRate(symbol string, options *FundingRateSubscriptionOptions) (func() error, error) {
	streamName := fmt.Sprintf("%s@%s", strings.ToLower(symbol), WSStreamFundingRate)

	wsCallback := func(data []byte) error {
		fundingRateData, err := ParseFundingRateData(data)
		if err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse funding rate data: %w", err))
			}
			return fmt.Errorf("failed to parse funding rate data: %w", err)
		}

		if options.fundingRateCallback != nil {
			return options.fundingRateCallback(fundingRateData)
		}
		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
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

// SubscribeToUserDataStream subscribes to user data stream with subscription options
func (c *WSStreamClient) SubscribeToUserDataStream(options *UserDataSubscriptionOptions) (func() error, error) {
	// Get listen key from REST API
	userDataStream, err := c.restClient.CreateUserDataStream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create user data stream: %w", err)
	}

	streamName := userDataStream.ListenKey

	wsCallback := func(data []byte) error {
		// Parse the event type first
		var baseEvent struct {
			EventType string `json:"e"`
		}
		if err := json.Unmarshal(data, &baseEvent); err != nil {
			if options.errorCallback != nil {
				options.errorCallback(fmt.Errorf("failed to parse event type: %w", err))
			}
			return fmt.Errorf("failed to parse event type: %w", err)
		}

		switch baseEvent.EventType {
		case "executionReport":
			if options.executionReportCallback != nil {
				executionReport, err := ParseExecutionReport(data)
				if err != nil {
					if options.errorCallback != nil {
						options.errorCallback(fmt.Errorf("failed to parse execution report: %w", err))
					}
					return fmt.Errorf("failed to parse execution report: %w", err)
				}
				return options.executionReportCallback(executionReport)
			}
		case "outboundAccountPosition":
			if options.accountUpdateCallback != nil {
				accountUpdate, err := ParseOutboundAccountPosition(data)
				if err != nil {
					if options.errorCallback != nil {
						options.errorCallback(fmt.Errorf("failed to parse account update: %w", err))
					}
					return fmt.Errorf("failed to parse account update: %w", err)
				}
				return options.accountUpdateCallback(accountUpdate)
			}
		case "balanceUpdate":
			if options.balanceUpdateCallback != nil {
				balanceUpdate, err := ParseBalanceUpdate(data)
				if err != nil {
					if options.errorCallback != nil {
						options.errorCallback(fmt.Errorf("failed to parse balance update: %w", err))
					}
					return fmt.Errorf("failed to parse balance update: %w", err)
				}
				return options.balanceUpdateCallback(balanceUpdate)
			}
		default:
			// Unknown event type, log it but don't error
			log.Printf("Unknown user data event type: %s", baseEvent.EventType)
		}

		return nil
	}

	return c.subscribeToStreamWithOptions(streamName, wsCallback, options.SubscriptionOptions)
}

// SubscribeToUserDataStreamWithListenKey subscribes to user data stream with a provided listen key
func (c *WSStreamClient) SubscribeToUserDataStreamWithListenKey(
	listenKey string,
	callback WebSocketCallback,
) (func() error, error) {
	return c.subscribeToStream(listenKey, callback)
}

// handleUserDataStreamReconnect handles reconnection for user data streams
func (c *WSStreamClient) handleUserDataStreamReconnect(listenKey *string, reconnectChan chan struct{}) error {
	go func() {
		for {
			select {
			case <-reconnectChan:
				// Close the old stream
				if err := c.restClient.CloseUserDataStream(context.Background(), *listenKey); err != nil {
					log.Printf("Failed to close old user data stream: %v", err)
				}

				// Create a new listen key
				userDataStream, err := c.restClient.CreateUserDataStream(context.Background())
				if err != nil {
					log.Printf("Failed to create new user data stream: %v", err)
					continue
				}

				*listenKey = userDataStream.ListenKey
				log.Printf("Reconnected user data stream with new listen key: %s...", (*listenKey)[:8])

				// Resubscribe to the new stream
				// Note: This is a simplified implementation. In a real scenario,
				// you might want to store the callback and resubscribe automatically
			}
		}
	}()

	return nil
}

// subscribeToStreamWithReconnect subscribes to a stream with reconnection support
func (c *WSStreamClient) subscribeToStreamWithReconnect(
	streamName string,
	callback WebSocketCallback,
	reconnectChan chan struct{},
) (func() error, error) {
	// Create a new WebSocket client with callbacks
	client := NewWSClient(c.config,
		WithOnMessage(func(data []byte) {
			if err := callback(data); err != nil {
				log.Printf("error in WebSocket callback: %v", err)
			}
		}),
		WithOnError(func(err error) {
			log.Printf("WebSocket error for stream %s: %v", streamName, err)
			reconnectChan <- struct{}{}
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

	c.mu.Lock()
	c.clients[streamName] = client
	c.callbacks[streamName] = callback
	c.mu.Unlock()

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

// subscribeToStreamWithOptions is a helper method that handles subscription with options
func (c *WSStreamClient) subscribeToStreamWithOptions(streamName string, callback WebSocketCallback, options *SubscriptionOptions) (func() error, error) {
	// Create a wrapper callback that handles the options callbacks
	wrappedCallback := func(data []byte) error {
		if err := callback(data); err != nil {
			if options.errorCallback != nil {
				options.errorCallback(err)
			}
			return err
		}
		return nil
	}

	// Subscribe to the stream
	unsubscribe, err := c.subscribeToStream(streamName, wrappedCallback)
	if err != nil {
		return nil, err
	}

	// If we have a connect callback, call it
	if options.connectCallback != nil {
		options.connectCallback()
	}

	// Return a wrapped unsubscribe function that calls the disconnect callback
	return func() error {
		if options.disconnectCallback != nil {
			options.disconnectCallback()
		}
		return unsubscribe()
	}, nil
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

// Parse functions for user data stream events
func ParseListenKeyExpiredEvent(data []byte) (*WSListenKeyExpiredEvent, error) {
	var event WSListenKeyExpiredEvent
	err := json.Unmarshal(data, &event)
	return &event, err
}

func ParseAccountUpdateEvent(data []byte) (*WSAccountUpdateEvent, error) {
	var event WSAccountUpdateEvent
	err := json.Unmarshal(data, &event)
	return &event, err
}

func ParseMarginCallEvent(data []byte) (*WSMarginCallEvent, error) {
	var event WSMarginCallEvent
	err := json.Unmarshal(data, &event)
	return &event, err
}

func ParseOrderTradeUpdateEvent(data []byte) (*WSOrderTradeUpdateEvent, error) {
	var event WSOrderTradeUpdateEvent
	err := json.Unmarshal(data, &event)
	return &event, err
}

func ParseTradeLiteEvent(data []byte) (*WSTradeLiteEvent, error) {
	var event WSTradeLiteEvent
	err := json.Unmarshal(data, &event)
	return &event, err
}

func ParseAccountConfigUpdateEvent(data []byte) (*WSAccountConfigUpdateEvent, error) {
	var event WSAccountConfigUpdateEvent
	err := json.Unmarshal(data, &event)
	return &event, err
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
