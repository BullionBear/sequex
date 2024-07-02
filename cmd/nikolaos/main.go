package main

/*
Nicholas is a backtest trading bot associated with alex and demetrius.
*/

import (
	"context"
	"flag"
	"log"

	"github.com/BullionBear/crypto-trade/domain/alpha"
	"github.com/BullionBear/crypto-trade/domain/config"
	"github.com/BullionBear/crypto-trade/domain/feedclient"
	"github.com/BullionBear/crypto-trade/domain/models"
	"github.com/BullionBear/crypto-trade/domain/reporter"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// New raw resources
	// grpc client
	feedConfig := nikoConfig.GrpcClient
	feedClient := feedclient.NewFeedClient(feedConfig.Host, feedConfig.Port)
	defer feedClient.Close()
	// mongo client
	mongoUri := nikoConfig.MongoUri
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// Dependency injection wrapped resources
	// alpha
	alpha := alpha.NewAlpha()
	reporter := reporter.NewReporter(client)
	reporter.Record("alpha", alpha)
	feedClient.SubscribeKlines(func(event *models.Kline) {
		alpha.Append(event)
		lm := alpha.LongCloseMovingAvg.Mean()
		sm := alpha.ShortCloseMovingAvg.Mean()
		logrus.Infof("Long moving average: %f, Short moving average: %f", lm, sm)
		// logrus.Infof("Received kline %+v", event)
	})
	doneC := make(chan struct{})
	<-doneC
}
