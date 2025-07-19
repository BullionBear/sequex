package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
	fmt.Println("Hello, World!")
	c := binance.NewClient(binance.DefaultConfig())
	res, err := c.GetServerTime(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
