package binance

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"
)

func TestUserDataStreamClient_New(t *testing.T) {
	client := NewUserDataStreamClient(nil)

	if client == nil {
		t.Fatal("expected client, got nil")
	}

	if client.config == nil {
		t.Error("expected config, got nil")
	}

	if client.client == nil {
		t.Error("expected REST client, got nil")
	}

	if client.keepAliveStop == nil {
		t.Error("expected keepAliveStop channel, got nil")
	}

	if client.IsConnected() {
		t.Error("client should not be connected initially")
	}
}

func TestUserDataStreamClient_EventHandlers(t *testing.T) {
	client := NewUserDataStreamClient(TestnetConfig())

	var accountUpdateReceived bool
	var balanceUpdateReceived bool
	var executionReportReceived bool
	var listStatusReceived bool

	// Add event handlers
	client.OnAccountUpdate(func(event *WSAccountUpdate) {
		accountUpdateReceived = true
	})

	client.OnBalanceUpdate(func(event *WSBalanceUpdate) {
		balanceUpdateReceived = true
	})

	client.OnExecutionReport(func(event *WSExecutionReport) {
		executionReportReceived = true
	})

	client.OnListStatus(func(event *WSListStatus) {
		listStatusReceived = true
	})

	// Simulate account update event
	accountData := WSAccountUpdate{
		EventType:    WSEventAccountUpdate,
		EventTime:    time.Now().UnixNano() / int64(time.Millisecond),
		LastUpdateID: 12345,
		Balances: []WSAccountBalance{
			{Asset: "BTC", Free: "1.00000000", Locked: "0.00000000"},
			{Asset: "USDT", Free: "1000.00000000", Locked: "0.00000000"},
		},
	}

	data, _ := json.Marshal(accountData)
	client.handleAccountUpdate(data)

	// Simulate balance update event
	balanceData := WSBalanceUpdate{
		EventType:    WSEventBalanceUpdate,
		EventTime:    time.Now().UnixNano() / int64(time.Millisecond),
		Asset:        "BTC",
		BalanceDelta: "0.01000000",
		ClearTime:    time.Now().UnixNano() / int64(time.Millisecond),
	}

	data, _ = json.Marshal(balanceData)
	client.handleBalanceUpdate(data)

	// Simulate execution report event
	executionData := WSExecutionReport{
		EventType:          WSEventExecutionReport,
		EventTime:          time.Now().UnixNano() / int64(time.Millisecond),
		Symbol:             "BTCUSDT",
		ClientOrderID:      "test123",
		Side:               "BUY",
		OrderType:          "LIMIT",
		CurrentOrderStatus: "NEW",
		OrderID:            12345,
		TransactionTime:    time.Now().UnixNano() / int64(time.Millisecond),
		OrderCreationTime:  time.Now().UnixNano() / int64(time.Millisecond),
	}

	data, _ = json.Marshal(executionData)
	client.handleExecutionReport(data)

	// Simulate list status event
	listData := WSListStatus{
		EventType:       WSEventListStatus,
		EventTime:       time.Now().UnixNano() / int64(time.Millisecond),
		Symbol:          "BTCUSDT",
		OrderListID:     12345,
		TransactionTime: time.Now().UnixNano() / int64(time.Millisecond),
	}

	data, _ = json.Marshal(listData)
	client.handleListStatus(data)

	// Wait for handlers to be called
	time.Sleep(100 * time.Millisecond)

	// Verify handlers were called
	if !accountUpdateReceived {
		t.Error("account update handler was not called")
	}

	if !balanceUpdateReceived {
		t.Error("balance update handler was not called")
	}

	if !executionReportReceived {
		t.Error("execution report handler was not called")
	}

	if !listStatusReceived {
		t.Error("list status handler was not called")
	}
}

