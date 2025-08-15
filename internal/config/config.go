package config

import (
	"fmt"
	"os"

	"github.com/BullionBear/sequex/pkg/node"

	"gopkg.in/yaml.v3"
)

// Config represents the merged configuration structure
type SrvConfig struct {
	App struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"app"`
	EventBus struct {
		Url string `yaml:"url"`
	} `yaml:"eventbus"`
	Node node.NodeConfig `yaml:"node"`
}

// LoadConfig loads the merged configuration from file
func LoadConfig[T any](configPath string) (*T, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config T
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}
