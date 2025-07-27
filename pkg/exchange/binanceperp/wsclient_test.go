package binanceperp

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestWSClient_SubscribeKline(t *testing.T) {
	// Test configuration (Binance WebSocket expects lowercase symbols)
	symbol := "btcusdt"
	interval := "1m"
	timeout := 5 * time.Second

	// Create WSClient
	client := NewWSClient(&WSConfig{
		BaseWSUrl:      MainnetWSBaseUrl,
		ReconnectDelay: 1 * time.Second,  // Faster reconnect for tests
		PingInterval:   30 * time.Second, // Longer ping interval for tests
		MaxReconnects:  3,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var klineCount int64
	var disconnectCount int64

	// Store received klines for validation
	var lastKline WSKline

	// Create subscription options with callbacks that count invocations
	options := &KlineSubscriptionOptions{}
	options.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		}).
		WithReconnect(func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		}).
		WithError(func(err error) {
			atomic.AddInt64(&errorCount, 1)
			t.Errorf("OnError called unexpectedly: %v", err)
		}).
		WithKline(func(kline WSKline) {
			count := atomic.AddInt64(&klineCount, 1)
			lastKline = kline
			t.Logf("OnKline called #%d: Symbol=%s, Interval=%s, Open=%s, Close=%s, IsClosed=%t",
				count, kline.Symbol, kline.Interval, kline.Open, kline.Close, kline.IsClosed)
		}).
		WithDisconnect(func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		})

	// Subscribe to kline stream
	unsubscribe, err := client.SubscribeKline(symbol, interval, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to kline stream: %v", err)
	}

	// Wait for the specified timeout
	t.Logf("Waiting %v for kline data...", timeout)
	time.Sleep(timeout)

	// Unsubscribe to trigger OnDisconnect
	unsubscribe()

	// Wait a bit for disconnect to be processed
	time.Sleep(200 * time.Millisecond)

	// Verify callback invocation counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalReconnectCount := atomic.LoadInt64(&reconnectCount)
	finalErrorCount := atomic.LoadInt64(&errorCount)
	finalKlineCount := atomic.LoadInt64(&klineCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("Callback invocation counts:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnReconnect: %d", finalReconnectCount)
	t.Logf("  OnError: %d", finalErrorCount)
	t.Logf("  OnKline: %d", finalKlineCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Test requirements verification:

	// 1. OnConnect should be called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called 1 time, got %d", finalConnectCount)
	}

	// 2. OnKline should work and deserialize correctly, count >= 1
	if finalKlineCount < 1 {
		t.Errorf("Expected OnKline to be called at least 1 time, got %d", finalKlineCount)
	}

	// Verify the last kline data was deserialized correctly
	if finalKlineCount > 0 {
		if lastKline.Symbol != "BTCUSDT" && lastKline.Symbol != "btcusdt" {
			t.Errorf("Expected symbol BTCUSDT/btcusdt, got: %s", lastKline.Symbol)
		}

		if lastKline.Interval != interval {
			t.Errorf("Expected interval %s, got: %s", interval, lastKline.Interval)
		}

		if lastKline.Open == "" {
			t.Error("Open price should not be empty")
		}

		if lastKline.Close == "" {
			t.Error("Close price should not be empty")
		}

		if lastKline.High == "" {
			t.Error("High price should not be empty")
		}

		if lastKline.Low == "" {
			t.Error("Low price should not be empty")
		}

		if lastKline.Volume == "" {
			t.Error("Volume should not be empty")
		}

		if lastKline.StartTime <= 0 {
			t.Error("StartTime should be positive")
		}

		if lastKline.CloseTime <= 0 {
			t.Error("CloseTime should be positive")
		}

		t.Logf("✓ Kline data deserialized correctly:")
		t.Logf("  Symbol: %s", lastKline.Symbol)
		t.Logf("  Interval: %s", lastKline.Interval)
		t.Logf("  Open: %s", lastKline.Open)
		t.Logf("  Close: %s", lastKline.Close)
		t.Logf("  High: %s", lastKline.High)
		t.Logf("  Low: %s", lastKline.Low)
		t.Logf("  Volume: %s", lastKline.Volume)
		t.Logf("  IsClosed: %t", lastKline.IsClosed)
	}

	// 3. No OnError should be called
	if finalErrorCount > 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
	}

	// 4. OnDisconnect should be called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called 1 time, got %d", finalDisconnectCount)
	}

	// Verify no reconnects occurred during normal operation
	if finalReconnectCount > 0 {
		t.Logf("Note: %d reconnects occurred (this may be normal depending on network conditions)", finalReconnectCount)
	}
}

