package entity

import (
	"fmt"
	"strings"
)

type TypeMetadata struct {
	Package     string
	TypeName    string
	TypeKind    string
	FieldName   string
	Nodes       []*TypeMetadata
	Tags        map[string][]string
	IsReference bool
}

func (tm *TypeMetadata) String(indentLevel int) string {
	var sb strings.Builder
	indent := strings.Repeat("  ", indentLevel)

	sb.WriteString(fmt.Sprintf("%sPackage: %s\n", indent, tm.Package))
	sb.WriteString(fmt.Sprintf("%sTypeName: %s\n", indent, tm.TypeName))
	sb.WriteString(fmt.Sprintf("%sTypeKind: %s\n", indent, tm.TypeKind))
	sb.WriteString(fmt.Sprintf("%sFieldName: %s\n", indent, tm.FieldName))
	sb.WriteString(fmt.Sprintf("%sIsReference: %t\n", indent, tm.IsReference))

	if len(tm.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("%sTags:\n", indent))
		for key, values := range tm.Tags {
			sb.WriteString(fmt.Sprintf("%s  %s: %v\n", indent, key, values))
		}
	}

	if len(tm.Nodes) > 0 {
		sb.WriteString(fmt.Sprintf("%sNodes:\n", indent))
		for _, node := range tm.Nodes {
			sb.WriteString(node.String(indentLevel + 1))
		}
	}

	return sb.String()
}
