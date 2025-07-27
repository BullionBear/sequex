package binanceperp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is the Binance Perpetual Futures API client.
type Client struct {
	cfg *Config
}

// NewClient creates a new Binance Perpetual Futures API client.
func NewClient(cfg *Config) *Client {
	return &Client{cfg: cfg}
}

// GetServerTime tests connectivity to the Rest API and gets the current server time.
func (c *Client) GetServerTime(ctx context.Context) (Response[GetServerTimeResponse], error) {
	body, status, err := doUnsignedGet(c.cfg, PathGetServerTime, nil)
	if err != nil {
		return Response[GetServerTimeResponse]{}, err
	}
	if status != http.StatusOK {
		return Response[GetServerTimeResponse]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	var resp GetServerTimeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[GetServerTimeResponse]{}, err
	}
	return Response[GetServerTimeResponse]{Code: 0, Message: "success", Data: &resp}, nil
}

// GetDepth queries symbol orderbook.
func (c *Client) GetDepth(ctx context.Context, req GetDepthRequest) (Response[GetDepthResponse], error) {
	params := map[string]string{
		"symbol": req.Symbol,
	}
	if req.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}

	body, status, err := doUnsignedGet(c.cfg, PathGetDepth, params)
	if err != nil {
		return Response[GetDepthResponse]{}, err
	}
	if status != http.StatusOK {
		return Response[GetDepthResponse]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	var resp GetDepthResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[GetDepthResponse]{}, err
	}
	return Response[GetDepthResponse]{Code: 0, Message: "success", Data: &resp}, nil
}

// GetRecentTrades gets recent market trades.
func (c *Client) GetRecentTrades(ctx context.Context, req GetRecentTradesRequest) (Response[[]RecentTrade], error) {
	params := map[string]string{
		"symbol": req.Symbol,
	}
	if req.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}

	body, status, err := doUnsignedGet(c.cfg, PathGetRecentTrades, params)
	if err != nil {
		return Response[[]RecentTrade]{}, err
	}
	if status != http.StatusOK {
		return Response[[]RecentTrade]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	var trades []RecentTrade
	if err := json.Unmarshal(body, &trades); err != nil {
		return Response[[]RecentTrade]{}, err
	}
	return Response[[]RecentTrade]{Code: 0, Message: "success", Data: &trades}, nil
}

// GetAggTrades gets compressed, aggregate market trades.
func (c *Client) GetAggTrades(ctx context.Context, req GetAggTradesRequest) (Response[[]AggTrade], error) {
	params := map[string]string{
		"symbol": req.Symbol,
	}
	if req.FromId > 0 {
		params["fromId"] = fmt.Sprintf("%d", req.FromId)
	}
	if req.StartTime > 0 {
		params["startTime"] = fmt.Sprintf("%d", req.StartTime)
	}
	if req.EndTime > 0 {
		params["endTime"] = fmt.Sprintf("%d", req.EndTime)
	}
	if req.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}

	body, status, err := doUnsignedGet(c.cfg, PathGetAggTrades, params)
	if err != nil {
		return Response[[]AggTrade]{}, err
	}
	if status != http.StatusOK {
		return Response[[]AggTrade]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	var trades []AggTrade
	if err := json.Unmarshal(body, &trades); err != nil {
		return Response[[]AggTrade]{}, err
	}
	return Response[[]AggTrade]{Code: 0, Message: "success", Data: &trades}, nil
}

