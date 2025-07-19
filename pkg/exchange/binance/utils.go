package binance

import (
	"net/url"
	"strconv"
	"time"
)

// GetCurrentTimestamp returns the current Unix timestamp in milliseconds
func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// FormatTimestamp formats a time.Time to Unix timestamp in milliseconds
func FormatTimestamp(t time.Time) string {
	return strconv.FormatInt(t.UnixNano()/int64(time.Millisecond), 10)
}

// ParseTimestamp parses a Unix timestamp in milliseconds to time.Time
func ParseTimestamp(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}

// BuildQueryString builds a query string from parameters
func BuildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	return values.Encode()
}

// MergeQueryParams merges multiple url.Values into one
func MergeQueryParams(params ...url.Values) url.Values {
	result := url.Values{}

	for _, p := range params {
		for key, values := range p {
			for _, value := range values {
				result.Add(key, value)
			}
		}
	}

	return result
}

// StringToFloat64 safely converts string to float64, returns 0 on error
func StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// StringToInt64 safely converts string to int64, returns 0 on error
func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// Float64ToString converts float64 to string with proper precision
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// Int64ToString converts int64 to string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// IsValidSymbol checks if a symbol string is valid (non-empty and uppercase)
func IsValidSymbol(symbol string) bool {
	if len(symbol) == 0 {
		return false
	}

	// Check if all characters are uppercase letters
	for _, r := range symbol {
		if r < 'A' || r > 'Z' {
			return false
		}
	}

	return true
}

// NormalizeSymbol normalizes a symbol to uppercase
func NormalizeSymbol(symbol string) string {
	// Simple uppercase conversion
	result := ""
	for _, r := range symbol {
		if r >= 'a' && r <= 'z' {
			result += string(r - 32) // Convert to uppercase
		} else {
			result += string(r)
		}
	}
	return result
}

// ValidateInterval checks if the given interval is valid
func ValidateInterval(interval string) bool {
	validIntervals := []string{
		Interval1m, Interval3m, Interval5m, Interval15m, Interval30m,
		Interval1h, Interval2h, Interval4h, Interval6h, Interval8h, Interval12h,
		Interval1d, Interval3d, Interval1w, Interval1M,
	}

	for _, valid := range validIntervals {
		if interval == valid {
			return true
		}
	}

	return false
}

// ValidateOrderType checks if the given order type is valid
func ValidateOrderType(orderType string) bool {
	validTypes := []string{
		OrderTypeLimit, OrderTypeMarket, OrderTypeStopLoss,
		OrderTypeStopLossLimit, OrderTypeTakeProfit, OrderTypeTakeProfitLimit,
		OrderTypeLimitMaker,
	}

	for _, valid := range validTypes {
		if orderType == valid {
			return true
		}
	}

	return false
}

// ValidateTimeInForce checks if the given time in force is valid
func ValidateTimeInForce(timeInForce string) bool {
	validTifs := []string{TimeInForceGTC, TimeInForceIOC, TimeInForceFOK}

	for _, valid := range validTifs {
		if timeInForce == valid {
			return true
		}
	}

	return false
}

// ValidateSide checks if the given order side is valid
func ValidateSide(side string) bool {
	return side == SideBuy || side == SideSell
}
