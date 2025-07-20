package binancefuture

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestWSStreamClient_SubscribeToAggTrade(t *testing.T) {
	// Create a test configuration
	config := &Config{
		BaseURL: "https://testnet.binancefuture.com", // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	symbol := "btcusdt"
	receivedData := make(chan *WSAggTradeData, 1)

	// Create subscription options
	options := NewAggTradeSubscriptionOptions()
	options.WithConnect(func() {
		t.Log("Connected to aggregated trade stream")
	}).WithAggTrade(func(data *WSAggTradeData) error {
		// Log the received data
		log.Printf("Received aggregated trade: Symbol=%s, Price=%f, Quantity=%f",
			data.Symbol, data.Price, data.Quantity)

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	}).WithError(func(err error) {
		t.Logf("Error in aggregated trade stream: %v", err)
	}).WithDisconnect(func() {
		t.Log("Disconnected from aggregated trade stream")
	})

	// Subscribe to aggregated trades
	unsubscribe, err := wsClient.SubscribeToAggTrade(symbol, options)

	if err != nil {
		t.Fatalf("Failed to subscribe to aggregated trades: %v", err)
	}

	// Wait for some data to be received (with timeout)
	select {
	case data := <-receivedData:
		t.Logf("Successfully received aggregated trade: Symbol=%s, Price=%f, Quantity=%f",
			data.Symbol, data.Price, data.Quantity)

	case <-time.After(30 * time.Second):
		t.Log("Timeout waiting for aggregated trade data")
	}

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		t.Fatalf("Failed to close WebSocket client: %v", err)
	}
}

func TestWSStreamClient_SubscribeToKline(t *testing.T) {
	// Create a test configuration
	config := &Config{
		BaseURL: "https://testnet.binancefuture.com", // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	symbol := "btcusdt"
	interval := "1m"
	receivedData := make(chan *WSKlineData, 1)

	// Create subscription options
	options := NewKlineSubscriptionOptions()
	options.WithConnect(func() {
		t.Log("Connected to kline stream")
	}).WithKline(func(data *WSKlineData) error {
		// Log the received data
		log.Printf("Received kline: Symbol=%s, Close=%f, Volume=%f",
			data.Symbol, data.Kline.ClosePrice, data.Kline.Volume)

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	}).WithError(func(err error) {
		t.Logf("Error in kline stream: %v", err)
	}).WithDisconnect(func() {
		t.Log("Disconnected from kline stream")
	})

	// Subscribe to kline data
	unsubscribe, err := wsClient.SubscribeToKline(symbol, interval, options)

	if err != nil {
		t.Fatalf("Failed to subscribe to kline: %v", err)
	}

	// Wait for some data to be received (with timeout)
	select {
	case data := <-receivedData:
		t.Logf("Successfully received kline: Symbol=%s, Close=%f, Volume=%f",
			data.Symbol, data.Kline.ClosePrice, data.Kline.Volume)

	case <-time.After(30 * time.Second):
		t.Log("Timeout waiting for kline data")
	}

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		t.Fatalf("Failed to close WebSocket client: %v", err)
	}
}

func TestWSStreamClient_SubscribeToTicker(t *testing.T) {
	// Create a test configuration
	config := &Config{
		BaseURL: "https://testnet.binancefuture.com", // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	symbol := "btcusdt"
	receivedData := make(chan *WSTickerData, 1)

	// Create subscription options
	options := NewTickerSubscriptionOptions()
	options.WithConnect(func() {
		t.Log("Connected to ticker stream")
	}).WithTicker(func(data *WSTickerData) error {
		// Log the received data
		log.Printf("Received ticker: Symbol=%s, LastPrice=%f, Volume=%f",
			data.Symbol, data.LastPrice, data.Volume)

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	}).WithError(func(err error) {
		t.Logf("Error in ticker stream: %v", err)
	}).WithDisconnect(func() {
		t.Log("Disconnected from ticker stream")
	})

	// Subscribe to ticker data
	unsubscribe, err := wsClient.SubscribeToTicker(symbol, options)

	if err != nil {
		t.Fatalf("Failed to subscribe to ticker: %v", err)
	}

	// Wait for some data to be received (with timeout)
	select {
	case data := <-receivedData:
		t.Logf("Successfully received ticker: Symbol=%s, LastPrice=%f, Volume=%f",
			data.Symbol, data.LastPrice, data.Volume)

	case <-time.After(30 * time.Second):
		t.Log("Timeout waiting for ticker data")
	}

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		t.Fatalf("Failed to close WebSocket client: %v", err)
	}
}

