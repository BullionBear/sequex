package binance

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"
)

func TestBuildStreamNames(t *testing.T) {
	tests := []struct {
		name     string
		function func() string
		expected string
	}{
		{
			name:     "Kline stream name",
			function: func() string { return BuildKlineStreamName("BTCUSDT", Interval1m) },
			expected: "BTCUSDT@kline_1m",
		},
		{
			name:     "Ticker stream name",
			function: func() string { return BuildTickerStreamName("ETHUSDT") },
			expected: "ETHUSDT@ticker",
		},
		{
			name:     "Trade stream name",
			function: func() string { return BuildTradeStreamName("ADAUSDT") },
			expected: "ADAUSDT@trade",
		},
		{
			name:     "Book ticker stream name",
			function: func() string { return BuildBookTickerStreamName("BNBUSDT") },
			expected: "BNBUSDT@bookTicker",
		},
		{
			name:     "Aggregate trade stream name",
			function: func() string { return BuildAggTradeStreamName("LTCUSDT") },
			expected: "LTCUSDT@aggTrade",
		},
		{
			name:     "Depth stream name with levels",
			function: func() string { return BuildDepthStreamName("SOLUSDT", 5) },
			expected: "SOLUSDT@depth5",
		},
		{
			name:     "Depth stream name without levels",
			function: func() string { return BuildDepthStreamName("SOLUSDT", 0) },
			expected: "SOLUSDT@depth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestWSKlineData_TimeMethods(t *testing.T) {
	kline := &WSKlineData{
		OpenTime:  1640995200000, // 2022-01-01 00:00:00 UTC
		CloseTime: 1640995260000, // 2022-01-01 00:01:00 UTC
	}

	openTime := kline.GetOpenTime()
	closeTime := kline.GetCloseTime()

	expectedOpen := time.Unix(1640995200, 0)
	expectedClose := time.Unix(1640995260, 0)

	if !openTime.Equal(expectedOpen) {
		t.Errorf("expected open time %v, got %v", expectedOpen, openTime)
	}

	if !closeTime.Equal(expectedClose) {
		t.Errorf("expected close time %v, got %v", expectedClose, closeTime)
	}
}

func TestWSError_Error(t *testing.T) {
	err := &WSError{
		Code: -1003,
		Msg:  "Too many requests",
	}

	expected := "Too many requests"
	if err.Error() != expected {
		t.Errorf("expected error message %s, got %s", expected, err.Error())
	}
}

func TestWSClient_NewWSClient(t *testing.T) {
	client := NewWSClient(nil)

	if client == nil {
		t.Fatal("expected client, got nil")
	}

	if client.config == nil {
		t.Error("expected config, got nil")
	}

	if client.conn == nil {
		t.Error("expected connection, got nil")
	}

	if client.subscriptions == nil {
		t.Error("expected subscriptions map, got nil")
	}

	if len(client.GetSubscriptions()) != 0 {
		t.Error("expected empty subscriptions, got some")
	}
}

func TestWSClient_SubscriptionManagement(t *testing.T) {
	client := NewWSClient(TestnetConfig())

	// Test kline subscription validation
	err := client.SubscribeKline([]string{}, Interval1m)
	if err == nil {
		t.Error("expected error for empty symbols")
	}

	err = client.SubscribeKline([]string{"BTCUSDT"}, "invalid")
	if err == nil {
		t.Error("expected error for invalid interval")
	}

	// Test ticker subscription validation
	err = client.SubscribeTicker([]string{})
	if err == nil {
		t.Error("expected error for empty symbols")
	}

	// Test trade subscription validation
	err = client.SubscribeTrade([]string{})
	if err == nil {
		t.Error("expected error for empty symbols")
	}
}

func TestWSClient_EventHandlers(t *testing.T) {
	client := NewWSClient(TestnetConfig())

	var klineReceived bool
	var tickerReceived bool
	var tradeReceived bool

	// Add event handlers
	client.OnKline(func(event *WSKlineEvent) {
		klineReceived = true
	})

	client.OnTicker(func(event *WSTickerEvent) {
		tickerReceived = true
	})

	client.OnTrade(func(event *WSTradeEvent) {
		tradeReceived = true
	})

	// Simulate kline event
	klineData := WSKlineEvent{
		EventType: "kline",
		EventTime: time.Now().UnixNano() / int64(time.Millisecond),
		Symbol:    "BTCUSDT",
		Kline: WSKlineData{
			Symbol:     "BTCUSDT",
			OpenTime:   time.Now().UnixNano() / int64(time.Millisecond),
			CloseTime:  time.Now().UnixNano() / int64(time.Millisecond),
			Interval:   Interval1m,
			OpenPrice:  "45000.00",
			ClosePrice: "45100.00",
			HighPrice:  "45200.00",
			LowPrice:   "44900.00",
		},
	}

	data, _ := json.Marshal(klineData)
	client.handleKlineEvent(data)

	// Simulate ticker event
	tickerData := WSTickerEvent{
		EventType: "24hrTicker",
		EventTime: time.Now().UnixNano() / int64(time.Millisecond),
		Symbol:    "BTCUSDT",
		LastPrice: "45000.00",
	}

	data, _ = json.Marshal(tickerData)
	client.handleTickerEvent(data)

	// Simulate trade event
	tradeData := WSTradeEvent{
		EventType: "trade",
		EventTime: time.Now().UnixNano() / int64(time.Millisecond),
		Symbol:    "BTCUSDT",
		Price:     "45000.00",
		Quantity:  "0.001",
	}

	data, _ = json.Marshal(tradeData)
	client.handleTradeEvent(data)

	// Wait for handlers to be called
	time.Sleep(100 * time.Millisecond)

	// Verify handlers were called
	if !klineReceived {
		t.Error("kline handler was not called")
	}

	if !tickerReceived {
		t.Error("ticker handler was not called")
	}

	if !tradeReceived {
		t.Error("trade handler was not called")
	}
}

func TestWSClient_MessageParsing(t *testing.T) {
	client := NewWSClient(TestnetConfig())

	var receivedEvent *WSKlineEvent
	var mu sync.Mutex

	client.OnKline(func(event *WSKlineEvent) {
		mu.Lock()
		defer mu.Unlock()
		receivedEvent = event
	})

	// Test direct event parsing
	klineJSON := `{
		"e": "kline",
		"E": 1640995200000,
		"s": "BTCUSDT",
		"k": {
			"t": 1640995200000,
			"T": 1640995260000,
			"s": "BTCUSDT",
			"i": "1m",
			"f": 100,
			"L": 200,
			"o": "45000.00",
			"c": "45100.00",
			"h": "45200.00",
			"l": "44900.00",
			"v": "1.000000",
			"n": 50,
			"x": true,
			"q": "45050.00",
			"V": "0.500000",
			"Q": "22525.00",
			"B": "0"
		}
	}`

	client.handleKlineEvent([]byte(klineJSON))

	// Wait for handler
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if receivedEvent == nil {
		t.Fatal("no kline event received")
	}

	if receivedEvent.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", receivedEvent.Symbol)
	}

	if receivedEvent.Kline.OpenPrice != "45000.00" {
		t.Errorf("expected open price 45000.00, got %s", receivedEvent.Kline.OpenPrice)
	}

	if receivedEvent.Kline.Interval != "1m" {
		t.Errorf("expected interval 1m, got %s", receivedEvent.Kline.Interval)
	}
}

func TestWSClient_StreamMessageParsing(t *testing.T) {
	client := NewWSClient(TestnetConfig())

	var receivedEvent *WSKlineEvent
	var mu sync.Mutex

	client.OnKline(func(event *WSKlineEvent) {
		mu.Lock()
		defer mu.Unlock()
		receivedEvent = event
	})

	// Test stream message parsing
	streamKlineJSON := `{
		"e": "kline",
		"E": 1640995200000,
		"s": "BTCUSDT",
		"k": {
			"t": 1640995200000,
			"T": 1640995260000,
			"s": "BTCUSDT",
			"i": "1m",
			"f": 100,
			"L": 200,
			"o": "46000.00",
			"c": "46100.00",
			"h": "46200.00",
			"l": "45900.00",
			"v": "2.000000",
			"n": 75,
			"x": false,
			"q": "92050.00",
			"V": "1.000000",
			"Q": "46050.00",
			"B": "0"
		}
	}`

	client.handleKlineEvent([]byte(streamKlineJSON))

	// Wait for handler
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if receivedEvent == nil {
		t.Fatal("no kline event received")
	}

	if receivedEvent.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", receivedEvent.Symbol)
	}

	if receivedEvent.Kline.OpenPrice != "46000.00" {
		t.Errorf("expected open price 46000.00, got %s", receivedEvent.Kline.OpenPrice)
	}
}

