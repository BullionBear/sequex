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
