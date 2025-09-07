package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Package-level variable that holds our configured logger instance.
// It starts with a disabled logger to be safe until it's initialized.
var Log zerolog.Logger = zerolog.New(nil).Level(zerolog.Disabled)

// InitLogger initializes the global logger with the desired configuration.
// This function should be called once, from main().
func InitLogger(isDevelopment bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro // Use microsecond precision
	zerolog.SetGlobalLevel(zerolog.InfoLevel)             // Set default global level

	// Human-friendly output for local development
	outputWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05.000000", // Microsecond precision
	}
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // More verbose in dev

	// For development, use the console writer
	Log = zerolog.New(outputWriter).
		With().
		Timestamp().
		Caller().
		Logger()
}

// Get returns the global logger instance.
// This is useful if you need to pass the logger to other libraries that don't use this package directly.
func Get() *zerolog.Logger {
	return &Log
}
