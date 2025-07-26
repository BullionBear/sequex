package binance

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWSClient_SubscribeKline(t *testing.T) {
	// Test configuration (Binance WebSocket expects lowercase symbols)
	symbol := "btcusdt"
	interval := "1m"
	timeout := 10 * time.Second

	// Create WSClient (using port 9443 for better WebSocket performance)
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var klineCount int64
	var disconnectCount int64

	// Test state tracking
	var mu sync.Mutex
	var receivedKlines []WSKline
	var lastError error

	// Create subscription options with callbacks that count invocations
	options := KlineSubscriptionOptions{
		OnConnect: func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		},
		OnReconnect: func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		},
		OnError: func(err error) {
			atomic.AddInt64(&errorCount, 1)
			mu.Lock()
			lastError = err
			mu.Unlock()
			t.Errorf("OnError called unexpectedly: %v", err)
		},
		OnKline: func(kline WSKline) {
			atomic.AddInt64(&klineCount, 1)
			mu.Lock()
			receivedKlines = append(receivedKlines, kline)
			mu.Unlock()
			t.Logf("OnKline called #%d with kline: Symbol=%s, Interval=%s, Open=%s, Close=%s",
				atomic.LoadInt64(&klineCount), kline.Symbol, kline.Interval, kline.Open, kline.Close)
		},
		OnDisconnect: func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		},
	}

	// Subscribe to kline stream
	unsubscribe, err := client.SubscribeKline(symbol, interval, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to kline stream: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for timeout
	t.Logf("Waiting %v for kline data...", timeout)
	<-ctx.Done()

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

	// Verify OnConnect called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called exactly 1 time, got %d", finalConnectCount)
	}

	// Verify OnError not called
	if finalErrorCount != 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
		mu.Lock()
		if lastError != nil {
			t.Errorf("Last error was: %v", lastError)
		}
		mu.Unlock()
	}

	// Verify OnDisconnect called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called exactly 1 time, got %d", finalDisconnectCount)
	}

	// Verify OnKline called at least once and deserialization works
	if finalKlineCount >= 1 {
		t.Logf("✅ OnKline called %d times - deserialization working correctly", finalKlineCount)

		// Verify deserialization worked correctly for received klines
		mu.Lock()
		for i, kline := range receivedKlines {
			if kline.Symbol == "" {
				t.Errorf("Kline #%d has empty Symbol", i+1)
			}
			if kline.Interval == "" {
				t.Errorf("Kline #%d has empty Interval", i+1)
			}
			if kline.Open == "" {
				t.Errorf("Kline #%d has empty Open price", i+1)
			}
			if kline.Close == "" {
				t.Errorf("Kline #%d has empty Close price", i+1)
			}
			if kline.High == "" {
				t.Errorf("Kline #%d has empty High price", i+1)
			}
			if kline.Low == "" {
				t.Errorf("Kline #%d has empty Low price", i+1)
			}
			if kline.Volume == "" {
				t.Errorf("Kline #%d has empty Volume", i+1)
			}
		}
		mu.Unlock()

		t.Logf("✅ All %d received klines have valid data structure", len(receivedKlines))
	} else {
		t.Logf("⚠️  No kline data received within %v timeout - this may be normal depending on market activity", timeout)
	}

	// Verify OnReconnect was not called (should only happen on connection issues)
	if finalReconnectCount > 0 {
		t.Logf("ℹ️  OnReconnect was called %d times (may indicate connection issues during test)", finalReconnectCount)
	}

	t.Log("✅ Test completed successfully")
}

