package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BullionBear/sequex/pkg/exchange/binanceperp"
)

func main() {
	fmt.Println("=== Binance Perpetual Futures REST API Example ===")
	restAPIExample()
}

func restAPIExample() {
	// Configure your API credentials here
	cfg := &binanceperp.Config{
		BaseURL: binanceperp.MainnetBaseUrl,
		// For signed requests, you would also need:
		// APIKey:    "your_api_key",
		// APISecret: "your_api_secret",
	}

	client := binanceperp.NewClient(cfg)

	// Example 1: Get Server Time (unsigned request)
	fmt.Println("\n--- Get Server Time ---")
	timeResp, err := client.GetServerTime(context.Background())
	if err != nil {
		log.Printf("GetServerTime error: %v", err)
		return
	}

	if timeResp.Code != 0 {
		log.Printf("GetServerTime failed with code %d: %s", timeResp.Code, timeResp.Message)
		return
	}

	fmt.Printf("Server Time: %d\n", timeResp.Data.ServerTime)
	fmt.Printf("Response Code: %d\n", timeResp.Code)
	fmt.Printf("Response Message: %s\n", timeResp.Message)
}
