package main

import (
	"os"
	"time"

	"github.com/BullionBear/sequex/pkg/log"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
	logger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
		log.WithOutput(os.Stdout),
	)

	logger.Info("Starting WebSocket trade subscription...")

	wsClient := binance.NewWSClient(&binance.WSConfig{})
	unsubscribe, err := wsClient.SubscribeTrade("BTCUSDT", binance.TradeSubscriptionOptions{
		OnConnect: func() {
			logger.Info("WebSocket connected successfully")
		},
		OnError: func(err error) {
			logger.Error("WebSocket error", log.Error(err))
		},
		OnDisconnect: func() {
			logger.Info("WebSocket disconnected")
		},
		OnTrade: func(trade binance.WSTrade) {
			logger.Info("Trade received",
				log.String("symbol", trade.Symbol),
				log.String("price", trade.Price),
				log.String("quantity", trade.Quantity),
				log.Int64("tradeId", trade.TradeId),
			)
		},
	})

	if err != nil {
		logger.Fatal("Failed to subscribe to trade",
			log.Error(err),
		)
	}

	logger.Info("Subscription created successfully, waiting for trade data...")

	// Keep the program running for a while to receive trade data
	time.Sleep(10 * time.Second)

	unsubscribe()
	logger.Info("Unsubscribed from trade")
}
