package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
	ctx := context.Background()
	c := binance.NewClient(&binance.Config{
		APIKey:    "NapwpjlsGTEgmc4cMgj8oA7zHzuUeAgRRj5hu0ZAKXA6XFYl3KYgiQd3YV9eVYrb",
		APISecret: "x3zyS1epGBKsz7KT4TqzIGFCKPLmsdFnEkt5EUioTgasgGaj8uXzhqDXIsRrXDMc",
		BaseURL:   "https://testnet.binance.vision",
	})
	// Single symbol
	price, err := c.GetTickerPrice(ctx, "BTCUSDT")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BTC price: %+v\n", price)
}
