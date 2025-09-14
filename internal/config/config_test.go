package config

import (
	"fmt"
	"log"
	"os"
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
			name:  "basic NATS connection with stream and subject",
			input: "nats://127.0.0.1:4222?stream=feed&subject=test",
			expected: &ConnectionConfig{
				Host:     "127.0.0.1",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{"stream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name:  "NATS with @ prefix and stream",
			input: "@nats://127.0.0.1:4222?stream=feed&subject=test",
			expected: &ConnectionConfig{
				Host:     "127.0.0.1",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{"stream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name:  "NATS with JetStream",
			input: "nats://user:pass@127.0.0.1:4022?stream=feed&subject=trade.btcusdt",
			expected: &ConnectionConfig{
				Host:     "127.0.0.1",
				Port:     4022,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"stream": "feed", "subject": "trade.btcusdt"},
			},
			expectError: false,
		},
		{
			name:  "NATS with credentials and stream",
			input: "nats://user:pass@localhost:4222?stream=feed&subject=test",
			expected: &ConnectionConfig{
				Host:     "localhost",
				Port:     4222,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"stream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name:  "NATS with default port and stream",
			input: "nats://localhost?stream=feed&subject=test",
			expected: &ConnectionConfig{
				Host:     "localhost",
				Port:     4222,
				Username: "",
				Password: "",
				Params:   map[string]string{"stream": "feed", "subject": "test"},
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
			name:        "stream scheme not supported",
			input:       "stream://localhost:4222?stream=test",
			expected:    nil,
			expectError: true,
			errorMsg:    "unsupported connection scheme: stream",
		},
		{
			name:        "tls scheme not supported",
			input:       "tls://localhost:4222?stream=feed",
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
			name:        "NATS without stream parameter",
			input:       "nats://127.0.0.1:4222?subject=test",
			expected:    nil,
			expectError: true,
			errorMsg:    "stream parameter is required",
		},
		{
			name:        "NATS with empty stream parameter",
			input:       "nats://127.0.0.1:4222?stream=&subject=test",
			expected:    nil,
			expectError: true,
			errorMsg:    "stream parameter cannot be empty",
		},
		{
			name:        "NATS without subject parameter",
			input:       "nats://127.0.0.1:4222?stream=feed",
			expected:    nil,
			expectError: true,
			errorMsg:    "subject parameter is required",
		},
		{
			name:        "NATS with empty subject parameter",
			input:       "nats://127.0.0.1:4222?stream=feed&subject=",
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
			name: "basic NATS with stream and subject",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"stream": "feed", "subject": "test"},
			},
			expected: "nats://localhost:4222?stream=feed&subject=test",
		},
		{
			name: "NATS with credentials and params",
			config: &ConnectionConfig{
				Host:     "localhost",
				Port:     4222,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"stream": "feed", "subject": "test"},
			},
			expected: "nats://user:pass@localhost:4222?stream=feed&subject=test",
		},
		{
			name: "NATS with JetStream",
			config: &ConnectionConfig{
				Host:     "127.0.0.1",
				Port:     4022,
				Username: "user",
				Password: "pass",
				Params:   map[string]string{"stream": "feed", "subject": "trade.btcusdt"},
			},
			expected: "nats://user:pass@127.0.0.1:4022?stream=feed&subject=trade.btcusdt",
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
			name: "valid NATS config with stream and subject",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"stream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name: "valid NATS with JetStream config",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"stream": "feed", "subject": "test"},
			},
			expectError: false,
		},
		{
			name: "empty host",
			config: &ConnectionConfig{
				Host:   "",
				Port:   4222,
				Params: map[string]string{"stream": "feed", "subject": "test"},
			},
			expectError: true,
			errorMsg:    "host cannot be empty",
		},
		{
			name: "invalid port - too low",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   0,
				Params: map[string]string{"stream": "feed", "subject": "test"},
			},
			expectError: true,
			errorMsg:    "port must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   65536,
				Params: map[string]string{"stream": "feed", "subject": "test"},
			},
			expectError: true,
			errorMsg:    "port must be between 1 and 65535",
		},
		{
			name: "NATS without stream parameter",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"subject": "test"},
			},
			expectError: true,
			errorMsg:    "stream parameter is required",
		},
		{
			name: "NATS with empty stream parameter",
			config: &ConnectionConfig{
				Host:   "localhost",
				Port:   4222,
				Params: map[string]string{"stream": "", "subject": "test"},
			},
			expectError: true,
			errorMsg:    "stream parameter cannot be empty",
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
		Params: map[string]string{"stream": "feed", "subject": "test"},
	}

	expected := "nats://localhost:4222?stream=feed&subject=test"
	result := config.String()

	if result != expected {
		t.Errorf("String() = %s, expected %s", result, expected)
	}
}

