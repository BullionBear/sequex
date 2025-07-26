package exchange

type Connector interface {
	GetBalance() (Balance, error)
	GetOpenOrders() ([]Order, error)
	GetOrder(orderID string) (Order, error)
	GetMyTrades(symbol string) ([]MyTrade, error)

	CreateLimitBuy(symbol string, price, quantity string, timeInForce string) (Order, error)
	CreateLimitSell(symbol string, price, quantity string, timeInForce string) (Order, error)
	CreateMarketBuy(symbol string, quantity string) error
	CreateMarketSell(symbol string, quantity string) error
	CreateLimitMakerBuy(symbol string, price, quantity string) (Order, error)
	CreateLimitMakerSell(symbol string, price, quantity string) (Order, error)
	CreateStopLimitBuy(symbol string, price, quantity string, timeInForce string) (Order, error)
	CreateStopLimitSell(symbol string, price, quantity string, timeInForce string) (Order, error)

	CancelOrder(orderID string) error
	CancelAllOrders() error

	GetOrderBook(symbol string, limit int) (OrderBook, error)
	GetTrades(symbol string) ([]Trade, error)
}
