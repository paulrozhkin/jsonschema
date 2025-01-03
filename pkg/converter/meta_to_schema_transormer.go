package converter

import "github.com/paulrozhkin/jsonschema/pkg/entity"

type Converter interface {
	Convert(metadata *entity.TypeMetadata) (*entity.JsonSchema, error)
}

type MetaToSchemaConverter struct{}

func NewMetaToSchemaConverter() *MetaToSchemaConverter {
	return &MetaToSchemaConverter{}
}

func (c *MetaToSchemaConverter) Convert(metadata *entity.TypeMetadata) (*entity.JsonSchema, error) {
	return &entity.JsonSchema{}, nil
}
