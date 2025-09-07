package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	exchange string
	dataType string
	natsURIs string
	verbose  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "feed <exchange> <data-type> <nats-uris>",
	Short: "Feed market data from exchanges to NATS",
	Long: `Feed is a scalable CLI tool for streaming market data from various exchanges
to NATS message brokers. It supports multiple exchanges and data types.

Examples:
  feed binance trades nats://localhost:4222
  feed binance klines nats://localhost:4222,nats://localhost:4223
  feed binance depth nats://localhost:4222 --verbose`,
	Args: cobra.ExactArgs(3),
	Run:  runFeed,
}

// runFeed executes the main feed logic
func runFeed(cmd *cobra.Command, args []string) {
	// Parse arguments
	exchange = args[0]
	dataType = args[1]
	natsURIs = args[2]

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
	fmt.Printf("Verbose: %t\n", verbose)
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

func init() {
	// Add flags
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Add examples
	rootCmd.Example = `  # Stream trades from Binance to local NATS
  feed binance trades nats://localhost:4222

  # Stream klines with multiple NATS servers
  feed binance klines nats://localhost:4222,nats://localhost:4223

  # Stream depth data with verbose output
  feed binance depth nats://localhost:4222 --verbose`
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
