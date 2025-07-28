package exchange

import "context"

type Connector interface {
	GetBalance(ctx context.Context) (Response[[]Balance], error)
	ListOpenOrders(ctx context.Context) (Response[[]Order], error)
	GetOrder(ctx context.Context, orderID string) (Response[Order], error)
	GetMyTrades(ctx context.Context, symbol Symbol) (Response[[]MyTrade], error)

	CreateLimitOrder(ctx context.Context, symbol Symbol, side string, price, quantity string, timeInForce string) (Response[Order], error)
	CreateLimitMakerOrder(ctx context.Context, symbol Symbol, side string, price, quantity string) (Response[Order], error)
	CreateStopLimitOrder(ctx context.Context, symbol Symbol, side string, price, quantity string, timeInForce string) (Response[Order], error)
	CreateMarketOrder(ctx context.Context, symbol Symbol, side string, quantity string) (Response[Order], error)

	CancelOrder(ctx context.Context, orderID string) (Response[Order], error)
	CancelAllOrders(ctx context.Context) (Response[Order], error)

	GetOrderBook(ctx context.Context, symbol string, limit int) (Response[OrderBook], error)
	GetTrades(ctx context.Context, symbol string) (Response[[]Trade], error)
	GetPositions(ctx context.Context) (Response[[]Position], error)
	/*
	   SubscribeDepthDiff(symbol string, limit int) (Response[OrderBook], error)
	   SubscribeAggTrades(symbol string) (Response[[]Trade], error)
	   SubscribeMyTrades(symbol string) (Response[[]MyTrade], error)
	*/
}
