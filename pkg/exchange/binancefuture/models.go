package binancefuture

import (
	"encoding/json"
	"fmt"
	"time"
)

// ServerTimeResponse represents the response from /fapi/v1/time endpoint
type ServerTimeResponse struct {
	ServerTime int64 `json:"serverTime"`
}

// GetTime returns the server time as a Go time.Time
func (s *ServerTimeResponse) GetTime() time.Time {
	return time.Unix(0, s.ServerTime*int64(time.Millisecond))
}

// PingResponse represents the response from /fapi/v1/ping endpoint
type PingResponse struct{}

// ExchangeInfoResponse represents the response from /fapi/v1/exchangeInfo endpoint
type ExchangeInfoResponse struct {
	Timezone   string       `json:"timezone"`
	ServerTime int64        `json:"serverTime"`
	RateLimits []RateLimit  `json:"rateLimits"`
	Symbols    []SymbolInfo `json:"symbols"`
}

// RateLimit represents rate limiting information
type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

// SymbolInfo represents trading symbol information
type SymbolInfo struct {
	Symbol                 string         `json:"symbol"`
	Status                 string         `json:"status"`
	BaseAsset              string         `json:"baseAsset"`
	BaseAssetPrecision     int            `json:"baseAssetPrecision"`
	QuoteAsset             string         `json:"quoteAsset"`
	QuoteAssetPrecision    int            `json:"quoteAssetPrecision"`
	OrderTypes             []string       `json:"orderTypes"`
	IcebergAllowed         bool           `json:"icebergAllowed"`
	OcoAllowed             bool           `json:"ocoAllowed"`
	IsSpotTradingAllowed   bool           `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed bool           `json:"isMarginTradingAllowed"`
	Filters                []SymbolFilter `json:"filters"`
	Permissions            []string       `json:"permissions"`
}

// SymbolFilter represents symbol trading filters
type SymbolFilter struct {
	FilterType          string `json:"filterType"`
	MinPrice            string `json:"minPrice,omitempty"`
	MaxPrice            string `json:"maxPrice,omitempty"`
	TickSize            string `json:"tickSize,omitempty"`
	MinQty              string `json:"minQty,omitempty"`
	MaxQty              string `json:"maxQty,omitempty"`
	StepSize            string `json:"stepSize,omitempty"`
	MinNotional         string `json:"minNotional,omitempty"`
	ApplyToMarket       bool   `json:"applyToMarket,omitempty"`
	AvgPriceMins        int    `json:"avgPriceMins,omitempty"`
	Limit               int    `json:"limit,omitempty"`
	MaxNumAlgoOrders    int    `json:"maxNumAlgoOrders,omitempty"`
	MaxNumOrders        int    `json:"maxNumOrders,omitempty"`
	MaxNumIcebergOrders int    `json:"maxNumIcebergOrders,omitempty"`
}

// TickerPriceResponse represents the response from /fapi/v1/ticker/price endpoint
type TickerPriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// Ticker24hrResponse represents the response from /fapi/v1/ticker/24hr endpoint
type Ticker24hrResponse struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstId            int64  `json:"firstId"`
	LastId             int64  `json:"lastId"`
	Count              int64  `json:"count"`
}

// OrderBookResponse represents the response from /fapi/v1/depth endpoint
type OrderBookResponse struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// KlineData represents a single kline (candlestick) data point
// Binance returns klines as arrays, so we need custom unmarshaling
type KlineData struct {
	OpenTime                 int64
	Open                     string
	High                     string
	Low                      string
	Close                    string
	Volume                   string
	CloseTime                int64
	QuoteAssetVolume         string
	NumberOfTrades           int64
	TakerBuyBaseAssetVolume  string
	TakerBuyQuoteAssetVolume string
}

// UnmarshalJSON implements custom JSON unmarshaling for KlineData
func (k *KlineData) UnmarshalJSON(data []byte) error {
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	if len(arr) < 11 {
		return fmt.Errorf("kline data must have at least 11 elements, got %d", len(arr))
	}

	// Parse each field from the array
	if openTime, ok := arr[0].(float64); ok {
		k.OpenTime = int64(openTime)
	}
	if open, ok := arr[1].(string); ok {
		k.Open = open
	}
	if high, ok := arr[2].(string); ok {
		k.High = high
	}
	if low, ok := arr[3].(string); ok {
		k.Low = low
	}
	if close, ok := arr[4].(string); ok {
		k.Close = close
	}
	if volume, ok := arr[5].(string); ok {
		k.Volume = volume
	}
	if closeTime, ok := arr[6].(float64); ok {
		k.CloseTime = int64(closeTime)
	}
	if quoteAssetVolume, ok := arr[7].(string); ok {
		k.QuoteAssetVolume = quoteAssetVolume
	}
	if numberOfTrades, ok := arr[8].(float64); ok {
		k.NumberOfTrades = int64(numberOfTrades)
	}
	if takerBuyBaseAssetVolume, ok := arr[9].(string); ok {
		k.TakerBuyBaseAssetVolume = takerBuyBaseAssetVolume
	}
	if takerBuyQuoteAssetVolume, ok := arr[10].(string); ok {
		k.TakerBuyQuoteAssetVolume = takerBuyQuoteAssetVolume
	}

	return nil
}

// KlineResponse represents the response from /fapi/v1/klines endpoint
type KlineResponse []KlineData

// TradeResponse represents the response from /fapi/v1/trades endpoint
type TradeResponse struct {
	Id           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

// MarkPriceResponse represents the response from /fapi/v1/premiumIndex endpoint
type MarkPriceResponse struct {
	Symbol               string `json:"symbol"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	InterestRate         string `json:"interestRate"`
	Time                 int64  `json:"time"`
}

