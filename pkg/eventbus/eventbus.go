package eventbus

import (
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type EventBus struct {
	nc *nats.Conn
}

// MessageFactory is a function that creates a new protobuf message instance
type MessageFactory func() proto.Message

func NewEventBus(nc *nats.Conn) *EventBus {
	return &EventBus{nc: nc}
}

func (e *EventBus) Emit(topic string, event proto.Message) error {
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	return e.nc.Publish(topic, data)
}

func (e *EventBus) On(topic string, messageFactory MessageFactory, callback func(event proto.Message)) error {
	_, err := e.nc.Subscribe(topic, func(m *nats.Msg) {
		event := messageFactory()
		err := proto.Unmarshal(m.Data, event)
		if err != nil {
			return
		}
		callback(event)
	})
	return err
}

func (e *EventBus) RegisterRPC(topic string, messageFactory MessageFactory, callback func(event proto.Message) proto.Message) error {
	_, err := e.nc.Subscribe(topic, func(m *nats.Msg) {
		event := messageFactory()
		err := proto.Unmarshal(m.Data, event)
		if err != nil {
			return
		}
		response := callback(event)
		data, err := proto.Marshal(response)
		if err != nil {
			return
		}
		m.Respond(data)
	})
	return err
}

func (e *EventBus) CallRPC(topic string, event proto.Message, responseFactory MessageFactory, timeout time.Duration) (proto.Message, error) {
	data, err := proto.Marshal(event)
	if err != nil {
		return nil, err
	}
	msg, err := e.nc.Request(topic, data, timeout)
	if err != nil {
		return nil, err
	}
	response := responseFactory()
	err = proto.Unmarshal(msg.Data, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
