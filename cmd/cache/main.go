package main

import (
	"time"

	"github.com/BullionBear/sequex/internal/model/sqx"
	"github.com/BullionBear/sequex/pkg/logger"
	"github.com/nats-io/nats.go"
)

func main() {
	natUrls := "nats://localhost:4222"
	natConn, err := nats.Connect(natUrls)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}
	defer natConn.Close()
	logger.Log.Info().Msg("Connected to NATS server")

	js, err := natConn.JetStream()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to create JetStream context")
	}

	consumer, err := js.ConsumerInfo("TRADE", "TRADE_PUBSUB")
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to get consumer info")
	}
	logger.Log.Info().Interface("consumer", consumer).Msg("Consumer info retrieved")

	sub, err := js.Subscribe("trade.btcusdt", func(msg *nats.Msg) {
		trade := sqx.Trade{}
		err := sqx.Unmarshal(msg.Data, &trade)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to unmarshal trade data")
			return
		}
		logger.Log.Info().Interface("trade", trade).Msg("Received trade")
	})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to subscribe to trade messages")
	}
	defer sub.Unsubscribe()
	logger.Log.Info().Msg("Subscribed to trade.btcusdt")

	logger.Log.Info().Msg("Cache service running for 10 seconds...")
	time.Sleep(30 * time.Second)
	logger.Log.Info().Msg("Cache service shutting down")
}
