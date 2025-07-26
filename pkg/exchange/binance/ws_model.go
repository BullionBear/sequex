package binance

// WSKlineEvent represents the complete kline/candlestick WebSocket event
type WSKlineEvent struct {
	EventType string  `json:"e"` // Event type
	EventTime int64   `json:"E"` // Event time
	Symbol    string  `json:"s"` // Symbol
	KlineData WSKline `json:"k"` // Kline data
}

// WSKline represents the kline/candlestick data within the WebSocket event
type WSKline struct {
	StartTime                int64  `json:"t"` // Kline start time
	CloseTime                int64  `json:"T"` // Kline close time
	Symbol                   string `json:"s"` // Symbol
	Interval                 string `json:"i"` // Interval
	FirstTradeId             int64  `json:"f"` // First trade ID
	LastTradeId              int64  `json:"L"` // Last trade ID
	Open                     string `json:"o"` // Open price
	Close                    string `json:"c"` // Close price
	High                     string `json:"h"` // High price
	Low                      string `json:"l"` // Low price
	Volume                   string `json:"v"` // Base asset volume
	NumberOfTrades           int    `json:"n"` // Number of trades
	IsClosed                 bool   `json:"x"` // Is this kline closed?
	QuoteAssetVolume         string `json:"q"` // Quote asset volume
	TakerBuyBaseAssetVolume  string `json:"V"` // Taker buy base asset volume
	TakerBuyQuoteAssetVolume string `json:"Q"` // Taker buy quote asset volume
	Ignore                   string `json:"B"` // Ignore
}

// WSAggTradeEvent represents the complete aggregate trade WebSocket event
type WSAggTradeEvent struct {
	EventType    string `json:"e"` // Event type
	EventTime    int64  `json:"E"` // Event time
	Symbol       string `json:"s"` // Symbol
	AggTradeId   int64  `json:"a"` // Aggregate trade ID
	Price        string `json:"p"` // Price
	Quantity     string `json:"q"` // Quantity
	FirstTradeId int64  `json:"f"` // First trade ID
	LastTradeId  int64  `json:"l"` // Last trade ID
	TradeTime    int64  `json:"T"` // Trade time
	IsBuyerMaker bool   `json:"m"` // Is the buyer the market maker?
	Ignore       bool   `json:"M"` // Ignore
}

// WSAggTrade represents aggregate trade data (alias for event for consistency with kline pattern)
type WSAggTrade = WSAggTradeEvent

// WSTradeEvent represents the complete raw trade WebSocket event
type WSTradeEvent struct {
	EventType    string `json:"e"` // Event type
	EventTime    int64  `json:"E"` // Event time
	Symbol       string `json:"s"` // Symbol
	TradeId      int64  `json:"t"` // Trade ID
	Price        string `json:"p"` // Price
	Quantity     string `json:"q"` // Quantity
	TradeTime    int64  `json:"T"` // Trade time
	IsBuyerMaker bool   `json:"m"` // Is the buyer the market maker?
	Ignore       bool   `json:"M"` // Ignore
}

// WSTrade represents raw trade data (alias for event for consistency with other patterns)
type WSTrade = WSTradeEvent

// PriceLevel represents a single price level in the order book [price, quantity]
type PriceLevel [2]string

// WSDepthEvent represents the complete partial book depth WebSocket event
type WSDepthEvent struct {
	LastUpdateId int64        `json:"lastUpdateId"` // Last update ID
	Bids         []PriceLevel `json:"bids"`         // Bids to be updated [price, quantity]
	Asks         []PriceLevel `json:"asks"`         // Asks to be updated [price, quantity]
}

// WSDepth represents partial book depth data (alias for event for consistency with other patterns)
type WSDepth = WSDepthEvent

