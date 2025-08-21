package strategy

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// EventType represents the type of event that can be handled by a strategy
type EventType string

// Event represents a generic event that can be processed by strategy handlers
type Event interface {
	Type() EventType
	Data() interface{}
}

// EventHandler is a function that processes events of a specific type
type EventHandler func(ctx context.Context, event Event) error

// StrategyState represents the current state of a strategy
type StrategyState string

const (
	StateInitialized StrategyState = "initialized"
	StateRunning     StrategyState = "running"
	StateStopped     StrategyState = "stopped"
	StateError       StrategyState = "error"
)

// StrategyConfig holds configuration for a strategy
type StrategyConfig struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Enabled     bool                   `json:"enabled"`
}

// StrategyStats holds runtime statistics for a strategy
type StrategyStats struct {
	EventsProcessed int64  `json:"events_processed"`
	Errors          int64  `json:"errors"`
	LastEventTime   string `json:"last_event_time"`
	Uptime          string `json:"uptime"`
}

// Strategy is the main interface that all strategies must implement
type Strategy interface {
	// Lifecycle methods
	Initialize(ctx context.Context, config StrategyConfig) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error

	// State management
	GetState() StrategyState
	GetStats() StrategyStats

	// Event handling
	RegisterEventHandler(eventType EventType, handler EventHandler) error
	UnregisterEventHandler(eventType EventType) error
	ProcessEvent(ctx context.Context, event Event) error

	// Serialization
	Serialize() ([]byte, error)
	Deserialize(data []byte) error

	// Configuration
	GetConfig() StrategyConfig
	UpdateConfig(config StrategyConfig) error
}

// BaseStrategy provides a default implementation of common strategy functionality
type BaseStrategy struct {
	mu          sync.RWMutex
	config      StrategyConfig
	state       StrategyState
	stats       StrategyStats
	handlers    map[EventType]EventHandler
	ctx         context.Context
	cancel      context.CancelFunc
	initialized bool
}

// NewBaseStrategy creates a new base strategy instance
func NewBaseStrategy() *BaseStrategy {
	return &BaseStrategy{
		handlers: make(map[EventType]EventHandler),
		state:    StateInitialized,
		stats:    StrategyStats{},
	}
}

// Initialize sets up the strategy with the given configuration
func (s *BaseStrategy) Initialize(ctx context.Context, config StrategyConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.initialized {
		return fmt.Errorf("strategy already initialized")
	}

	s.config = config
	s.state = StateInitialized
	s.initialized = true

	return nil
}

// Run starts the strategy execution
func (s *BaseStrategy) Run(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.initialized {
		return fmt.Errorf("strategy not initialized")
	}

	if s.state == StateRunning {
		return fmt.Errorf("strategy already running")
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.state = StateRunning

	return nil
}

// Shutdown gracefully stops the strategy
func (s *BaseStrategy) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateRunning {
		return nil
	}

	if s.cancel != nil {
		s.cancel()
	}

	s.state = StateStopped
	return nil
}

// GetState returns the current state of the strategy
func (s *BaseStrategy) GetState() StrategyState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// GetStats returns the current statistics of the strategy
func (s *BaseStrategy) GetStats() StrategyStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stats
}

// RegisterEventHandler registers a handler for a specific event type
func (s *BaseStrategy) RegisterEventHandler(eventType EventType, handler EventHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	s.handlers[eventType] = handler
	return nil
}

// UnregisterEventHandler removes a handler for a specific event type
func (s *BaseStrategy) UnregisterEventHandler(eventType EventType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.handlers[eventType]; !exists {
		return fmt.Errorf("no handler registered for event type: %s", eventType)
	}

	delete(s.handlers, eventType)
	return nil
}

// ProcessEvent processes an incoming event using the registered handler
func (s *BaseStrategy) ProcessEvent(ctx context.Context, event Event) error {
	s.mu.RLock()
	handler, exists := s.handlers[event.Type()]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no handler registered for event type: %s", event.Type())
	}

	// Update stats
	s.mu.Lock()
	s.stats.EventsProcessed++
	s.mu.Unlock()

	// Process the event
	if err := handler(ctx, event); err != nil {
		s.mu.Lock()
		s.stats.Errors++
		s.mu.Unlock()
		return fmt.Errorf("error processing event: %w", err)
	}

	return nil
}

// Serialize converts the strategy state to JSON
func (s *BaseStrategy) Serialize() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := struct {
		Config StrategyConfig `json:"config"`
		State  StrategyState  `json:"state"`
		Stats  StrategyStats  `json:"stats"`
	}{
		Config: s.config,
		State:  s.state,
		Stats:  s.stats,
	}

	return json.Marshal(data)
}

// Deserialize loads the strategy state from JSON
func (s *BaseStrategy) Deserialize(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var serialized struct {
		Config StrategyConfig `json:"config"`
		State  StrategyState  `json:"state"`
		Stats  StrategyStats  `json:"stats"`
	}

	if err := json.Unmarshal(data, &serialized); err != nil {
		return fmt.Errorf("failed to deserialize strategy data: %w", err)
	}

	s.config = serialized.Config
	s.state = serialized.State
	s.stats = serialized.Stats
	s.initialized = true

	return nil
}

// GetConfig returns the current configuration
func (s *BaseStrategy) GetConfig() StrategyConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// UpdateConfig updates the strategy configuration
func (s *BaseStrategy) UpdateConfig(config StrategyConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state == StateRunning {
		return fmt.Errorf("cannot update config while strategy is running")
	}

	s.config = config
	return nil
}

// GetContext returns the strategy's context
func (s *BaseStrategy) GetContext() context.Context {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ctx
}

// IsRunning checks if the strategy is currently running
func (s *BaseStrategy) IsRunning() bool {
	return s.GetState() == StateRunning
}

// IsInitialized checks if the strategy has been initialized
func (s *BaseStrategy) IsInitialized() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.initialized
}
