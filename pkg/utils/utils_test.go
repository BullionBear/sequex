package utils

import "testing"

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "item exists in slice",
			slice:    []string{"binance", "okx", "bybit"},
			item:     "binance",
			expected: true,
		},
		{
			name:     "item does not exist in slice",
			slice:    []string{"binance", "okx", "bybit"},
			item:     "coinbase",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			item:     "binance",
			expected: false,
		},
		{
			name:     "single item slice - match",
			slice:    []string{"binance"},
			item:     "binance",
			expected: true,
		},
		{
			name:     "single item slice - no match",
			slice:    []string{"binance"},
			item:     "okx",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("Contains(%v, %s) = %v, expected %v", tt.slice, tt.item, result, tt.expected)
			}
		})
	}
}
