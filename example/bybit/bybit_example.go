package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/bybit"
)

func main() {
	// Create Bybit client with testnet configuration
	config := bybit.DefaultConfig()

	config = config.WithAPIKey(os.Getenv("BYBIT_API_KEY")).WithAPISecret(os.Getenv("BYBIT_API_SECRET"))
	client := bybit.NewClient(config)

	ctx := context.Background()

	// Example 1: Get server time
	fmt.Println("=== Server Time ===")
	serverTime, err := client.GetServerTime(ctx)
	if err != nil {
		log.Fatalf("Failed to get server time: %v", err)
	}
	fmt.Printf("Server time: %s\n", serverTime.Result.TimeSecond)

	// Example 2: Get kline data for BTCUSD inverse perpetual
	fmt.Println("\n=== Kline Data ===")
	klineReq := &bybit.KlineRequest{
		Category: "inverse",
		Symbol:   "BTCUSD",
		Interval: "60",          // 1 hour interval
		Start:    1670601600000, // Start time
		End:      1670608800000, // End time
	}

	klineResp, err := client.GetKline(ctx, klineReq)
	if err != nil {
		log.Fatalf("Failed to get kline data: %v", err)
	}

	fmt.Printf("Retrieved %d kline records for %s (%s)\n",
		len(klineResp.Result.List),
		klineResp.Result.Symbol,
		klineResp.Result.Category)

	// Example 3: Get parsed kline data
	fmt.Println("\n=== Parsed Kline Data ===")
	klineData, err := client.GetKlineData(ctx, klineReq)
	if err != nil {
		log.Fatalf("Failed to get parsed kline data: %v", err)
	}

	for i, kline := range klineData {
		fmt.Printf("Kline %d: Time=%s, Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.2f\n",
			i+1,
			kline.Timestamp.Format(time.RFC3339),
			kline.OpenPrice,
			kline.HighPrice,
			kline.LowPrice,
			kline.ClosePrice,
			kline.Volume)
	}

	// Example 4: Get ticker information
	fmt.Println("\n=== Ticker Information ===")
	tickerResp, err := client.GetTickers(ctx, "spot", "ADAUSDT")
	if err != nil {
		log.Fatalf("Failed to get tickers: %v", err)
	}

	if len(tickerResp.Result.List) > 0 {
		ticker := tickerResp.Result.List[0]
		fmt.Printf("Symbol: %s\n", ticker.Symbol)
		fmt.Printf("Last Price: %s\n", ticker.LastPrice)
		fmt.Printf("24h High: %s\n", ticker.HighPrice24h)
		fmt.Printf("24h Low: %s\n", ticker.LowPrice24h)
		fmt.Printf("24h Volume: %s\n", ticker.Volume24h)
		fmt.Printf("24h Turnover: %s\n", ticker.Turnover24h)
	}

	// Example 5: Get recent kline data (last 10 records)
	fmt.Println("\n=== Recent Kline Data ===")
	recentReq := &bybit.KlineRequest{
		Category: "inverse",
		Symbol:   "BTCUSD",
		Interval: "15", // 15 minutes
		Limit:    10,   // Last 10 records
	}

	recentKlineData, err := client.GetKlineData(ctx, recentReq)
	if err != nil {
		log.Fatalf("Failed to get recent kline data: %v", err)
	}

	fmt.Printf("Recent %d kline records:\n", len(recentKlineData))
	for i, kline := range recentKlineData {
		fmt.Printf("  %d. %s - Close: %.2f, Volume: %.2f\n",
			i+1,
			kline.Timestamp.Format("15:04"),
			kline.ClosePrice,
			kline.Volume)
	}

	// Example 6: Get account information
	fmt.Println("\n=== Account Information ===")
	accountReq := &bybit.GetAccountRequest{
		AccountType: bybit.AccountTypeUnified,
	}
	accountResp, err := client.GetAccount(ctx, accountReq)
	if err != nil {
		log.Printf("Failed to get account info: %v", err)
	} else {
		for _, account := range accountResp.Result.List {
			fmt.Printf("Account Type: %s\n", account.AccountType)
			fmt.Printf("Total Wallet Balance: %s\n", account.TotalWalletBalance)
			fmt.Printf("Total Available Balance: %s\n", account.TotalAvailableBalance)
			fmt.Printf("Total Equity: %s\n", account.TotalEquity)
			fmt.Printf("Total Perp UPL: %s\n", account.TotalPerpUPL)
			fmt.Printf("Account LTV: %s\n", account.AccountLTV)
			fmt.Printf("Account IM Rate: %s\n", account.AccountIMRate)
			fmt.Printf("Account MM Rate: %s\n", account.AccountMMRate)

			fmt.Println("Coin Details:")
			for _, coin := range account.Coin {
				fmt.Printf("  Coin: %s\n", coin.Coin)
				fmt.Printf("  Wallet Balance: %s\n", coin.WalletBalance)
				fmt.Printf("  Available to Withdraw: %s\n", coin.AvailableToWithdraw)
				fmt.Printf("  Available to Borrow: %s\n", coin.AvailableToBorrow)
				fmt.Printf("  Unrealised PnL: %s\n", coin.UnrealisedPnl)
				fmt.Printf("  Cum Realised PnL: %s\n", coin.CumRealisedPnl)
				fmt.Printf("  Locked: %s\n", coin.Locked)
				fmt.Printf("  USD Value: %s\n", coin.UsdValue)
				fmt.Printf("  Margin Collateral: %t\n", coin.MarginCollateral)
				fmt.Printf("  Collateral Switch: %t\n", coin.CollateralSwitch)
			}
		}
	}

	// Example 7: Create a limit buy order (100 USD contract value)
	fmt.Println("\n=== Creating Limit Buy Order ===")
	createOrderReq := &bybit.CreateOrderRequest{
		Category:    "spot",
		Symbol:      "ADAUSDT",
		Side:        "Buy",
		OrderType:   "Limit",
		Qty:         "10",
		Price:       "0.8", // Set price below current market price
		TimeInForce: "GTC",
	}

	createOrderResp, err := client.CreateOrder(ctx, createOrderReq)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
	} else {
		fmt.Printf("Order created successfully!\n")
		fmt.Printf("Order ID: %s\n", createOrderResp.Result.OrderId)
		fmt.Printf("Order Link ID: %s\n", createOrderResp.Result.OrderLinkId)

		// Example 8: Get order information (UTA 2.0)
		fmt.Println("\n=== Order Information ===")
		getOrderReq := &bybit.GetOrderRequest{
			Category: "spot",
			Symbol:   "ADAUSDT",
			OrderId:  createOrderResp.Result.OrderId,
		}

		getOrderResp, err := client.GetOrder(ctx, getOrderReq)
		if err != nil {
			log.Printf("Failed to get order: %v", err)
		} else {
			if len(getOrderResp.Result.List) > 0 {
				order := getOrderResp.Result.List[0]
				fmt.Printf("Order ID: %s\n", order.OrderId)
				fmt.Printf("Status: %s\n", order.OrderStatus)
				fmt.Printf("Created Time: %s\n", order.CreatedTime)
				fmt.Printf("Updated Time: %s\n", order.UpdatedTime)
				fmt.Printf("Executed Qty: %s\n", order.CumExecQty)
				fmt.Printf("Executed Value: %s\n", order.CumExecValue)
			} else {
				fmt.Printf("No orders found\n")
			}
		}

		// Example 8b: Get single order information
		fmt.Println("\n=== Single Order Information ===")
		getSingleOrderReq := &bybit.GetOrderRequest{
			Category: "spot",
			Symbol:   "ADAUSDT",
			OrderId:  createOrderResp.Result.OrderId,
		}

		getSingleOrderResp, err := client.GetSingleOrder(ctx, getSingleOrderReq)
		if err != nil {
			log.Printf("Failed to get single order: %v", err)
		} else {
			fmt.Printf("Order ID: %s\n", getSingleOrderResp.Result.OrderId)
			fmt.Printf("Status: %s\n", getSingleOrderResp.Result.OrderStatus)
			fmt.Printf("Created Time: %s\n", getSingleOrderResp.Result.CreatedTime)
			fmt.Printf("Updated Time: %s\n", getSingleOrderResp.Result.UpdatedTime)
			fmt.Printf("Executed Qty: %s\n", getSingleOrderResp.Result.CumExecQty)
			fmt.Printf("Executed Value: %s\n", getSingleOrderResp.Result.CumExecValue)
		}

		// Example 9: Cancel the order
		fmt.Println("\n=== Canceling Order ===")
		cancelOrderReq := &bybit.CancelOrderRequest{
			Category: "spot",
			Symbol:   "ADAUSDT",
			OrderId:  createOrderResp.Result.OrderId,
		}

		cancelOrderResp, err := client.CancelOrder(ctx, cancelOrderReq)
		if err != nil {
			log.Printf("Failed to cancel order: %v", err)
		} else {
			fmt.Printf("Order canceled successfully!\n")
			fmt.Printf("Order ID: %s\n", cancelOrderResp.Result.OrderId)
			fmt.Printf("Status: %s\n", cancelOrderResp.Result.OrderStatus)
		}
	}

	// Example 10: Create a market sell order (100 USD contract value)
	fmt.Println("\n=== Creating Market Sell Order ===")
	marketSellReq := &bybit.CreateOrderRequest{
		Category:  "spot",
		Symbol:    "ADAUSDT",
		Side:      "Sell",
		OrderType: "Market",
		Qty:       "100", // 100 USD contract value
	}

	marketSellResp, err := client.CreateOrder(ctx, marketSellReq)
	if err != nil {
		log.Printf("Failed to create market sell order: %v", err)
	} else {
		fmt.Printf("Market sell order created successfully!\n")
		fmt.Printf("Order ID: %s\n", marketSellResp.Result.OrderId)
		fmt.Printf("Order Link ID: %s\n", marketSellResp.Result.OrderLinkId)
	}
}
