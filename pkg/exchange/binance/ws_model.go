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
	conn    *BinanceWSConn
	options interface{} // Can be KlineSubscriptionOptions, AggTradeSubscriptionOptions, TradeSubscriptionOptions, DepthSubscriptionOptions, or DepthUpdateSubscriptionOptions
	state   ConnectionState
}