func TestWSClient_SubscribeKline_MultipleIntervals(t *testing.T) {
	// Test subscribing to multiple intervals for the same symbol
	symbol := "btcusdt"
	intervals := []string{"1m", "5m"}
	timeout := 3 * time.Second

	client := NewWSClient(nil) // Use default config

	var totalKlineCount int64
	var connectCount int64
	var disconnectCount int64

	unsubscribeFuncs := make([]func(), len(intervals))

	// Subscribe to multiple intervals
	for i, interval := range intervals {
		options := &KlineSubscriptionOptions{}
		options.
			WithConnect(func() {
				atomic.AddInt64(&connectCount, 1)
				t.Logf("OnConnect called for %s@kline_%s", symbol, interval)
			}).
			WithKline(func(kline WSKline) {
				count := atomic.AddInt64(&totalKlineCount, 1)
				t.Logf("OnKline #%d: %s@%s", count, kline.Symbol, kline.Interval)
			}).
			WithDisconnect(func() {
				atomic.AddInt64(&disconnectCount, 1)
				t.Logf("OnDisconnect called for %s@kline_%s", symbol, interval)
			}).
			WithError(func(err error) {
				t.Errorf("OnError called for %s@kline_%s: %v", symbol, interval, err)
			})

		unsubscribe, err := client.SubscribeKline(symbol, interval, options)
		if err != nil {
			t.Fatalf("Failed to subscribe to %s@kline_%s: %v", symbol, interval, err)
		}
		unsubscribeFuncs[i] = unsubscribe
	}

	// Verify subscription count
	if client.GetSubscriptionCount() != len(intervals) {
		t.Errorf("Expected %d subscriptions, got %d", len(intervals), client.GetSubscriptionCount())
	}

	// Wait for data
	t.Logf("Waiting %v for kline data from %d streams...", timeout, len(intervals))
	time.Sleep(timeout)

	// Unsubscribe from all
	for _, unsubscribe := range unsubscribeFuncs {
		unsubscribe()
	}

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify final counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalKlineCount := atomic.LoadInt64(&totalKlineCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("Multi-interval test results:")
	t.Logf("  Intervals: %v", intervals)
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnKline: %d", finalKlineCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Should have one connect per interval
	if finalConnectCount != int64(len(intervals)) {
		t.Errorf("Expected OnConnect to be called %d times, got %d", len(intervals), finalConnectCount)
	}

	// Should have received klines
	if finalKlineCount < 1 {
		t.Errorf("Expected to receive at least 1 kline, got %d", finalKlineCount)
	}

	// Should have one disconnect per interval
	if finalDisconnectCount != int64(len(intervals)) {
		t.Errorf("Expected OnDisconnect to be called %d times, got %d", len(intervals), finalDisconnectCount)
	}

	// Verify all subscriptions are cleaned up
	if client.GetSubscriptionCount() != 0 {
		t.Errorf("Expected 0 subscriptions after cleanup, got %d", client.GetSubscriptionCount())
	}
}

func TestWSClient_SubscribeKline_DuplicateSubscription(t *testing.T) {
	// Test that duplicate subscriptions are rejected
	symbol := "btcusdt"
	interval := "1m"

	client := NewWSClient(nil)

	options := &KlineSubscriptionOptions{}
	options.WithKline(func(kline WSKline) {
		// Do nothing
	})

	// First subscription should succeed
	unsubscribe1, err1 := client.SubscribeKline(symbol, interval, options)
	if err1 != nil {
		t.Fatalf("First subscription failed: %v", err1)
	}

	// Second subscription to same stream should fail
	_, err2 := client.SubscribeKline(symbol, interval, options)
	if err2 == nil {
		t.Fatal("Expected second subscription to same stream to fail")
	}

	t.Logf("✓ Duplicate subscription correctly rejected: %v", err2)

	// Clean up
	unsubscribe1()
	time.Sleep(100 * time.Millisecond)

	// Now a new subscription should succeed
	unsubscribe3, err3 := client.SubscribeKline(symbol, interval, options)
	if err3 != nil {
		t.Fatalf("Subscription after cleanup failed: %v", err3)
	}

	// Clean up
	unsubscribe3()
}

func TestWSClient_SubscribeKline_InvalidSymbol(t *testing.T) {
	// Test behavior with potentially invalid symbol (though it may still connect)
	symbol := "invalidcoin"
	interval := "1m"
	timeout := 2 * time.Second

	client := NewWSClient(nil)

	var errorCount int64
	var connectCount int64

	options := &KlineSubscriptionOptions{}
	options.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called for invalid symbol")
		}).
		WithError(func(err error) {
			atomic.AddInt64(&errorCount, 1)
			t.Logf("OnError called (expected for invalid symbol): %v", err)
		}).
		WithKline(func(kline WSKline) {
			t.Logf("Unexpected kline received for invalid symbol: %+v", kline)
		})

	unsubscribe, err := client.SubscribeKline(symbol, interval, options)
	if err != nil {
		// Connection might fail immediately for invalid symbols
		t.Logf("Expected connection error for invalid symbol: %v", err)
		return
	}

	// Wait a bit to see if errors occur
	time.Sleep(timeout)

	// Clean up
	unsubscribe()

	finalErrorCount := atomic.LoadInt64(&errorCount)
	finalConnectCount := atomic.LoadInt64(&connectCount)

	t.Logf("Invalid symbol test results:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnError: %d", finalErrorCount)

	// Note: Some invalid symbols might still connect but just not receive data
	// This test mainly ensures the system doesn't crash with invalid inputs
}

