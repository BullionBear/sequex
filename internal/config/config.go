package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Account struct {
	Name      string `yaml:"name"`
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
}

type Config struct {
	Accounts map[string][]Account `yaml:"accounts"`
	Market   map[string][]string  `yaml:"market"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
