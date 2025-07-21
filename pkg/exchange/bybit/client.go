package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Client represents a Bybit API client
type Client struct {
	config         *Config
	requestService *RequestService
}

// NewClient creates a new Bybit API client
func NewClient(config *Config) *Client {
	return &Client{
		config:         config,
		requestService: NewRequestService(config),
	}
}

// GetKline retrieves kline/candlestick data
func (c *Client) GetKline(ctx context.Context, req *KlineRequest) (*KlineResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("category", req.Category)
	params.Set("symbol", req.Symbol)
	params.Set("interval", req.Interval)

	if req.Start > 0 {
		params.Set("start", fmt.Sprintf("%d", req.Start))
	}
	if req.End > 0 {
		params.Set("end", fmt.Sprintf("%d", req.End))
	}
	if req.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", req.Limit))
	}

	// Make the request
	resp, err := c.requestService.DoUnsignedRequest(ctx, EndpointKlines, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get kline data: %w", err)
	}

	// Parse the response
	var klineResp KlineResponse
	if err := json.Unmarshal(resp, &klineResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kline response: %w", err)
	}

	// Check for API errors
	if klineResp.RetCode != 0 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", klineResp.RetCode, klineResp.RetMsg)
	}

	return &klineResp, nil
}

// GetKlineData retrieves kline data and returns parsed KlineData structs
func (c *Client) GetKlineData(ctx context.Context, req *KlineRequest) ([]*KlineData, error) {
	klineResp, err := c.GetKline(ctx, req)
	if err != nil {
		return nil, err
	}

	// Parse each kline data
	var klineData []*KlineData
	for _, data := range klineResp.Result.List {
		kline, err := ParseKlineData(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse kline data: %w", err)
		}
		klineData = append(klineData, kline)
	}

	return klineData, nil
}

// GetServerTime retrieves server time
func (c *Client) GetServerTime(ctx context.Context) (*ServerTimeResponse, error) {
	// Make the request
	resp, err := c.requestService.DoUnsignedRequest(ctx, EndpointServerTime, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get server time: %w", err)
	}

	// Parse the response
	var serverTimeResp ServerTimeResponse
	if err := json.Unmarshal(resp, &serverTimeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal server time response: %w", err)
	}

	// Check for API errors
	if serverTimeResp.RetCode != 0 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", serverTimeResp.RetCode, serverTimeResp.RetMsg)
	}

	return &serverTimeResp, nil
}

// GetTickers retrieves ticker information
func (c *Client) GetTickers(ctx context.Context, category, symbol string) (*TickerResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("category", category)
	if symbol != "" {
		params.Set("symbol", symbol)
	}

	// Make the request
	resp, err := c.requestService.DoUnsignedRequest(ctx, EndpointTicker24hr, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickers: %w", err)
	}

	// Parse the response
	var tickerResp TickerResponse
	if err := json.Unmarshal(resp, &tickerResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticker response: %w", err)
	}

	// Check for API errors
	if tickerResp.RetCode != 0 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", tickerResp.RetCode, tickerResp.RetMsg)
	}

	return &tickerResp, nil
}

// Trading Methods

// CreateOrder creates a new order
func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// Convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create order request: %w", err)
	}

	// Make the signed POST request
	resp, err := c.requestService.DoSignedPOSTRequest(ctx, EndpointNewOrder, nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Parse the response
	var createOrderResp CreateOrderResponse
	if err := json.Unmarshal(resp, &createOrderResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create order response: %w", err)
	}

	// Check for API errors
	if createOrderResp.RetCode != 0 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", createOrderResp.RetCode, createOrderResp.RetMsg)
	}

	return &createOrderResp, nil
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CancelOrderResponse, error) {
	// Convert request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cancel order request: %w", err)
	}

	// Make the signed POST request
	resp, err := c.requestService.DoSignedPOSTRequest(ctx, EndpointCancelOrder, nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	// Parse the response
	var cancelOrderResp CancelOrderResponse
	if err := json.Unmarshal(resp, &cancelOrderResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel order response: %w", err)
	}

	// Check for API errors
	if cancelOrderResp.RetCode != 0 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", cancelOrderResp.RetCode, cancelOrderResp.RetMsg)
	}

	return &cancelOrderResp, nil
}

// GetOrder retrieves order information (UTA 2.0)
func (c *Client) GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderListResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("category", req.Category)
	if req.Symbol != "" {
		params.Set("symbol", req.Symbol)
	}
	if req.OrderId != "" {
		params.Set("orderId", req.OrderId)
	}
	if req.OrderLinkId != "" {
		params.Set("orderLinkId", req.OrderLinkId)
	}
	if req.SettleCoin != "" {
		params.Set("settleCoin", req.SettleCoin)
	}
	if req.OrderFilter != "" {
		params.Set("orderFilter", req.OrderFilter)
	}
	if req.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", req.Limit))
	}
	if req.Cursor != "" {
		params.Set("cursor", req.Cursor)
	}

	// Make the signed GET request
	resp, err := c.requestService.DoSignedGETRequest(ctx, EndpointOrderStatus, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Parse the response
	var getOrderResp GetOrderListResponse
	if err := json.Unmarshal(resp, &getOrderResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal get order response: %w", err)
	}

	// Check for API errors
	if getOrderResp.RetCode != 0 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", getOrderResp.RetCode, getOrderResp.RetMsg)
	}

	return &getOrderResp, nil
}

// GetSingleOrder retrieves a single order by orderId or orderLinkId
func (c *Client) GetSingleOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error) {
	// For single order lookup, we need either orderId or orderLinkId
	if req.OrderId == "" && req.OrderLinkId == "" {
		return nil, fmt.Errorf("either orderId or orderLinkId must be provided")
	}

	// Build query parameters
	params := url.Values{}
	params.Set("category", req.Category)
	if req.Symbol != "" {
		params.Set("symbol", req.Symbol)
	}
	if req.OrderId != "" {
		params.Set("orderId", req.OrderId)
	}
	if req.OrderLinkId != "" {
		params.Set("orderLinkId", req.OrderLinkId)
	}

	// Make the signed GET request
	resp, err := c.requestService.DoSignedGETRequest(ctx, EndpointOrderStatus, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get single order: %w", err)
	}

	// Parse the response
	var getOrderResp GetOrderResponse
	if err := json.Unmarshal(resp, &getOrderResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal get single order response: %w", err)
	}

	// Check for API errors
	if getOrderResp.RetCode != 0 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", getOrderResp.RetCode, getOrderResp.RetMsg)
	}

	return &getOrderResp, nil
}

// GetAccount retrieves account information
func (c *Client) GetAccount(ctx context.Context, accountType string) (*AccountResponse, error) {
	// Build query parameters
	params := url.Values{}
	if accountType != "" {
		params.Set("accountType", accountType)
	}

	// Make the signed GET request
	resp, err := c.requestService.DoSignedGETRequest(ctx, EndpointAccount, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Parse the response
	var accountResp AccountResponse
	if err := json.Unmarshal(resp, &accountResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	// Check for API errors
	if accountResp.RetCode != 0 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", accountResp.RetCode, accountResp.RetMsg)
	}

	return &accountResp, nil
}
