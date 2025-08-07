package exchange

import "fmt"

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    *T     `json:"data,omitempty"`
}

type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type Symbol struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

func (s *Symbol) String() string {
	return fmt.Sprintf("%s%s", s.Base, s.Quote)
}

type Order struct {
	Symbol      Symbol      `json:"symbol"`
	OrderID     string      `json:"orderId"`
	Price       string      `json:"price"`
	OrigQty     string      `json:"origQty"`
	Executed    string      `json:"executed"`
	Status      OrderStatus `json:"status"`
	TimeInForce TimeInForce `json:"timeInForce"`
	Type        OrderType   `json:"type"`
	Side        OrderSide   `json:"side"`
	CreatedAt   int64       `json:"time"`
}

type OrderResponseAck struct {
	Symbol  Symbol `json:"symbol"`
	OrderID string `json:"orderId"`
}

type Depth struct {
	Symbol string     `json:"symbol"`
	Asks   [][]string `json:"asks"`
	Bids   [][]string `json:"bids"`
}

type Trade struct {
	Symbol    Symbol    `json:"symbol"`
	ID        int64     `json:"id"`
	Price     string    `json:"price"`
	Qty       string    `json:"qty"`
	Time      int64     `json:"time"`
	TakerSide OrderSide `json:"takerSide"`
}

type GetMyTradesRequest struct {
	Symbol    Symbol `json:"symbol"`
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

type MyTrade struct {
	Symbol          Symbol    `json:"symbol"`
	TradeID         string    `json:"tradeId"`
	OrderID         string    `json:"orderId"`
	Price           string    `json:"price"`
	Quantity        string    `json:"qty"`
	Commission      string    `json:"commission"`
	CommissionAsset string    `json:"commissionAsset"`
	Side            OrderSide `json:"side"`
	IsMaker         bool      `json:"isMaker"`
	Time            int64     `json:"time"`
}

type Position struct {
}

type Empty struct{}

// Stream subscription options
type TradeSubscriptionOptions struct {
	OnConnect    func()            // Called when connection is established
	OnReconnect  func()            // Called when connection is reestablished
	OnError      func(err error)   // Called when an error occurs
	OnTrade      func(trade Trade) // Called when trade data is received
	OnDisconnect func()            // Called when connection is disconnected
}

func (o *TradeSubscriptionOptions) WithConnect(onConnect func()) *TradeSubscriptionOptions {
	o.OnConnect = onConnect
	return o
}

func (o *TradeSubscriptionOptions) WithReconnect(onReconnect func()) *TradeSubscriptionOptions {
	o.OnReconnect = onReconnect
	return o
}

func (o *TradeSubscriptionOptions) WithError(onError func(err error)) *TradeSubscriptionOptions {
	o.OnError = onError
	return o
}

func (o *TradeSubscriptionOptions) WithTrade(onTrade func(trade Trade)) *TradeSubscriptionOptions {
	o.OnTrade = onTrade
	return o
}

func (o *TradeSubscriptionOptions) WithDisconnect(onDisconnect func()) *TradeSubscriptionOptions {
	o.OnDisconnect = onDisconnect
	return o
}
