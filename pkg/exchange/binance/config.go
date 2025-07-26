package binance

type Config struct {
	// API credentials
	APIKey    string
	APISecret string

	// API endpoints
	BaseURL string
}

type WSConfig struct {
	// API credentials
	APIKey    string
	APISecret string

	// API endpoints
	BaseURL string
}
