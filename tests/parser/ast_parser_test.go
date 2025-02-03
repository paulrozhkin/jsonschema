package parser

import (
	"github.com/paulrozhkin/jsonschema/pkg/parser"
	"github.com/paulrozhkin/jsonschema/tests/base"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAstFromPackageAndTypeNames(t *testing.T) {
	expectedMetadata := base.ExpectedSettingsMetadata()
	astParser := parser.NewAstParser("Settings", "github.com/paulrozhkin/jsonschema/tests/base")
	result, err := astParser.Parse()
	assert.Nil(t, err)
	assert.Equal(t, expectedMetadata, result)
}