// ExampleParseConnectionString demonstrates how to parse various connection strings
func ExampleParseConnectionString() {
	// Example 1: Basic NATS connection with stream
	connStr1 := "nats://127.0.0.1:4222?stream=feed&subject=test"
	config1, err := ParseConnectionString(connStr1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Host: %s, Port: %d, Stream: %s, Subject: %s\n",
		config1.Host, config1.Port, config1.GetParam("stream", ""), config1.GetParam("subject", ""))

	// Example 2: NATS with Stream using new format
	connStr2 := "nats://user:pass@127.0.0.1:4022?stream=feed&subject=trade.btcusdt"
	config2, err := ParseConnectionString(connStr2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Host: %s, Port: %d, Stream: %s, Subject: %s\n",
		config2.Host, config2.Port,
		config2.GetParam("stream", ""), config2.GetParam("subject", ""))

	// Example 3: NATS with credentials and stream
	connStr3 := "nats://user:pass@localhost:4222?stream=feed&subject=test"
	config3, err := ParseConnectionString(connStr3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Host: %s, Port: %d, Username: %s, Stream: %s, Subject: %s\n",
		config3.Host, config3.Port, config3.Username,
		config3.GetParam("stream", ""), config3.GetParam("subject", ""))

	// Example 4: Validate the configuration
	if err := config3.Validate(); err != nil {
		log.Fatal("Validation failed:", err)
	}

	// Example 5: Convert back to NATS URL
	natsURL := config3.ToNATSURL()
	fmt.Printf("NATS URL: %s\n", natsURL)

	// Output:
	// Host: 127.0.0.1, Port: 4222, Stream: feed, Subject: test
	// Host: 127.0.0.1, Port: 4022, Stream: feed, Subject: trade.btcusdt
	// Host: localhost, Port: 4222, Username: user, Stream: feed, Subject: test
	// NATS URL: nats://user:pass@localhost:4222?stream=feed&subject=test
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		jsonContent string
		expected    *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			jsonContent: `{
				"exchange": "binance",
				"instrument": "spot",
				"symbol": "BTC-USDT",
				"type": "trade",
				"nats": {
					"uris": "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224",
					"stream": "TRADE",
					"subject": "trade.binance.spot.btcusdt"
				}
			}`,
			expected: &Config{
				Exchange:   "binance",
				Instrument: "spot",
				Symbol:     "BTC-USDT",
				Type:       "trade",
				NATS: NATSConfig{
					URIs:    "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224",
					Stream:  "TRADE",
					Subject: "trade.binance.spot.btcusdt",
				},
			},
			expectError: false,
		},
		{
			name: "single NATS URI",
			jsonContent: `{
				"exchange": "binance",
				"instrument": "perp",
				"symbol": "ETH-USDT",
				"type": "trade",
				"nats": {
					"uris": "nats://localhost:4222",
					"stream": "TRADE",
					"subject": "trade.binance.perp.ethusdt"
				}
			}`,
			expected: &Config{
				Exchange:   "binance",
				Instrument: "perp",
				Symbol:     "ETH-USDT",
				Type:       "trade",
				NATS: NATSConfig{
					URIs:    "nats://localhost:4222",
					Stream:  "TRADE",
					Subject: "trade.binance.perp.ethusdt",
				},
			},
			expectError: false,
		},
		{
			name:        "empty exchange",
			jsonContent: `{"exchange": "", "instrument": "spot", "symbol": "BTC-USDT", "type": "trade", "nats": {"uris": "nats://localhost:4222", "stream": "TRADE", "subject": "test"}}`,
			expected:    nil,
			expectError: true,
			errorMsg:    "exchange cannot be empty",
		},
		{
			name:        "empty instrument",
			jsonContent: `{"exchange": "binance", "instrument": "", "symbol": "BTC-USDT", "type": "trade", "nats": {"uris": "nats://localhost:4222", "stream": "TRADE", "subject": "test"}}`,
			expected:    nil,
			expectError: true,
			errorMsg:    "instrument cannot be empty",
		},
		{
			name:        "empty symbol",
			jsonContent: `{"exchange": "binance", "instrument": "spot", "symbol": "", "type": "trade", "nats": {"uris": "nats://localhost:4222", "stream": "TRADE", "subject": "test"}}`,
			expected:    nil,
			expectError: true,
			errorMsg:    "symbol cannot be empty",
		},
		{
			name:        "empty type",
			jsonContent: `{"exchange": "binance", "instrument": "spot", "symbol": "BTC-USDT", "type": "", "nats": {"uris": "nats://localhost:4222", "stream": "TRADE", "subject": "test"}}`,
			expected:    nil,
			expectError: true,
			errorMsg:    "type cannot be empty",
		},
		{
			name:        "empty NATS URIs",
			jsonContent: `{"exchange": "binance", "instrument": "spot", "symbol": "BTC-USDT", "type": "trade", "nats": {"uris": "", "stream": "TRADE", "subject": "test"}}`,
			expected:    nil,
			expectError: true,
			errorMsg:    "nats.uris cannot be empty",
		},
		{
			name:        "empty NATS stream",
			jsonContent: `{"exchange": "binance", "instrument": "spot", "symbol": "BTC-USDT", "type": "trade", "nats": {"uris": "nats://localhost:4222", "stream": "", "subject": "test"}}`,
			expected:    nil,
			expectError: true,
			errorMsg:    "nats.stream cannot be empty",
		},
		{
			name:        "empty NATS subject",
			jsonContent: `{"exchange": "binance", "instrument": "spot", "symbol": "BTC-USDT", "type": "trade", "nats": {"uris": "nats://localhost:4222", "stream": "TRADE", "subject": ""}}`,
			expected:    nil,
			expectError: true,
			errorMsg:    "nats.subject cannot be empty",
		},
		{
			name:        "invalid NATS URI scheme",
			jsonContent: `{"exchange": "binance", "instrument": "spot", "symbol": "BTC-USDT", "type": "trade", "nats": {"uris": "http://localhost:4222", "stream": "TRADE", "subject": "test"}}`,
			expected:    nil,
			expectError: true,
			errorMsg:    "invalid NATS URI scheme",
		},
		{
			name:        "invalid JSON",
			jsonContent: `{"exchange": "binance", "instrument": "spot"`,
			expected:    nil,
			expectError: true,
			errorMsg:    "failed to parse config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpFile, err := os.CreateTemp("", "config-test-*.json")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write test content to file
			if _, err := tmpFile.WriteString(tt.jsonContent); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Test LoadConfig
			result, err := LoadConfig(tmpFile.Name())

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
			if result.Exchange != tt.expected.Exchange {
				t.Errorf("expected Exchange %v, got %v", tt.expected.Exchange, result.Exchange)
			}
			if result.Instrument != tt.expected.Instrument {
				t.Errorf("expected Instrument %v, got %v", tt.expected.Instrument, result.Instrument)
			}
			if result.Symbol != tt.expected.Symbol {
				t.Errorf("expected Symbol %v, got %v", tt.expected.Symbol, result.Symbol)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("expected Type %v, got %v", tt.expected.Type, result.Type)
			}

			// Compare NATS config
			if result.NATS.URIs != tt.expected.NATS.URIs {
				t.Errorf("expected NATS.URIs %v, got %v", tt.expected.NATS.URIs, result.NATS.URIs)
			}
			if result.NATS.Stream != tt.expected.NATS.Stream {
				t.Errorf("expected NATS.Stream %v, got %v", tt.expected.NATS.Stream, result.NATS.Stream)
			}
			if result.NATS.Subject != tt.expected.NATS.Subject {
				t.Errorf("expected NATS.Subject %v, got %v", tt.expected.NATS.Subject, result.NATS.Subject)
			}
		})
	}
}

