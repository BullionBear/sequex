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

type OrderBook struct {
}

type Trade struct {
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
