package strategy

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ExampleEvent represents a simple event implementation
type ExampleEvent struct {
	eventType EventType
	data      interface{}
}

func (e *ExampleEvent) Type() EventType {
	return e.eventType
}

func (e *ExampleEvent) Data() interface{} {
	return e.data
}

// ExampleStrategy demonstrates how to implement a custom strategy
type ExampleStrategy struct {
	*BaseStrategy
	counter int
}

// NewExampleStrategy creates a new example strategy
func NewExampleStrategy() *ExampleStrategy {
	return &ExampleStrategy{
		BaseStrategy: NewBaseStrategy(),
		counter:      0,
	}
}

// Initialize overrides the base initialization to add custom logic
func (s *ExampleStrategy) Initialize(ctx context.Context, config StrategyConfig) error {
	// Call base initialization
	if err := s.BaseStrategy.Initialize(ctx, config); err != nil {
		return err
	}

	// Register default event handlers
	s.registerDefaultHandlers()

	log.Printf("Example strategy '%s' initialized", config.Name)
	return nil
}

// Run overrides the base run method to add custom logic
func (s *ExampleStrategy) Run(ctx context.Context) error {
	if err := s.BaseStrategy.Run(ctx); err != nil {
		return err
	}

	log.Printf("Example strategy '%s' started", s.GetConfig().Name)

	// Start a background goroutine for periodic tasks
	go s.backgroundTask(ctx)

	return nil
}

// Shutdown overrides the base shutdown method to add cleanup logic
func (s *ExampleStrategy) Shutdown(ctx context.Context) error {
	log.Printf("Shutting down example strategy '%s'", s.GetConfig().Name)
	return s.BaseStrategy.Shutdown(ctx)
}

// registerDefaultHandlers registers the default event handlers
func (s *ExampleStrategy) registerDefaultHandlers() {
	// Register a handler for "tick" events
	s.RegisterEventHandler("tick", s.handleTickEvent)

	// Register a handler for "trade" events
	s.RegisterEventHandler("trade", s.handleTradeEvent)

	// Register a handler for "signal" events
	s.RegisterEventHandler("signal", s.handleSignalEvent)
}

// handleTickEvent processes tick events
func (s *ExampleStrategy) handleTickEvent(ctx context.Context, event Event) error {
	s.counter++
	log.Printf("Processing tick event #%d: %v", s.counter, event.Data())
	return nil
}

// handleTradeEvent processes trade events
func (s *ExampleStrategy) handleTradeEvent(ctx context.Context, event Event) error {
	log.Printf("Processing trade event: %v", event.Data())
	return nil
}

// handleSignalEvent processes signal events
func (s *ExampleStrategy) handleSignalEvent(ctx context.Context, event Event) error {
	log.Printf("Processing signal event: %v", event.Data())
	return nil
}

// backgroundTask runs periodic background tasks
func (s *ExampleStrategy) backgroundTask(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Create and process a tick event
			tickEvent := &ExampleEvent{
				eventType: "tick",
				data:      fmt.Sprintf("Tick at %s", time.Now().Format(time.RFC3339)),
			}

			if err := s.ProcessEvent(ctx, tickEvent); err != nil {
				log.Printf("Error processing tick event: %v", err)
			}
		}
	}
}

// Example usage function
func ExampleUsage() {
	// Create a new strategy
	strategy := NewExampleStrategy()

	// Create configuration
	config := StrategyConfig{
		Name:        "Example Strategy",
		Description: "A simple example strategy",
		Parameters: map[string]interface{}{
			"interval": 5,
			"enabled":  true,
		},
		Enabled: true,
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the strategy
	if err := strategy.Initialize(ctx, config); err != nil {
		log.Fatalf("Failed to initialize strategy: %v", err)
	}

	// Register additional custom handler
	strategy.RegisterEventHandler("custom", func(ctx context.Context, event Event) error {
		log.Printf("Custom event handler: %v", event.Data())
		return nil
	})

	// Run the strategy
	if err := strategy.Run(ctx); err != nil {
		log.Fatalf("Failed to run strategy: %v", err)
	}

	// Simulate some events
	events := []Event{
		&ExampleEvent{eventType: "trade", data: "Buy 100 BTC at 50000"},
		&ExampleEvent{eventType: "signal", data: "RSI oversold"},
		&ExampleEvent{eventType: "custom", data: "Custom event data"},
	}

	for _, event := range events {
		if err := strategy.ProcessEvent(ctx, event); err != nil {
			log.Printf("Error processing event: %v", err)
		}
	}

	// Let it run for a bit
	time.Sleep(10 * time.Second)

	// Shutdown the strategy
	if err := strategy.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down strategy: %v", err)
	}

	// Print final stats
	stats := strategy.GetStats()
	log.Printf("Final stats: Events processed: %d, Errors: %d",
		stats.EventsProcessed, stats.Errors)
}

