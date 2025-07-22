package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is the Binance Spot API client.
type Client struct {
	cfg *Config
}

// NewClient creates a new Binance Spot API client.
func NewClient(cfg *Config) *Client {
	return &Client{cfg: cfg}
}

// CreateOrder places a new order on Binance Spot.
func (c *Client) CreateOrder(ctx context.Context, req CreateOrderRequest) (Response[CreateOrderResponse], error) {
	params := map[string]string{
		"symbol": req.Symbol,
		"side":   req.Side,
		"type":   req.Type,
	}
	if req.TimeInForce != "" {
		params["timeInForce"] = req.TimeInForce
	}
	if req.Quantity != "" {
		params["quantity"] = req.Quantity
	}
	if req.QuoteOrderQty != "" {
		params["quoteOrderQty"] = req.QuoteOrderQty
	}
	if req.Price != "" {
		params["price"] = req.Price
	}
	if req.NewClientOrderId != "" {
		params["newClientOrderId"] = req.NewClientOrderId
	}
	if req.StrategyId != 0 {
		params["strategyId"] = fmt.Sprintf("%d", req.StrategyId)
	}
	if req.StrategyType != 0 {
		params["strategyType"] = fmt.Sprintf("%d", req.StrategyType)
	}
	if req.StopPrice != "" {
		params["stopPrice"] = req.StopPrice
	}
	if req.TrailingDelta != 0 {
		params["trailingDelta"] = fmt.Sprintf("%d", req.TrailingDelta)
	}
	if req.IcebergQty != "" {
		params["icebergQty"] = req.IcebergQty
	}
	if req.NewOrderRespType != "" {
		params["newOrderRespType"] = req.NewOrderRespType
	}
	if req.SelfTradePreventionMode != "" {
		params["selfTradePreventionMode"] = req.SelfTradePreventionMode
	}
	if req.RecvWindow != 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}

	body, status, err := doSignedRequest(c.cfg, http.MethodPost, PathCreateOrder, params)
	if err != nil {
		return Response[CreateOrderResponse]{}, err
	}
	if status < 200 || status >= 300 {
		// Try to parse error response
		var errResp Response[CreateOrderResponse]
		_ = json.Unmarshal(body, &errResp)
		if errResp.Message == "" {
			errResp.Message = string(body)
		}
		return errResp, fmt.Errorf("binance error: %s", errResp.Message)
	}
	var resp CreateOrderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[CreateOrderResponse]{}, err
	}
	return Response[CreateOrderResponse]{
		Code:    0,
		Message: "success",
		Data:    &resp,
	}, nil
}
