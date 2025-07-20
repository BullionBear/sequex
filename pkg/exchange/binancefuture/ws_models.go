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
