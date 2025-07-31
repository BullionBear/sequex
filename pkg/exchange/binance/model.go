package binance

// Response is the unified response wrapper for all endpoints.
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    *T     `json:"data,omitempty"`
}

// CreateOrderRequest defines the parameters for creating a new order.
type CreateOrderRequest struct {
	Symbol                  string // required
	Side                    string // required (BUY/SELL)
	Type                    string // required (LIMIT/MARKET/etc)
	TimeInForce             string // optional
	Quantity                string // optional
	QuoteOrderQty           string // optional
	Price                   string // optional
	NewClientOrderId        string // optional
	StrategyId              int64  // optional
	StrategyType            int    // optional
	StopPrice               string // optional
	TrailingDelta           int64  // optional
	IcebergQty              string // optional
	NewOrderRespType        string // optional (ACK/RESULT/FULL)
	SelfTradePreventionMode string // optional
	RecvWindow              int64  // optional
}

// CreateOrderResponse is the unified order response (FULL type, superset of all response types).
type CreateOrderResponse struct {
	Symbol                  string      `json:"symbol"`
	OrderId                 int64       `json:"orderId"`
	OrderListId             int64       `json:"orderListId"`
	ClientOrderId           string      `json:"clientOrderId"`
	TransactTime            int64       `json:"transactTime"`
	Price                   string      `json:"price"`
	OrigQty                 string      `json:"origQty"`
	ExecutedQty             string      `json:"executedQty"`
	CummulativeQuoteQty     string      `json:"cummulativeQuoteQty"`
	Status                  string      `json:"status"`
	TimeInForce             string      `json:"timeInForce"`
	Type                    string      `json:"type"`
	Side                    string      `json:"side"`
	WorkingTime             int64       `json:"workingTime"`
	SelfTradePreventionMode string      `json:"selfTradePreventionMode"`
	OrigQuoteOrderQty       string      `json:"origQuoteOrderQty"`
	Fills                   []OrderFill `json:"fills,omitempty"`
	StopPrice               string      `json:"stopPrice,omitempty"`
	IcebergQty              string      `json:"icebergQty,omitempty"`
	PreventedMatchId        int64       `json:"preventedMatchId,omitempty"`
	PreventedQuantity       string      `json:"preventedQuantity,omitempty"`
	StrategyId              int64       `json:"strategyId,omitempty"`
	StrategyType            int         `json:"strategyType,omitempty"`
	TrailingDelta           int64       `json:"trailingDelta,omitempty"`
	TrailingTime            int64       `json:"trailingTime,omitempty"`
	UsedSor                 bool        `json:"usedSor,omitempty"`
	WorkingFloor            string      `json:"workingFloor,omitempty"`
}

type OrderFill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	TradeId         int64  `json:"tradeId"`
}

