package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ConnectionType represents the type of connection
type ConnectionType string

const (
	ConnectionTypeNATS      ConnectionType = "nats"
	ConnectionTypeJetStream ConnectionType = "jetstream"
	ConnectionTypeTLS       ConnectionType = "tls"
)

// ConnectionConfig represents a parsed connection string configuration
type ConnectionConfig struct {
	Type     ConnectionType
	Host     string
	Port     int
	Username string
	Password string
	Params   map[string]string
}

// ParseConnectionString parses a connection string and returns a ConnectionConfig
// Examples:
//   - nats://127.0.0.1:4222?subject=test
//   - jetstream://prod:4222?stream=cache&subject=binance.trades
//   - nats://user:pass@localhost:4222?subject=test
func ParseConnectionString(connStr string) (*ConnectionConfig, error) {
	if connStr == "" {
		return nil, fmt.Errorf("connection string cannot be empty")
	}

	// Handle the @ prefix if present
	connStr = strings.TrimPrefix(connStr, "@")

	// Parse the URL
	u, err := url.Parse(connStr)
	if err != nil {
		return nil, fmt.Errorf("invalid connection string format: %w", err)
	}

	// Determine connection type
	var connType ConnectionType
	switch u.Scheme {
	case "nats":
		connType = ConnectionTypeNATS
	case "jetstream":
		connType = ConnectionTypeJetStream
	case "tls":
		connType = ConnectionTypeTLS
	default:
		return nil, fmt.Errorf("unsupported connection type: %s. Supported types: nats, jetstream, tls", u.Scheme)
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

	return &ConnectionConfig{
		Type:     connType,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Params:   params,
	}, nil
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
	scheme := string(c.Type)
	if c.Type == ConnectionTypeJetStream {
		scheme = "nats" // JetStream uses nats:// scheme
	}

	// Build user info if credentials are present
	var userInfo string
	if c.Username != "" {
		userInfo = c.Username
		if c.Password != "" {
			userInfo += ":" + c.Password
		}
		userInfo += "@"
	}

	// Build query string
	var queryParts []string
	for key, value := range c.Params {
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

	// Validate connection type specific requirements
	switch c.Type {
	case ConnectionTypeJetStream:
		// JetStream typically requires a stream parameter
		if _, hasStream := c.Params["stream"]; !hasStream {
			return fmt.Errorf("jetstream connection requires 'stream' parameter")
		}
	}

	return nil
}
