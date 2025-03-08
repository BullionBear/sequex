package inprocq

import "github.com/BullionBear/sequex/pkg/mq"

var _ mq.MessageQueue = (*InprocQueue)(nil)

type InprocQueue struct {
}

func NewInprocQueue() *InprocQueue {
	return &InprocQueue{}
}

func (q *InprocQueue) Publish(topic string, msg mq.Message) error {
	return nil
}

func (q *InprocQueue) Subscribe(topic string) (<-chan mq.Message, error) {
	return nil, nil
}

func (q *InprocQueue) Unsubscribe(topic string, ch <-chan mq.Message) error {
	return nil
}
