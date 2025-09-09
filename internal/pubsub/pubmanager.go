package pubsub

import (
	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/pkg/logger"
	"github.com/nats-io/nats.go"
)

type PubManager struct {
	publishers []*Publisher
}

func NewPubManager(connConfigs []*config.ConnectionConfig) (*PubManager, error) {
	publishers := make([]*Publisher, 0)
	for _, connConfig := range connConfigs {
		natsConn, err := nats.Connect(connConfig.ToNATSURL())
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to connect to NATS")
			return nil, err
		}
		publisher, err := NewPublisher(natsConn, connConfig.GetParam("stream", ""), connConfig.GetParam("subject", ""))
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to create publisher")
			return nil, err
		}
		publishers = append(publishers, publisher)
	}
	return &PubManager{
		publishers: publishers,
	}, nil
}

func (p *PubManager) Publish(data []byte, headers map[string]string) error {
	for _, publisher := range p.publishers {
		err := publisher.Publish(data, headers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PubManager) Close() {
	for _, publisher := range p.publishers {
		publisher.Close()
	}
}
