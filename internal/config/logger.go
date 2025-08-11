package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
)

// LoggerSingleton provides a thread-safe singleton logger instance
type LoggerSingleton struct {
	logger log.Logger
	mu     sync.RWMutex
	config LoggerConfig
}

var (
	globalLogger *LoggerSingleton
	once         sync.Once
)

// GetLogger returns the global logger instance
func GetLogger() log.Logger {
	if globalLogger == nil {
		// Return a default logger if not initialized
		return log.New(
			log.WithLevel(log.LevelInfo),
			log.WithEncoder(log.NewTextEncoder()),
		)
	}
	globalLogger.mu.RLock()
	defer globalLogger.mu.RUnlock()
	return globalLogger.logger
}

// InitializeLogger initializes the global logger singleton with the given configuration
func InitializeLogger(config LoggerConfig) error {
	var initErr error
	once.Do(func() {
		globalLogger = &LoggerSingleton{
			config: config,
		}
		initErr = globalLogger.initialize()
	})
	return initErr
}

// initialize creates the logger instance based on the configuration
func (ls *LoggerSingleton) initialize() error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// Parse log level
	level, err := parseLogLevel(ls.config.Level)
	if err != nil {
		return fmt.Errorf("invalid log level '%s': %w", ls.config.Level, err)
	}

	// Create logger options
	opts := []log.Option{
		log.WithLevel(level),
	}

	// Set encoder based on format
	switch ls.config.Format {
	case "json":
		opts = append(opts, log.WithEncoder(log.NewJSONEncoder()))
	case "text", "":
		opts = append(opts, log.WithEncoder(log.NewTextEncoder()))
	default:
		return fmt.Errorf("unsupported log format: %s", ls.config.Format)
	}

	// Add time rotation if path is specified
	if ls.config.Path != "" {
		opts = append(opts, log.WithTimeRotation("./logs", ls.config.Path, 24*time.Hour, 7))
	}

	// Create the logger
	ls.logger = log.New(opts...)

	// Log the initialization
	ls.logger.Info("Logger initialized",
		log.String("format", ls.config.Format),
		log.String("level", ls.config.Level),
		log.String("path", ls.config.Path),
	)

	return nil
}

// ReconfigureLogger allows reconfiguring the logger at runtime
func ReconfigureLogger(config LoggerConfig) error {
	if globalLogger == nil {
		return InitializeLogger(config)
	}

	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()

	// Store old config for logging
	oldConfig := globalLogger.config
	globalLogger.config = config

	// Reinitialize the logger
	if err := globalLogger.initialize(); err != nil {
		// Restore old config on error
		globalLogger.config = oldConfig
		return err
	}

	globalLogger.logger.Info("Logger reconfigured",
		log.String("old_format", oldConfig.Format),
		log.String("new_format", config.Format),
		log.String("old_level", oldConfig.Level),
		log.String("new_level", config.Level),
	)

	return nil
}

// GetConfig returns the current logger configuration
func GetLoggerConfig() LoggerConfig {
	if globalLogger == nil {
		return LoggerConfig{}
	}
	globalLogger.mu.RLock()
	defer globalLogger.mu.RUnlock()
	return globalLogger.config
}

// parseLogLevel converts string level to log.Level
func parseLogLevel(level string) (log.Level, error) {
	switch level {
	case "debug", "DEBUG":
		return log.LevelDebug, nil
	case "info", "INFO":
		return log.LevelInfo, nil
	case "warn", "WARN":
		return log.LevelWarn, nil
	case "error", "ERROR":
		return log.LevelError, nil
	case "fatal", "FATAL":
		return log.LevelFatal, nil
	case "":
		return log.LevelInfo, nil // Default to INFO
	default:
		return log.LevelInfo, fmt.Errorf("unknown log level: %s", level)
	}
}

// Convenience functions for common logging operations

// Debug logs a debug message
func Debug(msg string, fields ...log.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...log.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...log.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...log.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...log.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// With creates a new logger with additional fields
func With(fields ...log.Field) log.Logger {
	return GetLogger().With(fields...)
}
