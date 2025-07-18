package binance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError(t *testing.T) {
	t.Run("ErrorInterface", func(t *testing.T) {
		apiErr := &APIError{
			Code: -1002,
			Msg:  "Unauthorized",
		}

		// Test that it implements error interface
		var err error = apiErr
		assert.NotNil(t, err)

		expected := "Binance API error -1002: Unauthorized"
		assert.Equal(t, expected, apiErr.Error())
	})

	t.Run("DifferentErrorCodes", func(t *testing.T) {
		testCases := []struct {
			code     int
			msg      string
			expected string
		}{
			{-1000, "Unknown error", "Binance API error -1000: Unknown error"},
			{-1021, "Invalid timestamp", "Binance API error -1021: Invalid timestamp"},
			{-1022, "Invalid signature", "Binance API error -1022: Invalid signature"},
			{-2010, "NEW_ORDER_REJECTED", "Binance API error -2010: NEW_ORDER_REJECTED"},
		}

		for _, tc := range testCases {
			apiErr := &APIError{
				Code: tc.code,
				Msg:  tc.msg,
			}
			assert.Equal(t, tc.expected, apiErr.Error())
		}
	})
}

func TestIsRetryableError(t *testing.T) {
	t.Run("RetryableErrors", func(t *testing.T) {
		retryableErrors := []*APIError{
			{Code: ErrCodeTooManyRequests, Msg: "Too many requests"},
			{Code: ErrCodeTimeout, Msg: "Timeout"},
			{Code: ErrCodeDisconnected, Msg: "Disconnected"},
		}

		for _, apiErr := range retryableErrors {
			assert.True(t, IsRetryableError(apiErr), "Error %d should be retryable", apiErr.Code)
		}
	})

	t.Run("NonRetryableErrors", func(t *testing.T) {
		nonRetryableErrors := []*APIError{
			{Code: ErrCodeUnauthorized, Msg: "Unauthorized"},
			{Code: ErrCodeInvalidSignature, Msg: "Invalid signature"},
			{Code: ErrCodeInvalidTimestamp, Msg: "Invalid timestamp"},
			{Code: ErrCodeBadSymbol, Msg: "Invalid symbol"},
			{Code: ErrCodeInvalidOrderType, Msg: "Invalid order type"},
		}

		for _, apiErr := range nonRetryableErrors {
			assert.False(t, IsRetryableError(apiErr), "Error %d should not be retryable", apiErr.Code)
		}
	})

	t.Run("NonAPIError", func(t *testing.T) {
		regularErr := assert.AnError
		assert.False(t, IsRetryableError(regularErr))
	})

	t.Run("NilError", func(t *testing.T) {
		assert.False(t, IsRetryableError(nil))
	})

	t.Run("UnknownErrorCode", func(t *testing.T) {
		unknownErr := &APIError{
			Code: -9999,
			Msg:  "Unknown error code",
		}
		assert.False(t, IsRetryableError(unknownErr))
	})
}

func TestErrorConstants(t *testing.T) {
	// Test that all error constants are negative
	errorCodes := []int{
		ErrCodeUnknown,
		ErrCodeDisconnected,
		ErrCodeUnauthorized,
		ErrCodeTooManyRequests,
		ErrCodeUnexpectedResp,
		ErrCodeTimeout,
		ErrCodeUnknownOrderComposition,
		ErrCodeTooManyOrders,
		ErrCodeServiceShuttingDown,
		ErrCodeUnsupportedOperation,
		ErrCodeInvalidTimestamp,
		ErrCodeInvalidSignature,
		ErrCodeIllegalChars,
		ErrCodeTooManyParameters,
		ErrCodeMandatoryParamEmpty,
		ErrCodeUnknownParam,
		ErrCodeUnreadParameters,
		ErrCodeParamEmpty,
		ErrCodeParamNotRequired,
		ErrCodeNoDepth,
		ErrCodeTIFNotRequired,
		ErrCodeInvalidTIF,
		ErrCodeInvalidOrderType,
		ErrCodeInvalidSide,
		ErrCodeEmptyNewClOrdID,
		ErrCodeEmptyOrgClOrdID,
		ErrCodeBadInterval,
		ErrCodeBadSymbol,
		ErrCodeInvalidListenKey,
		ErrCodeMoreThanXXHours,
		ErrCodeOptionalParamsBadCombo,
		ErrCodeInvalidParameter,
	}

	for _, code := range errorCodes {
		assert.Negative(t, code, "Error code %d should be negative", code)
	}

	// Test that specific constants have expected values
	assert.Equal(t, -1000, ErrCodeUnknown)
	assert.Equal(t, -1003, ErrCodeTooManyRequests)
	assert.Equal(t, -1021, ErrCodeInvalidTimestamp)
	assert.Equal(t, -1022, ErrCodeInvalidSignature)
	assert.Equal(t, -1121, ErrCodeBadSymbol)
}
