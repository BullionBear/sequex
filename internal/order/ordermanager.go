package order

import (
	"fmt"
	"sync"

	"github.com/BullionBear/sequex/pkg/log"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/BullionBear/sequex/internal/orderbook"
)

type BinanceOrderManager struct {
	binanceOrderManagerMap sync.Map
	ordbookManager         *orderbook.BinanceOrderBookManager
	logger                 *log.Logger
}

func NewBinanceOrderManager(ordbookManager *orderbook.BinanceOrderBookManager, logger *log.Logger) *BinanceOrderManager {
	return &BinanceOrderManager{
		binanceOrderManagerMap: sync.Map{},
		ordbookManager:         ordbookManager,
		logger:                 logger,
	}
}

func (bom *BinanceOrderManager) Register(accountName, apiKey, apiSecret string) error {
	if _, exists := bom.binanceOrderManagerMap.Load(accountName); exists {
		return fmt.Errorf("account %s already registered", accountName)
	}
	orderExecutor := NewBinanceOrderExecutor(apiKey, apiSecret, bom.ordbookManager, bom.logger.WithKV(log.KV{Key: "account", Value: accountName}))
	bom.binanceOrderManagerMap.Store(accountName, orderExecutor)
	return nil
}

func (bom *BinanceOrderManager) PlaceMarketOrder(accountName string, symbol string, qty decimal.Decimal) (*OrderResponse, error) {
	orderExecutor, ok := bom.binanceOrderManagerMap.Load(accountName)
	if !ok {
		return nil, fmt.Errorf("account %s not registered", accountName)
	}
	order := MarketOrder{
		ClientOrderID: uuid.New().String(),
		Symbol:        symbol,
		Side:          SideBuy,
		Quantity:      qty,
	}
	return orderExecutor.(*BinanceOrderExecutor).PlaceOrder(order)
}

func (bom *BinanceOrderManager) PlaceLimitOrder(accountName string, symbol string, qty decimal.Decimal, price decimal.Decimal) (*OrderResponse, error) {
	orderExecutor, ok := bom.binanceOrderManagerMap.Load(accountName)
	if !ok {
		return nil, fmt.Errorf("account %s not registered", accountName)
	}
	order := LimitOrder{
		ClientOrderID: uuid.New().String(),
		Symbol:        symbol,
		Side:          SideBuy,
		Quantity:      qty,
		Price:         price,
		TimeInForce:   TimeInForceGTC,
	}
	return orderExecutor.(*BinanceOrderExecutor).PlaceOrder(&order)
}

func (bom *BinanceOrderManager) PlaceStopMarketOrder(accountName string, symbol string, qty decimal.Decimal, stopPrice decimal.Decimal) (*OrderResponse, error) {
	orderExecutor, ok := bom.binanceOrderManagerMap.Load(accountName)
	if !ok {
		return nil, fmt.Errorf("account %s not registered", accountName)
	}
	order := StopMarketOrder{
		ClientOrderID: uuid.New().String(),
		Symbol:        symbol,
		Side:          SideBuy,
		Quantity:      qty,
		StopPrice:     stopPrice,
	}
	return orderExecutor.(*BinanceOrderExecutor).PlaceOrder(order)
}
