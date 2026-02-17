// examples/custom-validator/main.go
package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/KareemAzab21/envvalidator"
	"github.com/KareemAzab21/envvalidator/validators"
)

// Config with custom validation rules
type Config struct {
	// Username must be alphanumeric
	Username string `env:"USERNAME" validate:"required,alphanumeric,min:3,max:20"`

	// Country code must be uppercase 2 letters
	CountryCode string `env:"COUNTRY_CODE" default:"US" validate:"uppercase,len:2"`

	// Workers must be an even number
	Workers int `env:"WORKERS" default:"4" validate:"even,range:2-100"`

	// Email with custom domain validation
	Email string `env:"EMAIL" validate:"required,email,domain:example.com"`

	// Password with strength validation
	Password string `env:"PASSWORD" validate:"required,strong_password"`

	// Hex color code
	Color string `env:"COLOR" default:"#FF5733" validate:"hexcolor"`
}

func main() {
	fmt.Println("=== EnvValidator Custom Validator Example ===\n")

	// Register custom validators
	customValidators := map[string]validators.ValidatorFunc{
		// Alphanumeric: only letters and numbers
		"alphanumeric": func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, str)
			if !matched {
				return fmt.Errorf("must contain only letters and numbers")
			}
			return nil
		},

		// Uppercase: all characters must be uppercase
		"uppercase": func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			if str != strings.ToUpper(str) {
				return fmt.Errorf("must be uppercase")
			}
			return nil
		},

		// Exact length
		"len": func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			var expectedLen int
			if _, err := fmt.Sscanf(param, "%d", &expectedLen); err != nil {
				return fmt.Errorf("invalid len parameter: %s", param)
			}
			if len(str) != expectedLen {
				return fmt.Errorf("length must be exactly %d characters (got %d)", expectedLen, len(str))
			}
			return nil
		},

		// Even number
		"even": func(value interface{}, param string) error {
			var num int64
			switch v := value.(type) {
			case int:
				num = int64(v)
			case int8:
				num = int64(v)
			case int16:
				num = int64(v)
			case int32:
				num = int64(v)
			case int64:
				num = v
			default:
				return nil
			}

			if num%2 != 0 {
				return fmt.Errorf("must be an even number (got %d)", num)
			}
			return nil
		},

		// Email domain validation
		"domain": func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			parts := strings.Split(str, "@")
			if len(parts) != 2 {
				return nil // Let email validator handle format
			}
			if parts[1] != param {
				return fmt.Errorf("email must be from domain %q (got %q)", param, parts[1])
			}
			return nil
		},

		// Strong password: min 8 chars, uppercase, lowercase, number, special char
		"strong_password": func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}

			if len(str) < 8 {
				return fmt.Errorf("password must be at least 8 characters")
			}

			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(str)
			hasLower := regexp.MustCompile(`[a-z]`).MatchString(str)
			hasNumber := regexp.MustCompile(`[0-9]`).MatchString(str)
			hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(str)

			if !hasUpper {
				return fmt.Errorf("password must contain at least one uppercase letter")
			}
			if !hasLower {
				return fmt.Errorf("password must contain at least one lowercase letter")
			}
			if !hasNumber {
				return fmt.Errorf("password must contain at least one number")
			}
			if !hasSpecial {
				return fmt.Errorf("password must contain at least one special character")
			}

			return nil
		},

		// Hex color code
		"hexcolor": func(value interface{}, param string) error {
			str, ok := value.(string)
			if !ok {
				return nil
			}
			matched, _ := regexp.MatchString(`^#[0-9A-Fa-f]{6}$`, str)
			if !matched {
				return fmt.Errorf("must be a valid hex color code (e.g., #FF5733)")
			}
			return nil
		},
	}

	var cfg Config

	// Load with custom validators
	err := envvalidator.Load(&cfg,
		envvalidator.WithCustomValidators(customValidators),
	)

	if err != nil {
		fmt.Println("❌ Validation errors found:\n")

		var validationErrs envvalidator.ValidationErrors
		if errors.As(err, &validationErrs) {
			for i, e := range validationErrs {
				fmt.Printf("  %d. %s (%s): %v\n", i+1, e.Field, e.EnvVar, e.Err)
			}
		} else {
			fmt.Printf("  Error: %v\n", err)
		}

		fmt.Println("\n" + strings.Repeat("=", 60))
		printUsageInstructions()
		return
	}

	// Print the loaded configuration
	fmt.Println("✅ Configuration loaded and validated successfully!\n")

	fmt.Println("Configuration:")
	fmt.Printf("  Username:     %s\n", cfg.Username)
	fmt.Printf("  Country Code: %s\n", cfg.CountryCode)
	fmt.Printf("  Workers:      %d\n", cfg.Workers)
	fmt.Printf("  Email:        %s\n", cfg.Email)
	fmt.Printf("  Password:     %s (hidden)\n", strings.Repeat("*", len(cfg.Password)))
	fmt.Printf("  Color:        %s\n", cfg.Color)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("\n🎉 All custom validators passed!")
	fmt.Println("\nCustom validators used:")
	fmt.Println("  ✓ alphanumeric     - Username contains only letters and numbers")
	fmt.Println("  ✓ uppercase        - Country code is uppercase")
	fmt.Println("  ✓ len:2            - Country code is exactly 2 characters")
	fmt.Println("  ✓ even             - Workers is an even number")
	fmt.Println("  ✓ domain           - Email is from example.com domain")
	fmt.Println("  ✓ strong_password  - Password meets complexity requirements")
	fmt.Println("  ✓ hexcolor         - Color is a valid hex color code")

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("\n💡 Try modifying the environment variables to see validation in action!")
}

