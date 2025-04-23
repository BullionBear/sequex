package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	orderapi "github.com/BullionBear/sequex/api/order"
	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/internal/order"
	"github.com/BullionBear/sequex/internal/orderbook"
	pb "github.com/BullionBear/sequex/pkg/protobuf/order" // Correct import path
	"google.golang.org/grpc"
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
	orderbookManager := orderbook.NewBinanceOrderBookManager()
	for _, symbol := range conf.Market["binance"] {
		orderbookManager.CreateOrderBook(symbol, orderbook.UpdateSpeed1s)
		defer orderbookManager.CloseOrderBook(symbol)
	}
	orderManagers := make(map[string]*order.BinanceOrderManager)
	for _, account := range conf.Accounts["binance"] {
		orderManager := order.NewBinanceOrderManager(account.APIKey, account.APISecret, orderbookManager)
		orderManagers[account.Name] = orderManager
	}
	orderService := orderapi.NewBinanceOrderService(orderManagers["scylla"])
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBinanceOrderServiceServer(s, orderService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