func TestWSExecutionReport_HelperMethods(t *testing.T) {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	tests := []struct {
		name            string
		executionType   string
		orderStatus     string
		expectedNew     bool
		expectedFilled  bool
		expectedCancel  bool
		expectedPartial bool
		expectedTrade   bool
	}{
		{
			name:          "New Order",
			executionType: "NEW",
			orderStatus:   "NEW",
			expectedNew:   true,
		},
		{
			name:           "Filled Order",
			executionType:  "TRADE",
			orderStatus:    "FILLED",
			expectedFilled: true,
			expectedTrade:  true,
		},
		{
			name:           "Canceled Order",
			executionType:  "CANCELED",
			orderStatus:    "CANCELED",
			expectedCancel: true,
		},
		{
			name:            "Partially Filled",
			executionType:   "TRADE",
			orderStatus:     "PARTIALLY_FILLED",
			expectedPartial: true,
			expectedTrade:   true,
		},
		{
			name:           "Trade Execution",
			executionType:  "TRADE",
			orderStatus:    "FILLED",
			expectedTrade:  true,
			expectedFilled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &WSExecutionReport{
				CurrentExecutionType: tt.executionType,
				CurrentOrderStatus:   tt.orderStatus,
				TransactionTime:      now,
				OrderCreationTime:    now,
			}

			if report.IsNewOrder() != tt.expectedNew {
				t.Errorf("IsNewOrder() = %v, expected %v", report.IsNewOrder(), tt.expectedNew)
			}

			if report.IsFilled() != tt.expectedFilled {
				t.Errorf("IsFilled() = %v, expected %v", report.IsFilled(), tt.expectedFilled)
			}

			if report.IsCanceled() != tt.expectedCancel {
				t.Errorf("IsCanceled() = %v, expected %v", report.IsCanceled(), tt.expectedCancel)
			}

			if report.IsPartiallyFilled() != tt.expectedPartial {
				t.Errorf("IsPartiallyFilled() = %v, expected %v", report.IsPartiallyFilled(), tt.expectedPartial)
			}

			if report.IsTrade() != tt.expectedTrade {
				t.Errorf("IsTrade() = %v, expected %v", report.IsTrade(), tt.expectedTrade)
			}

			// Test time methods
			expectedTime := time.Unix(0, now*int64(time.Millisecond))
			if !report.GetTransactionTime().Equal(expectedTime) {
				t.Errorf("GetTransactionTime() = %v, expected %v", report.GetTransactionTime(), expectedTime)
			}

			if !report.GetOrderCreationTime().Equal(expectedTime) {
				t.Errorf("GetOrderCreationTime() = %v, expected %v", report.GetOrderCreationTime(), expectedTime)
			}
		})
	}
}

func TestUserDataStreamClient_MessageParsing(t *testing.T) {
	client := NewUserDataStreamClient(TestnetConfig())

	var receivedEvent *WSExecutionReport
	var mu sync.Mutex

	client.OnExecutionReport(func(event *WSExecutionReport) {
		mu.Lock()
		defer mu.Unlock()
		receivedEvent = event
	})

	// Test execution report parsing
	executionJSON := `{
		"e": "executionReport",
		"E": 1499405658658,
		"s": "ETHBTC",
		"c": "mUvoqJxFIILMdfAW5iGSOW",
		"S": "BUY",
		"o": "LIMIT",
		"f": "GTC",
		"q": "1.00000000",
		"p": "0.10264410",
		"P": "0.00000000",
		"F": "0.00000000",
		"g": -1,
		"C": "",
		"x": "NEW",
		"X": "NEW",
		"r": "NONE",
		"i": 4293153,
		"l": "0.00000000",
		"z": "0.00000000",
		"L": "0.00000000",
		"n": "0",
		"N": null,
		"T": 1499405658657,
		"t": -1,
		"I": 8641984,
		"w": true,
		"m": false,
		"M": false,
		"O": 1499405658657,
		"Z": "0.00000000",
		"Y": "0.00000000",
		"Q": "0.00000000"
	}`

	client.handleExecutionReport([]byte(executionJSON))

	// Wait for handler
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if receivedEvent == nil {
		t.Fatal("no execution report event received")
	}

	if receivedEvent.Symbol != "ETHBTC" {
		t.Errorf("expected symbol ETHBTC, got %s", receivedEvent.Symbol)
	}

	if receivedEvent.Side != "BUY" {
		t.Errorf("expected side BUY, got %s", receivedEvent.Side)
	}

	if receivedEvent.OrderType != "LIMIT" {
		t.Errorf("expected order type LIMIT, got %s", receivedEvent.OrderType)
	}

	if receivedEvent.CurrentOrderStatus != "NEW" {
		t.Errorf("expected order status NEW, got %s", receivedEvent.CurrentOrderStatus)
	}

	if !receivedEvent.IsNewOrder() {
		t.Error("should be a new order")
	}
}

