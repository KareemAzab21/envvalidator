// loader_test.go
package envvalidator

import (
	"errors"
	"os"
	"testing"
)

// Helper function to set and clean up environment variables
func setEnv(t *testing.T, key, value string) {
	t.Helper()
	os.Setenv(key, value)
	t.Cleanup(func() {
		os.Unsetenv(key)
	})
}

// Test loading string values
func TestLoadString(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "simple string",
			envValue: "hello",
			expected: "hello",
		},
		{
			name:     "string with spaces",
			envValue: "hello world",
			expected: "hello world",
		},
		{
			name:     "empty string",
			envValue: "",
			expected: "",
		},
		{
			name:     "string with special characters",
			envValue: "hello@world!123",
			expected: "hello@world!123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type Config struct {
				Value string `env:"TEST_STRING"`
			}

			if tt.envValue != "" {
				setEnv(t, "TEST_STRING", tt.envValue)
			}

			var cfg Config
			err := Load(&cfg)

			if tt.envValue == "" {
				// Empty string means env var not set, field should remain zero value
				if cfg.Value != "" {
					t.Errorf("expected empty string, got %q", cfg.Value)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Value != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, cfg.Value)
			}
		})
	}
}

// Test loading integer values
func TestLoadInt(t *testing.T) {
	tests := []struct {
		name      string
		envValue  string
		expected  int
		shouldErr bool
	}{
		{
			name:     "positive integer",
			envValue: "42",
			expected: 42,
		},
		{
			name:     "negative integer",
			envValue: "-42",
			expected: -42,
		},
		{
			name:     "zero",
			envValue: "0",
			expected: 0,
		},
		{
			name:      "invalid integer",
			envValue:  "not-a-number",
			shouldErr: true,
		},
		{
			name:      "float as integer",
			envValue:  "42.5",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type Config struct {
				Value int `env:"TEST_INT"`
			}

			setEnv(t, "TEST_INT", tt.envValue)

			var cfg Config
			err := Load(&cfg)

			if tt.shouldErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Value != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, cfg.Value)
			}
		})
	}
}

// Test loading different integer types
func TestLoadIntTypes(t *testing.T) {
	type Config struct {
		Int8Val  int8  `env:"INT8_VAL"`
		Int16Val int16 `env:"INT16_VAL"`
		Int32Val int32 `env:"INT32_VAL"`
		Int64Val int64 `env:"INT64_VAL"`
		UintVal  uint  `env:"UINT_VAL"`
	}

	setEnv(t, "INT8_VAL", "127")
	setEnv(t, "INT16_VAL", "32767")
	setEnv(t, "INT32_VAL", "2147483647")
	setEnv(t, "INT64_VAL", "9223372036854775807")
	setEnv(t, "UINT_VAL", "42")

	var cfg Config
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Int8Val != 127 {
		t.Errorf("expected 127, got %d", cfg.Int8Val)
	}
	if cfg.Int16Val != 32767 {
		t.Errorf("expected 32767, got %d", cfg.Int16Val)
	}
	if cfg.Int32Val != 2147483647 {
		t.Errorf("expected 2147483647, got %d", cfg.Int32Val)
	}
	if cfg.UintVal != 42 {
		t.Errorf("expected 42, got %d", cfg.UintVal)
	}
}

// Test loading boolean values
func TestLoadBool(t *testing.T) {
	tests := []struct {
		name      string
		envValue  string
		expected  bool
		shouldErr bool
	}{
		{name: "true lowercase", envValue: "true", expected: true},
		{name: "true uppercase", envValue: "TRUE", expected: true},
		{name: "true mixed case", envValue: "True", expected: true},
		{name: "true as 1", envValue: "1", expected: true},
		{name: "true as t", envValue: "t", expected: true},
		{name: "true as T", envValue: "T", expected: true},
		{name: "false lowercase", envValue: "false", expected: false},
		{name: "false uppercase", envValue: "FALSE", expected: false},
		{name: "false as 0", envValue: "0", expected: false},
		{name: "false as f", envValue: "f", expected: false},
		{name: "invalid bool", envValue: "yes", shouldErr: true},
		{name: "invalid bool number", envValue: "2", shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type Config struct {
				Value bool `env:"TEST_BOOL"`
			}

			setEnv(t, "TEST_BOOL", tt.envValue)

			var cfg Config
			err := Load(&cfg)

			if tt.shouldErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, cfg.Value)
			}
		})
	}
}

