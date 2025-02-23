package main

import (
	"fmt"
	"time"
)

// Event struct to define an event
type Event struct {
	Name string
	Data interface{}
}

type EventBus struct {
	subscribers map[string][]chan Event
}

// NewEventBus initializes an event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan Event),
	}
}

// Subscribe allows a consumer to listen for specific events
func (eb *EventBus) Subscribe(eventName string) chan Event {
	ch := make(chan Event, 1)
	eb.subscribers[eventName] = append(eb.subscribers[eventName], ch)
	return ch
}

// Publish sends an event to all subscribers
func (eb *EventBus) Publish(event Event) {
	if chans, found := eb.subscribers[event.Name]; found {
		for _, ch := range chans {
			go func(c chan Event) {
				c <- event
			}(ch)
		}
	}
}

func eventHandler(name string, ch chan Event) {
	for event := range ch {
		fmt.Printf("Handler %s received event: %s with data: %v\n", name, event.Name, event.Data)
	}
}

func main() {
	eventBus := NewEventBus()

	// Subscribers
	sub1 := eventBus.Subscribe("order_created")
	sub2 := eventBus.Subscribe("order_created")

	// Run event handlers
	go eventHandler("A", sub1)
	go eventHandler("B", sub2)

	// Publish an event
	eventBus.Publish(Event{Name: "order_created", Data: "Order ID 12345"})

	// Allow time for goroutines to process events
	time.Sleep(time.Second)
}
