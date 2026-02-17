package validators

import (
	"fmt"
	"regexp"
	"strings"
)

// validateMin validates minimum string length
// Usage: validate:"min:5"
func validateMin(value interface{}, param string) error {
	str, ok := value.(string)
	if !ok {
		return nil // Skip non-string types
	}

	var minLen int
	if _, err := fmt.Sscanf(param, "%d", &minLen); err != nil {
		return fmt.Errorf("invalid min parameter: %s", param)
	}

	if len(str) < minLen {
		return fmt.Errorf("length must be at least %d characters (got %d)", minLen, len(str))
	}

	return nil
}

// validateMax validates maximum string length
// Usage: validate:"max:100"
func validateMax(value interface{}, param string) error {
	str, ok := value.(string)
	if !ok {
		return nil // Skip non-string types
	}

	var maxLen int
	if _, err := fmt.Sscanf(param, "%d", &maxLen); err != nil {
		return fmt.Errorf("invalid max parameter: %s", param)
	}

	if len(str) > maxLen {
		return fmt.Errorf("length must be at most %d characters (got %d)", maxLen, len(str))
	}

	return nil
}

// validateOneOf validates that value is one of the allowed values
// Usage: validate:"oneof:dev staging prod"
func validateOneOf(value interface{}, param string) error {
	str, ok := value.(string)
	if !ok {
		return nil // Skip non-string types
	}

	// Split allowed values by space
	allowedValues := strings.Fields(param)
	if len(allowedValues) == 0 {
		return fmt.Errorf("oneof requires at least one allowed value")
	}

	// Check if value is in allowed list
	for _, allowed := range allowedValues {
		if str == allowed {
			return nil
		}
	}

	return fmt.Errorf("must be one of [%s], got %q", strings.Join(allowedValues, ", "), str)
}

// validateEmail validates email format (basic validation)
// Usage: validate:"email"
func validateEmail(value interface{}, param string) error {
	str, ok := value.(string)
	if !ok {
		return nil // Skip non-string types
	}

	// Skip empty strings (use required for that)
	if str == "" {
		return nil
	}

	// Basic email regex pattern
	// This is a simple pattern - for production, consider using a library
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(str) {
		return fmt.Errorf("invalid email format: %q", str)
	}

	return nil
}

// validateURL validates URL format (basic validation)
// Usage: validate:"url"
func validateURL(value interface{}, param string) error {
	str, ok := value.(string)
	if !ok {
		return nil // Skip non-string types
	}

	// Skip empty strings
	if str == "" {
		return nil
	}

	// Basic URL validation - must start with http:// or https://
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=]+$`)

	if !urlRegex.MatchString(str) {
		return fmt.Errorf("invalid URL format: %q (must start with http:// or https://)", str)
	}

	return nil
}