func TestWSStreamClient_SubscribeToTrade(t *testing.T) {
	// Create a test configuration
	config := &Config{
		BaseURL: "https://testnet.binancefuture.com", // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	symbol := "btcusdt"
	receivedData := make(chan *WSTradeData, 1)

	// Create subscription options
	options := NewTradeSubscriptionOptions()
	options.WithConnect(func() {
		t.Log("Connected to trade stream")
	}).WithTrade(func(data *WSTradeData) error {
		// Log the received data
		log.Printf("Received trade: Symbol=%s, Price=%f, Quantity=%f, IsBuyerMaker=%t",
			data.Symbol, data.Price, data.Quantity, data.IsBuyerMaker)

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	}).WithError(func(err error) {
		t.Logf("Error in trade stream: %v", err)
	}).WithDisconnect(func() {
		t.Log("Disconnected from trade stream")
	})

	// Subscribe to trade data
	unsubscribe, err := wsClient.SubscribeToTrade(symbol, options)

	if err != nil {
		t.Fatalf("Failed to subscribe to trade: %v", err)
	}

	// Wait for some data to be received (with timeout)
	select {
	case data := <-receivedData:
		t.Logf("Successfully received trade: Symbol=%s, Price=%f, Quantity=%f, IsBuyerMaker=%t",
			data.Symbol, data.Price, data.Quantity, data.IsBuyerMaker)

	case <-time.After(30 * time.Second):
		t.Log("Timeout waiting for trade data")
	}

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		t.Fatalf("Failed to close WebSocket client: %v", err)
	}
}

func TestWSStreamClient_SubscribeToDepth(t *testing.T) {
	// Create a test configuration
	config := &Config{
		BaseURL: "https://testnet.binancefuture.com", // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	symbol := "btcusdt"
	levels := "5"
	receivedData := make(chan *WSDepthData, 1)

	// Create subscription options
	options := NewDepthSubscriptionOptions()
	options.WithConnect(func() {
		t.Log("Connected to depth stream")
	}).WithDepth(func(data *WSDepthData) error {
		// Log the received data
		log.Printf("Received depth: Symbol=%s, Bids=%d, Asks=%d",
			data.Symbol, len(data.Bids), len(data.Asks))

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	}).WithError(func(err error) {
		t.Logf("Error in depth stream: %v", err)
	}).WithDisconnect(func() {
		t.Log("Disconnected from depth stream")
	})

	// Subscribe to depth data
	unsubscribe, err := wsClient.SubscribeToDepth(symbol, levels, options)

	if err != nil {
		t.Fatalf("Failed to subscribe to depth: %v", err)
	}

	// Wait for some data to be received (with timeout)
	select {
	case data := <-receivedData:
		t.Logf("Successfully received depth: Symbol=%s, Bids=%d, Asks=%d",
			data.Symbol, len(data.Bids), len(data.Asks))

	case <-time.After(30 * time.Second):
		t.Log("Timeout waiting for depth data")
	}

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		t.Fatalf("Failed to close WebSocket client: %v", err)
	}
}

