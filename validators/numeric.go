package validators

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// validateRange validates that a numeric value is within a range
// Usage: validate:"range:1000-9999"
func validateRange(value interface{}, param string) error {
	// Convert value to int64
	numValue, err := toInt64(value)
	if err != nil {
		return nil // Skip non-numeric types
	}

	// Parse range parameter (e.g., "1000-9999")
	parts := strings.Split(param, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range format: %s (expected format: min-max)", param)
	}

	min, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid range minimum: %s", parts[0])
	}

	max, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid range maximum: %s", parts[1])
	}

	if numValue < min || numValue > max {
		return fmt.Errorf("value must be between %d and %d (got %d)", min, max, numValue)
	}

	return nil
}

// validateMinValue validates minimum numeric value
// Usage: validate:"minvalue:0"
func validateMinValue(value interface{}, param string) error {
	numValue, err := toInt64(value)
	if err != nil {
		return nil // Skip non-numeric types
	}

	minValue, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid minvalue parameter: %s", param)
	}

	if numValue < minValue {
		return fmt.Errorf("value must be at least %d (got %d)", minValue, numValue)
	}

	return nil
}

// validateMaxValue validates maximum numeric value
// Usage: validate:"maxvalue:100"
func validateMaxValue(value interface{}, param string) error {
	numValue, err := toInt64(value)
	if err != nil {
		return nil // Skip non-numeric types
	}

	maxValue, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid maxvalue parameter: %s", param)
	}

	if numValue > maxValue {
		return fmt.Errorf("value must be at most %d (got %d)", maxValue, numValue)
	}

	return nil
}

// toInt64 converts various numeric types to int64
func toInt64(value interface{}) (int64, error) {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(val.Float()), nil
	default:
		return 0, fmt.Errorf("not a numeric type: %v", val.Kind())
	}
}