// GetKlines gets kline/candlestick bars for a symbol.
func (c *Client) GetKlines(ctx context.Context, req GetKlinesRequest) (Response[[]Kline], error) {
	params := map[string]string{
		"symbol":   req.Symbol,
		"interval": req.Interval,
	}
	if req.StartTime > 0 {
		params["startTime"] = fmt.Sprintf("%d", req.StartTime)
	}
	if req.EndTime > 0 {
		params["endTime"] = fmt.Sprintf("%d", req.EndTime)
	}
	if req.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}

	body, status, err := doUnsignedGet(c.cfg, PathGetKlines, params)
	if err != nil {
		return Response[[]Kline]{}, err
	}
	if status != http.StatusOK {
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

// GetMarkPrice gets mark price and funding rate data.
func (c *Client) GetMarkPrice(ctx context.Context, req GetMarkPriceRequest) (Response[[]MarkPrice], error) {
	params := map[string]string{}
	if req.Symbol != "" {
		params["symbol"] = req.Symbol
	}

	body, status, err := doUnsignedGet(c.cfg, PathGetMarkPrice, params)
	if err != nil {
		return Response[[]MarkPrice]{}, err
	}
	if status != http.StatusOK {
		return Response[[]MarkPrice]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}

	// Handle both single object and array responses
	var markPrices []MarkPrice

	// Try to unmarshal as array first
	if err := json.Unmarshal(body, &markPrices); err != nil {
		// If that fails, try to unmarshal as single object
		var singleMarkPrice MarkPrice
		if err := json.Unmarshal(body, &singleMarkPrice); err != nil {
			return Response[[]MarkPrice]{}, err
		}
		markPrices = []MarkPrice{singleMarkPrice}
	}

	return Response[[]MarkPrice]{Code: 0, Message: "success", Data: &markPrices}, nil
}

// GetPriceTicker gets latest price for a symbol or symbols.
func (c *Client) GetPriceTicker(ctx context.Context, req GetPriceTickerRequest) (Response[[]PriceTicker], error) {
	params := map[string]string{}
	if req.Symbol != "" {
		params["symbol"] = req.Symbol
	}

	body, status, err := doUnsignedGet(c.cfg, PathGetPriceTicker, params)
	if err != nil {
		return Response[[]PriceTicker]{}, err
	}
	if status != http.StatusOK {
		return Response[[]PriceTicker]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}

	// Handle both single object and array responses
	var priceTickers []PriceTicker

	// Try to unmarshal as array first
	if err := json.Unmarshal(body, &priceTickers); err != nil {
		// If that fails, try to unmarshal as single object
		var singlePriceTicker PriceTicker
		if err := json.Unmarshal(body, &singlePriceTicker); err != nil {
			return Response[[]PriceTicker]{}, err
		}
		priceTickers = []PriceTicker{singlePriceTicker}
	}

	return Response[[]PriceTicker]{Code: 0, Message: "success", Data: &priceTickers}, nil
}

// GetBookTicker gets best price/qty on the order book for a symbol or symbols.
func (c *Client) GetBookTicker(ctx context.Context, req GetBookTickerRequest) (Response[[]BookTicker], error) {
	params := map[string]string{}
	if req.Symbol != "" {
		params["symbol"] = req.Symbol
	}

	body, status, err := doUnsignedGet(c.cfg, PathGetBookTicker, params)
	if err != nil {
		return Response[[]BookTicker]{}, err
	}
	if status != http.StatusOK {
		return Response[[]BookTicker]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}

	// Handle both single object and array responses
	var bookTickers []BookTicker

	// Try to unmarshal as array first
	if err := json.Unmarshal(body, &bookTickers); err != nil {
		// If that fails, try to unmarshal as single object
		var singleBookTicker BookTicker
		if err := json.Unmarshal(body, &singleBookTicker); err != nil {
			return Response[[]BookTicker]{}, err
		}
		bookTickers = []BookTicker{singleBookTicker}
	}

	return Response[[]BookTicker]{Code: 0, Message: "success", Data: &bookTickers}, nil
}

// GetAccountBalance queries account balance info (USER_DATA - signed endpoint).
func (c *Client) GetAccountBalance(ctx context.Context, req GetAccountBalanceRequest) (Response[[]AccountBalance], error) {
	params := map[string]string{}
	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}

	body, status, err := doSignedRequest(c.cfg, "GET", PathGetAccountBalance, params)
	if err != nil {
		return Response[[]AccountBalance]{}, err
	}
	if status != http.StatusOK {
		// For signed requests, check if the response contains an error message
		var errResp Response[[]AccountBalance]
		if json.Unmarshal(body, &errResp) == nil && errResp.Code != 0 {
			return errResp, fmt.Errorf("api error: %d - %s", errResp.Code, errResp.Message)
		}
		return Response[[]AccountBalance]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}

	var balances []AccountBalance
	if err := json.Unmarshal(body, &balances); err != nil {
		return Response[[]AccountBalance]{}, err
	}

	return Response[[]AccountBalance]{Code: 0, Message: "success", Data: &balances}, nil
}