// FundingRateResponse represents the response from /fapi/v1/fundingRate endpoint
type FundingRateResponse struct {
	Symbol      string `json:"symbol"`
	FundingRate string `json:"fundingRate"`
	FundingTime int64  `json:"fundingTime"`
}

// OpenInterestResponse represents the response from /fapi/v1/openInterest endpoint
type OpenInterestResponse struct {
	OpenInterest string `json:"openInterest"`
	Symbol       string `json:"symbol"`
	Time         int64  `json:"time"`
}

// TickerPriceResult struct for handling both single and array responses
type TickerPriceResult struct {
	Single *TickerPriceResponse
	Array  []TickerPriceResponse
}

// Ticker24hrResult struct for handling both single and array responses
type Ticker24hrResult struct {
	Single *Ticker24hrResponse
	Array  []Ticker24hrResponse
}

// AccountResponse represents the response from /fapi/v2/account endpoint
type AccountResponse struct {
	FeeTier                     int64      `json:"feeTier"`
	CanTrade                    bool       `json:"canTrade"`
	CanDeposit                  bool       `json:"canDeposit"`
	CanWithdraw                 bool       `json:"canWithdraw"`
	UpdateTime                  int64      `json:"updateTime"`
	AccountType                 string     `json:"accountType"`
	TotalWalletBalance          string     `json:"totalWalletBalance"`
	TotalUnrealizedProfit       string     `json:"totalUnrealizedProfit"`
	TotalMarginBalance          string     `json:"totalMarginBalance"`
	TotalPositionInitialMargin  string     `json:"totalPositionInitialMargin"`
	TotalOpenOrderInitialMargin string     `json:"totalOpenOrderInitialMargin"`
	TotalCrossWalletBalance     string     `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl             string     `json:"totalCrossUnPnl"`
	AvailableBalance            string     `json:"availableBalance"`
	MaxWithdrawAmount           string     `json:"maxWithdrawAmount"`
	Assets                      []Asset    `json:"assets"`
	Positions                   []Position `json:"positions"`
}

// Asset represents account asset information
type Asset struct {
	Asset                  string `json:"asset"`
	WalletBalance          string `json:"walletBalance"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	MarginBalance          string `json:"marginBalance"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	AvailableBalance       string `json:"availableBalance"`
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
}

// Position represents position information
type Position struct {
	Symbol                 string `json:"symbol"`
	InitialMargin          string `json:"initialMargin"`
	MaintMargin            string `json:"maintMargin"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	Leverage               string `json:"leverage"`
	Isolated               bool   `json:"isolated"`
	EntryPrice             string `json:"entryPrice"`
	MaxNotional            string `json:"maxNotional"`
	BidNotional            string `json:"bidNotional"`
	AskNotional            string `json:"askNotional"`
	MarkPrice              string `json:"markPrice"`
	PositionAmt            string `json:"positionAmt"`
	PositionSide           string `json:"positionSide"`
	UpdateTime             int64  `json:"updateTime"`
}

