package binance

import (
	"context"
	"log"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

var (
	symbolMap = map[string]binance.Symbol{}
)

func init() {
	symbolMap = make(map[string]binance.Symbol)
	binanceClient := binance.NewClient(binance.NewMainnetConfig("", ""))
	exchInfo, err := binanceClient.GetExchangeInfo(context.Background(), binance.ExchangeInfoRequest{})
	if err != nil {
		log.Fatalf("Failed to get exchange info: %v", err)
	}
	for _, symbol := range exchInfo.Data.Symbols {
		symbolMap[symbol.Symbol] = symbol
	}
}
