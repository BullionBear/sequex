package order

import (
	"sync"

	"github.com/BullionBear/sequex/internal/orderbook"
	"github.com/adshao/go-binance/v2"
	"github.com/google/uuid"
)

type BinanceOrderManager struct {
	ordbook *orderbook.BinanceOrderBookManager
	svc     *binance.OrderCreateWsService
	orders  sync.Map
}

func NewBinanceOrderManager(apiKey, apiSecret string) *BinanceOrderManager {
	svc, err := binance.NewOrderCreateWsService(apiKey, apiSecret)
	if err != nil {
		panic(err)
	}
	return &BinanceOrderManager{
		ordbook: orderbook.NewBinanceOrderBookManager(),
		svc:     svc,
		orders:  sync.Map{},
	}
}

func (bom *BinanceOrderManager) MarketOrder(o MarketOrder) (string, error) {
	req := binance.NewOrderCreateWsRequest().
		Symbol(o.Symbol).
		Side(binance.SideType(o.Side.String())).
		NewClientOrderID(o.ClientOrderID).
		Type(binance.OrderTypeMarket).
		Quantity(o.Quantity.String())
	resp, err := bom.svc.SyncDo(uuid.NewString(), req)
	if err != nil {
		return "", err
	}
	return resp.Result.ClientOrderID, nil
}

func (bom *BinanceOrderManager) LimitOrder(o LimitOrder) (string, error) {
	req := binance.NewOrderCreateWsRequest().
		Symbol(o.Symbol).
		Side(binance.SideType(o.Side.String())).
		NewClientOrderID(o.ClientOrderID).
		Type(binance.OrderTypeLimit).
		Quantity(o.Quantity.String()).
		Price(o.Price.String()).
		TimeInForce(binance.TimeInForceType(o.TimeInForce.String()))
	resp, err := bom.svc.SyncDo(uuid.NewString(), req)
	if err != nil {
		return "", err
	}
	return resp.Result.ClientOrderID, nil
}
