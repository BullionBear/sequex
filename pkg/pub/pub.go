package pub

import (
	"strconv"

	"github.com/BullionBear/sequex/pkg/eds"
)

type Publisher struct {
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Register(eventBus *eds.EventBus) {
	// Publish an event
	for i := 0; i < 10; i++ {
		// Convert i to string ID and publish
		eventBus.Publish(eds.Event{ID: strconv.Itoa(i), Name: eds.KLineEvent, Data: "BTC/USD"})
	}
}
