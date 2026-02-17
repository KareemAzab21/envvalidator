package envvalidator

import (
	"fmt"
)

// Validator is the core struct that handles loading and validation
// It holds configuration and state for the validation process
type Validator struct {
	// prefix is prepended to all environment variable names
	// Example: if prefix is "APP_", it looks for APP_PORT instead of PORT
	prefix string

	// requiredByDefault determines if fields are required unless marked optional
	// false = fields optional by default (need "required" tag)
	// true = fields required by default (need "optional" tag)
	requiredByDefault bool

	// customValidators holds user-defined validation functions
	// Key is the validator name, value is the validation function
	customValidators map[string]ValidatorFunc

	// failFast stops on first error if true, collects all errors if false
	failFast bool
}

// Load loads environment variables into the provided config struct and validates them
//
// The config parameter must be a pointer to a struct with appropriate tags:
//   - `env:"VAR_NAME"` - specifies the environment variable name
//   - `validate:"rule1,rule2"` - specifies validation rules
//   - `default:"value"` - specifies default value if env var not set
//
// Example:
//
//	type Config struct {
//	    Port int `env:"PORT" validate:"required,range:1000-9999"`
//	    Host string `env:"HOST" default:"localhost"`
//	}
//
//	var cfg Config
//	if err := envvalidator.Load(&cfg); err != nil {
//	    log.Fatal(err)
//	}
//
// Options can be passed to customize behavior:
//
//	Load(&cfg, WithPrefix("APP_"), WithRequiredByDefault(true))
func Load(config interface{}, opts ...Option) error {
	// Create a new validator with default settings
	v := &Validator{
		prefix:            "",
		requiredByDefault: false,
		customValidators:  make(map[string]ValidatorFunc),
		failFast:          false,
	}

	// Apply all provided options
	// This modifies the validator based on user preferences
	for _, opt := range opts {
		opt(v)
	}

	// Phase 2: This will load env vars into the struct
	// For now, it's a placeholder
	if err := v.load(config); err != nil {
		return err
	}

	// Phase 3: This will validate the loaded values
	// For now, it's a placeholder
	if err := v.validate(config); err != nil {
		return err
	}

	return nil
}

// MustLoad is like Load but panics if an error occurs
// Useful for initialization code where you want to fail fast
//
// Example:
//
//	var cfg Config
//	envvalidator.MustLoad(&cfg)
//	// If this succeeds, cfg is guaranteed to be valid
func MustLoad(config interface{}, opts ...Option) {
	if err := Load(config, opts...); err != nil {
		panic(fmt.Sprintf("envvalidator: failed to load config: %v", err))
	}
}

// validate is an internal method that validates loaded values
// Will be implemented in Phase 3
func (v *Validator) validate(config interface{}) error {
	// TODO: Phase 3 - implement validation logic
	return nil
}
