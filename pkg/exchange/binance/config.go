package binance

import (
	"time"
)

// Config represents the configuration for Binance API client
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

// DefaultConfig returns a default configuration for Binance API
func DefaultConfig() *Config {
	return &Config{
		BaseURL:    "https://api.binance.com",
		Timeout:    30 * time.Second,
		UseTestnet: false,
	}
}

// TestnetConfig returns a configuration for Binance testnet
func TestnetConfig() *Config {
	config := DefaultConfig()
	config.BaseURL = "https://testnet.binance.vision"
	config.UseTestnet = true
	return config
}
