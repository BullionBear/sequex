package jetstream

import "github.com/nats-io/nats.go"

type Publisher struct {
	nats    *nats.Conn
	js      *nats.JetStreamContext
	subject string
}

func NewPublisher(nats *nats.Conn, js *nats.JetStreamContext, subject string) *Publisher {
	return &Publisher{nats: nats, js: js, subject: subject}
}

func (p *Publisher) Publish(data []byte) error {
	_, err := (*p.js).Publish(p.subject, data)
	return err
}
