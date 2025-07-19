package binance

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestClient_GetAccount_WithTestCredentials(t *testing.T) {
	if !HasTestCredentials() {
		t.Skip("Skipping test: no test credentials available")
	}

	client, err := CreateTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	ctx := context.Background()
	account, err := client.GetAccount(ctx)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}

	// Basic assertions
	if account == nil {
		t.Fatal("Expected account response, got nil")
	}

	// Account should have some balances
	if len(account.Balances) == 0 {
		t.Error("Expected some balances, got none")
	}

	// Account should be able to trade (on testnet)
	if !account.CanTrade {
		t.Error("Expected account to be able to trade")
	}

	t.Logf("Account retrieved successfully with %d balances", len(account.Balances))

	// Log some balances for debugging
	for i, balance := range account.Balances {
		if i >= 5 { // Only log first 5 balances
			break
		}
		if balance.Free != "0.00000000" || balance.Locked != "0.00000000" {
			t.Logf("Balance: %s - Free: %s, Locked: %s", balance.Asset, balance.Free, balance.Locked)
		}
	}
}

func TestClient_PlaceOrder_LimitBuy_WithTestCredentials(t *testing.T) {
	if !HasTestCredentials() {
		t.Skip("Skipping test: no test credentials available")
	}

	client, err := CreateTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	ctx := context.Background()

	// First, get the current BTC price
	result, err := client.GetTickerPrice(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("Failed to get BTC price: %v", err)
	}

	if !result.IsSingle() {
		t.Fatal("Expected single ticker result")
	}
	ticker := result.GetSingle()
	currentPrice, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		t.Fatalf("Failed to parse current price: %v", err)
	}

	// Place a buy order 10% below current market price (unlikely to fill)
	orderPrice := currentPrice * 0.9
	orderPriceStr := strconv.FormatFloat(orderPrice, 'f', 2, 64)

	// Create a limit buy order for BTCUSDT with a realistic low price
	orderReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",       // Small quantity
		Price:       orderPriceStr, // 10% below market, unlikely to fill
	}

	orderResp, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("Failed to place order: %v", err)
	}

	// Basic assertions
	if orderResp == nil {
		t.Fatal("Expected order response, got nil")
	}

	if orderResp.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", orderResp.Symbol)
	}

	if orderResp.Side != SideBuy {
		t.Errorf("Expected side BUY, got %s", orderResp.Side)
	}

	if orderResp.Type != OrderTypeLimit {
		t.Errorf("Expected type LIMIT, got %s", orderResp.Type)
	}

	if orderResp.OrderId == 0 {
		t.Error("Expected order ID, got 0")
	}

	t.Logf("Order placed successfully: ID=%d, Symbol=%s, Side=%s, Status=%s",
		orderResp.OrderId, orderResp.Symbol, orderResp.Side, orderResp.Status)

	// Clean up: cancel the order
	t.Run("Cancel_Order", func(t *testing.T) {
		cancelResp, err := client.CancelOrder(ctx, orderResp.Symbol, orderResp.OrderId)
		if err != nil {
			t.Errorf("Failed to cancel order: %v", err)
			return
		}

		if cancelResp.OrderId != orderResp.OrderId {
			t.Errorf("Expected cancelled order ID %d, got %d", orderResp.OrderId, cancelResp.OrderId)
		}

		if cancelResp.Status != OrderStatusCanceled {
			t.Errorf("Expected status CANCELED, got %s", cancelResp.Status)
		}

		t.Logf("Order cancelled successfully: ID=%d, Status=%s", cancelResp.OrderId, cancelResp.Status)
	})
}

func TestClient_PlaceOrder_MarketBuy_WithTestCredentials(t *testing.T) {
	if !HasTestCredentials() {
		t.Skip("Skipping test: no test credentials available")
	}

	client, err := CreateTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	ctx := context.Background()

	// Create a small market buy order for BTCUSDT
	orderReq := &NewOrderRequest{
		Symbol:        "BTCUSDT",
		Side:          SideBuy,
		Type:          OrderTypeMarket,
		QuoteOrderQty: "10.00", // Buy $10 worth of BTC
	}

	orderResp, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("Failed to place market order: %v", err)
	}

	// Basic assertions
	if orderResp == nil {
		t.Fatal("Expected order response, got nil")
	}

	if orderResp.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", orderResp.Symbol)
	}

	if orderResp.Side != SideBuy {
		t.Errorf("Expected side BUY, got %s", orderResp.Side)
	}

	if orderResp.Type != OrderTypeMarket {
		t.Errorf("Expected type MARKET, got %s", orderResp.Type)
	}

	// Market orders should be filled immediately (or rejected if insufficient balance)
	if orderResp.Status != OrderStatusFilled && orderResp.Status != OrderStatusRejected {
		t.Errorf("Expected status FILLED or REJECTED for market order, got %s", orderResp.Status)
	}

	t.Logf("Market order result: ID=%d, Symbol=%s, Status=%s, ExecutedQty=%s",
		orderResp.OrderId, orderResp.Symbol, orderResp.Status, orderResp.ExecutedQty)

	if len(orderResp.Fills) > 0 {
		t.Logf("Order fills: %d fills", len(orderResp.Fills))
		for i, fill := range orderResp.Fills {
			t.Logf("Fill %d: Price=%s, Qty=%s, Commission=%s", i+1, fill.Price, fill.Qty, fill.Commission)
		}
	}
}