// NewOrderRequest represents a new order request
type NewOrderRequest struct {
	Symbol           string `json:"symbol"`
	Side             string `json:"side"`
	PositionSide     string `json:"positionSide,omitempty"`
	Type             string `json:"type"`
	TimeInForce      string `json:"timeInForce,omitempty"`
	Quantity         string `json:"quantity,omitempty"`
	ReduceOnly       bool   `json:"reduceOnly,omitempty"`
	Price            string `json:"price,omitempty"`
	NewClientOrderId string `json:"newClientOrderId,omitempty"`
	StopPrice        string `json:"stopPrice,omitempty"`
	WorkingType      string `json:"workingType,omitempty"`
	PriceProtect     bool   `json:"priceProtect,omitempty"`
	NewOrderRespType string `json:"newOrderRespType,omitempty"`
	ClosePosition    bool   `json:"closePosition,omitempty"`
	ActivationPrice  string `json:"activationPrice,omitempty"`
	CallbackRate     string `json:"callbackRate,omitempty"`
	RecvWindow       int64  `json:"recvWindow,omitempty"`
	Timestamp        int64  `json:"timestamp"`
}

// NewOrderResponse represents the response from placing a new order
type NewOrderResponse struct {
	ClientOrderId       string `json:"clientOrderId"`
	CumQty              string `json:"cumQty"`
	CumQuote            string `json:"cumQuote"`
	ExecutedQty         string `json:"executedQty"`
	OrderId             int64  `json:"orderId"`
	AvgPrice            string `json:"avgPrice"`
	OrigQty             string `json:"origQty"`
	Price               string `json:"price"`
	ReduceOnly          bool   `json:"reduceOnly"`
	Side                string `json:"side"`
	PositionSide        string `json:"positionSide"`
	Status              string `json:"status"`
	StopPrice           string `json:"stopPrice"`
	ClosePosition       bool   `json:"closePosition"`
	Symbol              string `json:"symbol"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	OrigType            string `json:"origType"`
	ActivatePrice       string `json:"activatePrice"`
	PriceRate           string `json:"priceRate"`
	UpdateTime          int64  `json:"updateTime"`
	WorkingType         string `json:"workingType"`
	PriceProtect        bool   `json:"priceProtect"`
	PriceMatch          string `json:"priceMatch"`
	SelfTradePrevention string `json:"selfTradePrevention"`
	GoodTillDate        int64  `json:"goodTillDate"`
	Fills               []Fill `json:"fills,omitempty"`
}

// Fill represents order fill information
type Fill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	TradeId         int64  `json:"tradeId"`
}

// OrderResponse represents the response from getting order information
type OrderResponse struct {
	ClientOrderId       string `json:"clientOrderId"`
	CumQty              string `json:"cumQty"`
	CumQuote            string `json:"cumQuote"`
	ExecutedQty         string `json:"executedQty"`
	OrderId             int64  `json:"orderId"`
	AvgPrice            string `json:"avgPrice"`
	OrigQty             string `json:"origQty"`
	Price               string `json:"price"`
	ReduceOnly          bool   `json:"reduceOnly"`
	Side                string `json:"side"`
	PositionSide        string `json:"positionSide"`
	Status              string `json:"status"`
	StopPrice           string `json:"stopPrice"`
	ClosePosition       bool   `json:"closePosition"`
	Symbol              string `json:"symbol"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	OrigType            string `json:"origType"`
	ActivatePrice       string `json:"activatePrice"`
	PriceRate           string `json:"priceRate"`
	UpdateTime          int64  `json:"updateTime"`
	WorkingType         string `json:"workingType"`
	PriceProtect        bool   `json:"priceProtect"`
	PriceMatch          string `json:"priceMatch"`
	SelfTradePrevention string `json:"selfTradePrevention"`
	GoodTillDate        int64  `json:"goodTillDate"`
}

