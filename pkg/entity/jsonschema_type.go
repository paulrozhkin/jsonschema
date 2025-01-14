package entity

import (
	"encoding/json"
	"errors"
)

type JSONSchemaDataType string

const (
	JSONSchemaUnknown JSONSchemaDataType = "Unknown"
	JSONSchemaString  JSONSchemaDataType = "string"
	JSONSchemaNumber  JSONSchemaDataType = "number"
	JSONSchemaInteger JSONSchemaDataType = "integer"
	JSONSchemaObject  JSONSchemaDataType = "object"
	JSONSchemaArray   JSONSchemaDataType = "array"
	JSONSchemaBoolean JSONSchemaDataType = "boolean"
	JSONSchemaNull    JSONSchemaDataType = "null"
)

// JSONSchemaType is an alias for []string to represent the "type" field
type JSONSchemaType []JSONSchemaDataType

// MarshalJSON handles custom marshaling for JSONSchemaType
func (t JSONSchemaType) MarshalJSON() ([]byte, error) {
	if len(t) == 0 {
		return []byte(""), nil
	} else if len(t) == 1 {
		// If there's only one type, marshal as a single string
		return json.Marshal(t[0])
	}
	// Otherwise, marshal as an array of strings
	return json.Marshal([]JSONSchemaDataType(t))
}

// UnmarshalJSON handles custom unmarshaling for JSONSchemaType
func (t *JSONSchemaType) UnmarshalJSON(data []byte) error {
	var single JSONSchemaDataType
	var array []JSONSchemaDataType

	// Try unmarshaling as a single string
	if err := json.Unmarshal(data, &single); err == nil {
		*t = []JSONSchemaDataType{single}
		return nil
	}

	// Try unmarshaling as an array of strings
	if err := json.Unmarshal(data, &array); err == nil {
		*t = array
		return nil
	}

	return errors.New("invalid type: expected string or array of strings")
}