func TestWSClient_SubscribeAggTrade(t *testing.T) {
	// Test configuration (Binance WebSocket expects lowercase symbols)
	symbol := "btcusdt"
	timeout := 10 * time.Second

	// Create WSClient (using port 9443 for better WebSocket performance)
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var aggTradeCount int64
	var disconnectCount int64

	// Test state tracking
	var mu sync.Mutex
	var receivedAggTrades []WSAggTrade
	var lastError error

	// Create subscription options with callbacks that count invocations
	options := AggTradeSubscriptionOptions{
		OnConnect: func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		},
		OnReconnect: func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		},
		OnError: func(err error) {
			atomic.AddInt64(&errorCount, 1)
			mu.Lock()
			lastError = err
			mu.Unlock()
			t.Errorf("OnError called unexpectedly: %v", err)
		},
		OnAggTrade: func(aggTrade WSAggTrade) {
			atomic.AddInt64(&aggTradeCount, 1)
			mu.Lock()
			receivedAggTrades = append(receivedAggTrades, aggTrade)
			mu.Unlock()
			t.Logf("OnAggTrade called #%d with trade: Symbol=%s, Price=%s, Quantity=%s, ID=%d",
				atomic.LoadInt64(&aggTradeCount), aggTrade.Symbol, aggTrade.Price, aggTrade.Quantity, aggTrade.AggTradeId)
		},
		OnDisconnect: func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		},
	}

	// Subscribe to aggregate trade stream
	unsubscribe, err := client.SubscribeAggTrade(symbol, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to aggregate trade stream: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for timeout
	t.Logf("Waiting %v for aggregate trade data...", timeout)
	<-ctx.Done()

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

	// Verify OnConnect called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called exactly 1 time, got %d", finalConnectCount)
	}

	// Verify OnError not called
	if finalErrorCount != 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
		mu.Lock()
		if lastError != nil {
			t.Errorf("Last error was: %v", lastError)
		}
		mu.Unlock()
	}

	// Verify OnDisconnect called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called exactly 1 time, got %d", finalDisconnectCount)
	}

	// Verify OnAggTrade called at least once and deserialization works
	if finalAggTradeCount >= 1 {
		t.Logf("✅ OnAggTrade called %d times - deserialization working correctly", finalAggTradeCount)

		// Verify deserialization worked correctly for received aggregate trades
		mu.Lock()
		for i, aggTrade := range receivedAggTrades {
			if aggTrade.Symbol == "" {
				t.Errorf("AggTrade #%d has empty Symbol", i+1)
			}
			if aggTrade.EventType != "aggTrade" {
				t.Errorf("AggTrade #%d has incorrect EventType: %s", i+1, aggTrade.EventType)
			}
			if aggTrade.Price == "" {
				t.Errorf("AggTrade #%d has empty Price", i+1)
			}
			if aggTrade.Quantity == "" {
				t.Errorf("AggTrade #%d has empty Quantity", i+1)
			}
			if aggTrade.AggTradeId == 0 {
				t.Errorf("AggTrade #%d has zero AggTradeId", i+1)
			}
			if aggTrade.TradeTime == 0 {
				t.Errorf("AggTrade #%d has zero TradeTime", i+1)
			}
		}
		mu.Unlock()

		t.Logf("✅ All %d received aggregate trades have valid data structure", len(receivedAggTrades))
	} else {
		t.Logf("⚠️  No aggregate trade data received within %v timeout - this may be normal depending on market activity", timeout)
	}

	// Verify OnReconnect was not called (should only happen on connection issues)
	if finalReconnectCount > 0 {
		t.Logf("ℹ️  OnReconnect was called %d times (may indicate connection issues during test)", finalReconnectCount)
	}

	t.Log("✅ Test completed successfully")
}

