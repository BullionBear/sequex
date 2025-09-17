package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// NATSConfig represents NATS connection configuration
type NATSConfig struct {
	URIs    string `json:"uris"`
	Stream  string `json:"stream"`
	Subject string `json:"subject"`
}

// Config represents the main configuration structure
type Config struct {
	Exchange   string     `json:"exchange"`
	Instrument string     `json:"instrument"`
	Symbol     string     `json:"symbol"`
	Type       string     `json:"type"`
	NATS       NATSConfig `json:"nats"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		return nil, fmt.Errorf("config file path cannot be empty")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", filePath, err)
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration in %s: %w", filePath, err)
	}

	return &config, nil
}

// Validate validates the main configuration
func (c *Config) Validate() error {
	if c.Exchange == "" {
		return fmt.Errorf("exchange cannot be empty")
	}

	if c.Instrument == "" {
		return fmt.Errorf("instrument cannot be empty")
	}

	if c.Symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}

	if c.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}

	// Validate NATS configuration
	return c.NATS.Validate()
}

// Validate validates the NATS configuration
func (n *NATSConfig) Validate() error {
	if n.URIs == "" {
		return fmt.Errorf("nats.uris cannot be empty")
	}

	if n.Stream == "" {
		return fmt.Errorf("nats.stream cannot be empty")
	}

	if n.Subject == "" {
		return fmt.Errorf("nats.subject cannot be empty")
	}

	// Validate that URIs are valid NATS URLs
	uris := strings.Split(n.URIs, ",")
	for i, uri := range uris {
		uri = strings.TrimSpace(uri)
		if uri == "" {
			continue
		}

		parsedURL, err := url.Parse(uri)
		if err != nil {
			return fmt.Errorf("invalid NATS URI at index %d: %w", i, err)
		}

		if parsedURL.Scheme != "nats" {
			return fmt.Errorf("invalid NATS URI scheme at index %d: expected 'nats', got '%s'", i, parsedURL.Scheme)
		}

		if parsedURL.Hostname() == "" {
			return fmt.Errorf("invalid NATS URI at index %d: hostname cannot be empty", i)
		}
	}

	return nil
}

// GetNATSURIs returns the NATS URIs
func (n *NATSConfig) GetNATSURIs() []string {
	return strings.Split(n.URIs, ",")
}
