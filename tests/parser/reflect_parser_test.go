package parser

import (
	"github.com/paulrozhkin/jsonschema/pkg/parser"
	"github.com/paulrozhkin/jsonschema/tests/base"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReflectFromType(t *testing.T) {
	expectedMetadata := base.ExpectedSettingsMetadata()
	obj := &base.Settings{}
	reflectParser := parser.NewReflectParser(obj)
	result, err := reflectParser.Parse()
	require.Nil(t, err)
	require.EqualValues(t, expectedMetadata, result)
}
