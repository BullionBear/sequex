package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func wsExample() {
	// Create configuration
	config := binance.DefaultConfig()

	// Create WebSocket stream client
	client := binance.NewWSStreamClient(config)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Starting Binance WebSocket example...")
	fmt.Println("Press Ctrl+C to stop")

	// Subscribe to BTCUSDT ticker
	tickerCallback := func(data []byte) error {
		tickerData, err := binance.ParseTickerData(data)
		if err != nil {
			log.Printf("Error parsing ticker data: %v", err)
			return err
		}

		fmt.Printf("Ticker - Symbol: %s, Price: %.2f, Change: %.2f%%\n",
			tickerData.Symbol, tickerData.LastPrice, tickerData.PriceChangePercent)
		return nil
	}

	unsubscribeTicker, err := client.SubscribeToTicker("BTCUSDT", tickerCallback)
	if err != nil {
		log.Fatalf("Failed to subscribe to ticker: %v", err)
	}
	defer unsubscribeTicker()

	// Subscribe to BTCUSDT kline data
	klineCallback := func(data []byte) error {
		klineData, err := binance.ParseKlineData(data)
		if err != nil {
			log.Printf("Error parsing kline data: %v", err)
			return err
		}

		if klineData.IsKlineClosed {
			fmt.Printf("Kline - Symbol: %s, Interval: %s, Open: %.2f, Close: %.2f, High: %.2f, Low: %.2f\n",
				klineData.Symbol, klineData.Kline.Interval,
				klineData.Kline.OpenPrice, klineData.Kline.ClosePrice,
				klineData.Kline.HighPrice, klineData.Kline.LowPrice)
		}
		return nil
	}

	unsubscribeKline, err := client.SubscribeToKline("BTCUSDT", "1m", klineCallback)
	if err != nil {
		log.Fatalf("Failed to subscribe to kline: %v", err)
	}
	defer unsubscribeKline()

	// Subscribe to BTCUSDT trade data
	tradeCallback := func(data []byte) error {
		tradeData, err := binance.ParseTradeData(data)
		if err != nil {
			log.Printf("Error parsing trade data: %v", err)
			return err
		}

		fmt.Printf("Trade - Symbol: %s, Price: %.2f, Quantity: %.4f, Time: %s\n",
			tradeData.Symbol, tradeData.Price, tradeData.Quantity,
			time.Unix(0, tradeData.TradeTime*int64(time.Millisecond)).Format("15:04:05"))
		return nil
	}

	unsubscribeTrade, err := client.SubscribeToTrade("BTCUSDT", tradeCallback)
	if err != nil {
		log.Fatalf("Failed to subscribe to trade: %v", err)
	}
	defer unsubscribeTrade()

	// Subscribe to BTCUSDT book ticker
	bookTickerCallback := func(data []byte) error {
		bookTickerData, err := binance.ParseBookTickerData(data)
		if err != nil {
			log.Printf("Error parsing book ticker data: %v", err)
			return err
		}

		fmt.Printf("Book Ticker - Symbol: %s, Bid: %.2f (%.4f), Ask: %.2f (%.4f)\n",
			bookTickerData.Symbol, bookTickerData.BidPrice, bookTickerData.BidQty,
			bookTickerData.AskPrice, bookTickerData.AskQty)
		return nil
	}

	unsubscribeBookTicker, err := client.SubscribeToBookTicker("BTCUSDT", bookTickerCallback)
	if err != nil {
		log.Fatalf("Failed to subscribe to book ticker: %v", err)
	}
	defer unsubscribeBookTicker()

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nShutting down...")

	// Close all connections
	if err := client.Close(); err != nil {
		log.Printf("Error closing client: %v", err)
	}

	fmt.Println("WebSocket example completed")
}

// Example function showing how to use the WebSocket client programmatically
func ExampleUsage() {
	config := binance.DefaultConfig()
	client := binance.NewWSStreamClient(config)

	// Subscribe to multiple streams
	callback := func(data []byte) error {
		// Try to parse as different types
		if ticker, err := binance.ParseTickerData(data); err == nil {
			fmt.Printf("Received ticker: %s = %.2f\n", ticker.Symbol, ticker.LastPrice)
			return nil
		}

		if kline, err := binance.ParseKlineData(data); err == nil {
			fmt.Printf("Received kline: %s %s\n", kline.Symbol, kline.Kline.Interval)
			return nil
		}

		if trade, err := binance.ParseTradeData(data); err == nil {
			fmt.Printf("Received trade: %s %.2f\n", trade.Symbol, trade.Price)
			return nil
		}

		// If none of the above, just print raw data
		fmt.Printf("Received raw data: %s\n", string(data))
		return nil
	}

	// Subscribe to multiple streams
	streams := []string{
		"btcusdt@ticker",
		"btcusdt@kline_1m",
		"btcusdt@trade",
	}

	unsubscribe, err := client.SubscribeToCombinedStreams(streams, callback)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Do something with the data...
	time.Sleep(30 * time.Second)

	// Unsubscribe when done
	unsubscribe()
}

func main() {
	wsExample()
}
