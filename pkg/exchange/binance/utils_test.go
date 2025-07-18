package binance

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMustParseDecimal(t *testing.T) {
	t.Run("ValidDecimal", func(t *testing.T) {
		result := mustParseDecimal("123.456")
		expected := decimal.NewFromFloat(123.456)
		assert.True(t, result.Equal(expected))
	})

	t.Run("ValidInteger", func(t *testing.T) {
		result := mustParseDecimal("789")
		expected := decimal.NewFromInt(789)
		assert.True(t, result.Equal(expected))
	})

	t.Run("ValidZero", func(t *testing.T) {
		result := mustParseDecimal("0")
		assert.True(t, result.IsZero())
	})

	t.Run("ValidScientificNotation", func(t *testing.T) {
		result := mustParseDecimal("1.23e-4")
		expected := decimal.NewFromFloat(0.000123)
		assert.True(t, result.Equal(expected))
	})

	t.Run("InvalidDecimal", func(t *testing.T) {
		// Should return zero instead of panicking
		result := mustParseDecimal("invalid")
		assert.True(t, result.IsZero())
	})

	t.Run("EmptyString", func(t *testing.T) {
		result := mustParseDecimal("")
		assert.True(t, result.IsZero())
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		result := mustParseDecimal("-42.5")
		expected := decimal.NewFromFloat(-42.5)
		assert.True(t, result.Equal(expected))
	})
}

func TestParseDecimal(t *testing.T) {
	t.Run("ValidDecimal", func(t *testing.T) {
		result, err := parseDecimal("123.456")
		assert.NoError(t, err)
		expected := decimal.NewFromFloat(123.456)
		assert.True(t, result.Equal(expected))
	})

	t.Run("ValidInteger", func(t *testing.T) {
		result, err := parseDecimal("789")
		assert.NoError(t, err)
		expected := decimal.NewFromInt(789)
		assert.True(t, result.Equal(expected))
	})

	t.Run("ValidZero", func(t *testing.T) {
		result, err := parseDecimal("0")
		assert.NoError(t, err)
		assert.True(t, result.IsZero())
	})

	t.Run("ValidScientificNotation", func(t *testing.T) {
		result, err := parseDecimal("1.23e-4")
		assert.NoError(t, err)
		expected := decimal.NewFromFloat(0.000123)
		assert.True(t, result.Equal(expected))
	})

	t.Run("InvalidDecimal", func(t *testing.T) {
		result, err := parseDecimal("invalid")
		assert.Error(t, err)
		assert.True(t, result.IsZero())
	})

	t.Run("EmptyString", func(t *testing.T) {
		result, err := parseDecimal("")
		assert.Error(t, err)
		assert.True(t, result.IsZero())
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		result, err := parseDecimal("-42.5")
		assert.NoError(t, err)
		expected := decimal.NewFromFloat(-42.5)
		assert.True(t, result.Equal(expected))
	})

	t.Run("LargeNumber", func(t *testing.T) {
		result, err := parseDecimal("999999999999999999.999999999999999999")
		assert.NoError(t, err)
		assert.True(t, result.IsPositive())
	})
}
