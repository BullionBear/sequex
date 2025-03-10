package mq

import "github.com/BullionBear/sequex/pkg/message"

type MessageQueue interface {
	Publish(topic string, msg message.Message) error
	Subscribe(topic string) (<-chan message.Message, error)
	Unsubscribe(topic string, ch <-chan message.Message) error
}