// WSDepthUpdateEvent represents the complete differential depth WebSocket event
type WSDepthUpdateEvent struct {
	EventType     string       `json:"e"` // Event type ("depthUpdate")
	EventTime     int64        `json:"E"` // Event time
	Symbol        string       `json:"s"` // Symbol
	FirstUpdateId int64        `json:"U"` // First update ID in event
	FinalUpdateId int64        `json:"u"` // Final update ID in event
	BidUpdates    []PriceLevel `json:"b"` // Bids to be updated [price, quantity]
	AskUpdates    []PriceLevel `json:"a"` // Asks to be updated [price, quantity]
}

// WSDepthUpdate represents differential depth data (alias for event for consistency with other patterns)
type WSDepthUpdate = WSDepthUpdateEvent

// KlineSubscriptionOptions defines the callback functions for kline subscription
type KlineSubscriptionOptions struct {
	OnConnect    func()              // Called when connection is established
	OnReconnect  func()              // Called when connection is reestablished
	OnError      func(err error)     // Called when an error occurs
	OnKline      func(kline WSKline) // Called when kline data is received
	OnDisconnect func()              // Called when connection is disconnected
}

// AggTradeSubscriptionOptions defines the callback functions for aggregate trade subscription
type AggTradeSubscriptionOptions struct {
	OnConnect    func()                    // Called when connection is established
	OnReconnect  func()                    // Called when connection is reestablished
	OnError      func(err error)           // Called when an error occurs
	OnAggTrade   func(aggTrade WSAggTrade) // Called when aggregate trade data is received
	OnDisconnect func()                    // Called when connection is disconnected
}

// TradeSubscriptionOptions defines the callback functions for raw trade subscription
type TradeSubscriptionOptions struct {
	OnConnect    func()              // Called when connection is established
	OnReconnect  func()              // Called when connection is reestablished
	OnError      func(err error)     // Called when an error occurs
	OnTrade      func(trade WSTrade) // Called when trade data is received
	OnDisconnect func()              // Called when connection is disconnected
}

// DepthSubscriptionOptions defines the callback functions for partial book depth subscription
type DepthSubscriptionOptions struct {
	OnConnect    func()              // Called when connection is established
	OnReconnect  func()              // Called when connection is reestablished
	OnError      func(err error)     // Called when an error occurs
	OnDepth      func(depth WSDepth) // Called when depth data is received
	OnDisconnect func()              // Called when connection is disconnected
}

// DepthUpdateSubscriptionOptions defines the callback functions for differential depth subscription
type DepthUpdateSubscriptionOptions struct {
	OnConnect     func()                     // Called when connection is established
	OnReconnect   func()                     // Called when connection is reestablished
	OnError       func(err error)            // Called when an error occurs
	OnDepthUpdate func(update WSDepthUpdate) // Called when depth update data is received
	OnDisconnect  func()                     // Called when connection is disconnected
}

// ConnectionState represents the current state of a WebSocket subscription
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateReconnecting
)

// Subscription represents an active WebSocket subscription
type Subscription struct {
	id      string
	conn    WSConnection
	options interface{} // Can be KlineSubscriptionOptions, AggTradeSubscriptionOptions, TradeSubscriptionOptions, DepthSubscriptionOptions, DepthUpdateSubscriptionOptions, or UserDataSubscriptionOptions
	state   ConnectionState
}

// User Data Stream Event Models

// WSOutboundAccountPositionEvent represents the outboundAccountPosition WebSocket event
type WSOutboundAccountPositionEvent struct {
	EventType             string           `json:"e"` // Event type
	EventTime             int64            `json:"E"` // Event Time
	LastAccountUpdateTime int64            `json:"u"` // Time of last account update
	BalanceArray          []AccountBalance `json:"B"` // Balances Array
}

// AccountBalance represents a balance entry in outboundAccountPosition event
type AccountBalance struct {
	Asset  string `json:"a"` // Asset
	Free   string `json:"f"` // Free
	Locked string `json:"l"` // Locked
}

