package order

import (
	"context"
	"sync"

	"github.com/BullionBear/sequex/internal/orderbook"
	"github.com/adshao/go-binance/v2"
	evbus "github.com/asaskevich/EventBus"
	"github.com/google/uuid"
)

type BinanceOrderManager struct {
	ordbook  *orderbook.BinanceOrderBookManager
	client   *binance.Client
	orders   sync.Map
	eventBus evbus.Bus
}

func NewBinanceOrderManager(apiKey, apiSecret string) *BinanceOrderManager {
	client := binance.NewClient(apiKey, apiSecret)
	return &BinanceOrderManager{
		ordbook:  orderbook.NewBinanceOrderBookManager(),
		client:   client,
		orders:   sync.Map{},
		eventBus: evbus.New(),
	}
}

func (bom *BinanceOrderManager) MarketOrder(o MarketOrder) (string, error) {
	_, err := bom.client.NewCreateOrderService().
		Symbol(o.Symbol).
		Side(binance.SideType(o.Side.String())).
		NewClientOrderID(o.ClientOrderID).
		Type(binance.OrderTypeMarket).
		Quantity(o.Quantity.String()).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return o.ClientOrderID, nil
}

func (bom *BinanceOrderManager) LimitOrder(o LimitOrder) (string, error) {
	localID := uuid.New().String()
	_, err := bom.client.NewCreateOrderService().
		Symbol(o.Symbol).
		Side(binance.SideType(o.Side.String())).
		NewClientOrderID(o.ClientOrderID).
		Type(binance.OrderTypeLimit).
		Quantity(o.Quantity.String()).
		Price(o.Price.String()).
		TimeInForce(binance.TimeInForceType(o.TimeInForce.String())).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return localID, nil
}
