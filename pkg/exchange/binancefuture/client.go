package binancefuture

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
)

// Client represents the Binance Futures API client
type Client struct {
	config         *Config
	requestService *RequestService
	logger         *slog.Logger
}

// NewClient creates a new Binance Futures API client
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	requestService := NewRequestService(config)
	logger := slog.Default().With("component", "binance-futures-client")

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

// GetServerTime gets the server time from Binance Futures API
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
func (c *Client) GetTickerPrice(ctx context.Context, symbol string) (*TickerPriceResult, error) {
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

	result := &TickerPriceResult{}

	// If symbol is specified, return single ticker, otherwise return array
	if symbol != "" {
		var ticker TickerPriceResponse
		if err := json.Unmarshal(respBody, &ticker); err != nil {
			c.logger.Error("failed to parse ticker price response", "error", err)
			return nil, fmt.Errorf("failed to parse ticker price response: %w", err)
		}
		result.Single = &ticker
		c.logger.Debug("ticker price retrieved successfully", "symbol", symbol, "price", ticker.Price)
	} else {
		var tickers []TickerPriceResponse
		if err := json.Unmarshal(respBody, &tickers); err != nil {
			c.logger.Error("failed to parse ticker prices response", "error", err)
			return nil, fmt.Errorf("failed to parse ticker prices response: %w", err)
		}
		result.Array = tickers
		c.logger.Debug("ticker prices retrieved successfully", "count", len(tickers))
	}

	return result, nil
}

// GetTicker24hr gets 24hr ticker price change statistics
// If symbol is empty, returns tickers for all symbols
func (c *Client) GetTicker24hr(ctx context.Context, symbol string) (*Ticker24hrResult, error) {
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

	result := &Ticker24hrResult{}

	// If symbol is specified, return single ticker, otherwise return array
	if symbol != "" {
		var ticker Ticker24hrResponse
		if err := json.Unmarshal(respBody, &ticker); err != nil {
			c.logger.Error("failed to parse 24hr ticker response", "error", err)
			return nil, fmt.Errorf("failed to parse 24hr ticker response: %w", err)
		}
		result.Single = &ticker
		c.logger.Debug("24hr ticker retrieved successfully", "symbol", symbol)
	} else {
		var tickers []Ticker24hrResponse
		if err := json.Unmarshal(respBody, &tickers); err != nil {
			c.logger.Error("failed to parse 24hr tickers response", "error", err)
			return nil, fmt.Errorf("failed to parse 24hr tickers response: %w", err)
		}
		result.Array = tickers
		c.logger.Debug("24hr tickers retrieved successfully", "count", len(tickers))
	}

	return result, nil
}

// GetKlines gets kline/candlestick data for a symbol
func (c *Client) GetKlines(ctx context.Context, symbol, interval string, limit int) (*KlineResponse, error) {
	c.logger.Debug("getting klines", "symbol", symbol, "interval", interval, "limit", limit)

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("interval", interval)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointKlines,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get klines", "error", err, "symbol", symbol, "interval", interval)
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}

	var klines KlineResponse
	if err := json.Unmarshal(respBody, &klines); err != nil {
		c.logger.Error("failed to parse klines response", "error", err)
		return nil, fmt.Errorf("failed to parse klines response: %w", err)
	}

	c.logger.Debug("klines retrieved successfully", "symbol", symbol, "interval", interval, "count", len(klines))
	return &klines, nil
}

// GetOrderBook gets order book for a symbol
func (c *Client) GetOrderBook(ctx context.Context, symbol string, limit int) (*OrderBookResponse, error) {
	c.logger.Debug("getting order book", "symbol", symbol, "limit", limit)

	params := url.Values{}
	params.Set("symbol", symbol)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointOrderBook,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get order book", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}

	var orderBook OrderBookResponse
	if err := json.Unmarshal(respBody, &orderBook); err != nil {
		c.logger.Error("failed to parse order book response", "error", err)
		return nil, fmt.Errorf("failed to parse order book response: %w", err)
	}

	c.logger.Debug("order book retrieved successfully", "symbol", symbol, "bids", len(orderBook.Bids), "asks", len(orderBook.Asks))
	return &orderBook, nil
}

