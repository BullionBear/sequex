package exchange

type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type Order struct {
	Symbol      string `json:"symbol"`
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
