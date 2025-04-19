package main

import (
	"flag"
	"fmt"

	"github.com/shopspring/decimal"
)

func main() {
	// Define command-line arguments
	symbol := flag.String("symbol", "BTCUSDT", "Trading pair symbol (e.g., BTCUSDT)")
	side := flag.String("side", "BUY", "Order side (BUY or SELL)")
	orderType := flag.String("type", "LIMIT", "Order type (e.g., LIMIT, MARKET)")
	orderTTL := flag.Int("order-ttl", 60, "Order time-to-live in seconds")

	positionTTL := flag.Int("position-ttl", 3600, "Position time-to-live in seconds")
	stopLoss := flag.String("stop-loss", "0.0", "Stop-loss price")
	stopProfit := flag.String("stop-profit", "1000000.0", "Stop-profit price")
	price := flag.String("price", "0.0", "Price for the order")
	trailingStopLoss := flag.String("trailing-stoploss", "0.0", "Trailing stop-loss value")

	// Parse the flags
	flag.Parse()

	// Convert stopLoss, stopProfit, and price to decimal.Decimal
	stopLossDecimal, err := decimal.NewFromString(*stopLoss)
	if err != nil {
		fmt.Printf("Invalid stop-loss value: %v\n", err)
		return
	}

	stopProfitDecimal, err := decimal.NewFromString(*stopProfit)
	if err != nil {
		fmt.Printf("Invalid stop-profit value: %v\n", err)
		return
	}

	priceDecimal, err := decimal.NewFromString(*price)
	if err != nil {
		fmt.Printf("Invalid price value: %v\n", err)
		return
	}

	// Convert trailingStopLoss to decimal.Decimal
	trailingStopLossDecimal, err := decimal.NewFromString(*trailingStopLoss)
	if err != nil {
		fmt.Printf("Invalid trailing-stoploss value: %v\n", err)
		return
	}

	// Print parsed arguments for debugging
	fmt.Printf("Symbol: %s\n", *symbol)
	fmt.Printf("Side: %s\n", *side)
	fmt.Printf("Type: %s\n", *orderType)
	fmt.Printf("Order TTL: %d seconds\n", *orderTTL)
	fmt.Printf("Position TTL: %d seconds\n", *positionTTL)
	fmt.Printf("Price: %s\n", priceDecimal.String())
	fmt.Printf("Stop Loss: %s\n", stopLossDecimal.String())
	fmt.Printf("Stop Profit: %s\n", stopProfitDecimal.String())
	fmt.Printf("Trailing StopLoss: %s\n", trailingStopLossDecimal.String())
}
