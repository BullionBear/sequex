package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BullionBear/sequex/env"
)

var (
	exchange string
	dataType string
	natsURIs string
	version  bool
)

// runFeed executes the main feed logic
func runFeed() {
	// Validate inputs
	if err := validateInputs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print configuration
	printConfiguration()

	// TODO: Implement actual feed logic here
	fmt.Println("Feed command executed successfully!")
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
	fmt.Println("=== Feed Configuration ===")
	fmt.Printf("Exchange: %s\n", exchange)
	fmt.Printf("Data Type: %s\n", dataType)
	fmt.Printf("NATS URIs: %s\n", natsURIs)
	fmt.Println("==========================")
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
	flag.BoolVar(&version, "version", false, "Show version information")
	flag.BoolVar(&version, "v", false, "Show version information (shorthand)")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Feed is a scalable CLI tool for streaming market data from various exchanges
to NATS message brokers. It supports multiple exchanges and data types.

Usage:
  feed <exchange> <data-type> <nats-uris> [flags]

Examples:
  feed binance trades nats://localhost:4222
  feed binance klines nats://localhost:4222,nats://localhost:4223
  feed binance depth nats://localhost:4222

Flags:
`)
		flag.PrintDefaults()
	}

	// Parse flags
	flag.Parse()

	// Handle version flag
	if version {
		fmt.Printf("Version: %s\nBuild Time: %s\nCommit Hash: %s\n",
			env.Version, env.BuildTime, env.CommitHash)
		return
	}

	// Check if we have the required positional arguments
	args := flag.Args()
	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "Error: exactly 3 arguments required: <exchange> <data-type> <nats-uris>\n\n")
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
