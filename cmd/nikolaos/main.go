package main

/*
Nicholas is a backtest trading bot associated with alex and demetrius.
*/

import (
	"flag"

	"github.com/BullionBear/crypto-trade/domain/config"
	"github.com/sirupsen/logrus"
)

func main() {
	// Read the configuration file
	configPath := flag.String("config", "", "path of nikolaos config file")
	flag.Parse()

	if *configPath == "" {
		logrus.Fatal("Please provide a path to the configuration file")
	}
	nikoConfig, err := config.LoadAlexConfig(*configPath)
	if err != nil {
		logrus.Fatal("Can't read config: ", err)
	}
	logrus.Infof("Load config with %+v", *nikoConfig)

	// New resources

	// NewAlex

}