// CreateOrder sends in a new order (TRADE - signed endpoint).
func (c *Client) CreateOrder(ctx context.Context, req CreateOrderRequest) (Response[CreateOrderResponse], error) {
	params := map[string]string{
		"symbol": req.Symbol,
		"side":   req.Side,
		"type":   req.Type,
	}

	// Add optional parameters
	if req.PositionSide != "" {
		params["positionSide"] = req.PositionSide
	}
	if req.TimeInForce != "" {
		params["timeInForce"] = req.TimeInForce
	}
	if req.Quantity != "" {
		params["quantity"] = req.Quantity
	}
	if req.ReduceOnly != "" {
		params["reduceOnly"] = req.ReduceOnly
	}
	if req.Price != "" {
		params["price"] = req.Price
	}
	if req.NewClientOrderId != "" {
		params["newClientOrderId"] = req.NewClientOrderId
	}
	if req.StopPrice != "" {
		params["stopPrice"] = req.StopPrice
	}
	if req.ClosePosition != "" {
		params["closePosition"] = req.ClosePosition
	}
	if req.ActivationPrice != "" {
		params["activationPrice"] = req.ActivationPrice
	}
	if req.CallbackRate != "" {
		params["callbackRate"] = req.CallbackRate
	}
	if req.WorkingType != "" {
		params["workingType"] = req.WorkingType
	}
	if req.PriceProtect != "" {
		params["priceProtect"] = req.PriceProtect
	}
	if req.NewOrderRespType != "" {
		params["newOrderRespType"] = req.NewOrderRespType
	}
	if req.PriceMatch != "" {
		params["priceMatch"] = req.PriceMatch
	}
	if req.SelfTradePreventionMode != "" {
		params["selfTradePreventionMode"] = req.SelfTradePreventionMode
	}
	if req.GoodTillDate > 0 {
		params["goodTillDate"] = fmt.Sprintf("%d", req.GoodTillDate)
	}
	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}

	body, status, err := doSignedRequest(c.cfg, "POST", PathCreateOrder, params)
	if err != nil {
		return Response[CreateOrderResponse]{}, err
	}
	if status != http.StatusOK && status != http.StatusCreated {
		// For signed requests, check if the response contains an error message
		var errResp Response[CreateOrderResponse]
		if json.Unmarshal(body, &errResp) == nil && errResp.Code != 0 {
			return errResp, fmt.Errorf("api error: %d - %s", errResp.Code, errResp.Message)
		}
		return Response[CreateOrderResponse]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}

	var order CreateOrderResponse
	if err := json.Unmarshal(body, &order); err != nil {
		return Response[CreateOrderResponse]{}, err
	}

	return Response[CreateOrderResponse]{Code: 0, Message: "success", Data: &order}, nil
}

// CancelOrder cancels an active order (TRADE - signed endpoint).
func (c *Client) CancelOrder(ctx context.Context, req CancelOrderRequest) (Response[CancelOrderResponse], error) {
	params := map[string]string{
		"symbol": req.Symbol,
	}

	// Either orderId or origClientOrderId must be sent
	if req.OrderId > 0 {
		params["orderId"] = fmt.Sprintf("%d", req.OrderId)
	}
	if req.OrigClientOrderId != "" {
		params["origClientOrderId"] = req.OrigClientOrderId
	}
	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}

	body, status, err := doSignedRequest(c.cfg, "DELETE", PathCancelOrder, params)
	if err != nil {
		return Response[CancelOrderResponse]{}, err
	}
	if status != http.StatusOK {
		// For signed requests, check if the response contains an error message
		var errResp Response[CancelOrderResponse]
		if json.Unmarshal(body, &errResp) == nil && errResp.Code != 0 {
			return errResp, fmt.Errorf("api error: %d - %s", errResp.Code, errResp.Message)
		}
		return Response[CancelOrderResponse]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}

	var order CancelOrderResponse
	if err := json.Unmarshal(body, &order); err != nil {
		return Response[CancelOrderResponse]{}, err
	}

	return Response[CancelOrderResponse]{Code: 0, Message: "success", Data: &order}, nil
}

// CancelAllOrders cancels all open orders for a symbol (TRADE - signed endpoint).
func (c *Client) CancelAllOrders(ctx context.Context, req CancelAllOrdersRequest) (Response[CancelAllOrdersResponse], error) {
	params := map[string]string{
		"symbol": req.Symbol,
	}

	if req.RecvWindow > 0 {
		params["recvWindow"] = fmt.Sprintf("%d", req.RecvWindow)
	}

	body, status, err := doSignedRequest(c.cfg, "DELETE", PathCancelAllOrders, params)
	if err != nil {
		return Response[CancelAllOrdersResponse]{}, err
	}
	if status != http.StatusOK {
		// For signed requests, check if the response contains an error message
		var errResp Response[CancelAllOrdersResponse]
		if json.Unmarshal(body, &errResp) == nil && errResp.Code != 0 {
			return errResp, fmt.Errorf("api error: %d - %s", errResp.Code, errResp.Message)
		}
		return Response[CancelAllOrdersResponse]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}

	var cancelResp CancelAllOrdersResponse
	if err := json.Unmarshal(body, &cancelResp); err != nil {
		return Response[CancelAllOrdersResponse]{}, err
	}

	return Response[CancelAllOrdersResponse]{Code: 0, Message: "success", Data: &cancelResp}, nil
}
