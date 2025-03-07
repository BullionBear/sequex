package eventbus

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Event types
	SpotKlineEvent = iota
	SpotTradeEvent
	SpotOrderEvent

	PerpKlineEvent
	PerpTradeEvent
	PerpOrderEvent
)

type Event struct {
	Type      int       `type:"int"`
	UUID      uuid.UUID `uuid:"uuid"`
	Timestamp time.Time `timestamp:"time"`
	Data      byte      `data:"byte"`
}

func NewEvent(eventType int, uuid uuid.UUID, timestamp time.Time, data byte) *Event {
	return &Event{
		Type:      eventType,
		UUID:      uuid,
		Timestamp: timestamp,
		Data:      data,
	}
}