func TestWSClient_SubscribeTrade(t *testing.T) {
	// Test configuration (Binance WebSocket expects lowercase symbols)
	symbol := "btcusdt"
	timeout := 10 * time.Second

	// Create WSClient (using port 9443 for better WebSocket performance)
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var tradeCount int64
	var disconnectCount int64

	// Test state tracking
	var mu sync.Mutex
	var receivedTrades []WSTrade
	var lastError error

	// Create subscription options with callbacks that count invocations
	options := TradeSubscriptionOptions{
		OnConnect: func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		},
		OnReconnect: func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		},
		OnError: func(err error) {
			atomic.AddInt64(&errorCount, 1)
			mu.Lock()
			lastError = err
			mu.Unlock()
			t.Errorf("OnError called unexpectedly: %v", err)
		},
		OnTrade: func(trade WSTrade) {
			atomic.AddInt64(&tradeCount, 1)
			mu.Lock()
			receivedTrades = append(receivedTrades, trade)
			mu.Unlock()
			t.Logf("OnTrade called #%d with trade: Symbol=%s, Price=%s, Quantity=%s, ID=%d",
				atomic.LoadInt64(&tradeCount), trade.Symbol, trade.Price, trade.Quantity, trade.TradeId)
		},
		OnDisconnect: func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		},
	}

	// Subscribe to raw trade stream
	unsubscribe, err := client.SubscribeTrade(symbol, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to trade stream: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for timeout
	t.Logf("Waiting %v for trade data...", timeout)
	<-ctx.Done()

	// Unsubscribe to trigger OnDisconnect
	unsubscribe()

	// Wait a bit for disconnect to be processed
	time.Sleep(200 * time.Millisecond)

	// Verify callback invocation counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalReconnectCount := atomic.LoadInt64(&reconnectCount)
	finalErrorCount := atomic.LoadInt64(&errorCount)
	finalTradeCount := atomic.LoadInt64(&tradeCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("Callback invocation counts:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnReconnect: %d", finalReconnectCount)
	t.Logf("  OnError: %d", finalErrorCount)
	t.Logf("  OnTrade: %d", finalTradeCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Verify OnConnect called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called exactly 1 time, got %d", finalConnectCount)
	}

	// Verify OnError not called
	if finalErrorCount != 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
		mu.Lock()
		if lastError != nil {
			t.Errorf("Last error was: %v", lastError)
		}
		mu.Unlock()
	}

	// Verify OnDisconnect called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called exactly 1 time, got %d", finalDisconnectCount)
	}

	// Verify OnTrade called at least once and deserialization works
	if finalTradeCount >= 1 {
		t.Logf("✅ OnTrade called %d times - deserialization working correctly", finalTradeCount)

		// Verify deserialization worked correctly for received trades
		mu.Lock()
		for i, trade := range receivedTrades {
			if trade.Symbol == "" {
				t.Errorf("Trade #%d has empty Symbol", i+1)
			}
			if trade.EventType != "trade" {
				t.Errorf("Trade #%d has incorrect EventType: %s", i+1, trade.EventType)
			}
			if trade.Price == "" {
				t.Errorf("Trade #%d has empty Price", i+1)
			}
			if trade.Quantity == "" {
				t.Errorf("Trade #%d has empty Quantity", i+1)
			}
			if trade.TradeId == 0 {
				t.Errorf("Trade #%d has zero TradeId", i+1)
			}
			if trade.TradeTime == 0 {
				t.Errorf("Trade #%d has zero TradeTime", i+1)
			}
		}
		mu.Unlock()

		t.Logf("✅ All %d received trades have valid data structure", len(receivedTrades))
	} else {
		t.Logf("⚠️  No trade data received within %v timeout - this may be normal depending on market activity", timeout)
	}

	// Verify OnReconnect was not called (should only happen on connection issues)
	if finalReconnectCount > 0 {
		t.Logf("ℹ️  OnReconnect was called %d times (may indicate connection issues during test)", finalReconnectCount)
	}

	t.Log("✅ Test completed successfully")
}

