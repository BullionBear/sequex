package binancefuture

import (
	"time"
)

// Config represents the configuration for Binance Futures API client
type Config struct {
	// API credentials
	APIKey    string
	APISecret string

	// API endpoints
	BaseURL string

	// HTTP client configuration
	Timeout time.Duration

	// Testnet flag
	UseTestnet bool
}

// DefaultConfig returns a default configuration for Binance Futures API
func DefaultConfig() *Config {
	return &Config{
		BaseURL:    BaseURLFutures,
		Timeout:    30 * time.Second,
		UseTestnet: false,
	}
}

// TestnetConfig returns a configuration for Binance Futures testnet
func TestnetConfig() *Config {
	config := DefaultConfig()
	config.BaseURL = BaseURLFuturesTestnet
	config.UseTestnet = true
	return config
}
