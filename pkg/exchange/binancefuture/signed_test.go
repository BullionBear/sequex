package binancefuture

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestSignedEndpoints tests signed endpoints that require API credentials
// These tests will be skipped if API credentials are not provided via environment variables

func getSignedTestConfig() *Config {
	apiKey := os.Getenv("BINANCE_FUTURES_API_KEY")
	apiSecret := os.Getenv("BINANCE_FUTURES_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		return nil
	}

	config := TestnetConfig()
	config.APIKey = apiKey
	config.APISecret = apiSecret
	return config
}

func TestGetAccount(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	account, err := client.GetAccount(ctx)
	if err != nil {
		t.Fatalf("failed to get account: %v", err)
	}

	if account == nil {
		t.Fatal("account response should not be nil")
	}

	// Check basic account fields
	if account.UpdateTime <= 0 {
		t.Error("update time should be positive")
	}

	// Check assets
	if account.Assets == nil {
		t.Error("assets should not be nil")
	}

	// Check positions
	if account.Positions == nil {
		t.Error("positions should not be nil")
	}

	t.Logf("Account has %d assets and %d positions", len(account.Assets), len(account.Positions))
}

func TestPlaceAndCancelOrder(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Place a limit order that's unlikely to be filled (very low price)
	orderReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001", // Minimum quantity
		Price:       "1000",  // Very low price, unlikely to be filled
	}

	order, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("failed to place order: %v", err)
	}

	if order == nil {
		t.Fatal("order response should not be nil")
	}

	if order.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", order.Symbol)
	}

	if order.Side != SideBuy {
		t.Errorf("expected side BUY, got %s", order.Side)
	}

	if order.Type != OrderTypeLimit {
		t.Errorf("expected type LIMIT, got %s", order.Type)
	}

	t.Logf("Order placed successfully: OrderId=%d, Status=%s", order.OrderId, order.Status)

	// Cancel the order
	cancelResp, err := client.CancelOrder(ctx, "BTCUSDT", order.OrderId)
	if err != nil {
		t.Fatalf("failed to cancel order: %v", err)
	}

	if cancelResp == nil {
		t.Fatal("cancel order response should not be nil")
	}

	if cancelResp.OrderId != order.OrderId {
		t.Errorf("expected order ID %d, got %d", order.OrderId, cancelResp.OrderId)
	}

	t.Logf("Order canceled successfully: OrderId=%d, Status=%s", cancelResp.OrderId, cancelResp.Status)
}

func TestGetOrder(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First place an order
	orderReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		Price:       "1000",
	}

	placedOrder, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("failed to place order: %v", err)
	}

	// Get the order details
	order, err := client.GetOrder(ctx, "BTCUSDT", placedOrder.OrderId)
	if err != nil {
		t.Fatalf("failed to get order: %v", err)
	}

	if order == nil {
		t.Fatal("order response should not be nil")
	}

	if order.OrderId != placedOrder.OrderId {
		t.Errorf("expected order ID %d, got %d", placedOrder.OrderId, order.OrderId)
	}

	if order.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", order.Symbol)
	}

	// Clean up - cancel the order
	client.CancelOrder(ctx, "BTCUSDT", placedOrder.OrderId)
}

func TestGetOpenOrders(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all open orders
	orders, err := client.GetOpenOrders(ctx, "")
	if err != nil {
		t.Fatalf("failed to get open orders: %v", err)
	}

	if orders == nil {
		t.Fatal("open orders response should not be nil")
	}

	t.Logf("Found %d open orders", len(orders))

	// If there are open orders, test getting orders for a specific symbol
	if len(orders) > 0 {
		symbol := orders[0].Symbol
		symbolOrders, err := client.GetOpenOrders(ctx, symbol)
		if err != nil {
			t.Fatalf("failed to get open orders for symbol %s: %v", symbol, err)
		}

		if symbolOrders == nil {
			t.Fatal("symbol open orders response should not be nil")
		}

		// All returned orders should be for the specified symbol
		for _, order := range symbolOrders {
			if order.Symbol != symbol {
				t.Errorf("expected symbol %s, got %s", symbol, order.Symbol)
			}
		}
	}
}

func TestGetUserTrades(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get recent user trades for BTCUSDT
	trades, err := client.GetUserTrades(ctx, "BTCUSDT", 10)
	if err != nil {
		t.Fatalf("failed to get user trades: %v", err)
	}

	if trades == nil {
		t.Fatal("user trades response should not be nil")
	}

	t.Logf("Found %d user trades for BTCUSDT", len(trades))

	// Check trade structure if there are trades
	if len(trades) > 0 {
		trade := trades[0]
		if trade.Symbol != "BTCUSDT" {
			t.Errorf("expected symbol BTCUSDT, got %s", trade.Symbol)
		}
		if trade.Id <= 0 {
			t.Error("trade ID should be positive")
		}
		if trade.OrderId <= 0 {
			t.Error("order ID should be positive")
		}
		if trade.Price == "" {
			t.Error("trade price should not be empty")
		}
		if trade.Qty == "" {
			t.Error("trade quantity should not be empty")
		}
	}
}

