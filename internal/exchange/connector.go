package exchange

type SpotConnector interface {
	GetBalance() (Balance, error)
	GetOpenOrders() ([]Order, error)
	GetOrder(orderID string) (Order, error)
	GetMyTrades(symbol string) ([]MyTrade, error)

	CreateLimitOrder(symbol string, side string, price, quantity string, timeInForce string) (Order, error)
	CreateLimitMakerOrder(symbol string, side string, price, quantity string) (Order, error)
	CreateStopLimitOrder(symbol string, side string, price, quantity string, timeInForce string) (Order, error)
	CreateMarketOrder(symbol string, side string, quantity string) (Order, error)

	CancelOrder(orderID string) error
	CancelAllOrders() error

	GetOrderBook(symbol string, limit int) (OrderBook, error)
	GetTrades(symbol string) ([]Trade, error)
}

type PerpConnector interface {
	GetBalance() (Balance, error)
	GetOpenOrders() ([]Order, error)
	GetOrder(orderID string) (Order, error)
	GetMyTrades(symbol string) ([]MyTrade, error)

	CreateLimitOrder(symbol string, side string, price, quantity string, timeInForce string) (Order, error)
	CreateLimitMakerOrder(symbol string, side string, price, quantity string) (Order, error)
	CreateStopLimitOrder(symbol string, side string, price, quantity string, timeInForce string) (Order, error)
	CreateMarketOrder(symbol string, side string, quantity string) (Order, error)

	CancelOrder(orderID string) error
	CancelAllOrders() error

	GetOrderBook(symbol string, limit int) (OrderBook, error)
	GetTrades(symbol string) ([]Trade, error)
	GetPositions() ([]Position, error)
}

type UnifiedConnector interface {
	GetBalance() (Balance, error)
	GetOpenOrders() ([]Order, error)
	GetOrder(orderID string) (Order, error)
	GetMyTrades(symbol string) ([]MyTrade, error)
}
