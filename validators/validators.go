package validators

import (
	"fmt"
	"reflect"
)

// ValidatorFunc is a function that validates a field value
// It receives the field value and the validation parameter (if any)
// Returns nil if valid, error if invalid
type ValidatorFunc func(value interface{}, param string) error

// Registry holds all registered validators
var registry = make(map[string]ValidatorFunc)

// Register adds a validator to the registry
func Register(name string, fn ValidatorFunc) {
	registry[name] = fn
}

// Get retrieves a validator from the registry
func Get(name string) (ValidatorFunc, bool) {
	fn, ok := registry[name]
	return fn, ok
}

// ApplyValidators applies a list of validation rules to a field value
func ApplyValidators(fieldValue reflect.Value, rules []ValidationRule) error {
	// Get the actual value as interface{}
	value := fieldValue.Interface()

	for _, rule := range rules {
		// Skip "required" - it's handled during loading
		if rule.Name == "required" {
			continue
		}

		// Get the validator function
		validatorFn, ok := Get(rule.Name)
		if !ok {
			return fmt.Errorf("unknown validator: %s", rule.Name)
		}

		// Apply the validator
		if err := validatorFn(value, rule.Param); err != nil {
			return err
		}
	}

	return nil
}

// ValidationRule represents a single validation rule
type ValidationRule struct {
	Name  string // Rule name (e.g., "min", "max", "email")
	Param string // Rule parameter (e.g., "5" for "min:5")
}

// init registers all built-in validators
func init() {
	// String validators
	Register("min", validateMin)
	Register("max", validateMax)
	Register("oneof", validateOneOf)
	Register("email", validateEmail)
	Register("url", validateURL)

	// Numeric validators
	Register("range", validateRange)
	Register("minvalue", validateMinValue)
	Register("maxvalue", validateMaxValue)
}