// Integration test for user data stream (requires real credentials)
// DISABLED: Testnet WebSocket has connectivity issues
func testUserDataStreamClient_RealConnection_WithTestCredentials(t *testing.T) {
	if !HasTestCredentials() {
		t.Skip("Skipping real user data stream test: no test credentials available")
	}

	config, err := LoadTestConfig()
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	client := NewUserDataStreamClient(config)

	var accountUpdateCount int
	var executionReportCount int
	var mu sync.Mutex

	// Set up event handlers
	client.OnAccountUpdate(func(event *WSAccountUpdate) {
		mu.Lock()
		defer mu.Unlock()
		accountUpdateCount++
		t.Logf("Account Update: %d balances, LastUpdateID: %d",
			len(event.Balances), event.LastUpdateID)

		for _, balance := range event.Balances {
			if balance.Free != "0.00000000" || balance.Locked != "0.00000000" {
				t.Logf("  %s: Free=%s, Locked=%s",
					balance.Asset, balance.Free, balance.Locked)
			}
		}
	})

	client.OnExecutionReport(func(event *WSExecutionReport) {
		mu.Lock()
		defer mu.Unlock()
		executionReportCount++
		t.Logf("Execution Report: %s %s %s OrderID:%d Status:%s",
			event.Symbol, event.Side, event.OrderType, event.OrderID, event.CurrentOrderStatus)

		if event.IsNewOrder() {
			t.Logf("  New order created")
		}
		if event.IsTrade() {
			t.Logf("  Trade executed: %s @ %s", event.LastExecutedQuantity, event.LastExecutedPrice)
		}
		if event.IsCanceled() {
			t.Logf("  Order canceled")
		}
		if event.IsFilled() {
			t.Logf("  Order filled")
		}
	})

	client.OnBalanceUpdate(func(event *WSBalanceUpdate) {
		t.Logf("Balance Update: %s Delta=%s", event.Asset, event.BalanceDelta)
	})

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	if !client.IsConnected() {
		t.Fatal("User data stream should be connected")
	}

	t.Log("✅ User data stream connected successfully")

	// Wait for some events (or until we get initial account update)
	t.Log("📡 Waiting for user data stream events...")
	time.Sleep(5 * time.Second)

	// Create a small test order to generate execution reports
	restClient := NewClient(config)

	// Place a limit order at a low price (won't fill immediately)
	orderReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        "BUY",
		Type:        "LIMIT",
		TimeInForce: "GTC",
		Quantity:    "0.001",
		Price:       "50000.00", // Low price, won't fill
	}

	t.Log("📊 Placing test order at low price to generate execution reports")

	orderResp, err := restClient.PlaceOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("Failed to place order: %v", err)
	}

	t.Logf("✅ Order placed: ID=%d", orderResp.OrderId)

	// Wait for execution report
	time.Sleep(3 * time.Second)

	// Cancel the order to generate more events
	_, err = restClient.CancelOrder(ctx, "BTCUSDT", orderResp.OrderId)
	if err != nil {
		t.Logf("Failed to cancel order (may have filled): %v", err)
	} else {
		t.Log("📝 Order canceled")
	}

	// Wait for cancel event
	time.Sleep(2 * time.Second)

	// Check that we received events
	mu.Lock()
	aCount := accountUpdateCount
	eCount := executionReportCount
	mu.Unlock()

	t.Logf("📈 Received %d account updates, %d execution reports", aCount, eCount)

	if aCount == 0 {
		t.Error("❌ No account updates received")
	} else {
		t.Logf("✅ Received %d account updates", aCount)
	}

	if eCount == 0 {
		t.Error("❌ No execution reports received")
	} else {
		t.Logf("✅ Received %d execution reports", eCount)
	}

	// Disconnect
	err = client.Disconnect()
	if err != nil {
		t.Errorf("Failed to disconnect: %v", err)
	}

	if client.IsConnected() {
		t.Error("User data stream should be disconnected")
	}

	t.Log("🎉 User data stream test completed successfully!")
}

func TestClient_UserDataStreamEndpoints_WithTestCredentials(t *testing.T) {
	if !HasTestCredentials() {
		t.Skip("Skipping user data stream endpoints test: no test credentials available")
	}

	config, err := LoadTestConfig()
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test creating user data stream
	streamResp, err := client.CreateUserDataStream(ctx)
	if err != nil {
		t.Fatalf("Failed to create user data stream: %v", err)
	}

	if streamResp.ListenKey == "" {
		t.Error("Expected listen key, got empty string")
	}

	t.Logf("✅ Created user data stream with listen key: %s...", streamResp.ListenKey[:8])

	// Test keeping alive user data stream
	err = client.KeepAliveUserDataStream(ctx, streamResp.ListenKey)
	if err != nil {
		t.Errorf("Failed to keep alive user data stream: %v", err)
	} else {
		t.Log("✅ Successfully kept user data stream alive")
	}

	// Test closing user data stream
	err = client.CloseUserDataStream(ctx, streamResp.ListenKey)
	if err != nil {
		t.Errorf("Failed to close user data stream: %v", err)
	} else {
		t.Log("✅ Successfully closed user data stream")
	}
}
