package order

import (
	"context"
	"sync"

	"github.com/BullionBear/sequex/internal/orderbook"
	"github.com/adshao/go-binance/v2"
	evbus "github.com/asaskevich/EventBus"
)

type BinanceOrderManager struct {
	ordbook  *orderbook.BinanceOrderBookManager
	client   *binance.Client
	orders   sync.Map
	eventBus evbus.Bus
}

func NewBinanceOrderManager(apiKey, apiSecret string, ordbookManager *orderbook.BinanceOrderBookManager) *BinanceOrderManager {
	client := binance.NewClient(apiKey, apiSecret)
	return &BinanceOrderManager{
		ordbook:  ordbookManager,
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
	return o.ClientOrderID, nil
}

func (bom *BinanceOrderManager) StopMarketOrder(o StopMarketOrder) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	bom.orders.Store(o.ClientOrderID, cancel)
	go func() {
		doneC := make(chan struct{})
		unsubscribe, err := bom.ordbook.SubscribeBestDepth(o.Symbol, func(ask, bid orderbook.PriceLevel) {
			if o.OnBestDepth(ask, bid) {
				close(doneC)
			}
		})
		if err != nil {
			return
		}
		select {
		case <-doneC:
			unsubscribe()
			bom.orders.Delete(o.ClientOrderID)
			_, err = bom.MarketOrder(MarketOrder{
				Symbol:        o.Symbol,
				ClientOrderID: o.ClientOrderID,
				Side:          o.Side,
				Quantity:      o.Quantity,
			})
			if err != nil {
				return
			}
			return
		case <-ctx.Done():
			unsubscribe()
			bom.orders.Delete(o.ClientOrderID)
			return
		}
	}()

	return o.ClientOrderID, nil
}
