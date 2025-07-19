package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
)

// Client represents the Binance API client
type Client struct {
	config         *Config
	requestService *RequestService
	logger         *slog.Logger
}

// NewClient creates a new Binance API client
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	requestService := NewRequestService(config)
	logger := slog.Default().With("component", "binance-client")

	return &Client{
		config:         config,
		requestService: requestService,
		logger:         logger,
	}
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *Config {
	return c.config
}

// GetServerTime gets the server time from Binance API
// This is a public endpoint that doesn't require authentication
func (c *Client) GetServerTime(ctx context.Context) (*ServerTimeResponse, error) {
	c.logger.Debug("getting server time")

	// Make request to server time endpoint
	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointServerTime,
		nil,
	)
	if err != nil {
		c.logger.Error("failed to get server time", "error", err)
		// Return API errors directly without wrapping
		if _, ok := err.(*APIError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get server time: %w", err)
	}

	// Parse response
	var serverTime ServerTimeResponse
	if err := json.Unmarshal(respBody, &serverTime); err != nil {
		c.logger.Error("failed to parse server time response", "error", err, "body", string(respBody))
		return nil, fmt.Errorf("failed to parse server time response: %w", err)
	}

	c.logger.Debug("server time retrieved successfully", "serverTime", serverTime.ServerTime)
	return &serverTime, nil
}

// Ping tests connectivity to the REST API
func (c *Client) Ping(ctx context.Context) error {
	c.logger.Debug("pinging server")

	_, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointPing,
		nil,
	)
	if err != nil {
		c.logger.Error("ping failed", "error", err)
		// Return API errors directly without wrapping
		if _, ok := err.(*APIError); ok {
			return err
		}
		return fmt.Errorf("ping failed: %w", err)
	}

	c.logger.Debug("ping successful")
	return nil
}

// GetExchangeInfo gets current exchange trading rules and symbol information
func (c *Client) GetExchangeInfo(ctx context.Context) (*ExchangeInfoResponse, error) {
	c.logger.Debug("getting exchange info")

	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointExchangeInfo,
		nil,
	)
	if err != nil {
		c.logger.Error("failed to get exchange info", "error", err)
		return nil, fmt.Errorf("failed to get exchange info: %w", err)
	}

	var exchangeInfo ExchangeInfoResponse
	if err := json.Unmarshal(respBody, &exchangeInfo); err != nil {
		c.logger.Error("failed to parse exchange info response", "error", err)
		return nil, fmt.Errorf("failed to parse exchange info response: %w", err)
	}

	c.logger.Debug("exchange info retrieved successfully", "symbolCount", len(exchangeInfo.Symbols))
	return &exchangeInfo, nil
}

// GetTickerPrice gets symbol price ticker
// If symbol is empty, returns price tickers for all symbols
func (c *Client) GetTickerPrice(ctx context.Context, symbol string) (interface{}, error) {
	c.logger.Debug("getting ticker price", "symbol", symbol)

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}

	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointTickerPrice,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get ticker price", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get ticker price: %w", err)
	}

	// If symbol is specified, return single ticker, otherwise return array
	if symbol != "" {
		var ticker TickerPriceResponse
		if err := json.Unmarshal(respBody, &ticker); err != nil {
			c.logger.Error("failed to parse ticker price response", "error", err)
			return nil, fmt.Errorf("failed to parse ticker price response: %w", err)
		}
		c.logger.Debug("ticker price retrieved successfully", "symbol", symbol, "price", ticker.Price)
		return &ticker, nil
	} else {
		var tickers []TickerPriceResponse
		if err := json.Unmarshal(respBody, &tickers); err != nil {
			c.logger.Error("failed to parse ticker prices response", "error", err)
			return nil, fmt.Errorf("failed to parse ticker prices response: %w", err)
		}
		c.logger.Debug("ticker prices retrieved successfully", "count", len(tickers))
		return tickers, nil
	}
}

