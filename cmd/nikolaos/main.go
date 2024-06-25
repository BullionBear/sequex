package main

/*
Nicholas is a backtest trading bot associated with alex and demetrius.
*/

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/BullionBear/crypto-trade/api/gen/feed"
	"github.com/BullionBear/crypto-trade/domain/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
	feedClient := feed.NewFeedClient(conn)
	stream, err := feedClient.SubscribeKline(context.Background(), &emptypb.Empty{})
	if err != nil {
		logrus.Fatalf("could not get config: %v", err)
	}
	for {
		kline, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				logrus.Info("Stream closed by server")
				return
			} else {
				logrus.Info("Error receiving from kline stream: %v", status.Convert(err).Message())
			}
		}
		logrus.Infof("Received kline: %+v", kline)
	}

}
