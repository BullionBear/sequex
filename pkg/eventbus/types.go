package eventbus

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var (
	_ Event[KlineData]
	_ Event[ExecuteData]
)

// Event is a generic struct to hold different types of event data
type Event[T any] struct {
	ID   uuid.UUID
	Data T
}

// KlineData represents data for Kline events
type KlineData struct {
	Open   float64
	Close  float64
	Volume float64
}

// ExecuteData represents data for execution events
type ExecuteData struct {
	Price float64
	Size  float64
}

// EventType is used to distinguish different types of events
type EventType string

const (
	KlineEvent   EventType = "kline_event"
	ExecuteEvent EventType = "execute_event"
)

// EventBus manages subscriptions and publishing of events
type EventBus struct {
	handlers map[EventType][]interface{}
	mu       sync.RWMutex
}

// NewEventBus creates a new EventBus instance
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]interface{}),
	}
}

// Subscribe registers a handler function for a specific event type
func Subscribe[T any](eb *EventBus, eventType EventType, handler func(Event[T])) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Unsubscribe removes a handler function for a specific event type
func Unsubscribe[T any](eb *EventBus, eventType EventType, handler func(Event[T])) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	handlers, ok := eb.handlers[eventType]
	if !ok {
		return
	}

	// Find and remove the handler
	for i, h := range handlers {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// Publish sends an event to all subscribed handlers for the event's type
func Publish[T any](eb *EventBus, eventType EventType, event Event[T]) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers, ok := eb.handlers[eventType]
	if !ok {
		return
	}

	// Execute each handler in a separate goroutine
	for _, h := range handlers {
		if handler, ok := h.(func(Event[T])); ok {
			go handler(event)
		}
	}
}
