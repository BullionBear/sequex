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