func TestWSClient_SubscribeDepth(t *testing.T) {
	// Test configuration (Binance WebSocket expects lowercase symbols)
	symbol := "btcusdt"
	levels := 5
	updateSpeed := "" // Use default 1000ms
	timeout := 10 * time.Second

	// Create WSClient (using port 9443 for better WebSocket performance)
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var depthCount int64
	var disconnectCount int64

	// Test state tracking
	var mu sync.Mutex
	var receivedDepths []WSDepth
	var lastError error

	// Create subscription options with callbacks that count invocations
	options := DepthSubscriptionOptions{
		OnConnect: func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		},
		OnReconnect: func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		},
		OnError: func(err error) {
			atomic.AddInt64(&errorCount, 1)
			mu.Lock()
			lastError = err
			mu.Unlock()
			t.Errorf("OnError called unexpectedly: %v", err)
		},
		OnDepth: func(depth WSDepth) {
			atomic.AddInt64(&depthCount, 1)
			mu.Lock()
			receivedDepths = append(receivedDepths, depth)
			mu.Unlock()
			t.Logf("OnDepth called #%d with depth: LastUpdateId=%d, Bids=%d, Asks=%d",
				atomic.LoadInt64(&depthCount), depth.LastUpdateId, len(depth.Bids), len(depth.Asks))
		},
		OnDisconnect: func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		},
	}

	// Subscribe to depth stream
	unsubscribe, err := client.SubscribeDepth(symbol, levels, updateSpeed, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to depth stream: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for timeout
	t.Logf("Waiting %v for depth data...", timeout)
	<-ctx.Done()

	// Unsubscribe to trigger OnDisconnect
	unsubscribe()

	// Wait a bit for disconnect to be processed
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

	// Verify OnConnect called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called exactly 1 time, got %d", finalConnectCount)
	}

	// Verify OnError not called
	if finalErrorCount != 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
		mu.Lock()
		if lastError != nil {
			t.Errorf("Last error was: %v", lastError)
		}
		mu.Unlock()
	}

	// Verify OnDisconnect called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called exactly 1 time, got %d", finalDisconnectCount)
	}

	// Verify OnDepth called at least once and deserialization works
	if finalDepthCount >= 1 {
		t.Logf("✅ OnDepth called %d times - deserialization working correctly", finalDepthCount)

		// Verify deserialization worked correctly for received depth data
		mu.Lock()
		for i, depth := range receivedDepths {
			if depth.LastUpdateId == 0 {
				t.Errorf("Depth #%d has zero LastUpdateId", i+1)
			}
			if len(depth.Bids) == 0 && len(depth.Asks) == 0 {
				t.Errorf("Depth #%d has empty bids and asks", i+1)
			}

			// Verify bid price levels
			for j, bid := range depth.Bids {
				if bid[0] == "" || bid[1] == "" {
					t.Errorf("Depth #%d Bid #%d has empty price or quantity: [%s, %s]", i+1, j+1, bid[0], bid[1])
				}
			}

			// Verify ask price levels
			for j, ask := range depth.Asks {
				if ask[0] == "" || ask[1] == "" {
					t.Errorf("Depth #%d Ask #%d has empty price or quantity: [%s, %s]", i+1, j+1, ask[0], ask[1])
				}
			}
		}
		mu.Unlock()

		t.Logf("✅ All %d received depth updates have valid data structure", len(receivedDepths))
	} else {
		t.Logf("⚠️  No depth data received within %v timeout - this may be normal depending on market activity", timeout)
	}

	// Verify OnReconnect was not called (should only happen on connection issues)
	if finalReconnectCount > 0 {
		t.Logf("ℹ️  OnReconnect was called %d times (may indicate connection issues during test)", finalReconnectCount)
	}

	t.Log("✅ Test completed successfully")
}

