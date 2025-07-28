package binance

type Config struct {
	// API credentials
	APIKey    string
	APISecret string

	// API endpoints
	BaseURL string
}

func NewConfig(apiKey, apiSecret, baseURL string) *Config {
	return &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   baseURL,
	}
}

func NewMainnetConfig(apiKey, apiSecret string) *Config {
	return NewConfig(apiKey, apiSecret, MainnetBaseUrl)
}

func NewTestnetConfig(apiKey, apiSecret string) *Config {
	return NewConfig(apiKey, apiSecret, TestnetBaseUrl)
}

type WSConfig struct {
	// API credentials
	APIKey    string
	APISecret string

	// API endpoints
	BaseWsURL   string
	BaseRestURL string
}

func NewMainnetWSConfig(apiKey, apiSecret string) *WSConfig {
	return &WSConfig{
		APIKey:      apiKey,
		APISecret:   apiSecret,
		BaseWsURL:   MainnetWSBaseUrl,
		BaseRestURL: MainnetBaseUrl,
	}
}

func NewTestnetWSConfig(apiKey, apiSecret string) *WSConfig {
	return &WSConfig{
		APIKey:      apiKey,
		APISecret:   apiSecret,
		BaseWsURL:   TestnetWSBaseUrl,
		BaseRestURL: TestnetBaseUrl,
	}
}
