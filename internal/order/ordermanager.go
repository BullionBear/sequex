package order

import (
	"fmt"
	"sync"

	"github.com/BullionBear/sequex/pkg/log"
	"github.com/shopspring/decimal"

	"github.com/BullionBear/sequex/internal/orderbook"
)

type BinanceOrderManager struct {
	binanceOrderManagerMap sync.Map
	ordbookManager         *orderbook.BinanceOrderBookManager
	logger                 *log.Logger
}

func NewBinanceOrderManager(ordbookManager *orderbook.BinanceOrderBookManager) *BinanceOrderManager {
	return &BinanceOrderManager{
		binanceOrderManagerMap: sync.Map{},
		ordbookManager:         ordbookManager,
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

func (bom *BinanceOrderManager) PlaceMarketOrder(accountName string, symbol string, qty decimal.Decimal) (OrderResponse, error) {

}
