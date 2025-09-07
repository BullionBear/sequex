package log

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Level represents the logging level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// Logger interface that matches the expected usage
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Infof(format string, args ...interface{})
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
}

// Field represents a log field
type Field interface{}

// loggerImpl implements the Logger interface using zerolog
type loggerImpl struct {
	logger zerolog.Logger
	fields []Field
}

// New creates a new logger with the given options
func New(opts ...Option) Logger {
	config := &config{
		level:   LevelInfo,
		output:  os.Stdout,
		encoder: &textEncoder{},
	}

	for _, opt := range opts {
		opt(config)
	}

	// Convert our level to zerolog level
	var zLevel zerolog.Level
	switch config.level {
	case LevelDebug:
		zLevel = zerolog.DebugLevel
	case LevelInfo:
		zLevel = zerolog.InfoLevel
	case LevelWarn:
		zLevel = zerolog.WarnLevel
	case LevelError:
		zLevel = zerolog.ErrorLevel
	case LevelFatal:
		zLevel = zerolog.FatalLevel
	default:
		zLevel = zerolog.InfoLevel
	}

	// Create zerolog logger
	var zLogger zerolog.Logger
	if config.encoder != nil {
		// Use console writer for text encoding
		zLogger = zerolog.New(zerolog.ConsoleWriter{
			Out:        config.output,
			TimeFormat: "15:04:05.000000",
		}).Level(zLevel).With().Timestamp().Caller().Logger()
	} else {
		// Use JSON output
		zLogger = zerolog.New(config.output).Level(zLevel).With().Timestamp().Caller().Logger()
	}

	return &loggerImpl{logger: zLogger, fields: make([]Field, 0)}
}

// DefaultLogger creates a default logger with standard configuration
func DefaultLogger(opts ...Option) Logger {
	return New(opts...)
}

// Debug logs a debug message
func (l *loggerImpl) Debug(msg string, fields ...Field) {
	event := l.logger.Debug()

	// Apply stored fields first
	for _, field := range l.fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}

	// Then apply new fields
	for _, field := range fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}
	event.Msg(msg)
}

// Info logs an info message
func (l *loggerImpl) Info(msg string, fields ...Field) {
	event := l.logger.Info()

	// Apply stored fields first
	for _, field := range l.fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}

	// Then apply new fields
	for _, field := range fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}
	event.Msg(msg)
}

// Infof logs a formatted info message
func (l *loggerImpl) Infof(format string, args ...interface{}) {
	event := l.logger.Info()

	// Apply stored fields first
	for _, field := range l.fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}

	event.Msgf(format, args...)
}

// Warn logs a warning message
func (l *loggerImpl) Warn(msg string, fields ...Field) {
	event := l.logger.Warn()

	// Apply stored fields first
	for _, field := range l.fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}

	// Then apply new fields
	for _, field := range fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}
	event.Msg(msg)
}

// Error logs an error message
func (l *loggerImpl) Error(msg string, fields ...Field) {
	event := l.logger.Error()

	// Apply stored fields first
	for _, field := range l.fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}

	// Then apply new fields
	for _, field := range fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *loggerImpl) Fatal(msg string, fields ...Field) {
	event := l.logger.Fatal()

	// Apply stored fields first
	for _, field := range l.fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}

	// Then apply new fields
	for _, field := range fields {
		if f, ok := field.(func(*zerolog.Event) *zerolog.Event); ok {
			event = f(event)
		}
	}
	event.Msg(msg)
}

// With creates a new logger with additional context fields
func (l *loggerImpl) With(fields ...Field) Logger {
	// Combine existing fields with new fields
	combinedFields := make([]Field, 0, len(l.fields)+len(fields))
	combinedFields = append(combinedFields, l.fields...)
	combinedFields = append(combinedFields, fields...)

	// Return a new logger with the combined fields
	return &loggerImpl{
		logger: l.logger,
		fields: combinedFields,
	}
}

// config holds the logger configuration
type config struct {
	level   Level
	output  io.Writer
	encoder Encoder
}

// Option is a function that configures the logger
type Option func(*config)

// WithLevel sets the logging level
func WithLevel(level Level) Option {
	return func(c *config) {
		c.level = level
	}
}

// WithOutput sets the output writer
func WithOutput(w io.Writer) Option {
	return func(c *config) {
		c.output = w
	}
}

// WithEncoder sets the encoder
func WithEncoder(encoder Encoder) Option {
	return func(c *config) {
		c.encoder = encoder
	}
}

// WithTimeRotation sets up time-based log rotation (placeholder for compatibility)
func WithTimeRotation(dir, filename string, rotationTime time.Duration, maxAge int) Option {
	// For now, just return a no-op option since we're using zerolog
	// In a real implementation, you might want to use lumberjack or similar
	return func(c *config) {
		// No-op for now
	}
}

// Encoder interface for different output formats
type Encoder interface{}

// textEncoder represents a text encoder
type textEncoder struct{}

// NewTextEncoder creates a new text encoder
func NewTextEncoder() Encoder {
	return &textEncoder{}
}

// Helper functions for creating fields

// String creates a string field
func String(key, value string) func(*zerolog.Event) *zerolog.Event {
	return func(e *zerolog.Event) *zerolog.Event {
		return e.Str(key, value)
	}
}

// Error creates an error field
func Error(err error) func(*zerolog.Event) *zerolog.Event {
	return func(e *zerolog.Event) *zerolog.Event {
		return e.Err(err)
	}
}

// Int creates an integer field
func Int(key string, value int) func(*zerolog.Event) *zerolog.Event {
	return func(e *zerolog.Event) *zerolog.Event {
		return e.Int(key, value)
	}
}

// Bool creates a boolean field
func Bool(key string, value bool) func(*zerolog.Event) *zerolog.Event {
	return func(e *zerolog.Event) *zerolog.Event {
		return e.Bool(key, value)
	}
}

// Float64 creates a float64 field
func Float64(key string, value float64) func(*zerolog.Event) *zerolog.Event {
	return func(e *zerolog.Event) *zerolog.Event {
		return e.Float64(key, value)
	}
}

// Int64 creates an int64 field
func Int64(key string, value int64) func(*zerolog.Event) *zerolog.Event {
	return func(e *zerolog.Event) *zerolog.Event {
		return e.Int64(key, value)
	}
}