// CancelOrderResponse represents the response from canceling an order
type CancelOrderResponse struct {
	ClientOrderId       string `json:"clientOrderId"`
	CumQty              string `json:"cumQty"`
	CumQuote            string `json:"cumQuote"`
	ExecutedQty         string `json:"executedQty"`
	OrderId             int64  `json:"orderId"`
	OrigQty             string `json:"origQty"`
	Price               string `json:"price"`
	ReduceOnly          bool   `json:"reduceOnly"`
	Side                string `json:"side"`
	PositionSide        string `json:"positionSide"`
	Status              string `json:"status"`
	StopPrice           string `json:"stopPrice"`
	ClosePosition       bool   `json:"closePosition"`
	Symbol              string `json:"symbol"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	OrigType            string `json:"origType"`
	ActivatePrice       string `json:"activatePrice"`
	PriceRate           string `json:"priceRate"`
	UpdateTime          int64  `json:"updateTime"`
	WorkingType         string `json:"workingType"`
	PriceProtect        bool   `json:"priceProtect"`
	PriceMatch          string `json:"priceMatch"`
	SelfTradePrevention string `json:"selfTradePrevention"`
	GoodTillDate        int64  `json:"goodTillDate"`
}

// UserTradeResponse represents the response from getting user trades
type UserTradeResponse struct {
	Symbol          string `json:"symbol"`
	Id              int64  `json:"id"`
	OrderId         int64  `json:"orderId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	RealizedPnl     string `json:"realizedPnl"`
	Side            string `json:"side"`
	PositionSide    string `json:"positionSide"`
	Fee             string `json:"fee"`
	FeeAsset        string `json:"feeAsset"`
	Time            int64  `json:"time"`
	MatchingOrderId int64  `json:"matchingOrderId"`
	Maker           bool   `json:"maker"`
}

// PositionRiskResponse represents the response from /fapi/v3/positionRisk endpoint
type PositionRiskResponse struct {
	EntryPrice             string `json:"entryPrice"`
	MarginType             string `json:"marginType"`
	IsAutoAddMargin        string `json:"isAutoAddMargin"`
	IsolatedMargin         string `json:"isolatedMargin"`
	Leverage               string `json:"leverage"`
	LiquidationPrice       string `json:"liquidationPrice"`
	MarkPrice              string `json:"markPrice"`
	MaxNotionalValue       string `json:"maxNotionalValue"`
	NetUnrealizedPnl       string `json:"netUnrealizedPnl"`
	Notional               string `json:"notional"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	PositionAmt            string `json:"positionAmt"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	PositionSide           string `json:"positionSide"`
	Symbol                 string `json:"symbol"`
	UnrealizedPnl          string `json:"unrealizedPnl"`
	UpdateTime             int64  `json:"updateTime"`
}

// HasPosition returns true if the position has a non-zero amount
func (p *PositionRiskResponse) HasPosition() bool {
	return p.PositionAmt != "0" && p.PositionAmt != ""
}

// IsLong returns true if the position is long
func (p *PositionRiskResponse) IsLong() bool {
	return p.HasPosition() && p.PositionSide == PositionSideLong
}

// IsShort returns true if the position is short
func (p *PositionRiskResponse) IsShort() bool {
	return p.HasPosition() && p.PositionSide == PositionSideShort
}

// GetUpdateTime returns the update time as a Go time.Time
func (p *PositionRiskResponse) GetUpdateTime() time.Time {
	return time.Unix(0, p.UpdateTime*int64(time.Millisecond))
}

// PositionSideResponse represents the response from /fapi/v1/positionSide/dual endpoint
type PositionSideResponse struct {
	DualSidePosition bool `json:"dualSidePosition"`
}

// LeverageResponse represents the response from /fapi/v1/leverage endpoint
type LeverageResponse struct {
	Leverage         int    `json:"leverage"`
	MaxNotionalValue string `json:"maxNotionalValue"`
	Symbol           string `json:"symbol"`
}

// Helper methods for TickerPriceResult
func (t *TickerPriceResult) IsSingle() bool {
	return t.Single != nil
}

func (t *TickerPriceResult) IsArray() bool {
	return t.Array != nil
}

func (t *TickerPriceResult) GetSingle() *TickerPriceResponse {
	return t.Single
}

func (t *TickerPriceResult) GetArray() []TickerPriceResponse {
	return t.Array
}

// Helper methods for Ticker24hrResult
func (t *Ticker24hrResult) IsSingle() bool {
	return t.Single != nil
}

func (t *Ticker24hrResult) IsArray() bool {
	return t.Array != nil
}

func (t *Ticker24hrResult) GetSingle() *Ticker24hrResponse {
	return t.Single
}

func (t *Ticker24hrResult) GetArray() []Ticker24hrResponse {
	return t.Array
}

// UserDataStreamResponse represents the response for user data stream operations
type UserDataStreamResponse struct {
	ListenKey string `json:"listenKey"`
}

// UserDataStreamRequest represents the request for user data stream operations
type UserDataStreamRequest struct {
	ListenKey string `json:"listenKey,omitempty"`
}
