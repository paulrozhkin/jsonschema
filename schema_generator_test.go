package jsonschema

import (
	"github.com/paulrozhkin/jsonschema/tests/base"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGenerateSchemaFromType(t *testing.T) {
	generator, err := FromTypeToJsonSchema(base.Settings{})
	require.NoError(t, err)
	compareSchemaOutput(t, generator, "./tests/output/settings.json")
}

func compareSchemaOutput(t *testing.T, generator *SchemaGenerator, filename string) {
	t.Helper()
	expectedJSON, err := os.ReadFile(filename)
	require.NoError(t, err)

	actualJSON, err := generator.ToJson()
	require.NoError(t, err)
	require.JSONEq(t, string(expectedJSON), string(actualJSON))
}
