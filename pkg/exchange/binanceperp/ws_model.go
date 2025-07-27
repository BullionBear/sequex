package binanceperp

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

// KlineSubscriptionOptions defines the callback functions for kline subscription
type KlineSubscriptionOptions struct {
	onConnect    func()              // Called when connection is established
	onReconnect  func()              // Called when connection is reestablished
	onError      func(err error)     // Called when an error occurs
	onKline      func(kline WSKline) // Called when kline data is received
	onDisconnect func()              // Called when connection is disconnected
}

// WithConnect sets the OnConnect callback using chain method
func (k *KlineSubscriptionOptions) WithConnect(onConnect func()) *KlineSubscriptionOptions {
	k.onConnect = onConnect
	return k
}

// WithReconnect sets the OnReconnect callback using chain method
func (k *KlineSubscriptionOptions) WithReconnect(onReconnect func()) *KlineSubscriptionOptions {
	k.onReconnect = onReconnect
	return k
}

// WithError sets the OnError callback using chain method
func (k *KlineSubscriptionOptions) WithError(onError func(error)) *KlineSubscriptionOptions {
	k.onError = onError
	return k
}

// WithKline sets the OnKline callback using chain method
func (k *KlineSubscriptionOptions) WithKline(onKline func(WSKline)) *KlineSubscriptionOptions {
	k.onKline = onKline
	return k
}

// WithDisconnect sets the OnDisconnect callback using chain method
func (k *KlineSubscriptionOptions) WithDisconnect(onDisconnect func()) *KlineSubscriptionOptions {
	k.onDisconnect = onDisconnect
	return k
}

// ConnectionState represents the current state of a WebSocket subscription
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateReconnecting
)

// WSAggTradeEvent represents the complete aggregate trade WebSocket event
type WSAggTradeEvent struct {
	EventType    string `json:"e"` // Event type
	EventTime    int64  `json:"E"` // Event time
	Symbol       string `json:"s"` // Symbol
	AggTradeID   int64  `json:"a"` // Aggregate trade ID
	Price        string `json:"p"` // Price
	Quantity     string `json:"q"` // Quantity
	FirstTradeID int64  `json:"f"` // First trade ID
	LastTradeID  int64  `json:"l"` // Last trade ID
	TradeTime    int64  `json:"T"` // Trade time
	IsBuyerMaker bool   `json:"m"` // Is the buyer the market maker?
}

// WSAggTrade represents aggregate trade data (alias for event for consistency)
type WSAggTrade = WSAggTradeEvent

// AggTradeSubscriptionOptions defines the callback functions for aggregate trade subscription
type AggTradeSubscriptionOptions struct {
	onConnect    func()                    // Called when connection is established
	onReconnect  func()                    // Called when connection is reestablished
	onError      func(err error)           // Called when an error occurs
	onAggTrade   func(aggTrade WSAggTrade) // Called when aggregate trade data is received
	onDisconnect func()                    // Called when connection is disconnected
}

// WithConnect sets the OnConnect callback using chain method
func (a *AggTradeSubscriptionOptions) WithConnect(onConnect func()) *AggTradeSubscriptionOptions {
	a.onConnect = onConnect
	return a
}

// WithReconnect sets the OnReconnect callback using chain method
func (a *AggTradeSubscriptionOptions) WithReconnect(onReconnect func()) *AggTradeSubscriptionOptions {
	a.onReconnect = onReconnect
	return a
}

// WithError sets the OnError callback using chain method
func (a *AggTradeSubscriptionOptions) WithError(onError func(error)) *AggTradeSubscriptionOptions {
	a.onError = onError
	return a
}

// WithAggTrade sets the OnAggTrade callback using chain method
func (a *AggTradeSubscriptionOptions) WithAggTrade(onAggTrade func(WSAggTrade)) *AggTradeSubscriptionOptions {
	a.onAggTrade = onAggTrade
	return a
}

// WithDisconnect sets the OnDisconnect callback using chain method
func (a *AggTradeSubscriptionOptions) WithDisconnect(onDisconnect func()) *AggTradeSubscriptionOptions {
	a.onDisconnect = onDisconnect
	return a
}

