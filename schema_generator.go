package jsonschema

import (
	"errors"
	"github.com/paulrozhkin/jsonschema/pkg/converter"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
	"github.com/paulrozhkin/jsonschema/pkg/parser"
)

var ErrParserNotFound = errors.New("parser not found")
var ErrConverterNotFound = errors.New("metadata to jsonschema converter not found")

type AfterParseFunc func(*entity.TypeMetadata) error
type AfterConvertFunc func(schema *entity.JsonSchema) error

type SchemaGenerator struct {
	Parser       parser.Parser
	Converter    converter.Converter
	AfterParse   AfterParseFunc
	AfterConvert AfterConvertFunc
	jsonSchema   *entity.JsonSchema
	Config       entity.Config
}

func DefaultGenerator() *SchemaGenerator {
	return &SchemaGenerator{
		Converter: converter.NewMetaToSchemaConverter(),
		Config: entity.Config{
			SchemaVersion: entity.Draft202012,
			FieldNameTag:  nil,
			KeyNamer:      nil,
		},
	}
}

func FromTypeToJsonSchema(obj any) (*SchemaGenerator, error) {
	generator := DefaultGenerator()
	generator.Parser = parser.NewReflectParser(obj)
	return generator, generator.Generate()
}

func FromFilesToJsonSchema() (*SchemaGenerator, error) {
	panic("implement me")
}

func (g *SchemaGenerator) Generate() error {
	// Parse go structs from any source to type metadata
	if g.Parser == nil {
		return ErrParserNotFound
	}
	metadata, err := g.Parser.Parse()
	if err != nil {
		return err
	}
	if g.AfterParse != nil {
		err = g.AfterParse(metadata)
		if err != nil {
			return err
		}
	}

	// Convert go type metadata to jsonschema struct
	if g.Converter == nil {
		return ErrConverterNotFound
	}
	jsonSchema, err := g.Converter.Convert(metadata)
	if err != nil {
		return err
	}
	if g.AfterConvert != nil {
		err = g.AfterConvert(jsonSchema)
		if err != nil {
			return err
		}
	}

	g.jsonSchema = jsonSchema
	return nil
}

func (g *SchemaGenerator) ToJson() string {
	panic("implement me")
}
