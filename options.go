package envvalidator

import (
	"github.com/KareemAzab21/envvalidator/validators"
)

// Option is a function that modifies a Validator
// This is the functional options pattern, allowing flexible configuration
type Option func(*Validator)

// WithPrefix adds a prefix to all environment variable names
// This is useful when you want to namespace your environment variables
// to avoid conflicts with other applications or system variables.
//
// Example: WithPrefix("APP_") will look for APP_PORT instead of PORT
//
// Usage:
//
//	type Config struct {
//	    Port int `env:"PORT"`
//	}
//	Load(&cfg, WithPrefix("MYAPP_"))
//
// If struct has `env:"PORT"`, it will look for MYAPP_PORT in the environment.
// The prefix is automatically added to all env tags in the struct.
func WithPrefix(prefix string) Option {
	return func(v *Validator) {
		v.prefix = prefix
	}
}

// WithRequiredByDefault sets whether fields without explicit "required" tag are required
// Default behavior is false (fields are optional unless marked as required)
//
// When true: All fields must have values unless explicitly tagged with "optional"
// When false: Fields are optional unless explicitly tagged with "required"
//
// Usage:
//
//	type Config struct {
//	    Port int    `env:"PORT"`           // Required when WithRequiredByDefault(true)
//	    Host string `env:"HOST"`           // Required when WithRequiredByDefault(true)
//	    Name string `env:"NAME" optional:""` // Always optional (not implemented yet)
//	}
//	Load(&cfg, WithRequiredByDefault(true))
//
// This is useful for strict configuration where you want to ensure
// all fields are explicitly set in the environment.
func WithRequiredByDefault(required bool) Option {
	return func(v *Validator) {
		v.requiredByDefault = required
	}
}

// WithFailFast stops validation on the first error instead of collecting all errors
// Default behavior is false (collect all validation errors)
//
// When true: Returns immediately after encountering the first validation error
// When false: Collects all validation errors and returns them together
//
// Usage:
//
//	type Config struct {
//	    Port int    `env:"PORT" validate:"range:1000-9999"`
//	    Host string `env:"HOST" validate:"required,min:3"`
//	}
//	Load(&cfg, WithFailFast(true))
//
// With fail-fast enabled, if PORT validation fails, HOST won't be validated.
// Without fail-fast, you'll get all validation errors at once, which is better
// for user experience as they can fix all issues in one go.
func WithFailFast(failFast bool) Option {
	return func(v *Validator) {
		v.failFast = failFast
	}
}

// WithCustomValidator registers a custom validation function
// This allows you to add your own validation rules beyond the built-in ones.
//
// The validator function receives:
//   - value: The field value to validate (as interface{})
//   - param: The parameter from the validation tag (e.g., "5" in "validate:myvalidator:5")
//
// The function should return nil if validation passes, or an error describing
// why validation failed.
//
// Usage:
//
//	// Define a custom validator
//	uppercaseValidator := func(value interface{}, param string) error {
//	    str, ok := value.(string)
//	    if !ok {
//	        return nil // Skip non-string types
//	    }
//	    if str != strings.ToUpper(str) {
//	        return fmt.Errorf("value must be uppercase")
//	    }
//	    return nil
//	}
//
//	// Register and use it
//	type Config struct {
//	    Code string `env:"CODE" validate:"uppercase"`
//	}
//	Load(&cfg, WithCustomValidator("uppercase", uppercaseValidator))
//
// Custom validators are registered globally and can be used across
// multiple Load() calls. They can also override built-in validators
// if you register them with the same name.
func WithCustomValidator(name string, fn validators.ValidatorFunc) Option {
	return func(v *Validator) {
		// Register the validator in the global registry
		validators.Register(name, fn)
	}
}

// WithCustomValidators registers multiple custom validation functions at once
// This is a convenience function for registering many validators in one call.
//
// Usage:
//
//	customValidators := map[string]validators.ValidatorFunc{
//	    "uppercase": func(value interface{}, param string) error {
//	        str, ok := value.(string)
//	        if !ok {
//	            return nil
//	        }
//	        if str != strings.ToUpper(str) {
//	            return fmt.Errorf("must be uppercase")
//	        }
//	        return nil
//	    },
//	    "lowercase": func(value interface{}, param string) error {
//	        str, ok := value.(string)
//	        if !ok {
//	            return nil
//	        }
//	        if str != strings.ToLower(str) {
//	            return fmt.Errorf("must be lowercase")
//	        }
//	        return nil
//	    },
//	}
//
//	type Config struct {
//	    Code string `env:"CODE" validate:"uppercase"`
//	    Name string `env:"NAME" validate:"lowercase"`
//	}
//	Load(&cfg, WithCustomValidators(customValidators))
//
// This is equivalent to calling WithCustomValidator() multiple times,
// but more convenient when you have many custom validators to register.
func WithCustomValidators(customValidators map[string]validators.ValidatorFunc) Option {
	return func(v *Validator) {
		// Register all validators in the global registry
		for name, fn := range customValidators {
			validators.Register(name, fn)
		}
	}
}

// CombineOptions combines multiple options into a single option
// This is useful when you have a standard set of options you want to reuse.
//
// Usage:
//
//	// Define a standard configuration
//	standardOpts := CombineOptions(
//	    WithPrefix("APP_"),
//	    WithRequiredByDefault(true),
//	    WithFailFast(false),
//	)
//
//	// Use it
//	Load(&cfg, standardOpts)
//
//	// Or combine with additional options
//	Load(&cfg, standardOpts, WithCustomValidator("myvalidator", myFunc))
//
// This helps maintain consistency across your application when loading
// different configuration structs.
func CombineOptions(opts ...Option) Option {
	return func(v *Validator) {
		for _, opt := range opts {
			opt(v)
		}
	}
}
