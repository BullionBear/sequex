package inprocbus

import (
	"sync"

	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/message"
)

// Ensure InprocQueue implements the mq.MessageQueue interface.
var _ eventbus.EventBus = (*InprocBus)(nil)

// InprocBus is a basic implementation of the EventBus interface using sync.Map.
type InprocBus struct {
	subscribers sync.Map // key: topic, value: []func(Message)
}

// NewEventBus creates and returns a new InprocBus instance.
func NewEventBus() *InprocBus {
	return &InprocBus{}
}

// Publish sends a message to all subscribers of a given topic.
func (bus *InprocBus) Publish(topic string, msg message.Message) error {
	if handlers, ok := bus.subscribers.Load(topic); ok {
		for _, handler := range handlers.([]func(message.Message)) {
			go handler(msg) // Send message asynchronously to each handler.
		}
	}
	return nil
}

// Subscribe adds a new handler for a specific topic and returns an unsubscribe function.
func (bus *InprocBus) Subscribe(topic string, handler func(message.Message)) (func(), error) {
	var newHandlers []func(message.Message)

	actual, _ := bus.subscribers.LoadOrStore(topic, []func(message.Message){handler})
	handlers := actual.([]func(message.Message))

	newHandlers = append(handlers, handler)
	bus.subscribers.Store(topic, newHandlers)

	unsubscribe := func() {
		if actual, ok := bus.subscribers.Load(topic); ok {
			handlers := actual.([]func(message.Message))
			var updatedHandlers []func(message.Message)
			for _, h := range handlers {
				if &h != &handler {
					updatedHandlers = append(updatedHandlers, h)
				}
			}
			if len(updatedHandlers) > 0 {
				bus.subscribers.Store(topic, updatedHandlers)
			} else {
				bus.subscribers.Delete(topic)
			}
		}
	}

	return unsubscribe, nil
}