// GetTrades gets recent trades for a symbol
func (c *Client) GetTrades(ctx context.Context, symbol string, limit int) ([]TradeResponse, error) {
	c.logger.Debug("getting trades", "symbol", symbol, "limit", limit)

	params := url.Values{}
	params.Set("symbol", symbol)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointTrades,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get trades", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}

	var trades []TradeResponse
	if err := json.Unmarshal(respBody, &trades); err != nil {
		c.logger.Error("failed to parse trades response", "error", err)
		return nil, fmt.Errorf("failed to parse trades response: %w", err)
	}

	c.logger.Debug("trades retrieved successfully", "symbol", symbol, "count", len(trades))
	return trades, nil
}

// GetMarkPrice gets mark price for a symbol
// If symbol is empty, returns mark prices for all symbols
func (c *Client) GetMarkPrice(ctx context.Context, symbol string) ([]MarkPriceResponse, error) {
	c.logger.Debug("getting mark price", "symbol", symbol)

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}

	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointMarkPrice,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get mark price", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get mark price: %w", err)
	}

	var markPrices []MarkPriceResponse

	// If symbol is specified, return single mark price, otherwise return array
	if symbol != "" {
		var markPrice MarkPriceResponse
		if err := json.Unmarshal(respBody, &markPrice); err != nil {
			c.logger.Error("failed to parse mark price response", "error", err)
			return nil, fmt.Errorf("failed to parse mark price response: %w", err)
		}
		markPrices = []MarkPriceResponse{markPrice}
	} else {
		if err := json.Unmarshal(respBody, &markPrices); err != nil {
			c.logger.Error("failed to parse mark prices response", "error", err)
			return nil, fmt.Errorf("failed to parse mark prices response: %w", err)
		}
	}

	c.logger.Debug("mark price retrieved successfully", "symbol", symbol, "count", len(markPrices))
	return markPrices, nil
}

// GetFundingRate gets funding rate for a symbol
func (c *Client) GetFundingRate(ctx context.Context, symbol string, limit int) ([]FundingRateResponse, error) {
	c.logger.Debug("getting funding rate", "symbol", symbol, "limit", limit)

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointFundingRate,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get funding rate", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get funding rate: %w", err)
	}

	var fundingRates []FundingRateResponse
	if err := json.Unmarshal(respBody, &fundingRates); err != nil {
		c.logger.Error("failed to parse funding rate response", "error", err)
		return nil, fmt.Errorf("failed to parse funding rate response: %w", err)
	}

	c.logger.Debug("funding rate retrieved successfully", "symbol", symbol, "count", len(fundingRates))
	return fundingRates, nil
}

// GetOpenInterest gets open interest for a symbol
func (c *Client) GetOpenInterest(ctx context.Context, symbol string) (*OpenInterestResponse, error) {
	c.logger.Debug("getting open interest", "symbol", symbol)

	params := url.Values{}
	params.Set("symbol", symbol)

	respBody, err := c.requestService.DoUnsignedRequest(
		ctx,
		MethodGET,
		EndpointOpenInterest,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get open interest", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get open interest: %w", err)
	}

	var openInterest OpenInterestResponse
	if err := json.Unmarshal(respBody, &openInterest); err != nil {
		c.logger.Error("failed to parse open interest response", "error", err)
		return nil, fmt.Errorf("failed to parse open interest response: %w", err)
	}

	c.logger.Debug("open interest retrieved successfully", "symbol", symbol, "openInterest", openInterest.OpenInterest)
	return &openInterest, nil
}

// GetAccount gets current account information
func (c *Client) GetAccount(ctx context.Context) (*AccountResponse, error) {
	c.logger.Debug("getting account information")

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodGET,
		EndpointAccount,
		nil,
	)
	if err != nil {
		c.logger.Error("failed to get account", "error", err)
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	var account AccountResponse
	if err := json.Unmarshal(respBody, &account); err != nil {
		c.logger.Error("failed to parse account response", "error", err)
		return nil, fmt.Errorf("failed to parse account response: %w", err)
	}

	c.logger.Debug("account retrieved successfully", "assets", len(account.Assets), "positions", len(account.Positions))
	return &account, nil
}

