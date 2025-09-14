package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
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

// GetNATSURIs returns a slice of individual NATS URIs
func (n *NATSConfig) GetNATSURIs() []string {
	uris := strings.Split(n.URIs, ",")
	var cleanURIs []string
	for _, uri := range uris {
		uri = strings.TrimSpace(uri)
		if uri != "" {
			cleanURIs = append(cleanURIs, uri)
		}
	}
	return cleanURIs
}

// ConnectionConfig represents a parsed connection string configuration
type ConnectionConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Params   map[string]string
}

// ParseConnectionString parses a connection string and returns a ConnectionConfig
// Examples:
//   - nats://127.0.0.1:4222?stream=feed&subject=test
//   - nats://user:pass@127.0.0.1:4022?stream=feed&subject=trade.btcusdt
//   - @nats://user:pass@localhost:4222?stream=feed&subject=test (with @ prefix for auth)
func ParseConnectionString(connStr string) (*ConnectionConfig, error) {
	if connStr == "" {
		return nil, fmt.Errorf("connection string cannot be empty")
	}

	// Handle the @ prefix if present (indicates username/password authentication)
	connStr = strings.TrimPrefix(connStr, "@")

	// Parse the URL
	u, err := url.Parse(connStr)
	if err != nil {
		return nil, fmt.Errorf("invalid connection string format: %w", err)
	}

	// Validate that only nats:// scheme is supported
	if u.Scheme != "nats" {
		return nil, fmt.Errorf("unsupported connection scheme: %s. Only nats:// is supported", u.Scheme)
	}

	// Parse host and port
	host := u.Hostname()
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	port := 4222 // Default NATS port
	if u.Port() != "" {
		var err error
		port, err = strconv.Atoi(u.Port())
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %w", err)
		}
	}

	// Parse credentials
	username := u.User.Username()
	password, _ := u.User.Password()

	// Parse query parameters
	params := make(map[string]string)
	for key, values := range u.Query() {
		if len(values) > 0 {
			params[key] = values[0] // Take the first value if multiple are provided
		}
	}

	config := &ConnectionConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Params:   params,
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// GetParam returns a query parameter value, with an optional default
func (c *ConnectionConfig) GetParam(key, defaultValue string) string {
	if value, exists := c.Params[key]; exists {
		return value
	}
	return defaultValue
}

// GetIntParam returns a query parameter as an integer, with an optional default
func (c *ConnectionConfig) GetIntParam(key string, defaultValue int) (int, error) {
	if value, exists := c.Params[key]; exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("invalid integer parameter '%s': %w", key, err)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

// GetBoolParam returns a query parameter as a boolean, with an optional default
func (c *ConnectionConfig) GetBoolParam(key string, defaultValue bool) (bool, error) {
	if value, exists := c.Params[key]; exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, fmt.Errorf("invalid boolean parameter '%s': %w", key, err)
		}
		return boolValue, nil
	}
	return defaultValue, nil
}

// ToNATSURL converts the connection config back to a NATS-compatible URL
func (c *ConnectionConfig) ToNATSURL() string {
	scheme := "nats"

	// Build user info if credentials are present
	var userInfo string
	if c.Username != "" {
		userInfo = c.Username
		if c.Password != "" {
			userInfo += ":" + c.Password
		}
		userInfo += "@"
	}

	// Build query string with sorted parameters for consistent output
	var keys []string
	for key := range c.Params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var queryParts []string
	for _, key := range keys {
		value := c.Params[key]
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
	}
	queryString := ""
	if len(queryParts) > 0 {
		queryString = "?" + strings.Join(queryParts, "&")
	}

	return fmt.Sprintf("%s://%s%s:%d%s", scheme, userInfo, c.Host, c.Port, queryString)
}

// String returns a string representation of the connection config
func (c *ConnectionConfig) String() string {
	return c.ToNATSURL()
}

// Validate performs validation on the connection configuration
func (c *ConnectionConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}

	// Stream parameter is mandatory for all connections
	streamValue, hasStream := c.Params["stream"]
	if !hasStream {
		return fmt.Errorf("stream parameter is required")
	}
	if streamValue == "" {
		return fmt.Errorf("stream parameter cannot be empty")
	}

	// Subject parameter is mandatory for all connections
	subjectValue, hasSubject := c.Params["subject"]
	if !hasSubject {
		return fmt.Errorf("subject parameter is required")
	}
	if subjectValue == "" {
		return fmt.Errorf("subject parameter cannot be empty")
	}

	return nil
}
