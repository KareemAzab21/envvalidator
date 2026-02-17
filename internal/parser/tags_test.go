// internal/parser/tags_test.go
package parser

import (
	"reflect"
	"testing"
)

func TestParseRule(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Rule
	}{
		{
			name:  "simple rule without parameter",
			input: "required",
			expected: Rule{
				Name:  "required",
				Param: "",
			},
		},
		{
			name:  "rule with parameter",
			input: "min:5",
			expected: Rule{
				Name:  "min",
				Param: "5",
			},
		},
		{
			name:  "rule with complex parameter",
			input: "range:1000-9999",
			expected: Rule{
				Name:  "range",
				Param: "1000-9999",
			},
		},
		{
			name:  "rule with spaces",
			input: " max : 100 ",
			expected: Rule{
				Name:  "max",
				Param: "100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseRule(tt.input)
			if result.Name != tt.expected.Name {
				t.Errorf("expected name %q, got %q", tt.expected.Name, result.Name)
			}
			if result.Param != tt.expected.Param {
				t.Errorf("expected param %q, got %q", tt.expected.Param, result.Param)
			}
		})
	}
}

func TestParseValidationTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Rule
	}{
		{
			name:     "empty tag",
			input:    "",
			expected: []Rule{},
		},
		{
			name:  "single rule",
			input: "required",
			expected: []Rule{
				{Name: "required", Param: ""},
			},
		},
		{
			name:  "multiple rules",
			input: "required,min:5,max:10",
			expected: []Rule{
				{Name: "required", Param: ""},
				{Name: "min", Param: "5"},
				{Name: "max", Param: "10"},
			},
		},
		{
			name:  "rules with spaces",
			input: " required , min:5 , max:10 ",
			expected: []Rule{
				{Name: "required", Param: ""},
				{Name: "min", Param: "5"},
				{Name: "max", Param: "10"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseValidationTag(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d rules, got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i].Name != tt.expected[i].Name {
					t.Errorf("rule %d: expected name %q, got %q", i, tt.expected[i].Name, result[i].Name)
				}
				if result[i].Param != tt.expected[i].Param {
					t.Errorf("rule %d: expected param %q, got %q", i, tt.expected[i].Param, result[i].Param)
				}
			}
		})
	}
}

func TestParseField(t *testing.T) {
	type TestStruct struct {
		Port     int    `env:"PORT" validate:"required,range:1000-9999" default:"8080"`
		Host     string `env:"HOST" default:"localhost"`
		Debug    bool   `env:"DEBUG"`
		Optional string `env:"OPTIONAL" validate:"min:5"`
	}

	typ := reflect.TypeOf(TestStruct{})

	t.Run("field with all tags", func(t *testing.T) {
		field := typ.Field(0) // Port
		info := ParseField(field)

		if info.Name != "Port" {
			t.Errorf("expected name Port, got %s", info.Name)
		}
		if info.EnvName != "PORT" {
			t.Errorf("expected env name PORT, got %s", info.EnvName)
		}
		if info.DefaultValue != "8080" {
			t.Errorf("expected default 8080, got %s", info.DefaultValue)
		}
		if !info.HasDefault {
			t.Error("expected HasDefault to be true")
		}
		if !info.Required {
			t.Error("expected Required to be true")
		}
		if len(info.Validations) != 2 {
			t.Errorf("expected 2 validations, got %d", len(info.Validations))
		}
	})

	t.Run("field with default only", func(t *testing.T) {
		field := typ.Field(1) // Host
		info := ParseField(field)

		if info.EnvName != "HOST" {
			t.Errorf("expected env name HOST, got %s", info.EnvName)
		}
		if info.DefaultValue != "localhost" {
			t.Errorf("expected default localhost, got %s", info.DefaultValue)
		}
		if info.Required {
			t.Error("expected Required to be false")
		}
	})

	t.Run("field without default", func(t *testing.T) {
		field := typ.Field(2) // Debug
		info := ParseField(field)

		if info.HasDefault {
			t.Error("expected HasDefault to be false")
		}
		if info.DefaultValue != "" {
			t.Errorf("expected empty default, got %s", info.DefaultValue)
		}
	})
}

func TestFieldInfo_HasRule(t *testing.T) {
	info := FieldInfo{
		Validations: []Rule{
			{Name: "required", Param: ""},
			{Name: "min", Param: "5"},
		},
	}

	if !info.HasRule("required") {
		t.Error("expected to have 'required' rule")
	}
	if !info.HasRule("min") {
		t.Error("expected to have 'min' rule")
	}
	if info.HasRule("max") {
		t.Error("expected not to have 'max' rule")
	}
}

func TestFieldInfo_GetRule(t *testing.T) {
	info := FieldInfo{
		Validations: []Rule{
			{Name: "min", Param: "5"},
			{Name: "max", Param: "10"},
		},
	}

	rule := info.GetRule("min")
	if rule == nil {
		t.Fatal("expected to find 'min' rule")
	}
	if rule.Param != "5" {
		t.Errorf("expected param '5', got '%s'", rule.Param)
	}

	rule = info.GetRule("nonexistent")
	if rule != nil {
		t.Error("expected nil for nonexistent rule")
	}
}
