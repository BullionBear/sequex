package binance

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

// Type-specific callback functions
type KlineCallback func(data *WSKlineData) error
type TickerCallback func(data *WSTickerData) error
type MiniTickerCallback func(data *WSMiniTickerData) error
type BookTickerCallback func(data *WSBookTickerData) error
type PartialDepthCallback func(data *WSPartialDepthData) error
type DiffDepthCallback func(data *WSDiffDepthData) error
type DepthCallback func(data *WSDepthData) error
type TradeCallback func(data *WSTradeData) error
type AggTradeCallback func(data *WSAggTradeData) error
type OutboundAccountPositionCallback func(data *WSOutboundAccountPosition) error
type BalanceUpdateCallback func(data *WSBalanceUpdate) error
type ExecutionReportCallback func(data *WSExecutionReport) error

// webSocketCallback represents a callback function for WebSocket events (internal use only)
type webSocketCallback func(data []byte) error

// WebSocketStream represents a WebSocket stream configuration
type WebSocketStream struct {
	StreamName string
	Callback   webSocketCallback
}

// WebSocketConnection represents a WebSocket connection configuration
type WebSocketConnection struct {
	URL      string
	Streams  []WebSocketStream
	IsActive bool
	Close    chan struct{}
}

// SubscriptionOptions represents the base subscription options with common callbacks
type SubscriptionOptions struct {
	onConnect    func()
	onReconnect  func()
	onDisconnect func()
	onError      func(error)
}

// KlineSubscriptionOptions represents subscription options for kline data
type KlineSubscriptionOptions struct {
	SubscriptionOptions
	onKline KlineCallback
}

// TickerSubscriptionOptions represents subscription options for ticker data
type TickerSubscriptionOptions struct {
	SubscriptionOptions
	onTicker TickerCallback
}

// MiniTickerSubscriptionOptions represents subscription options for mini ticker data
type MiniTickerSubscriptionOptions struct {
	SubscriptionOptions
	onMiniTicker MiniTickerCallback
}

// BookTickerSubscriptionOptions represents subscription options for book ticker data
type BookTickerSubscriptionOptions struct {
	SubscriptionOptions
	onBookTicker BookTickerCallback
}

// DepthSubscriptionOptions represents subscription options for depth data
type DepthSubscriptionOptions struct {
	SubscriptionOptions
	onDepth DepthCallback
}

// TradeSubscriptionOptions represents subscription options for trade data
type TradeSubscriptionOptions struct {
	SubscriptionOptions
	onTrade TradeCallback
}

// AggTradeSubscriptionOptions represents subscription options for aggregated trade data
type AggTradeSubscriptionOptions struct {
	SubscriptionOptions
	onAggTrade AggTradeCallback
}

// UserDataSubscriptionOptions represents subscription options for user data streams
type UserDataSubscriptionOptions struct {
	SubscriptionOptions
	onExecutionReport ExecutionReportCallback
	onAccountUpdate   OutboundAccountPositionCallback
	onBalanceUpdate   BalanceUpdateCallback
}

// Base subscription option methods
func (o *SubscriptionOptions) WithConnect(callback func()) *SubscriptionOptions {
	o.onConnect = callback
	return o
}

func (o *SubscriptionOptions) WithReconnect(callback func()) *SubscriptionOptions {
	o.onReconnect = callback
	return o
}

func (o *SubscriptionOptions) WithDisconnect(callback func()) *SubscriptionOptions {
	o.onDisconnect = callback
	return o
}

func (o *SubscriptionOptions) WithError(callback func(error)) *SubscriptionOptions {
	o.onError = callback
	return o
}

// Kline subscription option methods
func (o *KlineSubscriptionOptions) WithConnect(callback func()) *KlineSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

func (o *KlineSubscriptionOptions) WithReconnect(callback func()) *KlineSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

func (o *KlineSubscriptionOptions) WithDisconnect(callback func()) *KlineSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

func (o *KlineSubscriptionOptions) WithError(callback func(error)) *KlineSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

func (o *KlineSubscriptionOptions) WithKline(callback KlineCallback) *KlineSubscriptionOptions {
	o.onKline = callback
	return o
}

// Ticker subscription option methods
func (o *TickerSubscriptionOptions) WithConnect(callback func()) *TickerSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

func (o *TickerSubscriptionOptions) WithReconnect(callback func()) *TickerSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

func (o *TickerSubscriptionOptions) WithDisconnect(callback func()) *TickerSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

func (o *TickerSubscriptionOptions) WithError(callback func(error)) *TickerSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

