package binanceperp

import (
	"context"
	"encoding/json"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

// Use WSAggTradeEvent from ws_model.go

func TestBinancePerpWSConn_AggTradePayload(t *testing.T) {
	// Test configuration for Binance perpetual futures (lowercase symbols)
	symbol := "btcusdt"
	streamName := symbol + "@aggTrade"
	timeout := 10 * time.Second

	// Create message channel for testing
	msgCh := make(chan WSAggTradeEvent, 1)
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
			var event WSAggTradeEvent
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

func TestBinancePerpUserDataStream_ConnectionLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	// Check for API credentials
	apiKey := os.Getenv("BINANCEPERP_API_KEY")
	apiSecret := os.Getenv("BINANCEPERP_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCEPERP_API_KEY or BINANCEPERP_API_SECRET not set; skipping user data stream test.")
	}

	// Create REST client for listen key management
	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)

	// Create WebSocket config
	wsConfig := &WSConfig{
		BaseWSUrl:      MainnetWSBaseUrl,
		ReconnectDelay: 1 * time.Second,
		PingInterval:   30 * time.Second,
		MaxReconnects:  3,
	}

	// Track connection events
	var connectCount int64
	var messageCount int64
	var disconnectCount int64
	var errorCount int64

	// Create subscription with callbacks
	subscription := &Subscription{}
	subscription.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("✓ User data stream connected")
		}).
		WithMessage(func(data []byte) {
			count := atomic.AddInt64(&messageCount, 1)
			t.Logf("User data message #%d received (%d bytes)", count, len(data))
		}).
		WithError(func(err error) {
			atomic.AddInt64(&errorCount, 1)
			t.Logf("User data stream error: %v", err)
		}).
		WithClose(func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("✓ User data stream disconnected")
		})

	// Create user data stream
	userDataStream := NewBinancePerpUserDataStream(client, wsConfig, subscription)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := userDataStream.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect user data stream: %v", err)
	}

	// Verify connected
	if !userDataStream.IsConnected() {
		t.Error("User data stream should be connected")
	}

	// Wait a bit to receive any potential messages
	t.Log("Waiting 5 seconds for user data...")
	time.Sleep(5 * time.Second)

	// Disconnect
	err = userDataStream.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect user data stream: %v", err)
	}

	// Verify disconnected
	if userDataStream.IsConnected() {
		t.Error("User data stream should be disconnected")
	}

	// Wait for disconnect to be processed
	time.Sleep(500 * time.Millisecond)

	// Verify callback counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalMessageCount := atomic.LoadInt64(&messageCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)
	finalErrorCount := atomic.LoadInt64(&errorCount)

	t.Logf("Callback counts:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnMessage: %d", finalMessageCount)
	t.Logf("  OnError: %d", finalErrorCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Verify basic lifecycle
	if finalConnectCount != 1 {
		t.Errorf("Expected 1 connect event, got %d", finalConnectCount)
	}

	if finalDisconnectCount != 1 {
		t.Errorf("Expected 1 disconnect event, got %d", finalDisconnectCount)
	}

	// Errors should be minimal for normal operation
	if finalErrorCount > 1 {
		t.Logf("Note: %d errors occurred (this may be normal for connection management)", finalErrorCount)
	}

	t.Log("✓ User data stream lifecycle test completed successfully")
}
