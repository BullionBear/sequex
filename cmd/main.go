package main

import (
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
)

func main() {
	fmt.Printf("Hello, world.\n")
	symbol := "BTCUSDT"
	doneC, stopC, err := binance.WsKlineServe(symbol, "1m",
		func(event *binance.WsKlineEvent) { fmt.Printf("Kline: %+v\n", event) },
		func(err error) { fmt.Printf("Error: %v\n", err) })
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer func() {
		stopC <- struct{}{}
		<-doneC
	}()
	time.Sleep(20 * time.Second)

}
