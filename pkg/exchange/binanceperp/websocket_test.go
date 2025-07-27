package binanceperp

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

// AggTradeEvent represents the aggregate trade WebSocket event for perpetual futures
type AggTradeEvent struct {
	EventType    string `json:"e"` // Event type
	EventTime    int64  `json:"E"` // Event time
	Symbol       string `json:"s"` // Symbol
	AggTradeID   int64  `json:"a"` // Aggregate trade ID
	Price        string `json:"p"` // Price
	Quantity     string `json:"q"` // Quantity
	FirstTradeID int64  `json:"f"` // First trade ID
	LastTradeID  int64  `json:"l"` // Last trade ID
	TradeTime    int64  `json:"T"` // Trade time
	IsBuyerMaker bool   `json:"m"` // Is the buyer the market maker?
}

func TestBinancePerpWSConn_AggTradePayload(t *testing.T) {
	// Test configuration for Binance perpetual futures (lowercase symbols)
	symbol := "btcusdt"
	streamName := symbol + "@aggTrade"
	timeout := 10 * time.Second

	// Create message channel for testing
	msgCh := make(chan AggTradeEvent, 1)
	errorCh := make(chan error, 1)

	// Create subscription with test callbacks
	subscription := &Subscription{}
	subscription.
		WithConnect(func() {
			t.Log("WebSocket connected successfully")
		}).
		WithError(func(err error) {
			t.Logf("WebSocket error: %v", err)
			select {
			case errorCh <- err:
			default:
			}
		}).
		WithMessage(func(data []byte) {
			var event AggTradeEvent
			if err := json.Unmarshal(data, &event); err == nil && event.EventType == "aggTrade" {
				t.Logf("Received raw message: %s", string(data))
				select {
				case msgCh <- event:
				default:
				}
			} else if err != nil {
				t.Logf("Failed to unmarshal aggTrade: %v, raw: %s", err, string(data))
			}
		}).
		WithClose(func() {
			t.Log("WebSocket connection closed")
		})

	// Create WebSocket configuration
	config := &WSConfig{
		BaseWSUrl:      MainnetWSBaseUrl,
		ReconnectDelay: 1 * time.Second,  // Faster reconnect for tests
		PingInterval:   30 * time.Second, // Longer ping interval for tests
		MaxReconnects:  3,
	}

	// Create WebSocket connection
	conn := NewBinancePerpWSConn(config, subscription)

	// Connect to the WebSocket
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.Connect(ctx, streamName); err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	// Verify connection state
	if !conn.IsConnected() {
		t.Fatal("Connection should be established")
	}

	t.Logf("Connected to stream: %s", streamName)
	t.Logf("Waiting %v for aggTrade data...", timeout)

	// Wait for message or timeout
	select {
	case event := <-msgCh:
		// Validate the received aggTrade event
		if event.EventType != "aggTrade" {
			t.Fatalf("Expected event type 'aggTrade', got: %s", event.EventType)
		}

		if event.Symbol != "BTCUSDT" && event.Symbol != "btcusdt" {
			t.Fatalf("Expected symbol BTCUSDT/btcusdt, got: %s", event.Symbol)
		}

		if event.Price == "" {
			t.Fatal("Price should not be empty")
		}

		if event.Quantity == "" {
			t.Fatal("Quantity should not be empty")
		}

		if event.EventTime <= 0 {
			t.Fatal("EventTime should be positive")
		}

		if event.TradeTime <= 0 {
			t.Fatal("TradeTime should be positive")
		}

		if event.AggTradeID <= 0 {
			t.Fatal("AggTradeID should be positive")
		}

		t.Logf("✓ Received valid aggTrade event:")
		t.Logf("  Symbol: %s", event.Symbol)
		t.Logf("  Price: %s", event.Price)
		t.Logf("  Quantity: %s", event.Quantity)
		t.Logf("  AggTradeID: %d", event.AggTradeID)
		t.Logf("  IsBuyerMaker: %t", event.IsBuyerMaker)
		t.Logf("  EventTime: %d", event.EventTime)
		t.Logf("  TradeTime: %d", event.TradeTime)

	case err := <-errorCh:
		t.Fatalf("Received unexpected error: %v", err)

	case <-time.After(timeout):
		t.Fatalf("Timeout: did not receive aggTrade event within %v", timeout)
	}

	// Test disconnection
	if err := conn.Disconnect(); err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}

	// Verify connection state after disconnect
	time.Sleep(100 * time.Millisecond) // Give time for graceful shutdown
	if conn.IsConnected() {
		t.Fatal("Connection should be closed after disconnect")
	}
}

