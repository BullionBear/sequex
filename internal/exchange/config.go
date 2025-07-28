package exchange

type Config struct {
	MarketType   MarketType
	Credentials  map[string]string
	WalletType   WalletType
	MarginType   MarginType
	PositionType PositionType
}