func TestWSClient_SubscribeAggTrade(t *testing.T) {
	// Test configuration (Binance WebSocket expects lowercase symbols)
	symbol := "btcusdt"
	timeout := 5 * time.Second

	// Create WSClient
	client := NewWSClient(&WSConfig{
		BaseWSUrl:      MainnetWSBaseUrl,
		ReconnectDelay: 1 * time.Second,  // Faster reconnect for tests
		PingInterval:   30 * time.Second, // Longer ping interval for tests
		MaxReconnects:  3,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var aggTradeCount int64
	var disconnectCount int64

	// Store received aggTrades for validation
	var lastAggTrade WSAggTrade

	// Create subscription options with callbacks that count invocations
	options := &AggTradeSubscriptionOptions{}
	options.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		}).
		WithReconnect(func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		}).
		WithError(func(err error) {
			atomic.AddInt64(&errorCount, 1)
			t.Errorf("OnError called unexpectedly: %v", err)
		}).
		WithAggTrade(func(aggTrade WSAggTrade) {
			count := atomic.AddInt64(&aggTradeCount, 1)
			lastAggTrade = aggTrade
			t.Logf("OnAggTrade called #%d: Symbol=%s, Price=%s, Quantity=%s, IsBuyerMaker=%t",
				count, aggTrade.Symbol, aggTrade.Price, aggTrade.Quantity, aggTrade.IsBuyerMaker)
		}).
		WithDisconnect(func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		})

	// Subscribe to aggTrade stream
	unsubscribe, err := client.SubscribeAggTrade(symbol, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to aggTrade stream: %v", err)
	}

	// Wait for the specified timeout
	t.Logf("Waiting %v for aggTrade data...", timeout)
	time.Sleep(timeout)

	// Unsubscribe to trigger OnDisconnect
	unsubscribe()

	// Wait a bit for disconnect to be processed
	time.Sleep(200 * time.Millisecond)

	// Verify callback invocation counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalReconnectCount := atomic.LoadInt64(&reconnectCount)
	finalErrorCount := atomic.LoadInt64(&errorCount)
	finalAggTradeCount := atomic.LoadInt64(&aggTradeCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("Callback invocation counts:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnReconnect: %d", finalReconnectCount)
	t.Logf("  OnError: %d", finalErrorCount)
	t.Logf("  OnAggTrade: %d", finalAggTradeCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Test requirements verification:

	// 1. OnConnect should be called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called 1 time, got %d", finalConnectCount)
	}

	// 2. OnAggTrade should work and deserialize correctly, count >= 1
	if finalAggTradeCount < 1 {
		t.Errorf("Expected OnAggTrade to be called at least 1 time, got %d", finalAggTradeCount)
	}

	// Verify the last aggTrade data was deserialized correctly
	if finalAggTradeCount > 0 {
		if lastAggTrade.Symbol != "BTCUSDT" && lastAggTrade.Symbol != "btcusdt" {
			t.Errorf("Expected symbol BTCUSDT/btcusdt, got: %s", lastAggTrade.Symbol)
		}

		if lastAggTrade.EventType != "aggTrade" {
			t.Errorf("Expected event type 'aggTrade', got: %s", lastAggTrade.EventType)
		}

		if lastAggTrade.Price == "" {
			t.Error("Price should not be empty")
		}

		if lastAggTrade.Quantity == "" {
			t.Error("Quantity should not be empty")
		}

		if lastAggTrade.EventTime <= 0 {
			t.Error("EventTime should be positive")
		}

		if lastAggTrade.TradeTime <= 0 {
			t.Error("TradeTime should be positive")
		}

		if lastAggTrade.AggTradeID <= 0 {
			t.Error("AggTradeID should be positive")
		}

		if lastAggTrade.FirstTradeID <= 0 {
			t.Error("FirstTradeID should be positive")
		}

		if lastAggTrade.LastTradeID <= 0 {
			t.Error("LastTradeID should be positive")
		}

		t.Logf("✓ AggTrade data deserialized correctly:")
		t.Logf("  Symbol: %s", lastAggTrade.Symbol)
		t.Logf("  EventType: %s", lastAggTrade.EventType)
		t.Logf("  Price: %s", lastAggTrade.Price)
		t.Logf("  Quantity: %s", lastAggTrade.Quantity)
		t.Logf("  AggTradeID: %d", lastAggTrade.AggTradeID)
		t.Logf("  FirstTradeID: %d", lastAggTrade.FirstTradeID)
		t.Logf("  LastTradeID: %d", lastAggTrade.LastTradeID)
		t.Logf("  IsBuyerMaker: %t", lastAggTrade.IsBuyerMaker)
		t.Logf("  EventTime: %d", lastAggTrade.EventTime)
		t.Logf("  TradeTime: %d", lastAggTrade.TradeTime)
	}

	// 3. No OnError should be called
	if finalErrorCount > 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
	}

	// 4. OnDisconnect should be called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called 1 time, got %d", finalDisconnectCount)
	}

	// Verify no reconnects occurred during normal operation
	if finalReconnectCount > 0 {
		t.Logf("Note: %d reconnects occurred (this may be normal depending on network conditions)", finalReconnectCount)
	}
}

