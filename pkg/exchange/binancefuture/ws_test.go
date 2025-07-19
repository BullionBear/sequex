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