// WSTickerEvent represents the complete 24hr ticker statistics WebSocket event
type WSTickerEvent struct {
	EventType          string `json:"e"` // Event type
	EventTime          int64  `json:"E"` // Event time
	Symbol             string `json:"s"` // Symbol
	PriceChange        string `json:"p"` // Price change
	PriceChangePercent string `json:"P"` // Price change percent
	WeightedAvgPrice   string `json:"w"` // Weighted average price
	LastPrice          string `json:"c"` // Last price
	LastQuantity       string `json:"Q"` // Last quantity
	OpenPrice          string `json:"o"` // Open price
	HighPrice          string `json:"h"` // High price
	LowPrice           string `json:"l"` // Low price
	Volume             string `json:"v"` // Total traded base asset volume
	QuoteVolume        string `json:"q"` // Total traded quote asset volume
	OpenTime           int64  `json:"O"` // Statistics open time
	CloseTime          int64  `json:"C"` // Statistics close time
	FirstTradeId       int64  `json:"F"` // First trade ID
	LastTradeId        int64  `json:"L"` // Last trade ID
	Count              int64  `json:"n"` // Total number of trades
}

// WSTicker represents 24hr ticker data (alias for event for consistency)
type WSTicker = WSTickerEvent

// TickerSubscriptionOptions defines the callback functions for ticker subscription
type TickerSubscriptionOptions struct {
	onConnect    func()                // Called when connection is established
	onReconnect  func()                // Called when connection is reestablished
	onError      func(err error)       // Called when an error occurs
	onTicker     func(ticker WSTicker) // Called when ticker data is received
	onDisconnect func()                // Called when connection is disconnected
}

// WithConnect sets the OnConnect callback using chain method
func (t *TickerSubscriptionOptions) WithConnect(onConnect func()) *TickerSubscriptionOptions {
	t.onConnect = onConnect
	return t
}

// WithReconnect sets the OnReconnect callback using chain method
func (t *TickerSubscriptionOptions) WithReconnect(onReconnect func()) *TickerSubscriptionOptions {
	t.onReconnect = onReconnect
	return t
}

// WithError sets the OnError callback using chain method
func (t *TickerSubscriptionOptions) WithError(onError func(error)) *TickerSubscriptionOptions {
	t.onError = onError
	return t
}

// WithTicker sets the OnTicker callback using chain method
func (t *TickerSubscriptionOptions) WithTicker(onTicker func(WSTicker)) *TickerSubscriptionOptions {
	t.onTicker = onTicker
	return t
}

// WithDisconnect sets the OnDisconnect callback using chain method
func (t *TickerSubscriptionOptions) WithDisconnect(onDisconnect func()) *TickerSubscriptionOptions {
	t.onDisconnect = onDisconnect
	return t
}

// WSLiquidationEvent represents the complete liquidation order WebSocket event
type WSLiquidationEvent struct {
	EventType string             `json:"e"` // Event type
	EventTime int64              `json:"E"` // Event time
	Order     WSLiquidationOrder `json:"o"` // Liquidation order data
}

// WSLiquidationOrder represents the liquidation order data within the WebSocket event
type WSLiquidationOrder struct {
	Symbol               string `json:"s"`  // Symbol
	Side                 string `json:"S"`  // Side
	OrderType            string `json:"o"`  // Order Type
	TimeInForce          string `json:"f"`  // Time in Force
	OriginalQuantity     string `json:"q"`  // Original Quantity
	Price                string `json:"p"`  // Price
	AveragePrice         string `json:"ap"` // Average Price
	OrderStatus          string `json:"X"`  // Order Status
	LastFilledQuantity   string `json:"l"`  // Order Last Filled Quantity
	FilledAccumulatedQty string `json:"z"`  // Order Filled Accumulated Quantity
	TradeTime            int64  `json:"T"`  // Order Trade Time
}

// WSLiquidation represents liquidation data (alias for event for consistency)
type WSLiquidation = WSLiquidationEvent

