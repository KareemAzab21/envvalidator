// examples/basic/main.go
package main

import (
	"fmt"
	"log"

	"github.com/KareemAzab21/envvalidator"
)

// Config represents a basic application configuration
type Config struct {
	// Server settings
	Port int    `env:"PORT" default:"8080"`
	Host string `env:"HOST" default:"localhost"`

	// Application settings
	AppName string `env:"APP_NAME" default:"MyApp"`
	Debug   bool   `env:"DEBUG" default:"false"`
}

func main() {
	fmt.Println("=== EnvValidator Basic Example ===\n")

	// Create a config instance
	var cfg Config

	// Load environment variables into the config struct
	if err := envvalidator.Load(&cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Print the loaded configuration
	fmt.Println("Configuration loaded successfully!")
	fmt.Printf("  App Name: %s\n", cfg.AppName)
	fmt.Printf("  Server:   %s:%d\n", cfg.Host, cfg.Port)
	fmt.Printf("  Debug:    %v\n", cfg.Debug)

	fmt.Println("\n✅ Basic example completed!")
	fmt.Println("\nTry setting environment variables:")
	fmt.Println("  export PORT=3000")
	fmt.Println("  export HOST=0.0.0.0")
	fmt.Println("  export APP_NAME=\"My Awesome App\"")
	fmt.Println("  export DEBUG=true")
	fmt.Println("  go run examples/basic/main.go")
}
