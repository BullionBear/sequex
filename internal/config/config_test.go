package config

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestParseConnectionString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *ConnectionConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:  "basic NATS connection",
			input: "nats://127.0.0.1:4222",
			expected: &ConnectionConfig{
				Type:     ConnectionTypeNATS,
				Host:     "127.0.0.1",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{},
			},
			expectError: false,
		},
		{
			name:  "NATS with @ prefix",
			input: "@nats://127.0.0.1:4222?subject=test",
			expected: &ConnectionConfig{
				Type:     ConnectionTypeNATS,
				Host:     "127.0.0.1",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{"subject": "test"},
			},
			expectError: false,
		},
		{
			name:  "JetStream connection",
			input: "jetstream://prod:4222?stream=cache&subject=binance.trades",
			expected: &ConnectionConfig{
				Type:     ConnectionTypeJetStream,
				Host:     "prod",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{"stream": "cache", "subject": "binance.trades"},
			},
			expectError: false,
		},
		{
			name:  "NATS with credentials",
			input: "nats://user:pass@localhost:4222?subject=test",
			expected: &ConnectionConfig{
				Type:     ConnectionTypeNATS,
				Host:     "localhost",
				Port:     4222,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"subject": "test"},
			},
			expectError: false,
		},
		{
			name:  "NATS with default port",
			input: "nats://localhost",
			expected: &ConnectionConfig{
				Type:     ConnectionTypeNATS,
				Host:     "localhost",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{},
			},
			expectError: false,
		},
		{
			name:  "TLS connection",
			input: "tls://secure.example.com:4222?subject=secure.test",
			expected: &ConnectionConfig{
				Type:     ConnectionTypeTLS,
				Host:     "secure.example.com",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{"subject": "secure.test"},
			},
			expectError: false,
		},
		{
			name:        "empty connection string",
			input:       "",
			expected:    nil,
			expectError: true,
			errorMsg:    "connection string cannot be empty",
		},
		{
			name:        "invalid scheme",
			input:       "http://localhost:4222",
			expected:    nil,
			expectError: true,
			errorMsg:    "unsupported connection type: http",
		},
		{
			name:        "invalid URL format",
			input:       "nats://[invalid-url",
			expected:    nil,
			expectError: true,
			errorMsg:    "invalid connection string format",
		},
		{
			name:        "invalid port",
			input:       "nats://localhost:invalid",
			expected:    nil,
			expectError: true,
			errorMsg:    "invalid connection string format",
		},
		{
			name:        "empty host",
			input:       "nats://:4222",
			expected:    nil,
			expectError: true,
			errorMsg:    "host cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseConnectionString(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("expected result but got nil")
				return
			}

			// Compare fields
			if result.Type != tt.expected.Type {
				t.Errorf("expected Type %v, got %v", tt.expected.Type, result.Type)
			}
			if result.Host != tt.expected.Host {
				t.Errorf("expected Host %v, got %v", tt.expected.Host, result.Host)
			}
			if result.Port != tt.expected.Port {
				t.Errorf("expected Port %v, got %v", tt.expected.Port, result.Port)
			}
			if result.Username != tt.expected.Username {
				t.Errorf("expected Username %v, got %v", tt.expected.Username, result.Username)
			}
			if result.Password != tt.expected.Password {
				t.Errorf("expected Password %v, got %v", tt.expected.Password, result.Password)
			}

			// Compare params
			if len(result.Params) != len(tt.expected.Params) {
				t.Errorf("expected %d params, got %d", len(tt.expected.Params), len(result.Params))
			}
			for key, expectedValue := range tt.expected.Params {
				if actualValue, exists := result.Params[key]; !exists {
					t.Errorf("expected param '%s' not found", key)
				} else if actualValue != expectedValue {
					t.Errorf("expected param '%s' to be '%s', got '%s'", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestConnectionConfig_GetParam(t *testing.T) {
	config := &ConnectionConfig{
		Params: map[string]string{
			"subject": "test.subject",
		},
	}

	tests := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{"subject", "default.subject", "test.subject"},
		{"nonexistent", "default.value", "default.value"},
	}

	for _, tt := range tests {
		result := config.GetParam(tt.key, tt.defaultValue)
		if result != tt.expected {
			t.Errorf("GetParam(%s, %s) = %s, expected %s", tt.key, tt.defaultValue, result, tt.expected)
		}
	}
}

func TestConnectionConfig_GetIntParam(t *testing.T) {
	config := &ConnectionConfig{
		Params: map[string]string{
			"port": "8080",
		},
	}

	tests := []struct {
		key          string
		defaultValue int
		expected     int
		expectError  bool
	}{
		{"port", 3000, 8080, false},
		{"nonexistent", 100, 100, false},
		{"invalid", 0, 0, true},
	}

	// Add invalid parameter
	config.Params["invalid"] = "not-a-number"

	for _, tt := range tests {
		result, err := config.GetIntParam(tt.key, tt.defaultValue)

		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for key '%s' but got none", tt.key)
			}
			continue
		}

		if err != nil {
			t.Errorf("unexpected error for key '%s': %v", tt.key, err)
			continue
		}

		if result != tt.expected {
			t.Errorf("GetIntParam(%s, %d) = %d, expected %d", tt.key, tt.defaultValue, result, tt.expected)
		}
	}
}

