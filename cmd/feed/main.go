package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/internal/pubsub"
	"github.com/BullionBear/sequex/pkg/logger"
	"github.com/nats-io/nats.go"
)

// runFeed executes the main feed logic
func runFeed(exchange string, dataType string, natsURIs string) {
	// Output version information
	logger.Log.Info().
		Str("version", env.Version).
		Str("buildTime", env.BuildTime).
		Str("commitHash", env.CommitHash).
		Msg("Feed started")

	printConfiguration(exchange, dataType, natsURIs)

	// Validate inputs
	connConfigs, err := parseNatsURIs(natsURIs)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Validation failed")
		os.Exit(1)
	}

	// Print configuration
	natsConn, err := nats.Connect(connConfigs[0].ToNATSURL())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to connect to NATS")
		os.Exit(1)
	}

	defer natsConn.Close()
	js, err := natsConn.JetStream()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create JetStream context")
		os.Exit(1)
	}

	streamInfo, err := js.StreamInfo(connConfigs[0].GetParam("stream", ""))
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get stream info")
		os.Exit(1)
	}

	logger.Log.Info().Msg("Stream info:")
	logger.Log.Info().Msg(streamInfo.Config.Name)

	publisher, err := pubsub.NewPublisher(natsConn, streamInfo.Config.Name, connConfigs[0].GetParam("subject", ""))
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create publisher")
		os.Exit(1)
	}
	publisher.Publish([]byte("Hello, world!"))

	// TODO: Implement actual feed logic here
	logger.Log.Info().Msg("Feed command executed successfully!")
}

// validateInputs validates the command line arguments
func parseNatsURIs(natsURIs string) ([]*config.ConnectionConfig, error) {
	// Validate NATS URIs using the connection string parser
	if natsURIs == "" {
		return nil, fmt.Errorf("NATS URIs cannot be empty")
	}

	var connConfigs []*config.ConnectionConfig

	// Parse and validate each URI
	uris := strings.Split(natsURIs, ",")
	for _, uri := range uris {
		uri = strings.TrimSpace(uri)
		if uri == "" {
			continue
		}

		connConfig, err := config.ParseConnectionString(uri)
		if err != nil {
			return nil, fmt.Errorf("invalid connection string '%s': %w", uri, err)
		}

		// Validate the parsed configuration
		if err := connConfig.Validate(); err != nil {
			return nil, fmt.Errorf("invalid connection configuration for '%s': %w", uri, err)
		}

		// Log parsed configuration details
		logger.Log.Debug().
			Str("uri", uri).
			Str("host", connConfig.Host).
			Int("port", connConfig.Port).
			Str("username", connConfig.Username).
			Interface("params", connConfig.Params).
			Msg("Parsed connection string")

		connConfigs = append(connConfigs, connConfig)
	}

	return connConfigs, nil
}

// printConfiguration prints the parsed configuration
func printConfiguration(exchange string, dataType string, natsURIs string) {
	logger.Log.Info().
		Str("exchange", exchange).
		Str("dataType", dataType).
		Str("natsURIs", natsURIs).
		Msg("Feed Configuration")
}

func main() {
	// Define flags

	// Custom usage function
	flag.Usage = func() {
		logger.Log.Info().Msg(`Feed is a scalable CLI tool for streaming market data from various exchanges
to NATS message brokers. It supports multiple exchanges and data types.

Usage:
  feed <exchange> <data-type> <nats-uris>

Examples:
  feed binance trades 'nats://localhost:4222?stream=feed&subject=test'
  feed binance klines 'nats://localhost:4223?stream=feed&subject=test,nats://localhost:4224?stream=feed&subject=test'
`)
		flag.PrintDefaults()
	}

	// Parse flags
	flag.Parse()

	// Check if we have the required positional arguments
	args := flag.Args()
	if len(args) != 3 {
		logger.Log.Error().Msg("exactly 3 arguments required: <exchange> <data-type> <nats-uris>")
		flag.Usage()
		os.Exit(1)
	}

	// Parse positional arguments
	exchange := args[0]
	dataType := args[1]
	natsURIs := args[2]

	// Run the main logic
	runFeed(exchange, dataType, natsURIs)
}
