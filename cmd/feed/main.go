package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/pkg/logger"
)

var (
	exchange string
	dataType string
	natsURIs string
)

// runFeed executes the main feed logic
func runFeed() {
	// Output version information
	logger.Log.Info().
		Str("version", env.Version).
		Str("buildTime", env.BuildTime).
		Str("commitHash", env.CommitHash).
		Msg("Feed started")

	// Validate inputs
	if err := validateInputs(); err != nil {
		logger.Log.Error().Err(err).Msg("Validation failed")
		os.Exit(1)
	}

	// Print configuration
	printConfiguration()

	// TODO: Implement actual feed logic here
	logger.Log.Info().Msg("Feed command executed successfully!")
}

// validateInputs validates the command line arguments
func validateInputs() error {
	// Validate exchange
	validExchanges := []string{"binance", "binanceperp", "okx", "bybit"}
	if !contains(validExchanges, exchange) {
		return fmt.Errorf("invalid exchange '%s'. Supported exchanges: %s",
			exchange, strings.Join(validExchanges, ", "))
	}

	// Validate data type
	validDataTypes := []string{"trades", "klines", "depth", "ticker", "book"}
	if !contains(validDataTypes, dataType) {
		return fmt.Errorf("invalid data type '%s'. Supported data types: %s",
			dataType, strings.Join(validDataTypes, ", "))
	}

	// Validate NATS URIs
	if natsURIs == "" {
		return fmt.Errorf("NATS URIs cannot be empty")
	}

	// Check if URIs contain valid NATS protocol
	uris := strings.Split(natsURIs, ",")
	for _, uri := range uris {
		uri = strings.TrimSpace(uri)
		if !strings.HasPrefix(uri, "nats://") && !strings.HasPrefix(uri, "tls://") {
			return fmt.Errorf("invalid NATS URI '%s'. Must start with 'nats://' or 'tls://'", uri)
		}
	}

	return nil
}

// printConfiguration prints the parsed configuration
func printConfiguration() {
	logger.Log.Info().
		Str("exchange", exchange).
		Str("dataType", dataType).
		Str("natsURIs", natsURIs).
		Msg("Feed Configuration")
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
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
  feed binance trades nats://localhost:4222
  feed binance klines nats://localhost:4222,nats://localhost:4223
  feed binance depth nats://localhost:4222
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
	exchange = args[0]
	dataType = args[1]
	natsURIs = args[2]

	// Run the main logic
	runFeed()
}