// PlaceOrder places a new order
func (c *Client) PlaceOrder(ctx context.Context, req *NewOrderRequest) (*NewOrderResponse, error) {
	c.logger.Debug("placing order", "symbol", req.Symbol, "side", req.Side, "type", req.Type)

	// Validate request
	if err := c.validateOrderRequest(req); err != nil {
		return nil, fmt.Errorf("invalid order request: %w", err)
	}

	// Convert request to parameters
	params := c.orderRequestToParams(req)

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodPOST,
		EndpointNewOrder,
		params,
	)
	if err != nil {
		c.logger.Error("failed to place order", "error", err, "symbol", req.Symbol)
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	var order NewOrderResponse
	if err := json.Unmarshal(respBody, &order); err != nil {
		c.logger.Error("failed to parse order response", "error", err)
		return nil, fmt.Errorf("failed to parse order response: %w", err)
	}

	c.logger.Debug("order placed successfully", "symbol", order.Symbol, "orderId", order.OrderId, "status", order.Status)
	return &order, nil
}

// GetOrder gets order information
func (c *Client) GetOrder(ctx context.Context, symbol string, orderID int64) (*OrderResponse, error) {
	c.logger.Debug("getting order", "symbol", symbol, "orderId", orderID)

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
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var order OrderResponse
	if err := json.Unmarshal(respBody, &order); err != nil {
		c.logger.Error("failed to parse order response", "error", err)
		return nil, fmt.Errorf("failed to parse order response: %w", err)
	}

	c.logger.Debug("order retrieved successfully", "symbol", order.Symbol, "orderId", order.OrderId, "status", order.Status)
	return &order, nil
}

// CancelOrder cancels an order
func (c *Client) CancelOrder(ctx context.Context, symbol string, orderID int64) (*CancelOrderResponse, error) {
	c.logger.Debug("canceling order", "symbol", symbol, "orderId", orderID)

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
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	var order CancelOrderResponse
	if err := json.Unmarshal(respBody, &order); err != nil {
		c.logger.Error("failed to parse cancel order response", "error", err)
		return nil, fmt.Errorf("failed to parse cancel order response: %w", err)
	}

	c.logger.Debug("order canceled successfully", "symbol", order.Symbol, "orderId", order.OrderId, "status", order.Status)
	return &order, nil
}

// GetOpenOrders gets open orders for a symbol
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
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}

	var orders []OrderResponse
	if err := json.Unmarshal(respBody, &orders); err != nil {
		c.logger.Error("failed to parse open orders response", "error", err)
		return nil, fmt.Errorf("failed to parse open orders response: %w", err)
	}

	c.logger.Debug("open orders retrieved successfully", "symbol", symbol, "count", len(orders))
	return orders, nil
}

// GetUserTrades gets user trades for a symbol
func (c *Client) GetUserTrades(ctx context.Context, symbol string, limit int) ([]UserTradeResponse, error) {
	c.logger.Debug("getting user trades", "symbol", symbol, "limit", limit)

	params := url.Values{}
	params.Set("symbol", symbol)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodGET,
		EndpointMyTrades,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get user trades", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get user trades: %w", err)
	}

	var trades []UserTradeResponse
	if err := json.Unmarshal(respBody, &trades); err != nil {
		c.logger.Error("failed to parse user trades response", "error", err)
		return nil, fmt.Errorf("failed to parse user trades response: %w", err)
	}

	c.logger.Debug("user trades retrieved successfully", "symbol", symbol, "count", len(trades))
	return trades, nil
}

// GetPositionRisk gets position risk information
// If symbol is empty, returns position risk for all symbols
func (c *Client) GetPositionRisk(ctx context.Context, symbol string) ([]PositionRiskResponse, error) {
	c.logger.Debug("getting position risk", "symbol", symbol)

	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodGET,
		EndpointPositionRisk,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get position risk", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get position risk: %w", err)
	}

	var positionRisks []PositionRiskResponse
	if err := json.Unmarshal(respBody, &positionRisks); err != nil {
		c.logger.Error("failed to parse position risk response", "error", err)
		return nil, fmt.Errorf("failed to parse position risk response: %w", err)
	}

	c.logger.Debug("position risk retrieved successfully", "symbol", symbol, "count", len(positionRisks))
	return positionRisks, nil
}

