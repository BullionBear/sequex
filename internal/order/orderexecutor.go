package order

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/BullionBear/sequex/pkg/log"
	"github.com/shopspring/decimal"

	"github.com/BullionBear/sequex/internal/orderbook"
	"github.com/adshao/go-binance/v2"
	"github.com/google/uuid"
)

type BinanceOrderExecutor struct {
	ordbook *orderbook.BinanceOrderBookManager
	svc     *binance.OrderCreateWsService
	orders  sync.Map
	logger  *log.Logger
}

func NewBinanceOrderExecutor(apiKey, apiSecret string, orderbookManager *orderbook.BinanceOrderBookManager, logger *log.Logger) *BinanceOrderExecutor {
	svc, err := binance.NewOrderCreateWsService(apiKey, apiSecret)
	if err != nil {
		logger.Fatal("Error creating Binance OrderCreateWsService: %v", err)
	}
	return &BinanceOrderExecutor{
		ordbook: orderbookManager,
		svc:     svc,
		orders:  sync.Map{},
		logger:  logger,
	}
}

func (bom *BinanceOrderExecutor) PlaceOrder(o *Order) (*OrderResponse, error) {
	switch o.GetType() {
	case OrderTypeMarket:
		marketOrder, ok := o.(*MarketOrder)
		if !ok {
			return nil, errors.New("invalid order type")
		}
		return bom.PlaceMarketOrder(marketOrder)
	case OrderTypeLimit:
		limitOrder, ok := o.(*LimitOrder)
		if !ok {
			return nil, errors.New("invalid order type")
		}
		return bom.PlaceLimitOrder(limitOrder)
}

func (bom *BinanceOrderExecutor) MarketOrder(o *MarketOrder) (*OrderResponse, error) {
	clientID := uuid.NewString()
	req := binance.NewOrderCreateWsRequest().
		Symbol(symbol).
		Side(binance.SideType(side.String())).
		NewClientOrderID(clientID).
		Type(binance.OrderTypeMarket).
		Quantity(quantity.String())
	resp, err := bom.svc.SyncDo(uuid.NewString(), req)
	bom.logger.Info("Market order response: %v", resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}
	orderID := fmt.Sprintf("%d", resp.Result.OrderID)
	return &OrderResponse{
		OrderID:        func(s string) *string { return &s }(orderID),
		CliendtOrderID: &resp.Result.ClientOrderID,
		Status:         toOrderStatus(resp.Result.Status),
		Symbol:         resp.Result.Symbol,
	}, nil
}

func (bom *BinanceOrderExecutor) PlaceLimitOrder(o *LimitOrder) (*OrderResponse, error) {
	clientID := uuid.NewString()
	req := binance.NewOrderCreateWsRequest().
		Symbol(symbol).
		Side(binance.SideType(side.String())).
		NewClientOrderID(clientID).
		Type(binance.OrderTypeLimit).
		Quantity(quantity.String()).
		Price(price.String()).
		TimeInForce(binance.TimeInForceType(tif.String())).
		NewOrderRespType(binance.NewOrderRespTypeRESULT)
	resp, err := bom.svc.SyncDo(uuid.NewString(), req)
	bom.logger.Info("Limit order response: %+v", resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}
	orderID := fmt.Sprintf("%d", resp.Result.OrderID)
	return &OrderResponse{
		OrderID:        func(s string) *string { return &s }(orderID),
		CliendtOrderID: &resp.Result.ClientOrderID,
		Status:         toOrderStatus(resp.Result.Status),
		Symbol:         resp.Result.Symbol,
	}, nil
}

func (bom *BinanceOrderExecutor) CancelOrder(symbol, clientID string) error {
	req := binance.NewOrderCancelWsRequest().
		Symbol(symbol).
		NewClientOrderID(clientID)
	resp, err := bom.svc.SyncDo(uuid.NewString(), req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	return nil
}


type OrderResponse struct {
	OrderID        *string     `json:"order_id"`
	CliendtOrderID *string     `json:"client_order_id"`
	Status         OrderStatus `json:"status"`
	Symbol         string      `json:"symbol"`
}
