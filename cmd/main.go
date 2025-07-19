package main

import (
	"fmt"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
	c := binance.NewWSStreamClient(binance.DefaultConfig())
	unsubscribe, _ := c.SubscribeToDiffDepthWithCallback("ETHUSDT", "@100ms", func(data *binance.WSDiffDepthData) error {
		fmt.Println(data)
		return nil
	})
	time.Sleep(10 * time.Second)
	unsubscribe()
}
