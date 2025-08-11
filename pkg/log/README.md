# Go Lightweight Structured Logger

A lightweight, thread-safe structured logger for Go applications with time-based log rotation support.

## Features

- **Structured Logging**: JSON and text format support with key-value pairs
- **Time-Based Rotation**: Automatic log file rotation with configurable intervals
- **Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Contextual Information**: Automatic inclusion of timestamp, filename, function name, and line number
- **Format String Support**: Backward compatibility with `fmt.Sprintf`-style logging
- **Thread-Safe**: Concurrent logging without data races
- **Lightweight**: Minimal dependencies, optimized for production use
- **Configurable**: Flexible output options and encoding formats

## Quick Start

```go
package main

import (
    "github.com/BullionBear/sequex/pkg/log"
)

func main() {
    // Create a basic logger
    logger := log.New(
        log.WithLevel(log.LevelInfo),
        log.WithEncoder(log.NewJSONEncoder()),
    )

    // Log with structured fields
    logger.Info("User logged in", 
        log.String("user_id", "12345"),
        log.String("ip", "192.168.1.100"),
    )
}
```

## Installation

The logger is part of the `sequex` project. Import it in your Go code:

```go
import "github.com/BullionBear/sequex/pkg/log"
```

## Basic Usage

### Creating a Logger

```go
// Basic logger with JSON output to stdout
logger := log.New(
    log.WithLevel(log.LevelDebug),
    log.WithEncoder(log.NewJSONEncoder()),
)

// Text format logger
logger := log.New(
    log.WithLevel(log.LevelInfo),
    log.WithEncoder(log.NewTextEncoder()),
)
```

### Logging Methods

```go
// Structured logging with fields
logger.Info("User action", 
    log.String("user_id", "12345"),
    log.Int("age", 30),
    log.Float64("score", 95.5),
    log.Bool("active", true),
    log.Error(fmt.Errorf("some error")),
)

// Format string support
logger.Infof("Processing %d items for user %s", 42, "alice")
logger.Errorf("Failed to connect to %s:%d", "localhost", 5432)
```

### Log Levels

```go
logger.Debug("Debug information")
logger.Info("Information message")
logger.Warn("Warning message")
logger.Error("Error message")
logger.Fatal("Fatal error - exits program")
```

## Advanced Features

### Time-Based Log Rotation

```go
logger := log.New(
    log.WithLevel(log.LevelDebug),
    log.WithEncoder(log.NewJSONEncoder()),
    log.WithTimeRotation(
        "./logs",           // Log directory
        "app.log",          // Log filename
        24*time.Hour,       // Rotate daily
        7,                  // Keep 7 backup files
    ),
)
```

### Contextual Fields

```go
// Create a logger with persistent fields
userLogger := logger.With(
    log.String("service", "user-service"),
    log.String("instance", "prod-01"),
)

// All subsequent logs will include these fields
userLogger.Info("User created", log.String("user_id", "67890"))
```

### Dynamic Configuration

```go
logger := log.New(log.WithLevel(log.LevelWarn))

// Change log level at runtime
logger.SetLevel(log.LevelDebug)

// Change output at runtime
logger.SetOutput(os.Stderr)
```

## Field Types

The logger supports various field types:

```go
log.String("key", "value")
log.Int("key", 42)
log.Int64("key", 123456789)
log.Float64("key", 3.14)
log.Bool("key", true)
log.Error(err)  // Special field for errors
log.Any("key", anyValue)  // For any type
```

## Output Formats

### JSON Format

```json
{
  "timestamp": "2023-10-01T12:00:00Z",
  "level": "INFO",
  "message": "User logged in",
  "file": "main.go",
  "function": "main",
  "line": 42,
  "user_id": "12345",
  "ip": "192.168.1.100"
}
```

### Text Format

```
2023-10-01T12:00:00Z INFO main.go:42 > User logged in user_id=12345 ip=192.168.1.100
```

## Configuration Options

### Logger Options

- `WithLevel(level)`: Set the minimum log level
- `WithOutput(writer)`: Set the output writer (stdout, stderr, file, etc.)
- `WithEncoder(encoder)`: Set the output format (JSON or text)
- `WithTimeRotation(dir, filename, interval, maxBackups)`: Enable time-based rotation
- `WithCallerSkip(skip)`: Set the number of call frames to skip for caller info

### Log Levels

- `LevelDebug`: Debug messages
- `LevelInfo`: Information messages
- `LevelWarn`: Warning messages
- `LevelError`: Error messages
- `LevelFatal`: Fatal errors (exits program)

## Examples

### Web Application Logging

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    logger := log.New(log.WithLevel(log.LevelInfo))
    
    logger.Info("Request received",
        log.String("method", r.Method),
        log.String("path", r.URL.Path),
        log.String("ip", r.RemoteAddr),
    )
    
    // Process request...
    
    duration := time.Since(start)
    logger.Info("Request completed",
        log.String("method", r.Method),
        log.String("path", r.URL.Path),
        log.Int("status_code", 200),
        log.Float64("duration_ms", float64(duration.Microseconds())/1000),
    )
}
```

### Database Operations

```go
func createUser(user *User) error {
    logger := log.New(log.WithLevel(log.LevelDebug))
    
    logger.Debug("Creating user",
        log.String("email", user.Email),
        log.String("username", user.Username),
    )
    
    err := db.Create(user).Error
    if err != nil {
        logger.Error("Failed to create user",
            log.String("email", user.Email),
            log.Error(err),
        )
        return err
    }
    
    logger.Info("User created successfully",
        log.String("user_id", user.ID),
        log.String("email", user.Email),
    )
    
    return nil
}
```

### Error Handling

```go
func processData(data []byte) error {
    logger := log.New(log.WithLevel(log.LevelError))
    
    var result map[string]interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        logger.Error("Failed to parse JSON",
            log.Int("data_size", len(data)),
            log.Error(err),
        )
        return err
    }
    
    // Process data...
    
    return nil
}
```

## Performance Considerations

- The logger is designed for high-performance applications
- JSON encoding is optimized for speed
- Thread-safe operations use minimal locking
- Time-based rotation is efficient with lazy evaluation

## Testing

Run the tests to verify functionality:

```bash
go test ./pkg/log
```

Run benchmarks to check performance:

```bash
go test -bench=. ./pkg/log
```

## Thread Safety

The logger is fully thread-safe and can be used concurrently from multiple goroutines without additional synchronization.

## Error Handling

- The logger handles encoding errors gracefully with fallback output
- Time-based rotation errors are logged but don't stop logging
- Invalid field values are handled safely

## Best Practices

1. **Use structured logging** for machine-readable logs
2. **Include contextual fields** for better debugging
3. **Set appropriate log levels** for different environments
4. **Use time-based rotation** for production applications
5. **Include error details** using the `log.Error()` field helper
6. **Use contextual loggers** for related operations
7. **Monitor log performance** in high-throughput applications

## License

This logger is part of the sequex project and follows the same license terms.
