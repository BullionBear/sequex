package main

import (
	"flag"
	"net"

	orderapi "github.com/BullionBear/sequex/api/order"
	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/internal/order"
	"github.com/BullionBear/sequex/internal/orderbook"
	"github.com/BullionBear/sequex/pkg/log"
	pb "github.com/BullionBear/sequex/pkg/protobuf/order" // Correct import path
	"google.golang.org/grpc"
)

func main() {
	// -c flag to specify the configuration file
	logger, err := log.NewLogger(log.InfoLevel, "stdout", "sequex.log")
	if err != nil {
		panic("Error creating logger " + err.Error())
	}
	defer logger.Close()

	path := flag.String("c", "config.yml", "Path to the configuration file")
	flag.Parse()

	// Use the flag value
	logger.Info("Config file: %s", *path)
	// Load the configuration
	conf, err := config.LoadConfig(*path)
	if err != nil {
		logger.Fatal("Error loading config " + err.Error())
	}
	orderbookManager := orderbook.NewBinanceOrderBookManager(logger.WithKV(log.KV{Key: "market", Value: "binance"}))
	for _, symbol := range conf.Market["binance"] {
		orderbookManager.CreateOrderBook(symbol, orderbook.UpdateSpeed1s)
		defer orderbookManager.CloseOrderBook(symbol)
	}
	orderManagers := make(map[string]*order.BinanceOrderManager)
	for _, account := range conf.Accounts["binance"] {
		logger.Info("Creating order manager for account %s, %s", account.APIKey, account.APISecret)
		orderManager := order.NewBinanceOrderManager(account.APIKey, account.APISecret, orderbookManager, logger.WithKV(log.KV{Key: "account", Value: account.Name}))
		orderManagers[account.Name] = orderManager
	}
	orderService := orderapi.NewBinanceOrderService(orderManagers["scylla"], logger.WithKV(log.KV{Key: "orderapi", Value: "scylla"}))
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBinanceOrderServiceServer(s, orderService)
	logger.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Info("failed to serve: %v", err)
	}
}