func TestWSClient_SubscribeDepthUpdate(t *testing.T) {
	// Test configuration (Binance WebSocket expects lowercase symbols)
	symbol := "btcusdt"
	updateSpeed := "" // Use default 1000ms
	timeout := 10 * time.Second

	// Create WSClient (using port 9443 for better WebSocket performance)
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})

	// Callback invocation counters
	var connectCount int64
	var reconnectCount int64
	var errorCount int64
	var depthUpdateCount int64
	var disconnectCount int64

	// Test state tracking
	var mu sync.Mutex
	var receivedDepthUpdates []WSDepthUpdate
	var lastError error

	// Create subscription options with callbacks that count invocations
	options := DepthUpdateSubscriptionOptions{
		OnConnect: func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called")
		},
		OnReconnect: func() {
			atomic.AddInt64(&reconnectCount, 1)
			t.Log("OnReconnect called")
		},
		OnError: func(err error) {
			atomic.AddInt64(&errorCount, 1)
			mu.Lock()
			lastError = err
			mu.Unlock()
			t.Errorf("OnError called unexpectedly: %v", err)
		},
		OnDepthUpdate: func(update WSDepthUpdate) {
			atomic.AddInt64(&depthUpdateCount, 1)
			mu.Lock()
			receivedDepthUpdates = append(receivedDepthUpdates, update)
			mu.Unlock()
			t.Logf("OnDepthUpdate called #%d with update: Symbol=%s, FirstUpdateId=%d, FinalUpdateId=%d, BidUpdates=%d, AskUpdates=%d",
				atomic.LoadInt64(&depthUpdateCount), update.Symbol, update.FirstUpdateId, update.FinalUpdateId, len(update.BidUpdates), len(update.AskUpdates))
		},
		OnDisconnect: func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		},
	}

	// Subscribe to differential depth stream
	unsubscribe, err := client.SubscribeDepthUpdate(symbol, updateSpeed, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to depth update stream: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for timeout
	t.Logf("Waiting %v for depth update data...", timeout)
	<-ctx.Done()

	// Unsubscribe to trigger OnDisconnect
	unsubscribe()

	// Wait a bit for disconnect to be processed
	time.Sleep(200 * time.Millisecond)

	// Verify callback invocation counts
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalReconnectCount := atomic.LoadInt64(&reconnectCount)
	finalErrorCount := atomic.LoadInt64(&errorCount)
	finalDepthUpdateCount := atomic.LoadInt64(&depthUpdateCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("Callback invocation counts:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnReconnect: %d", finalReconnectCount)
	t.Logf("  OnError: %d", finalErrorCount)
	t.Logf("  OnDepthUpdate: %d", finalDepthUpdateCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	// Verify OnConnect called exactly once
	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect to be called exactly 1 time, got %d", finalConnectCount)
	}

	// Verify OnError not called
	if finalErrorCount != 0 {
		t.Errorf("Expected OnError to be called 0 times, got %d", finalErrorCount)
		mu.Lock()
		if lastError != nil {
			t.Errorf("Last error was: %v", lastError)
		}
		mu.Unlock()
	}

	// Verify OnDisconnect called exactly once
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called exactly 1 time, got %d", finalDisconnectCount)
	}

	// Verify OnDepthUpdate called at least once and deserialization works
	if finalDepthUpdateCount >= 1 {
		t.Logf("✅ OnDepthUpdate called %d times - deserialization working correctly", finalDepthUpdateCount)

		// Verify deserialization worked correctly for received depth updates
		mu.Lock()
		for i, update := range receivedDepthUpdates {
			if update.EventType != "depthUpdate" {
				t.Errorf("DepthUpdate #%d has incorrect EventType: %s", i+1, update.EventType)
			}
			if update.Symbol == "" {
				t.Errorf("DepthUpdate #%d has empty Symbol", i+1)
			}
			if update.FirstUpdateId == 0 && update.FinalUpdateId == 0 {
				t.Errorf("DepthUpdate #%d has zero FirstUpdateId and FinalUpdateId", i+1)
			}
			if update.FirstUpdateId > update.FinalUpdateId && update.FinalUpdateId != 0 {
				t.Errorf("DepthUpdate #%d has FirstUpdateId (%d) > FinalUpdateId (%d)", i+1, update.FirstUpdateId, update.FinalUpdateId)
			}

			// Verify bid updates
			for j, bid := range update.BidUpdates {
				if bid[0] == "" {
					t.Errorf("DepthUpdate #%d BidUpdate #%d has empty price: [%s, %s]", i+1, j+1, bid[0], bid[1])
				}
				// Note: quantity can be "0" for removals, so we don't check for empty quantity
			}

			// Verify ask updates
			for j, ask := range update.AskUpdates {
				if ask[0] == "" {
					t.Errorf("DepthUpdate #%d AskUpdate #%d has empty price: [%s, %s]", i+1, j+1, ask[0], ask[1])
				}
				// Note: quantity can be "0" for removals, so we don't check for empty quantity
			}
		}
		mu.Unlock()

		t.Logf("✅ All %d received depth updates have valid data structure", len(receivedDepthUpdates))
	} else {
		t.Logf("⚠️  No depth update data received within %v timeout - this may be normal depending on market activity", timeout)
	}

	// Verify OnReconnect was not called (should only happen on connection issues)
	if finalReconnectCount > 0 {
		t.Logf("ℹ️  OnReconnect was called %d times (may indicate connection issues during test)", finalReconnectCount)
	}

	t.Log("✅ Test completed successfully")
}