// GetTicker24hr gets 24hr ticker price change statistics
// If symbol is empty, returns tickers for all symbols
func (c *Client) GetTicker24hr(ctx context.Context, symbol string) (interface{}, error) {
	c.logger.Debug("getting 24hr ticker", "symbol", symbol)

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}

	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointTicker24hr,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get 24hr ticker", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get 24hr ticker: %w", err)
	}

	// If symbol is specified, return single ticker, otherwise return array
	if symbol != "" {
		var ticker Ticker24hrResponse
		if err := json.Unmarshal(respBody, &ticker); err != nil {
			c.logger.Error("failed to parse 24hr ticker response", "error", err)
			return nil, fmt.Errorf("failed to parse 24hr ticker response: %w", err)
		}
		c.logger.Debug("24hr ticker retrieved successfully", "symbol", symbol)
		return &ticker, nil
	} else {
		var tickers []Ticker24hrResponse
		if err := json.Unmarshal(respBody, &tickers); err != nil {
			c.logger.Error("failed to parse 24hr tickers response", "error", err)
			return nil, fmt.Errorf("failed to parse 24hr tickers response: %w", err)
		}
		c.logger.Debug("24hr tickers retrieved successfully", "count", len(tickers))
		return tickers, nil
	}
}

// GetAccount gets current account information (requires signature)
func (c *Client) GetAccount(ctx context.Context) (*AccountResponse, error) {
	c.logger.Debug("getting account information")

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodGET,
		EndpointAccount,
		nil,
	)
	if err != nil {
		c.logger.Error("failed to get account information", "error", err)
		// Return API errors directly without wrapping
		if _, ok := err.(*APIError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get account information: %w", err)
	}

	var account AccountResponse
	if err := json.Unmarshal(respBody, &account); err != nil {
		c.logger.Error("failed to parse account response", "error", err)
		return nil, fmt.Errorf("failed to parse account response: %w", err)
	}

	c.logger.Debug("account information retrieved successfully", "balanceCount", len(account.Balances))
	return &account, nil
}

