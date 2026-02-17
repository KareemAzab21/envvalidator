package envvalidator

import (
	"fmt"
	"reflect"
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

	// failFast stops on first error if true, collects all errors if false
	// When true: returns immediately after first validation error
	// When false: collects all validation errors and returns them together
	failFast bool
}

// New creates a new Validator with the given options
// This is useful when you want to reuse the same validator configuration
// for multiple Load operations.
//
// Example:
//
//	v := envvalidator.New(
//	    envvalidator.WithPrefix("APP_"),
//	    envvalidator.WithRequiredByDefault(true),
//	)
//	v.Load(&cfg1)
//	v.Load(&cfg2)
func New(opts ...Option) *Validator {
	v := &Validator{
		prefix:            "",
		requiredByDefault: false,
		failFast:          false,
	}

	// Apply all options
	for _, opt := range opts {
		opt(v)
	}

	return v
}

// Load loads environment variables into the provided config struct and validates them
//
// The config parameter must be a pointer to a struct with appropriate tags:
//   - `env:"VAR_NAME"` - specifies the environment variable name
//   - `validate:"rule1,rule2"` - specifies validation rules (comma-separated)
//   - `default:"value"` - specifies default value if env var not set
//
// Supported validation rules:
//   - required: field must have a value
//   - min:N: minimum string length or numeric value
//   - max:N: maximum string length or numeric value
//   - range:MIN-MAX: numeric value must be in range
//   - minvalue:N: minimum numeric value
//   - maxvalue:N: maximum numeric value
//   - oneof:val1 val2 val3: value must be one of the listed values
//   - email: must be valid email format
//   - url: must be valid URL format (http:// or https://)
//
// Supported types:
//   - string
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - bool
//   - float32, float64
//   - time.Duration
//   - []string (comma-separated values)
//
// Example:
//
//	type Config struct {
//	    Port        int    `env:"PORT" default:"8080" validate:"range:1000-9999"`
//	    Host        string `env:"HOST" default:"localhost" validate:"min:1"`
//	    Environment string `env:"ENV" default:"dev" validate:"oneof:dev staging prod"`
//	    APIKey      string `env:"API_KEY" validate:"required,min:32"`
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
//
// This will look for APP_PORT, APP_HOST, etc. and make all fields required by default.
func Load(config interface{}, opts ...Option) error {
	// Create a new validator with the given options
	v := New(opts...)

	// Validate that config is a pointer to a struct
	if err := v.validateInput(config); err != nil {
		return err
	}

	// Load environment variables into the struct
	if err := v.load(config); err != nil {
		return err
	}

	// Validate the loaded values
	if err := v.validate(config); err != nil {
		return err
	}

	return nil
}

// MustLoad is like Load but panics if an error occurs
// Useful for initialization code where you want to fail fast
//
// This is typically used in main() or init() functions where
// you want the application to crash immediately if configuration
// is invalid, rather than handling the error.
//
// Example:
//
//	func main() {
//	    var cfg Config
//	    envvalidator.MustLoad(&cfg)
//	    // If this succeeds, cfg is guaranteed to be valid
//	    // If it fails, the application will panic with a descriptive message
//
//	    startServer(cfg)
//	}
//
// The panic message will include the full error details, making it
// easy to diagnose configuration issues during startup.
func MustLoad(config interface{}, opts ...Option) {
	if err := Load(config, opts...); err != nil {
		panic(fmt.Sprintf("envvalidator: failed to load config: %v", err))
	}
}

// validateInput checks that the config parameter is valid
// It ensures config is a non-nil pointer to a struct
func (v *Validator) validateInput(config interface{}) error {
	// Check if config is nil
	if config == nil {
		return ErrNotStructPointer
	}

	// Get the reflection value
	val := reflect.ValueOf(config)

	// Must be a pointer
	if val.Kind() != reflect.Ptr {
		return ErrNotStructPointer
	}

	// Must not be a nil pointer
	if val.IsNil() {
		return ErrNotStructPointer
	}

	// The pointer must point to a struct
	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	return nil
}
