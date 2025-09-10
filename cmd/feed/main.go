package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/internal/adapter"
	_ "github.com/BullionBear/sequex/internal/adapter/init"
	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/internal/model/sqx"
	"github.com/BullionBear/sequex/internal/pubsub"
	"github.com/BullionBear/sequex/pkg/logger"
	"github.com/BullionBear/sequex/pkg/shutdown"
)

// runFeed executes the main feed logic
func runFeed(exchange string, instrument string, symbol string, dataType string, natsURIs string) {
	// Output version information
	logger.Log.Info().
		Str("version", env.Version).
		Str("buildTime", env.BuildTime).
		Str("commitHash", env.CommitHash).
		Msg("Feed started")

	printConfiguration(exchange, instrument, symbol, dataType, natsURIs)
	sqxExchange := sqx.NewExchange(exchange)
	if sqxExchange == sqx.ExchangeUnknown {
		logger.Log.Error().Msg("Invalid exchange")
		os.Exit(1)
	}

	sqxInstrumentType := sqx.NewInstrumentType(instrument)
	if sqxInstrumentType == sqx.InstrumentTypeUnknown {
		logger.Log.Error().Msg("Invalid instrument")
		os.Exit(1)
	}

	sqxSymbol, err := sqx.NewSymbolFromStr(symbol)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create symbol")
		os.Exit(1)
	}

	sqxDataType := sqx.NewDataType(dataType)
	if sqxDataType == sqx.DataTypeUnknown {
		logger.Log.Error().Msg("Invalid data type")
		os.Exit(1)
	}

	shutdown := shutdown.NewShutdown(logger.Log)

	connConfigs, err := parseNatsURIs(natsURIs)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Validation failed")
		os.Exit(1)
	}

	pubManager, err := pubsub.NewPubManager(connConfigs)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create pub manager")
		os.Exit(1)
	}
	defer pubManager.Close()

	switch sqxDataType {
	case sqx.DataTypeTrade:
		adapter, err := adapter.CreateTradeAdapter(sqxExchange)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to create adapter")
			os.Exit(1)
		}
		unsubscribe, err := adapter.Subscribe(sqxSymbol, sqxInstrumentType, func(trade sqx.Trade) error {
			_, err := trade.Marshal()
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to marshal trade")
				return err
			}
			logger.Log.Info().Msgf("Publishing trade: %s", trade.IdStr())
			/*return pubManager.Publish(data, map[string]string{
				"Nats-Msg-Id": trade.IdStr(),
			})*/
			return nil
		})
		shutdown.HookShutdownCallback("unsubscribe", unsubscribe, 10*time.Second)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to subscribe to adapter")
			os.Exit(1)
		}

	case sqx.DataTypeDepth:
		logger.Log.Error().Msg("Depth data type not supported")
		os.Exit(1)
	}

	shutdown.WaitForShutdown(syscall.SIGINT, syscall.SIGTERM)
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
func printConfiguration(exchange string, instrument string, symbol string, dataType string, natsURIs string) {
	logger.Log.Info().
		Str("exchange", exchange).
		Str("instrument", instrument).
		Str("symbol", symbol).
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
  feed <exchange> <instrument> <symbol> <data-type> <nats-uris>

Examples:
  feed binance spot BTCUSDT trade 'nats://localhost:4222?stream=feed&subject=test'
  feed binance futures ETHUSDT depth 'nats://localhost:4223?stream=feed&subject=test,nats://localhost:4224?stream=feed&subject=test'
`)
		flag.PrintDefaults()
	}

	// Parse flags
	flag.Parse()

	// Check if we have the required positional arguments
	args := flag.Args()
	if len(args) != 5 {
		logger.Log.Error().Msg("exactly 5 arguments required: <exchange> <instrument> <symbol> <data-type> <nats-uris>")
		flag.Usage()
		os.Exit(1)
	}

	// Parse positional arguments
	exchange := args[0]
	instrument := args[1]
	symbol := args[2]
	dataType := args[3]
	natsURIs := args[4]

	// Run the main logic
	runFeed(exchange, instrument, symbol, dataType, natsURIs)
}
