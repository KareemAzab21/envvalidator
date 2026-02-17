package envvalidator

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinel errors - predefined errors that can be checked with errors.Is()
var (
	// ErrNotStruct is returned when the config is not a struct
	ErrNotStruct = errors.New("config must be a struct")

	// ErrNotStructPointer is returned when config is not a pointer to struct
	ErrNotStructPointer = errors.New("config must be a pointer to struct")

	// ErrRequired is returned when a required field is missing
	ErrRequired = errors.New("field is required but not set")

	// ErrUnsupportedType is returned when a field type is not supported
	ErrUnsupportedType = errors.New("unsupported field type")

	// ErrInvalidValue is returned when a value doesn't pass validation
	ErrInvalidValue = errors.New("invalid value")

	// ErrInvalidTag is returned when a struct tag is malformed
	ErrInvalidTag = errors.New("invalid struct tag format")
)

// ValidationError represents a single validation error for a specific field
// This provides context about WHICH field failed and WHY
type ValidationError struct {
	Field  string // Name of the struct field that failed
	EnvVar string // Name of the environment variable
	Value  string // The actual value that failed validation
	Tag    string // The validation tag that failed (e.g., "min:5")
	Err    error  // The underlying error

}

// Error implements the error interface
// Returns a human-readable error message
func (e *ValidationError) Error() string {
	if e.EnvVar != "" {
		return fmt.Sprintf("validation failed for field '%s'(env: %s) %v", e.Field, e.EnvVar, e.Err)
	}
	return fmt.Sprintf("validation failed for field '%s': %v", e.Field, e.Err)
}

// Unwrap allows errors.Is() and errors.As() to work with wrapped errors
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// ValidationErrors is a collection of multiple validation errors
// Useful when you want to collect ALL errors instead of failing on first error
type ValidationErrors []ValidationError

// Error implements the error interface for multiple errors
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}

	if len(e) == 1 {
		return e[0].Error()
	}

	// For multiple errors, show first error + count
	return fmt.Sprintf("%s (and %d more error(s))", e[0].Error(), len(e)-1)
}

// Errors returns all error messages as a slice of strings
// Useful for displaying all errors to users
func (e ValidationErrors) Errors() []string {
	errs := make([]string, len(e))
	for i, err := range e {
		errs[i] = err.Error()
	}
	return errs
}

// Error returns all errors as a formatted string
func (e ValidationErrors) String() string {
	return strings.Join(e.Errors(), "\n")
}

// Add appends a new validation error
func (e *ValidationErrors) Add(err ValidationError) {
	*e = append(*e, err)
}

// HasErrors returns true if there are any errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}
