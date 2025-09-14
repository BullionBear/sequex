package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/internal/model/sqx"
	"github.com/BullionBear/sequex/pkg/logger"

	"github.com/nats-io/nats.go"
)

func main() {
	// Parse command line flags
	var n = flag.Int("n", 100, "Number of latest data points to query")
	var configFile = flag.String("c", "", "Configuration file path (optional, uses hardcoded values if not provided)")
	flag.Usage = func() {
		logger.Log.Info().Msg(`Cache is a CLI tool for querying trade data from the NATS cache consumer.

The cache creates ephemeral consumers positioned to fetch the latest N messages from
the TRADE stream. Ephemeral consumers are automatically cleaned up after use.
Messages are NAKed (not ACKed) to keep them available for subsequent cache queries.

Usage:
  cache [-n N] [-c <config-file>]

Examples:
  cache                                           # Query latest 100 data points (default)
  cache -n 10                                     # Query latest 10 data points
  cache -n 500                                    # Query latest 500 data points
  cache -c config/trade-binance-spot-btcusdt.json # Use config file for NATS connection

Prerequisites:
  - NATS server must be running with JetStream enabled
  - TRADE stream and consumer must exist
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

	// Determine NATS connection details and stream/subject info
	var natsURIs, streamName, subject string
	if *configFile != "" {
		// Load configuration from file
		cfg, err := config.LoadConfig(*configFile)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to load config")
			os.Exit(1)
		}
		natsURIs = cfg.NATS.URIs
		streamName = cfg.NATS.Stream
		subject = cfg.NATS.Subject
		logger.Log.Info().
			Str("configFile", *configFile).
			Str("natsURIs", natsURIs).
			Str("stream", streamName).
			Str("subject", subject).
			Msg("Using configuration from file")
	} else {
		// Use hardcoded defaults for backward compatibility
		natsURIs = "nats://localhost:4222"
		streamName = "TRADE"
		subject = "trade.btcusdt"
		logger.Log.Info().
			Str("natsURIs", natsURIs).
			Str("stream", streamName).
			Str("subject", subject).
			Msg("Using hardcoded configuration")
	}

	nc, err := nats.Connect(natsURIs)
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
	streamInfo, err := js.StreamInfo(streamName)
	if err != nil {
		log.Fatal("Stream does not exist. Please run './script.sh create' to set up the infrastructure:", err)
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
	sub, err := js.PullSubscribe(subject, "", nats.StartSequence(startSeq))
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
