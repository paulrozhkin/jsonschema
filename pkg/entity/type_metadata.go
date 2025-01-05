package entity

type TypeMetadata struct {
	Package     string
	TypeName    string
	TypeKind    string
	FieldName   string
	Nodes       []*TypeMetadata
	Tags        map[string][]string
	IsReference bool
}
