package exchange

type Constructor[C IsolatedSpotConnector] func(credentials Credentials) C

var (
	IsolatedSpotConnectorFactories = map[MarketType]Constructor[IsolatedSpotConnector]{}
)

func Register[C IsolatedSpotConnector](marketType MarketType, constructor Constructor[C]) {
	IsolatedSpotConnectorFactories[marketType] = func(credentials Credentials) IsolatedSpotConnector {
		return constructor(credentials)
	}
}