// LiquidationSubscriptionOptions defines the callback functions for liquidation subscription
type LiquidationSubscriptionOptions struct {
	onConnect     func()                          // Called when connection is established
	onReconnect   func()                          // Called when connection is reestablished
	onError       func(err error)                 // Called when an error occurs
	onLiquidation func(liquidation WSLiquidation) // Called when liquidation data is received
	onDisconnect  func()                          // Called when connection is disconnected
}

// WithConnect sets the OnConnect callback using chain method
func (l *LiquidationSubscriptionOptions) WithConnect(onConnect func()) *LiquidationSubscriptionOptions {
	l.onConnect = onConnect
	return l
}

// WithReconnect sets the OnReconnect callback using chain method
func (l *LiquidationSubscriptionOptions) WithReconnect(onReconnect func()) *LiquidationSubscriptionOptions {
	l.onReconnect = onReconnect
	return l
}

// WithError sets the OnError callback using chain method
func (l *LiquidationSubscriptionOptions) WithError(onError func(error)) *LiquidationSubscriptionOptions {
	l.onError = onError
	return l
}

// WithLiquidation sets the OnLiquidation callback using chain method
func (l *LiquidationSubscriptionOptions) WithLiquidation(onLiquidation func(WSLiquidation)) *LiquidationSubscriptionOptions {
	l.onLiquidation = onLiquidation
	return l
}

// WithDisconnect sets the OnDisconnect callback using chain method
func (l *LiquidationSubscriptionOptions) WithDisconnect(onDisconnect func()) *LiquidationSubscriptionOptions {
	l.onDisconnect = onDisconnect
	return l
}

// WSDepthEvent represents the partial book depth WebSocket event
type WSDepthEvent struct {
	EventType       string     `json:"e"`  // Event type
	EventTime       int64      `json:"E"`  // Event time
	TransactionTime int64      `json:"T"`  // Transaction time
	Symbol          string     `json:"s"`  // Symbol
	FirstUpdateID   int64      `json:"U"`  // First update ID in event
	FinalUpdateID   int64      `json:"u"`  // Final update ID in event
	PrevUpdateID    int64      `json:"pu"` // Final update Id in last stream
	Bids            [][]string `json:"b"`  // Bids to be updated [["price", "quantity"], ...]
	Asks            [][]string `json:"a"`  // Asks to be updated [["price", "quantity"], ...]
}

// WSDepth represents depth data (alias for event for consistency)
type WSDepth = WSDepthEvent

// DepthUpdateSpeed represents the update speed for depth streams
type DepthUpdateSpeed string

const (
	DepthUpdate100ms DepthUpdateSpeed = "100ms"
	DepthUpdate250ms DepthUpdateSpeed = "250ms" // Default
	DepthUpdate500ms DepthUpdateSpeed = "500ms"
)

// DepthLevel represents valid depth levels
type DepthLevel int

const (
	DepthLevel5  DepthLevel = 5
	DepthLevel10 DepthLevel = 10
	DepthLevel20 DepthLevel = 20
)

// DepthSubscriptionOptions defines the callback functions for depth subscription
type DepthSubscriptionOptions struct {
	onConnect    func()              // Called when connection is established
	onReconnect  func()              // Called when connection is reestablished
	onError      func(err error)     // Called when an error occurs
	onDepth      func(depth WSDepth) // Called when depth data is received
	onDisconnect func()              // Called when connection is disconnected
}

// WithConnect sets the OnConnect callback using chain method
func (d *DepthSubscriptionOptions) WithConnect(onConnect func()) *DepthSubscriptionOptions {
	d.onConnect = onConnect
	return d
}

// WithReconnect sets the OnReconnect callback using chain method
func (d *DepthSubscriptionOptions) WithReconnect(onReconnect func()) *DepthSubscriptionOptions {
	d.onReconnect = onReconnect
	return d
}

// WithError sets the OnError callback using chain method
func (d *DepthSubscriptionOptions) WithError(onError func(error)) *DepthSubscriptionOptions {
	d.onError = onError
	return d
}

// WithDepth sets the OnDepth callback using chain method
func (d *DepthSubscriptionOptions) WithDepth(onDepth func(WSDepth)) *DepthSubscriptionOptions {
	d.onDepth = onDepth
	return d
}

