package main

import (
	"flag"
	"time"

	"github.com/BullionBear/crypto-trade/internal/alpha"
	"github.com/BullionBear/crypto-trade/internal/config"
	"github.com/sirupsen/logrus"
)

// Config represents the structure of the configuration file.
func main() {
	logrus.SetLevel(logrus.InfoLevel)
	// Define a flag for the configuration file path.
	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	cgf, err := config.LoadConfig(*configPath)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"config path": *configPath}).Fatalln("Can't read config")
	}

	// Log the configuration to demonstrate that it's been loaded.
	logrus.Info("Loaded configuration", cgf)

	janus := alpha.NewJanus()
	go janus.Start()

	// Simulate pushing a model to Janus for processing
	go func() {
		janus.Channel <- alpha.Kline{StartTime: 123, EndTime: 456}
	}()
	// Receive the processed data
	processedData := <-janus.OutputChannel()
	logrus.Println(processedData.Alpha)

	// Give some time for the example to run
	time.Sleep(1 * time.Second)
}
