package envvalidator

import (
	"reflect"

	"github.com/KareemAzab21/envvalidator/internal/parser"
	"github.com/KareemAzab21/envvalidator/validators"
)

// validate validates all fields in the config struct
func (v *Validator) validate(config interface{}) error {
	val := reflect.ValueOf(config).Elem()
	typ := val.Type()

	var errors ValidationErrors

	// Iterate through all struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Parse field tags
		fieldInfo := parser.ParseField(fieldType)

		// Skip fields without env tag
		if fieldInfo.EnvName == "" {
			continue
		}

		// Validate the field
		if err := v.validateField(field, fieldInfo); err != nil {
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

// validateField validates a single field using its validation rules
func (v *Validator) validateField(field reflect.Value, info parser.FieldInfo) error {
	// Skip if no validations
	if len(info.Validations) == 0 {
		return nil
	}

	// Convert parser.Rule to validators.ValidationRule
	var rules []validators.ValidationRule
	for _, rule := range info.Validations {
		rules = append(rules, validators.ValidationRule{
			Name:  rule.Name,
			Param: rule.Param,
		})
	}

	// Apply validators
	return validators.ApplyValidators(field, rules)
}