func TestWSClient_SubscribeKlineAndAggTrade(t *testing.T) {
	// Test subscribing to both Kline and AggTrade streams for the same symbol
	symbol := "btcusdt"
	interval := "1m"
	timeout := 3 * time.Second

	client := NewWSClient(nil) // Use default config

	var totalKlineCount int64
	var totalAggTradeCount int64
	var connectCount int64
	var disconnectCount int64

	// Subscribe to Kline stream
	klineOptions := &KlineSubscriptionOptions{}
	klineOptions.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Logf("OnConnect called for %s@kline_%s", symbol, interval)
		}).
		WithKline(func(kline WSKline) {
			count := atomic.AddInt64(&totalKlineCount, 1)
			t.Logf("OnKline #%d: %s@%s", count, kline.Symbol, kline.Interval)
		}).
		WithDisconnect(func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Logf("OnDisconnect called for %s@kline_%s", symbol, interval)
		}).
		WithError(func(err error) {
			t.Errorf("OnError called for %s@kline_%s: %v", symbol, interval, err)
		})

	unsubscribeKline, err := client.SubscribeKline(symbol, interval, klineOptions)
	if err != nil {
		t.Fatalf("Failed to subscribe to %s@kline_%s: %v", symbol, interval, err)
	}

	// Subscribe to AggTrade stream
	aggTradeOptions := &AggTradeSubscriptionOptions{}
	aggTradeOptions.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Logf("OnConnect called for %s@aggTrade", symbol)
		}).
		WithAggTrade(func(aggTrade WSAggTrade) {
			count := atomic.AddInt64(&totalAggTradeCount, 1)
			t.Logf("OnAggTrade #%d: %s@%s", count, aggTrade.Symbol, aggTrade.Price)
		}).
		WithDisconnect(func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Logf("OnDisconnect called for %s@aggTrade", symbol)
		}).
		WithError(func(err error) {
			t.Errorf("OnError called for %s@aggTrade: %v", symbol, err)
		})

	unsubscribeAggTrade, err := client.SubscribeAggTrade(symbol, aggTradeOptions)
	if err != nil {
		t.Fatalf("Failed to subscribe to %s@aggTrade: %v", symbol, err)
	}

	// Verify subscription count
	if client.GetSubscriptionCount() != 2 {
		t.Errorf("Expected 2 subscriptions, got %d", client.GetSubscriptionCount())
	}

	// Wait for data
	t.Logf("Waiting %v for data from both streams...", timeout)
	time.Sleep(timeout)

	// Unsubscribe from both
	unsubscribeKline()
	unsubscribeAggTrade()

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify final counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalKlineCount := atomic.LoadInt64(&totalKlineCount)
	finalAggTradeCount := atomic.LoadInt64(&totalAggTradeCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("Mixed stream test results:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnKline: %d", finalKlineCount)
	t.Logf("  OnAggTrade: %d", finalAggTradeCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Should have two connects (one for each stream)
	if finalConnectCount != 2 {
		t.Errorf("Expected OnConnect to be called 2 times, got %d", finalConnectCount)
	}

	// Should have received data from both streams
	if finalKlineCount < 1 {
		t.Errorf("Expected to receive at least 1 kline, got %d", finalKlineCount)
	}

	if finalAggTradeCount < 1 {
		t.Errorf("Expected to receive at least 1 aggTrade, got %d", finalAggTradeCount)
	}

	// Should have two disconnects (one for each stream)
	if finalDisconnectCount != 2 {
		t.Errorf("Expected OnDisconnect to be called 2 times, got %d", finalDisconnectCount)
	}

	// Verify all subscriptions are cleaned up
	if client.GetSubscriptionCount() != 0 {
		t.Errorf("Expected 0 subscriptions after cleanup, got %d", client.GetSubscriptionCount())
	}
}

func TestWSClient_SubscribeAggTrade_DuplicateSubscription(t *testing.T) {
	// Test that duplicate AggTrade subscriptions are rejected
	symbol := "btcusdt"

	client := NewWSClient(nil)

	options := &AggTradeSubscriptionOptions{}
	options.WithAggTrade(func(aggTrade WSAggTrade) {
		// Do nothing
	})

	// First subscription should succeed
	unsubscribe1, err1 := client.SubscribeAggTrade(symbol, options)
	if err1 != nil {
		t.Fatalf("First subscription failed: %v", err1)
	}

	// Second subscription to same stream should fail
	_, err2 := client.SubscribeAggTrade(symbol, options)
	if err2 == nil {
		t.Fatal("Expected second subscription to same stream to fail")
	}

	t.Logf("✓ Duplicate AggTrade subscription correctly rejected: %v", err2)

	// Clean up
	unsubscribe1()
	time.Sleep(100 * time.Millisecond)

	// Now a new subscription should succeed
	unsubscribe3, err3 := client.SubscribeAggTrade(symbol, options)
	if err3 != nil {
		t.Fatalf("Subscription after cleanup failed: %v", err3)
	}

	// Clean up
	unsubscribe3()
}

func TestWSClient_SubscribeTicker(t *testing.T) {
	// Test configuration (Binance WebSocket expects lowercase symbols)
	symbol := "btcusdt"
	timeout := 5 * time.Second

	// Create WSClient
	client := NewWSClient(&WSConfig{
		BaseWSUrl:      MainnetWSBaseUrl,
		ReconnectDelay: 1 * time.Second,  // Faster reconnect for tests
		PingInterval:   30 * time.Second, // Longer ping interval for tests
		MaxReconnects:  3,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var tickerCount int64
	var disconnectCount int64

	// Store received tickers for validation
	var lastTicker WSTicker

	// Create subscription options with callbacks that count invocations
	options := &TickerSubscriptionOptions{}
	options.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		}).
		WithReconnect(func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		}).
		WithError(func(err error) {
			atomic.AddInt64(&errorCount, 1)
			t.Errorf("OnError called unexpectedly: %v", err)
		}).
		WithTicker(func(ticker WSTicker) {
			count := atomic.AddInt64(&tickerCount, 1)
			lastTicker = ticker
			t.Logf("OnTicker called #%d: Symbol=%s, LastPrice=%s, PriceChange=%s%%",
				count, ticker.Symbol, ticker.LastPrice, ticker.PriceChangePercent)
		}).
		WithDisconnect(func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		})

	// Subscribe to ticker stream
	unsubscribe, err := client.SubscribeTicker(symbol, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to ticker stream: %v", err)
	}

	// Wait for the specified timeout
	t.Logf("Waiting %v for ticker data...", timeout)
	time.Sleep(timeout)

	// Unsubscribe to trigger OnDisconnect
	unsubscribe()

	// Wait a bit for disconnect to be processed
	time.Sleep(200 * time.Millisecond)

	// Verify callback invocation counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalReconnectCount := atomic.LoadInt64(&reconnectCount)
	finalErrorCount := atomic.LoadInt64(&errorCount)
	finalTickerCount := atomic.LoadInt64(&tickerCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("Callback invocation counts:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnReconnect: %d", finalReconnectCount)
	t.Logf("  OnError: %d", finalErrorCount)
	t.Logf("  OnTicker: %d", finalTickerCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Test requirements verification:

	// 1. OnConnect should be called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called 1 time, got %d", finalConnectCount)
	}

	// 2. OnTicker should work and deserialize correctly, count >= 1
	if finalTickerCount < 1 {
		t.Errorf("Expected OnTicker to be called at least 1 time, got %d", finalTickerCount)
	}

	// Verify the last ticker data was deserialized correctly
	if finalTickerCount > 0 {
		if lastTicker.Symbol != "BTCUSDT" && lastTicker.Symbol != "btcusdt" {
			t.Errorf("Expected symbol BTCUSDT/btcusdt, got: %s", lastTicker.Symbol)
		}

		if lastTicker.EventType != "24hrTicker" {
			t.Errorf("Expected event type '24hrTicker', got: %s", lastTicker.EventType)
		}

		if lastTicker.LastPrice == "" {
			t.Error("LastPrice should not be empty")
		}

		if lastTicker.OpenPrice == "" {
			t.Error("OpenPrice should not be empty")
		}

		if lastTicker.HighPrice == "" {
			t.Error("HighPrice should not be empty")
		}

		if lastTicker.LowPrice == "" {
			t.Error("LowPrice should not be empty")
		}

		if lastTicker.Volume == "" {
			t.Error("Volume should not be empty")
		}

		if lastTicker.PriceChange == "" {
			t.Error("PriceChange should not be empty")
		}

		if lastTicker.PriceChangePercent == "" {
			t.Error("PriceChangePercent should not be empty")
		}

		if lastTicker.EventTime <= 0 {
			t.Error("EventTime should be positive")
		}

		if lastTicker.OpenTime <= 0 {
			t.Error("OpenTime should be positive")
		}

		if lastTicker.CloseTime <= 0 {
			t.Error("CloseTime should be positive")
		}

		if lastTicker.Count <= 0 {
			t.Error("Count should be positive")
		}

		t.Logf("✓ Ticker data deserialized correctly:")
		t.Logf("  Symbol: %s", lastTicker.Symbol)
		t.Logf("  EventType: %s", lastTicker.EventType)
		t.Logf("  LastPrice: %s", lastTicker.LastPrice)
		t.Logf("  OpenPrice: %s", lastTicker.OpenPrice)
		t.Logf("  HighPrice: %s", lastTicker.HighPrice)
		t.Logf("  LowPrice: %s", lastTicker.LowPrice)
		t.Logf("  Volume: %s", lastTicker.Volume)
		t.Logf("  PriceChange: %s", lastTicker.PriceChange)
		t.Logf("  PriceChangePercent: %s", lastTicker.PriceChangePercent)
		t.Logf("  WeightedAvgPrice: %s", lastTicker.WeightedAvgPrice)
		t.Logf("  Count: %d", lastTicker.Count)
	}

	// 3. No OnError should be called
	if finalErrorCount > 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
	}

	// 4. OnDisconnect should be called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called 1 time, got %d", finalDisconnectCount)
	}

	// Verify no reconnects occurred during normal operation
	if finalReconnectCount > 0 {
		t.Logf("Note: %d reconnects occurred (this may be normal depending on network conditions)", finalReconnectCount)
	}
}

