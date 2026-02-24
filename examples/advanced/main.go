// examples/advanced/main.go
package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/KareemAzab21/envvalidator"
)

// ServerConfig represents advanced server configuration
type ServerConfig struct {
	// Server settings with validation
	Port         int           `env:"PORT" default:"8080" validate:"range:1000-65535"`
	Host         string        `env:"HOST" default:"0.0.0.0"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" default:"30s"`

	// Application settings
	Environment string `env:"ENV" default:"dev" validate:"oneof:dev staging prod"`
	Debug       bool   `env:"DEBUG" default:"false"`
	LogLevel    string `env:"LOG_LEVEL" default:"info" validate:"oneof:debug info warn error"`

	// Security
	APIKey    string `env:"API_KEY" validate:"required,min:32"`
	SecretKey string `env:"SECRET_KEY" validate:"required,min:16"`

	// Database
	DatabaseURL  string        `env:"DATABASE_URL" validate:"required,url"`
	MaxConns     int           `env:"MAX_CONNS" default:"10" validate:"range:1-100"`
	MinConns     int           `env:"MIN_CONNS" default:"2" validate:"range:1-50"`
	ConnLifetime time.Duration `env:"CONN_LIFETIME" default:"5m"`

	// External services
	RedisURL     string `env:"REDIS_URL" validate:"required,url"`
	CacheEnabled bool   `env:"CACHE_ENABLED" default:"true"`

	// Feature flags
	Features []string `env:"FEATURES" default:"feature1,feature2,feature3"`

	// Metrics
	MetricsEnabled bool `env:"METRICS_ENABLED" default:"true"`
	MetricsPort    int  `env:"METRICS_PORT" default:"9090" validate:"range:1000-9999"`
}

func main() {
	fmt.Println("=== EnvValidator Advanced Example ===")

	var cfg ServerConfig

	// Load with options
	err := envvalidator.Load(&cfg,
		envvalidator.WithPrefix("APP_"),           // Add APP_ prefix to all env vars
		envvalidator.WithRequiredByDefault(false), // Fields optional unless marked required
		envvalidator.WithFailFast(false),          // Collect all errors
	)

	if err != nil {
		// Handle validation errors
		var validationErrs envvalidator.ValidationErrors
		if errors.As(err, &validationErrs) {
			fmt.Println("❌ Validation errors found:")
			for i, e := range validationErrs {
				fmt.Printf("  %d. Field '%s' (%s): %v\n", i+1, e.Field, e.EnvVar, e.Err)
			}
			fmt.Println("\n💡 Fix the errors above and try again.")
			return
		}

		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Print the loaded configuration
	fmt.Println("✅ Configuration loaded successfully!")

	fmt.Println("Server Configuration:")
	fmt.Printf("  Host:          %s\n", cfg.Host)
	fmt.Printf("  Port:          %d\n", cfg.Port)
	fmt.Printf("  Read Timeout:  %v\n", cfg.ReadTimeout)
	fmt.Printf("  Write Timeout: %v\n", cfg.WriteTimeout)

	fmt.Println("\nApplication Settings:")
	fmt.Printf("  Environment:   %s\n", cfg.Environment)
	fmt.Printf("  Debug:         %v\n", cfg.Debug)
	fmt.Printf("  Log Level:     %s\n", cfg.LogLevel)

	fmt.Println("\nSecurity:")
	fmt.Printf("  API Key:       %s... (hidden)\n", maskString(cfg.APIKey, 8))
	fmt.Printf("  Secret Key:    %s... (hidden)\n", maskString(cfg.SecretKey, 8))

	fmt.Println("\nDatabase:")
	fmt.Printf("  URL:           %s\n", cfg.DatabaseURL)
	fmt.Printf("  Max Conns:     %d\n", cfg.MaxConns)
	fmt.Printf("  Min Conns:     %d\n", cfg.MinConns)
	fmt.Printf("  Conn Lifetime: %v\n", cfg.ConnLifetime)

	fmt.Println("\nExternal Services:")
	fmt.Printf("  Redis URL:     %s\n", cfg.RedisURL)
	fmt.Printf("  Cache Enabled: %v\n", cfg.CacheEnabled)

	fmt.Println("\nFeature Flags:")
	fmt.Printf("  Features:      %v\n", cfg.Features)

	fmt.Println("\nMetrics:")
	fmt.Printf("  Enabled:       %v\n", cfg.MetricsEnabled)
	fmt.Printf("  Port:          %d\n", cfg.MetricsPort)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("\n💡 To run this example, set the following environment variables:")
	fmt.Println("\nRequired:")
	fmt.Println("  export APP_API_KEY=\"your-secret-api-key-minimum-32-characters-long\"")
	fmt.Println("  export APP_SECRET_KEY=\"your-secret-key-16+\"")
	fmt.Println("  export APP_DATABASE_URL=\"https://db.example.com\"")
	fmt.Println("  export APP_REDIS_URL=\"https://redis.example.com\"")
	fmt.Println("\nOptional (with defaults):")
	fmt.Println("  export APP_PORT=8080")
	fmt.Println("  export APP_ENV=prod")
	fmt.Println("  export APP_LOG_LEVEL=debug")
	fmt.Println("  export APP_FEATURES=\"feature1,feature2,feature3\"")
	fmt.Println("\nThen run:")
	fmt.Println("  go run examples/advanced/main.go")
}

// maskString masks a string showing only the first n characters
func maskString(s string, show int) string {
	if len(s) <= show {
		return s
	}
	return s[:show] + strings.Repeat("*", len(s)-show)
}