// WSBalanceUpdateEvent represents the balanceUpdate WebSocket event
type WSBalanceUpdateEvent struct {
	EventType    string `json:"e"` // Event Type
	EventTime    int64  `json:"E"` // Event Time
	Asset        string `json:"a"` // Asset
	BalanceDelta string `json:"d"` // Balance Delta
	ClearTime    int64  `json:"T"` // Clear Time
}

// WSExecutionReportEvent represents the executionReport WebSocket event
type WSExecutionReportEvent struct {
	EventType                         string  `json:"e"` // Event type
	EventTime                         int64   `json:"E"` // Event time
	Symbol                            string  `json:"s"` // Symbol
	ClientOrderId                     string  `json:"c"` // Client order ID
	Side                              string  `json:"S"` // Side
	OrderType                         string  `json:"o"` // Order type
	TimeInForce                       string  `json:"f"` // Time in force
	OrderQuantity                     string  `json:"q"` // Order quantity
	OrderPrice                        string  `json:"p"` // Order price
	StopPrice                         string  `json:"P"` // Stop price
	IcebergQuantity                   string  `json:"F"` // Iceberg quantity
	OrderListId                       int64   `json:"g"` // OrderListId
	OriginalClientOrderId             string  `json:"C"` // Original client order ID
	CurrentExecutionType              string  `json:"x"` // Current execution type
	CurrentOrderStatus                string  `json:"X"` // Current order status
	OrderRejectReason                 string  `json:"r"` // Order reject reason
	OrderId                           int64   `json:"i"` // Order ID
	LastExecutedQuantity              string  `json:"l"` // Last executed quantity
	CumulativeFilledQuantity          string  `json:"z"` // Cumulative filled quantity
	LastExecutedPrice                 string  `json:"L"` // Last executed price
	CommissionAmount                  string  `json:"n"` // Commission amount
	CommissionAsset                   *string `json:"N"` // Commission asset
	TransactionTime                   int64   `json:"T"` // Transaction time
	TradeId                           int64   `json:"t"` // Trade ID
	PreventedMatchId                  int64   `json:"v"` // Prevented Match Id
	ExecutionId                       int64   `json:"I"` // Execution Id
	IsOrderOnBook                     bool    `json:"w"` // Is the order on the book?
	IsMakerSide                       bool    `json:"m"` // Is this trade the maker side?
	Ignore                            bool    `json:"M"` // Ignore
	OrderCreationTime                 int64   `json:"O"` // Order creation time
	CumulativeQuoteAssetTransactedQty string  `json:"Z"` // Cumulative quote asset transacted quantity
	LastQuoteAssetTransactedQty       string  `json:"Y"` // Last quote asset transacted quantity
	QuoteOrderQty                     string  `json:"Q"` // Quote Order Quantity
	WorkingTime                       int64   `json:"W"` // Working Time
	SelfTradePreventionMode           string  `json:"V"` // SelfTradePreventionMode
}

// WSListenKeyExpiredEvent represents the listenKeyExpired WebSocket event
type WSListenKeyExpiredEvent struct {
	EventType string `json:"e"`         // Event type
	EventTime int64  `json:"E"`         // Event time
	ListenKey string `json:"listenKey"` // Listen key that expired
}

// UserDataSubscriptionOptions defines the callback functions for user data subscription
type UserDataSubscriptionOptions struct {
	OnConnect          func()                                     // Called when connection is established
	OnReconnect        func()                                     // Called when connection is reestablished
	OnError            func(err error)                            // Called when an error occurs
	OnAccountPosition  func(event WSOutboundAccountPositionEvent) // Called when account position update is received
	OnBalanceUpdate    func(event WSBalanceUpdateEvent)           // Called when balance update is received
	OnExecutionReport  func(event WSExecutionReportEvent)         // Called when execution report is received
	OnListenKeyExpired func(event WSListenKeyExpiredEvent)        // Called when listen key expires (internal use)
	OnDisconnect       func()                                     // Called when connection is disconnected
}
