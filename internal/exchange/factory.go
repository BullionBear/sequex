package exchange

import "fmt"

type Constructor[C IsolatedSpotConnector] func(credentials Credentials) C

var (
	IsolatedSpotConnectorFactories = map[MarketType]Constructor[IsolatedSpotConnector]{}
)

func Register[C IsolatedSpotConnector](marketType MarketType, constructor Constructor[C]) {
	IsolatedSpotConnectorFactories[marketType] = func(credentials Credentials) IsolatedSpotConnector {
		return constructor(credentials)
	}
}

func NewConnector[C IsolatedSpotConnector](marketType MarketType, credentials Credentials) (C, error) {
	constructor, ok := IsolatedSpotConnectorFactories[marketType]
	if !ok {
		var zero C
		return zero, fmt.Errorf("unsupported market type: %s", marketType)
	}
	return constructor(credentials).(C), nil
}