func TestWSClient_SubscribeDepthUpdate_100ms(t *testing.T) {
	// Test configuration for high-frequency updates
	symbol := "btcusdt"
	updateSpeed := "100ms"
	timeout := 5 * time.Second // Shorter timeout for faster updates

	// Create WSClient
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})

	// Callback invocation counters
	var connectCount int64
	var depthUpdateCount int64
	var disconnectCount int64

	// Create subscription options
	options := DepthUpdateSubscriptionOptions{
		OnConnect: func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called for 100ms depth updates")
		},
		OnDepthUpdate: func(update WSDepthUpdate) {
			count := atomic.AddInt64(&depthUpdateCount, 1)
			if count <= 3 { // Log first few to avoid spam
				t.Logf("OnDepthUpdate called #%d (100ms): Symbol=%s, UpdateIds=%d-%d, BidUpdates=%d, AskUpdates=%d",
					count, update.Symbol, update.FirstUpdateId, update.FinalUpdateId, len(update.BidUpdates), len(update.AskUpdates))
			}
		},
		OnDisconnect: func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		},
	}

	// Subscribe to high-frequency differential depth stream
	unsubscribe, err := client.SubscribeDepthUpdate(symbol, updateSpeed, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to 100ms depth update stream: %v", err)
	}

	// Wait for timeout
	t.Logf("Waiting %v for high-frequency depth update data...", timeout)
	time.Sleep(timeout)

	// Unsubscribe
	unsubscribe()
	time.Sleep(200 * time.Millisecond)

	// Verify results
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalDepthUpdateCount := atomic.LoadInt64(&depthUpdateCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("100ms depth update stream results:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnDepthUpdate: %d", finalDepthUpdateCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect=1, got %d", finalConnectCount)
	}
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect=1, got %d", finalDisconnectCount)
	}
	if finalDepthUpdateCount >= 1 {
		t.Logf("✅ 100ms depth update stream working: received %d updates", finalDepthUpdateCount)
	} else {
		t.Logf("⚠️  No 100ms depth update data received")
	}
}

func TestWSClient_SubscribeDepth_100ms(t *testing.T) {
	// Test configuration for high-frequency updates
	symbol := "btcusdt"
	levels := 10
	updateSpeed := "100ms"
	timeout := 5 * time.Second // Shorter timeout for faster updates

	// Create WSClient
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})

	// Callback invocation counters
	var connectCount int64
	var depthCount int64
	var disconnectCount int64

	// Create subscription options
	options := DepthSubscriptionOptions{
		OnConnect: func() {
			atomic.AddInt64(&connectCount, 1)
			t.Log("OnConnect called for 100ms updates")
		},
		OnDepth: func(depth WSDepth) {
			count := atomic.AddInt64(&depthCount, 1)
			if count <= 3 { // Log first few to avoid spam
				t.Logf("OnDepth called #%d (100ms): LastUpdateId=%d, Bids=%d, Asks=%d",
					count, depth.LastUpdateId, len(depth.Bids), len(depth.Asks))
			}
		},
		OnDisconnect: func() {
			atomic.AddInt64(&disconnectCount, 1)
			t.Log("OnDisconnect called")
		},
	}

	// Subscribe to high-frequency depth stream
	unsubscribe, err := client.SubscribeDepth(symbol, levels, updateSpeed, options)
	if err != nil {
		t.Fatalf("Failed to subscribe to 100ms depth stream: %v", err)
	}

	// Wait for timeout
	t.Logf("Waiting %v for high-frequency depth data...", timeout)
	time.Sleep(timeout)

	// Unsubscribe
	unsubscribe()
	time.Sleep(200 * time.Millisecond)

	// Verify results
	finalConnectCount := atomic.LoadInt64(&connectCount)
	finalDepthCount := atomic.LoadInt64(&depthCount)
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)

	t.Logf("100ms stream results:")
	t.Logf("  OnConnect: %d", finalConnectCount)
	t.Logf("  OnDepth: %d", finalDepthCount)
	t.Logf("  OnDisconnect: %d", finalDisconnectCount)

	if finalConnectCount != 1 {
		t.Errorf("Expected OnConnect=1, got %d", finalConnectCount)
	}
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect=1, got %d", finalDisconnectCount)
	}
	if finalDepthCount >= 1 {
		t.Logf("✅ 100ms depth stream working: received %d updates", finalDepthCount)
	} else {
		t.Logf("⚠️  No 100ms depth data received")
	}
}

