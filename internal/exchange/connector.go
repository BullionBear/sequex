package exchange

import (
	"context"

	"github.com/shopspring/decimal"
)

type IsolatedConnector interface {
	GetBalance(ctx context.Context) (Response[[]Balance], error)
	ListOpenOrders(ctx context.Context) (Response[[]Order], error)
	QueryOrder(ctx context.Context, symbol Symbol, orderID string) (Response[Order], error)
	GetMyTrades(ctx context.Context, req GetMyTradesRequest) (Response[[]MyTrade], error)

	CreateLimitOrder(ctx context.Context, symbol Symbol, side OrderSide, price, quantity decimal.Decimal, timeInForce TimeInForce) (Response[string], error)
	CreateLimitMakerOrder(ctx context.Context, symbol Symbol, side OrderSide, price, quantity decimal.Decimal) (Response[string], error)
	CreateStopOrder(ctx context.Context, symbol Symbol, side OrderSide, price, quantity decimal.Decimal) (Response[string], error)
	CreateMarketOrder(ctx context.Context, symbol Symbol, side OrderSide, quantity decimal.Decimal) (Response[string], error)

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
