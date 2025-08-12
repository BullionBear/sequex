package config

import (
	"io"
	"os"

	"github.com/BullionBear/sequex/pkg/log"
)

func CreateLogger(loggerConfig LoggerConfig) (log.Logger, error) {
	level, err := parseLevel(loggerConfig.Level)
	if err != nil {
		return nil, err
	}
	encoder, err := parseFormat(loggerConfig.Format)
	if err != nil {
		return nil, err
	}
	output, err := parseWriter(loggerConfig.Path)
	if err != nil {
		return nil, err
	}

	return log.New(
		log.WithLevel(level),
		log.WithEncoder(encoder),
		log.WithOutput(output),
		log.WithCallerSkip(2),
	), nil
}

func parseLevel(level string) (log.Level, error) {
	switch level {
	case "debug", "Debug", "DEBUG":
		return log.LevelDebug, nil
	case "info", "Info", "INFO":
		return log.LevelInfo, nil
	case "warn", "Warn", "WARN":
		return log.LevelWarn, nil
	case "error", "Error", "ERROR":
		return log.LevelError, nil
	case "fatal", "Fatal", "FATAL":
		return log.LevelFatal, nil
	}
	return log.LevelInfo, nil
}

func parseFormat(format string) (log.Encoder, error) {
	switch format {
	case "text":
		return log.NewTextEncoder(), nil
	case "json":
		return log.NewJSONEncoder(), nil
	}
	return log.NewTextEncoder(), nil
}

func parseWriter(path string) (io.Writer, error) {
	if path == "" {
		return os.Stdout, nil
	}
	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}
