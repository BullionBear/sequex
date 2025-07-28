package exchange

type Connector interface {
	GetBalance() (Response[Balance], error)
	GetOpenOrders() (Response[[]Order], error)
	GetOrder(orderID string) (Response[Order], error)
	GetMyTrades(symbol string) (Response[[]MyTrade], error)

	CreateLimitOrder(symbol string, side string, price, quantity string, timeInForce string) (Response[Order], error)
	CreateLimitMakerOrder(symbol string, side string, price, quantity string) (Response[Order], error)
	CreateStopLimitOrder(symbol string, side string, price, quantity string, timeInForce string) (Response[Order], error)
	CreateMarketOrder(symbol string, side string, quantity string) (Response[Order], error)

	CancelOrder(orderID string) (Response[Order], error)
	CancelAllOrders() (Response[Order], error)

	GetOrderBook(symbol string, limit int) (Response[OrderBook], error)
	GetTrades(symbol string) (Response[[]Trade], error)
	GetPositions() (Response[[]Position], error)
}
