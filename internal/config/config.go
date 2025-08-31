package config

type FeedConfig struct {
	Exchange   string              `yaml:"exchange"`
	Type       string              `yaml:"type"`
	Symbol     string              `yaml:"symbol"`
	Instrument string              `yaml:"instrument"`
	Nats       NatsJetStreamConfig `yaml:"nats"`
}

type NatsJetStreamConfig struct {
	URL string `yaml:"url"`
}
