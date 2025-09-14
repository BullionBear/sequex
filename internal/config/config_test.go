package config

import (
	"os"
	"strings"
	"testing"
)

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
