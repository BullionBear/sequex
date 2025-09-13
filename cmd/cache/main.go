package main

import (
	"context"
	"fmt"
	"time"

	"github.com/BullionBear/sequex/internal/model/sqx"
	"github.com/BullionBear/sequex/pkg/logger"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})
	defer rdb.Close()

	// Test Redis connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	logger.Log.Info().Msg("Connected to Redis server")

	natUrls := "nats://localhost:4222"
	natConn, err := nats.Connect(natUrls)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}
	defer natConn.Close()
	logger.Log.Info().Msg("Connected to NATS server")
	/*
		js, err := natConn.JetStream()
		if err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to create JetStream context")
		}

		// Verify consumer exists (should be created by script.sh)

			consumerInfo, err := js.ConsumerInfo("TRADE", "TRADE_PUBSUB")
			if err != nil {
				logger.Log.Fatal().Err(err).Msg("Consumer 'TRADE_PUBSUB' does not exist. Please run setup scripts first.")
			}
			logger.Log.Info().Interface("consumer", consumerInfo).Msg("Consumer info retrieved")
	*/
	// Subscribe to the deliver_subject configured for this push consumer
	// For push consumers, we use regular NATS subscription to the deliver_subject
	sub, err := natConn.Subscribe("fanout.btcusdt", func(msg *nats.Msg) {
		// Unmarshal once and use for both Redis time series and logging
		trade := sqx.Trade{}
		err := sqx.Unmarshal(msg.Data, &trade)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to unmarshal trade")
			return
		}

		// Add raw message data to Redis
		err = addRawDataToRedis(ctx, rdb, msg.Data, trade)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to add raw data to Redis")
			return
		}

		logger.Log.Info().
			Str("subject", msg.Subject).
			Interface("trade", trade).
			Msg("Received message")

		// No acknowledgment needed since this consumer has ack_policy: "none"
	})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to subscribe to deliver subject")
	}
	defer sub.Unsubscribe()
	logger.Log.Info().Msg("Successfully subscribed to fanout.btcusdt")
	time.Sleep(30 * time.Second)
}

// addRawDataToRedis stores raw message bytes in Redis time series bucket
func addRawDataToRedis(ctx context.Context, rdb *redis.Client, data []byte, trade sqx.Trade) error {
	// Create time series key based on exchange, instrument, and symbol
	tsKey := fmt.Sprintf("ts:trades:%s:%s:%s",
		trade.Exchange.String(),
		trade.InstrumentType.String(),
		trade.Symbol.String())

	// Store the raw bytes using timestamp as the time series key
	// Note: Redis time series only supports numeric values, so we'll store as a regular key-value
	// with a timestamp-based key for time series-like behavior
	timestampKey := fmt.Sprintf("%s:%d", tsKey, trade.Timestamp)
	err := rdb.Set(ctx, timestampKey, data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to add raw data to Redis: %w", err)
	}

	logger.Log.Debug().
		Str("key", tsKey).
		Int64("timestamp", trade.Timestamp).
		Int("data_size", len(data)).
		Str("trade_id", trade.IdStr()).
		Msg("Added raw trade data to Redis time series")

	return nil
}
