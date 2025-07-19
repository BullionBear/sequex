package binancefuture

import (
	"context"
	"fmt"
	"log"
	"time"
)

func ExampleNewClient() {
	// Create a client with testnet configuration
	config := TestnetConfig()
	client := NewClient(config)

	fmt.Printf("Client created with base URL: %s\n", client.GetConfig().BaseURL)
	// Output: Client created with base URL: https://testnet.binancefuture.com
}

func ExampleClient_GetServerTime() {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverTime, err := client.GetServerTime(ctx)
	if err != nil {
		log.Fatalf("Failed to get server time: %v", err)
	}

	fmt.Printf("Server time: %v\n", serverTime.GetTime())
}

func ExampleClient_GetExchangeInfo() {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	exchangeInfo, err := client.GetExchangeInfo(ctx)
	if err != nil {
		log.Fatalf("Failed to get exchange info: %v", err)
	}

	fmt.Printf("Exchange has %d symbols\n", len(exchangeInfo.Symbols))

	// Find BTCUSDT symbol info
	for _, symbol := range exchangeInfo.Symbols {
		if symbol.Symbol == "BTCUSDT" {
			fmt.Printf("BTCUSDT base asset: %s, quote asset: %s\n", symbol.BaseAsset, symbol.QuoteAsset)
			break
		}
	}
}

func ExampleClient_GetTickerPrice() {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get price for a specific symbol
	result, err := client.GetTickerPrice(ctx, "BTCUSDT")
	if err != nil {
		log.Fatalf("Failed to get ticker price: %v", err)
	}

	if result.IsSingle() {
		ticker := result.GetSingle()
		fmt.Printf("BTCUSDT price: %s\n", ticker.Price)
	}
}

func ExampleClient_GetKlines() {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get 1-hour klines for BTCUSDT
	klines, err := client.GetKlines(ctx, "BTCUSDT", Interval1h, 5)
	if err != nil {
		log.Fatalf("Failed to get klines: %v", err)
	}

	if len(*klines) > 0 {
		latest := (*klines)[len(*klines)-1]
		fmt.Printf("Latest BTCUSDT 1h kline - Open: %s, High: %s, Low: %s, Close: %s\n",
			latest.Open, latest.High, latest.Low, latest.Close)
	}
}

func ExampleClient_GetOrderBook() {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get order book for BTCUSDT (top 5 levels)
	orderBook, err := client.GetOrderBook(ctx, "BTCUSDT", 5)
	if err != nil {
		log.Fatalf("Failed to get order book: %v", err)
	}

	if len(orderBook.Bids) > 0 {
		bestBid := orderBook.Bids[0]
		fmt.Printf("Best bid: %s @ %s\n", bestBid[1], bestBid[0])
	}

	if len(orderBook.Asks) > 0 {
		bestAsk := orderBook.Asks[0]
		fmt.Printf("Best ask: %s @ %s\n", bestAsk[1], bestAsk[0])
	}
}

func ExampleClient_GetMarkPrice() {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get mark price for BTCUSDT
	markPrices, err := client.GetMarkPrice(ctx, "BTCUSDT")
	if err != nil {
		log.Fatalf("Failed to get mark price: %v", err)
	}

	if len(markPrices) > 0 {
		markPrice := markPrices[0]
		fmt.Printf("BTCUSDT mark price: %s, funding rate: %s\n",
			markPrice.MarkPrice, markPrice.LastFundingRate)
	}
}

func ExampleClient_GetFundingRate() {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get funding rate history for BTCUSDT
	fundingRates, err := client.GetFundingRate(ctx, "BTCUSDT", 3)
	if err != nil {
		log.Fatalf("Failed to get funding rate: %v", err)
	}

	for _, rate := range fundingRates {
		fmt.Printf("Funding rate at %v: %s\n",
			time.Unix(rate.FundingTime/1000, 0), rate.FundingRate)
	}
}

func ExampleClient_GetAccount() {
	config := TestnetConfig()
	// Set API credentials for signed endpoints
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	account, err := client.GetAccount(ctx)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}

	fmt.Printf("Account has %d assets and %d positions\n", len(account.Assets), len(account.Positions))

	// Print total wallet balance
	fmt.Printf("Total wallet balance: %s\n", account.TotalWalletBalance)

	// Print positions with non-zero amounts
	for _, position := range account.Positions {
		if position.PositionAmt != "0" {
			fmt.Printf("Position %s: %s @ %s (P&L: %s)\n",
				position.Symbol, position.PositionAmt, position.EntryPrice, position.UnrealizedProfit)
		}
	}
}