// OrderBookDepthResponse models the response for the /api/v3/depth endpoint.
type OrderBookDepthResponse struct {
	LastUpdateId int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// RecentTrade models a single trade in the /api/v3/trades response.
type RecentTrade struct {
	ID           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

// AggTrade models a single aggregate trade in the /api/v3/aggTrades response.
type AggTrade struct {
	AggTradeId   int64  `json:"a"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	FirstTradeId int64  `json:"f"`
	LastTradeId  int64  `json:"l"`
	Timestamp    int64  `json:"T"`
	IsBuyerMaker bool   `json:"m"`
	IsBestMatch  bool   `json:"M"`
}

// Kline models a single candlestick/kline in the /api/v3/klines response.
type Kline struct {
	OpenTime                 int64  `json:"openTime"`
	Open                     string `json:"open"`
	High                     string `json:"high"`
	Low                      string `json:"low"`
	Close                    string `json:"close"`
	Volume                   string `json:"volume"`
	CloseTime                int64  `json:"closeTime"`
	QuoteAssetVolume         string `json:"quoteAssetVolume"`
	NumberOfTrades           int    `json:"numberOfTrades"`
	TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume"`
	Ignore                   string `json:"ignore"`
}

// PriceTicker models a single price ticker in the /api/v3/ticker/price response.
type PriceTicker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// CancelOrderRequest models the request for cancelling an order.
type CancelOrderRequest struct {
	Symbol             string
	OrderId            int64
	OrigClientOrderId  string
	NewClientOrderId   string
	CancelRestrictions string
	RecvWindow         int64
}

// CancelOrderResponse models the response for cancelling an order.
type CancelOrderResponse struct {
	Symbol                  string `json:"symbol"`
	OrigClientOrderId       string `json:"origClientOrderId"`
	OrderId                 int64  `json:"orderId"`
	OrderListId             int64  `json:"orderListId"`
	ClientOrderId           string `json:"clientOrderId"`
	TransactTime            int64  `json:"transactTime"`
	Price                   string `json:"price"`
	OrigQty                 string `json:"origQty"`
	ExecutedQty             string `json:"executedQty"`
	OrigQuoteOrderQty       string `json:"origQuoteOrderQty"`
	CummulativeQuoteQty     string `json:"cummulativeQuoteQty"`
	Status                  string `json:"status"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	Side                    string `json:"side"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
}

// CancelAllOrdersRequest models the request for cancelling all open orders on a symbol.
type CancelAllOrdersRequest struct {
	Symbol     string
	RecvWindow int64
}

// QueryOrderRequest models the request for querying an order's status.
type QueryOrderRequest struct {
	Symbol            string
	OrderId           int64
	OrigClientOrderId string
	RecvWindow        int64
}

// QueryOrderResponse models the response for querying an order's status.
type QueryOrderResponse struct {
	Symbol                  string `json:"symbol"`
	OrderId                 int64  `json:"orderId"`
	OrderListId             int64  `json:"orderListId"`
	ClientOrderId           string `json:"clientOrderId"`
	Price                   string `json:"price"`
	OrigQty                 string `json:"origQty"`
	ExecutedQty             string `json:"executedQty"`
	CummulativeQuoteQty     string `json:"cummulativeQuoteQty"`
	Status                  string `json:"status"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	Side                    string `json:"side"`
	StopPrice               string `json:"stopPrice"`
	IcebergQty              string `json:"icebergQty"`
	Time                    int64  `json:"time"`
	UpdateTime              int64  `json:"updateTime"`
	IsWorking               bool   `json:"isWorking"`
	WorkingTime             int64  `json:"workingTime"`
	OrigQuoteOrderQty       string `json:"origQuoteOrderQty"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
}

// ListOpenOrdersRequest models the request for listing open orders.
type ListOpenOrdersRequest struct {
	Symbol     string
	RecvWindow int64
}

// GetAccountInfoRequest models the request for getting account information.
type GetAccountInfoRequest struct {
	OmitZeroBalances bool
	RecvWindow       int64
}

// CommissionRates models the commission rates in the account info response.
type CommissionRates struct {
	Maker  string `json:"maker"`
	Taker  string `json:"taker"`
	Buyer  string `json:"buyer"`
	Seller string `json:"seller"`
}

// Balance models a single asset balance in the account info response.
type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// GetAccountInfoResponse models the response for getting account information.
type GetAccountInfoResponse struct {
	MakerCommission            int             `json:"makerCommission"`
	TakerCommission            int             `json:"takerCommission"`
	BuyerCommission            int             `json:"buyerCommission"`
	SellerCommission           int             `json:"sellerCommission"`
	CommissionRates            CommissionRates `json:"commissionRates"`
	CanTrade                   bool            `json:"canTrade"`
	CanWithdraw                bool            `json:"canWithdraw"`
	CanDeposit                 bool            `json:"canDeposit"`
	Brokered                   bool            `json:"brokered"`
	RequireSelfTradePrevention bool            `json:"requireSelfTradePrevention"`
	PreventSor                 bool            `json:"preventSor"`
	UpdateTime                 int64           `json:"updateTime"`
	AccountType                string          `json:"accountType"`
	Balances                   []Balance       `json:"balances"`
	Permissions                []string        `json:"permissions"`
	Uid                        int64           `json:"uid"`
}

// GetAccountTradesRequest models the request for getting account trades.
type GetAccountTradesRequest struct {
	Symbol     string
	OrderId    int64
	StartTime  int64
	EndTime    int64
	FromId     int64
	Limit      int
	RecvWindow int64
}

// AccountTrade models a single trade in the account trade list response.
type AccountTrade struct {
	Symbol          string `json:"symbol"`
	Id              int64  `json:"id"`
	OrderId         int64  `json:"orderId"`
	OrderListId     int64  `json:"orderListId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	QuoteQty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
}

// UserDataStreamResponse models the response for starting a user data stream.
type UserDataStreamResponse struct {
	ListenKey string `json:"listenKey"`
}

// EmptyResponse models empty responses for keepalive and close stream operations.
type EmptyResponse struct{}

// ExchangeInfoRequest models the request for getting exchange information.
type ExchangeInfoRequest struct {
	Symbol             string
	Symbols            []string
	Permissions        []string
	ShowPermissionSets bool
	SymbolStatus       string
}

// RateLimit models a single rate limit in the exchange info response.
type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

// Filter models a single filter in the exchange info response.
type Filter struct {
	FilterType string `json:"filterType"`
	// PRICE_FILTER fields
	MinPrice string `json:"minPrice,omitempty"`
	MaxPrice string `json:"maxPrice,omitempty"`
	TickSize string `json:"tickSize,omitempty"`
	// PERCENT_PRICE fields
	MultiplierUp   string `json:"multiplierUp,omitempty"`
	MultiplierDown string `json:"multiplierDown,omitempty"`
	AvgPriceMins   int    `json:"avgPriceMins,omitempty"`
	// PERCENT_PRICE_BY_SIDE fields
	BidMultiplierUp   string `json:"bidMultiplierUp,omitempty"`
	BidMultiplierDown string `json:"bidMultiplierDown,omitempty"`
	AskMultiplierUp   string `json:"askMultiplierUp,omitempty"`
	AskMultiplierDown string `json:"askMultiplierDown,omitempty"`
	// LOT_SIZE fields
	MinQty   string `json:"minQty,omitempty"`
	MaxQty   string `json:"maxQty,omitempty"`
	StepSize string `json:"stepSize,omitempty"`
	// MIN_NOTIONAL fields
	MinNotional   string `json:"minNotional,omitempty"`
	ApplyToMarket bool   `json:"applyToMarket,omitempty"`
	// NOTIONAL fields
	ApplyMinToMarket bool   `json:"applyMinToMarket,omitempty"`
	ApplyMaxToMarket bool   `json:"applyMaxToMarket,omitempty"`
	MaxNotional      string `json:"maxNotional,omitempty"`
	// ICEBERG_PARTS fields
	Limit int `json:"limit,omitempty"`
	// MAX_NUM_ORDERS fields
	MaxNumOrders int `json:"maxNumOrders,omitempty"`
	// MAX_NUM_ALGO_ORDERS fields
	MaxNumAlgoOrders int `json:"maxNumAlgoOrders,omitempty"`
	// MAX_NUM_ICEBERG_ORDERS fields
	MaxNumIcebergOrders int `json:"maxNumIcebergOrders,omitempty"`
	// MAX_POSITION fields
	MaxPosition string `json:"maxPosition,omitempty"`
	// TRAILING_DELTA fields
	MinTrailingAboveDelta int `json:"minTrailingAboveDelta,omitempty"`
	MaxTrailingAboveDelta int `json:"maxTrailingAboveDelta,omitempty"`
	MinTrailingBelowDelta int `json:"minTrailingBelowDelta,omitempty"`
	MaxTrailingBelowDelta int `json:"maxTrailingBelowDelta,omitempty"`
}

// Symbol models a single symbol in the exchange info response.
type Symbol struct {
	Symbol                          string     `json:"symbol"`
	Status                          string     `json:"status"`
	BaseAsset                       string     `json:"baseAsset"`
	BaseAssetPrecision              int        `json:"baseAssetPrecision"`
	QuoteAsset                      string     `json:"quoteAsset"`
	QuotePrecision                  int        `json:"quotePrecision"`
	QuoteAssetPrecision             int        `json:"quoteAssetPrecision"`
	BaseCommissionPrecision         int        `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision        int        `json:"quoteCommissionPrecision"`
	OrderTypes                      []string   `json:"orderTypes"`
	IcebergAllowed                  bool       `json:"icebergAllowed"`
	OcoAllowed                      bool       `json:"ocoAllowed"`
	OtoAllowed                      bool       `json:"otoAllowed"`
	QuoteOrderQtyMarketAllowed      bool       `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop               bool       `json:"allowTrailingStop"`
	CancelReplaceAllowed            bool       `json:"cancelReplaceAllowed"`
	AmendAllowed                    bool       `json:"amendAllowed"`
	IsSpotTradingAllowed            bool       `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed          bool       `json:"isMarginTradingAllowed"`
	Filters                         []Filter   `json:"filters"`
	Permissions                     []string   `json:"permissions"`
	PermissionSets                  [][]string `json:"permissionSets"`
	DefaultSelfTradePreventionMode  string     `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string   `json:"allowedSelfTradePreventionModes"`
}

// SOR models a single SOR (Smart Order Router) in the exchange info response.
type SOR struct {
	BaseAsset string   `json:"baseAsset"`
	Symbols   []string `json:"symbols"`
}

// ExchangeInfoResponse models the response for getting exchange information.
type ExchangeInfoResponse struct {
	Timezone        string      `json:"timezone"`
	ServerTime      int64       `json:"serverTime"`
	RateLimits      []RateLimit `json:"rateLimits"`
	ExchangeFilters []Filter    `json:"exchangeFilters"`
	Symbols         []Symbol    `json:"symbols"`
	Sors            []SOR       `json:"sors,omitempty"`
}