// Test loading float values
func TestLoadFloat(t *testing.T) {
	tests := []struct {
		name      string
		envValue  string
		expected  float64
		shouldErr bool
	}{
		{name: "positive float", envValue: "3.14", expected: 3.14},
		{name: "negative float", envValue: "-2.5", expected: -2.5},
		{name: "integer as float", envValue: "42", expected: 42.0},
		{name: "scientific notation", envValue: "1.23e10", expected: 1.23e10},
		{name: "invalid float", envValue: "not-a-float", shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type Config struct {
				Value float64 `env:"TEST_FLOAT"`
			}

			setEnv(t, "TEST_FLOAT", tt.envValue)

			var cfg Config
			err := Load(&cfg)

			if tt.shouldErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Value != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, cfg.Value)
			}
		})
	}
}

// Test required fields
func TestRequiredFields(t *testing.T) {
	t.Run("required field present", func(t *testing.T) {
		type Config struct {
			Required string `env:"REQUIRED_FIELD" validate:"required"`
		}

		setEnv(t, "REQUIRED_FIELD", "value")

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Required != "value" {
			t.Errorf("expected 'value', got %q", cfg.Required)
		}
	})

	t.Run("required field missing", func(t *testing.T) {
		type Config struct {
			Required string `env:"REQUIRED_FIELD" validate:"required"`
		}

		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected error for missing required field, got nil")
		}

		// Check if it's ValidationErrors (plural) or ValidationError (singular)
		var validationErrs ValidationErrors
		var validationErr *ValidationError

		if errors.As(err, &validationErrs) {
			// Multiple errors - get the first one
			if len(validationErrs) == 0 {
				t.Fatal("expected at least one validation error")
			}
			validationErr = &validationErrs[0]
		} else if errors.As(err, &validationErr) {
			// Single error - use it directly
		} else {
			t.Fatalf("expected ValidationError or ValidationErrors, got %T", err)
		}

		if !errors.Is(validationErr.Err, ErrRequired) {
			t.Errorf("expected ErrRequired, got %v", validationErr.Err)
		}

		if validationErr.Field != "Required" {
			t.Errorf("expected field name 'Required', got %q", validationErr.Field)
		}
	})

	t.Run("optional field missing", func(t *testing.T) {
		type Config struct {
			Optional string `env:"OPTIONAL_FIELD"`
		}

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error for optional field: %v", err)
		}

		if cfg.Optional != "" {
			t.Errorf("expected empty string, got %q", cfg.Optional)
		}
	})
}

// Test missing environment variables

func TestMissingEnvVars(t *testing.T) {
	t.Run("missing optional field", func(t *testing.T) {
		type Config struct {
			Optional string `env:"MISSING_VAR"`
		}

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should be zero value
		if cfg.Optional != "" {
			t.Errorf("expected empty string, got %q", cfg.Optional)
		}
	})

	t.Run("missing field with default", func(t *testing.T) {
		type Config struct {
			WithDefault string `env:"MISSING_VAR" default:"default-value"`
		}

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.WithDefault != "default-value" {
			t.Errorf("expected 'default-value', got %q", cfg.WithDefault)
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		type Config struct {
			Required string `env:"MISSING_REQUIRED" validate:"required"`
		}

		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected error for missing required field")
		}

		// The error should contain information about the required field
		// We check the error message contains "required"
		errMsg := err.Error()
		if !contains(errMsg, "required") && !contains(errMsg, "Required") {
			t.Errorf("expected error message to mention 'required', got: %v", err)
		}

		// Also verify it's a validation error type
		var validationErrs ValidationErrors
		var validationErr *ValidationError

		if !errors.As(err, &validationErrs) && !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError or ValidationErrors, got %T", err)
		}
	})
}