func ExampleClient_PlaceOrder() {
	config := TestnetConfig()
	// Set API credentials for signed endpoints
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Place a limit buy order
	orderReq := &NewOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		Type:        OrderTypeLimit,
		TimeInForce: TimeInForceGTC,
		Quantity:    "0.001",
		Price:       "50000",
	}

	order, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}

	fmt.Printf("Order placed: ID=%d, Status=%s, Price=%s, Quantity=%s\n",
		order.OrderId, order.Status, order.Price, order.OrigQty)
}

func ExampleClient_GetOpenOrders() {
	config := TestnetConfig()
	// Set API credentials for signed endpoints
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all open orders
	orders, err := client.GetOpenOrders(ctx, "")
	if err != nil {
		log.Fatalf("Failed to get open orders: %v", err)
	}

	fmt.Printf("Found %d open orders\n", len(orders))

	for _, order := range orders {
		fmt.Printf("Order %d: %s %s %s @ %s (Status: %s)\n",
			order.OrderId, order.Symbol, order.Side, order.OrigQty, order.Price, order.Status)
	}
}

func ExampleClient_GetUserTrades() {
	config := TestnetConfig()
	// Set API credentials for signed endpoints
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get recent user trades
	trades, err := client.GetUserTrades(ctx, "BTCUSDT", 5)
	if err != nil {
		log.Fatalf("Failed to get user trades: %v", err)
	}

	fmt.Printf("Found %d recent trades for BTCUSDT\n", len(trades))

	for _, trade := range trades {
		fmt.Printf("Trade %d: %s %s @ %s (P&L: %s)\n",
			trade.Id, trade.Side, trade.Qty, trade.Price, trade.RealizedPnl)
	}
}

func ExampleClient_GetPositionRisk() {
	config := TestnetConfig()
	// Set API credentials for signed endpoints
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get position risk for all symbols
	positionRisks, err := client.GetPositionRisk(ctx, "")
	if err != nil {
		log.Fatalf("Failed to get position risk: %v", err)
	}

	fmt.Printf("Found %d position risk entries\n", len(positionRisks))

	// Print positions with non-zero amounts
	for _, pr := range positionRisks {
		if pr.HasPosition() {
			side := "LONG"
			if pr.IsShort() {
				side = "SHORT"
			}
			fmt.Printf("Position %s (%s): %s @ %s (P&L: %s, Liquidation: %s)\n",
				pr.Symbol, side, pr.PositionAmt, pr.EntryPrice, pr.UnrealizedPnl, pr.LiquidationPrice)
		}
	}
}

func ExampleClient_GetPositionSide() {
	config := TestnetConfig()
	// Set API credentials for signed endpoints
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current position side mode
	positionSide, err := client.GetPositionSide(ctx)
	if err != nil {
		log.Fatalf("Failed to get position side: %v", err)
	}

	if positionSide.DualSidePosition {
		fmt.Println("Position side mode: Dual Side (can hold both long and short positions)")
	} else {
		fmt.Println("Position side mode: Single Side (can only hold one position side)")
	}
}

func ExampleClient_ChangePositionSide() {
	config := TestnetConfig()
	// Set API credentials for signed endpoints
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Change to dual side position mode
	positionSide, err := client.ChangePositionSide(ctx, true)
	if err != nil {
		log.Fatalf("Failed to change position side: %v", err)
	}

	fmt.Printf("Position side changed to: dualSidePosition=%t\n", positionSide.DualSidePosition)
}

func ExampleClient_GetLeverage() {
	config := TestnetConfig()
	// Set API credentials for signed endpoints
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current leverage for BTCUSDT
	leverage, err := client.GetLeverage(ctx, "BTCUSDT")
	if err != nil {
		log.Fatalf("Failed to get leverage: %v", err)
	}

	fmt.Printf("Current leverage for %s: %dx (Max notional: %s)\n",
		leverage.Symbol, leverage.Leverage, leverage.MaxNotionalValue)
}

func ExampleClient_ChangeLeverage() {
	config := TestnetConfig()
	// Set API credentials for signed endpoints
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Change leverage to 10x for BTCUSDT
	leverage, err := client.ChangeLeverage(ctx, "BTCUSDT", 10)
	if err != nil {
		log.Fatalf("Failed to change leverage: %v", err)
	}

	fmt.Printf("Leverage changed for %s: %dx (Max notional: %s)\n",
		leverage.Symbol, leverage.Leverage, leverage.MaxNotionalValue)
}