func (o *TickerSubscriptionOptions) WithTicker(callback TickerCallback) *TickerSubscriptionOptions {
	o.onTicker = callback
	return o
}

// MiniTicker subscription option methods
func (o *MiniTickerSubscriptionOptions) WithConnect(callback func()) *MiniTickerSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

func (o *MiniTickerSubscriptionOptions) WithReconnect(callback func()) *MiniTickerSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

func (o *MiniTickerSubscriptionOptions) WithDisconnect(callback func()) *MiniTickerSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

func (o *MiniTickerSubscriptionOptions) WithError(callback func(error)) *MiniTickerSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

func (o *MiniTickerSubscriptionOptions) WithMiniTicker(callback MiniTickerCallback) *MiniTickerSubscriptionOptions {
	o.onMiniTicker = callback
	return o
}

// BookTicker subscription option methods
func (o *BookTickerSubscriptionOptions) WithConnect(callback func()) *BookTickerSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

func (o *BookTickerSubscriptionOptions) WithReconnect(callback func()) *BookTickerSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

func (o *BookTickerSubscriptionOptions) WithDisconnect(callback func()) *BookTickerSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

func (o *BookTickerSubscriptionOptions) WithError(callback func(error)) *BookTickerSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

func (o *BookTickerSubscriptionOptions) WithBookTicker(callback BookTickerCallback) *BookTickerSubscriptionOptions {
	o.onBookTicker = callback
	return o
}

// Depth subscription option methods
func (o *DepthSubscriptionOptions) WithConnect(callback func()) *DepthSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

func (o *DepthSubscriptionOptions) WithReconnect(callback func()) *DepthSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

func (o *DepthSubscriptionOptions) WithDisconnect(callback func()) *DepthSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

func (o *DepthSubscriptionOptions) WithError(callback func(error)) *DepthSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

func (o *DepthSubscriptionOptions) WithDepth(callback DepthCallback) *DepthSubscriptionOptions {
	o.onDepth = callback
	return o
}

// Trade subscription option methods
func (o *TradeSubscriptionOptions) WithConnect(callback func()) *TradeSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

func (o *TradeSubscriptionOptions) WithReconnect(callback func()) *TradeSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

func (o *TradeSubscriptionOptions) WithDisconnect(callback func()) *TradeSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

func (o *TradeSubscriptionOptions) WithError(callback func(error)) *TradeSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

func (o *TradeSubscriptionOptions) WithTrade(callback TradeCallback) *TradeSubscriptionOptions {
	o.onTrade = callback
	return o
}

// AggTrade subscription option methods
func (o *AggTradeSubscriptionOptions) WithConnect(callback func()) *AggTradeSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

func (o *AggTradeSubscriptionOptions) WithReconnect(callback func()) *AggTradeSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

func (o *AggTradeSubscriptionOptions) WithDisconnect(callback func()) *AggTradeSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

func (o *AggTradeSubscriptionOptions) WithError(callback func(error)) *AggTradeSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

func (o *AggTradeSubscriptionOptions) WithAggTrade(callback AggTradeCallback) *AggTradeSubscriptionOptions {
	o.onAggTrade = callback
	return o
}

// UserData subscription option methods
func (o *UserDataSubscriptionOptions) WithConnect(callback func()) *UserDataSubscriptionOptions {
	o.SubscriptionOptions.WithConnect(callback)
	return o
}

func (o *UserDataSubscriptionOptions) WithReconnect(callback func()) *UserDataSubscriptionOptions {
	o.SubscriptionOptions.WithReconnect(callback)
	return o
}

func (o *UserDataSubscriptionOptions) WithDisconnect(callback func()) *UserDataSubscriptionOptions {
	o.SubscriptionOptions.WithDisconnect(callback)
	return o
}

func (o *UserDataSubscriptionOptions) WithError(callback func(error)) *UserDataSubscriptionOptions {
	o.SubscriptionOptions.WithError(callback)
	return o
}

func (o *UserDataSubscriptionOptions) WithExecutionReport(callback ExecutionReportCallback) *UserDataSubscriptionOptions {
	o.onExecutionReport = callback
	return o
}

func (o *UserDataSubscriptionOptions) WithAccountUpdate(callback OutboundAccountPositionCallback) *UserDataSubscriptionOptions {
	o.onAccountUpdate = callback
	return o
}

func (o *UserDataSubscriptionOptions) WithBalanceUpdate(callback BalanceUpdateCallback) *UserDataSubscriptionOptions {
	o.onBalanceUpdate = callback
	return o
}
