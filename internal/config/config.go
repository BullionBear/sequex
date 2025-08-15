package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	PathGlobalConfig = "~/.sequex/config.yml"
)

type GlobalConfig struct {
	EventBus struct {
		Url string `yaml:"url"`
	} `yaml:"eventbus"`
}

// expandTilde expands the tilde (~) to the user's home directory
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path // Return original path if we can't get home dir
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

// createDefaultConfig creates a default configuration file at the specified path
func createDefaultConfig(configPath string) error {
	// Expand tilde in path
	expandedPath := expandTilde(configPath)

	// Create directory if it doesn't exist
	dir := filepath.Dir(expandedPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create default config
	defaultConfig := GlobalConfig{}
	defaultConfig.EventBus.Url = "nats://localhost:4222"

	// Marshal to YAML
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(expandedPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}

	return nil
}

// LoadConfig loads the merged configuration from file
func LoadConfig[T any](configPath string) (*T, error) {
	// Expand tilde in path
	expandedPath := expandTilde(configPath)

	// Check if file exists, if not create default config
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		// Only create default config for GlobalConfig type
		// We can detect this by checking if the configPath matches PathGlobalConfig
		if configPath == PathGlobalConfig {
			if err := createDefaultConfig(configPath); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("config file does not exist: %s", expandedPath)
		}
	}

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config T
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}
