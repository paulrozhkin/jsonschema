package converter

import (
	"fmt"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
)

type Converter interface {
	Convert(config entity.Config, metadata *entity.TypeMetadata) (*entity.JSONSchema, error)
}

type MetaToSchemaConverter struct{}

func NewMetaToSchemaConverter() *MetaToSchemaConverter {
	return &MetaToSchemaConverter{}
}

func (c *MetaToSchemaConverter) Convert(config entity.Config, metadata *entity.TypeMetadata) (*entity.JSONSchema, error) {
	schema := entity.NewJSONSchema().
		SetSchema(config.SchemaVersion).
		SetID(c.getIdFromRootType(metadata))

	for _, node := range metadata.Nodes {
		dataType := typeKindToJsonSchemaType(node.TypeKind)
		if dataType == entity.JSONSchemaNumber {
			schema.AddProperty(getFieldName(node), entity.NewNumberSchema())
		}
	}
	return schema, nil
}

func getFieldName(metadata *entity.TypeMetadata) string {
	if jsonTags, ok := metadata.Tags["json"]; ok {
		return jsonTags[0]
	}
	return metadata.FieldName
}

func typeKindToJsonSchemaType(typeKind string) entity.JSONSchemaDataType {
	switch typeKind {
	case "string":
		return entity.JSONSchemaString
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return entity.JSONSchemaInteger
	case "float32", "float64":
		return entity.JSONSchemaNumber
	case "bool":
		return entity.JSONSchemaBoolean
	case "array", "slice":
		return entity.JSONSchemaArray
	case "struct":
		return entity.JSONSchemaObject
	case "map":
		return entity.JSONSchemaObject
	}
	return entity.JSONSchemaUnknown
}

func (c *MetaToSchemaConverter) getIdFromRootType(rootMetadata *entity.TypeMetadata) string {
	// Attempt to set the schema ID
	return fmt.Sprintf("https://%s/%s", rootMetadata.Package, rootMetadata.TypeName)
	//if !r.Anonymous && s.ID == EmptyID {
	//	baseSchemaID := r.BaseSchemaID
	//	if baseSchemaID == EmptyID {
	//		id := ID("https://" + t.PkgPath())
	//		if err := id.Validate(); err == nil {
	//			// it's okay to silently ignore URL errors
	//			baseSchemaID = id
	//		}
	//	}
	//	if baseSchemaID != EmptyID {
	//		s.ID = baseSchemaID.Add(ToSnakeCase(name))
	//	}
	//}
}
