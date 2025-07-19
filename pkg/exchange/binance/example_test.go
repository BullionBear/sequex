package binance_test

import (
	"context"
	"fmt"
	"log"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func ExampleClient_GetServerTime() {
	// Create a new client with default configuration
	client := binance.NewClient(nil)

	// Create a context
	ctx := context.Background()

	// Get server time
	serverTime, err := client.GetServerTime(ctx)
	if err != nil {
		log.Fatalf("Failed to get server time: %v", err)
	}

	// Print the server time
	fmt.Printf("Server timestamp: %d\n", serverTime.ServerTime)
	fmt.Printf("Server time: %v\n", serverTime.GetTime())
}

func ExampleClient_Ping() {
	// Create a new client with default configuration
	client := binance.NewClient(nil)

	// Create a context
	ctx := context.Background()

	// Test connectivity
	err := client.Ping(ctx)
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	fmt.Println("Ping successful - connection established")
}

func ExampleClient_GetTickerPrice() {
	// Create a new client with default configuration
	client := binance.NewClient(nil)

	// Create a context
	ctx := context.Background()

	// Get price for a specific symbol
	result, err := client.GetTickerPrice(ctx, "BTCUSDT")
	if err != nil {
		log.Fatalf("Failed to get ticker price: %v", err)
	}

	ticker := result.(*binance.TickerPriceResponse)
	fmt.Printf("Symbol: %s, Price: %s\n", ticker.Symbol, ticker.Price)
}

func ExampleNewClient() {
	// Create client with default configuration
	client1 := binance.NewClient(nil)
	fmt.Printf("Default client base URL: %s\n", client1.GetConfig().BaseURL)

	// Create client with custom configuration
	config := binance.DefaultConfig()
	config.APIKey = "your-api-key"
	config.APISecret = "your-api-secret"

	client2 := binance.NewClient(config)
	fmt.Printf("Custom client base URL: %s\n", client2.GetConfig().BaseURL)

	// Create testnet client
	testnetConfig := binance.TestnetConfig()
	testnetConfig.APIKey = "your-testnet-api-key"
	testnetConfig.APISecret = "your-testnet-api-secret"

	client3 := binance.NewClient(testnetConfig)
	fmt.Printf("Testnet client base URL: %s\n", client3.GetConfig().BaseURL)
}
