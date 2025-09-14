package main

import (
	"flag"
	"log"
	"time"

	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/internal/model/sqx"
	"github.com/BullionBear/sequex/pkg/logger"

	"github.com/nats-io/nats.go"
)

func main() {
	// Parse command line flags
	var n = flag.Int("n", 100, "Number of latest data points to query")
	flag.Usage = func() {
		logger.Log.Info().Msg(`Cache is a CLI tool for querying trade data from the NATS cache consumer.

The cache creates ephemeral consumers positioned to fetch the latest N messages from
the TRADE stream. Ephemeral consumers are automatically cleaned up after use.
Messages are NAKed (not ACKed) to keep them available for subsequent cache queries.

Usage:
  cache [-n N]

Examples:
  cache           # Query latest 100 data points (default)
  cache -n 10     # Query latest 10 data points
  cache -n 500    # Query latest 500 data points

Prerequisites:
  - NATS server must be running with JetStream enabled
  - TRADE stream and TRADE_CACHE_BTCUSDT consumer must exist
  - Run './script.sh create' to set up the required infrastructure
`)
		flag.PrintDefaults()
	}
	flag.Parse()

	logger.Log.Info().
		Str("version", env.Version).
		Str("buildTime", env.BuildTime).
		Str("commitHash", env.CommitHash).
		Int("dataPoints", *n).
		Msg("Cache started")
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("Failed to create JetStream context:", err)
	}

	// Get stream info to calculate starting position for latest N messages
	streamInfo, err := js.StreamInfo("TRADE")
	if err != nil {
		log.Fatal("TRADE stream does not exist. Please run './script.sh create' to set up the infrastructure:", err)
	}

	// Calculate starting sequence for the latest N messages
	startSeq := streamInfo.State.LastSeq - uint64(*n) + 1
	if startSeq < 1 {
		startSeq = 1
	}

	logger.Log.Info().Msgf("Stream has %d messages, fetching latest %d starting from sequence %d",
		streamInfo.State.Msgs, *n, startSeq)

	// Create an ephemeral consumer starting from the calculated position for latest N messages
	// Ephemeral consumers are automatically cleaned up when the subscription is closed
	sub, err := js.PullSubscribe("trade.btcusdt", "", nats.StartSequence(startSeq))
	if err != nil {
		log.Fatal("Failed to create pull subscription for cache:", err)
	}
	defer sub.Unsubscribe()

	// Fetch the latest N messages
	msgs, err := sub.Fetch(*n, nats.MaxWait(5*time.Second))
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch messages from cache")
	} else {
		logger.Log.Info().Msgf("Fetched %d messages from cache", len(msgs))

		for i, msg := range msgs {
			trade := &sqx.Trade{}
			err := sqx.Unmarshal(msg.Data, trade)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to unmarshal trade")
				// NAK the message to keep it available for future cache queries
				msg.Nak()
				continue
			}
			logger.Log.Info().Msgf("[%d] Trade: %s", i+1, trade.IdStr())
			// NAK instead of ACK to keep the message available as cache data
			// This ensures the same data remains accessible for subsequent cache queries
			msg.Nak()
		}
	}

}
