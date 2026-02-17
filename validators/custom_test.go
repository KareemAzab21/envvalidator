package validators

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestCustomValidators(t *testing.T) {
	t.Run("register and use custom validator", func(t *testing.T) {
		// Register a custom validator that checks if string is uppercase
		Register("uppercase", func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			if str != strings.ToUpper(str) {
				return fmt.Errorf("value must be uppercase, got %q", str)
			}
			return nil
		})

		// Test valid uppercase
		val := reflect.ValueOf("HELLO")
		rules := []ValidationRule{{Name: "uppercase", Param: ""}}
		err := ApplyValidators(val, rules)
		if err != nil {
			t.Errorf("unexpected error for uppercase string: %v", err)
		}

		// Test invalid lowercase
		val = reflect.ValueOf("hello")
		err = ApplyValidators(val, rules)
		if err == nil {
			t.Error("expected error for lowercase string")
		}
	})

	t.Run("custom validator with parameter", func(t *testing.T) {
		// Register a validator that checks if string starts with a prefix
		Register("startswith", func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			if !strings.HasPrefix(str, param) {
				return fmt.Errorf("value must start with %q, got %q", param, str)
			}
			return nil
		})

		// Test valid prefix
		val := reflect.ValueOf("hello world")
		rules := []ValidationRule{{Name: "startswith", Param: "hello"}}
		err := ApplyValidators(val, rules)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Test invalid prefix
		val = reflect.ValueOf("world hello")
		err = ApplyValidators(val, rules)
		if err == nil {
			t.Error("expected error for wrong prefix")
		}
	})

	t.Run("custom numeric validator", func(t *testing.T) {
		// Register a validator that checks if number is even
		Register("even", func(value interface{}, param string) error {
			num, err := toInt64(value)
			if err != nil {
				return nil // Skip non-numeric
			}
			if num%2 != 0 {
				return fmt.Errorf("value must be even, got %d", num)
			}
			return nil
		})

		// Test even number
		val := reflect.ValueOf(42)
		rules := []ValidationRule{{Name: "even", Param: ""}}
		err := ApplyValidators(val, rules)
		if err != nil {
			t.Errorf("unexpected error for even number: %v", err)
		}

		// Test odd number
		val = reflect.ValueOf(43)
		err = ApplyValidators(val, rules)
		if err == nil {
			t.Error("expected error for odd number")
		}
	})

	t.Run("custom validator - alphanumeric", func(t *testing.T) {
		Register("alphanumeric", func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			for _, r := range str {
				if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
					return fmt.Errorf("value must be alphanumeric, got %q", str)
				}
			}
			return nil
		})

		// Valid
		val := reflect.ValueOf("Hello123")
		rules := []ValidationRule{{Name: "alphanumeric", Param: ""}}
		err := ApplyValidators(val, rules)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Invalid - contains special chars
		val = reflect.ValueOf("Hello@123")
		err = ApplyValidators(val, rules)
		if err == nil {
			t.Error("expected error for non-alphanumeric string")
		}
	})

	t.Run("custom validator - contains", func(t *testing.T) {
		Register("contains", func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			if !strings.Contains(str, param) {
				return fmt.Errorf("value must contain %q, got %q", param, str)
			}
			return nil
		})

		// Valid
		val := reflect.ValueOf("hello world")
		rules := []ValidationRule{{Name: "contains", Param: "world"}}
		err := ApplyValidators(val, rules)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Invalid
		val = reflect.ValueOf("hello there")
		err = ApplyValidators(val, rules)
		if err == nil {
			t.Error("expected error when substring not found")
		}
	})

	t.Run("custom validator - length equals", func(t *testing.T) {
		Register("len", func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			var expectedLen int
			if _, err := fmt.Sscanf(param, "%d", &expectedLen); err != nil {
				return fmt.Errorf("invalid len parameter: %s", param)
			}
			if len(str) != expectedLen {
				return fmt.Errorf("length must be exactly %d, got %d", expectedLen, len(str))
			}
			return nil
		})

		// Valid
		val := reflect.ValueOf("hello")
		rules := []ValidationRule{{Name: "len", Param: "5"}}
		err := ApplyValidators(val, rules)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Invalid
		val = reflect.ValueOf("hi")
		err = ApplyValidators(val, rules)
		if err == nil {
			t.Error("expected error for wrong length")
		}
	})

	t.Run("custom validator - divisible by", func(t *testing.T) {
		Register("divisibleby", func(value interface{}, param string) error {
			num, err := toInt64(value)
			if err != nil {
				return nil
			}
			var divisor int64
			if _, err := fmt.Sscanf(param, "%d", &divisor); err != nil {
				return fmt.Errorf("invalid divisibleby parameter: %s", param)
			}
			if divisor == 0 {
				return fmt.Errorf("divisor cannot be zero")
			}
			if num%divisor != 0 {
				return fmt.Errorf("value must be divisible by %d, got %d", divisor, num)
			}
			return nil
		})

		// Valid - 10 is divisible by 5
		val := reflect.ValueOf(10)
		rules := []ValidationRule{{Name: "divisibleby", Param: "5"}}
		err := ApplyValidators(val, rules)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Invalid - 11 is not divisible by 5
		val = reflect.ValueOf(11)
		err = ApplyValidators(val, rules)
		if err == nil {
			t.Error("expected error for non-divisible number")
		}
	})

	t.Run("get non-existent validator", func(t *testing.T) {
		_, ok := Get("nonexistent_validator_xyz")
		if ok {
			t.Error("expected false for non-existent validator")
		}
	})

	t.Run("register nil validator", func(t *testing.T) {
		// This should not panic
		Register("nil_validator", nil)

		fn, ok := Get("nil_validator")
		if !ok {
			t.Error("expected to find registered validator")
		}
		if fn != nil {
			t.Error("expected nil validator function")
		}
	})
}

func TestCustomValidatorIntegration(t *testing.T) {
	t.Run("multiple custom validators on same field", func(t *testing.T) {
		// Register validators
		Register("notempty", func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			if strings.TrimSpace(str) == "" {
				return fmt.Errorf("value cannot be empty")
			}
			return nil
		})

		Register("lowercase", func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			if str != strings.ToLower(str) {
				return fmt.Errorf("value must be lowercase")
			}
			return nil
		})

		// Test with both validators
		val := reflect.ValueOf("hello")
		rules := []ValidationRule{
			{Name: "notempty", Param: ""},
			{Name: "lowercase", Param: ""},
		}
		err := ApplyValidators(val, rules)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Test failure on second validator
		val = reflect.ValueOf("HELLO")
		err = ApplyValidators(val, rules)
		if err == nil {
			t.Error("expected error for uppercase string")
		}
	})
}
