package parser

import (
	"fmt"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
	"reflect"
	"strings"
)

type ReflectParser struct {
	parsingObj any
}

func NewReflectParser(parsingObj any) *ReflectParser {
	return &ReflectParser{parsingObj: parsingObj}
}

func (p *ReflectParser) Parse() (schemaMetadata *entity.JsonSchemaMetadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unknown error in reflect Parse: %v", r)
		}
	}()
	// Create jsonschema metadata
	schemaMetadata = entity.NewJsonSchemaMetadata()

	// Get the type of the object
	objType := reflect.TypeOf(p.parsingObj)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	// Initialize metadata parsing
	rootMetadata, err := parseTypeMetadata(schemaMetadata, objType)
	if err != nil {
		return nil, err
	}
	schemaMetadata.Root = rootMetadata
	return schemaMetadata, nil
}

func parseTypeMetadata(schemaMetadata *entity.JsonSchemaMetadata,
	t reflect.Type) (metadata *entity.DataTypeMetadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unknown error in parseTypeMetadata: %v", r)
		}
	}()
	typeKind := t.Kind()
	metadata = entity.NewDataTypeMetadata(t.PkgPath(), t.Name(), typeKind.String(), typeKind == reflect.Ptr)
	// Only process structs
	if t.Kind() == reflect.Struct {
		// If data type metadata created then return it
		if dataTypeMetadata, ok := schemaMetadata.Types[metadata.ID()]; ok {
			return dataTypeMetadata, nil
		}
		schemaMetadata.Types[metadata.ID()] = metadata

		// Else create a new metadata for type
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldType := field.Type
			isPointer := fieldType.Kind() == reflect.Ptr
			if isPointer {
				fieldType = fieldType.Elem()
			}

			// Extract tags
			tags := extractTags(field.Tag)

			// Recursively parse nested structs
			nodeTypeMetadata, parseErr := parseTypeMetadata(schemaMetadata, fieldType)
			if err != nil {
				return nil, parseErr
			}

			// Populate metadata for each field
			nodeTypeKind := fieldType.Kind()
			var nodeMetadata *entity.DataTypeMetadata
			if nodeTypeKind == reflect.Struct {
				nodeMetadata = entity.NewDataTypeRefMetadata(nodeTypeMetadata)
			} else {
				nodeMetadata = entity.NewDataTypeMetadata(fieldType.PkgPath(), fieldType.Name(),
					nodeTypeKind.String(), nodeTypeKind == reflect.Ptr)
				nodeMetadata.Nodes = nodeTypeMetadata.Nodes
			}
			nodeMetadata.Tags = tags
			nodeMetadata.FieldName = field.Name
			nodeMetadata.IsPointer = isPointer
			metadata.Nodes = append(metadata.Nodes, nodeMetadata)
		}
	}
	return metadata, nil
}

func extractTags(tag reflect.StructTag) map[string][]string {
	var tags map[string][]string
	for _, key := range strings.Split(string(tag), " ") {
		// Parse each tag in the format `key:"value"`
		parts := strings.SplitN(key, ":", 2)
		if len(parts) == 2 {
			tagKey := parts[0]
			tagValue := strings.Trim(parts[1], `"`)
			if tags == nil {
				tags = make(map[string][]string)
			}
			tags[tagKey] = strings.Split(tagValue, ",")
		}
	}
	return tags
}
