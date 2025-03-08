package mq

type Message struct {
	ID        string
	CreatedAt int64
	Content   interface{}
}

type MessageQueue interface {
	Publish(topic string, msg Message) error
	Subscribe(topic string) (<-chan Message, error)
	Unsubscribe(topic string, ch <-chan Message) error
}
