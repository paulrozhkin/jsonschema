package converter

import (
	"fmt"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
)

type Converter interface {
	Convert(config entity.Config, metadata *entity.JsonSchemaMetadata) (*entity.JSONSchema, error)
}

type MetaToSchemaConverter struct{}

func NewMetaToSchemaConverter() *MetaToSchemaConverter {
	return &MetaToSchemaConverter{}
}

func (c *MetaToSchemaConverter) Convert(config entity.Config, metadata *entity.JsonSchemaMetadata) (*entity.JSONSchema, error) {
	schema := entity.NewJSONSchema().
		SetSchema(config.SchemaVersion).
		SetID(c.getIdFromRootType(metadata.Root))

	definitions, err := createDefinitions(metadata.Types)
	if err != nil {
		return nil, err
	}
	schema.Defs = definitions
	rootObject := schema.Defs[metadata.Root.ID()].(*entity.ObjectSchema)
	delete(schema.Defs, metadata.Root.ID())
	schema.Properties = rootObject.Properties

	//for _, node := range metadata.Root.Nodes {
	//	dataType := typeKindToJsonSchemaType(node.TypeKind)
	//	if dataType == entity.JSONSchemaNumber {
	//		schema.AddProperty(getFieldName(node), entity.NewNumberSchema())
	//	}
	//}
	return schema, nil
}

func createDefinitions(dataTypeDefinitions map[string]*entity.DataTypeMetadata) (map[string]entity.DataType, error) {
	if len(dataTypeDefinitions) == 0 {
		return nil, nil
	}
	definitions := make(map[string]entity.DataType)
	for _, dataTypeMetadata := range dataTypeDefinitions {
		dataType := typeKindToJsonSchemaType(dataTypeMetadata.TypeKind)
		if dataType != entity.JSONSchemaObject {
			return nil, fmt.Errorf("invalid data type for definisions: %s. Only struct supported", dataTypeMetadata.TypeKind)
		}
		_, err := transformObjectToObjectSchema(definitions, dataTypeMetadata)
		if err != nil {
			return nil, err
		}
	}
	return definitions, nil
}

func transformObjectToObjectSchema(definitions map[string]entity.DataType,
	dataTypeMetadata *entity.DataTypeMetadata) (*entity.ObjectSchema, error) {
	if dataTypeMetadata.Ref != nil {
		dataTypeMetadata = dataTypeMetadata.Ref
	}
	if objectSchema, ok := definitions[dataTypeMetadata.ID()]; ok {
		return objectSchema.(*entity.ObjectSchema), nil
	}

	objectSchema := entity.NewObjectSchema()
	definitions[dataTypeMetadata.ID()] = objectSchema

	for _, node := range dataTypeMetadata.Nodes {
		dataType := typeKindToJsonSchemaType(node.TypeKind)
		switch dataType {
		case entity.JSONSchemaNumber:
			objectSchema.AddProperty(getFieldName(node), entity.NewNumberSchema())
		case entity.JSONSchemaString:
			objectSchema.AddProperty(getFieldName(node), entity.NewStringSchema())
		case entity.JSONSchemaBoolean:
			objectSchema.AddProperty(getFieldName(node), entity.NewBooleanSchema())
		case entity.JSONSchemaInteger:
			objectSchema.AddProperty(getFieldName(node), entity.NewIntegerSchema())
		case entity.JSONSchemaUnknown:
			if node.Ref != nil {
				objectSchema.AddProperty(getFieldName(node), entity.NewJSONSchema().SetRef(node.Ref.ID()))
			} else {
				return nil, fmt.Errorf("invalid object field %s for %s (%s)", dataTypeMetadata.TypeName,
					node.TypeName, node.TypeKind)
			}
			//	objectSchemaNode, err := transformObjectToObjectSchema(definitions, node)
		//	if err != nil {
		//		return nil, err
		//	}
		//	objectSchema.AddProperty(getFieldName(node), objectSchemaNode)
		//case entity.JSONSchemaArray:
		default:
			return nil, fmt.Errorf("invalid object field %s for %s (%s)", dataTypeMetadata.TypeName,
				node.TypeName, node.TypeKind)
		}
	}
	return objectSchema, nil
}

func getFieldName(metadata *entity.DataTypeMetadata) string {
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

func (c *MetaToSchemaConverter) getIdFromRootType(rootMetadata *entity.DataTypeMetadata) string {
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
