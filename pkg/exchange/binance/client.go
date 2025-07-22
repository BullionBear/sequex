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

// CancelOrder cancels an active order on Binance Spot.
func (c *Client) CancelOrder(ctx context.Context, req CancelOrderRequest) (Response[CancelOrderResponse], error) {
	params := map[string]string{
		"symbol": req.Symbol,
	}
	if req.OrderId > 0 {
		params["orderId"] = fmt.Sprintf("%d", req.OrderId)
	}
	if req.OrigClientOrderId != "" {
		params["origClientOrderId"] = req.OrigClientOrderId
	}
	if req.NewClientOrderId != "" {
		params["newClientOrderId"] = req.NewClientOrderId
	}
	if req.CancelRestrictions != "" {
		params["cancelRestrictions"] = req.CancelRestrictions
	}
	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}
	body, status, err := doSignedRequest(c.cfg, http.MethodDelete, PathCancelOrder, params)
	if err != nil {
		return Response[CancelOrderResponse]{}, err
	}
	if status < 200 || status >= 300 {
		var errResp Response[CancelOrderResponse]
		_ = json.Unmarshal(body, &errResp)
		if errResp.Message == "" {
			errResp.Message = string(body)
		}
		return errResp, fmt.Errorf("binance error: %s", errResp.Message)
	}
	var resp CancelOrderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[CancelOrderResponse]{}, err
	}
	return Response[CancelOrderResponse]{Code: 0, Message: "success", Data: &resp}, nil
}

// CancelAllOrders cancels all open orders on a symbol on Binance Spot.
func (c *Client) CancelAllOrders(ctx context.Context, req CancelAllOrdersRequest) (Response[[]CancelOrderResponse], error) {
	params := map[string]string{
		"symbol": req.Symbol,
	}
	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}
	body, status, err := doSignedRequest(c.cfg, http.MethodDelete, PathCancelAllOrders, params)
	if err != nil {
		return Response[[]CancelOrderResponse]{}, err
	}
	if status < 200 || status >= 300 {
		var errResp Response[[]CancelOrderResponse]
		_ = json.Unmarshal(body, &errResp)
		if errResp.Message == "" {
			errResp.Message = string(body)
		}
		return errResp, fmt.Errorf("binance error: %s", errResp.Message)
	}
	var resp []CancelOrderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[[]CancelOrderResponse]{}, err
	}
	return Response[[]CancelOrderResponse]{Code: 0, Message: "success", Data: &resp}, nil
}

// GetDepth retrieves the order book depth for a symbol.
func (c *Client) GetDepth(ctx context.Context, symbol string, limit int) (Response[OrderBookDepthResponse], error) {
	params := map[string]string{"symbol": symbol}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	body, status, err := doUnsignedGet(c.cfg, PathGetDepth, params)
	if err != nil {
		return Response[OrderBookDepthResponse]{}, err
	}
	if status < 200 || status >= 300 {
		return Response[OrderBookDepthResponse]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	var resp OrderBookDepthResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[OrderBookDepthResponse]{}, err
	}
	return Response[OrderBookDepthResponse]{Code: 0, Message: "success", Data: &resp}, nil
}