// Helper function to check if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Test with prefix option
func TestWithPrefix(t *testing.T) {
	t.Run("load with prefix", func(t *testing.T) {
		type Config struct {
			Port int    `env:"PORT"`
			Host string `env:"HOST"`
		}

		setEnv(t, "APP_PORT", "8080")
		setEnv(t, "APP_HOST", "localhost")

		var cfg Config
		err := Load(&cfg, WithPrefix("APP_"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Port != 8080 {
			t.Errorf("expected 8080, got %d", cfg.Port)
		}
		if cfg.Host != "localhost" {
			t.Errorf("expected 'localhost', got %q", cfg.Host)
		}
	})

	t.Run("prefix not applied to non-prefixed vars", func(t *testing.T) {
		type Config struct {
			Port int `env:"PORT"`
		}

		// Set without prefix
		setEnv(t, "PORT", "8080")

		var cfg Config
		// Try to load with prefix - should fail/be empty
		err := Load(&cfg, WithPrefix("APP_"))

		// Should not error (field is optional)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should be zero value since APP_PORT doesn't exist
		if cfg.Port != 0 {
			t.Errorf("expected 0, got %d", cfg.Port)
		}
	})
}

// Test default values
func TestDefaultValues(t *testing.T) {
	t.Run("use default when env var not set", func(t *testing.T) {
		type Config struct {
			Port  int    `env:"PORT" default:"8080"`
			Host  string `env:"HOST" default:"localhost"`
			Debug bool   `env:"DEBUG" default:"false"`
		}

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Port != 8080 {
			t.Errorf("expected 8080, got %d", cfg.Port)
		}
		if cfg.Host != "localhost" {
			t.Errorf("expected 'localhost', got %q", cfg.Host)
		}
		if cfg.Debug != false {
			t.Errorf("expected false, got %v", cfg.Debug)
		}
	})

	t.Run("env var overrides default", func(t *testing.T) {
		type Config struct {
			Port int `env:"PORT" default:"8080"`
		}

		setEnv(t, "PORT", "9000")

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Port != 9000 {
			t.Errorf("expected 9000 (env var), got %d", cfg.Port)
		}
	})

	t.Run("invalid default value", func(t *testing.T) {
		type Config struct {
			Port int `env:"PORT" default:"not-a-number"`
		}

		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected error for invalid default value")
		}
	})
}

// Test RequiredByDefault option
func TestRequiredByDefault(t *testing.T) {
	t.Run("all fields required by default", func(t *testing.T) {
		type Config struct {
			Port int    `env:"PORT"`
			Host string `env:"HOST"`
		}

		setEnv(t, "PORT", "8080")
		// HOST is missing

		var cfg Config
		err := Load(&cfg, WithRequiredByDefault(true))
		if err == nil {
			t.Fatal("expected error for missing required field")
		}

		// Should mention HOST field
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Field != "Host" {
				t.Errorf("expected error for 'Host' field, got %q", validationErr.Field)
			}
		}
	})

	t.Run("required by default with all fields set", func(t *testing.T) {
		type Config struct {
			Port int    `env:"PORT"`
			Host string `env:"HOST"`
		}

		setEnv(t, "PORT", "8080")
		setEnv(t, "HOST", "localhost")

		var cfg Config
		err := Load(&cfg, WithRequiredByDefault(true))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Port != 8080 {
			t.Errorf("expected 8080, got %d", cfg.Port)
		}
		if cfg.Host != "localhost" {
			t.Errorf("expected 'localhost', got %q", cfg.Host)
		}
	})
}

// Test multiple fields
func TestMultipleFields(t *testing.T) {
	type Config struct {
		AppName  string  `env:"APP_NAME" validate:"required"`
		Port     int     `env:"PORT" default:"8080"`
		Host     string  `env:"HOST" default:"localhost"`
		Debug    bool    `env:"DEBUG" default:"false"`
		Timeout  int     `env:"TIMEOUT"`
		MaxConns uint    `env:"MAX_CONNS" default:"100"`
		Rate     float64 `env:"RATE" default:"1.5"`
	}

	setEnv(t, "APP_NAME", "myapp")
	setEnv(t, "PORT", "9000")
	setEnv(t, "DEBUG", "true")
	setEnv(t, "MAX_CONNS", "200")

	var cfg Config
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.AppName != "myapp" {
		t.Errorf("AppName: expected 'myapp', got %q", cfg.AppName)
	}
	if cfg.Port != 9000 {
		t.Errorf("Port: expected 9000, got %d", cfg.Port)
	}
	if cfg.Host != "localhost" {
		t.Errorf("Host: expected 'localhost', got %q", cfg.Host)
	}
	if cfg.Debug != true {
		t.Errorf("Debug: expected true, got %v", cfg.Debug)
	}
	if cfg.Timeout != 0 {
		t.Errorf("Timeout: expected 0, got %d", cfg.Timeout)
	}
	if cfg.MaxConns != 200 {
		t.Errorf("MaxConns: expected 200, got %d", cfg.MaxConns)
	}
	if cfg.Rate != 1.5 {
		t.Errorf("Rate: expected 1.5, got %f", cfg.Rate)
	}
}