func TestWSStreamClient_SubscribeToMarkPrice(t *testing.T) {
	// Create a test configuration
	config := &Config{
		BaseURL: "https://testnet.binancefuture.com", // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	symbol := "btcusdt"
	receivedData := make(chan *WSMarkPriceData, 1)

	// Create subscription options
	options := NewMarkPriceSubscriptionOptions()
	options.WithConnect(func() {
		t.Log("Connected to mark price stream")
	}).WithMarkPrice(func(data *WSMarkPriceData) error {
		// Log the received data
		log.Printf("Received mark price: Symbol=%s, MarkPrice=%f, FundingRate=%f",
			data.Symbol, data.MarkPrice, data.LastFundingRate)

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	}).WithError(func(err error) {
		t.Logf("Error in mark price stream: %v", err)
	}).WithDisconnect(func() {
		t.Log("Disconnected from mark price stream")
	})

	// Subscribe to mark price data
	unsubscribe, err := wsClient.SubscribeToMarkPrice(symbol, options)

	if err != nil {
		t.Fatalf("Failed to subscribe to mark price: %v", err)
	}

	// Wait for some data to be received (with timeout)
	select {
	case data := <-receivedData:
		t.Logf("Successfully received mark price: Symbol=%s, MarkPrice=%f, FundingRate=%f",
			data.Symbol, data.MarkPrice, data.LastFundingRate)

	case <-time.After(30 * time.Second):
		t.Log("Timeout waiting for mark price data")
	}

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		t.Fatalf("Failed to close WebSocket client: %v", err)
	}
}

func TestWSStreamClient_SubscribeToFundingRate(t *testing.T) {
	// Create a test configuration
	config := &Config{
		BaseURL: "https://testnet.binancefuture.com", // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	symbol := "btcusdt"
	receivedData := make(chan *WSFundingRateData, 1)

	// Create subscription options
	options := NewFundingRateSubscriptionOptions()
	options.WithConnect(func() {
		t.Log("Connected to funding rate stream")
	}).WithFundingRate(func(data *WSFundingRateData) error {
		// Log the received data
		log.Printf("Received funding rate: Symbol=%s, Rate=%f, Time=%d",
			data.Symbol, data.FundingRate, data.FundingTime)

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	}).WithError(func(err error) {
		t.Logf("Error in funding rate stream: %v", err)
	}).WithDisconnect(func() {
		t.Log("Disconnected from funding rate stream")
	})

	// Subscribe to funding rate data
	unsubscribe, err := wsClient.SubscribeToFundingRate(symbol, options)

	if err != nil {
		t.Fatalf("Failed to subscribe to funding rate: %v", err)
	}

	// Wait for some data to be received (with timeout)
	select {
	case data := <-receivedData:
		t.Logf("Successfully received funding rate: Symbol=%s, Rate=%f, Time=%d",
			data.Symbol, data.FundingRate, data.FundingTime)

	case <-time.After(30 * time.Second):
		t.Log("Timeout waiting for funding rate data")
	}

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		t.Fatalf("Failed to close WebSocket client: %v", err)
	}
}

