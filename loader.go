package envvalidator

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/KareemAzab21/envvalidator/internal/parser"
)

// load loads environment variables into the config struct
// This is called by Load() before validation
func (v *Validator) load(config interface{}) error {
	// Validate that config is a pointer to struct
	val := reflect.ValueOf(config)
	if val.Kind() != reflect.Ptr {
		return ErrNotStructPointer
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	typ := val.Type()
	var errors ValidationErrors

	// Iterate through all struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields (can't set them)
		if !field.CanSet() {
			continue
		}

		// Parse field tags
		fieldInfo := parser.ParseField(fieldType)

		// Skip fields without env tag
		if fieldInfo.EnvName == "" {
			continue
		}

		// Load the field value
		if err := v.loadField(field, fieldInfo); err != nil {
			errors.Add(ValidationError{
				Field:  fieldInfo.Name,
				EnvVar: fieldInfo.EnvName,
				Err:    err,
			})

			// If fail-fast is enabled, return immediately
			if v.failFast {
				return &errors[0]
			}
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// loadField loads a single field from environment variable
func (v *Validator) loadField(field reflect.Value, info parser.FieldInfo) error {
	// Build the full environment variable name (with prefix if set)
	envKey := info.EnvName
	if v.prefix != "" {
		envKey = v.prefix + envKey
	}

	// Get the environment variable value
	envValue := os.Getenv(envKey)

	// Handle missing environment variable
	if envValue == "" {
		return v.handleMissingValue(field, info)
	}

	// Convert and set the value
	return v.setField(field, envValue, info.EnvName)
}

// handleMissingValue handles the case when an environment variable is not set
func (v *Validator) handleMissingValue(field reflect.Value, info parser.FieldInfo) error {
	// If there's a default value, use it
	if info.HasDefault {
		return v.setField(field, info.DefaultValue, info.EnvName)
	}

	// Check if field is required
	isRequired := info.Required || v.requiredByDefault

	// If required and no default, return error
	if isRequired {
		return ErrRequired
	}

	// Optional field with no value - leave as zero value
	return nil
}

// setField converts a string value to the appropriate type and sets the field
func (v *Validator) setField(field reflect.Value, value string, envName string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.setIntField(field, value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.setUintField(field, value)

	case reflect.Bool:
		return v.setBoolField(field, value)

	case reflect.Float32, reflect.Float64:
		return v.setFloatField(field, value)

	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedType, field.Kind())
	}
}

// setIntField converts and sets an integer field
func (v *Validator) setIntField(field reflect.Value, value string) error {
	bitSize := field.Type().Bits()
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return fmt.Errorf("invalid integer value %q: %w", value, err)
	}

	field.SetInt(intVal)
	return nil
}

// setUintField converts and sets an unsigned integer field
func (v *Validator) setUintField(field reflect.Value, value string) error {
	bitSize := field.Type().Bits()
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return fmt.Errorf("invalid unsigned integer value %q: %w", value, err)
	}

	field.SetUint(uintVal)
	return nil
}

// setBoolField converts and sets a boolean field
// Accepts: "1", "t", "T", "true", "TRUE", "True", "0", "f", "F", "false", "FALSE", "False"
func (v *Validator) setBoolField(field reflect.Value, value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("invalid boolean value %q: %w", value, err)
	}

	field.SetBool(boolVal)
	return nil
}

// setFloatField converts and sets a float field
func (v *Validator) setFloatField(field reflect.Value, value string) error {
	bitSize := field.Type().Bits()
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return fmt.Errorf("invalid float value %q: %w", value, err)
	}

	field.SetFloat(floatVal)
	return nil
}
