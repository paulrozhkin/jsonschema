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

func NewDataTypeMetadata(packageName, typeName, typeKind string, isReference bool) *DataTypeMetadata {
	return &DataTypeMetadata{
		Package:   packageName,
		TypeName:  typeName,
		TypeKind:  typeKind,
		IsPointer: isReference,
	}
}

func (m *DataTypeMetadata) ID() string {
	return fmt.Sprintf("%s#%s", m.Package, m.TypeName)
}