// GetRecentTrades retrieves recent trades for a symbol.
func (c *Client) GetRecentTrades(ctx context.Context, symbol string, limit int) (Response[[]RecentTrade], error) {
	params := map[string]string{"symbol": symbol}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	body, status, err := doUnsignedGet(c.cfg, PathGetRecentTrades, params)
	if err != nil {
		return Response[[]RecentTrade]{}, err
	}
	if status < 200 || status >= 300 {
		return Response[[]RecentTrade]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	var trades []RecentTrade
	if err := json.Unmarshal(body, &trades); err != nil {
		return Response[[]RecentTrade]{}, err
	}
	return Response[[]RecentTrade]{Code: 0, Message: "success", Data: &trades}, nil
}

// GetAggTrades retrieves compressed, aggregate trades for a symbol.
func (c *Client) GetAggTrades(ctx context.Context, symbol string, fromId int64, startTime, endTime int64, limit int) (Response[[]AggTrade], error) {
	params := map[string]string{"symbol": symbol}
	if fromId > 0 {
		params["fromId"] = fmt.Sprintf("%d", fromId)
	}
	if startTime > 0 {
		params["startTime"] = fmt.Sprintf("%d", startTime)
	}
	if endTime > 0 {
		params["endTime"] = fmt.Sprintf("%d", endTime)
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	body, status, err := doUnsignedGet(c.cfg, PathGetAggTrades, params)
	if err != nil {
		return Response[[]AggTrade]{}, err
	}
	if status < 200 || status >= 300 {
		return Response[[]AggTrade]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	var trades []AggTrade
	if err := json.Unmarshal(body, &trades); err != nil {
		return Response[[]AggTrade]{}, err
	}
	return Response[[]AggTrade]{Code: 0, Message: "success", Data: &trades}, nil
}

// GetCandles retrieves kline/candlestick bars for a symbol.
func (c *Client) GetCandles(ctx context.Context, symbol, interval string, startTime, endTime int64, timeZone string, limit int) (Response[[]Kline], error) {
	params := map[string]string{"symbol": symbol, "interval": interval}
	if startTime > 0 {
		params["startTime"] = fmt.Sprintf("%d", startTime)
	}
	if endTime > 0 {
		params["endTime"] = fmt.Sprintf("%d", endTime)
	}
	if timeZone != "" {
		params["timeZone"] = timeZone
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	body, status, err := doUnsignedGet(c.cfg, PathGetKlines, params)
	if err != nil {
		return Response[[]Kline]{}, err
	}
	if status < 200 || status >= 300 {
		return Response[[]Kline]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	var raw [][]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return Response[[]Kline]{}, err
	}
	klines := make([]Kline, 0, len(raw))
	for _, k := range raw {
		if len(k) < 12 {
			continue
		}
		klines = append(klines, Kline{
			OpenTime:                 int64(k[0].(float64)),
			Open:                     k[1].(string),
			High:                     k[2].(string),
			Low:                      k[3].(string),
			Close:                    k[4].(string),
			Volume:                   k[5].(string),
			CloseTime:                int64(k[6].(float64)),
			QuoteAssetVolume:         k[7].(string),
			NumberOfTrades:           int(k[8].(float64)),
			TakerBuyBaseAssetVolume:  k[9].(string),
			TakerBuyQuoteAssetVolume: k[10].(string),
			Ignore:                   k[11].(string),
		})
	}
	return Response[[]Kline]{Code: 0, Message: "success", Data: &klines}, nil
}

// GetPriceTicker retrieves the latest price for a symbol or symbols.
func (c *Client) GetPriceTicker(ctx context.Context, symbols ...string) (Response[[]PriceTicker], error) {
	params := map[string]string{}
	if len(symbols) == 1 {
		params["symbol"] = symbols[0]
	} else if len(symbols) > 1 {
		b, err := json.Marshal(symbols)
		if err != nil {
			return Response[[]PriceTicker]{}, err
		}
		params["symbols"] = string(b)
	}
	body, status, err := doUnsignedGet(c.cfg, PathGetPriceTicker, params)
	if err != nil {
		return Response[[]PriceTicker]{}, err
	}
	if status < 200 || status >= 300 {
		return Response[[]PriceTicker]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	// Always unmarshal as []PriceTicker
	var tickers []PriceTicker
	if body[0] == '{' {
		// Single object, wrap in array
		var single PriceTicker
		if err := json.Unmarshal(body, &single); err != nil {
			return Response[[]PriceTicker]{}, err
		}
		tickers = append(tickers, single)
	} else {
		if err := json.Unmarshal(body, &tickers); err != nil {
			return Response[[]PriceTicker]{}, err
		}
	}
	return Response[[]PriceTicker]{Code: 0, Message: "success", Data: &tickers}, nil
}

// QueryOrder queries the status of an order on Binance Spot.
func (c *Client) QueryOrder(ctx context.Context, req QueryOrderRequest) (Response[QueryOrderResponse], error) {
	params := map[string]string{
		"symbol": req.Symbol,
	}
	if req.OrderId > 0 {
		params["orderId"] = fmt.Sprintf("%d", req.OrderId)
	}
	if req.OrigClientOrderId != "" {
		params["origClientOrderId"] = req.OrigClientOrderId
	}
	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}
	body, status, err := doSignedRequest(c.cfg, http.MethodGet, PathQueryOrder, params)
	if err != nil {
		return Response[QueryOrderResponse]{}, err
	}
	if status < 200 || status >= 300 {
		var errResp Response[QueryOrderResponse]
		_ = json.Unmarshal(body, &errResp)
		if errResp.Message == "" {
			errResp.Message = string(body)
		}
		return errResp, fmt.Errorf("binance error: %s", errResp.Message)
	}
	var resp QueryOrderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[QueryOrderResponse]{}, err
	}
	return Response[QueryOrderResponse]{Code: 0, Message: "success", Data: &resp}, nil
}

// GetAccountInfo retrieves current account information on Binance Spot.
func (c *Client) GetAccountInfo(ctx context.Context, req GetAccountInfoRequest) (Response[GetAccountInfoResponse], error) {
	params := map[string]string{}
	if req.OmitZeroBalances {
		params["omitZeroBalances"] = "true"
	}
	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}
	body, status, err := doSignedRequest(c.cfg, http.MethodGet, PathGetAccountInfo, params)
	if err != nil {
		return Response[GetAccountInfoResponse]{}, err
	}
	if status < 200 || status >= 300 {
		var errResp Response[GetAccountInfoResponse]
		_ = json.Unmarshal(body, &errResp)
		if errResp.Message == "" {
			errResp.Message = string(body)
		}
		return errResp, fmt.Errorf("binance error: %s", errResp.Message)
	}
	var resp GetAccountInfoResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[GetAccountInfoResponse]{}, err
	}
	return Response[GetAccountInfoResponse]{Code: 0, Message: "success", Data: &resp}, nil
}

// ListOpenOrders retrieves all open orders for a symbol or all symbols on Binance Spot.
func (c *Client) ListOpenOrders(ctx context.Context, req ListOpenOrdersRequest) (Response[[]QueryOrderResponse], error) {
	params := map[string]string{}
	if req.Symbol != "" {
		params["symbol"] = req.Symbol
	}
	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}
	body, status, err := doSignedRequest(c.cfg, http.MethodGet, PathListOpenOrders, params)
	if err != nil {
		return Response[[]QueryOrderResponse]{}, err
	}
	if status < 200 || status >= 300 {
		var errResp Response[[]QueryOrderResponse]
		_ = json.Unmarshal(body, &errResp)
		if errResp.Message == "" {
			errResp.Message = string(body)
		}
		return errResp, fmt.Errorf("binance error: %s", errResp.Message)
	}
	var resp []QueryOrderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[[]QueryOrderResponse]{}, err
	}
	return Response[[]QueryOrderResponse]{Code: 0, Message: "success", Data: &resp}, nil
}

