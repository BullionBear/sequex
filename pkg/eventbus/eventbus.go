package eventbus

import (
	"fmt"
	"sync"
)

// EventBus is a struct that manages event subscriptions and publishing
type EventBus struct {
	handlers map[EventType][]interface{} // Handlers for different event types
	mu       sync.RWMutex                // Mutex for safe concurrent access
}

// NewEventBus creates a new EventBus instance
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]interface{}),
	}
}

// Subscribe registers a handler function for a specific event type
func (eb *EventBus) Subscribe(eventType EventType, handler interface{}) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// Store the handler in the handlers map
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Unsubscribe removes a handler function for a specific event type
func (eb *EventBus) Unsubscribe(eventType EventType, handler interface{}) {
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
func (eb *EventBus) Publish(event Event[interface{}]) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers, ok := eb.handlers[event.Type]
	if !ok {
		return
	}

	// Execute each handler in a separate goroutine
	for _, h := range handlers {
		if handler, ok := h.(func(Event[interface{}])); ok {
			go handler(event)
		}
	}
}