func TestConnectionConfig_GetBoolParam(t *testing.T) {
	config := &ConnectionConfig{
		Params: map[string]string{
			"enabled":  "true",
			"disabled": "false",
			"invalid":  "maybe",
		},
	}

	tests := []struct {
		key          string
		defaultValue bool
		expected     bool
		expectError  bool
	}{
		{"enabled", false, true, false},
		{"disabled", true, false, false},
		{"nonexistent", true, true, false},
		{"invalid", false, false, true},
	}

	for _, tt := range tests {
		result, err := config.GetBoolParam(tt.key, tt.defaultValue)

		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for key '%s' but got none", tt.key)
			}
			continue
		}

		if err != nil {
			t.Errorf("unexpected error for key '%s': %v", tt.key, err)
			continue
		}

		if result != tt.expected {
			t.Errorf("GetBoolParam(%s, %t) = %t, expected %t", tt.key, tt.defaultValue, result, tt.expected)
		}
	}
}

func TestConnectionConfig_ToNATSURL(t *testing.T) {
	tests := []struct {
		name     string
		config   *ConnectionConfig
		expected string
	}{
		{
			name: "basic NATS",
			config: &ConnectionConfig{
				Type: ConnectionTypeNATS,
				Host: "localhost",
				Port: 4222,
			},
			expected: "nats://localhost:4222",
		},
		{
			name: "NATS with credentials and params",
			config: &ConnectionConfig{
				Type:     ConnectionTypeNATS,
				Host:     "localhost",
				Port:     4222,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"subject": "test"},
			},
			expected: "nats://user:pass@localhost:4222?subject=test",
		},
		{
			name: "JetStream",
			config: &ConnectionConfig{
				Type:   ConnectionTypeJetStream,
				Host:   "prod",
				Port:   4222,
				Params: map[string]string{"stream": "cache", "subject": "binance.trades"},
			},
			expected: "nats://prod:4222?stream=cache&subject=binance.trades",
		},
		{
			name: "TLS",
			config: &ConnectionConfig{
				Type:   ConnectionTypeTLS,
				Host:   "secure.example.com",
				Port:   4222,
				Params: map[string]string{"subject": "secure.test"},
			},
			expected: "tls://secure.example.com:4222?subject=secure.test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ToNATSURL()

			if result != tt.expected {
				t.Errorf("ToNATSURL() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestConnectionConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *ConnectionConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid NATS config",
			config: &ConnectionConfig{
				Type: ConnectionTypeNATS,
				Host: "localhost",
				Port: 4222,
			},
			expectError: false,
		},
		{
			name: "valid JetStream config",
			config: &ConnectionConfig{
				Type:   ConnectionTypeJetStream,
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"stream": "test"},
			},
			expectError: false,
		},
		{
			name: "empty host",
			config: &ConnectionConfig{
				Type: ConnectionTypeNATS,
				Host: "",
				Port: 4222,
			},
			expectError: true,
			errorMsg:    "host cannot be empty",
		},
		{
			name: "invalid port - too low",
			config: &ConnectionConfig{
				Type: ConnectionTypeNATS,
				Host: "localhost",
				Port: 0,
			},
			expectError: true,
			errorMsg:    "port must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			config: &ConnectionConfig{
				Type: ConnectionTypeNATS,
				Host: "localhost",
				Port: 65536,
			},
			expectError: true,
			errorMsg:    "port must be between 1 and 65535",
		},
		{
			name: "JetStream without stream",
			config: &ConnectionConfig{
				Type: ConnectionTypeJetStream,
				Host: "localhost",
				Port: 4222,
			},
			expectError: true,
			errorMsg:    "jetstream connection requires 'stream' parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestConnectionConfig_String(t *testing.T) {
	config := &ConnectionConfig{
		Type:   ConnectionTypeNATS,
		Host:   "localhost",
		Port:   4222,
		Params: map[string]string{"subject": "test"},
	}

	expected := "nats://localhost:4222?subject=test"
	result := config.String()

	if result != expected {
		t.Errorf("String() = %s, expected %s", result, expected)
	}
}

// ExampleParseConnectionString demonstrates how to parse various connection strings
func ExampleParseConnectionString() {
	// Example 1: Basic NATS connection
	connStr1 := "nats://127.0.0.1:4222?subject=test"
	config1, err := ParseConnectionString(connStr1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Type: %s, Host: %s, Port: %d, Subject: %s\n",
		config1.Type, config1.Host, config1.Port, config1.GetParam("subject", ""))

	// Example 2: JetStream connection with @ prefix
	connStr2 := "@jetstream://prod:4222?stream=cache&subject=binance.trades"
	config2, err := ParseConnectionString(connStr2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Type: %s, Host: %s, Port: %d, Stream: %s, Subject: %s\n",
		config2.Type, config2.Host, config2.Port,
		config2.GetParam("stream", ""), config2.GetParam("subject", ""))

	// Example 3: NATS with credentials
	connStr3 := "nats://user:pass@localhost:4222?subject=test"
	config3, err := ParseConnectionString(connStr3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Type: %s, Host: %s, Port: %d, Username: %s, Subject: %s\n",
		config3.Type, config3.Host, config3.Port, config3.Username,
		config3.GetParam("subject", ""))

	// Example 4: Validate the configuration
	if err := config3.Validate(); err != nil {
		log.Fatal("Validation failed:", err)
	}

	// Example 5: Convert back to NATS URL
	natsURL := config3.ToNATSURL()
	fmt.Printf("NATS URL: %s\n", natsURL)

	// Output:
	// Type: nats, Host: 127.0.0.1, Port: 4222, Subject: test
	// Type: jetstream, Host: prod, Port: 4222, Stream: cache, Subject: binance.trades
	// Type: nats, Host: localhost, Port: 4222, Username: user, Subject: test
	// NATS URL: nats://user:pass@localhost:4222?subject=test
}
