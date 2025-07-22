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
		BaseURL:   binance.MainnetBaseUrl,
	}
	client := binance.NewClient(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Example: GetDepth (order book)
	depthResp, err := client.GetDepth(ctx, "ADAUSDT", 5)
	if err != nil || depthResp.Code != 0 {
		fmt.Printf("GetDepth error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", depthResp.Message)
	} else {
		fmt.Printf("Order book: %+v\n", depthResp.Data)
	}

	// Example: GetRecentTrades
	tradesResp, err := client.GetRecentTrades(ctx, "ADAUSDT", 5)
	if err != nil || tradesResp.Code != 0 {
		fmt.Printf("GetRecentTrades error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", tradesResp.Message)
	} else {
		fmt.Printf("Recent trades: %+v\n", tradesResp.Data)
	}

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