func TestBinancePerpWSConn_ConnectionLifecycle(t *testing.T) {
	// Test connection lifecycle management
	connectCount := 0
	disconnectCount := 0
	errorCount := 0

	subscription := &Subscription{}
	subscription.
		WithConnect(func() {
			connectCount++
			t.Log("OnConnect called")
		}).
		WithClose(func() {
			disconnectCount++
			t.Log("OnClose called")
		}).
		WithError(func(err error) {
			errorCount++
			t.Logf("OnError called: %v", err)
		}).
		WithMessage(func(data []byte) {
			// Just log that we received a message
			t.Log("OnMessage called")
		})

	config := &WSConfig{
		BaseWSUrl:      MainnetWSBaseUrl,
		ReconnectDelay: 100 * time.Millisecond,
		PingInterval:   1 * time.Second,
		MaxReconnects:  1,
	}

	conn := NewBinancePerpWSConn(config, subscription)

	// Test initial state
	if conn.IsConnected() {
		t.Fatal("Connection should not be established initially")
	}

	// Test connection
	ctx := context.Background()
	if err := conn.Connect(ctx, "btcusdt@aggTrade"); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Verify connection established
	if !conn.IsConnected() {
		t.Fatal("Connection should be established after Connect()")
	}

	// Give some time for potential messages
	time.Sleep(1 * time.Second)

	// Test disconnection
	if err := conn.Disconnect(); err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify callbacks were called
	if connectCount != 1 {
		t.Errorf("Expected OnConnect to be called once, got %d", connectCount)
	}

	if disconnectCount != 1 {
		t.Errorf("Expected OnClose to be called once, got %d", disconnectCount)
	}

	// Verify final state
	if conn.IsConnected() {
		t.Fatal("Connection should be closed after Disconnect()")
	}

	t.Logf("Lifecycle test completed - Connect: %d, Disconnect: %d, Errors: %d",
		connectCount, disconnectCount, errorCount)
}

func TestBinancePerpWSConn_InvalidStream(t *testing.T) {
	// Test behavior with invalid stream name
	errorCh := make(chan error, 1)

	subscription := &Subscription{}
	subscription.
		WithError(func(err error) {
			select {
			case errorCh <- err:
			default:
			}
		}).
		WithMessage(func(data []byte) {
			t.Log("Unexpected message received for invalid stream")
		})

	config := &WSConfig{
		BaseWSUrl:      MainnetWSBaseUrl,
		ReconnectDelay: 100 * time.Millisecond,
		PingInterval:   1 * time.Second,
		MaxReconnects:  0, // No reconnects for this test
	}

	conn := NewBinancePerpWSConn(config, subscription)

	// Try to connect to invalid stream
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := conn.Connect(ctx, "invalid@stream@name")

	// Either the connection should fail immediately, or we should receive an error callback
	if err == nil {
		// If connection succeeded, wait for error callback
		select {
		case <-errorCh:
			t.Log("Received expected error for invalid stream")
		case <-time.After(3 * time.Second):
			// Some invalid streams might not immediately error, that's okay
			t.Log("No immediate error for invalid stream (acceptable)")
		}
	} else {
		t.Logf("Expected connection error for invalid stream: %v", err)
	}

	// Clean up
	conn.Disconnect()
}
