package binance

import (
	"context"
	"strconv"
	"testing"
)

func TestClient_GetRealBTCPrice_AndPlaceOrder(t *testing.T) {
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

	ticker := result.(*TickerPriceResponse)
	t.Logf("Current BTC price: %s", ticker.Price)

	// Parse the current price
	currentPrice, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		t.Fatalf("Failed to parse current price: %v", err)
	}

	// Place a buy order 10% below current market price (very unlikely to fill)
	orderPrice := currentPrice * 0.9
	orderPriceStr := strconv.FormatFloat(orderPrice, 'f', 2, 64)

	t.Logf("Placing order at price: %s (10%% below market)", orderPriceStr)

	// Create a limit buy order
	orderReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001", // Small quantity
		Price:       orderPriceStr,
	}

	orderResp, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("Failed to place order: %v", err)
	}

	// Basic assertions
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

	t.Logf("Order placed successfully: ID=%d, Symbol=%s, Side=%s, Status=%s, Price=%s",
		orderResp.OrderId, orderResp.Symbol, orderResp.Side, orderResp.Status, orderResp.Price)

	// Test getting the order status
	order, err := client.GetOrder(ctx, orderResp.Symbol, orderResp.OrderId)
	if err != nil {
		t.Errorf("Failed to get order status: %v", err)
	} else {
		t.Logf("Order status: ID=%d, Status=%s, Price=%s, OrigQty=%s",
			order.OrderId, order.Status, order.Price, order.OrigQty)
	}

	// Test getting open orders
	openOrders, err := client.GetOpenOrders(ctx, "BTCUSDT")
	if err != nil {
		t.Errorf("Failed to get open orders: %v", err)
	} else {
		t.Logf("Found %d open orders for BTCUSDT", len(openOrders))

		// Find our order in the list
		var foundOrder *OrderResponse
		for i := range openOrders {
			if openOrders[i].OrderId == orderResp.OrderId {
				foundOrder = &openOrders[i]
				break
			}
		}

		if foundOrder != nil {
			t.Logf("Found our order in open orders: ID=%d, Status=%s", foundOrder.OrderId, foundOrder.Status)
		} else {
			t.Error("Could not find our order in open orders list")
		}
	}

	// Clean up: cancel the order
	cancelResp, err := client.CancelOrder(ctx, orderResp.Symbol, orderResp.OrderId)
	if err != nil {
		t.Errorf("Failed to cancel order: %v", err)
	} else {
		if cancelResp.OrderId != orderResp.OrderId {
			t.Errorf("Expected cancelled order ID %d, got %d", orderResp.OrderId, cancelResp.OrderId)
		}

		if cancelResp.Status != OrderStatusCanceled {
			t.Errorf("Expected status CANCELED, got %s", cancelResp.Status)
		}

		t.Logf("Order cancelled successfully: ID=%d, Status=%s", cancelResp.OrderId, cancelResp.Status)
	}
}