// PlaceOrder places a new order (requires signature)
func (c *Client) PlaceOrder(ctx context.Context, req *NewOrderRequest) (*NewOrderResponse, error) {
	c.logger.Debug("placing new order", "symbol", req.Symbol, "side", req.Side, "type", req.Type)

	// Validate required fields
	if err := c.validateOrderRequest(req); err != nil {
		return nil, fmt.Errorf("invalid order request: %w", err)
	}

	// Set timestamp if not provided
	if req.Timestamp == 0 {
		req.Timestamp = GetCurrentTimestamp()
	}

	// Convert request to query parameters
	params := c.orderRequestToParams(req)

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodPOST,
		EndpointNewOrder,
		params,
	)
	if err != nil {
		c.logger.Error("failed to place order", "error", err, "symbol", req.Symbol)
		// Return API errors directly without wrapping
		if _, ok := err.(*APIError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	var orderResp NewOrderResponse
	if err := json.Unmarshal(respBody, &orderResp); err != nil {
		c.logger.Error("failed to parse order response", "error", err)
		return nil, fmt.Errorf("failed to parse order response: %w", err)
	}

	c.logger.Debug("order placed successfully", "orderId", orderResp.OrderId, "symbol", orderResp.Symbol)
	return &orderResp, nil
}

// GetOrder gets order status (requires signature)
func (c *Client) GetOrder(ctx context.Context, symbol string, orderID int64) (*OrderResponse, error) {
	c.logger.Debug("getting order status", "symbol", symbol, "orderId", orderID)

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("orderId", strconv.FormatInt(orderID, 10))

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodGET,
		EndpointOrderStatus,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get order", "error", err, "symbol", symbol, "orderId", orderID)
		// Return API errors directly without wrapping
		if _, ok := err.(*APIError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var order OrderResponse
	if err := json.Unmarshal(respBody, &order); err != nil {
		c.logger.Error("failed to parse order response", "error", err)
		return nil, fmt.Errorf("failed to parse order response: %w", err)
	}

	c.logger.Debug("order retrieved successfully", "orderId", order.OrderId, "status", order.Status)
	return &order, nil
}

// CancelOrder cancels an active order (requires signature)
func (c *Client) CancelOrder(ctx context.Context, symbol string, orderID int64) (*CancelOrderResponse, error) {
	c.logger.Debug("cancelling order", "symbol", symbol, "orderId", orderID)

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("orderId", strconv.FormatInt(orderID, 10))

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodDELETE,
		EndpointCancelOrder,
		params,
	)
	if err != nil {
		c.logger.Error("failed to cancel order", "error", err, "symbol", symbol, "orderId", orderID)
		// Return API errors directly without wrapping
		if _, ok := err.(*APIError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	var cancelResp CancelOrderResponse
	if err := json.Unmarshal(respBody, &cancelResp); err != nil {
		c.logger.Error("failed to parse cancel order response", "error", err)
		return nil, fmt.Errorf("failed to parse cancel order response: %w", err)
	}

	c.logger.Debug("order cancelled successfully", "orderId", cancelResp.OrderId, "symbol", cancelResp.Symbol)
	return &cancelResp, nil
}

// GetOpenOrders gets all open orders for a symbol (requires signature)
func (c *Client) GetOpenOrders(ctx context.Context, symbol string) ([]OrderResponse, error) {
	c.logger.Debug("getting open orders", "symbol", symbol)

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodGET,
		EndpointOpenOrders,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get open orders", "error", err, "symbol", symbol)
		// Return API errors directly without wrapping
		if _, ok := err.(*APIError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}

	var orders []OrderResponse
	if err := json.Unmarshal(respBody, &orders); err != nil {
		c.logger.Error("failed to parse open orders response", "error", err)
		return nil, fmt.Errorf("failed to parse open orders response: %w", err)
	}

	c.logger.Debug("open orders retrieved successfully", "count", len(orders))
	return orders, nil
}

// validateOrderRequest validates the order request parameters
func (c *Client) validateOrderRequest(req *NewOrderRequest) error {
	if req.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if !ValidateSide(req.Side) {
		return fmt.Errorf("invalid side: %s", req.Side)
	}
	if !ValidateOrderType(req.Type) {
		return fmt.Errorf("invalid order type: %s", req.Type)
	}

	// For LIMIT orders, price and timeInForce are required
	if req.Type == OrderTypeLimit {
		if req.Price == "" {
			return fmt.Errorf("price is required for LIMIT orders")
		}
		if req.TimeInForce == "" {
			return fmt.Errorf("timeInForce is required for LIMIT orders")
		}
		if !ValidateTimeInForce(req.TimeInForce) {
			return fmt.Errorf("invalid timeInForce: %s", req.TimeInForce)
		}
	}

	// Either quantity or quoteOrderQty must be specified
	if req.Quantity == "" && req.QuoteOrderQty == "" {
		return fmt.Errorf("either quantity or quoteOrderQty must be specified")
	}

	return nil
}

// orderRequestToParams converts NewOrderRequest to url.Values
func (c *Client) orderRequestToParams(req *NewOrderRequest) url.Values {
	params := url.Values{}

	params.Set("symbol", req.Symbol)
	params.Set("side", req.Side)
	params.Set("type", req.Type)

	if req.TimeInForce != "" {
		params.Set("timeInForce", req.TimeInForce)
	}
	if req.Quantity != "" {
		params.Set("quantity", req.Quantity)
	}
	if req.QuoteOrderQty != "" {
		params.Set("quoteOrderQty", req.QuoteOrderQty)
	}
	if req.Price != "" {
		params.Set("price", req.Price)
	}
	if req.NewClientOrderId != "" {
		params.Set("newClientOrderId", req.NewClientOrderId)
	}
	if req.StopPrice != "" {
		params.Set("stopPrice", req.StopPrice)
	}
	if req.IcebergQty != "" {
		params.Set("icebergQty", req.IcebergQty)
	}
	if req.NewOrderRespType != "" {
		params.Set("newOrderRespType", req.NewOrderRespType)
	}
	if req.RecvWindow > 0 {
		params.Set("recvWindow", strconv.FormatInt(req.RecvWindow, 10))
	}

	return params
}
