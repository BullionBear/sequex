package main

import (
	"fmt"
	"time"

	"github.com/BullionBear/sequex/pkg/eds"
	"github.com/BullionBear/sequex/pkg/pub"
)

func eventHandler(name string, ch chan eds.Event) {
	for event := range ch {
		fmt.Printf("Handler %s received event: %s with data: %v\n", name, event.Name, event.Data)
	}
}

func main() {
	eventBus := eds.NewEventBus()

	// Subscribers
	sub1 := eventBus.Subscribe(eds.KLineEvent)
	sub2 := eventBus.Subscribe(eds.KLineEvent)
	// Run event handlers
	go eventHandler("A", sub1)
	go eventHandler("B", sub2)

	// Publisher
	publisher := pub.NewPublisher()
	publisher.Register(eventBus)

	// Publish an event
	// eventBus.Publish(eds.Event{ID: "123", Name: eds.KLineEvent, Data: "BTC/USD"})

	// Allow time for goroutines to process events
	time.Sleep(time.Second)
}
