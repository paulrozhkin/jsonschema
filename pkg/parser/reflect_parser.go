package parser

import (
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

func (p *ReflectParser) Parse() (*entity.TypeMetadata, error) {
	// Get the type of the object
	objType := reflect.TypeOf(p.parsingObj)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	// Initialize metadata parsing
	metadata := parseTypeMetadata(objType)
	return metadata, nil
}

func parseTypeMetadata(t reflect.Type) *entity.TypeMetadata {
	metadata := &entity.TypeMetadata{
		Package:     t.PkgPath(),
		TypeName:    t.Name(),
		TypeKind:    t.Kind().String(),
		Nodes:       nil,
		Tags:        nil,
		IsReference: t.Kind() == reflect.Ptr,
	}

	// Only process structs
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			// Extract tags
			tags := extractTags(field.Tag)

			// Recursively parse nested structs
			nodeMetadata := parseTypeMetadata(fieldType)

			// Populate metadata for each field
			metadata.Nodes = append(metadata.Nodes, entity.TypeMetadata{
				FieldName:   field.Name,
				Package:     fieldType.PkgPath(),
				TypeName:    fieldType.Name(),
				TypeKind:    fieldType.Kind().String(),
				Nodes:       nodeMetadata.Nodes,
				Tags:        tags,
				IsReference: field.Type.Kind() == reflect.Ptr,
			})
		}
	}
	return metadata
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
