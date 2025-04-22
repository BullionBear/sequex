package main

import (
	"flag"
	"fmt"

	"github.com/BullionBear/sequex/internal/config"
)

func main() {
	// -c flag to specify the configuration file
	path := flag.String("c", "config.yml", "Path to the configuration file")
	flag.Parse()

	// Use the flag value
	fmt.Println("Config file:", *path)
	// Load the configuration
	conf, err := config.LoadConfig(*path)
	if err != nil {
		panic("Error loading config " + err.Error())
	}
	// orderbookManager := orderbook.NewBinanceOrderBookManager()
	for _, symbol := range conf.Market["binance"] {
		fmt.Printf("Symbol: %s\n", symbol)
	}
}
