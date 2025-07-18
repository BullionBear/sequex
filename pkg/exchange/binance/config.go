package binance

// Config represents Binance API configuration
type Config struct {
	Name      string `yaml:"name" json:"name"`
	APIKey    string `yaml:"api_key" json:"api_key"`
	APISecret string `yaml:"api_secret" json:"api_secret"`
	Sandbox   bool   `yaml:"sandbox" json:"sandbox"`
	Timeout   int    `yaml:"timeout" json:"timeout"` // in seconds
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Name:    "default",
		Sandbox: false,
		Timeout: 30,
	}
}

// IsValid checks if configuration has required fields
func (c *Config) IsValid() bool {
	return c.APIKey != "" && c.APISecret != ""
}

// GetBaseURL returns the appropriate base URL
func (c *Config) GetBaseURL() string {
	if c.Sandbox {
		return SandboxBaseURL
	}
	return BaseURL
}