func TestWSClient_SubscribeLiquidation_ConnectionOnly(t *testing.T) {
	// Test liquidation subscription focusing on connection lifecycle only
	// Since liquidations are event-driven and may not occur during testing
	symbol := "ethusdt"
	timeout := 2 * time.Second

	client := NewWSClient(nil) // Use default config

	var connectCount int64
	var disconnectCount int64
	var errorCount int64

	options := &LiquidationSubscriptionOptions{}
	options.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("✓ Connected to liquidation stream")
		}).
		WithLiquidation(func(liquidation WSLiquidation) {
			t.Logf("Liquidation received: %+v", liquidation.Order)
		}).
		WithDisconnect(func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("✓ Disconnected from liquidation stream")
		}).
		WithError(func(err error) {
			atomic.AddInt64(&errorCount, 1)
			t.Errorf("OnError called: %v", err)
		})

	unsubscribe, err := client.SubscribeLiquidation(symbol, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to liquidation stream: %v", err)
	}

	// Brief wait then disconnect
	time.Sleep(timeout)
	unsubscribe()
	time.Sleep(100 * time.Millisecond)

	// Verify basic lifecycle
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)
	finalErrorCount := atomic.LoadInt64(&errorCount)

	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called 1 time, got %d", finalConnectCount)
	}

	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called 1 time, got %d", finalDisconnectCount)
	}

	if finalErrorCount > 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
	}

	t.Log("✓ Liquidation stream connection lifecycle test passed")
}

