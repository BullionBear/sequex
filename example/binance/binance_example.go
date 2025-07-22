package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		fmt.Println("Please set BINANCE_API_KEY and BINANCE_API_SECRET environment variables.")
		return
	}
	cfg := &binance.Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   binance.BinanceMainnetBaseUrl,
	}
	client := binance.NewClient(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderReq := binance.CreateOrderRequest{
		Symbol:           "ADAUSDT",
		Side:             binance.OrderSideSell,
		Type:             binance.OrderTypeMarket,
		QuoteOrderQty:    "10", // Buy $100 worth of BTC
		NewOrderRespType: "RESULT",
	}
	resp, err := client.CreateOrder(ctx, orderReq)
	if err != nil || resp.Code != 0 {
		fmt.Printf("CreateOrder error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", resp.Message)
		return
	}
	fmt.Printf("Order placed successfully: %+v\n", resp.Data)
}
