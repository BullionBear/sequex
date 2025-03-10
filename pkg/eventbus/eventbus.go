package eventbus

import (
	"github.com/BullionBear/sequex/pkg/message"
)

type EventBus interface {
	Publish(topic string, msg message.Message) error
	Subscribe(topic string, handler func(message.Message)) (unsubscribe func(), err error)
}