func TestGetPositionRisk(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get position risk for all symbols
	positionRisks, err := client.GetPositionRisk(ctx, "")
	if err != nil {
		t.Fatalf("failed to get position risk: %v", err)
	}

	if positionRisks == nil {
		t.Fatal("position risk response should not be nil")
	}

	t.Logf("Found %d position risk entries", len(positionRisks))

	// Check position risk structure if there are entries
	if len(positionRisks) > 0 {
		positionRisk := positionRisks[0]
		if positionRisk.Symbol == "" {
			t.Error("symbol should not be empty")
		}
		if positionRisk.UpdateTime <= 0 {
			t.Error("update time should be positive")
		}
	}

	// Test getting position risk for a specific symbol
	if len(positionRisks) > 0 {
		symbol := positionRisks[0].Symbol
		symbolPositionRisks, err := client.GetPositionRisk(ctx, symbol)
		if err != nil {
			t.Fatalf("failed to get position risk for symbol %s: %v", symbol, err)
		}

		if symbolPositionRisks == nil {
			t.Fatal("symbol position risk response should not be nil")
		}

		// All returned position risks should be for the specified symbol
		for _, pr := range symbolPositionRisks {
			if pr.Symbol != symbol {
				t.Errorf("expected symbol %s, got %s", symbol, pr.Symbol)
			}
		}
	}
}

func TestGetPositionSide(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current position side mode
	positionSide, err := client.GetPositionSide(ctx)
	if err != nil {
		t.Fatalf("failed to get position side: %v", err)
	}

	if positionSide == nil {
		t.Fatal("position side response should not be nil")
	}

	t.Logf("Current position side mode: dualSidePosition=%t", positionSide.DualSidePosition)
}

func TestChangePositionSide(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First get current position side mode
	currentPositionSide, err := client.GetPositionSide(ctx)
	if err != nil {
		t.Fatalf("failed to get current position side: %v", err)
	}

	// Change to opposite mode
	newDualSidePosition := !currentPositionSide.DualSidePosition
	changedPositionSide, err := client.ChangePositionSide(ctx, newDualSidePosition)
	if err != nil {
		t.Fatalf("failed to change position side: %v", err)
	}

	if changedPositionSide == nil {
		t.Fatal("changed position side response should not be nil")
	}

	if changedPositionSide.DualSidePosition != newDualSidePosition {
		t.Errorf("expected dualSidePosition %t, got %t", newDualSidePosition, changedPositionSide.DualSidePosition)
	}

	t.Logf("Position side changed successfully: dualSidePosition=%t", changedPositionSide.DualSidePosition)

	// Change back to original mode
	originalPositionSide, err := client.ChangePositionSide(ctx, currentPositionSide.DualSidePosition)
	if err != nil {
		t.Fatalf("failed to change position side back: %v", err)
	}

	if originalPositionSide.DualSidePosition != currentPositionSide.DualSidePosition {
		t.Errorf("expected dualSidePosition %t, got %t", currentPositionSide.DualSidePosition, originalPositionSide.DualSidePosition)
	}

	t.Logf("Position side restored: dualSidePosition=%t", originalPositionSide.DualSidePosition)
}

func TestGetLeverage(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current leverage for BTCUSDT
	leverage, err := client.GetLeverage(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("failed to get leverage: %v", err)
	}

	if leverage == nil {
		t.Fatal("leverage response should not be nil")
	}

	if leverage.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", leverage.Symbol)
	}

	if leverage.Leverage <= 0 {
		t.Error("leverage should be positive")
	}

	t.Logf("Current leverage for %s: %d", leverage.Symbol, leverage.Leverage)
}

func TestChangeLeverage(t *testing.T) {
	config := getSignedTestConfig()
	if config == nil {
		t.Skip("Skipping signed endpoint test: BINANCE_FUTURES_API_KEY and BINANCE_FUTURES_API_SECRET environment variables not set")
	}

	client := NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First get current leverage
	currentLeverage, err := client.GetLeverage(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("failed to get current leverage: %v", err)
	}

	// Change to a different leverage (use 10 if current is not 10, otherwise use 20)
	newLeverage := 10
	if currentLeverage.Leverage == 10 {
		newLeverage = 20
	}

	changedLeverage, err := client.ChangeLeverage(ctx, "BTCUSDT", newLeverage)
	if err != nil {
		t.Fatalf("failed to change leverage: %v", err)
	}

	if changedLeverage == nil {
		t.Fatal("changed leverage response should not be nil")
	}

	if changedLeverage.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", changedLeverage.Symbol)
	}

	if changedLeverage.Leverage != newLeverage {
		t.Errorf("expected leverage %d, got %d", newLeverage, changedLeverage.Leverage)
	}

	t.Logf("Leverage changed successfully: %s -> %d", changedLeverage.Symbol, changedLeverage.Leverage)

	// Change back to original leverage
	originalLeverage, err := client.ChangeLeverage(ctx, "BTCUSDT", currentLeverage.Leverage)
	if err != nil {
		t.Fatalf("failed to change leverage back: %v", err)
	}

	if originalLeverage.Leverage != currentLeverage.Leverage {
		t.Errorf("expected leverage %d, got %d", currentLeverage.Leverage, originalLeverage.Leverage)
	}

	t.Logf("Leverage restored: %s -> %d", originalLeverage.Symbol, originalLeverage.Leverage)
}
