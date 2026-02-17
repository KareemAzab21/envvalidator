// Package envvalidator provides a powerful, type-safe way to load and validate
// environment variables into Go structs.
//
// # Features
//
//   - Type-safe environment variable loading
//   - Comprehensive validation rules
//   - Custom validator support
//   - Default values
//   - Detailed error messages
//   - Zero dependencies
//
// # Quick Start
//
// Define a configuration struct with tags:
//
//	type Config struct {
//	    Port        int    `env:"PORT" default:"8080" validate:"range:1000-9999"`
//	    Host        string `env:"HOST" default:"localhost"`
//	    Environment string `env:"ENV" default:"dev" validate:"oneof:dev staging prod"`
//	    APIKey      string `env:"API_KEY" validate:"required,min:32"`
//	}
//
// Load and validate:
//
//	var cfg Config
//	if err := envvalidator.Load(&cfg); err != nil {
//	    log.Fatal(err)
//	}
//
// # Supported Types
//
// The following types are supported:
//
//   - string
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - bool
//   - float32, float64
//   - time.Duration
//   - []string (comma-separated)
//
// # Struct Tags
//
// Three tags are supported:
//
//   - env: Specifies the environment variable name
//   - default: Provides a default value if not set
//   - validate: Comma-separated list of validation rules
//
// Example:
//
//	type Config struct {
//	    Port int `env:"PORT" default:"8080" validate:"range:1000-9999"`
//	}
//
// # Built-in Validators
//
// String validators:
//   - required: Field must have a value
//   - min:N: Minimum string length
//   - max:N: Maximum string length
//   - oneof:a b c: Value must be one of the listed options
//   - email: Must be valid email format
//   - url: Must be valid URL (http/https)
//
// Numeric validators:
//   - range:MIN-MAX: Value must be in range
//   - minvalue:N: Minimum numeric value
//   - maxvalue:N: Maximum numeric value
//
// # Options
//
// Customize behavior with options:
//
//	envvalidator.Load(&cfg,
//	    envvalidator.WithPrefix("APP_"),
//	    envvalidator.WithRequiredByDefault(true),
//	    envvalidator.WithFailFast(false),
//	)
//
// # Custom Validators
//
// Register custom validators:
//
//	uppercaseValidator := func(value interface{}, param string) error {
//	    str, ok := value.(string)
//	    if !ok {
//	        return nil
//	    }
//	    if str != strings.ToUpper(str) {
//	        return fmt.Errorf("must be uppercase")
//	    }
//	    return nil
//	}
//
//	envvalidator.Load(&cfg,
//	    envvalidator.WithCustomValidator("uppercase", uppercaseValidator),
//	)
//
// # Error Handling
//
// Handle validation errors:
//
//	err := envvalidator.Load(&cfg)
//	if err != nil {
//	    var validationErrs envvalidator.ValidationErrors
//	    if errors.As(err, &validationErrs) {
//	        for _, e := range validationErrs {
//	            fmt.Printf("Field: %s, Error: %v\n", e.Field, e.Err)
//	        }
//	    }
//	}
//
// Or use MustLoad to panic on error:
//
//	envvalidator.MustLoad(&cfg)
//
// # Examples
//
// See the examples directory for complete working examples:
//   - examples/basic: Simple usage
//   - examples/advanced: All features
//   - examples/custom-validator: Custom validators
package envvalidator
