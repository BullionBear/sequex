package main

import (
	"fmt"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
)

func main() {
	// Example 1: Basic logger with JSON output to stdout
	fmt.Println("=== Example 1: Basic JSON Logger ===")
	basicLogger := log.New(
		log.WithLevel(log.LevelDebug),
		log.WithEncoder(log.NewJSONEncoder()),
	)

	basicLogger.Info("Application started",
		log.String("version", "1.0.0"),
		log.String("environment", "development"),
	)
	basicLogger.Debug("Debug information", log.Int("debug_level", 5))
	basicLogger.Warn("Warning message", log.String("component", "database"))
	basicLogger.Error("Error occurred", log.Error(fmt.Errorf("connection timeout")))

	// Example 2: Text format logger
	fmt.Println("\n=== Example 2: Text Format Logger ===")
	textLogger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
	)

	textLogger.Info("User action",
		log.String("user_id", "12345"),
		log.String("action", "login"),
		log.String("ip", "192.168.1.100"),
	)
	textLogger.Warn("Performance warning",
		log.Float64("response_time", 2.5),
		log.String("endpoint", "/api/users"),
	)

	// Example 3: Logger with time-based rotation
	fmt.Println("\n=== Example 3: Time-Based Rotation Logger ===")
	rotationLogger := log.New(
		log.WithLevel(log.LevelDebug),
		log.WithEncoder(log.NewJSONEncoder()),
		log.WithTimeRotation(
			"./logs",     // Log directory
			"app.log",    // Log filename
			24*time.Hour, // Rotate daily
			7,            // Keep 7 backup files
		),
	)

	rotationLogger.Info("Logging with rotation enabled",
		log.String("feature", "time_rotation"),
		log.Int("max_backups", 7),
	)

	// Example 4: Logger with contextual fields
	fmt.Println("\n=== Example 4: Contextual Fields ===")
	userLogger := log.New(
		log.WithLevel(log.LevelDebug),
		log.WithEncoder(log.NewJSONEncoder()),
	).With(
		log.String("service", "user-service"),
		log.String("instance", "prod-01"),
	)

	userLogger.Info("User created",
		log.String("user_id", "67890"),
		log.String("email", "user@example.com"),
	)
	userLogger.Error("Failed to send email",
		log.String("user_id", "67890"),
		log.Error(fmt.Errorf("SMTP server unavailable")),
	)

	// Example 5: Format string support
	fmt.Println("\n=== Example 5: Format String Support ===")
	formatLogger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
	)

	formatLogger.Infof("Processing %d items for user %s", 42, "alice")
	formatLogger.Warnf("Database query took %.2f seconds", 1.23)
	formatLogger.Errorf("Failed to connect to %s:%d", "localhost", 5432)

	// Example 6: Dynamic level and output changes
	fmt.Println("\n=== Example 6: Dynamic Configuration ===")
	dynamicLogger := log.New(
		log.WithLevel(log.LevelWarn), // Start with WARN level
		log.WithEncoder(log.NewTextEncoder()),
	)

	dynamicLogger.Debug("This debug message won't appear")
	dynamicLogger.Info("This info message won't appear")
	dynamicLogger.Warn("This warning will appear")

	// Change level dynamically
	dynamicLogger.SetLevel(log.LevelDebug)
	dynamicLogger.Debug("Now debug messages will appear")

	// Example 7: Error handling with structured logging
	fmt.Println("\n=== Example 7: Error Handling ===")
	errorLogger := log.New(
		log.WithLevel(log.LevelError),
		log.WithEncoder(log.NewJSONEncoder()),
	)

	// Simulate different types of errors
	errors := []error{
		fmt.Errorf("database connection failed"),
		fmt.Errorf("invalid input parameter"),
		fmt.Errorf("network timeout"),
	}

	for i, err := range errors {
		errorLogger.Error("Operation failed",
			log.Int("attempt", i+1),
			log.Error(err),
			log.String("operation", "data_processing"),
			log.Bool("retry", i < len(errors)-1),
		)
	}

	// Example 8: Performance monitoring
	fmt.Println("\n=== Example 8: Performance Monitoring ===")
	perfLogger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewJSONEncoder()),
	)

	start := time.Now()
	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	duration := time.Since(start)

	perfLogger.Info("Request completed",
		log.String("endpoint", "/api/data"),
		log.String("method", "GET"),
		log.Int("status_code", 200),
		log.Float64("duration_ms", float64(duration.Microseconds())/1000),
		log.Int("response_size", 1024),
	)

	// Example 9: Business logic logging
	fmt.Println("\n=== Example 9: Business Logic Logging ===")
	businessLogger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
	).With(
		log.String("module", "order_processing"),
	)

	// Simulate order processing
	orderID := "ORD-12345"
	businessLogger.Info("Order received",
		log.String("order_id", orderID),
		log.Float64("amount", 99.99),
		log.String("currency", "USD"),
	)

	businessLogger.Info("Payment processed",
		log.String("order_id", orderID),
		log.String("payment_method", "credit_card"),
		log.String("transaction_id", "TXN-67890"),
	)

	businessLogger.Info("Order shipped",
		log.String("order_id", orderID),
		log.String("tracking_number", "TRK-11111"),
		log.String("carrier", "FedEx"),
	)

	// Example 10: Fatal logging (simulated)
	fmt.Println("\n=== Example 10: Fatal Logging (simulated) ===")
	// fatalLogger := log.New(
	// 	log.WithLevel(log.LevelFatal),
	// 	log.WithEncoder(log.NewTextEncoder()),
	// )

	// Uncomment the next line to see fatal logging in action
	// fatalLogger.Fatal("Critical system failure", log.String("component", "database"))

	fmt.Println("Fatal logging would exit the program here")
	fmt.Println("All examples completed successfully!")
}