// WithDisconnect sets the OnDisconnect callback using chain method
func (d *DepthSubscriptionOptions) WithDisconnect(onDisconnect func()) *DepthSubscriptionOptions {
	d.onDisconnect = onDisconnect
	return d
}

// DiffDepthSubscriptionOptions defines the callback functions for differential depth subscription
// Note: Reuses WSDepthEvent structure but represents order book changes rather than snapshots
type DiffDepthSubscriptionOptions struct {
	onConnect    func()                  // Called when connection is established
	onReconnect  func()                  // Called when connection is reestablished
	onError      func(err error)         // Called when an error occurs
	onDiffDepth  func(diffDepth WSDepth) // Called when differential depth data is received
	onDisconnect func()                  // Called when connection is disconnected
}

// User Data Stream Events

// WSListenKeyExpiredEvent represents a listen key expiration event (handled internally)
type WSListenKeyExpiredEvent struct {
	EventType string `json:"e"`         // Event type ("listenKeyExpired")
	EventTime int64  `json:"E"`         // Event time
	ListenKey string `json:"listenKey"` // Expired listen key
}

// WSAccountUpdateEvent represents an account update event
type WSAccountUpdateEvent struct {
	EventType       string              `json:"e"` // Event type ("ACCOUNT_UPDATE")
	EventTime       int64               `json:"E"` // Event time
	TransactionTime int64               `json:"T"` // Transaction time
	UpdateData      WSAccountUpdateData `json:"a"` // Update data
}

// WSAccountUpdateData represents the update data in account update event
type WSAccountUpdateData struct {
	EventReasonType string              `json:"m"` // Event reason type
	Balances        []WSAccountBalance  `json:"B"` // Balances
	Positions       []WSAccountPosition `json:"P"` // Positions
}

// WSAccountBalance represents balance information in account update
type WSAccountBalance struct {
	Asset              string `json:"a"`  // Asset
	WalletBalance      string `json:"wb"` // Wallet Balance
	CrossWalletBalance string `json:"cw"` // Cross Wallet Balance
	BalanceChange      string `json:"bc"` // Balance Change except PnL and Commission
}

// WSAccountPosition represents position information in account update
type WSAccountPosition struct {
	Symbol              string `json:"s"`   // Symbol
	PositionAmount      string `json:"pa"`  // Position Amount
	EntryPrice          string `json:"ep"`  // Entry Price
	BreakEvenPrice      string `json:"bep"` // Breakeven Price
	AccumulatedRealized string `json:"cr"`  // (Pre-fee) Accumulated Realized
	UnrealizedPnL       string `json:"up"`  // Unrealized PnL
	MarginType          string `json:"mt"`  // Margin Type
	IsolatedWallet      string `json:"iw"`  // Isolated Wallet (if isolated position)
	PositionSide        string `json:"ps"`  // Position Side
}

// WSMarginCallEvent represents a margin call event
type WSMarginCallEvent struct {
	EventType          string                 `json:"e"`  // Event type ("MARGIN_CALL")
	EventTime          int64                  `json:"E"`  // Event time
	CrossWalletBalance string                 `json:"cw"` // Cross Wallet Balance
	Positions          []WSMarginCallPosition `json:"p"`  // Position(s) of Margin Call
}

// WSMarginCallPosition represents position information in margin call
type WSMarginCallPosition struct {
	Symbol                    string `json:"s"`  // Symbol
	PositionSide              string `json:"ps"` // Position Side
	PositionAmount            string `json:"pa"` // Position Amount
	MarginType                string `json:"mt"` // Margin Type
	IsolatedWallet            string `json:"iw"` // Isolated Wallet (if isolated position)
	MarkPrice                 string `json:"mp"` // Mark Price
	UnrealizedPnL             string `json:"up"` // Unrealized PnL
	MaintenanceMarginRequired string `json:"mm"` // Maintenance Margin Required
}

// WSOrderTradeUpdateEvent represents an order trade update event
type WSOrderTradeUpdateEvent struct {
	EventType       string             `json:"e"` // Event type ("ORDER_TRADE_UPDATE")
	EventTime       int64              `json:"E"` // Event time
	TransactionTime int64              `json:"T"` // Transaction time
	Order           WSOrderTradeUpdate `json:"o"` // Order information
}