func TestWSClient_SubscribeDepth_Basic(t *testing.T) {
	// Test basic depth subscription functionality
	symbol := "btcusdt"
	level := DepthLevel5
	updateSpeed := DepthUpdate250ms // Default speed
	timeout := 5 * time.Second

	client := NewWSClient(&WSConfig{
		BaseWSUrl:      MainnetWSBaseUrl,
		ReconnectDelay: 1 * time.Second,
		PingInterval:   30 * time.Second,
		MaxReconnects:  3,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var depthCount int64
	var disconnectCount int64

	// Store received depth data for validation
	var lastDepth WSDepth

	// Create subscription options with callbacks
	options := &DepthSubscriptionOptions{}
	options.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		}).
		WithReconnect(func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		}).
		WithError(func(err error) {
			atomic.AddInt64(&errorCount, 1)
			t.Errorf("OnError called unexpectedly: %v", err)
		}).
		WithDepth(func(depth WSDepth) {
			count := atomic.AddInt64(&depthCount, 1)
			lastDepth = depth
			t.Logf("OnDepth called #%d: Symbol=%s, Bids=%d, Asks=%d, UpdateIDs=%d-%d",
				count, depth.Symbol, len(depth.Bids), len(depth.Asks),
				depth.FirstUpdateID, depth.FinalUpdateID)
		}).
		WithDisconnect(func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		})

	// Subscribe to depth stream
	unsubscribe, err := client.SubscribeDepth(symbol, level, updateSpeed, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to depth stream: %v", err)
	}

	// Wait for the specified timeout
	t.Logf("Waiting %v for depth data...", timeout)
	time.Sleep(timeout)

	// Unsubscribe to trigger OnDisconnect
	unsubscribe()

	// Wait for disconnect to be processed
	time.Sleep(200 * time.Millisecond)

	// Verify callback invocation counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalReconnectCount := atomic.LoadInt64(&reconnectCount)
	finalErrorCount := atomic.LoadInt64(&errorCount)
	finalDepthCount := atomic.LoadInt64(&depthCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("Callback invocation counts:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnReconnect: %d", finalReconnectCount)
	t.Logf("  OnError: %d", finalErrorCount)
	t.Logf("  OnDepth: %d", finalDepthCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Test requirements verification:

	// 1. OnConnect should be called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called 1 time, got %d", finalConnectCount)
	}

	// 2. OnDepth should receive data and work correctly
	if finalDepthCount < 1 {
		t.Errorf("Expected OnDepth to be called at least 1 time, got %d", finalDepthCount)
	} else {
		// Verify the last depth data was deserialized correctly
		if lastDepth.Symbol != "BTCUSDT" && lastDepth.Symbol != "btcusdt" {
			t.Errorf("Expected symbol BTCUSDT/btcusdt, got: %s", lastDepth.Symbol)
		}

		if lastDepth.EventType != "depthUpdate" {
			t.Errorf("Expected event type 'depthUpdate', got: %s", lastDepth.EventType)
		}

		// Validate basic depth structure
		if len(lastDepth.Bids) == 0 && len(lastDepth.Asks) == 0 {
			t.Error("Expected at least some bids or asks")
		}

		// Verify depth level constraint (should have at most 'level' entries)
		if len(lastDepth.Bids) > int(level) {
			t.Errorf("Expected at most %d bids, got %d", level, len(lastDepth.Bids))
		}
		if len(lastDepth.Asks) > int(level) {
			t.Errorf("Expected at most %d asks, got %d", level, len(lastDepth.Asks))
		}

		// Validate bid/ask format [price, quantity]
		for i, bid := range lastDepth.Bids {
			if len(bid) != 2 {
				t.Errorf("Bid %d should have 2 elements [price, quantity], got %d", i, len(bid))
			}
			if bid[0] == "" || bid[1] == "" {
				t.Errorf("Bid %d has empty price or quantity: %v", i, bid)
			}
		}

		for i, ask := range lastDepth.Asks {
			if len(ask) != 2 {
				t.Errorf("Ask %d should have 2 elements [price, quantity], got %d", i, len(ask))
			}
			if ask[0] == "" || ask[1] == "" {
				t.Errorf("Ask %d has empty price or quantity: %v", i, ask)
			}
		}

		if lastDepth.EventTime <= 0 {
			t.Error("EventTime should be positive")
		}

		if lastDepth.TransactionTime <= 0 {
			t.Error("TransactionTime should be positive")
		}

		t.Logf("✓ Depth data deserialized correctly:")
		t.Logf("  Symbol: %s", lastDepth.Symbol)
		t.Logf("  EventType: %s", lastDepth.EventType)
		t.Logf("  Bids: %d entries", len(lastDepth.Bids))
		t.Logf("  Asks: %d entries", len(lastDepth.Asks))
		t.Logf("  UpdateIDs: %d -> %d (prev: %d)", lastDepth.FirstUpdateID, lastDepth.FinalUpdateID, lastDepth.PrevUpdateID)
		if len(lastDepth.Bids) > 0 {
			t.Logf("  Sample Bid: [%s, %s]", lastDepth.Bids[0][0], lastDepth.Bids[0][1])
		}
		if len(lastDepth.Asks) > 0 {
			t.Logf("  Sample Ask: [%s, %s]", lastDepth.Asks[0][0], lastDepth.Asks[0][1])
		}
	}

	// 3. No OnError should be called
	if finalErrorCount > 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
	}

	// 4. OnDisconnect should be called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called 1 time, got %d", finalDisconnectCount)
	}

	// Verify no reconnects occurred during normal operation
	if finalReconnectCount > 0 {
		t.Logf("Note: %d reconnects occurred (this may be normal depending on network conditions)", finalReconnectCount)
	}
}

func TestWSClient_SubscribeDepth_DifferentLevels(t *testing.T) {
	// Test different depth levels
	symbol := "ethusdt"
	timeout := 3 * time.Second

	client := NewWSClient(nil) // Use default config

	levels := []DepthLevel{DepthLevel5, DepthLevel10, DepthLevel20}

	for _, level := range levels {
		t.Run(fmt.Sprintf("Level%d", level), func(t *testing.T) {
			var connectCount int64
			var depthCount int64
			var disconnectCount int64

			options := &DepthSubscriptionOptions{}
			options.
				WithConnect(func() {
					atomic.AddInt64(&connectCount, 1)
					t.Logf("✓ Connected to depth%d stream", level)
				}).
				WithDepth(func(depth WSDepth) {
					atomic.AddInt64(&depthCount, 1)
					t.Logf("Depth%d: %d bids, %d asks", level, len(depth.Bids), len(depth.Asks))
				}).
				WithDisconnect(func() {
					atomic.AddInt64(&disconnectCount, 1)
					t.Logf("✓ Disconnected from depth%d stream", level)
				}).
				WithError(func(err error) {
					t.Errorf("OnError called: %v", err)
				})

			unsubscribe, err := client.SubscribeDepth(symbol, level, DepthUpdate250ms, options)
			if err != nil {
				t.Fatalf("Failed to subscribe to depth%d stream: %v", level, err)
			}

			time.Sleep(timeout)
			unsubscribe()
			time.Sleep(100 * time.Millisecond)

			// Verify basic lifecycle for each level
			finalConnectCount := atomic.LoadInt64(&connectCount)
			finalDepthCount := atomic.LoadInt64(&depthCount)
			finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

			if finalConnectCount != 1 {
				t.Errorf("Expected OnConnect=1, got %d", finalConnectCount)
			}

			if finalDepthCount < 1 {
				t.Errorf("Expected OnDepth>=1, got %d", finalDepthCount)
			}

			if finalDisconnectCount != 1 {
				t.Errorf("Expected OnDisconnect=1, got %d", finalDisconnectCount)
			}
		})
	}
}

