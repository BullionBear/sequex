package binancefuture

import (
	"encoding/json"
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

func ExampleWSStreamClient_SubscribeToUserDataStream() {
	// Create configuration
	config := &Config{
		APIKey:     "your-api-key",
		APISecret:  "your-api-secret",
		UseTestnet: true, // Use testnet for testing
	}

	// Create WebSocket stream client
	wsClient := NewWSStreamClient(config)

	// Define callback function for user data stream events
	callback := func(data []byte) error {
		// Parse the event type first
		var baseEvent struct {
			EventType string `json:"e"`
		}
		if err := json.Unmarshal(data, &baseEvent); err != nil {
			return fmt.Errorf("failed to parse event type: %w", err)
		}

		// Handle different event types
		switch baseEvent.EventType {
		case "listenKeyExpired":
			event, err := ParseListenKeyExpiredEvent(data)
			if err != nil {
				return fmt.Errorf("failed to parse listen key expired event: %w", err)
			}
			fmt.Printf("Listen key expired: %s\n", event.ListenKey)

		case "ACCOUNT_UPDATE":
			event, err := ParseAccountUpdateEvent(data)
			if err != nil {
				return fmt.Errorf("failed to parse account update event: %w", err)
			}
			fmt.Printf("Account update - Event reason: %s\n", event.UpdateData.EventReasonType)
			fmt.Printf("Balances updated: %d\n", len(event.UpdateData.Balances))
			fmt.Printf("Positions updated: %d\n", len(event.UpdateData.Positions))

		case "ORDER_TRADE_UPDATE":
			event, err := ParseOrderTradeUpdateEvent(data)
			if err != nil {
				return fmt.Errorf("failed to parse order trade update event: %w", err)
			}
			fmt.Printf("Order update - Symbol: %s, Side: %s, Status: %s\n",
				event.Order.Symbol, event.Order.Side, event.Order.OrderStatus)

		case "MARGIN_CALL":
			event, err := ParseMarginCallEvent(data)
			if err != nil {
				return fmt.Errorf("failed to parse margin call event: %w", err)
			}
			fmt.Printf("Margin call - Cross wallet balance: %s\n", event.CrossWalletBalance)
			fmt.Printf("Positions in margin call: %d\n", len(event.Positions))

		case "TRADE_LITE":
			event, err := ParseTradeLiteEvent(data)
			if err != nil {
				return fmt.Errorf("failed to parse trade lite event: %w", err)
			}
			fmt.Printf("Trade lite - Symbol: %s, Side: %s, Quantity: %s, Price: %s\n",
				event.Symbol, event.Side, event.Quantity, event.Price)

		case "ACCOUNT_CONFIG_UPDATE":
			event, err := ParseAccountConfigUpdateEvent(data)
			if err != nil {
				return fmt.Errorf("failed to parse account config update event: %w", err)
			}
			fmt.Printf("Account config update - Symbol: %s, Leverage: %d\n",
				event.AccountConfig.Symbol, event.AccountConfig.Leverage)

		default:
			fmt.Printf("Unknown event type: %s\n", baseEvent.EventType)
		}

		return nil
	}

	// Subscribe to user data stream
	unsubscribe, err := wsClient.SubscribeToUserDataStream(callback)
	if err != nil {
		fmt.Printf("Failed to subscribe to user data stream: %v\n", err)
		return
	}

	fmt.Println("Subscribed to user data stream")

	// Keep the connection alive for some time
	time.Sleep(30 * time.Second)

	// Unsubscribe and close the stream
	if err := unsubscribe(); err != nil {
		fmt.Printf("Failed to unsubscribe: %v\n", err)
	}

	fmt.Println("Unsubscribed from user data stream")
}
