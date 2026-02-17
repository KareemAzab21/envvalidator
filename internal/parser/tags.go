// internal/parser/tags.go
package parser

import (
	"reflect"
	"strings"
)

// FieldInfo contains all parsed information about a struct field
type FieldInfo struct {
	Name         string   // Struct field name (e.g., "Port")
	EnvName      string   // Environment variable name (e.g., "PORT")
	DefaultValue string   // Default value if env var not set
	Validations  []Rule   // List of validation rules
	Required     bool     // Whether field is required
	HasDefault   bool     // Whether a default value is specified
}

// Rule represents a single validation rule
type Rule struct {
	Name  string // Rule name (e.g., "min", "max", "range")
	Param string // Rule parameter (e.g., "5" for "min:5")
}

// ParseField extracts all tag information from a struct field
//
// Example struct field:
//   Port int `env:"PORT" validate:"required,range:1000-9999" default:"8080"`
//
// Returns FieldInfo with all parsed data
func ParseField(field reflect.StructField) FieldInfo {
	info := FieldInfo{
		Name:        field.Name,
		Validations: []Rule{},
	}

	// Parse env tag
	if envTag := field.Tag.Get("env"); envTag != "" {
		info.EnvName = envTag
	}

	// Parse default tag
	if defaultTag := field.Tag.Get("default"); defaultTag != "" {
		info.DefaultValue = defaultTag
		info.HasDefault = true
	}

	// Parse validate tag
	if validateTag := field.Tag.Get("validate"); validateTag != "" {
		info.Validations = ParseValidationTag(validateTag)
		
		// Check if "required" is in validations
		for _, rule := range info.Validations {
			if rule.Name == "required" {
				info.Required = true
				break
			}
		}
	}

	return info
}

// ParseValidationTag parses a validation tag string into individual rules
//
// Examples:
//   "required" → [{Name: "required", Param: ""}]
//   "min:5,max:10" → [{Name: "min", Param: "5"}, {Name: "max", Param: "10"}]
//   "range:1000-9999" → [{Name: "range", Param: "1000-9999"}]
func ParseValidationTag(tag string) []Rule {
	if tag == "" {
		return []Rule{}
	}

	var rules []Rule
	
	// Split by comma to get individual rules
	parts := strings.Split(tag, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		rule := ParseRule(part)
		rules = append(rules, rule)
	}

	return rules
}

// ParseRule parses a single validation rule
//
// Examples:
//   "required" → Rule{Name: "required", Param: ""}
//   "min:5" → Rule{Name: "min", Param: "5"}
//   "range:1000-9999" → Rule{Name: "range", Param: "1000-9999"}
func ParseRule(ruleStr string) Rule {
	// Check if rule has a parameter (contains ":")
	colonIdx := strings.Index(ruleStr, ":")
	
	if colonIdx == -1 {
		// No parameter (e.g., "required")
		return Rule{
			Name:  ruleStr,
			Param: "",
		}
	}

	// Has parameter (e.g., "min:5")
	return Rule{
		Name:  strings.TrimSpace(ruleStr[:colonIdx]),
		Param: strings.TrimSpace(ruleStr[colonIdx+1:]),
	}
}

// HasRule checks if a field has a specific validation rule
func (f *FieldInfo) HasRule(ruleName string) bool {
	for _, rule := range f.Validations {
		if rule.Name == ruleName {
			return true
		}
	}
	return false
}

// GetRule returns the rule with the given name, or nil if not found
func (f *FieldInfo) GetRule(ruleName string) *Rule {
	for i := range f.Validations {
		if f.Validations[i].Name == ruleName {
			return &f.Validations[i]
		}
	}
	return nil
}
