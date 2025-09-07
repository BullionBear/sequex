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
			name:  "basic NATS connection with jetstream and subject",
			input: "nats://127.0.0.1:4222?jetstream=feed&subject=test",
			expected: &ConnectionConfig{
				Host:     "127.0.0.1",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name:  "NATS with @ prefix and jetstream",
			input: "@nats://127.0.0.1:4222?jetstream=feed&subject=test",
			expected: &ConnectionConfig{
				Host:     "127.0.0.1",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name:  "NATS with JetStream",
			input: "nats://user:pass@127.0.0.1:4022?jetstream=feed&subject=trade.btcusdt",
			expected: &ConnectionConfig{
				Host:     "127.0.0.1",
				Port:     4022,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"jetstream": "feed", "subject": "trade.btcusdt"},
			},
			expectError: false,
		},
		{
			name:  "NATS with credentials and jetstream",
			input: "nats://user:pass@localhost:4222?jetstream=feed&subject=test",
			expected: &ConnectionConfig{
				Host:     "localhost",
				Port:     4222,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name:  "NATS with default port and jetstream",
			input: "nats://localhost?jetstream=feed&subject=test",
			expected: &ConnectionConfig{
				Host:     "localhost",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{"jetstream": "feed", "subject": "test"},
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
			errorMsg:    "unsupported connection scheme: http",
		},
		{
			name:        "jetstream scheme not supported",
			input:       "jetstream://localhost:4222?stream=test",
			expected:    nil,
			expectError: true,
			errorMsg:    "unsupported connection scheme: jetstream",
		},
		{
			name:        "tls scheme not supported",
			input:       "tls://localhost:4222?jetstream=feed",
			expected:    nil,
			expectError: true,
			errorMsg:    "unsupported connection scheme: tls",
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
		{
			name:        "NATS without jetstream parameter",
			input:       "nats://127.0.0.1:4222?subject=test",
			expected:    nil,
			expectError: true,
			errorMsg:    "jetstream parameter is required",
		},
		{
			name:        "NATS with empty jetstream parameter",
			input:       "nats://127.0.0.1:4222?jetstream=&subject=test",
			expected:    nil,
			expectError: true,
			errorMsg:    "jetstream parameter cannot be empty",
		},
		{
			name:        "NATS without subject parameter",
			input:       "nats://127.0.0.1:4222?jetstream=feed",
			expected:    nil,
			expectError: true,
			errorMsg:    "subject parameter is required",
		},
		{
			name:        "NATS with empty subject parameter",
			input:       "nats://127.0.0.1:4222?jetstream=feed&subject=",
			expected:    nil,
			expectError: true,
			errorMsg:    "subject parameter cannot be empty",
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
			name: "basic NATS with jetstream and subject",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expected: "nats://localhost:4222?jetstream=feed&subject=test",
		},
		{
			name: "NATS with credentials and params",
			config: &ConnectionConfig{
				Host:     "localhost",
				Port:     4222,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expected: "nats://user:pass@localhost:4222?jetstream=feed&subject=test",
		},
		{
			name: "NATS with JetStream",
			config: &ConnectionConfig{
				Host:     "127.0.0.1",
				Port:     4022,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"jetstream": "feed", "subject": "trade.btcusdt"},
			},
			expected: "nats://user:pass@127.0.0.1:4022?jetstream=feed&subject=trade.btcusdt",
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
			name: "valid NATS config with jetstream and subject",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name: "valid NATS with JetStream config",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name: "empty host",
			config: &ConnectionConfig{
				Host:   "",
				Port:   4222,
				Params: map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expectError: true,
			errorMsg:    "host cannot be empty",
		},
		{
			name: "invalid port - too low",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   0,
				Params: map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expectError: true,
			errorMsg:    "port must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   65536,
				Params: map[string]string{"jetstream": "feed", "subject": "test"},
			},
			expectError: true,
			errorMsg:    "port must be between 1 and 65535",
		},
		{
			name: "NATS without jetstream parameter",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"subject": "test"},
			},
			expectError: true,
			errorMsg:    "jetstream parameter is required",
		},
		{
			name: "NATS with empty jetstream parameter",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"jetstream": "", "subject": "test"},
			},
			expectError: true,
			errorMsg:    "jetstream parameter cannot be empty",
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
		Host:   "localhost",
		Port:   4222,
		Params: map[string]string{"jetstream": "feed", "subject": "test"},
	}

	expected := "nats://localhost:4222?jetstream=feed&subject=test"
	result := config.String()

	if result != expected {
		t.Errorf("String() = %s, expected %s", result, expected)
	}
}

// ExampleParseConnectionString demonstrates how to parse various connection strings
func ExampleParseConnectionString() {
	// Example 1: Basic NATS connection with jetstream
	connStr1 := "nats://127.0.0.1:4222?jetstream=feed&subject=test"
	config1, err := ParseConnectionString(connStr1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Host: %s, Port: %d, JetStream: %s, Subject: %s\n",
		config1.Host, config1.Port, config1.GetParam("jetstream", ""), config1.GetParam("subject", ""))

	// Example 2: NATS with JetStream using new format
	connStr2 := "nats://user:pass@127.0.0.1:4022?jetstream=feed&subject=trade.btcusdt"
	config2, err := ParseConnectionString(connStr2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Host: %s, Port: %d, JetStream: %s, Subject: %s\n",
		config2.Host, config2.Port,
		config2.GetParam("jetstream", ""), config2.GetParam("subject", ""))

	// Example 3: NATS with credentials and jetstream
	connStr3 := "nats://user:pass@localhost:4222?jetstream=feed&subject=test"
	config3, err := ParseConnectionString(connStr3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Host: %s, Port: %d, Username: %s, JetStream: %s, Subject: %s\n",
		config3.Host, config3.Port, config3.Username,
		config3.GetParam("jetstream", ""), config3.GetParam("subject", ""))

	// Example 4: Validate the configuration
	if err := config3.Validate(); err != nil {
		log.Fatal("Validation failed:", err)
	}

	// Example 5: Convert back to NATS URL
	natsURL := config3.ToNATSURL()
	fmt.Printf("NATS URL: %s\n", natsURL)

	// Output:
	// Host: 127.0.0.1, Port: 4222, JetStream: feed, Subject: test
	// Host: 127.0.0.1, Port: 4022, JetStream: feed, Subject: trade.btcusdt
	// Host: localhost, Port: 4222, Username: user, JetStream: feed, Subject: test
	// NATS URL: nats://user:pass@localhost:4222?jetstream=feed&subject=test
}