func TestWSClient_ResponseParsing(t *testing.T) {
	client := NewWSClient(TestnetConfig())

	// Test successful response
	successJSON := `{
		"result": null,
		"id": 1
	}`

	client.handleMessage([]byte(successJSON))

	// Test error response
	errorJSON := `{
		"error": {
			"code": -2011,
			"msg": "Unknown symbol"
		},
		"id": 2
	}`

	client.handleMessage([]byte(errorJSON))

	// No assertions here since we're just testing that parsing doesn't crash
}

// Integration test with real WebSocket connection (if credentials available)
// DISABLED: Testnet WebSocket has connectivity issues
func testWSClient_RealConnection_WithTestCredentials(t *testing.T) {
	if !HasTestCredentials() {
		t.Skip("Skipping real WebSocket test: no test credentials available")
	}

	config, err := LoadTestConfig()
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	client := NewWSClient(config)

	var klineCount int
	var mu sync.Mutex

	// Set up kline handler
	client.OnKline(func(event *WSKlineEvent) {
		mu.Lock()
		defer mu.Unlock()
		klineCount++
		t.Logf("Received kline: %s %s O:%s C:%s H:%s L:%s V:%s",
			event.Symbol,
			event.Kline.Interval,
			event.Kline.OpenPrice,
			event.Kline.ClosePrice,
			event.Kline.HighPrice,
			event.Kline.LowPrice,
			event.Kline.BaseAssetVolume)
	})

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Wait for connection
	time.Sleep(1 * time.Second)

	if !client.IsConnected() {
		t.Fatal("WebSocket should be connected")
	}

	// Subscribe to BTCUSDT 1m klines
	err = client.SubscribeKline([]string{"BTCUSDT"}, Interval1m)
	if err != nil {
		t.Fatalf("Failed to subscribe to klines: %v", err)
	}

	// Wait for some data
	t.Log("Waiting for kline data...")
	time.Sleep(10 * time.Second)

	// Check subscriptions
	subscriptions := client.GetSubscriptions()
	expectedSub := "BTCUSDT@kline_1m"
	found := false
	for _, sub := range subscriptions {
		if sub == expectedSub {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected subscription %s not found in %v", expectedSub, subscriptions)
	}

	// Check that we received some data
	mu.Lock()
	count := klineCount
	mu.Unlock()

	if count == 0 {
		t.Error("No kline data received")
	} else {
		t.Logf("Received %d kline events", count)
	}

	// Unsubscribe
	err = client.UnsubscribeKline([]string{"BTCUSDT"}, Interval1m)
	if err != nil {
		t.Errorf("Failed to unsubscribe: %v", err)
	}

	// Wait a moment
	time.Sleep(2 * time.Second)

	// Check subscriptions are cleared
	subscriptions = client.GetSubscriptions()
	for _, sub := range subscriptions {
		if sub == expectedSub {
			t.Errorf("Subscription %s should have been removed", expectedSub)
		}
	}

	// Disconnect
	err = client.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect: %v", err)
	}

	// Wait for disconnect
	time.Sleep(1 * time.Second)

	if client.IsConnected() {
		t.Error("WebSocket should be disconnected")
	}
}