// Test error cases
func TestErrorCases(t *testing.T) {
	t.Run("not a pointer", func(t *testing.T) {
		type Config struct {
			Port int `env:"PORT"`
		}

		var cfg Config
		err := Load(cfg) // Not a pointer!
		if err == nil {
			t.Fatal("expected error for non-pointer config")
		}

		if !errors.Is(err, ErrNotStructPointer) {
			t.Errorf("expected ErrNotStructPointer, got %v", err)
		}
	})

	t.Run("not a struct", func(t *testing.T) {
		var notStruct int
		err := Load(&notStruct)
		if err == nil {
			t.Fatal("expected error for non-struct config")
		}

		if !errors.Is(err, ErrNotStruct) {
			t.Errorf("expected ErrNotStruct, got %v", err)
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		type Config struct {
			Port int `env:"PORT"`
		}

		var cfg *Config
		err := Load(cfg)
		if err == nil {
			t.Fatal("expected error for nil pointer")
		}
	})
}

// Test MustLoad
func TestMustLoad(t *testing.T) {
	t.Run("MustLoad success", func(t *testing.T) {
		type Config struct {
			Port int `env:"PORT" default:"8080"`
		}

		var cfg Config
		// Should not panic
		MustLoad(&cfg)

		if cfg.Port != 8080 {
			t.Errorf("expected 8080, got %d", cfg.Port)
		}
	})

	t.Run("MustLoad panic on error", func(t *testing.T) {
		type Config struct {
			Required string `env:"REQUIRED" validate:"required"`
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected MustLoad to panic, but it didn't")
			}
		}()

		var cfg Config
		MustLoad(&cfg) // Should panic because REQUIRED is missing
	})
}

// Test fields without env tag
func TestFieldsWithoutEnvTag(t *testing.T) {
	type Config struct {
		Port       int    `env:"PORT" default:"8080"`
		NoEnvTag   string // No env tag - should be ignored
		unexported int    `env:"UNEXPORTED"` // Unexported - should be ignored
	}

	setEnv(t, "PORT", "9000")
	setEnv(t, "UNEXPORTED", "123")

	var cfg Config
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 9000 {
		t.Errorf("expected 9000, got %d", cfg.Port)
	}
	if cfg.NoEnvTag != "" {
		t.Errorf("expected empty string for NoEnvTag, got %q", cfg.NoEnvTag)
	}
	if cfg.unexported != 0 {
		t.Errorf("expected 0 for unexported, got %d", cfg.unexported)
	}
}

// Test fail-fast option
func TestFailFast(t *testing.T) {
	t.Run("fail fast on first error", func(t *testing.T) {
		type Config struct {
			Field1 string `env:"FIELD1" validate:"required"`
			Field2 string `env:"FIELD2" validate:"required"`
			Field3 string `env:"FIELD3" validate:"required"`
		}

		var cfg Config
		err := Load(&cfg, WithFailFast(true))
		if err == nil {
			t.Fatal("expected error")
		}

		// Should be a single ValidationError, not ValidationErrors
		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("expected ValidationError, got %T", err)
		}

		// Should only report first error
		if validationErr.Field != "Field1" {
			t.Errorf("expected first field error, got %q", validationErr.Field)
		}
	})

	t.Run("collect all errors without fail fast", func(t *testing.T) {
		type Config struct {
			Field1 string `env:"FIELD1" validate:"required"`
			Field2 string `env:"FIELD2" validate:"required"`
			Field3 string `env:"FIELD3" validate:"required"`
		}

		var cfg Config
		err := Load(&cfg, WithFailFast(false))
		if err == nil {
			t.Fatal("expected error")
		}

		// Should be ValidationErrors (multiple)
		var validationErrs ValidationErrors
		if !errors.As(err, &validationErrs) {
			t.Fatalf("expected ValidationErrors, got %T", err)
		}

		// Should have 3 errors
		if len(validationErrs) != 3 {
			t.Errorf("expected 3 errors, got %d", len(validationErrs))
		}
	})
}

// Test combination of options
func TestCombinedOptions(t *testing.T) {
	type Config struct {
		Port int    `env:"PORT"`
		Host string `env:"HOST"`
	}

	setEnv(t, "MYAPP_PORT", "8080")
	setEnv(t, "MYAPP_HOST", "localhost")

	var cfg Config
	err := Load(&cfg,
		WithPrefix("MYAPP_"),
		WithRequiredByDefault(true),
		WithFailFast(false),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected 8080, got %d", cfg.Port)
	}
	if cfg.Host != "localhost" {
		t.Errorf("expected 'localhost', got %q", cfg.Host)
	}
}