func TestWSClient_SubscribeDepth_DifferentUpdateSpeeds(t *testing.T) {
	// Test different update speeds
	symbol := "btcusdt"
	level := DepthLevel5
	timeout := 3 * time.Second

	client := NewWSClient(nil)

	speeds := []DepthUpdateSpeed{DepthUpdate100ms, DepthUpdate250ms, DepthUpdate500ms}

	for _, speed := range speeds {
		t.Run(fmt.Sprintf("Speed%s", speed), func(t *testing.T) {
			var connectCount int64
			var depthCount int64
			var disconnectCount int64

			options := &DepthSubscriptionOptions{}
			options.
				WithConnect(func() {
					atomic.AddInt64(&connectCount, 1)
					t.Logf("✓ Connected to depth stream @ %s", speed)
				}).
				WithDepth(func(depth WSDepth) {
					atomic.AddInt64(&depthCount, 1)
				}).
				WithDisconnect(func() {
					atomic.AddInt64(&disconnectCount, 1)
					t.Logf("✓ Disconnected from depth stream @ %s", speed)
				}).
				WithError(func(err error) {
					t.Errorf("OnError called: %v", err)
				})

			unsubscribe, err := client.SubscribeDepth(symbol, level, speed, options)
			if err != nil {
				t.Fatalf("Failed to subscribe to depth stream @ %s: %v", speed, err)
			}

			time.Sleep(timeout)
			unsubscribe()
			time.Sleep(100 * time.Millisecond)

			// Verify basic lifecycle for each speed
			finalConnectCount := atomic.LoadInt64(&connectCount)
			finalDepthCount := atomic.LoadInt64(&depthCount)
			finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

			if finalConnectCount != 1 {
				t.Errorf("Expected OnConnect=1, got %d", finalConnectCount)
			}

			if finalDepthCount < 1 {
				t.Errorf("Expected OnDepth>=1, got %d", finalDepthCount)
			}

			if finalDisconnectCount != 1 {
				t.Errorf("Expected OnDisconnect=1, got %d", finalDisconnectCount)
			}
		})
	}
}

func TestWSClient_SubscribeDepth_InvalidParameters(t *testing.T) {
	client := NewWSClient(nil)
	symbol := "btcusdt"

	options := &DepthSubscriptionOptions{}

	// Test invalid depth level
	_, err := client.SubscribeDepth(symbol, DepthLevel(999), DepthUpdate250ms, options)
	if err == nil {
		t.Error("Expected error for invalid depth level")
	}
	t.Logf("✓ Invalid depth level correctly rejected: %v", err)

	// Test invalid update speed
	_, err = client.SubscribeDepth(symbol, DepthLevel5, DepthUpdateSpeed("invalid"), options)
	if err == nil {
		t.Error("Expected error for invalid update speed")
	}
	t.Logf("✓ Invalid update speed correctly rejected: %v", err)
}

func TestWSClient_SubscribeDepth_DuplicateSubscription(t *testing.T) {
	// Test that duplicate depth subscriptions are rejected
	symbol := "btcusdt"
	level := DepthLevel10
	speed := DepthUpdate250ms

	client := NewWSClient(nil)

	options := &DepthSubscriptionOptions{}
	options.WithDepth(func(depth WSDepth) {
		// Do nothing
	})

	// First subscription should succeed
	unsubscribe1, err1 := client.SubscribeDepth(symbol, level, speed, options)
	if err1 != nil {
		t.Fatalf("First subscription failed: %v", err1)
	}

	// Second subscription to same stream should fail
	_, err2 := client.SubscribeDepth(symbol, level, speed, options)
	if err2 == nil {
		t.Fatal("Expected second subscription to same stream to fail")
	}

	t.Logf("✓ Duplicate depth subscription correctly rejected: %v", err2)

	// Clean up
	unsubscribe1()
	time.Sleep(100 * time.Millisecond)

	// Now a new subscription should succeed
	unsubscribe3, err3 := client.SubscribeDepth(symbol, level, speed, options)
	if err3 != nil {
		t.Fatalf("Subscription after cleanup failed: %v", err3)
	}

	// Clean up
	unsubscribe3()
}