// WSOrderTradeUpdate represents order information in order trade update
type WSOrderTradeUpdate struct {
	Symbol                    string `json:"s"`   // Symbol
	ClientOrderID             string `json:"c"`   // Client Order Id
	Side                      string `json:"S"`   // Side
	OrderType                 string `json:"o"`   // Order Type
	TimeInForce               string `json:"f"`   // Time in Force
	OriginalQuantity          string `json:"q"`   // Original Quantity
	OriginalPrice             string `json:"p"`   // Original Price
	AveragePrice              string `json:"ap"`  // Average Price
	StopPrice                 string `json:"sp"`  // Stop Price
	ExecutionType             string `json:"x"`   // Execution Type
	OrderStatus               string `json:"X"`   // Order Status
	OrderID                   int64  `json:"i"`   // Order Id
	LastFilledQuantity        string `json:"l"`   // Order Last Filled Quantity
	FilledAccumulatedQuantity string `json:"z"`   // Order Filled Accumulated Quantity
	LastFilledPrice           string `json:"L"`   // Last Filled Price
	CommissionAsset           string `json:"N"`   // Commission Asset
	Commission                string `json:"n"`   // Commission
	OrderTradeTime            int64  `json:"T"`   // Order Trade Time
	TradeID                   int64  `json:"t"`   // Trade Id
	BidsNotional              string `json:"b"`   // Bids Notional
	AskNotional               string `json:"a"`   // Ask Notional
	IsMakerSide               bool   `json:"m"`   // Is this trade the maker side?
	IsReduceOnly              bool   `json:"R"`   // Is this reduce only
	StopPriceWorkingType      string `json:"wt"`  // Stop Price Working Type
	OriginalOrderType         string `json:"ot"`  // Original Order Type
	PositionSide              string `json:"ps"`  // Position Side
	IsCloseAll                bool   `json:"cp"`  // If Close-All, pushed with conditional order
	ActivationPrice           string `json:"AP"`  // Activation Price, only pushed with TRAILING_STOP_MARKET order
	CallbackRate              string `json:"cr"`  // Callback Rate, only pushed with TRAILING_STOP_MARKET order
	IsPriceProtection         bool   `json:"pP"`  // If price protection is turned on
	RealizedProfit            string `json:"rp"`  // Realized Profit of the trade
	STPMode                   string `json:"V"`   // STP mode
	PriceMatchMode            string `json:"pm"`  // Price match mode
	GTDOrderAutoCancelTime    int64  `json:"gtd"` // TIF GTD order auto cancel time
}

// WSTradeLiteEvent represents a trade lite update event
type WSTradeLiteEvent struct {
	EventType          string `json:"e"` // Event type ("TRADE_LITE")
	EventTime          int64  `json:"E"` // Event time
	TransactionTime    int64  `json:"T"` // Transaction time
	Symbol             string `json:"s"` // Symbol
	OriginalQuantity   string `json:"q"` // Original Quantity
	OriginalPrice      string `json:"p"` // Original Price
	IsMakerSide        bool   `json:"m"` // Is this trade the maker side?
	ClientOrderID      string `json:"c"` // Client Order Id
	Side               string `json:"S"` // Side
	LastFilledPrice    string `json:"L"` // Last Filled Price
	LastFilledQuantity string `json:"l"` // Order Last Filled Quantity
	TradeID            int64  `json:"t"` // Trade Id
	OrderID            int64  `json:"i"` // Order Id
}

// UserDataSubscriptionOptions defines callbacks for user data stream events
type UserDataSubscriptionOptions struct {
	onConnect       func()                                    // Called when connection is established
	onReconnect     func()                                    // Called when connection is reestablished (includes unexpected disconnects and listen key refreshes)
	onError         func(err error)                           // Called when an error occurs
	onAccountUpdate func(accountUpdate WSAccountUpdateEvent)  // Called when account update is received
	onMarginCall    func(marginCall WSMarginCallEvent)        // Called when margin call is received
	onOrderUpdate   func(orderUpdate WSOrderTradeUpdateEvent) // Called when order trade update is received
	onTradeLite     func(tradeLite WSTradeLiteEvent)          // Called when trade lite update is received
	onDisconnect    func()                                    // Called when connection is disconnected
}

