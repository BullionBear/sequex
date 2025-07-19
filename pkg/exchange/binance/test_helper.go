package binance

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// TestConfig represents the structure of binance_test.yml
type TestConfig struct {
	Binance struct {
		Testnet struct {
			APIKey    string `yaml:"api_key"`
			APISecret string `yaml:"api_secret"`
			BaseURL   string `yaml:"base_url"`
			Timeout   int    `yaml:"timeout"`
		} `yaml:"testnet"`
	} `yaml:"binance"`
}

// LoadTestConfig loads configuration from binance_test.yml
func LoadTestConfig() (*Config, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Look for binance_test.yml in current directory and parent directories
	var configPath string
	dir := wd
	for {
		path := filepath.Join(dir, "binance_test.yml")
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}

	if configPath == "" {
		// If no config file found, return testnet config with empty credentials
		// This allows tests to run without real credentials (they will just fail gracefully)
		config := TestnetConfig()
		config.APIKey = ""
		config.APISecret = ""
		return config, nil
	}

	// Read and parse the YAML file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var testConfig TestConfig
	if err := yaml.Unmarshal(data, &testConfig); err != nil {
		return nil, err
	}

	// Create Config from test configuration
	config := &Config{
		APIKey:     testConfig.Binance.Testnet.APIKey,
		APISecret:  testConfig.Binance.Testnet.APISecret,
		BaseURL:    testConfig.Binance.Testnet.BaseURL,
		Timeout:    time.Duration(testConfig.Binance.Testnet.Timeout) * time.Second,
		UseTestnet: true,
	}

	return config, nil
}

// CreateTestClient creates a client configured for testnet with real credentials
func CreateTestClient() (*Client, error) {
	config, err := LoadTestConfig()
	if err != nil {
		return nil, err
	}

	return NewClient(config), nil
}

// HasTestCredentials checks if test credentials are available
func HasTestCredentials() bool {
	config, err := LoadTestConfig()
	if err != nil {
		return false
	}
	return config.APIKey != "" && config.APISecret != ""
}