func TestWSClient_SubscribeDiffDepth_Basic(t *testing.T) {
	// Test basic differential depth subscription functionality
	symbol := "btcusdt"
	updateSpeed := DepthUpdate250ms // Default speed
	timeout := 5 * time.Second

	client := NewWSClient(&WSConfig{
		BaseWSUrl:      MainnetWSBaseUrl,
		ReconnectDelay: 1 * time.Second,
		PingInterval:   30 * time.Second,
		MaxReconnects:  3,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var diffDepthCount int64
	var disconnectCount int64

	// Store received differential depth data for validation
	var lastDiffDepth WSDepth

	// Create subscription options with callbacks
	options := &DiffDepthSubscriptionOptions{}
	options.
		WithConnect(func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		}).
		WithReconnect(func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		}).
		WithError(func(err error) {
			atomic.AddInt64(&errorCount, 1)
			t.Errorf("OnError called unexpectedly: %v", err)
		}).
		WithDiffDepth(func(diffDepth WSDepth) {
			count := atomic.AddInt64(&diffDepthCount, 1)
			lastDiffDepth = diffDepth
			t.Logf("OnDiffDepth called #%d: Symbol=%s, Bids=%d, Asks=%d, UpdateIDs=%d-%d",
				count, diffDepth.Symbol, len(diffDepth.Bids), len(diffDepth.Asks),
				diffDepth.FirstUpdateID, diffDepth.FinalUpdateID)
		}).
		WithDisconnect(func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		})

	// Subscribe to differential depth stream
	unsubscribe, err := client.SubscribeDiffDepth(symbol, updateSpeed, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to differential depth stream: %v", err)
	}

	// Wait for the specified timeout
	t.Logf("Waiting %v for differential depth data...", timeout)
	time.Sleep(timeout)

	// Unsubscribe to trigger OnDisconnect
	unsubscribe()

	// Wait for disconnect to be processed
	time.Sleep(200 * time.Millisecond)

	// Verify callback invocation counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalReconnectCount := atomic.LoadInt64(&reconnectCount)
	finalErrorCount := atomic.LoadInt64(&errorCount)
	finalDiffDepthCount := atomic.LoadInt64(&diffDepthCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("Callback invocation counts:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnReconnect: %d", finalReconnectCount)
	t.Logf("  OnError: %d", finalErrorCount)
	t.Logf("  OnDiffDepth: %d", finalDiffDepthCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Test requirements verification:

	// 1. OnConnect should be called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called 1 time, got %d", finalConnectCount)
	}

	// 2. OnDiffDepth should receive data and work correctly
	if finalDiffDepthCount < 1 {
		t.Errorf("Expected OnDiffDepth to be called at least 1 time, got %d", finalDiffDepthCount)
	} else {
		// Verify the last differential depth data was deserialized correctly
		if lastDiffDepth.Symbol != "BTCUSDT" && lastDiffDepth.Symbol != "btcusdt" {
			t.Errorf("Expected symbol BTCUSDT/btcusdt, got: %s", lastDiffDepth.Symbol)
		}

		if lastDiffDepth.EventType != "depthUpdate" {
			t.Errorf("Expected event type 'depthUpdate', got: %s", lastDiffDepth.EventType)
		}

		// Validate bid/ask format [price, quantity] for differential updates
		for i, bid := range lastDiffDepth.Bids {
			if len(bid) != 2 {
				t.Errorf("Bid %d should have 2 elements [price, quantity], got %d", i, len(bid))
			}
			if bid[0] == "" {
				t.Errorf("Bid %d has empty price: %v", i, bid)
			}
			// Note: Quantity can be "0" in differential updates (removal)
		}

		for i, ask := range lastDiffDepth.Asks {
			if len(ask) != 2 {
				t.Errorf("Ask %d should have 2 elements [price, quantity], got %d", i, len(ask))
			}
			if ask[0] == "" {
				t.Errorf("Ask %d has empty price: %v", i, ask)
			}
			// Note: Quantity can be "0" in differential updates (removal)
		}

		if lastDiffDepth.EventTime <= 0 {
			t.Error("EventTime should be positive")
		}

		if lastDiffDepth.TransactionTime <= 0 {
			t.Error("TransactionTime should be positive")
		}

		// Validate update ID sequencing
		if lastDiffDepth.FinalUpdateID < lastDiffDepth.FirstUpdateID {
			t.Errorf("FinalUpdateID (%d) should be >= FirstUpdateID (%d)",
				lastDiffDepth.FinalUpdateID, lastDiffDepth.FirstUpdateID)
		}

		t.Logf("✓ Differential depth data deserialized correctly:")
		t.Logf("  Symbol: %s", lastDiffDepth.Symbol)
		t.Logf("  EventType: %s", lastDiffDepth.EventType)
		t.Logf("  Bid Updates: %d entries", len(lastDiffDepth.Bids))
		t.Logf("  Ask Updates: %d entries", len(lastDiffDepth.Asks))
		t.Logf("  UpdateIDs: %d -> %d (prev: %d)", lastDiffDepth.FirstUpdateID, lastDiffDepth.FinalUpdateID, lastDiffDepth.PrevUpdateID)
		if len(lastDiffDepth.Bids) > 0 {
			t.Logf("  Sample Bid Update: [%s, %s]", lastDiffDepth.Bids[0][0], lastDiffDepth.Bids[0][1])
		}
		if len(lastDiffDepth.Asks) > 0 {
			t.Logf("  Sample Ask Update: [%s, %s]", lastDiffDepth.Asks[0][0], lastDiffDepth.Asks[0][1])
		}
	}

	// 3. No OnError should be called
	if finalErrorCount > 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
	}

	// 4. OnDisconnect should be called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called 1 time, got %d", finalDisconnectCount)
	}

	// Verify no reconnects occurred during normal operation
	if finalReconnectCount > 0 {
		t.Logf("Note: %d reconnects occurred (this may be normal depending on network conditions)", finalReconnectCount)
	}
}

func TestWSClient_SubscribeDiffDepth_InvalidParameters(t *testing.T) {
	client := NewWSClient(nil)
	symbol := "btcusdt"

	options := &DiffDepthSubscriptionOptions{}

	// Test invalid update speed
	_, err := client.SubscribeDiffDepth(symbol, DepthUpdateSpeed("invalid"), options)
	if err == nil {
		t.Error("Expected error for invalid update speed")
	}
	t.Logf("✓ Invalid update speed correctly rejected: %v", err)

	// Test empty string defaults to 250ms
	unsubscribe, err := client.SubscribeDiffDepth(symbol, "", options)
	if err != nil {
		t.Errorf("Expected empty string to default to 250ms, got error: %v", err)
	} else {
		t.Log("✓ Empty string correctly defaults to 250ms")
		unsubscribe()
	}
}
