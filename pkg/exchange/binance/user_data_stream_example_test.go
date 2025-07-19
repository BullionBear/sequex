package binance

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example_userDataStream demonstrates how to use the user data stream client
// to receive real-time account and order updates.
func Example_userDataStream() {
	// Create user data stream client with testnet credentials
	// Note: This requires API key and secret for authentication
	config := TestnetConfig()
	config.APIKey = "your_api_key"
	config.APISecret = "your_api_secret"

	client := NewUserDataStreamClient(config)

	// Set up event handlers
	client.OnAccountUpdate(func(event *WSAccountUpdate) {
		fmt.Printf("Account Update: %d balances updated\n", len(event.Balances))
		for _, balance := range event.Balances {
			if balance.Free != "0.00000000" || balance.Locked != "0.00000000" {
				fmt.Printf("  %s: Free=%s, Locked=%s\n",
					balance.Asset, balance.Free, balance.Locked)
			}
		}
	})

	client.OnExecutionReport(func(event *WSExecutionReport) {
		fmt.Printf("Order Update: %s %s %s OrderID:%d Status:%s\n",
			event.Symbol, event.Side, event.OrderType, event.OrderID, event.CurrentOrderStatus)

		if event.IsNewOrder() {
			fmt.Println("  → New order created")
		}
		if event.IsTrade() {
			fmt.Printf("  → Trade executed: %s @ %s\n",
				event.LastExecutedQuantity, event.LastExecutedPrice)
		}
		if event.IsCanceled() {
			fmt.Println("  → Order canceled")
		}
		if event.IsFilled() {
			fmt.Println("  → Order completely filled")
		}
	})

	client.OnBalanceUpdate(func(event *WSBalanceUpdate) {
		fmt.Printf("Balance Update: %s Delta=%s\n", event.Asset, event.BalanceDelta)
	})

	// Connect to user data stream
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	fmt.Println("✅ Connected to user data stream")

	// Listen for events
	time.Sleep(30 * time.Second)

	// Disconnect
	err = client.Disconnect()
	if err != nil {
		log.Printf("Failed to disconnect: %v", err)
	}

	fmt.Println("🎉 User data stream example completed")
}

// Example_listenKeyManagement demonstrates manual listen key management
func Example_listenKeyManagement() {
	// Create REST client for manual listen key management
	config := TestnetConfig()
	config.APIKey = "your_api_key"
	config.APISecret = "your_api_secret"

	client := NewClient(config)
	ctx := context.Background()

	// Create listen key
	streamResp, err := client.CreateUserDataStream(ctx)
	if err != nil {
		log.Fatalf("Failed to create user data stream: %v", err)
	}

	fmt.Printf("Created listen key: %s\n", streamResp.ListenKey[:8]+"...")

	// Keep alive (should be done every 30-60 minutes)
	err = client.KeepAliveUserDataStream(ctx, streamResp.ListenKey)
	if err != nil {
		log.Printf("Failed to keep alive: %v", err)
	} else {
		fmt.Println("✅ Listen key kept alive")
	}

	// Close when done
	err = client.CloseUserDataStream(ctx, streamResp.ListenKey)
	if err != nil {
		log.Printf("Failed to close: %v", err)
	} else {
		fmt.Println("✅ Listen key closed")
	}

	fmt.Println("Listen key management example completed")
}
