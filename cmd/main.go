package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
	fmt.Println("Binance API Client Demo")
	fmt.Println("=======================")

	// Load configuration from file
	appConfig, err := binance.LoadConfig("config/config.example.yml")
	if err != nil {
		log.Printf("Warning: Could not load config file, using default config: %v", err)

		// Use default config for demo (you'll need to set your API credentials)
		config := binance.DefaultConfig()
		config.APIKey = "your_api_key_here"
		config.APISecret = "your_api_secret_here"
		config.Sandbox = true // Use testnet for demo

		demonstrateAPI(config)
		return
	}

	// Get Binance configuration
	config, err := appConfig.GetBinanceConfig()
	if err != nil {
		log.Fatalf("Failed to get Binance config: %v", err)
	}

	demonstrateAPI(config)
}

func demonstrateAPI(config *binance.Config) {
	// Create Binance client
	client, err := binance.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Binance client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test connectivity
	fmt.Println("\n1. Testing connectivity...")
	if err := client.Ping(ctx); err != nil {
		log.Printf("Ping failed: %v", err)
	} else {
		fmt.Println("✓ Connection successful")
	}

	// Get server time
	fmt.Println("\n2. Getting server time...")
	serverTime, err := client.GetServerTime(ctx)
	if err != nil {
		log.Printf("Failed to get server time: %v", err)
	} else {
		fmt.Printf("✓ Server time: %v\n", time.UnixMilli(serverTime))
	}

	// Get exchange info
	fmt.Println("\n3. Getting exchange info...")
	exchangeInfo, err := client.GetExchangeInfo(ctx)
	if err != nil {
		log.Printf("Failed to get exchange info: %v", err)
	} else {
		fmt.Printf("✓ Exchange info loaded, %d symbols available\n", len(exchangeInfo.Symbols))
	}

	// Get ticker for BTCUSDT
	fmt.Println("\n4. Getting 24hr ticker for BTCUSDT...")
	ticker, err := client.GetTicker24hr(ctx, "BTCUSDT")
	if err != nil {
		log.Printf("Failed to get ticker: %v", err)
	} else {
		fmt.Printf("✓ BTCUSDT Price: %s, 24h Change: %s%%\n",
			ticker.LastPrice.String(), ticker.PriceChangePercent.String())
	}

	// Get order book
	fmt.Println("\n5. Getting order book for BTCUSDT (limit 5)...")
	orderBook, err := client.GetOrderBook(ctx, "BTCUSDT", 5)
	if err != nil {
		log.Printf("Failed to get order book: %v", err)
	} else {
		fmt.Printf("✓ Order book loaded, %d bids, %d asks\n",
			len(orderBook.Bids), len(orderBook.Asks))
		if len(orderBook.Bids) > 0 {
			fmt.Printf("  Best bid: %s\n", orderBook.Bids[0][0])
		}
		if len(orderBook.Asks) > 0 {
			fmt.Printf("  Best ask: %s\n", orderBook.Asks[0][0])
		}
	}

	// Get recent trades
	fmt.Println("\n6. Getting recent trades for BTCUSDT (limit 5)...")
	trades, err := client.GetRecentTrades(ctx, "BTCUSDT", 5)
	if err != nil {
		log.Printf("Failed to get recent trades: %v", err)
	} else {
		fmt.Printf("✓ Retrieved %d recent trades\n", len(trades))
		for i, trade := range trades {
			if i < 3 { // Show first 3 trades
				fmt.Printf("  Trade %d: Price %s, Qty %s\n",
					trade.ID, trade.Price.String(), trade.Qty.String())
			}
		}
	}

	// Get klines (candlestick data)
	fmt.Println("\n7. Getting klines for BTCUSDT (1h interval, limit 5)...")
	klines, err := client.GetKlines(ctx, "BTCUSDT", "1h", 5, nil, nil)
	if err != nil {
		log.Printf("Failed to get klines: %v", err)
	} else {
		fmt.Printf("✓ Retrieved %d klines\n", len(klines))
		for i, kline := range klines {
			if i < 3 { // Show first 3 klines
				fmt.Printf("  Kline %d: Open %s, High %s, Low %s, Close %s\n",
					i+1, kline.Open.String(), kline.High.String(),
					kline.Low.String(), kline.Close.String())
			}
		}
	}

	// Account endpoints require valid API credentials
	if config.APIKey != "your_api_key_here" && config.APISecret != "your_api_secret_here" {
		fmt.Println("\n8. Testing authenticated endpoints...")

		// Get account info
		account, err := client.GetAccount(ctx)
		if err != nil {
			log.Printf("Failed to get account info: %v", err)
		} else {
			fmt.Printf("✓ Account info loaded, can trade: %v\n", account.CanTrade)
			fmt.Printf("  Account has %d balances\n", len(account.Balances))
		}

		// Get open orders
		openOrders, err := client.GetOpenOrders(ctx, "")
		if err != nil {
			log.Printf("Failed to get open orders: %v", err)
		} else {
			fmt.Printf("✓ Retrieved %d open orders\n", len(openOrders))
		}
	} else {
		fmt.Println("\n8. Skipping authenticated endpoints (API credentials not configured)")
	}

	fmt.Println("\n✅ Demo completed successfully!")
	fmt.Println("\nTo use authenticated endpoints, update your config file with valid API credentials.")
}
