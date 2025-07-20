package bybit

import (
	"time"
)

// Config represents the configuration for Bybit API client
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

// DefaultConfig returns a default configuration for Bybit API
func DefaultConfig() *Config {
	return &Config{
		BaseURL:    BaseURLMainnet,
		Timeout:    30 * time.Second,
		UseTestnet: false,
	}
}

// TestnetConfig returns a configuration for Bybit testnet
func TestnetConfig() *Config {
	config := DefaultConfig()
	config.BaseURL = BaseURLTestnet
	config.UseTestnet = true
	return config
}

// WithAPIKey sets the API key for the configuration
func (c *Config) WithAPIKey(apiKey string) *Config {
	c.APIKey = apiKey
	return c
}

// WithAPISecret sets the API secret for the configuration
func (c *Config) WithAPISecret(apiSecret string) *Config {
	c.APISecret = apiSecret
	return c
}

// WithBaseURL sets the base URL for the configuration
func (c *Config) WithBaseURL(baseURL string) *Config {
	c.BaseURL = baseURL
	return c
}

// WithTimeout sets the timeout for the configuration
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}