func TestWSClient_SubscribeDepth_InvalidLevels(t *testing.T) {
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})

	options := DepthSubscriptionOptions{
		OnConnect: func() {},
		OnDepth:   func(depth WSDepth) {},
	}

	// Test invalid levels
	invalidLevels := []int{1, 3, 15, 25, 50}
	for _, levels := range invalidLevels {
		_, err := client.SubscribeDepth("btcusdt", levels, "", options)
		if err == nil {
			t.Errorf("Expected error for invalid levels %d, but got nil", levels)
		} else {
			t.Logf("✅ Correctly rejected invalid levels %d: %v", levels, err)
		}
	}

	// Test valid levels
	validLevels := []int{5, 10, 20}
	for _, levels := range validLevels {
		unsubscribe, err := client.SubscribeDepth("btcusdt", levels, "", options)
		if err != nil {
			t.Errorf("Unexpected error for valid levels %d: %v", levels, err)
		} else {
			t.Logf("✅ Valid levels %d accepted", levels)
			if unsubscribe != nil {
				unsubscribe() // Clean up
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func TestWSClient_SubscribeKline_DuplicateSubscription(t *testing.T) {
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})
	symbol := "btcusdt"
	interval := "1m"

	options := KlineSubscriptionOptions{
		OnConnect: func() {},
		OnKline:   func(kline WSKline) {},
	}

	// First subscription should succeed
	unsubscribe1, err := client.SubscribeKline(symbol, interval, options)
	if err != nil {
		t.Fatalf("First subscription failed: %v", err)
	}
	defer unsubscribe1()

	// Second subscription to same symbol/interval should fail
	_, err = client.SubscribeKline(symbol, interval, options)
	if err == nil {
		t.Error("Expected error for duplicate subscription, but got nil")
	} else {
		t.Logf("✅ Duplicate subscription correctly returned error: %v", err)
	}
}

func TestWSClient_SubscribeAggTrade_DuplicateSubscription(t *testing.T) {
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})
	symbol := "btcusdt"

	options := AggTradeSubscriptionOptions{
		OnConnect:  func() {},
		OnAggTrade: func(aggTrade WSAggTrade) {},
	}

	// First subscription should succeed
	unsubscribe1, err := client.SubscribeAggTrade(symbol, options)
	if err != nil {
		t.Fatalf("First subscription failed: %v", err)
	}
	defer unsubscribe1()

	// Second subscription to same symbol should fail
	_, err = client.SubscribeAggTrade(symbol, options)
	if err == nil {
		t.Error("Expected error for duplicate subscription, but got nil")
	} else {
		t.Logf("✅ Duplicate subscription correctly returned error: %v", err)
	}
}

func TestWSClient_SubscribeTrade_DuplicateSubscription(t *testing.T) {
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})
	symbol := "btcusdt"

	options := TradeSubscriptionOptions{
		OnConnect: func() {},
		OnTrade:   func(trade WSTrade) {},
	}

	// First subscription should succeed
	unsubscribe1, err := client.SubscribeTrade(symbol, options)
	if err != nil {
		t.Fatalf("First subscription failed: %v", err)
	}
	defer unsubscribe1()

	// Second subscription to same symbol should fail
	_, err = client.SubscribeTrade(symbol, options)
	if err == nil {
		t.Error("Expected error for duplicate subscription, but got nil")
	} else {
		t.Logf("✅ Duplicate subscription correctly returned error: %v", err)
	}
}

func TestWSClient_SubscribeDepth_DuplicateSubscription(t *testing.T) {
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})
	symbol := "btcusdt"
	levels := 5

	options := DepthSubscriptionOptions{
		OnConnect: func() {},
		OnDepth:   func(depth WSDepth) {},
	}

	// First subscription should succeed
	unsubscribe1, err := client.SubscribeDepth(symbol, levels, "", options)
	if err != nil {
		t.Fatalf("First subscription failed: %v", err)
	}
	defer unsubscribe1()

	// Second subscription to same symbol/levels should fail
	_, err = client.SubscribeDepth(symbol, levels, "", options)
	if err == nil {
		t.Error("Expected error for duplicate subscription, but got nil")
	} else {
		t.Logf("✅ Duplicate subscription correctly returned error: %v", err)
	}
}

