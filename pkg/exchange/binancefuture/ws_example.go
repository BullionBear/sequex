package binancefuture

import (
	"fmt"
	"log"
	"time"
)

// ExampleWSStreamClient_SubscribeToAggTrade demonstrates how to subscribe to aggregated trades
func ExampleWSStreamClient_SubscribeToAggTrade() {
	// Create a configuration
	config := &Config{
		UseTestnet: true, // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Subscribe to aggregated trades for BTCUSDT
	symbol := "btcusdt"
	unsubscribe, err := wsClient.SubscribeToAggTrade(symbol, func(data []byte) error {
		// Parse the aggregated trade data
		aggTrade, err := ParseAggTradeData(data)
		if err != nil {
			log.Printf("Failed to parse aggregated trade data: %v", err)
			return err
		}

		// Print the trade information
		fmt.Printf("Aggregated Trade: Symbol=%s, Price=%.2f, Quantity=%.4f, IsBuyerMaker=%t\n",
			aggTrade.Symbol, aggTrade.Price, aggTrade.Quantity, aggTrade.IsBuyerMaker)

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to subscribe to aggregated trades: %v", err)
	}

	// Keep the connection alive for a few seconds to receive some data
	time.Sleep(10 * time.Second)

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		log.Printf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		log.Printf("Failed to close WebSocket client: %v", err)
	}
}

// ExampleWSStreamClient_SubscribeToAggTradeWithCallback demonstrates how to subscribe to aggregated trades with type-specific callback
func ExampleWSStreamClient_SubscribeToAggTradeWithCallback() {
	// Create a configuration
	config := &Config{
		UseTestnet: true, // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Subscribe to aggregated trades with type-specific callback
	symbol := "btcusdt"
	unsubscribe, err := wsClient.SubscribeToAggTradeWithCallback(symbol, func(data *WSAggTradeData) error {
		// Print the trade information
		fmt.Printf("Aggregated Trade: Symbol=%s, Price=%.2f, Quantity=%.4f, IsBuyerMaker=%t\n",
			data.Symbol, data.Price, data.Quantity, data.IsBuyerMaker)

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to subscribe to aggregated trades: %v", err)
	}

	// Keep the connection alive for a few seconds to receive some data
	time.Sleep(10 * time.Second)

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		log.Printf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		log.Printf("Failed to close WebSocket client: %v", err)
	}
}

// ExampleWSStreamClient_SubscribeToCombinedStreams demonstrates how to subscribe to multiple streams
func ExampleWSStreamClient_SubscribeToCombinedStreams() {
	// Create a configuration
	config := &Config{
		UseTestnet: true, // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Subscribe to multiple streams
	streams := []string{
		"btcusdt@aggTrade",
		"ethusdt@aggTrade",
		"bnbusdt@aggTrade",
	}

	unsubscribe, err := wsClient.SubscribeToCombinedStreams(streams, func(data []byte) error {
		// Parse the aggregated trade data
		aggTrade, err := ParseAggTradeData(data)
		if err != nil {
			log.Printf("Failed to parse aggregated trade data: %v", err)
			return err
		}

		// Print the trade information
		fmt.Printf("Combined Stream Trade: Symbol=%s, Price=%.2f, Quantity=%.4f\n",
			aggTrade.Symbol, aggTrade.Price, aggTrade.Quantity)

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to subscribe to combined streams: %v", err)
	}

	// Keep the connection alive for a few seconds to receive some data
	time.Sleep(10 * time.Second)

	// Unsubscribe
	if err := unsubscribe(); err != nil {
		log.Printf("Failed to unsubscribe: %v", err)
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		log.Printf("Failed to close WebSocket client: %v", err)
	}
}

// ExampleWSStreamClient_MultipleSubscriptions demonstrates how to manage multiple subscriptions
func ExampleWSStreamClient_MultipleSubscriptions() {
	// Create a configuration
	config := &Config{
		UseTestnet: true, // Use testnet for testing
	}

	// Create a new WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Subscribe to multiple different types of streams
	unsubscribes := make([]func() error, 0)

	// Subscribe to aggregated trades
	aggTradeUnsub, err := wsClient.SubscribeToAggTrade("btcusdt", func(data []byte) error {
		aggTrade, _ := ParseAggTradeData(data)
		fmt.Printf("AggTrade: %s %.2f\n", aggTrade.Symbol, aggTrade.Price)
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to aggregated trades: %v", err)
	}
	unsubscribes = append(unsubscribes, aggTradeUnsub)

	// Subscribe to ticker
	tickerUnsub, err := wsClient.SubscribeToTicker("btcusdt", func(data []byte) error {
		ticker, _ := ParseTickerData(data)
		fmt.Printf("Ticker: %s %.2f\n", ticker.Symbol, ticker.LastPrice)
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to ticker: %v", err)
	}
	unsubscribes = append(unsubscribes, tickerUnsub)

	// Subscribe to book ticker
	bookTickerUnsub, err := wsClient.SubscribeToBookTicker("btcusdt", func(data []byte) error {
		bookTicker, _ := ParseBookTickerData(data)
		fmt.Printf("BookTicker: %s Bid=%.2f Ask=%.2f\n",
			bookTicker.Symbol, bookTicker.BidPrice, bookTicker.AskPrice)
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to book ticker: %v", err)
	}
	unsubscribes = append(unsubscribes, bookTickerUnsub)

	// Print active streams
	fmt.Printf("Active streams: %v\n", wsClient.GetActiveStreams())

	// Keep the connection alive for a few seconds to receive some data
	time.Sleep(10 * time.Second)

	// Unsubscribe from all streams
	for _, unsub := range unsubscribes {
		if err := unsub(); err != nil {
			log.Printf("Failed to unsubscribe: %v", err)
		}
	}

	// Close the client
	if err := wsClient.Close(); err != nil {
		log.Printf("Failed to close WebSocket client: %v", err)
	}
}
