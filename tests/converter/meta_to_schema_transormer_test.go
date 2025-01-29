package converter

import (
	"github.com/paulrozhkin/jsonschema/pkg/converter"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
	"github.com/paulrozhkin/jsonschema/tests/base"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvertMetadataToJSONSchema(t *testing.T) {
	expectedMetadata := base.ExpectedSettingsMetadata()
	expectedJsonSchema := base.ExpectedSettingsJsonSchema()
	cfg := entity.Config{SchemaVersion: entity.Draft202012}

	schemaConverter := converter.NewMetaToSchemaConverter()
	result, err := schemaConverter.Convert(cfg, expectedMetadata)
	require.Nil(t, err)
	require.Equal(t, expectedJsonSchema, result)
}