func TestClient_GetOrder_WithTestCredentials(t *testing.T) {
	if !HasTestCredentials() {
		t.Skip("Skipping test: no test credentials available")
	}

	client, err := CreateTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	ctx := context.Background()

	// First, get the current BTC price
	result, err := client.GetTickerPrice(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("Failed to get BTC price: %v", err)
	}

	if !result.IsSingle() {
		t.Fatal("Expected single ticker result")
	}
	ticker := result.GetSingle()
	currentPrice, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		t.Fatalf("Failed to parse current price: %v", err)
	}

	// Place a buy order 10% below current market price
	orderPrice := currentPrice * 0.9
	orderPriceStr := strconv.FormatFloat(orderPrice, 'f', 2, 64)

	// First, place an order
	orderReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		Price:       orderPriceStr, // 10% below market
	}

	orderResp, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("Failed to place order: %v", err)
	}

	// Wait a moment for the order to be processed
	time.Sleep(500 * time.Millisecond)

	// Now get the order status
	order, err := client.GetOrder(ctx, orderResp.Symbol, orderResp.OrderId)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}

	// Assertions
	if order.OrderId != orderResp.OrderId {
		t.Errorf("Expected order ID %d, got %d", orderResp.OrderId, order.OrderId)
	}

	if order.Symbol != orderResp.Symbol {
		t.Errorf("Expected symbol %s, got %s", orderResp.Symbol, order.Symbol)
	}

	if order.Status != OrderStatusNew {
		t.Errorf("Expected status NEW, got %s", order.Status)
	}

	t.Logf("Order retrieved: ID=%d, Symbol=%s, Status=%s, Price=%s, OrigQty=%s",
		order.OrderId, order.Symbol, order.Status, order.Price, order.OrigQty)

	// Clean up
	_, err = client.CancelOrder(ctx, order.Symbol, order.OrderId)
	if err != nil {
		t.Errorf("Failed to clean up order: %v", err)
	}
}

func TestClient_GetOpenOrders_WithTestCredentials(t *testing.T) {
	if !HasTestCredentials() {
		t.Skip("Skipping test: no test credentials available")
	}

	client, err := CreateTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	ctx := context.Background()

	// Get initial open orders count
	initialOrders, err := client.GetOpenOrders(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("Failed to get initial open orders: %v", err)
	}

	initialCount := len(initialOrders)
	t.Logf("Initial open orders for BTCUSDT: %d", initialCount)

	// Get current BTC price for realistic order
	result, err := client.GetTickerPrice(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("Failed to get BTC price: %v", err)
	}

	if !result.IsSingle() {
		t.Fatal("Expected single ticker result")
	}
	ticker := result.GetSingle()
	currentPrice, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		t.Fatalf("Failed to parse current price: %v", err)
	}

	orderPrice := currentPrice * 0.9
	orderPriceStr := strconv.FormatFloat(orderPrice, 'f', 2, 64)

	// Place an order
	orderReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		Price:       orderPriceStr,
	}

	orderResp, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("Failed to place order: %v", err)
	}

	// Wait a moment for the order to appear
	time.Sleep(1 * time.Second)

	// Get open orders again
	orders, err := client.GetOpenOrders(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("Failed to get open orders: %v", err)
	}

	// Should have one more order now
	if len(orders) != initialCount+1 {
		t.Errorf("Expected %d open orders, got %d", initialCount+1, len(orders))
	}

	// Find our order
	var foundOrder *OrderResponse
	for i := range orders {
		if orders[i].OrderId == orderResp.OrderId {
			foundOrder = &orders[i]
			break
		}
	}

	if foundOrder == nil {
		t.Error("Could not find our placed order in open orders list")
	} else {
		t.Logf("Found our order in open orders: ID=%d, Status=%s", foundOrder.OrderId, foundOrder.Status)
	}

	// Clean up
	_, err = client.CancelOrder(ctx, orderResp.Symbol, orderResp.OrderId)
	if err != nil {
		t.Errorf("Failed to clean up order: %v", err)
	}
}

func TestClient_OrderValidation(t *testing.T) {
	client := NewClient(TestnetConfig())

	ctx := context.Background()

	// Test empty symbol
	_, err := client.PlaceOrder(ctx, &NewOrderRequest{
		Side: SideBuy,
		Type: OrderTypeLimit,
	})
	if err == nil {
		t.Error("Expected error for empty symbol")
	}

	// Test invalid side
	_, err = client.PlaceOrder(ctx, &NewOrderRequest{
		Symbol: "BTCUSDT",
		Side:   "INVALID",
		Type:   OrderTypeLimit,
	})
	if err == nil {
		t.Error("Expected error for invalid side")
	}

	// Test invalid order type
	_, err = client.PlaceOrder(ctx, &NewOrderRequest{
		Symbol: "BTCUSDT",
		Side:   SideBuy,
		Type:   "INVALID",
	})
	if err == nil {
		t.Error("Expected error for invalid order type")
	}

	// Test LIMIT order without price
	_, err = client.PlaceOrder(ctx, &NewOrderRequest{
		Symbol:   "BTCUSDT",
		Side:     SideBuy,
		Type:     OrderTypeLimit,
		Quantity: "0.001",
	})
	if err == nil {
		t.Error("Expected error for LIMIT order without price")
	}

	// Test order without quantity or quoteOrderQty
	_, err = client.PlaceOrder(ctx, &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		Price:       "50000.00",
		TimeInForce: TimeInForceGTC,
	})
	if err == nil {
		t.Error("Expected error for order without quantity or quoteOrderQty")
	}
}