// GetPositionSide gets current position side mode
func (c *Client) GetPositionSide(ctx context.Context) (*PositionSideResponse, error) {
	c.logger.Debug("getting position side mode")

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodGET,
		EndpointPositionSide,
		nil,
	)
	if err != nil {
		c.logger.Error("failed to get position side", "error", err)
		return nil, fmt.Errorf("failed to get position side: %w", err)
	}

	var positionSide PositionSideResponse
	if err := json.Unmarshal(respBody, &positionSide); err != nil {
		c.logger.Error("failed to parse position side response", "error", err)
		return nil, fmt.Errorf("failed to parse position side response: %w", err)
	}

	c.logger.Debug("position side retrieved successfully", "dualSidePosition", positionSide.DualSidePosition)
	return &positionSide, nil
}

// ChangePositionSide changes position side mode
func (c *Client) ChangePositionSide(ctx context.Context, dualSidePosition bool) (*PositionSideResponse, error) {
	c.logger.Debug("changing position side mode", "dualSidePosition", dualSidePosition)

	params := url.Values{}
	params.Set("dualSidePosition", strconv.FormatBool(dualSidePosition))

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodPOST,
		EndpointPositionSide,
		params,
	)
	if err != nil {
		c.logger.Error("failed to change position side", "error", err)
		return nil, fmt.Errorf("failed to change position side: %w", err)
	}

	var positionSide PositionSideResponse
	if err := json.Unmarshal(respBody, &positionSide); err != nil {
		c.logger.Error("failed to parse position side response", "error", err)
		return nil, fmt.Errorf("failed to parse position side response: %w", err)
	}

	c.logger.Debug("position side changed successfully", "dualSidePosition", positionSide.DualSidePosition)
	return &positionSide, nil
}

// GetLeverage gets current leverage for a symbol
func (c *Client) GetLeverage(ctx context.Context, symbol string) (*LeverageResponse, error) {
	c.logger.Debug("getting leverage", "symbol", symbol)

	params := url.Values{}
	params.Set("symbol", symbol)

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodGET,
		EndpointLeverage,
		params,
	)
	if err != nil {
		c.logger.Error("failed to get leverage", "error", err, "symbol", symbol)
		return nil, fmt.Errorf("failed to get leverage: %w", err)
	}

	var leverage LeverageResponse
	if err := json.Unmarshal(respBody, &leverage); err != nil {
		c.logger.Error("failed to parse leverage response", "error", err)
		return nil, fmt.Errorf("failed to parse leverage response: %w", err)
	}

	c.logger.Debug("leverage retrieved successfully", "symbol", leverage.Symbol, "leverage", leverage.Leverage)
	return &leverage, nil
}

// ChangeLeverage changes leverage for a symbol
func (c *Client) ChangeLeverage(ctx context.Context, symbol string, leverage int) (*LeverageResponse, error) {
	c.logger.Debug("changing leverage", "symbol", symbol, "leverage", leverage)

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("leverage", strconv.Itoa(leverage))

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodPOST,
		EndpointLeverage,
		params,
	)
	if err != nil {
		c.logger.Error("failed to change leverage", "error", err, "symbol", symbol, "leverage", leverage)
		return nil, fmt.Errorf("failed to change leverage: %w", err)
	}

	var leverageResp LeverageResponse
	if err := json.Unmarshal(respBody, &leverageResp); err != nil {
		c.logger.Error("failed to parse leverage response", "error", err)
		return nil, fmt.Errorf("failed to parse leverage response: %w", err)
	}

	c.logger.Debug("leverage changed successfully", "symbol", leverageResp.Symbol, "leverage", leverageResp.Leverage)
	return &leverageResp, nil
}

// CreateUserDataStream creates a new user data stream
func (c *Client) CreateUserDataStream(ctx context.Context) (*UserDataStreamResponse, error) {
	c.logger.Debug("creating user data stream")

	respBody, err := c.requestService.DoSignedRequest(
		ctx,
		MethodPOST,
		EndpointUserDataStream,
		nil,
	)
	if err != nil {
		c.logger.Error("failed to create user data stream", "error", err)
		return nil, fmt.Errorf("failed to create user data stream: %w", err)
	}

	var userDataStream UserDataStreamResponse
	if err := json.Unmarshal(respBody, &userDataStream); err != nil {
		c.logger.Error("failed to parse user data stream response", "error", err)
		return nil, fmt.Errorf("failed to parse user data stream response: %w", err)
	}

	c.logger.Debug("user data stream created successfully", "listenKey", userDataStream.ListenKey[:8]+"...")
	return &userDataStream, nil
}