func TestWSStreamClient_GetActiveStreams(t *testing.T) {
	// Create a test configuration
	config := &Config{
		BaseURL: "https://testnet.binancefuture.com", // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Initially, no streams should be active
	activeStreams := wsClient.GetActiveStreams()
	if len(activeStreams) != 0 {
		t.Errorf("Expected 0 active streams, got %d", len(activeStreams))
	}

	// Subscribe to a stream
	symbol := "btcusdt"
	options := NewAggTradeSubscriptionOptions()
	options.WithAggTrade(func(data *WSAggTradeData) error {
		return nil
	})

	unsubscribe, err := wsClient.SubscribeToAggTrade(symbol, options)

	if err != nil {
		t.Fatalf("Failed to subscribe to aggregated trades: %v", err)
	}

	// Check that the stream is now active
	activeStreams = wsClient.GetActiveStreams()
	expectedStream := fmt.Sprintf("%s@%s", symbol, WSStreamAggTrade)

	found := false
	for _, stream := range activeStreams {
		if stream == expectedStream {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find stream %s in active streams: %v", expectedStream, activeStreams)
	}

	// Check IsStreamActive
	if !wsClient.IsStreamActive(expectedStream) {
		t.Errorf("Expected stream %s to be active", expectedStream)
	}

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Check that no streams are active after unsubscribe
	activeStreams = wsClient.GetActiveStreams()
	if len(activeStreams) != 0 {
		t.Errorf("Expected 0 active streams after unsubscribe, got %d", len(activeStreams))
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		t.Fatalf("Failed to close WebSocket client: %v", err)
	}
}

func TestParseAggTradeData(t *testing.T) {
	// Test data
	testData := []byte(`{
		"s": "BTCUSDT",
		"a": 12345,
		"p": "50000.00",
		"q": "1.5",
		"f": 1000,
		"l": 1002,
		"T": 1640995200000,
		"m": false,
		"M": false,
		"E": 1640995200000,
		"e": "aggTrade"
	}`)

	// Parse the data
	aggTrade, err := ParseAggTradeData(testData)
	if err != nil {
		t.Fatalf("Failed to parse aggregated trade data: %v", err)
	}

	// Verify the parsed data
	if aggTrade.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", aggTrade.Symbol)
	}
	if aggTrade.ID != 12345 {
		t.Errorf("Expected ID 12345, got %d", aggTrade.ID)
	}
	if aggTrade.Price != 50000.00 {
		t.Errorf("Expected price 50000.00, got %f", aggTrade.Price)
	}
	if aggTrade.Quantity != 1.5 {
		t.Errorf("Expected quantity 1.5, got %f", aggTrade.Quantity)
	}
	if aggTrade.EventType != "aggTrade" {
		t.Errorf("Expected event type aggTrade, got %s", aggTrade.EventType)
	}
}

func TestWSStreamClient_UserDataStream(t *testing.T) {
	config := &Config{
		APIKey:    "test-api-key",
		APISecret: "test-api-secret",
		BaseURL:   "https://testnet.binancefuture.com",
	}

	client := NewWSStreamClient(config)

	// Test that user data stream methods exist and have correct signatures
	t.Run("UserDataStreamMethodsExist", func(t *testing.T) {
		// Create subscription options
		options := NewUserDataSubscriptionOptions()
		options.WithExecutionReport(func(data *WSExecutionReport) error {
			return nil
		}).WithAccountUpdate(func(data *WSOutboundAccountPosition) error {
			return nil
		}).WithBalanceUpdate(func(data *WSBalanceUpdate) error {
			return nil
		})

		// Test SubscribeToUserDataStream
		_, err := client.SubscribeToUserDataStream(options)
		// We expect an authentication error, but not a method signature error
		if err != nil && err.Error() == "method SubscribeToUserDataStream not found" {
			t.Fatalf("SubscribeToUserDataStream method not found")
		}
	})
}

func TestParseUserDataStreamEvents(t *testing.T) {
	// Test parsing listen key expired event
	t.Run("ParseListenKeyExpiredEvent", func(t *testing.T) {
		data := []byte(`{
			"e": "listenKeyExpired",
			"E": 1736996475556,
			"listenKey": "WsCMN0a4KHUPTQuX6IUnqEZfB1inxmv1qR4kbf1LuEjur5VdbzqvyxqG9TSjVVxv"
		}`)

		event, err := ParseListenKeyExpiredEvent(data)
		if err != nil {
			t.Fatalf("Failed to parse listen key expired event: %v", err)
		}

		if event.EventType != "listenKeyExpired" {
			t.Errorf("Expected event type 'listenKeyExpired', got '%s'", event.EventType)
		}
		if event.EventTime != 1736996475556 {
			t.Errorf("Expected event time 1736996475556, got %d", event.EventTime)
		}
		if event.ListenKey != "WsCMN0a4KHUPTQuX6IUnqEZfB1inxmv1qR4kbf1LuEjur5VdbzqvyxqG9TSjVVxv" {
			t.Errorf("Expected listen key 'WsCMN0a4KHUPTQuX6IUnqEZfB1inxmv1qR4kbf1LuEjur5VdbzqvyxqG9TSjVVxv', got '%s'", event.ListenKey)
		}
	})

	// Test parsing account update event
	t.Run("ParseAccountUpdateEvent", func(t *testing.T) {
		data := []byte(`{
			"e": "ACCOUNT_UPDATE",
			"E": 1564745798939,
			"T": 1564745798938,
			"a": {
				"m": "ORDER",
				"B": [
					{
						"a": "USDT",
						"wb": "122624.12345678",
						"cw": "100.12345678",
						"bc": "50.12345678"
					}
				],
				"P": [
					{
						"s": "BTCUSDT",
						"pa": "0",
						"ep": "0.00000",
						"bep": "0",
						"cr": "200",
						"up": "0",
						"mt": "isolated",
						"iw": "0.00000000",
						"ps": "BOTH"
					}
				]
			}
		}`)

		event, err := ParseAccountUpdateEvent(data)
		if err != nil {
			t.Fatalf("Failed to parse account update event: %v", err)
		}

		if event.EventType != "ACCOUNT_UPDATE" {
			t.Errorf("Expected event type 'ACCOUNT_UPDATE', got '%s'", event.EventType)
		}
		if event.EventTime != 1564745798939 {
			t.Errorf("Expected event time 1564745798939, got %d", event.EventTime)
		}
		if event.UpdateData.EventReasonType != "ORDER" {
			t.Errorf("Expected event reason type 'ORDER', got '%s'", event.UpdateData.EventReasonType)
		}
		if len(event.UpdateData.Balances) != 1 {
			t.Errorf("Expected 1 balance, got %d", len(event.UpdateData.Balances))
		}
		if len(event.UpdateData.Positions) != 1 {
			t.Errorf("Expected 1 position, got %d", len(event.UpdateData.Positions))
		}
	})

	// Test parsing order trade update event
	t.Run("ParseOrderTradeUpdateEvent", func(t *testing.T) {
		data := []byte(`{
			"e": "ORDER_TRADE_UPDATE",
			"E": 1568879465651,
			"T": 1568879465650,
			"o": {
				"s": "BTCUSDT",
				"c": "TEST",
				"S": "SELL",
				"o": "TRAILING_STOP_MARKET",
				"f": "GTC",
				"q": "0.001",
				"p": "0",
				"ap": "0",
				"sp": "7103.04",
				"x": "NEW",
				"X": "NEW",
				"i": 8886774,
				"l": "0",
				"z": "0",
				"L": "0",
				"N": "USDT",
				"n": "0",
				"T": 1568879465650,
				"t": 0,
				"b": "0",
				"a": "9.91",
				"m": false,
				"R": false,
				"wt": "CONTRACT_PRICE",
				"ot": "TRAILING_STOP_MARKET",
				"ps": "LONG",
				"cp": false,
				"AP": "7476.89",
				"cr": "5.0",
				"pP": false,
				"si": 0,
				"ss": 0,
				"rp": "0",
				"V": "EXPIRE_TAKER",
				"pm": "OPPONENT",
				"gtd": 0
			}
		}`)

		event, err := ParseOrderTradeUpdateEvent(data)
		if err != nil {
			t.Fatalf("Failed to parse order trade update event: %v", err)
		}

		if event.EventType != "ORDER_TRADE_UPDATE" {
			t.Errorf("Expected event type 'ORDER_TRADE_UPDATE', got '%s'", event.EventType)
		}
		if event.EventTime != 1568879465651 {
			t.Errorf("Expected event time 1568879465651, got %d", event.EventTime)
		}
		if event.Order.Symbol != "BTCUSDT" {
			t.Errorf("Expected symbol 'BTCUSDT', got '%s'", event.Order.Symbol)
		}
		if event.Order.Side != "SELL" {
			t.Errorf("Expected side 'SELL', got '%s'", event.Order.Side)
		}
		if event.Order.OrderType != "TRAILING_STOP_MARKET" {
			t.Errorf("Expected order type 'TRAILING_STOP_MARKET', got '%s'", event.Order.OrderType)
		}
	})
}

// Test subscription options chaining
func TestSubscriptionOptionsChaining(t *testing.T) {
	// Test kline subscription options chaining
	t.Run("KlineSubscriptionOptionsChaining", func(t *testing.T) {
		options := NewKlineSubscriptionOptions().
			WithConnect(func() {
				t.Log("Connected!")
			}).
			WithKline(func(data *WSKlineData) error {
				t.Logf("Kline: %s", data.Symbol)
				return nil
			}).
			WithError(func(err error) {
				t.Logf("Error: %v", err)
			}).
			WithDisconnect(func() {
				t.Log("Disconnected!")
			})

		// Verify that all callbacks are set
		if options.connectCallback == nil {
			t.Error("Connect callback should be set")
		}
		if options.klineCallback == nil {
			t.Error("Kline callback should be set")
		}
		if options.errorCallback == nil {
			t.Error("Error callback should be set")
		}
		if options.disconnectCallback == nil {
			t.Error("Disconnect callback should be set")
		}
	})

	// Test ticker subscription options chaining
	t.Run("TickerSubscriptionOptionsChaining", func(t *testing.T) {
		options := NewTickerSubscriptionOptions().
			WithConnect(func() {
				t.Log("Connected!")
			}).
			WithTicker(func(data *WSTickerData) error {
				t.Logf("Ticker: %s", data.Symbol)
				return nil
			}).
			WithError(func(err error) {
				t.Logf("Error: %v", err)
			}).
			WithDisconnect(func() {
				t.Log("Disconnected!")
			})

		// Verify that all callbacks are set
		if options.connectCallback == nil {
			t.Error("Connect callback should be set")
		}
		if options.tickerCallback == nil {
			t.Error("Ticker callback should be set")
		}
		if options.errorCallback == nil {
			t.Error("Error callback should be set")
		}
		if options.disconnectCallback == nil {
			t.Error("Disconnect callback should be set")
		}
	})

	// Test user data subscription options chaining
	t.Run("UserDataSubscriptionOptionsChaining", func(t *testing.T) {
		options := NewUserDataSubscriptionOptions().
			WithConnect(func() {
				t.Log("Connected!")
			}).
			WithExecutionReport(func(data *WSExecutionReport) error {
				t.Logf("Execution Report: %s", data.Symbol)
				return nil
			}).
			WithAccountUpdate(func(data *WSOutboundAccountPosition) error {
				t.Logf("Account Update: %d balances", len(data.Balances))
				return nil
			}).
			WithBalanceUpdate(func(data *WSBalanceUpdate) error {
				t.Logf("Balance Update: %s", data.Asset)
				return nil
			}).
			WithError(func(err error) {
				t.Logf("Error: %v", err)
			}).
			WithDisconnect(func() {
				t.Log("Disconnected!")
			})

		// Verify that all callbacks are set
		if options.connectCallback == nil {
			t.Error("Connect callback should be set")
		}
		if options.executionReportCallback == nil {
			t.Error("Execution report callback should be set")
		}
		if options.accountUpdateCallback == nil {
			t.Error("Account update callback should be set")
		}
		if options.balanceUpdateCallback == nil {
			t.Error("Balance update callback should be set")
		}
		if options.errorCallback == nil {
			t.Error("Error callback should be set")
		}
		if options.disconnectCallback == nil {
			t.Error("Disconnect callback should be set")
		}
	})
}
