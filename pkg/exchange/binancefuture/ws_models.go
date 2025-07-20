package binancefuture

import (
	"encoding/json"
)

// WebSocketMessage represents a generic WebSocket message
type WebSocketMessage struct {
	Type   string          `json:"type,omitempty"`
	Method string          `json:"method,omitempty"`
	Params []string        `json:"params,omitempty"`
	ID     int64           `json:"id,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *WSError        `json:"error,omitempty"`
}

// WSError represents a WebSocket error response
type WSError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// WSKlineData represents kline/candlestick data from WebSocket
type WSKlineData struct {
	Symbol            string  `json:"s"`
	Kline             WSKline `json:"k"`
	EventTime         int64   `json:"E"`
	EventType         string  `json:"e"`
	FirstTradeID      int64   `json:"f"`
	LastTradeID       int64   `json:"L"`
	IsKlineClosed     bool    `json:"x"`
	QuoteVolume       float64 `json:"q,string"`
	ActiveBuyVolume   float64 `json:"V,string"`
	ActiveBuyQuoteVol float64 `json:"Q,string"`
}

// WSKline represents individual kline data from WebSocket
type WSKline struct {
	StartTime         int64   `json:"t"`
	CloseTime         int64   `json:"T"`
	Symbol            string  `json:"s"`
	Interval          string  `json:"i"`
	FirstTradeID      int64   `json:"f"`
	LastTradeID       int64   `json:"L"`
	OpenPrice         float64 `json:"o,string"`
	ClosePrice        float64 `json:"c,string"`
	HighPrice         float64 `json:"h,string"`
	LowPrice          float64 `json:"l,string"`
	Volume            float64 `json:"v,string"`
	TradeCount        int64   `json:"n"`
	IsKlineClosed     bool    `json:"x"`
	QuoteVolume       float64 `json:"q,string"`
	ActiveBuyVolume   float64 `json:"V,string"`
	ActiveBuyQuoteVol float64 `json:"Q,string"`
	Ignore            float64 `json:"B,string"`
}

// WSTickerData represents 24hr ticker data from WebSocket
type WSTickerData struct {
	Symbol             string  `json:"s"`
	PriceChange        float64 `json:"P,string"`
	PriceChangePercent float64 `json:"p,string"`
	WeightedAvgPrice   float64 `json:"w,string"`
	PrevClosePrice     float64 `json:"x,string"`
	LastPrice          float64 `json:"c,string"`
	LastQty            float64 `json:"Q,string"`
	BidPrice           float64 `json:"b,string"`
	BidQty             float64 `json:"B,string"`
	AskPrice           float64 `json:"a,string"`
	AskQty             float64 `json:"A,string"`
	OpenPrice          float64 `json:"o,string"`
	HighPrice          float64 `json:"h,string"`
	LowPrice           float64 `json:"l,string"`
	Volume             float64 `json:"v,string"`
	QuoteVolume        float64 `json:"q,string"`
	OpenTime           int64   `json:"O"`
	CloseTime          int64   `json:"C"`
	FirstID            int64   `json:"F"`
	LastID             int64   `json:"L"`
	Count              int64   `json:"n"`
	EventTime          int64   `json:"E"`
	EventType          string  `json:"e"`
}

// WSMiniTickerData represents mini ticker data from WebSocket
type WSMiniTickerData struct {
	Symbol      string  `json:"s"`
	ClosePrice  float64 `json:"c,string"`
	OpenPrice   float64 `json:"o,string"`
	HighPrice   float64 `json:"h,string"`
	LowPrice    float64 `json:"l,string"`
	Volume      float64 `json:"v,string"`
	QuoteVolume float64 `json:"q,string"`
	EventTime   int64   `json:"E"`
	EventType   string  `json:"e"`
}

// WSBookTickerData represents book ticker data from WebSocket
type WSBookTickerData struct {
	Symbol          string  `json:"s"`
	BidPrice        float64 `json:"b,string"`
	BidQty          float64 `json:"B,string"`
	AskPrice        float64 `json:"a,string"`
	AskQty          float64 `json:"A,string"`
	EventTime       int64   `json:"E"`
	EventType       string  `json:"e"`
	UpdateID        int64   `json:"u"`
	TransactionTime int64   `json:"T"`
}

// WSDepthData represents order book depth data from WebSocket
type WSDepthData struct {
	Symbol        string     `json:"s"`
	EventTime     int64      `json:"E"`
	EventType     string     `json:"e"`
	FirstUpdateID int64      `json:"U"`
	FinalUpdateID int64      `json:"u"`
	Bids          [][]string `json:"b"`
	Asks          [][]string `json:"a"`
}

// WSPartialDepthData represents partial book depth data from WebSocket
type WSPartialDepthData struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// WSDiffDepthData represents diff depth data from WebSocket
type WSDiffDepthData struct {
	EventType     string     `json:"e"`
	EventTime     int64      `json:"E"`
	Symbol        string     `json:"s"`
	FirstUpdateID int64      `json:"U"`
	FinalUpdateID int64      `json:"u"`
	Bids          [][]string `json:"b"`
	Asks          [][]string `json:"a"`
}

// WSTradeData represents trade data from WebSocket
type WSTradeData struct {
	Symbol        string  `json:"s"`
	ID            int64   `json:"t"`
	Price         float64 `json:"p,string"`
	Quantity      float64 `json:"q,string"`
	BuyerOrderID  int64   `json:"b"`
	SellerOrderID int64   `json:"a"`
	TradeTime     int64   `json:"T"`
	IsBuyerMaker  bool    `json:"m"`
	Ignore        bool    `json:"M"`
	EventTime     int64   `json:"E"`
	EventType     string  `json:"e"`
}

// WSAggTradeData represents aggregated trade data from WebSocket
type WSAggTradeData struct {
	Symbol       string  `json:"s"`
	ID           int64   `json:"a"`
	Price        float64 `json:"p,string"`
	Quantity     float64 `json:"q,string"`
	FirstTradeID int64   `json:"f"`
	LastTradeID  int64   `json:"l"`
	TradeTime    int64   `json:"T"`
	IsBuyerMaker bool    `json:"m"`
	Ignore       bool    `json:"M"`
	EventTime    int64   `json:"E"`
	EventType    string  `json:"e"`
}

// WSMarkPriceData represents mark price data from WebSocket
type WSMarkPriceData struct {
	Symbol          string  `json:"s"`
	MarkPrice       float64 `json:"p,string"`
	IndexPrice      float64 `json:"i,string"`
	EstimatedPrice  float64 `json:"P,string"`
	LastFundingRate float64 `json:"r,string"`
	NextFundingTime int64   `json:"T"`
	EventTime       int64   `json:"E"`
	EventType       string  `json:"e"`
}

// WSFundingRateData represents funding rate data from WebSocket
type WSFundingRateData struct {
	Symbol      string  `json:"s"`
	FundingRate float64 `json:"r,string"`
	FundingTime int64   `json:"T"`
	EventTime   int64   `json:"E"`
	EventType   string  `json:"e"`
}

// WSBalance represents a balance in user data stream
type WSBalance struct {
	Asset  string `json:"a"`
	Free   string `json:"f"`
	Locked string `json:"l"`
}

// WSOutboundAccountPosition represents outbound account position data
type WSOutboundAccountPosition struct {
	EventType    string      `json:"e"`
	EventTime    int64       `json:"E"`
	LastUpdateID int64       `json:"u"`
	Balances     []WSBalance `json:"B"`
}

// WSBalanceUpdate represents balance update data
type WSBalanceUpdate struct {
	EventType    string `json:"e"`
	EventTime    int64  `json:"E"`
	Asset        string `json:"a"`
	BalanceDelta string `json:"d"`
	ClearTime    int64  `json:"T"`
}

// WSExecutionReport represents execution report data
type WSExecutionReport struct {
	EventType                    string `json:"e"`
	EventTime                    int64  `json:"E"`
	Symbol                       string `json:"s"`
	ClientOrderID                string `json:"c"`
	Side                         string `json:"S"`
	OrderType                    string `json:"o"`
	TimeInForce                  string `json:"f"`
	OrderQuantity                string `json:"q"`
	OrderPrice                   string `json:"p"`
	StopPrice                    string `json:"P"`
	IcebergQuantity              string `json:"F"`
	OrderListID                  int64  `json:"g"`
	OriginalClientOrderID        string `json:"C"`
	CurrentExecutionType         string `json:"x"`
	CurrentOrderStatus           string `json:"X"`
	OrderRejectReason            string `json:"r"`
	OrderID                      int64  `json:"i"`
	LastExecutedQuantity         string `json:"l"`
	CumulativeFilledQuantity     string `json:"z"`
	LastExecutedPrice            string `json:"L"`
	CommissionAmount             string `json:"n"`
	CommissionAsset              string `json:"N"`
	TransactionTime              int64  `json:"T"`
	TradeID                      int64  `json:"t"`
	PreventedMatchID             int64  `json:"v"`
	ExecutionID                  int64  `json:"I"`
	IsOrderOnBook                bool   `json:"w"`
	IsTradeMakerSide             bool   `json:"m"`
	Ignore                       bool   `json:"M"`
	OrderCreationTime            int64  `json:"O"`
	CumulativeQuoteAssetQuantity string `json:"Z"`
	LastQuoteAssetQuantity       string `json:"Y"`
	QuoteOrderQuantity           string `json:"Q"`
	WorkingTime                  int64  `json:"W"`
	SelfTradePreventionMode      string `json:"V"`
}

// Callback types for different WebSocket data types
type KlineCallback func(data *WSKlineData) error
type TickerCallback func(data *WSTickerData) error
type MiniTickerCallback func(data *WSMiniTickerData) error
type BookTickerCallback func(data *WSBookTickerData) error
type PartialDepthCallback func(data *WSPartialDepthData) error
type DiffDepthCallback func(data *WSDiffDepthData) error
type DepthCallback func(data *WSDepthData) error
type TradeCallback func(data *WSTradeData) error
type AggTradeCallback func(data *WSAggTradeData) error
type MarkPriceCallback func(data *WSMarkPriceData) error
type FundingRateCallback func(data *WSFundingRateData) error
type OutboundAccountPositionCallback func(data *WSOutboundAccountPosition) error
type BalanceUpdateCallback func(data *WSBalanceUpdate) error
type ExecutionReportCallback func(data *WSExecutionReport) error

// WebSocketCallback represents a generic WebSocket callback function
type WebSocketCallback func(data []byte) error

// WebSocketStream represents a WebSocket stream with its callback
type WebSocketStream struct {
	StreamName string
	Callback   WebSocketCallback
}

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection struct {
	URL      string
	Streams  []WebSocketStream
	IsActive bool
	Close    chan struct{}
}

// WSUserDataStreamEvent represents a generic user data stream event
type WSUserDataStreamEvent struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
}

// WSListenKeyExpiredEvent represents listen key expired event
type WSListenKeyExpiredEvent struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	ListenKey string `json:"listenKey"`
}

// WSAccountUpdateEvent represents account update event
type WSAccountUpdateEvent struct {
	EventType       string              `json:"e"`
	EventTime       int64               `json:"E"`
	TransactionTime int64               `json:"T"`
	UpdateData      WSAccountUpdateData `json:"a"`
}

// WSAccountUpdateData represents account update data
type WSAccountUpdateData struct {
	EventReasonType string           `json:"m"`
	Balances        []WSBalanceData  `json:"B"`
	Positions       []WSPositionData `json:"P"`
}

// WSBalanceData represents balance data in account update
type WSBalanceData struct {
	Asset              string `json:"a"`
	WalletBalance      string `json:"wb"`
	CrossWalletBalance string `json:"cw"`
	BalanceChange      string `json:"bc"`
}

// WSPositionData represents position data in account update
type WSPositionData struct {
	Symbol              string `json:"s"`
	PositionAmount      string `json:"pa"`
	EntryPrice          string `json:"ep"`
	BreakevenPrice      string `json:"bep"`
	AccumulatedRealized string `json:"cr"`
	UnrealizedPnL       string `json:"up"`
	MarginType          string `json:"mt"`
	IsolatedWallet      string `json:"iw"`
	PositionSide        string `json:"ps"`
}

// WSMarginCallEvent represents margin call event
type WSMarginCallEvent struct {
	EventType          string                 `json:"e"`
	EventTime          int64                  `json:"E"`
	CrossWalletBalance string                 `json:"cw"`
	Positions          []WSMarginCallPosition `json:"p"`
}

// WSMarginCallPosition represents position in margin call
type WSMarginCallPosition struct {
	Symbol                    string `json:"s"`
	PositionSide              string `json:"ps"`
	PositionAmount            string `json:"pa"`
	MarginType                string `json:"mt"`
	IsolatedWallet            string `json:"iw"`
	MarkPrice                 string `json:"mp"`
	UnrealizedPnL             string `json:"up"`
	MaintenanceMarginRequired string `json:"mm"`
}

// WSOrderTradeUpdateEvent represents order trade update event
type WSOrderTradeUpdateEvent struct {
	EventType       string                  `json:"e"`
	EventTime       int64                   `json:"E"`
	TransactionTime int64                   `json:"T"`
	Order           WSOrderTradeUpdateOrder `json:"o"`
}

// WSOrderTradeUpdateOrder represents order in trade update
type WSOrderTradeUpdateOrder struct {
	Symbol                    string `json:"s"`
	ClientOrderID             string `json:"c"`
	Side                      string `json:"S"`
	OrderType                 string `json:"o"`
	TimeInForce               string `json:"f"`
	OriginalQuantity          string `json:"q"`
	OriginalPrice             string `json:"p"`
	AveragePrice              string `json:"ap"`
	StopPrice                 string `json:"sp"`
	ExecutionType             string `json:"x"`
	OrderStatus               string `json:"X"`
	OrderID                   int64  `json:"i"`
	LastFilledQuantity        string `json:"l"`
	FilledAccumulatedQuantity string `json:"z"`
	LastFilledPrice           string `json:"L"`
	CommissionAsset           string `json:"N"`
	Commission                string `json:"n"`
	OrderTradeTime            int64  `json:"T"`
	TradeID                   int64  `json:"t"`
	BidsNotional              string `json:"b"`
	AsksNotional              string `json:"a"`
	IsMaker                   bool   `json:"m"`
	IsReduceOnly              bool   `json:"R"`
	WorkingType               string `json:"wt"`
	OriginalOrderType         string `json:"ot"`
	PositionSide              string `json:"ps"`
	IsCloseAll                bool   `json:"cp"`
	ActivationPrice           string `json:"AP"`
	CallbackRate              string `json:"cr"`
	PriceProtection           bool   `json:"pP"`
	RealizedProfit            string `json:"rp"`
	STPMode                   string `json:"V"`
	PriceMatchMode            string `json:"pm"`
	GoodTillDate              int64  `json:"gtd"`
	// Additional fields from the payload
	Ignore1 int `json:"si,omitempty"` // ignore field
	Ignore2 int `json:"ss,omitempty"` // ignore field
}

// WSTradeLiteEvent represents trade lite event
type WSTradeLiteEvent struct {
	EventType          string `json:"e"`
	EventTime          int64  `json:"E"`
	TransactionTime    int64  `json:"T"`
	Symbol             string `json:"s"`
	Quantity           string `json:"q"`
	Price              string `json:"p"`
	IsMaker            bool   `json:"m"`
	ClientOrderID      string `json:"c"`
	Side               string `json:"S"`
	LastFilledPrice    string `json:"L"`
	LastFilledQuantity string `json:"l"`
	TradeID            int64  `json:"t"`
	OrderID            int64  `json:"i"`
}

// WSAccountConfigUpdateEvent represents account config update event
type WSAccountConfigUpdateEvent struct {
	EventType       string          `json:"e"`
	EventTime       int64           `json:"E"`
	TransactionTime int64           `json:"T"`
	AccountConfig   WSAccountConfig `json:"ac"`
}

// WSAccountConfig represents account configuration
type WSAccountConfig struct {
	Symbol   string `json:"s"`
	Leverage int    `json:"l"`
}

// User data stream callback types
type ListenKeyExpiredCallback func(data *WSListenKeyExpiredEvent) error
type AccountUpdateCallback func(data *WSAccountUpdateEvent) error
type MarginCallCallback func(data *WSMarginCallEvent) error
type OrderTradeUpdateCallback func(data *WSOrderTradeUpdateEvent) error
type TradeLiteCallback func(data *WSTradeLiteEvent) error
type AccountConfigUpdateCallback func(data *WSAccountConfigUpdateEvent) error

// SubscriptionOptions represents base subscription options with common callbacks
type SubscriptionOptions struct {
	connectCallback    func()
	reconnectCallback  func()
	disconnectCallback func()
	errorCallback      func(error)
}

// WithConnect sets the connect callback
func (o *SubscriptionOptions) WithConnect(callback func()) *SubscriptionOptions {
	o.connectCallback = callback
	return o
}

// WithReconnect sets the reconnect callback
func (o *SubscriptionOptions) WithReconnect(callback func()) *SubscriptionOptions {
	o.reconnectCallback = callback
	return o
}

// WithDisconnect sets the disconnect callback
func (o *SubscriptionOptions) WithDisconnect(callback func()) *SubscriptionOptions {
	o.disconnectCallback = callback
	return o
}

// WithError sets the error callback
func (o *SubscriptionOptions) WithError(callback func(error)) *SubscriptionOptions {
	o.errorCallback = callback
	return o
}

// KlineSubscriptionOptions represents kline subscription options
type KlineSubscriptionOptions struct {
	*SubscriptionOptions
	klineCallback KlineCallback
}

// NewKlineSubscriptionOptions creates a new kline subscription options
func NewKlineSubscriptionOptions() *KlineSubscriptionOptions {
	return &KlineSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithKline sets the kline callback
func (o *KlineSubscriptionOptions) WithKline(callback KlineCallback) *KlineSubscriptionOptions {
	o.klineCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *KlineSubscriptionOptions) WithConnect(callback func()) *KlineSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *KlineSubscriptionOptions) WithReconnect(callback func()) *KlineSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *KlineSubscriptionOptions) WithDisconnect(callback func()) *KlineSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *KlineSubscriptionOptions) WithError(callback func(error)) *KlineSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

// TickerSubscriptionOptions represents ticker subscription options
type TickerSubscriptionOptions struct {
	*SubscriptionOptions
	tickerCallback TickerCallback
}

// NewTickerSubscriptionOptions creates a new ticker subscription options
func NewTickerSubscriptionOptions() *TickerSubscriptionOptions {
	return &TickerSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithTicker sets the ticker callback
func (o *TickerSubscriptionOptions) WithTicker(callback TickerCallback) *TickerSubscriptionOptions {
	o.tickerCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *TickerSubscriptionOptions) WithConnect(callback func()) *TickerSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *TickerSubscriptionOptions) WithReconnect(callback func()) *TickerSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *TickerSubscriptionOptions) WithDisconnect(callback func()) *TickerSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *TickerSubscriptionOptions) WithError(callback func(error)) *TickerSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

// MiniTickerSubscriptionOptions represents mini ticker subscription options
type MiniTickerSubscriptionOptions struct {
	*SubscriptionOptions
	miniTickerCallback MiniTickerCallback
}

// NewMiniTickerSubscriptionOptions creates a new mini ticker subscription options
func NewMiniTickerSubscriptionOptions() *MiniTickerSubscriptionOptions {
	return &MiniTickerSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithMiniTicker sets the mini ticker callback
func (o *MiniTickerSubscriptionOptions) WithMiniTicker(callback MiniTickerCallback) *MiniTickerSubscriptionOptions {
	o.miniTickerCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *MiniTickerSubscriptionOptions) WithConnect(callback func()) *MiniTickerSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *MiniTickerSubscriptionOptions) WithReconnect(callback func()) *MiniTickerSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *MiniTickerSubscriptionOptions) WithDisconnect(callback func()) *MiniTickerSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *MiniTickerSubscriptionOptions) WithError(callback func(error)) *MiniTickerSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

// BookTickerSubscriptionOptions represents book ticker subscription options
type BookTickerSubscriptionOptions struct {
	*SubscriptionOptions
	bookTickerCallback BookTickerCallback
}

// NewBookTickerSubscriptionOptions creates a new book ticker subscription options
func NewBookTickerSubscriptionOptions() *BookTickerSubscriptionOptions {
	return &BookTickerSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithBookTicker sets the book ticker callback
func (o *BookTickerSubscriptionOptions) WithBookTicker(callback BookTickerCallback) *BookTickerSubscriptionOptions {
	o.bookTickerCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *BookTickerSubscriptionOptions) WithConnect(callback func()) *BookTickerSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *BookTickerSubscriptionOptions) WithReconnect(callback func()) *BookTickerSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *BookTickerSubscriptionOptions) WithDisconnect(callback func()) *BookTickerSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *BookTickerSubscriptionOptions) WithError(callback func(error)) *BookTickerSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

// DepthSubscriptionOptions represents depth subscription options
type DepthSubscriptionOptions struct {
	*SubscriptionOptions
	depthCallback DepthCallback
}

// NewDepthSubscriptionOptions creates a new depth subscription options
func NewDepthSubscriptionOptions() *DepthSubscriptionOptions {
	return &DepthSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithDepth sets the depth callback
func (o *DepthSubscriptionOptions) WithDepth(callback DepthCallback) *DepthSubscriptionOptions {
	o.depthCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *DepthSubscriptionOptions) WithConnect(callback func()) *DepthSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *DepthSubscriptionOptions) WithReconnect(callback func()) *DepthSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *DepthSubscriptionOptions) WithDisconnect(callback func()) *DepthSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *DepthSubscriptionOptions) WithError(callback func(error)) *DepthSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

// TradeSubscriptionOptions represents trade subscription options
type TradeSubscriptionOptions struct {
	*SubscriptionOptions
	tradeCallback TradeCallback
}

// NewTradeSubscriptionOptions creates a new trade subscription options
func NewTradeSubscriptionOptions() *TradeSubscriptionOptions {
	return &TradeSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithTrade sets the trade callback
func (o *TradeSubscriptionOptions) WithTrade(callback TradeCallback) *TradeSubscriptionOptions {
	o.tradeCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *TradeSubscriptionOptions) WithConnect(callback func()) *TradeSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *TradeSubscriptionOptions) WithReconnect(callback func()) *TradeSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *TradeSubscriptionOptions) WithDisconnect(callback func()) *TradeSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *TradeSubscriptionOptions) WithError(callback func(error)) *TradeSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

// AggTradeSubscriptionOptions represents aggregated trade subscription options
type AggTradeSubscriptionOptions struct {
	*SubscriptionOptions
	aggTradeCallback AggTradeCallback
}

// NewAggTradeSubscriptionOptions creates a new aggregated trade subscription options
func NewAggTradeSubscriptionOptions() *AggTradeSubscriptionOptions {
	return &AggTradeSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithAggTrade sets the aggregated trade callback
func (o *AggTradeSubscriptionOptions) WithAggTrade(callback AggTradeCallback) *AggTradeSubscriptionOptions {
	o.aggTradeCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *AggTradeSubscriptionOptions) WithConnect(callback func()) *AggTradeSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *AggTradeSubscriptionOptions) WithReconnect(callback func()) *AggTradeSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *AggTradeSubscriptionOptions) WithDisconnect(callback func()) *AggTradeSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *AggTradeSubscriptionOptions) WithError(callback func(error)) *AggTradeSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

// MarkPriceSubscriptionOptions represents mark price subscription options
type MarkPriceSubscriptionOptions struct {
	*SubscriptionOptions
	markPriceCallback MarkPriceCallback
}

// NewMarkPriceSubscriptionOptions creates a new mark price subscription options
func NewMarkPriceSubscriptionOptions() *MarkPriceSubscriptionOptions {
	return &MarkPriceSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithMarkPrice sets the mark price callback
func (o *MarkPriceSubscriptionOptions) WithMarkPrice(callback MarkPriceCallback) *MarkPriceSubscriptionOptions {
	o.markPriceCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *MarkPriceSubscriptionOptions) WithConnect(callback func()) *MarkPriceSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *MarkPriceSubscriptionOptions) WithReconnect(callback func()) *MarkPriceSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *MarkPriceSubscriptionOptions) WithDisconnect(callback func()) *MarkPriceSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *MarkPriceSubscriptionOptions) WithError(callback func(error)) *MarkPriceSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

// FundingRateSubscriptionOptions represents funding rate subscription options
type FundingRateSubscriptionOptions struct {
	*SubscriptionOptions
	fundingRateCallback FundingRateCallback
}

// NewFundingRateSubscriptionOptions creates a new funding rate subscription options
func NewFundingRateSubscriptionOptions() *FundingRateSubscriptionOptions {
	return &FundingRateSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithFundingRate sets the funding rate callback
func (o *FundingRateSubscriptionOptions) WithFundingRate(callback FundingRateCallback) *FundingRateSubscriptionOptions {
	o.fundingRateCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *FundingRateSubscriptionOptions) WithConnect(callback func()) *FundingRateSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *FundingRateSubscriptionOptions) WithReconnect(callback func()) *FundingRateSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *FundingRateSubscriptionOptions) WithDisconnect(callback func()) *FundingRateSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *FundingRateSubscriptionOptions) WithError(callback func(error)) *FundingRateSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

// UserDataSubscriptionOptions represents user data subscription options
type UserDataSubscriptionOptions struct {
	*SubscriptionOptions
	executionReportCallback     ExecutionReportCallback
	accountUpdateCallback       OutboundAccountPositionCallback
	balanceUpdateCallback       BalanceUpdateCallback
	listenKeyExpiredCallback    ListenKeyExpiredCallback
	accountUpdateEventCallback  AccountUpdateCallback
	marginCallCallback          MarginCallCallback
	orderTradeUpdateCallback    OrderTradeUpdateCallback
	tradeLiteCallback           TradeLiteCallback
	accountConfigUpdateCallback AccountConfigUpdateCallback
}

// NewUserDataSubscriptionOptions creates a new user data subscription options
func NewUserDataSubscriptionOptions() *UserDataSubscriptionOptions {
	return &UserDataSubscriptionOptions{
		SubscriptionOptions: &SubscriptionOptions{},
	}
}

// WithExecutionReport sets the execution report callback
func (o *UserDataSubscriptionOptions) WithExecutionReport(callback ExecutionReportCallback) *UserDataSubscriptionOptions {
	o.executionReportCallback = callback
	return o
}

// WithAccountUpdate sets the account update callback
func (o *UserDataSubscriptionOptions) WithAccountUpdate(callback OutboundAccountPositionCallback) *UserDataSubscriptionOptions {
	o.accountUpdateCallback = callback
	return o
}

// WithBalanceUpdate sets the balance update callback
func (o *UserDataSubscriptionOptions) WithBalanceUpdate(callback BalanceUpdateCallback) *UserDataSubscriptionOptions {
	o.balanceUpdateCallback = callback
	return o
}

// WithListenKeyExpired sets the listen key expired callback
func (o *UserDataSubscriptionOptions) WithListenKeyExpired(callback ListenKeyExpiredCallback) *UserDataSubscriptionOptions {
	o.listenKeyExpiredCallback = callback
	return o
}

// WithAccountUpdateEvent sets the account update event callback
func (o *UserDataSubscriptionOptions) WithAccountUpdateEvent(callback AccountUpdateCallback) *UserDataSubscriptionOptions {
	o.accountUpdateEventCallback = callback
	return o
}

// WithMarginCall sets the margin call callback
func (o *UserDataSubscriptionOptions) WithMarginCall(callback MarginCallCallback) *UserDataSubscriptionOptions {
	o.marginCallCallback = callback
	return o
}

// WithOrderTradeUpdate sets the order trade update callback
func (o *UserDataSubscriptionOptions) WithOrderTradeUpdate(callback OrderTradeUpdateCallback) *UserDataSubscriptionOptions {
	o.orderTradeUpdateCallback = callback
	return o
}

// WithTradeLite sets the trade lite callback
func (o *UserDataSubscriptionOptions) WithTradeLite(callback TradeLiteCallback) *UserDataSubscriptionOptions {
	o.tradeLiteCallback = callback
	return o
}

// WithAccountConfigUpdate sets the account config update callback
func (o *UserDataSubscriptionOptions) WithAccountConfigUpdate(callback AccountConfigUpdateCallback) *UserDataSubscriptionOptions {
	o.accountConfigUpdateCallback = callback
	return o
}

// WithConnect sets the connect callback
func (o *UserDataSubscriptionOptions) WithConnect(callback func()) *UserDataSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

// WithReconnect sets the reconnect callback
func (o *UserDataSubscriptionOptions) WithReconnect(callback func()) *UserDataSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

// WithDisconnect sets the disconnect callback
func (o *UserDataSubscriptionOptions) WithDisconnect(callback func()) *UserDataSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

// WithError sets the error callback
func (o *UserDataSubscriptionOptions) WithError(callback func(error)) *UserDataSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}
