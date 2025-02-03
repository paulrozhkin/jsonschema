package entity

import "fmt"

type JsonSchemaMetadata struct {
	Types map[string]*DataTypeMetadata
	Root  *DataTypeMetadata
}

type DataTypeMetadata struct {
	Ref       *DataTypeMetadata
	Package   string
	TypeName  string
	TypeKind  string
	FieldName string
	Nodes     []*DataTypeMetadata
	Tags      map[string][]string
	IsPointer bool
}

func NewJsonSchemaMetadata() *JsonSchemaMetadata {
	return &JsonSchemaMetadata{Types: make(map[string]*DataTypeMetadata), Root: nil}
}

func NewDataTypeRefMetadata(ref *DataTypeMetadata) *DataTypeMetadata {
	return &DataTypeMetadata{
		Ref: ref,
	}
}

func NewDataTypeMetadata(packageName, typeName, typeKind string, isPointer bool) *DataTypeMetadata {
	return &DataTypeMetadata{
		Package:   packageName,
		TypeName:  typeName,
		TypeKind:  typeKind,
		IsPointer: isPointer,
	}
}

func (m *DataTypeMetadata) ID() string {
	return fmt.Sprintf("%s#%s", m.Package, m.TypeName)
}

func NewDataTypeMetadataWithBaseMetadata(metadata *DataTypeMetadata, packageName, typeName, typeKind string, isPointer bool) *DataTypeMetadata {
	if metadata == nil {
		return &DataTypeMetadata{
			Package:   packageName,
			TypeName:  typeName,
			TypeKind:  typeKind,
			IsPointer: isPointer,
		}
	}
	metadata.Package = packageName
	metadata.TypeName = typeName
	metadata.TypeKind = typeKind
	metadata.IsPointer = isPointer

	return metadata
}
