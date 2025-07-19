package binance

import (
	"time"
)

// ServerTimeResponse represents the response from /api/v3/time endpoint
type ServerTimeResponse struct {
	ServerTime int64 `json:"serverTime"`
}

// GetTime returns the server time as a Go time.Time
func (s *ServerTimeResponse) GetTime() time.Time {
	return time.Unix(0, s.ServerTime*int64(time.Millisecond))
}

// PingResponse represents the response from /api/v3/ping endpoint
type PingResponse struct{}

// ExchangeInfoResponse represents the response from /api/v3/exchangeInfo endpoint
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

// TickerPriceResponse represents the response from /api/v3/ticker/price endpoint
type TickerPriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// Ticker24hrResponse represents the response from /api/v3/ticker/24hr endpoint
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

// OrderBookResponse represents the response from /api/v3/depth endpoint
type OrderBookResponse struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// KlineResponse represents a single kline (candlestick) data
type KlineResponse []interface{}

// KlinesResponse represents the response from /api/v3/klines endpoint
type KlinesResponse []KlineResponse

// TradeResponse represents the response from /api/v3/trades endpoint
type TradeResponse struct {
	Id           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

// AccountResponse represents the response from /api/v3/account endpoint
type AccountResponse struct {
	MakerCommission  int64     `json:"makerCommission"`
	TakerCommission  int64     `json:"takerCommission"`
	BuyerCommission  int64     `json:"buyerCommission"`
	SellerCommission int64     `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	UpdateTime       int64     `json:"updateTime"`
	AccountType      string    `json:"accountType"`
	Balances         []Balance `json:"balances"`
	Permissions      []string  `json:"permissions"`
}

// Balance represents an account balance
type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// NewOrderRequest represents a request to place a new order
type NewOrderRequest struct {
	Symbol           string `json:"symbol"`
	Side             string `json:"side"`
	Type             string `json:"type"`
	TimeInForce      string `json:"timeInForce,omitempty"`
	Quantity         string `json:"quantity,omitempty"`
	QuoteOrderQty    string `json:"quoteOrderQty,omitempty"`
	Price            string `json:"price,omitempty"`
	NewClientOrderId string `json:"newClientOrderId,omitempty"`
	StopPrice        string `json:"stopPrice,omitempty"`
	IcebergQty       string `json:"icebergQty,omitempty"`
	NewOrderRespType string `json:"newOrderRespType,omitempty"`
	RecvWindow       int64  `json:"recvWindow,omitempty"`
	Timestamp        int64  `json:"timestamp"`
}

// NewOrderResponse represents the response from placing a new order
type NewOrderResponse struct {
	Symbol              string `json:"symbol"`
	OrderId             int64  `json:"orderId"`
	OrderListId         int64  `json:"orderListId"`
	ClientOrderId       string `json:"clientOrderId"`
	TransactTime        int64  `json:"transactTime"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	Fills               []Fill `json:"fills,omitempty"`
}

// Fill represents a trade fill
type Fill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	TradeId         int64  `json:"tradeId"`
}

// OrderResponse represents the response from querying an order
type OrderResponse struct {
	Symbol              string `json:"symbol"`
	OrderId             int64  `json:"orderId"`
	OrderListId         int64  `json:"orderListId"`
	ClientOrderId       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
}

// CancelOrderResponse represents the response from cancelling an order
type CancelOrderResponse struct {
	Symbol              string `json:"symbol"`
	OrigClientOrderId   string `json:"origClientOrderId"`
	OrderId             int64  `json:"orderId"`
	OrderListId         int64  `json:"orderListId"`
	ClientOrderId       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
}
