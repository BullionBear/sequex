package binance

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
	"github.com/patrickmn/go-cache"
)

var (
	symbolCache = cache.New(24*time.Hour, 36*time.Hour)
)

func init() {
	_, err := GetSymbol("BTCUSDT")
	if err != nil {
		log.Fatalf("Failed to get symbol: %v", err)
	}
}

func GetSymbol(symbol string) (binance.Symbol, error) {
	if cachedSymbol, found := symbolCache.Get(symbol); found {
		return cachedSymbol.(binance.Symbol), nil
	}
	binanceClient := binance.NewClient(binance.NewMainnetConfig("", ""))
	exchInfo, err := binanceClient.GetExchangeInfo(context.Background(), binance.ExchangeInfoRequest{
		Permissions:  []string{"SPOT"},
		SymbolStatus: "TRADING",
	})
	if err != nil {
		log.Fatalf("Failed to get exchange info: %v", err)
	}
	if exchInfo.Code != 0 {
		log.Fatalf("Failed to get exchange info: %v", exchInfo.Message)
	}
	for _, symbol := range exchInfo.Data.Symbols {
		symbolCache.Set(symbol.Symbol, symbol, cache.DefaultExpiration)
	}
	if cachedSymbol, found := symbolCache.Get(symbol); found {
		return cachedSymbol.(binance.Symbol), nil
	}
	return binance.Symbol{}, fmt.Errorf("symbol not found")
}
