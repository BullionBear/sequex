package main

/*
Nicholas is a backtest trading bot associated with alex and demetrius.
*/

import (
	"context"
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/BullionBear/crypto-trade/domain/alpha"
	"github.com/BullionBear/crypto-trade/domain/chronicler"
	"github.com/BullionBear/crypto-trade/domain/config"
	"github.com/BullionBear/crypto-trade/domain/feedclient"
	"github.com/BullionBear/crypto-trade/domain/models"
	"github.com/BullionBear/crypto-trade/domain/shutdown"
	"github.com/BullionBear/crypto-trade/domain/wallet"
	"github.com/BullionBear/crypto-trade/trade/nikolaos"
	"github.com/shopspring/decimal"
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
	shutdown.HookShutdownCallback("cleanup", func() {
		logrus.Info("receive signal, start cleanup resources")
	})

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
	// chronicler
	chronicler := chronicler.NewChronicler(client, "nikolaos")
	shutdown.HookShutdownCallback("close chronicler", chronicler.Close)
	// wallet
	wallet := wallet.NewWallet()
	for _, b := range nikoConfig.Balance {
		wallet.SetBalance(b.Coin, decimal.NewFromFloat(b.Amount))
	}
	// Nikolas Strategy
	niko := nikolaos.NewNikolaos(wallet, alpha, chronicler)

	// Run niko
	go feedClient.SubscribeKlines(func(event *models.Kline) {
		niko.MakeDecision(event)
	})

	shutdown.WaitForShutdown(os.Interrupt, syscall.SIGTERM)
	logrus.Info("Nikolaos is done!")
}
