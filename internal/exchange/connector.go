package exchange

type Connector interface {
	GetBalance() (Balance, error)
	GetOpenOrders() (OpenOrders, error)
	GetOrder(orderID string) (Order, error)
	GetMyTrades(symbol string) (MyTrade, error)
	CreateOrder(order Order) (Order, error)
	CancelOrder(orderID string) (Order, error)
	CancelAllOrders() (Order, error)

	GetOrderBook(symbol string, limit int) (OrderBook, error)
	GetTrade(symbol string) (Trade, error)
}
