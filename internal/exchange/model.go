package exchange

import "fmt"

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    T      `json:"data"`
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
	Symbol      Symbol `json:"symbol"`
	OrderID     int64  `json:"orderId"`
	Price       string `json:"price"`
	OrigQty     string `json:"origQty"`
	Executed    string `json:"executed"`
	Status      string `json:"status"`
	TimeInForce string `json:"timeInForce"`
	Type        string `json:"type"`
	Side        string `json:"side"`
	CreatedAt   int64  `json:"time"`
}

type OrderBook struct {
}

type Trade struct {
}

type MyTrade struct {
}

type Position struct {
}
