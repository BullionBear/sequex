package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
	c := binance.NewClient(&binance.Config{
		APIKey:    "NapwpjlsGTEgmc4cMgj8oA7zHzuUeAgRRj5hu0ZAKXA6XFYl3KYgiQd3YV9eVYrb",
		APISecret: "x3zyS1epGBKsz7KT4TqzIGFCKPLmsdFnEkt5EUioTgasgGaj8uXzhqDXIsRrXDMc",
		BaseURL:   "https://testnet.binance.vision",
	})
	res, err := c.GetServerTime(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	// get exchange info
	res2, err := c.GetExchangeInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res2)
}
