package converter

import (
	"encoding/json"
	"fmt"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
	"strings"
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
	rootObject := schema.Defs[metadata.Root.TypeName].(*entity.ObjectSchema)
	delete(schema.Defs, metadata.Root.TypeName)
	schema.Properties = rootObject.Properties
	schema.Required = rootObject.Required
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
	definitions[dataTypeMetadata.TypeName] = objectSchema
	for _, node := range dataTypeMetadata.Nodes {
		dataType := typeKindToJsonSchemaType(node.TypeKind)
		if !node.IsPointer {
			objectSchema.Required = append(objectSchema.Required, getFieldName(node))
		}
		switch dataType {
		case entity.JSONSchemaNumber:
			objectSchema.AddProperty(getFieldName(node), entity.NewNumberSchema())
		case entity.JSONSchemaString:
			objectSchema.AddProperty(getFieldName(node), entity.NewStringSchema())
		case entity.JSONSchemaBoolean:
			objectSchema.AddProperty(getFieldName(node), entity.NewBooleanSchema())
		case entity.JSONSchemaInteger:
			integerSchema, err := transformIntegerToIntegerSchema(node)
			if err != nil {
				return nil, err
			}
			objectSchema.AddProperty(getFieldName(node), integerSchema)
		case entity.JSONSchemaUnknown:
			if node.Ref != nil {
				schema := entity.NewJSONEmptySchema().SetRef(fmt.Sprintf("#/$defs/%s", node.Ref.TypeName))
				objectSchema.AddProperty(getFieldName(node), schema)
			} else {
				return nil, fmt.Errorf("invalid object field %s for %s (%s)", dataTypeMetadata.TypeName,
					node.TypeName, node.TypeKind)
			}
		default:
			return nil, fmt.Errorf("not supported type for object field %s for %s (%s)", dataTypeMetadata.TypeName,
				node.TypeName, node.TypeKind)
		}
	}
	return objectSchema, nil
}

func transformIntegerToIntegerSchema(dataTypeMetadata *entity.DataTypeMetadata) (*entity.IntegerSchema, error) {
	integerSchema := entity.NewIntegerSchema()
	if jsonschemaTags, ok := dataTypeMetadata.Tags["jsonschema"]; ok {
		err := numericalKeywords(integerSchema, jsonschemaTags)
		if err != nil {
			return nil, fmt.Errorf("Field %s: ", dataTypeMetadata.FieldName)
		}
	}
	return integerSchema, nil
}

// read struct tags for numerical type keywords
func numericalKeywords(schema *entity.IntegerSchema, tags []string) error {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "multipleOf":
				schema.MultipleOf, _ = toJSONNumber[int](val)
			case "minimum":
				schema.Minimum, _ = toJSONNumber[int](val)
			case "maximum":
				schema.Maximum, _ = toJSONNumber[int](val)
			case "exclusiveMaximum":
				schema.ExclusiveMaximum, _ = toJSONNumber[int](val)
			case "exclusiveMinimum":
				schema.ExclusiveMinimum, _ = toJSONNumber[int](val)
			case "default":
				if num, ok := toJSONNumber[int](val); ok {
					schema.Default = num
				}
			case "example":
				if num, ok := toJSONNumber[int](val); ok {
					schema.Examples = append(schema.Examples, num)
				}
			case "enum":
				if num, ok := toJSONNumber[int](val); ok {
					schema.Enum = append(schema.Enum, num)
				}
			default:
				return fmt.Errorf("invalid tag for integer schema %s", tag)
			}
		}
	}
	return nil
}

// toJSONNumber converts string to *json.Number.
// It'll aso return whether the number is valid.
func toJSONNumber[T int | float64](s string) (*T, bool) {
	var result *T
	num := json.Number(s)
	if val, err := num.Int64(); err == nil {
		tResult := T(val)
		result = &tResult
		return result, true
	}
	if val, err := num.Float64(); err == nil {
		tResult := T(val)
		result = &tResult
		return result, true
	}
	return result, false
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