// GetMyTrades retrieves trades for a specific account and symbol on Binance Spot.
func (c *Client) GetMyTrades(ctx context.Context, req GetAccountTradesRequest) (Response[[]AccountTrade], error) {
	params := map[string]string{"symbol": req.Symbol}
	if req.OrderId > 0 {
		params["orderId"] = fmt.Sprintf("%d", req.OrderId)
	}
	if req.StartTime > 0 {
		params["startTime"] = fmt.Sprintf("%d", req.StartTime)
	}
	if req.EndTime > 0 {
		params["endTime"] = fmt.Sprintf("%d", req.EndTime)
	}
	if req.FromId > 0 {
		params["fromId"] = fmt.Sprintf("%d", req.FromId)
	}
	if req.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}
	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}
	body, status, err := doSignedRequest(c.cfg, http.MethodGet, PathGetAccountTrades, params)
	if err != nil {
		return Response[[]AccountTrade]{}, err
	}
	if status < 200 || status >= 300 {
		var errResp Response[[]AccountTrade]
		_ = json.Unmarshal(body, &errResp)
		if errResp.Message == "" {
			errResp.Message = string(body)
		}
		return errResp, fmt.Errorf("binance error: %s", errResp.Message)
	}
	var resp []AccountTrade
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[[]AccountTrade]{}, err
	}
	return Response[[]AccountTrade]{Code: 0, Message: "success", Data: &resp}, nil
}