// printUsageInstructions prints instructions for running the example
func printUsageInstructions() {
	fmt.Println("\n💡 To run this example successfully, set these environment variables:\n")

	fmt.Println("Required variables:")
	fmt.Println("  export USERNAME=\"john123\"")
	fmt.Println("    ↳ Must be alphanumeric, 3-20 characters")
	fmt.Println()
	fmt.Println("  export EMAIL=\"user@example.com\"")
	fmt.Println("    ↳ Must be valid email from example.com domain")
	fmt.Println()
	fmt.Println("  export PASSWORD=\"SecurePass123!\"")
	fmt.Println("    ↳ Must be 8+ chars with uppercase, lowercase, number, and special char")
	fmt.Println()

	fmt.Println("Optional variables (have defaults):")
	fmt.Println("  export COUNTRY_CODE=\"US\"")
	fmt.Println("    ↳ Must be uppercase, exactly 2 characters (default: US)")
	fmt.Println()
	fmt.Println("  export WORKERS=\"8\"")
	fmt.Println("    ↳ Must be even number between 2-100 (default: 4)")
	fmt.Println()
	fmt.Println("  export COLOR=\"#FF5733\"")
	fmt.Println("    ↳ Must be valid hex color code (default: #FF5733)")
	fmt.Println()

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nExample commands to set all variables:")
	fmt.Println()
	fmt.Println("  export USERNAME=\"john123\"")
	fmt.Println("  export COUNTRY_CODE=\"US\"")
	fmt.Println("  export WORKERS=\"8\"")
	fmt.Println("  export EMAIL=\"user@example.com\"")
	fmt.Println("  export PASSWORD=\"SecurePass123!\"")
	fmt.Println("  export COLOR=\"#FF5733\"")
	fmt.Println()
	fmt.Println("Then run:")
	fmt.Println("  go run examples/custom-validator/main.go")
	fmt.Println()

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n🧪 Test different scenarios:")
	fmt.Println()
	fmt.Println("1. Invalid username (too short):")
	fmt.Println("   export USERNAME=\"ab\"")
	fmt.Println()
	fmt.Println("2. Invalid country code (lowercase):")
	fmt.Println("   export COUNTRY_CODE=\"us\"")
	fmt.Println()
	fmt.Println("3. Invalid workers (odd number):")
	fmt.Println("   export WORKERS=\"7\"")
	fmt.Println()
	fmt.Println("4. Invalid email domain:")
	fmt.Println("   export EMAIL=\"user@gmail.com\"")
	fmt.Println()
	fmt.Println("5. Weak password:")
	fmt.Println("   export PASSWORD=\"weak\"")
	fmt.Println()
	fmt.Println("6. Invalid color code:")
	fmt.Println("   export COLOR=\"red\"")
	fmt.Println()
}
