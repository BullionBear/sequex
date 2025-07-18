package binance

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// AppConfig represents the entire application configuration
type AppConfig struct {
	Accounts map[string][]Config `yaml:"accounts"`
	Market   map[string][]string `yaml:"market"`
}

// LoadConfig loads configuration from YAML file
func LoadConfig(filePath string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetBinanceConfig returns the first Binance configuration
func (ac *AppConfig) GetBinanceConfig() (*Config, error) {
	configs, exists := ac.Accounts["binance"]
	if !exists || len(configs) == 0 {
		return nil, fmt.Errorf("no Binance configuration found")
	}

	return &configs[0], nil
}

// GetBinanceConfigByName returns Binance configuration by name
func (ac *AppConfig) GetBinanceConfigByName(name string) (*Config, error) {
	configs, exists := ac.Accounts["binance"]
	if !exists {
		return nil, fmt.Errorf("no Binance configuration found")
	}

	for _, config := range configs {
		if config.Name == name {
			return &config, nil
		}
	}

	return nil, fmt.Errorf("Binance configuration with name '%s' not found", name)
}

// GetBinanceSymbols returns configured Binance trading symbols
func (ac *AppConfig) GetBinanceSymbols() []string {
	if symbols, exists := ac.Market["binance"]; exists {
		return symbols
	}
	return []string{}
}
