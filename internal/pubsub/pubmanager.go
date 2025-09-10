package pubsub

import (
	"time"

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
		// Configure NATS connection with proper timeouts for JetStream
		opts := []nats.Option{
			nats.Timeout(1 * time.Second),       // Connection timeout
			nats.ReconnectWait(2 * time.Second), // Reconnect wait time
			nats.MaxReconnects(-1),              // Unlimited reconnects
			nats.PingInterval(20 * time.Second), // Ping interval
			nats.MaxPingsOutstanding(3),         // Max outstanding pings
		}

		natsConn, err := nats.Connect(connConfig.ToNATSURL(), opts...)
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
