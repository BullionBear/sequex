package binance

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ExampleUsage demonstrates how to use the new WebSocket subscription pattern
func ExampleUsage() {
	// Create configuration
	config := &Config{
		APIKey:     "your-api-key",
		APISecret:  "your-api-secret",
		UseTestnet: true,
	}

	// Create WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Example 1: Subscribe to kline data with chainable options
	fmt.Println("=== Example 1: Kline Subscription ===")

	// Create kline subscription options
	klineOptions := &KlineSubscriptionOptions{}
	klineOptions.WithConnect(func() {
		fmt.Println("Connected to kline stream")
	}).WithKline(func(data *WSKlineData) error {
		fmt.Printf("Kline: Symbol=%s, Interval=%s, Close=%f, Volume=%f\n",
			data.Symbol, data.Kline.Interval, data.Kline.ClosePrice, data.Kline.Volume)
		return nil
	}).WithError(func(err error) {
		fmt.Printf("Kline error: %v\n", err)
	}).WithDisconnect(func() {
		fmt.Println("Disconnected from kline stream")
	})

	// Subscribe to kline data
	unsubscribeKline, err := wsClient.SubscribeToKline("BTCUSDT", "1m", klineOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to kline: %v", err)
	}

	// Example 2: Subscribe to ticker data with chainable options
	fmt.Println("\n=== Example 2: Ticker Subscription ===")

	tickerOptions := &TickerSubscriptionOptions{}
	tickerOptions.WithConnect(func() {
		fmt.Println("Connected to ticker stream")
	}).WithTicker(func(data *WSTickerData) error {
		fmt.Printf("Ticker: Symbol=%s, LastPrice=%f, Volume=%f, ChangePercent=%f%%\n",
			data.Symbol, data.LastPrice, data.Volume, data.PriceChangePercent)
		return nil
	}).WithError(func(err error) {
		fmt.Printf("Ticker error: %v\n", err)
	})

	unsubscribeTicker, err := wsClient.SubscribeToTicker("ETHUSDT", tickerOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to ticker: %v", err)
	}

	// Example 3: Subscribe to user data stream (requires API key)
	fmt.Println("\n=== Example 3: User Data Stream Subscription ===")

	// Create listen key (this requires API credentials)
	restClient := NewClient(config)
	userDataStream, err := restClient.CreateUserDataStream(context.Background())
	if err != nil {
		fmt.Printf("Failed to create user data stream: %v\n", err)
	} else {
		userDataOptions := &UserDataSubscriptionOptions{}
		userDataOptions.WithConnect(func() {
			fmt.Println("Connected to user data stream")
		}).WithExecutionReport(func(data *WSExecutionReport) error {
			fmt.Printf("Execution Report: Symbol=%s, Status=%s, Side=%s, Price=%s\n",
				data.Symbol, data.CurrentOrderStatus, data.Side, data.OrderPrice)
			return nil
		}).WithAccountUpdate(func(data *WSOutboundAccountPosition) error {
			fmt.Printf("Account Update: EventTime=%d, Balances=%d\n",
				data.EventTime, len(data.Balances))
			return nil
		}).WithBalanceUpdate(func(data *WSBalanceUpdate) error {
			fmt.Printf("Balance Update: Asset=%s, Delta=%s\n",
				data.Asset, data.BalanceDelta)
			return nil
		}).WithError(func(err error) {
			fmt.Printf("User data error: %v\n", err)
		})

		unsubscribeUserData, err := wsClient.SubscribeToUserDataStream(userDataStream.ListenKey, userDataOptions)
		if err != nil {
			fmt.Printf("Failed to subscribe to user data stream: %v\n", err)
		} else {
			// Keep the user data stream alive
			go func() {
				ticker := time.NewTicker(30 * time.Minute)
				defer ticker.Stop()
				for range ticker.C {
					if err := restClient.KeepAliveUserDataStream(context.Background(), userDataStream.ListenKey); err != nil {
						fmt.Printf("Failed to keep alive user data stream: %v\n", err)
					}
				}
			}()

			// Cleanup function for user data stream
			defer func() {
				unsubscribeUserData()
				restClient.CloseUserDataStream(context.Background(), userDataStream.ListenKey)
			}()
		}
	}

	// Example 5: Subscribe to all mini tickers
	fmt.Println("\n=== Example 5: All Mini Tickers ===")

	allMiniTickerOptions := &MiniTickerSubscriptionOptions{}
	allMiniTickerOptions.WithConnect(func() {
		fmt.Println("Connected to all mini tickers stream")
	}).WithMiniTicker(func(data *WSMiniTickerData) error {
		fmt.Printf("Mini Ticker: Symbol=%s, Close=%f, Volume=%f\n",
			data.Symbol, data.ClosePrice, data.Volume)
		return nil
	}).WithError(func(err error) {
		fmt.Printf("Mini ticker error: %v\n", err)
	})

	unsubscribeAllMiniTickers, err := wsClient.SubscribeToAllMiniTickers(allMiniTickerOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to all mini tickers: %v", err)
	}

	// Run for a while to see the data
	fmt.Println("\n=== Running for 30 seconds to see data ===")
	time.Sleep(30 * time.Second)

	// Cleanup
	fmt.Println("\n=== Cleaning up subscriptions ===")

	unsubscribeKline()
	unsubscribeTicker()
	unsubscribeAllMiniTickers()

	// Close all connections
	wsClient.Close()

	fmt.Println("All subscriptions cleaned up")
}

// ExampleMinimal demonstrates minimal usage with just the essential callbacks
func ExampleMinimal() {
	config := &Config{
		UseTestnet: true,
	}

	wsClient := NewWSStreamClient(config)

	// Minimal kline subscription - just the data callback
	klineOptions := &KlineSubscriptionOptions{}
	klineOptions.WithKline(func(data *WSKlineData) error {
		fmt.Printf("Kline: %s %s Close: %f\n",
			data.Symbol, data.Kline.Interval, data.Kline.ClosePrice)
		return nil
	})

	unsubscribe, err := wsClient.SubscribeToKline("BTCUSDT", "1m", klineOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Run for 10 seconds
	time.Sleep(10 * time.Second)

	unsubscribe()
	wsClient.Close()
}
