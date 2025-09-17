package main

import (
	"flag"
	"os"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/internal/adapter"
	_ "github.com/BullionBear/sequex/internal/adapter/init"
	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/internal/model/sqx"
	"github.com/BullionBear/sequex/pkg/logger"
	"github.com/BullionBear/sequex/pkg/shutdown"
	"github.com/nats-io/nats.go"
)

// runFeed executes the main feed logic
func runFeed(configFile string) {
	// Output version information
	logger.Log.Info().
		Str("version", env.Version).
		Str("buildTime", env.BuildTime).
		Str("commitHash", env.CommitHash).
		Msg("Feed started")

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to load config")
		os.Exit(1)
	}

	printConfiguration(cfg)
	sqxExchange := sqx.NewExchange(cfg.Exchange)
	if sqxExchange == sqx.ExchangeUnknown {
		logger.Log.Error().Msg("Invalid exchange")
		os.Exit(1)
	}

	sqxInstrumentType := sqx.NewInstrumentType(cfg.Instrument)
	if sqxInstrumentType == sqx.InstrumentTypeUnknown {
		logger.Log.Error().Msg("Invalid instrument")
		os.Exit(1)
	}

	sqxSymbol, err := sqx.NewSymbolFromStr(cfg.Symbol)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create symbol")
		os.Exit(1)
	}

	sqxDataType := sqx.NewDataType(cfg.Type)
	if sqxDataType == sqx.DataTypeUnknown {
		logger.Log.Error().Msg("Invalid data type")
		os.Exit(1)
	}

	shutdown := shutdown.NewShutdown(logger.Log)

	natsConn, err := nats.Connect(cfg.NATS.URIs)
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
	streamInfo, err := js.StreamInfo(cfg.NATS.Stream)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get stream info")
		os.Exit(1)
	}
	logger.Log.Info().Msgf("Stream info: %+v", streamInfo)
	subject := cfg.NATS.Subject
	switch sqxDataType {
	case sqx.DataTypeTrade:
		adapter, err := adapter.CreateTradeAdapter(sqxExchange)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to create adapter")
			os.Exit(1)
		}
		unsubscribe, err := adapter.Subscribe(sqxSymbol, sqxInstrumentType, func(trade sqx.Trade) error {
			data, err := trade.Marshal()
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to marshal trade")
				return err
			}
			header := nats.Header{
				"Nats-Msg-Id": []string{trade.IdStr()},
			}

			_, err = js.PublishMsg(&nats.Msg{
				Subject: subject,
				Data:    data,
				Header:  header,
			})
			return err
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

// printConfiguration prints the parsed configuration
func printConfiguration(cfg *config.Config) {
	logger.Log.Info().
		Str("exchange", cfg.Exchange).
		Str("instrument", cfg.Instrument).
		Str("symbol", cfg.Symbol).
		Str("dataType", cfg.Type).
		Str("natsURIs", cfg.NATS.URIs).
		Str("stream", cfg.NATS.Stream).
		Str("subject", cfg.NATS.Subject).
		Msg("Feed Configuration")
}

func main() {
	// Define flags
	var configFile string
	flag.StringVar(&configFile, "c", "", "Configuration file path (required)")

	// Custom usage function
	flag.Usage = func() {
		logger.Log.Info().Msg(`Feed is a scalable CLI tool for streaming market data from various exchanges
to NATS message brokers. It supports multiple exchanges and data types.

Usage:
  feed -c <config-file>

Examples:
  feed -c config/trade-binance-spot-btcusdt.json
`)
		flag.PrintDefaults()
	}

	// Parse flags
	flag.Parse()

	// Check if the required config file flag is provided
	if configFile == "" {
		logger.Log.Error().Msg("config file path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Run the main logic
	runFeed(configFile)
}
