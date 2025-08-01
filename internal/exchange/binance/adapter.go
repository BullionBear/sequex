package binance

import (
	"context"
	"strconv"

	"github.com/BullionBear/sequex/internal/exchange"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
	"github.com/shopspring/decimal"
)

var _ exchange.IsolatedSpotConnector = (*BinanceExchangeAdapter)(nil)

func NewBinanceAdapter(cfg exchange.Config) *BinanceExchangeAdapter {
	wsClient := binance.NewWSClient(&binance.WSConfig{
		APIKey:      cfg.APIKey,
		APISecret:   cfg.APISecret,
		BaseWsURL:   binance.MainnetWSBaseUrl,
		BaseRestURL: binance.MainnetBaseUrl,
	})
	restClient := wsClient.GetRestClient()
	return &BinanceExchangeAdapter{cfg: cfg, restClient: restClient, wsClient: wsClient}
}

type BinanceExchangeAdapter struct {
	cfg        exchange.Config
	restClient *binance.Client
	wsClient   *binance.WSClient
}

func (a *BinanceExchangeAdapter) GetBalance(ctx context.Context) (exchange.Response[[]exchange.Balance], error) {
	resp, err := a.restClient.GetAccountInfo(ctx, binance.GetAccountInfoRequest{
		OmitZeroBalances: true,
		RecvWindow:       5000,
	})
	if err != nil {
		return exchange.Response[[]exchange.Balance]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[[]exchange.Balance]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	balances := make([]exchange.Balance, len(resp.Data.Balances))
	for i, balance := range resp.Data.Balances {
		balances[i] = exchange.Balance{
			Asset:  balance.Asset,
			Free:   balance.Free,
			Locked: balance.Locked,
		}
	}
	return exchange.Response[[]exchange.Balance]{
		Code:    200,
		Message: "OK",
		Data:    &balances,
	}, nil
}

func (a *BinanceExchangeAdapter) ListOpenOrders(ctx context.Context, symbol exchange.Symbol) (exchange.Response[[]exchange.Order], error) {
	resp, err := a.restClient.ListOpenOrders(ctx, binance.ListOpenOrdersRequest{
		Symbol:     symbol.String(),
		RecvWindow: 5000,
	})
	if err != nil {
		return exchange.Response[[]exchange.Order]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[[]exchange.Order]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	orders := make([]exchange.Order, len(*resp.Data))
	for i, order := range *resp.Data {
		if symbol.String() != order.Symbol {
			continue
		}
		orders[i] = exchange.Order{
			Symbol:      symbol,
			OrderID:     strconv.FormatInt(order.OrderId, 10),
			Price:       order.Price,
			OrigQty:     order.OrigQty,
			Executed:    order.ExecutedQty,
			Status:      toExchangeOrderStatus(order.Status),
			TimeInForce: toExchangeTimeInForce(order.TimeInForce),
			Type:        toExchangeOrderType(order.Type),
			Side:        toExchangeOrderSide(order.Side),
		}
	}
	return exchange.Response[[]exchange.Order]{
		Code:    0,
		Message: "OK",
		Data:    &orders,
	}, nil
}

func (a *BinanceExchangeAdapter) QueryOrder(ctx context.Context, symbol exchange.Symbol, orderID string) (exchange.Response[exchange.Order], error) {
	orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return exchange.Response[exchange.Order]{}, err
	}
	resp, err := a.restClient.QueryOrder(ctx, binance.QueryOrderRequest{
		Symbol:     toBianceSymbol(symbol),
		OrderId:    orderIDInt,
		RecvWindow: 5000,
	})
	if err != nil {
		return exchange.Response[exchange.Order]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[exchange.Order]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	symbol, err = toExchangeSymbol(resp.Data.Symbol)
	if err != nil {
		return exchange.Response[exchange.Order]{}, err
	}
	order := exchange.Order{
		Symbol:      symbol,
		OrderID:     strconv.FormatInt(resp.Data.OrderId, 10),
		Price:       resp.Data.Price,
		OrigQty:     resp.Data.OrigQty,
		Executed:    resp.Data.ExecutedQty,
		Status:      toExchangeOrderStatus(resp.Data.Status),
		TimeInForce: toExchangeTimeInForce(resp.Data.TimeInForce),
		Type:        toExchangeOrderType(resp.Data.Type),
		Side:        toExchangeOrderSide(resp.Data.Side),
		CreatedAt:   resp.Data.Time,
	}
	return exchange.Response[exchange.Order]{
		Code:    0,
		Message: "OK",
		Data:    &order,
	}, nil
}
func (a *BinanceExchangeAdapter) GetMyTrades(ctx context.Context, req exchange.GetMyTradesRequest) (exchange.Response[[]exchange.MyTrade], error) {
	resp, err := a.restClient.GetMyTrades(ctx, binance.GetAccountTradesRequest{
		Symbol:     toBianceSymbol(req.Symbol),
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Limit:      req.Limit,
		RecvWindow: 5000,
	})
	if err != nil {
		return exchange.Response[[]exchange.MyTrade]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[[]exchange.MyTrade]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	trades := make([]exchange.MyTrade, len(*resp.Data))
	for i, trade := range *resp.Data {
		side := exchange.OrderSideSell
		if trade.IsBuyer {
			side = exchange.OrderSideBuy
		}
		trades[i] = exchange.MyTrade{
			Symbol:          req.Symbol,
			TradeID:         strconv.FormatInt(trade.Id, 10),
			OrderID:         strconv.FormatInt(trade.OrderId, 10),
			Price:           trade.Price,
			Quantity:        trade.Qty,
			Commission:      trade.Commission,
			CommissionAsset: trade.CommissionAsset,
			Side:            side,
			IsMaker:         trade.IsMaker,
			Time:            trade.Time,
		}
	}
	return exchange.Response[[]exchange.MyTrade]{
		Code:    0,
		Message: "OK",
		Data:    &trades,
	}, nil
}

func (a *BinanceExchangeAdapter) CreateLimitOrder(ctx context.Context, symbol exchange.Symbol, side exchange.OrderSide, price, quantity decimal.Decimal, timeInForce exchange.TimeInForce) (exchange.Response[exchange.OrderResponseAck], error) {
	resp, err := a.restClient.CreateOrder(ctx, binance.CreateOrderRequest{
		Symbol:           toBianceSymbol(symbol),
		Side:             toBinanceOrderSide(side),
		Type:             binance.OrderTypeLimit,
		Price:            price.String(),
		Quantity:         quantity.String(),
		TimeInForce:      toBinanceTimeInForce(timeInForce),
		NewOrderRespType: binance.NewOrderRespTypeAck,
	})
	if err != nil {
		return exchange.Response[exchange.OrderResponseAck]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[exchange.OrderResponseAck]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	orderID := strconv.FormatInt(resp.Data.OrderId, 10)
	return exchange.Response[exchange.OrderResponseAck]{
		Code:    0,
		Message: "OK",
		Data: &exchange.OrderResponseAck{
			Symbol:  symbol,
			OrderID: orderID,
		},
	}, nil
}

func (a *BinanceExchangeAdapter) CreateLimitMakerOrder(ctx context.Context, symbol exchange.Symbol, side exchange.OrderSide, price, quantity decimal.Decimal) (exchange.Response[exchange.OrderResponseAck], error) {
	resp, err := a.restClient.CreateOrder(ctx, binance.CreateOrderRequest{
		Symbol:           toBianceSymbol(symbol),
		Side:             toBinanceOrderSide(side),
		Price:            price.String(),
		Quantity:         quantity.String(),
		TimeInForce:      binance.TimeInForceGTC,
		Type:             binance.OrderTypeLimitMaker,
		NewOrderRespType: binance.NewOrderRespTypeAck,
	})
	if err != nil {
		return exchange.Response[exchange.OrderResponseAck]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[exchange.OrderResponseAck]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	orderID := strconv.FormatInt(resp.Data.OrderId, 10)
	return exchange.Response[exchange.OrderResponseAck]{
		Code:    0,
		Message: "OK",
		Data: &exchange.OrderResponseAck{
			Symbol:  symbol,
			OrderID: orderID,
		},
	}, nil
}

func (a *BinanceExchangeAdapter) CreateStopOrder(ctx context.Context, symbol exchange.Symbol, side exchange.OrderSide, price, quantity decimal.Decimal) (exchange.Response[exchange.OrderResponseAck], error) {
	resp, err := a.restClient.CreateOrder(ctx, binance.CreateOrderRequest{
		Symbol:           toBianceSymbol(symbol),
		Side:             toBinanceOrderSide(side),
		Price:            price.String(),
		Quantity:         quantity.String(),
		Type:             binance.OrderTypeStopLoss,
		NewOrderRespType: binance.NewOrderRespTypeAck,
	})
	if err != nil {
		return exchange.Response[exchange.OrderResponseAck]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[exchange.OrderResponseAck]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	orderID := strconv.FormatInt(resp.Data.OrderId, 10)
	return exchange.Response[exchange.OrderResponseAck]{
		Code:    0,
		Message: "OK",
		Data: &exchange.OrderResponseAck{
			Symbol:  symbol,
			OrderID: orderID,
		},
	}, nil
}

func (a *BinanceExchangeAdapter) CreateMarketOrder(ctx context.Context, symbol exchange.Symbol, side exchange.OrderSide, quantity decimal.Decimal) (exchange.Response[exchange.OrderResponseAck], error) {
	resp, err := a.restClient.CreateOrder(ctx, binance.CreateOrderRequest{
		Symbol:           toBianceSymbol(symbol),
		Side:             toBinanceOrderSide(side),
		Quantity:         quantity.String(),
		Type:             binance.OrderTypeMarket,
		NewOrderRespType: binance.NewOrderRespTypeAck,
	})
	if err != nil {
		return exchange.Response[exchange.OrderResponseAck]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[exchange.OrderResponseAck]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	orderID := strconv.FormatInt(resp.Data.OrderId, 10)
	return exchange.Response[exchange.OrderResponseAck]{
		Code:    0,
		Message: "OK",
		Data: &exchange.OrderResponseAck{
			Symbol:  symbol,
			OrderID: orderID,
		},
	}, nil
}

func (a *BinanceExchangeAdapter) CancelOrder(ctx context.Context, symbol exchange.Symbol, orderID string) (exchange.Response[exchange.Empty], error) {
	orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return exchange.Response[exchange.Empty]{}, err
	}
	resp, err := a.restClient.CancelOrder(ctx, binance.CancelOrderRequest{
		Symbol:     toBianceSymbol(symbol),
		OrderId:    orderIDInt,
		RecvWindow: 5000,
	})
	if err != nil {
		return exchange.Response[exchange.Empty]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[exchange.Empty]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return exchange.Response[exchange.Empty]{
		Code:    0,
		Message: "OK",
	}, nil
}

func (a *BinanceExchangeAdapter) CancelAllOrders(ctx context.Context, symbol exchange.Symbol) (exchange.Response[exchange.Empty], error) {
	resp, err := a.restClient.CancelAllOrders(ctx, binance.CancelAllOrdersRequest{
		Symbol:     toBianceSymbol(symbol),
		RecvWindow: 5000,
	})
	if err != nil {
		return exchange.Response[exchange.Empty]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[exchange.Empty]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return exchange.Response[exchange.Empty]{
		Code:    0,
		Message: "OK",
	}, nil
}

func (a *BinanceExchangeAdapter) GetDepth(ctx context.Context, symbol exchange.Symbol, limit int) (exchange.Response[exchange.Depth], error) {
	resp, err := a.restClient.GetDepth(ctx, toBianceSymbol(symbol), limit)
	if err != nil {
		return exchange.Response[exchange.Depth]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[exchange.Depth]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return exchange.Response[exchange.Depth]{
		Code:    0,
		Message: "OK",
		Data: &exchange.Depth{
			Symbol: symbol.String(),
			Asks:   resp.Data.Asks,
			Bids:   resp.Data.Bids,
		},
	}, nil
}

func (a *BinanceExchangeAdapter) GetTrades(ctx context.Context, symbol exchange.Symbol, limit int) (exchange.Response[[]exchange.Trade], error) {
	resp, err := a.restClient.GetRecentTrades(ctx, toBianceSymbol(symbol), limit)
	if err != nil {
		return exchange.Response[[]exchange.Trade]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[[]exchange.Trade]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	trades := make([]exchange.Trade, len(*resp.Data))
	for i, trade := range *resp.Data {
		side := exchange.OrderSideBuy
		if trade.IsBuyerMaker {
			side = exchange.OrderSideSell // if Buyer is maker, then Seller is taker
		}
		trades[i] = exchange.Trade{
			Symbol:    symbol,
			ID:        int64(trade.ID),
			Price:     trade.Price,
			Qty:       trade.Qty,
			Time:      trade.Time,
			TakerSide: side,
		}
	}
	return exchange.Response[[]exchange.Trade]{
		Code:    0,
		Message: "OK",
		Data:    &trades,
	}, nil
}
