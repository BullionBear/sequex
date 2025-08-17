package binance

import (
	"context"
	"fmt"
	"log"
)

var (
	symbolMap = make(map[string]Symbol)
)

func init() {
	binanceClient := NewClient(NewMainnetConfig("", ""))
	resp, err := binanceClient.GetExchangeInfo(context.Background(), ExchangeInfoRequest{
		Permissions:  []string{"SPOT"},
		SymbolStatus: "TRADING",
	})
	if err != nil {
		log.Fatalf("Failed to get exchange info: %v", err)
	}
	if resp.Code != 0 {
		log.Fatalf("Failed to get exchange info: %v", resp.Message)
	}
	for _, s := range resp.Data.Symbols {
		symbolMap[s.Symbol] = s
	}
}

func GetBaseAsset(symbol string) (string, error) {
	binanceSymbol, ok := symbolMap[symbol]
	if !ok {
		return "", fmt.Errorf("symbol %s not found", symbol)
	}
	return binanceSymbol.BaseAsset, nil
}

func GetQuoteAsset(symbol string) (string, error) {
	binanceSymbol, ok := symbolMap[symbol]
	if !ok {
		return "", fmt.Errorf("symbol %s not found", symbol)
	}
	return binanceSymbol.QuoteAsset, nil
}