// DISABLED: Testnet WebSocket has connectivity issues
func testWSClient_MultipleEventTypes(t *testing.T) {
	if !HasTestCredentials() {
		t.Skip("Skipping multi-event test: no test credentials available")
	}

	config, err := LoadTestConfig()
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	client := NewWSClient(config)

	var klineCount, tickerCount, tradeCount int
	var mu sync.Mutex

	// Set up event handlers
	client.OnKline(func(event *WSKlineEvent) {
		mu.Lock()
		defer mu.Unlock()
		klineCount++
		t.Logf("Kline: %s %s", event.Symbol, event.Kline.ClosePrice)
	})

	client.OnTicker(func(event *WSTickerEvent) {
		mu.Lock()
		defer mu.Unlock()
		tickerCount++
		t.Logf("Ticker: %s %s", event.Symbol, event.LastPrice)
	})

	client.OnTrade(func(event *WSTradeEvent) {
		mu.Lock()
		defer mu.Unlock()
		tradeCount++
		t.Logf("Trade: %s %s@%s", event.Symbol, event.Quantity, event.Price)
	})

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Subscribe to different event types
	err = client.SubscribeKline([]string{"BTCUSDT"}, Interval1m)
	if err != nil {
		t.Fatalf("Failed to subscribe to klines: %v", err)
	}

	err = client.SubscribeTicker([]string{"BTCUSDT"})
	if err != nil {
		t.Fatalf("Failed to subscribe to ticker: %v", err)
	}

	err = client.SubscribeTrade([]string{"BTCUSDT"})
	if err != nil {
		t.Fatalf("Failed to subscribe to trades: %v", err)
	}

	// Wait for data
	t.Log("Waiting for multi-stream data...")
	time.Sleep(15 * time.Second)

	// Check that we received data from all streams
	mu.Lock()
	kCount := klineCount
	tCount := tickerCount
	trCount := tradeCount
	mu.Unlock()

	t.Logf("Received: %d klines, %d tickers, %d trades", kCount, tCount, trCount)

	if kCount == 0 {
		t.Error("No kline data received")
	}

	if tCount == 0 {
		t.Error("No ticker data received")
	}

	if trCount == 0 {
		t.Error("No trade data received")
	}

	// Clean up
	client.Disconnect()
}
