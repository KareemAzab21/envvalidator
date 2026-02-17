package envvalidator

// Option is a function that modifies a Validator
// This is the functional options pattern
type Option func(*Validator)

// WithPrefix adds a prefix to all environment variable names
// Example: WithPrefix("APP_") will look for APP_PORT instead of PORT
//
// Usage:
//   Load(&cfg, WithPrefix("MYAPP_"))
//
// If struct has `env:"PORT"`, it will look for MYAPP_PORT
func WithPrefix(prefix string) Option {
	return func(v *Validator) {
		v.prefix = prefix
	}
}

// WithRequiredByDefault sets whether fields without "required" tag are required
// Default is false (fields are optional unless marked as required)
//
// Usage:
//   Load(&cfg, WithRequiredByDefault(true))
//
// When true: all fields must have values unless tagged with "optional"
// When false: fields are optional unless tagged with "required"
func WithRequiredByDefault(required bool) Option {
	return func(v *Validator) {
		v.requiredByDefault = required
	}
}

// WithCustomValidator registers a custom validation function
// This will be implemented in Phase 4, but we define it here for completeness
//
// Usage:
//   Load(&cfg, WithCustomValidator("custom", myValidatorFunc))
func WithCustomValidator(name string, fn ValidatorFunc) Option {
	return func(v *Validator) {
		if v.customValidators == nil {
			v.customValidators = make(map[string]ValidatorFunc)
		}
		v.customValidators[name] = fn
	}
}

// WithFailFast stops validation on first error instead of collecting all errors
// Default is false (collect all errors)
//
// Usage:
//   Load(&cfg, WithFailFast(true))
func WithFailFast(failFast bool) Option {
	return func(v *Validator) {
		v.failFast = failFast
	}
}

// ValidatorFunc is a function that validates a value
// Will be fully implemented in Phase 3
type ValidatorFunc func(value interface{}, param string) error
