package exchange

type Config struct {
	MarketType   MarketType
	Credentials  Credentials
	WalletType   WalletType
	MarginType   MarginType
	PositionType PositionType
}

type Credentials struct {
	APIKey     string
	APISecret  string
	Passphrase string
}
