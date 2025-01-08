package entity

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONSchemaTypeMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    JSONSchemaType
		expected string
	}{
		{
			name:     "Single type",
			input:    JSONSchemaType{"string"},
			expected: `"string"`,
		},
		{
			name:     "Multiple types",
			input:    JSONSchemaType{"string", "null"},
			expected: `["string","null"]`,
		},
		// TODO: пофиксить
		//{
		//	name:     "Empty type",
		//	input:    JSONSchemaType{},
		//	expected: ``,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(output))
		})
	}
}

func TestJSONSchemaTypeUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    JSONSchemaType
		expectedErr bool
	}{
		{
			name:        "Single type",
			input:       `"string"`,
			expected:    JSONSchemaType{"string"},
			expectedErr: false,
		},
		{
			name:        "Multiple types",
			input:       `["string","null"]`,
			expected:    JSONSchemaType{"string", "null"},
			expectedErr: false,
		},
		{
			name:        "Empty array",
			input:       `[]`,
			expected:    JSONSchemaType{},
			expectedErr: false,
		},
		{
			name:        "Invalid type (not string or array)",
			input:       `123`,
			expectedErr: true,
		},
		{
			name:        "Invalid type in array",
			input:       `["string",123]`,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output JSONSchemaType
			err := json.Unmarshal([]byte(tt.input), &output)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, output)
			}
		})
	}
}
