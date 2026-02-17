package validators

import (
	"reflect"
	"testing"
)

func TestValidateMin(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		param   string
		wantErr bool
	}{
		{"valid length", "hello", "3", false},
		{"exact length", "hello", "5", false},
		{"too short", "hi", "5", true},
		{"empty string", "", "1", true},
		{"non-string", 123, "5", false}, // Should skip
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMin(tt.value, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMax(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		param   string
		wantErr bool
	}{
		{"valid length", "hello", "10", false},
		{"exact length", "hello", "5", false},
		{"too long", "hello world", "5", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMax(tt.value, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMax() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOneOf(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		param   string
		wantErr bool
	}{
		{"valid value", "dev", "dev staging prod", false},
		{"another valid", "prod", "dev staging prod", false},
		{"invalid value", "test", "dev staging prod", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOneOf(tt.value, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOneOf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		param   string
		wantErr bool
	}{
		{"valid email", "user@example.com", "", false},
		{"valid with subdomain", "user@mail.example.com", "", false},
		{"invalid - no @", "userexample.com", "", true},
		{"invalid - no domain", "user@", "", true},
		{"empty string", "", "", false}, // Empty is OK (use required for that)
		{"invalid - no TLD", "user@example", "", true},
		{"valid with plus", "user+tag@example.com", "", false},
		{"non-string", 123, "", false}, // Should skip
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.value, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		param   string
		wantErr bool
	}{
		{"valid http", "http://example.com", "", false},
		{"valid https", "https://example.com", "", false},
		{"valid with path", "https://example.com/path/to/page", "", false},
		{"valid with query", "https://example.com?key=value", "", false},
		{"invalid - no protocol", "example.com", "", true},
		{"invalid - ftp protocol", "ftp://example.com", "", true},
		{"empty string", "", "", false}, // Empty is OK
		{"non-string", 123, "", false},  // Should skip
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.value, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRange(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		param   string
		wantErr bool
	}{
		{"valid int in range", 5000, "1000-9999", false},
		{"valid at min", 1000, "1000-9999", false},
		{"valid at max", 9999, "1000-9999", false},
		{"below range", 500, "1000-9999", true},
		{"above range", 10000, "1000-9999", true},
		{"valid int8", int8(50), "0-100", false},
		{"valid int64", int64(5000), "1000-9999", false},
		{"valid uint", uint(5000), "1000-9999", false},
		{"non-numeric", "hello", "1000-9999", false}, // Should skip
		{"invalid param format", 5000, "1000", true},
		{"invalid min", 5000, "abc-9999", true},
		{"invalid max", 5000, "1000-xyz", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRange(tt.value, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMinValue(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		param   string
		wantErr bool
	}{
		{"valid above min", 10, "5", false},
		{"valid at min", 5, "5", false},
		{"below min", 3, "5", true},
		{"valid negative", -5, "-10", false},
		{"below negative min", -15, "-10", true},
		{"non-numeric", "hello", "5", false}, // Should skip
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMinValue(tt.value, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMinValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMaxValue(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		param   string
		wantErr bool
	}{
		{"valid below max", 5, "10", false},
		{"valid at max", 10, "10", false},
		{"above max", 15, "10", true},
		{"valid negative", -15, "-10", false},
		{"above negative max", -5, "-10", true},
		{"non-numeric", "hello", "10", false}, // Should skip
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMaxValue(tt.value, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMaxValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    int64
		wantErr bool
	}{
		{"int", 42, 42, false},
		{"int8", int8(42), 42, false},
		{"int16", int16(42), 42, false},
		{"int32", int32(42), 42, false},
		{"int64", int64(42), 42, false},
		{"uint", uint(42), 42, false},
		{"uint8", uint8(42), 42, false},
		{"uint16", uint16(42), 42, false},
		{"uint32", uint32(42), 42, false},
		{"uint64", uint64(42), 42, false},
		{"float32", float32(42.7), 42, false},
		{"float64", float64(42.7), 42, false},
		{"string", "42", 0, true},
		{"bool", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toInt64(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("toInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("toInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterAndGet(t *testing.T) {
	// Test registering a custom validator
	customValidator := func(value interface{}, param string) error {
		return nil
	}

	Register("custom", customValidator)

	// Test getting the validator
	fn, ok := Get("custom")
	if !ok {
		t.Error("expected to find custom validator")
	}
	if fn == nil {
		t.Error("expected non-nil validator function")
	}

	// Test getting non-existent validator
	_, ok = Get("nonexistent")
	if ok {
		t.Error("expected not to find nonexistent validator")
	}
}

func TestApplyValidators(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		rules   []ValidationRule
		wantErr bool
	}{
		{
			name:  "valid string with min",
			value: "hello",
			rules: []ValidationRule{
				{Name: "min", Param: "3"},
			},
			wantErr: false,
		},
		{
			name:  "invalid string with min",
			value: "hi",
			rules: []ValidationRule{
				{Name: "min", Param: "5"},
			},
			wantErr: true,
		},
		{
			name:  "multiple valid rules",
			value: "hello",
			rules: []ValidationRule{
				{Name: "min", Param: "3"},
				{Name: "max", Param: "10"},
			},
			wantErr: false,
		},
		{
			name:  "skip required rule",
			value: "",
			rules: []ValidationRule{
				{Name: "required", Param: ""},
			},
			wantErr: false, // required is skipped in ApplyValidators
		},
		{
			name:  "unknown validator",
			value: "hello",
			rules: []ValidationRule{
				{Name: "unknown", Param: ""},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a reflect.Value from the test value
			val := reflect.ValueOf(tt.value)

			err := ApplyValidators(val, tt.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyValidators() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
