# EnvValidator

[![Go Reference](https://pkg.go.dev/badge/github.com/KareemAzab21/envvalidator.svg)](https://pkg.go.dev/github.com/KareemAzab21/envvalidator)
[![Go Report Card](https://goreportcard.com/badge/github.com/KareemAzab21/envvalidator)](https://goreportcard.com/report/github.com/KareemAzab21/envvalidator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen.svg)](https://github.com/KareemAzab21/envvalidator)

A powerful, type-safe Go library for loading and validating environment variables into structs with comprehensive validation rules.

## ✨ Features

- 🎯 **Type-safe** - Automatic type conversion with overflow protection
- ✅ **Rich Validation** - Built-in validators for common use cases
- 🔧 **Extensible** - Easy custom validator registration
- 📝 **Clear Errors** - Detailed error messages with field names
- 🚀 **Zero Dependencies** - Pure Go standard library
- 🎨 **Flexible** - Support for prefixes, defaults, and optional fields
- 🔒 **Production Ready** - Thoroughly tested with 95%+ coverage

---

## 📦 Installation

```bash
go get github.com/KareemAzab21/envvalidator
```

## 🚀 Quick Start

### 1️⃣ Define Your Configuration Struct

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/KareemAzab21/envvalidator"
)

type Config struct {
    // Server configuration
    Port        int    `env:"PORT" default:"8080" validate:"range:1000-9999"`
    Host        string `env:"HOST" default:"localhost"`
    
    // Application settings
    Environment string `env:"ENV" default:"dev" validate:"oneof:dev staging prod"`
    Debug       bool   `env:"DEBUG" default:"false"`
    
    // Security
    APIKey      string `env:"API_KEY" validate:"required,min:32"`
    
    // Database
    DatabaseURL string `env:"DATABASE_URL" validate:"required,url"`
    MaxConns    int    `env:"MAX_CONNS" default:"10" validate:"range:1-100"`
}

func main() {
    var cfg Config
    
    // Load and validate environment variables
    if err := envvalidator.Load(&cfg); err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Server starting on %s:%d\n", cfg.Host, cfg.Port)
    fmt.Printf("Environment: %s\n", cfg.Environment)
    fmt.Printf("Database: %s\n", cfg.DatabaseURL)
}
```

### 2️⃣ Set Environment Variables

```bash
export API_KEY="your-secret-api-key-here-minimum-32-characters"
export DATABASE_URL="https://db.example.com"
export ENV="prod"
export PORT="8080"
```

### 3️⃣ Run Your Application

```bash
go run main.go
```

**Output:**

```
Server starting on localhost:8080
Environment: prod
Database: https://db.example.com
```

## 📖 Table of Contents

- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Supported Types](#-supported-types)
- [Struct Tags](#️-struct-tags)
- [Built-in Validators](#-built-in-validators)
- [Options](#️-options)
- [Custom Validators](#custom-validators)
- [Error Handling](#error-handling)

## 🔢 Supported Types

EnvValidator supports the following Go types with automatic conversion:

| Type | Example Value | Environment Variable | Notes |
|------|-------|-----------|-------|
| string | "hello" | VALUE=hello | Direct string value |
| int | 42 | VALUE=42 | Standard integer |
| int8 | 127 | VALUE=127 | Range: -128 to 127 |
| int16 | 32767 | VALUE=32767 | Range: -32768 to 32767 |
| int32 | 2147483647 | VALUE=2147483647 | 32-bit integer |
| int64 | 9223372036854775807 | VALUE=9223372036854775807 | 64-bit integer |
| uint | 42 | VALUE=42 | Unsigned integer |
| uint8 | 255 | VALUE=255 | Range: 0 to 255 |
| uint16 | 65535 | VALUE=65535 | Range: 0 to 65535 |
| uint32 | 4294967295 | VALUE=4294967295 | 32-bit unsigned |
| uint64 | 18446744073709551615 | VALUE=18446744073709551615 | 64-bit unsigned |
| bool | true | VALUE=true | Accepts: true, false, 1, 0, t, f |
| float32 | 3.14 | VALUE=3.14 | 32-bit floating point |
| float64 | 3.141592653589793 | VALUE=3.141592653589793 | 64-bit floating point |
| time.Duration | 5s | VALUE=5s | Go duration format |
| []string | ["a","b"] | VALUE=a,b | Comma-separated |

## 🏷️ Struct Tags

### env Tag

Specifies the environment variable name.

```go
type Config struct {
    Port int `env:"PORT"`
    Host string `env:"SERVER_HOST"`
}
```

**Important:** Without env tag → field is ignored.

### default Tag

Provides a default value if environment variable is not set.

```go
type Config struct {
    Port    int    `env:"PORT" default:"8080"`
    Host    string `env:"HOST" default:"localhost"`
    Debug   bool   `env:"DEBUG" default:"false"`
}
```

**Behavior:**

- If environment variable is set → uses environment value
- If not set → uses default value
- If no default → uses zero value

### validate Tag

Specifies validation rules (comma-separated).

```go
type Config struct {
    APIKey   string `env:"API_KEY" validate:"required"`
    Username string `env:"USERNAME" validate:"required,min:3,max:20"`
    Port     int    `env:"PORT" default:"8080" validate:"range:1000-9999"`
}
```

## ✅ Built-in Validators

### String Validators

| Validator | Description | Example |
|-----------|-------------|---------|
| required | Must have value | validate:"required" |
| min:N | Minimum length | validate:"min:5" |
| max:N | Maximum length | validate:"max:100" |
| oneof:a b c | Must be one of options | validate:"oneof:dev staging prod" |
| email | Valid email format | validate:"email" |
| url | Must start with http:// or https:// | validate:"url" |

### Numeric Validators

| Validator | Description | Example |
|-----------|-------------|---------|
| range:MIN-MAX | Inclusive range | validate:"range:1000-9999" |
| minvalue:N | Minimum value | validate:"minvalue:0" |
| maxvalue:N | Maximum value | validate:"maxvalue:100" |

## ⚙️ Options

### WithPrefix(prefix string)

```go
envvalidator.Load(&cfg, envvalidator.WithPrefix("MYAPP_"))
```

### WithRequiredByDefault(bool)

```go
envvalidator.Load(&cfg, envvalidator.WithRequiredByDefault(true))
```

### WithFailFast(bool)

```go
envvalidator.Load(&cfg, envvalidator.WithFailFast(true))
```

### WithCustomValidator(name, func)

```go
uppercaseValidator := func(value interface{}, param string) error {
    str, ok := value.(string)
    if !ok {
        return nil
    }
    if str != strings.ToUpper(str) {
        return fmt.Errorf("must be uppercase")
    }
    return nil
}

envvalidator.Load(
    &cfg,
    envvalidator.WithCustomValidator("uppercase", uppercaseValidator),
)
```

## Custom Validators

## Error Handling

## Complete Examples

## Contributing

## License

## Support
