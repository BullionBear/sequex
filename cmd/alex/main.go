package main

import (
	"flag"

	"github.com/BullionBear/crypto-trade/domain/config"
	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		logrus.Fatal("Please provide a path to the configuration file")
	}
	alexConfig, err := config.LoadAlexConfig(*configPath)
	if err != nil {
		logrus.Fatal("Can't read config: ", err)
	}
	logrus.Infof("Load config with %+v", *alexConfig)

}