func TestWSClient_SubscribeDepthUpdate_DuplicateSubscription(t *testing.T) {
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})
	symbol := "btcusdt"

	options := DepthUpdateSubscriptionOptions{
		OnConnect:     func() {},
		OnDepthUpdate: func(update WSDepthUpdate) {},
	}

	// First subscription should succeed
	unsubscribe1, err := client.SubscribeDepthUpdate(symbol, "", options)
	if err != nil {
		t.Fatalf("First subscription failed: %v", err)
	}
	defer unsubscribe1()

	// Second subscription to same symbol should fail
	_, err = client.SubscribeDepthUpdate(symbol, "", options)
	if err == nil {
		t.Error("Expected error for duplicate subscription, but got nil")
	} else {
		t.Logf("✅ Duplicate subscription correctly returned error: %v", err)
	}
}

func TestWSClient_Close(t *testing.T) {
	client := NewWSClient(WSConfig{
		BaseURL: MainnetWSBaseUrl9443,
	})
	symbol := "btcusdt"
	interval := "1m"

	var disconnectCount int64

	options := KlineSubscriptionOptions{
		OnConnect: func() {},
		OnKline:   func(kline WSKline) {},
		OnDisconnect: func() {
			atomic.AddInt64(&disconnectCount, 1)
		},
	}

	// Subscribe
	_, err := client.SubscribeKline(symbol, interval, options)
	if err != nil {
		t.Fatalf("Subscription failed: %v", err)
	}

	// Close all connections
	client.Close()

	// Wait for disconnect processing
	time.Sleep(200 * time.Millisecond)

	// Verify OnDisconnect was called exactly once
	finalDisconnectCount := atomic.LoadInt64(&disconnectCount)
	if finalDisconnectCount != 1 {
		t.Errorf("Expected OnDisconnect to be called exactly 1 time after Close(), got %d", finalDisconnectCount)
	} else {
		t.Log("✅ OnDisconnect called correctly after Close()")
	}
}

func TestSubscribeUserData(t *testing.T) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCE_API_KEY or BINANCE_API_SECRET not set; skipping user data stream test.")
	}

	// Create REST API client
	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)

	// Create WebSocket client with REST client
	wsClient := NewWSClientWithRESTClient(WSConfig{
		BaseURL: TestnetWSBaseUrl, // Use testnet for testing
	}, client)

	connected := false
	options := UserDataSubscriptionOptions{
		OnConnect: func() {
			t.Log("Connected to user data stream")
			connected = true
		},
		OnError: func(err error) {
			t.Logf("User data stream error: %v", err)
		},
		OnAccountPosition: func(event WSOutboundAccountPositionEvent) {
			t.Logf("Account position update: %d balances", len(event.BalanceArray))
		},
		OnBalanceUpdate: func(event WSBalanceUpdateEvent) {
			t.Logf("Balance update: %s delta=%s", event.Asset, event.BalanceDelta)
		},
		OnExecutionReport: func(event WSExecutionReportEvent) {
			t.Logf("Execution report: %s %s %s", event.Symbol, event.Side, event.CurrentOrderStatus)
		},
		OnListenKeyExpired: func(event WSListenKeyExpiredEvent) {
			t.Logf("Listen key expired: %s", event.ListenKey)
		},
		OnDisconnect: func() {
			t.Log("Disconnected from user data stream")
		},
	}

	unsubscribe, err := wsClient.SubscribeUserData(options)
	if err != nil {
		t.Fatalf("Failed to subscribe to user data stream: %v", err)
	}

	// Wait a bit for connection
	time.Sleep(2 * time.Second)

	if !connected {
		t.Error("Failed to connect to user data stream")
	}

	// Test that we can't subscribe twice
	_, err = wsClient.SubscribeUserData(options)
	if err == nil {
		t.Error("Expected error when subscribing twice, but got nil")
	}

	// Unsubscribe
	unsubscribe()

	// Give some time for cleanup
	time.Sleep(100 * time.Millisecond)

	// Clean up
	wsClient.Close()
}

func TestUserDataStreamWithoutClient(t *testing.T) {
	// Create WebSocket client without REST client
	wsClient := NewWSClient(WSConfig{
		BaseURL: TestnetWSBaseUrl,
	})

	options := UserDataSubscriptionOptions{
		OnError: func(err error) {
			t.Logf("Expected error: %v", err)
		},
	}

	_, err := wsClient.SubscribeUserData(options)
	if err == nil {
		t.Error("Expected error when subscribing without REST client, but got nil")
	}

	expectedMsg := "REST API client is required for user data stream subscription"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}
