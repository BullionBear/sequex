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
	ordbook  *orderbook.BinanceOrderBookManager
	tradeSvc *binance.OrderCreateWsService
	orders   sync.Map
	logger   *log.Logger
}

func NewBinanceOrderExecutor(apiKey, apiSecret string, orderbookManager *orderbook.BinanceOrderBookManager, logger *log.Logger) *BinanceOrderExecutor {
	tradeSvc, err := binance.NewOrderCreateWsService(apiKey, apiSecret)
	if err != nil {
		logger.Fatal("Error creating Binance OrderCreateWsService: %v", err)
	}

	return &BinanceOrderExecutor{
		ordbook:  orderbookManager,
		tradeSvc: tradeSvc,
		orders:   sync.Map{},
		logger:   logger,
	}
}

func (bom *BinanceOrderExecutor) PlaceOrder(o Order) (*OrderResponse, error) {
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
	case OrderTypeStopMarket:
		stopMarketOrder, ok := o.(*StopMarketOrder)
		if !ok {
			return nil, errors.New("invalid order type")
		}
		return bom.PlaceStopMarketOrder(stopMarketOrder)
	default:
		return nil, errors.New("unsupported order type")
	}
}

func (bom *BinanceOrderExecutor) PlaceMarketOrder(o *MarketOrder) (*OrderResponse, error) {
	localID := uuid.NewString()
	req := binance.NewOrderCreateWsRequest().
		Symbol(o.Symbol).
		Side(binance.SideType(o.Side.String())).
		NewClientOrderID(localID).
		Type(binance.OrderTypeMarket).
		Quantity(o.Quantity.String())
	resp, err := bom.tradeSvc.SyncDo(uuid.NewString(), req)
	bom.logger.Info("Market order response: %v", resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}
	return &OrderResponse{
		SequexID: resp.Result.ClientOrderID,
		Status:   toOrderStatus(resp.Result.Status),
		Symbol:   resp.Result.Symbol,
	}, nil
}

func (bom *BinanceOrderExecutor) PlaceLimitOrder(o *LimitOrder) (*OrderResponse, error) {
	localID := uuid.NewString()
	req := binance.NewOrderCreateWsRequest().
		Symbol(o.Symbol).
		Side(binance.SideType(o.Side.String())).
		NewClientOrderID(localID).
		Type(binance.OrderTypeLimit).
		Quantity(o.Quantity.String()).
		Price(o.Price.String()).
		TimeInForce(binance.TimeInForceType(o.TimeInForce.String())).
		NewOrderRespType(binance.NewOrderRespTypeRESULT)
	bom.logger.Info("Request params %+v", req.GetParams())
	resp, err := bom.tradeSvc.SyncDo(uuid.NewString(), req)
	bom.logger.Info("Limit order response: %+v", resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}
	return &OrderResponse{
		SequexID: localID,
		Status:   toOrderStatus(resp.Result.Status),
		Symbol:   resp.Result.Symbol,
	}, nil
}

func (bom *BinanceOrderExecutor) PlaceStopMarketOrder(o *StopMarketOrder) (*OrderResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	localID := uuid.NewString()
	once := &sync.Once{}
	unsubscribe, err := bom.ordbook.SubscribeBestDepth(o.Symbol, func(ask, bid orderbook.PriceLevel) {
		if o.OnBestDepth(ask, bid) {
			once.Do(func() {
				cancel()
				bom.PlaceMarketOrder(&MarketOrder{
					Symbol:        o.Symbol,
					Side:          o.Side,
					Quantity:      o.Quantity,
					ClientOrderID: o.ClientOrderID,
				})
			})
		}
	})
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		unsubscribe()
	}()
	return &OrderResponse{
		SequexID: localID,
		Status:   OrderStatusLocalPending,
		Symbol:   o.Symbol,
	}, nil
}
