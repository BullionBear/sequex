package order

import (
	"context"
	"errors"
	"sync"

	"github.com/BullionBear/sequex/pkg/log"

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

func (bom *BinanceOrderExecutor) MarketOrder(o MarketOrder) (string, error) {
	req := binance.NewOrderCreateWsRequest().
		Symbol(o.Symbol).
		Side(binance.SideType(o.Side.String())).
		NewClientOrderID(o.ClientOrderID).
		Type(binance.OrderTypeMarket).
		Quantity(o.Quantity.String())
	resp, err := bom.svc.SyncDo(uuid.NewString(), req)
	bom.logger.Info("Market order response: %v", resp)
	if err != nil {
		return "", err
	}

	return resp.Result.ClientOrderID, nil
}

func (bom *BinanceOrderExecutor) LimitOrder(o LimitOrder) (string, error) {
	req := binance.NewOrderCreateWsRequest().
		Symbol(o.Symbol).
		Side(binance.SideType(o.Side.String())).
		NewClientOrderID(o.ClientOrderID).
		Type(binance.OrderTypeLimit).
		Quantity(o.Quantity.String()).
		Price(o.Price.String()).
		TimeInForce(binance.TimeInForceType(o.TimeInForce.String())).
		NewOrderRespType(binance.NewOrderRespTypeRESULT)
	resp, err := bom.svc.SyncDo(uuid.NewString(), req)
	bom.logger.Info("Limit order response: %+v", resp)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", errors.New(resp.Error.Message)
	}
	return resp.Result.ClientOrderID, nil
}

func (bom *BinanceOrderExecutor) StopMarketOrder(o StopMarketOrder) (string, error) {
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

type OrderResponse struct {
	OrderID        *string `json:"order_id"`
	CliendtOrderID *string `json:"client_order_id"`
	Status         string  `json:"status"`
	Symbol         string  `json:"symbol"`
}
