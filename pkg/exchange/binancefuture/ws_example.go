package binancefuture

import (
	"fmt"
	"log"
	"time"
)

// Example usage of the refactored WSStreamClient with subscription options

func ExampleWSStreamClientUsage() {
	// Create configuration
	config := &Config{
		BaseURL: "https://testnet.binancefuture.com",
		// Add your API key and secret if needed for user data streams
		// APIKey:    "your-api-key",
		// APISecret: "your-api-secret",
	}

	// Create WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Example 1: Subscribe to kline data with options
	fmt.Println("=== Example 1: Kline Subscription ===")
	klineOptions := NewKlineSubscriptionOptions()
	klineOptions.WithConnect(func() {
		fmt.Println("Connected to kline stream")
	}).WithKline(func(data *WSKlineData) error {
		fmt.Printf("Kline: Symbol=%s, Close=%f, Volume=%f\n",
			data.Symbol, data.Kline.ClosePrice, data.Kline.Volume)
		return nil
	}).WithError(func(err error) {
		fmt.Printf("Kline Error: %v\n", err)
	}).WithDisconnect(func() {
		fmt.Println("Disconnected from kline stream")
	})

	unsubscribeKline, err := wsClient.SubscribeToKline("BTCUSDT", "1m", klineOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to kline: %v", err)
	}

	// Example 2: Subscribe to ticker data with options
	fmt.Println("\n=== Example 2: Ticker Subscription ===")
	tickerOptions := NewTickerSubscriptionOptions()
	tickerOptions.WithConnect(func() {
		fmt.Println("Connected to ticker stream")
	}).WithTicker(func(data *WSTickerData) error {
		fmt.Printf("Ticker: Symbol=%s, LastPrice=%f, Volume=%f\n",
			data.Symbol, data.LastPrice, data.Volume)
		return nil
	}).WithError(func(err error) {
		fmt.Printf("Ticker Error: %v\n", err)
	})

	unsubscribeTicker, err := wsClient.SubscribeToTicker("ETHUSDT", tickerOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to ticker: %v", err)
	}

	// Example 3: Subscribe to trade data with options
	fmt.Println("\n=== Example 3: Trade Subscription ===")
	tradeOptions := NewTradeSubscriptionOptions()
	tradeOptions.WithConnect(func() {
		fmt.Println("Connected to trade stream")
	}).WithTrade(func(data *WSTradeData) error {
		fmt.Printf("Trade: Symbol=%s, Price=%f, Quantity=%f, IsBuyerMaker=%t\n",
			data.Symbol, data.Price, data.Quantity, data.IsBuyerMaker)
		return nil
	}).WithError(func(err error) {
		fmt.Printf("Trade Error: %v\n", err)
	})

	unsubscribeTrade, err := wsClient.SubscribeToTrade("BNBUSDT", tradeOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to trade: %v", err)
	}

	// Example 4: Subscribe to depth data with options
	fmt.Println("\n=== Example 4: Depth Subscription ===")
	depthOptions := NewDepthSubscriptionOptions()
	depthOptions.WithConnect(func() {
		fmt.Println("Connected to depth stream")
	}).WithDepth(func(data *WSDepthData) error {
		fmt.Printf("Depth: Symbol=%s, Bids=%d, Asks=%d\n",
			data.Symbol, len(data.Bids), len(data.Asks))
		return nil
	}).WithError(func(err error) {
		fmt.Printf("Depth Error: %v\n", err)
	})

	unsubscribeDepth, err := wsClient.SubscribeToDepth("ADAUSDT", "5", depthOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to depth: %v", err)
	}

	// Example 5: Subscribe to mark price data with options
	fmt.Println("\n=== Example 5: Mark Price Subscription ===")
	markPriceOptions := NewMarkPriceSubscriptionOptions()
	markPriceOptions.WithConnect(func() {
		fmt.Println("Connected to mark price stream")
	}).WithMarkPrice(func(data *WSMarkPriceData) error {
		fmt.Printf("Mark Price: Symbol=%s, MarkPrice=%f, FundingRate=%f\n",
			data.Symbol, data.MarkPrice, data.LastFundingRate)
		return nil
	}).WithError(func(err error) {
		fmt.Printf("Mark Price Error: %v\n", err)
	})

	unsubscribeMarkPrice, err := wsClient.SubscribeToMarkPrice("BTCUSDT", markPriceOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to mark price: %v", err)
	}

	// Example 6: Subscribe to funding rate data with options
	fmt.Println("\n=== Example 6: Funding Rate Subscription ===")
	fundingRateOptions := NewFundingRateSubscriptionOptions()
	fundingRateOptions.WithConnect(func() {
		fmt.Println("Connected to funding rate stream")
	}).WithFundingRate(func(data *WSFundingRateData) error {
		fmt.Printf("Funding Rate: Symbol=%s, Rate=%f, Time=%d\n",
			data.Symbol, data.FundingRate, data.FundingTime)
		return nil
	}).WithError(func(err error) {
		fmt.Printf("Funding Rate Error: %v\n", err)
	})

	unsubscribeFundingRate, err := wsClient.SubscribeToFundingRate("BTCUSDT", fundingRateOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to funding rate: %v", err)
	}

	// Example 7: Subscribe to user data stream (requires API credentials)
	fmt.Println("\n=== Example 7: User Data Stream Subscription ===")
	// Note: This requires API credentials to be set in the config
	/*
		userDataOptions := NewUserDataSubscriptionOptions()
		userDataOptions.WithConnect(func() {
			fmt.Println("Connected to user data stream")
		}).WithExecutionReport(func(data *WSExecutionReport) error {
			fmt.Printf("Execution Report: Symbol=%s, Status=%s, Side=%s\n",
				data.Symbol, data.CurrentOrderStatus, data.Side)
			return nil
		}).WithAccountUpdate(func(data *WSOutboundAccountPosition) error {
			fmt.Printf("Account Update: Balances=%d\n", len(data.Balances))
			return nil
		}).WithBalanceUpdate(func(data *WSBalanceUpdate) error {
			fmt.Printf("Balance Update: Asset=%s, Delta=%s\n",
				data.Asset, data.BalanceDelta)
			return nil
		}).WithError(func(err error) {
			fmt.Printf("User Data Error: %v\n", err)
		})

		unsubscribeUserData, err := wsClient.SubscribeToUserDataStream(userDataOptions)
		if err != nil {
			log.Fatalf("Failed to subscribe to user data stream: %v", err)
		}
	*/

	// Let the streams run for a while
	fmt.Println("\n=== Running streams for 30 seconds ===")
	time.Sleep(30 * time.Second)

	// Unsubscribe from all streams
	fmt.Println("\n=== Unsubscribing from streams ===")

	if err := unsubscribeKline(); err != nil {
		fmt.Printf("Error unsubscribing from kline: %v\n", err)
	}

	if err := unsubscribeTicker(); err != nil {
		fmt.Printf("Error unsubscribing from ticker: %v\n", err)
	}

	if err := unsubscribeTrade(); err != nil {
		fmt.Printf("Error unsubscribing from trade: %v\n", err)
	}

	if err := unsubscribeDepth(); err != nil {
		fmt.Printf("Error unsubscribing from depth: %v\n", err)
	}

	if err := unsubscribeMarkPrice(); err != nil {
		fmt.Printf("Error unsubscribing from mark price: %v\n", err)
	}

	if err := unsubscribeFundingRate(); err != nil {
		fmt.Printf("Error unsubscribing from funding rate: %v\n", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		fmt.Printf("Error closing client: %v\n", err)
	}

	fmt.Println("=== Example completed ===")
}

// Example of chaining subscription options
func ExampleSubscriptionOptionsChaining() {
	// Create options with chaining
	options := NewKlineSubscriptionOptions().
		WithConnect(func() {
			fmt.Println("Connected!")
		}).
		WithKline(func(data *WSKlineData) error {
			fmt.Printf("Kline: %s\n", data.Symbol)
			return nil
		}).
		WithError(func(err error) {
			fmt.Printf("Error: %v\n", err)
		}).
		WithDisconnect(func() {
			fmt.Println("Disconnected!")
		})

	// Use the options
	fmt.Printf("Options created with %d callbacks\n",
		len([]interface{}{
			options.connectCallback,
			options.klineCallback,
			options.errorCallback,
			options.disconnectCallback,
		}))
}
