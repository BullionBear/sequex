package main

/*
Nicholas is a backtest trading bot associated with alex and demetrius.
*/

import (
	"flag"
	"fmt"

	"github.com/BullionBear/crypto-trade/domain/config"
	"github.com/BullionBear/crypto-trade/domain/feedclient"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Read the configuration file
	configPath := flag.String("config", "", "path of nikolaos config file")
	flag.Parse()

	if *configPath == "" {
		logrus.Fatal("Please provide a path to the configuration file")
	}
	nikoConfig, err := config.LoadNikoConfig(*configPath)
	if err != nil {
		logrus.Fatal("Can't read config: ", err)
	}
	logrus.Infof("Load config with %+v", *nikoConfig)

	// New resources
	// grpc client
	grpcConfig := nikoConfig.GrpcClient
	srvAddr := fmt.Sprintf("%s:%d", grpcConfig.Host, grpcConfig.Port)
	conn, err := grpc.NewClient(srvAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	grpcClient := feedclient.NewFeedClient(conn)
	grpcClient.SubscribeKlines(func(event *feedclient.Kline) {
		logrus.Infof("Received kline: %+v", event)
	})

	// NewNiko

}