func TestLoadConfig_FileErrors(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty file path",
			filePath:    "",
			expectError: true,
			errorMsg:    "config file path cannot be empty",
		},
		{
			name:        "non-existent file",
			filePath:    "/non/existent/file.json",
			expectError: true,
			errorMsg:    "failed to read config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LoadConfig(tt.filePath)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
				if result != nil {
					t.Errorf("expected nil result but got %v", result)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNATSConfig_GetNATSURIs(t *testing.T) {
	tests := []struct {
		name     string
		config   *NATSConfig
		expected []string
	}{
		{
			name: "single URI",
			config: &NATSConfig{
				URIs: "nats://localhost:4222",
			},
			expected: []string{"nats://localhost:4222"},
		},
		{
			name: "multiple URIs",
			config: &NATSConfig{
				URIs: "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224",
			},
			expected: []string{"nats://localhost:4222", "nats://localhost:4223", "nats://localhost:4224"},
		},
		{
			name: "URIs with spaces",
			config: &NATSConfig{
				URIs: "nats://localhost:4222, nats://localhost:4223 , nats://localhost:4224",
			},
			expected: []string{"nats://localhost:4222", "nats://localhost:4223", "nats://localhost:4224"},
		},
		{
			name: "empty URIs filtered out",
			config: &NATSConfig{
				URIs: "nats://localhost:4222,,nats://localhost:4223,",
			},
			expected: []string{"nats://localhost:4222", "nats://localhost:4223"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetNATSURIs()

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d URIs, got %d", len(tt.expected), len(result))
				return
			}

			for i, expectedURI := range tt.expected {
				if result[i] != expectedURI {
					t.Errorf("expected URI[%d] to be '%s', got '%s'", i, expectedURI, result[i])
				}
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &Config{
				Exchange:   "binance",
				Instrument: "spot",
				Symbol:     "BTC-USDT",
				Type:       "trade",
				NATS: NATSConfig{
					URIs:    "nats://localhost:4222",
					Stream:  "TRADE",
					Subject: "trade.binance.spot.btcusdt",
				},
			},
			expectError: false,
		},
		{
			name: "empty exchange",
			config: &Config{
				Exchange:   "",
				Instrument: "spot",
				Symbol:     "BTC-USDT",
				Type:       "trade",
				NATS: NATSConfig{
					URIs:    "nats://localhost:4222",
					Stream:  "TRADE",
					Subject: "test",
				},
			},
			expectError: true,
			errorMsg:    "exchange cannot be empty",
		},
		{
			name: "invalid NATS config",
			config: &Config{
				Exchange:   "binance",
				Instrument: "spot",
				Symbol:     "BTC-USDT",
				Type:       "trade",
				NATS: NATSConfig{
					URIs:    "",
					Stream:  "TRADE",
					Subject: "test",
				},
			},
			expectError: true,
			errorMsg:    "nats.uris cannot be empty",
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

func TestNATSConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *NATSConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid single URI",
			config: &NATSConfig{
				URIs:    "nats://localhost:4222",
				Stream:  "TRADE",
				Subject: "trade.test",
			},
			expectError: false,
		},
		{
			name: "valid multiple URIs",
			config: &NATSConfig{
				URIs:    "nats://localhost:4222,nats://localhost:4223",
				Stream:  "TRADE",
				Subject: "trade.test",
			},
			expectError: false,
		},
		{
			name: "empty URIs",
			config: &NATSConfig{
				URIs:    "",
				Stream:  "TRADE",
				Subject: "trade.test",
			},
			expectError: true,
			errorMsg:    "nats.uris cannot be empty",
		},
		{
			name: "empty stream",
			config: &NATSConfig{
				URIs:    "nats://localhost:4222",
				Stream:  "",
				Subject: "trade.test",
			},
			expectError: true,
			errorMsg:    "nats.stream cannot be empty",
		},
		{
			name: "empty subject",
			config: &NATSConfig{
				URIs:    "nats://localhost:4222",
				Stream:  "TRADE",
				Subject: "",
			},
			expectError: true,
			errorMsg:    "nats.subject cannot be empty",
		},
		{
			name: "invalid URI scheme",
			config: &NATSConfig{
				URIs:    "http://localhost:4222",
				Stream:  "TRADE",
				Subject: "trade.test",
			},
			expectError: true,
			errorMsg:    "invalid NATS URI scheme",
		},
		{
			name: "invalid URI format",
			config: &NATSConfig{
				URIs:    "nats://[invalid",
				Stream:  "TRADE",
				Subject: "trade.test",
			},
			expectError: true,
			errorMsg:    "invalid NATS URI at index 0",
		},
		{
			name: "empty hostname",
			config: &NATSConfig{
				URIs:    "nats://:4222",
				Stream:  "TRADE",
				Subject: "trade.test",
			},
			expectError: true,
			errorMsg:    "hostname cannot be empty",
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