// WithConnect sets the OnConnect callback using chain method
func (dd *DiffDepthSubscriptionOptions) WithConnect(onConnect func()) *DiffDepthSubscriptionOptions {
	dd.onConnect = onConnect
	return dd
}

// WithReconnect sets the OnReconnect callback using chain method
func (dd *DiffDepthSubscriptionOptions) WithReconnect(onReconnect func()) *DiffDepthSubscriptionOptions {
	dd.onReconnect = onReconnect
	return dd
}

// WithError sets the OnError callback using chain method
func (dd *DiffDepthSubscriptionOptions) WithError(onError func(error)) *DiffDepthSubscriptionOptions {
	dd.onError = onError
	return dd
}

// WithDiffDepth sets the OnDiffDepth callback using chain method
func (dd *DiffDepthSubscriptionOptions) WithDiffDepth(onDiffDepth func(WSDepth)) *DiffDepthSubscriptionOptions {
	dd.onDiffDepth = onDiffDepth
	return dd
}

// WithDisconnect sets the OnDisconnect callback using chain method
func (dd *DiffDepthSubscriptionOptions) WithDisconnect(onDisconnect func()) *DiffDepthSubscriptionOptions {
	dd.onDisconnect = onDisconnect
	return dd
}

// WSSubscription represents an active WebSocket subscription
type WSSubscription struct {
	id      string
	conn    *BinancePerpWSConn
	options interface{} // Can be KlineSubscriptionOptions, AggTradeSubscriptionOptions, TickerSubscriptionOptions, LiquidationSubscriptionOptions, DepthSubscriptionOptions, DiffDepthSubscriptionOptions, or other subscription types
	state   ConnectionState
}

// WithConnect sets the OnConnect callback for user data subscription
func (o *UserDataSubscriptionOptions) WithConnect(onConnect func()) *UserDataSubscriptionOptions {
	o.onConnect = onConnect
	return o
}

// WithReconnect sets the OnReconnect callback for user data subscription
func (o *UserDataSubscriptionOptions) WithReconnect(onReconnect func()) *UserDataSubscriptionOptions {
	o.onReconnect = onReconnect
	return o
}

// WithError sets the OnError callback for user data subscription
func (o *UserDataSubscriptionOptions) WithError(onError func(error)) *UserDataSubscriptionOptions {
	o.onError = onError
	return o
}

// WithAccountUpdate sets the OnAccountUpdate callback for user data subscription
func (o *UserDataSubscriptionOptions) WithAccountUpdate(onAccountUpdate func(WSAccountUpdateEvent)) *UserDataSubscriptionOptions {
	o.onAccountUpdate = onAccountUpdate
	return o
}

// WithMarginCall sets the OnMarginCall callback for user data subscription
func (o *UserDataSubscriptionOptions) WithMarginCall(onMarginCall func(WSMarginCallEvent)) *UserDataSubscriptionOptions {
	o.onMarginCall = onMarginCall
	return o
}

// WithOrderUpdate sets the OnOrderUpdate callback for user data subscription
func (o *UserDataSubscriptionOptions) WithOrderUpdate(onOrderUpdate func(WSOrderTradeUpdateEvent)) *UserDataSubscriptionOptions {
	o.onOrderUpdate = onOrderUpdate
	return o
}

// WithTradeLite sets the OnTradeLite callback for user data subscription
func (o *UserDataSubscriptionOptions) WithTradeLite(onTradeLite func(WSTradeLiteEvent)) *UserDataSubscriptionOptions {
	o.onTradeLite = onTradeLite
	return o
}

// WithDisconnect sets the OnDisconnect callback for user data subscription
func (o *UserDataSubscriptionOptions) WithDisconnect(onDisconnect func()) *UserDataSubscriptionOptions {
	o.onDisconnect = onDisconnect
	return o
}