// KeepAliveUserDataStream extends the validity of a user data stream
func (c *Client) KeepAliveUserDataStream(ctx context.Context, listenKey string) error {
	c.logger.Debug("keeping alive user data stream", "listenKey", listenKey[:8]+"...")

	params := url.Values{}
	params.Set("listenKey", listenKey)

	_, err := c.requestService.DoSignedRequest(
		ctx,
		MethodPUT,
		EndpointUserDataStream,
		params,
	)
	if err != nil {
		c.logger.Error("failed to keep alive user data stream", "error", err)
		return fmt.Errorf("failed to keep alive user data stream: %w", err)
	}

	c.logger.Debug("user data stream kept alive successfully")
	return nil
}

// CloseUserDataStream closes a user data stream
func (c *Client) CloseUserDataStream(ctx context.Context, listenKey string) error {
	c.logger.Debug("closing user data stream", "listenKey", listenKey[:8]+"...")

	params := url.Values{}
	params.Set("listenKey", listenKey)

	_, err := c.requestService.DoSignedRequest(
		ctx,
		MethodDELETE,
		EndpointUserDataStream,
		params,
	)
	if err != nil {
		c.logger.Error("failed to close user data stream", "error", err)
		return fmt.Errorf("failed to close user data stream: %w", err)
	}

	c.logger.Debug("user data stream closed successfully")
	return nil
}

// validateOrderRequest validates the order request parameters
func (c *Client) validateOrderRequest(req *NewOrderRequest) error {
	if req.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if req.Side == "" {
		return fmt.Errorf("side is required")
	}
	if req.Type == "" {
		return fmt.Errorf("type is required")
	}

	// Validate side
	if req.Side != SideBuy && req.Side != SideSell {
		return fmt.Errorf("invalid side: %s", req.Side)
	}

	// Validate order type
	switch req.Type {
	case OrderTypeLimit:
		if req.Price == "" {
			return fmt.Errorf("price is required for limit orders")
		}
		if req.TimeInForce == "" {
			return fmt.Errorf("timeInForce is required for limit orders")
		}
	case OrderTypeMarket:
		// Market orders don't require price or timeInForce
	case OrderTypeStopLoss, OrderTypeStopLossLimit:
		if req.StopPrice == "" {
			return fmt.Errorf("stopPrice is required for stop orders")
		}
	case OrderTypeTakeProfit, OrderTypeTakeProfitLimit:
		if req.StopPrice == "" {
			return fmt.Errorf("stopPrice is required for take profit orders")
		}
	}

	// Validate time in force for limit orders
	if req.TimeInForce != "" {
		switch req.TimeInForce {
		case TimeInForceGTC, TimeInForceIOC, TimeInForceFOK, TimeInForceGTX:
			// Valid values
		default:
			return fmt.Errorf("invalid timeInForce: %s", req.TimeInForce)
		}
	}

	return nil
}

// orderRequestToParams converts NewOrderRequest to URL parameters
func (c *Client) orderRequestToParams(req *NewOrderRequest) url.Values {
	params := url.Values{}
	params.Set("symbol", req.Symbol)
	params.Set("side", req.Side)
	params.Set("type", req.Type)

	if req.PositionSide != "" {
		params.Set("positionSide", req.PositionSide)
	}
	if req.TimeInForce != "" {
		params.Set("timeInForce", req.TimeInForce)
	}
	if req.Quantity != "" {
		params.Set("quantity", req.Quantity)
	}
	if req.ReduceOnly {
		params.Set("reduceOnly", "true")
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
	if req.WorkingType != "" {
		params.Set("workingType", req.WorkingType)
	}
	if req.PriceProtect {
		params.Set("priceProtect", "true")
	}
	if req.NewOrderRespType != "" {
		params.Set("newOrderRespType", req.NewOrderRespType)
	}
	if req.ClosePosition {
		params.Set("closePosition", "true")
	}
	if req.ActivationPrice != "" {
		params.Set("activationPrice", req.ActivationPrice)
	}
	if req.CallbackRate != "" {
		params.Set("callbackRate", req.CallbackRate)
	}

	return params
}
