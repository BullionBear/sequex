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

	// Example: GetAggTrades
	aggTradesResp, err := client.GetAggTrades(ctx, "ADAUSDT", 0, 0, 0, 5)
	if err != nil || aggTradesResp.Code != 0 {
		fmt.Printf("GetAggTrades error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", aggTradesResp.Message)
	} else {
		fmt.Printf("Aggregate trades: %+v\n", aggTradesResp.Data)
	}

	// Example: GetCandles
	candlesResp, err := client.GetCandles(ctx, "ADAUSDT", "1m", 0, 0, "", 5)
	if err != nil || candlesResp.Code != 0 {
		fmt.Printf("GetCandles error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", candlesResp.Message)
	} else {
		fmt.Printf("Klines: %+v\n", candlesResp.Data)
	}

	// Example: GetPriceTicker (single symbol)
	priceResp, err := client.GetPriceTicker(ctx, "ADAUSDT")
	if err != nil || priceResp.Code != 0 {
		fmt.Printf("GetPriceTicker error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", priceResp.Message)
	} else {
		fmt.Printf("Price ticker (single): %+v\n", priceResp.Data)
	}

	// Example: GetPriceTicker (multiple symbols)
	multiResp, err := client.GetPriceTicker(ctx, "ADAUSDT", "BTCUSDT")
	if err != nil || multiResp.Code != 0 {
		fmt.Printf("GetPriceTicker (multi) error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", multiResp.Message)
	} else {
		fmt.Printf("Price ticker (multi): %+v\n", multiResp.Data)
	}

	// Example: CancelOrder (will fail unless you have an open order)
	cancelReq := binance.CancelOrderRequest{
		Symbol:  "ADAUSDT",
		OrderId: 123456789, // Replace with a real orderId
		// OrigClientOrderId: "yourClientOrderId", // Or use this instead
	}
	cancelResp, err := client.CancelOrder(ctx, cancelReq)
	if err != nil || cancelResp.Code != 0 {
		fmt.Printf("CancelOrder error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", cancelResp.Message)
	} else {
		fmt.Printf("CancelOrder result: %+v\n", cancelResp.Data)
	}

	// Example: CancelAllOrders (will fail unless you have open orders)
	cancelAllReq := binance.CancelAllOrdersRequest{
		Symbol: "ADAUSDT",
	}
	cancelAllResp, err := client.CancelAllOrders(ctx, cancelAllReq)
	if err != nil || cancelAllResp.Code != 0 {
		fmt.Printf("CancelAllOrders error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", cancelAllResp.Message)
	} else {
		fmt.Printf("CancelAllOrders result: %+v\n", cancelAllResp.Data)
	}

	// Example: QueryOrder (will fail unless you have a real order)
	queryReq := binance.QueryOrderRequest{
		Symbol:  "ADAUSDT",
		OrderId: 123456789, // Replace with a real orderId
		// OrigClientOrderId: "yourClientOrderId", // Or use this instead
	}
	queryResp, err := client.QueryOrder(ctx, queryReq)
	if err != nil || queryResp.Code != 0 {
		fmt.Printf("QueryOrder error: %v\n", err)
		fmt.Printf("Response Error: %+v\n", queryResp.Message)
	} else {
		fmt.Printf("QueryOrder result: %+v\n", queryResp.Data)
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
