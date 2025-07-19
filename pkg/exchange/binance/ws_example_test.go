package binance

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Example demonstrates how to use the WebSocket client
// to subscribe to kline (candlestick) data.
func Example() {
	// Create WebSocket client with default production config
	// Note: WebSocket streams are public and don't require API credentials
	config := DefaultConfig()
	client := NewWSClient(config)

	// Track received events
	var eventCount int
	var mu sync.Mutex

	// Set up kline event handler
	client.OnKline(func(event *WSKlineEvent) {
		mu.Lock()
		defer mu.Unlock()
		eventCount++

		fmt.Printf("Kline Event #%d:\n", eventCount)
		fmt.Printf("  Symbol: %s\n", event.Symbol)
		fmt.Printf("  Interval: %s\n", event.Kline.Interval)
		fmt.Printf("  Open: %s\n", event.Kline.OpenPrice)
		fmt.Printf("  Close: %s\n", event.Kline.ClosePrice)
		fmt.Printf("  High: %s\n", event.Kline.HighPrice)
		fmt.Printf("  Low: %s\n", event.Kline.LowPrice)
		fmt.Printf("  Volume: %s\n", event.Kline.BaseAssetVolume)
		fmt.Printf("  Closed: %t\n", event.Kline.IsClosed)
		fmt.Printf("  Time: %s\n", event.Kline.GetOpenTime().Format("15:04:05"))
		fmt.Println("---")
	})

	// Connect to WebSocket
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	fmt.Println("Connected to Binance WebSocket")

	// Subscribe to BTCUSDT 1-minute klines
	err = client.SubscribeKline([]string{"BTCUSDT"}, Interval1m)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	fmt.Println("Subscribed to BTCUSDT 1m klines")

	// Listen for events for 5 seconds
	time.Sleep(5 * time.Second)

	// Check active subscriptions
	subscriptions := client.GetSubscriptions()
	fmt.Printf("Active subscriptions: %v\n", subscriptions)

	// Unsubscribe
	err = client.UnsubscribeKline([]string{"BTCUSDT"}, Interval1m)
	if err != nil {
		log.Printf("Failed to unsubscribe: %v", err)
	}

	// Disconnect
	err = client.Disconnect()
	if err != nil {
		log.Printf("Failed to disconnect: %v", err)
	}

	mu.Lock()
	totalEvents := eventCount
	mu.Unlock()

	fmt.Printf("Received %d kline events\n", totalEvents)
	fmt.Println("WebSocket example completed")
}

// Example_multipleStreams demonstrates subscribing to multiple data streams
func Example_multipleStreams() {
	config := DefaultConfig()
	client := NewWSClient(config)

	// Set up handlers for different event types
	client.OnKline(func(event *WSKlineEvent) {
		fmt.Printf("📊 Kline: %s %s\n", event.Symbol, event.Kline.ClosePrice)
	})

	client.OnTicker(func(event *WSTickerEvent) {
		fmt.Printf("📈 Ticker: %s = %s (24h change: %s%%)\n",
			event.Symbol, event.LastPrice, event.PriceChangePercent)
	})

	client.OnTrade(func(event *WSTradeEvent) {
		fmt.Printf("💰 Trade: %s %s @ %s\n",
			event.Symbol, event.Quantity, event.Price)
	})

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Subscribe to multiple streams
	symbols := []string{"BTCUSDT", "ETHUSDT"}

	err = client.SubscribeKline(symbols, Interval1m)
	if err != nil {
		log.Fatalf("Failed to subscribe to klines: %v", err)
	}

	err = client.SubscribeTicker(symbols)
	if err != nil {
		log.Fatalf("Failed to subscribe to tickers: %v", err)
	}

	err = client.SubscribeTrade(symbols)
	if err != nil {
		log.Fatalf("Failed to subscribe to trades: %v", err)
	}

	fmt.Println("Subscribed to multiple streams, listening...")

	// Listen for events
	time.Sleep(10 * time.Second)

	// Show active subscriptions
	subscriptions := client.GetSubscriptions()
	fmt.Printf("Active subscriptions: %d\n", len(subscriptions))
	for _, sub := range subscriptions {
		fmt.Printf("  - %s\n", sub)
	}

	// Clean up
	client.Disconnect()
	fmt.Println("Multi-stream example completed")
}
