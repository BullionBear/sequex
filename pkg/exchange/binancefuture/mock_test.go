package binancefuture

import (
	"testing"
	"time"
)

func TestValidateOrderRequest(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	// Test valid order request
	validReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		Price:       "50000",
	}

	err := client.validateOrderRequest(validReq)
	if err != nil {
		t.Errorf("valid order request should not return error: %v", err)
	}

	// Test invalid order request - missing symbol
	invalidReq := &NewOrderRequest{
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		Price:       "50000",
	}

	err = client.validateOrderRequest(invalidReq)
	if err == nil {
		t.Error("invalid order request should return error")
	}

	// Test invalid order request - missing price for limit order
	invalidReq2 := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		// Missing price
	}

	err = client.validateOrderRequest(invalidReq2)
	if err == nil {
		t.Error("limit order without price should return error")
	}

	// Test invalid order request - missing timeInForce for limit order
	invalidReq3 := &NewOrderRequest{
		Symbol:   "BTCUSDT",
		Side:     SideBuy,
		Type:     OrderTypeLimit,
		Quantity: "0.001",
		Price:    "50000",
		// Missing timeInForce
	}

	err = client.validateOrderRequest(invalidReq3)
	if err == nil {
		t.Error("limit order without timeInForce should return error")
	}

	// Test invalid order request - missing stopPrice for stop loss order
	invalidReq4 := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeStopLoss,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		// Missing stopPrice
	}

	err = client.validateOrderRequest(invalidReq4)
	if err == nil {
		t.Error("stop loss order without stopPrice should return error")
	}

	// Test valid market order (no price or timeInForce required)
	validMarketReq := &NewOrderRequest{
		Symbol:   "BTCUSDT",
		Side:     SideBuy,
		Type:     OrderTypeMarket,
		Quantity: "0.001",
	}

	err = client.validateOrderRequest(validMarketReq)
	if err != nil {
		t.Errorf("valid market order should not return error: %v", err)
	}

	// Test invalid side
	invalidSideReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        "INVALID",
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		Price:       "50000",
	}

	err = client.validateOrderRequest(invalidSideReq)
	if err == nil {
		t.Error("invalid side should return error")
	}

	// Test invalid timeInForce
	invalidTimeInForceReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: "INVALID",
		Quantity:    "0.001",
		Price:       "50000",
	}

	err = client.validateOrderRequest(invalidTimeInForceReq)
	if err == nil {
		t.Error("invalid timeInForce should return error")
	}
}

func TestOrderRequestToParams(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	// Test basic order request
	req := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		Price:       "50000",
	}

	params := client.orderRequestToParams(req)

	if params.Get("symbol") != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", params.Get("symbol"))
	}
	if params.Get("side") != SideBuy {
		t.Errorf("expected side %s, got %s", SideBuy, params.Get("side"))
	}
	if params.Get("type") != OrderTypeLimit {
		t.Errorf("expected type %s, got %s", OrderTypeLimit, params.Get("type"))
	}
	if params.Get("timeInForce") != TimeInForceGTC {
		t.Errorf("expected timeInForce %s, got %s", TimeInForceGTC, params.Get("timeInForce"))
	}
	if params.Get("quantity") != "0.001" {
		t.Errorf("expected quantity 0.001, got %s", params.Get("quantity"))
	}
	if params.Get("price") != "50000" {
		t.Errorf("expected price 50000, got %s", params.Get("price"))
	}

	// Test order with optional parameters
	reqWithOpts := &NewOrderRequest{
		Symbol:           "BTCUSDT",
		Side:             SideSell,
		Type:             OrderTypeLimit,
		TimeInForce:      TimeInForceGTC,
		Quantity:         "0.001",
		Price:            "50000",
		PositionSide:     PositionSideShort,
		ReduceOnly:       true,
		NewClientOrderId: "test-order-123",
		StopPrice:        "45000",
		WorkingType:      WorkingTypeMarkPrice,
		PriceProtect:     true,
		NewOrderRespType: NewOrderRespTypeFull,
		ClosePosition:    true,
		ActivationPrice:  "48000",
		CallbackRate:     "0.1",
	}

	paramsWithOpts := client.orderRequestToParams(reqWithOpts)

	if paramsWithOpts.Get("positionSide") != PositionSideShort {
		t.Errorf("expected positionSide %s, got %s", PositionSideShort, paramsWithOpts.Get("positionSide"))
	}
	if paramsWithOpts.Get("reduceOnly") != "true" {
		t.Errorf("expected reduceOnly true, got %s", paramsWithOpts.Get("reduceOnly"))
	}
	if paramsWithOpts.Get("newClientOrderId") != "test-order-123" {
		t.Errorf("expected newClientOrderId test-order-123, got %s", paramsWithOpts.Get("newClientOrderId"))
	}
	if paramsWithOpts.Get("stopPrice") != "45000" {
		t.Errorf("expected stopPrice 45000, got %s", paramsWithOpts.Get("stopPrice"))
	}
	if paramsWithOpts.Get("workingType") != WorkingTypeMarkPrice {
		t.Errorf("expected workingType %s, got %s", WorkingTypeMarkPrice, paramsWithOpts.Get("workingType"))
	}
	if paramsWithOpts.Get("priceProtect") != "true" {
		t.Errorf("expected priceProtect true, got %s", paramsWithOpts.Get("priceProtect"))
	}
	if paramsWithOpts.Get("newOrderRespType") != NewOrderRespTypeFull {
		t.Errorf("expected newOrderRespType %s, got %s", NewOrderRespTypeFull, paramsWithOpts.Get("newOrderRespType"))
	}
	if paramsWithOpts.Get("closePosition") != "true" {
		t.Errorf("expected closePosition true, got %s", paramsWithOpts.Get("closePosition"))
	}
	if paramsWithOpts.Get("activationPrice") != "48000" {
		t.Errorf("expected activationPrice 48000, got %s", paramsWithOpts.Get("activationPrice"))
	}
	if paramsWithOpts.Get("callbackRate") != "0.1" {
		t.Errorf("expected callbackRate 0.1, got %s", paramsWithOpts.Get("callbackRate"))
	}
}

func TestNewOrderRequestValidation(t *testing.T) {
	// Test that NewOrderRequest can be created with all fields
	req := &NewOrderRequest{
		Symbol:           "BTCUSDT",
		Side:             SideBuy,
		PositionSide:     PositionSideLong,
		Type:             OrderTypeLimit,
		TimeInForce:      TimeInForceGTC,
		Quantity:         "0.001",
		ReduceOnly:       true,
		Price:            "50000",
		NewClientOrderId: "test-123",
		StopPrice:        "45000",
		WorkingType:      WorkingTypeMarkPrice,
		PriceProtect:     true,
		NewOrderRespType: NewOrderRespTypeFull,
		ClosePosition:    true,
		ActivationPrice:  "48000",
		CallbackRate:     "0.1",
		RecvWindow:       5000,
		Timestamp:        time.Now().UnixNano() / int64(time.Millisecond),
	}

	// Just verify the struct can be created without panicking
	if req.Symbol != "BTCUSDT" {
		t.Error("symbol should be BTCUSDT")
	}
	if req.Side != SideBuy {
		t.Error("side should be BUY")
	}
	if req.Type != OrderTypeLimit {
		t.Error("type should be LIMIT")
	}
}
