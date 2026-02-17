package envvalidator

import (
	"errors"
	"os"
	"testing"
)

func TestValidationIntegration(t *testing.T) {
	t.Run("string validation - min length", func(t *testing.T) {
		type Config struct {
			Name string `env:"NAME" validate:"min:5"`
		}

		os.Setenv("NAME", "John")
		defer os.Unsetenv("NAME")

		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected validation error for short name")
		}
	})

	t.Run("string validation - valid", func(t *testing.T) {
		type Config struct {
			Name string `env:"NAME" validate:"min:3,max:10"`
		}

		os.Setenv("NAME", "Alice")
		defer os.Unsetenv("NAME")

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Name != "Alice" {
			t.Errorf("expected Alice, got %s", cfg.Name)
		}
	})

	t.Run("numeric validation - range", func(t *testing.T) {
		type Config struct {
			Port int `env:"PORT" validate:"range:1000-9999"`
		}

		os.Setenv("PORT", "8080")
		defer os.Unsetenv("PORT")

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Port != 8080 {
			t.Errorf("expected 8080, got %d", cfg.Port)
		}
	})

	t.Run("numeric validation - out of range", func(t *testing.T) {
		type Config struct {
			Port int `env:"PORT" validate:"range:1000-9999"`
		}

		os.Setenv("PORT", "100")
		defer os.Unsetenv("PORT")

		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected validation error for port out of range")
		}
	})

	t.Run("email validation", func(t *testing.T) {
		type Config struct {
			Email string `env:"EMAIL" validate:"email"`
		}

		os.Setenv("EMAIL", "user@example.com")
		defer os.Unsetenv("EMAIL")

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("oneof validation", func(t *testing.T) {
		type Config struct {
			Environment string `env:"ENV" validate:"oneof:dev staging prod"`
		}

		os.Setenv("ENV", "dev")
		defer os.Unsetenv("ENV")

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Environment != "dev" {
			t.Errorf("expected dev, got %s", cfg.Environment)
		}
	})

	t.Run("oneof validation - invalid", func(t *testing.T) {
		type Config struct {
			Environment string `env:"ENV" validate:"oneof:dev staging prod"`
		}

		os.Setenv("ENV", "test")
		defer os.Unsetenv("ENV")

		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected validation error for invalid environment")
		}
	})

	t.Run("url validation", func(t *testing.T) {
		type Config struct {
			APIUrl string `env:"API_URL" validate:"url"`
		}

		os.Setenv("API_URL", "https://api.example.com")
		defer os.Unsetenv("API_URL")

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("url validation - invalid", func(t *testing.T) {
		type Config struct {
			APIUrl string `env:"API_URL" validate:"url"`
		}

		os.Setenv("API_URL", "not-a-url")
		defer os.Unsetenv("API_URL")

		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected validation error for invalid URL")
		}
	})

	t.Run("minvalue validation", func(t *testing.T) {
		type Config struct {
			Workers int `env:"WORKERS" validate:"minvalue:1"`
		}

		os.Setenv("WORKERS", "5")
		defer os.Unsetenv("WORKERS")

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Workers != 5 {
			t.Errorf("expected 5, got %d", cfg.Workers)
		}
	})

	t.Run("minvalue validation - invalid", func(t *testing.T) {
		type Config struct {
			Workers int `env:"WORKERS" validate:"minvalue:1"`
		}

		os.Setenv("WORKERS", "0")
		defer os.Unsetenv("WORKERS")

		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected validation error for workers below minimum")
		}
	})

	t.Run("maxvalue validation", func(t *testing.T) {
		type Config struct {
			Timeout int `env:"TIMEOUT" validate:"maxvalue:300"`
		}

		os.Setenv("TIMEOUT", "60")
		defer os.Unsetenv("TIMEOUT")

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Timeout != 60 {
			t.Errorf("expected 60, got %d", cfg.Timeout)
		}
	})

	t.Run("maxvalue validation - invalid", func(t *testing.T) {
		type Config struct {
			Timeout int `env:"TIMEOUT" validate:"maxvalue:300"`
		}

		os.Setenv("TIMEOUT", "500")
		defer os.Unsetenv("TIMEOUT")

		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected validation error for timeout above maximum")
		}
	})

	t.Run("multiple validations", func(t *testing.T) {
		type Config struct {
			Username string `env:"USERNAME" validate:"required,min:3,max:20"`
			Port     int    `env:"PORT" validate:"range:1000-9999"`
			Email    string `env:"EMAIL" validate:"email"`
		}

		os.Setenv("USERNAME", "alice")
		os.Setenv("PORT", "8080")
		os.Setenv("EMAIL", "alice@example.com")
		defer func() {
			os.Unsetenv("USERNAME")
			os.Unsetenv("PORT")
			os.Unsetenv("EMAIL")
		}()

		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Username != "alice" {
			t.Errorf("expected alice, got %s", cfg.Username)
		}
		if cfg.Port != 8080 {
			t.Errorf("expected 8080, got %d", cfg.Port)
		}
		if cfg.Email != "alice@example.com" {
			t.Errorf("expected alice@example.com, got %s", cfg.Email)
		}
	})

	t.Run("validation with default value", func(t *testing.T) {
		type Config struct {
			Environment string `env:"ENV" default:"dev" validate:"oneof:dev staging prod"`
		}

		// Don't set ENV - should use default
		var cfg Config
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Environment != "dev" {
			t.Errorf("expected dev, got %s", cfg.Environment)
		}
	})

	t.Run("validation with invalid default value", func(t *testing.T) {
		type Config struct {
			Environment string `env:"ENV" default:"invalid" validate:"oneof:dev staging prod"`
		}

		// Don't set ENV - should use default but fail validation
		var cfg Config
		err := Load(&cfg)
		if err == nil {
			t.Fatal("expected validation error for invalid default value")
		}
	})

	t.Run("fail fast on validation error", func(t *testing.T) {
		type Config struct {
			Field1 string `env:"FIELD1" validate:"min:10"`
			Field2 string `env:"FIELD2" validate:"min:10"`
		}

		os.Setenv("FIELD1", "short")
		os.Setenv("FIELD2", "short")
		defer func() {
			os.Unsetenv("FIELD1")
			os.Unsetenv("FIELD2")
		}()

		var cfg Config
		err := Load(&cfg, WithFailFast(true))
		if err == nil {
			t.Fatal("expected validation error")
		}

		// With fail-fast, should only get one error
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			// Single error - good
		} else {
			var validationErrs ValidationErrors
			if errors.As(err, &validationErrs) {
				if len(validationErrs) > 1 {
					t.Error("expected only one error with fail-fast enabled")
				}
			}
		}
	})

	t.Run("collect all validation errors", func(t *testing.T) {
		type Config struct {
			Field1 string `env:"FIELD1" validate:"min:10"`
			Field2 string `env:"FIELD2" validate:"min:10"`
		}

		os.Setenv("FIELD1", "short")
		os.Setenv("FIELD2", "short")
		defer func() {
			os.Unsetenv("FIELD1")
			os.Unsetenv("FIELD2")
		}()

		var cfg Config
		err := Load(&cfg) // No fail-fast
		if err == nil {
			t.Fatal("expected validation errors")
		}

		// Should collect both errors
		var validationErrs ValidationErrors
		if errors.As(err, &validationErrs) {
			if len(validationErrs) != 2 {
				t.Errorf("expected 2 errors, got %d", len(validationErrs))
			}
		}
	})

	t.Run("validation with prefix", func(t *testing.T) {
		type Config struct {
			Port int `env:"PORT" validate:"range:1000-9999"`
		}

		os.Setenv("APP_PORT", "8080")
		defer os.Unsetenv("APP_PORT")

		var cfg Config
		err := Load(&cfg, WithPrefix("APP_"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Port != 8080 {
			t.Errorf("expected 8080, got %d", cfg.Port)
		}
	})

	t.Run("skip validation for fields without env tag", func(t *testing.T) {
		type Config struct {
			Internal string // No env tag
			External string `env:"EXTERNAL" validate:"min:5"`
		}

		os.Setenv("EXTERNAL", "valid_value")
		defer os.Unsetenv("EXTERNAL")

		var cfg Config
		cfg.Internal = "x" // Set to short value
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Internal should remain unchanged
		if cfg.Internal != "x" {
			t.Errorf("expected x, got %s", cfg.Internal)
		}
	})
}

func TestComplexValidationScenarios(t *testing.T) {
	t.Run("real-world config example", func(t *testing.T) {
		type DatabaseConfig struct {
			Host     string `env:"DB_HOST" default:"localhost" validate:"min:1"`
			Port     int    `env:"DB_PORT" default:"5432" validate:"range:1024-65535"`
			Username string `env:"DB_USER" validate:"required,min:3"`
			Password string `env:"DB_PASS" validate:"required,min:8"`
			Database string `env:"DB_NAME" validate:"required,min:1"`
		}

		os.Setenv("DB_USER", "admin")
		os.Setenv("DB_PASS", "securepassword123")
		os.Setenv("DB_NAME", "myapp")
		defer func() {
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASS")
			os.Unsetenv("DB_NAME")
		}()

		var cfg DatabaseConfig
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Host != "localhost" {
			t.Errorf("expected localhost, got %s", cfg.Host)
		}
		if cfg.Port != 5432 {
			t.Errorf("expected 5432, got %d", cfg.Port)
		}
		if cfg.Username != "admin" {
			t.Errorf("expected admin, got %s", cfg.Username)
		}
	})

	t.Run("api config with validations", func(t *testing.T) {
		type APIConfig struct {
			BaseURL     string `env:"API_URL" validate:"required,url"`
			Timeout     int    `env:"API_TIMEOUT" default:"30" validate:"range:1-300"`
			Environment string `env:"API_ENV" default:"dev" validate:"oneof:dev staging prod"`
			APIKey      string `env:"API_KEY" validate:"required,min:32"`
		}

		os.Setenv("API_URL", "https://api.example.com")
		os.Setenv("API_KEY", "abcdef1234567890abcdef1234567890")
		defer func() {
			os.Unsetenv("API_URL")
			os.Unsetenv("API_KEY")
		}()

		var cfg APIConfig
		err := Load(&cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.BaseURL != "https://api.example.com" {
			t.Errorf("expected https://api.example.com, got %s", cfg.BaseURL)
		}
		if cfg.Timeout != 30 {
			t.Errorf("expected 30, got %d", cfg.Timeout)
		}
		if cfg.Environment != "dev" {
			t.Errorf("expected dev, got %s", cfg.Environment)
		}
	})
}
