package pubsub

import "github.com/nats-io/nats.go"

type Publisher struct {
	nats       *nats.Conn
	js         *nats.JetStreamContext
	streamName string
	subject    string
}

func NewPublisher(nats *nats.Conn, jetstream string, subject string) (*Publisher, error) {
	js, err := nats.JetStream()
	if err != nil {
		return nil, err
	}
	return &Publisher{nats: nats, js: &js, streamName: jetstream, subject: subject}, nil
}

func (p *Publisher) Publish(data []byte, headers map[string]string) error {
	header := nats.Header{}
	for k, v := range headers {
		header[k] = []string{v}
	}

	_, err := (*p.js).PublishMsg(&nats.Msg{
		Subject: p.subject,
		Data:    data,
		Header:  header,
	})
	return err
}
