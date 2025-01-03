package parser

import (
	"github.com/paulrozhkin/jsonschema/pkg/parser"
	"github.com/paulrozhkin/jsonschema/tests/parser/base"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReflectFromType(t *testing.T) {
	expectedMetadata := ExpectedSettingsMetadata()
	obj := &base.Settings{}
	reflectParser := parser.NewReflectParser(obj)
	result, err := reflectParser.Parse()
	assert.Nil(t, err)
	assert.Equal(t, expectedMetadata, result)
}
