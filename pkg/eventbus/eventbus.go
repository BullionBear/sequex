package eventbus

import "sync"

type EventBus struct {
	subscribers sync.Map
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: sync.Map{},
	}
}

type Subscriber func(event byte) error

func (e *EventBus) Subscribe(eventType int, subscriber Subscriber) {
	if _, ok := e.subscribers.Load(eventType); !ok {
		e.subscribers.Store(eventType, []Subscriber{})
	}

	subscribers, _ := e.subscribers.Load(eventType)
	subscribers = append(subscribers.([]Subscriber), subscriber)
	e.subscribers.Store(eventType, subscribers)
}

func (e *EventBus) Publish(eventType int, event byte) {
	if subscribers, ok := e.subscribers.Load(eventType); ok {
		for _, subscriber := range subscribers.([]Subscriber) {
			subscriber(event)
		}
	}
}
