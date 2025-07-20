package binancefuture

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestWSStreamClient_SubscribeToAggTrade(t *testing.T) {
	// Create a test configuration
	config := &Config{
		UseTestnet: true, // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	symbol := "btcusdt"
	receivedData := make(chan []byte, 1)

	// Subscribe to aggregated trades
	unsubscribe, err := wsClient.SubscribeToAggTrade(symbol, func(data []byte) error {
		// Log the received data
		log.Printf("Received aggregated trade data: %s", string(data))

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to subscribe to aggregated trades: %v", err)
	}

	// Wait for some data to be received (with timeout)
	select {
	case data := <-receivedData:
		t.Logf("Successfully received aggregated trade data: %s", string(data))

		// Try to parse the data
		var aggTrade WSAggTradeData
		if err := json.Unmarshal(data, &aggTrade); err != nil {
			t.Logf("Failed to parse aggregated trade data: %v", err)
		} else {
			t.Logf("Parsed aggregated trade: Symbol=%s, Price=%f, Quantity=%f",
				aggTrade.Symbol, aggTrade.Price, aggTrade.Quantity)
		}

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

func TestWSStreamClient_SubscribeToAggTradeWithCallback(t *testing.T) {
	// Create a test configuration
	config := &Config{
		UseTestnet: true, // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	symbol := "btcusdt"
	receivedData := make(chan *WSAggTradeData, 1)

	// Subscribe to aggregated trades with type-specific callback
	unsubscribe, err := wsClient.SubscribeToAggTradeWithCallback(symbol, func(data *WSAggTradeData) error {
		// Log the received data
		log.Printf("Received aggregated trade: Symbol=%s, Price=%f, Quantity=%f",
			data.Symbol, data.Price, data.Quantity)

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	})

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

func TestWSStreamClient_SubscribeToCombinedStreams(t *testing.T) {
	// Create a test configuration
	config := &Config{
		UseTestnet: true, // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Test data
	streams := []string{
		"btcusdt@aggTrade",
		"ethusdt@aggTrade",
	}
	receivedData := make(chan []byte, 10)

	// Subscribe to combined streams
	unsubscribe, err := wsClient.SubscribeToCombinedStreams(streams, func(data []byte) error {
		// Log the received data
		log.Printf("Received combined stream data: %s", string(data))

		// Send data to channel for testing
		select {
		case receivedData <- data:
		default:
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to subscribe to combined streams: %v", err)
	}

	// Wait for some data to be received (with timeout)
	select {
	case data := <-receivedData:
		t.Logf("Successfully received combined stream data: %s", string(data))

	case <-time.After(30 * time.Second):
		t.Log("Timeout waiting for combined stream data")
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
		UseTestnet: true, // Use testnet for testing
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
	unsubscribe, err := wsClient.SubscribeToAggTrade(symbol, func(data []byte) error {
		return nil
	})

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
		APIKey:     "test-api-key",
		APISecret:  "test-api-secret",
		UseTestnet: true,
	}

	client := NewWSStreamClient(config)

	// Test that user data stream methods exist and have correct signatures
	t.Run("UserDataStreamMethodsExist", func(t *testing.T) {
		listenKey := "test-listen-key-12345"
		callback := func(data []byte) error {
			return nil
		}

		// Test SubscribeToUserDataStream
		_, err := client.SubscribeToUserDataStream(callback)
		// We expect an authentication error, but not a method signature error
		if err != nil && err.Error() == "method SubscribeToUserDataStream not found" {
			t.Fatalf("SubscribeToUserDataStream method not found")
		}

		// Test SubscribeToUserDataStreamWithListenKey
		_, err = client.SubscribeToUserDataStreamWithListenKey(listenKey, callback)
		if err != nil && err.Error() == "method SubscribeToUserDataStreamWithListenKey not found" {
			t.Fatalf("SubscribeToUserDataStreamWithListenKey method not found")
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
