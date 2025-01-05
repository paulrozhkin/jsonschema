package base

import (
	"fmt"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
	"github.com/paulrozhkin/jsonschema/tests/additional"
	"reflect"
)

type Settings struct {
	ValInnerSettings additional.InnerSettings  `json:"valInnerSettings"`
	RefInnerSettings *additional.InnerSettings `json:"refInnerSettings,omitempty"`
	FloatValue       float32                   `json:"floatValue"`
}

func ExpectedSettingsMetadata() *entity.TypeMetadata {
	innerSettingsMetadata := []*entity.TypeMetadata{
		{FieldName: "StringValue", Package: "", TypeName: "string", TypeKind: "string", Nodes: nil, Tags: map[string][]string{"json": {"stringValue"}}, IsReference: false},
		{FieldName: "IntValue", Package: "", TypeName: "int", TypeKind: "int", Nodes: nil, Tags: map[string][]string{"json": {"intValue"}, "jsonschema": {"minimum=0", "maximum=10"}}, IsReference: false},
		{FieldName: "BoolValue", Package: "", TypeName: "bool", TypeKind: "bool", Nodes: nil, Tags: map[string][]string{"json": {"boolValue"}}, IsReference: false},
	}
	packageAdditional := reflect.TypeOf(additional.InnerSettings{}).PkgPath()
	packageBase := reflect.TypeOf(Settings{}).PkgPath()
	return &entity.TypeMetadata{
		TypeName: "Settings",
		Package:  packageBase,
		TypeKind: "struct",
		Nodes: []*entity.TypeMetadata{
			{FieldName: "ValInnerSettings", Package: packageAdditional, TypeName: "InnerSettings", TypeKind: "struct", Nodes: innerSettingsMetadata, Tags: map[string][]string{"json": {"valInnerSettings"}}, IsReference: false},
			{FieldName: "RefInnerSettings", Package: packageAdditional, TypeName: "InnerSettings", TypeKind: "struct", Nodes: innerSettingsMetadata, Tags: map[string][]string{"json": {"refInnerSettings", "omitempty"}}, IsReference: true},
			{FieldName: "FloatValue", Package: "", TypeName: "float32", TypeKind: "float32", Nodes: nil, Tags: map[string][]string{"json": {"floatValue"}}, IsReference: false},
		},
		Tags:        nil,
		IsReference: false,
	}
}

func ExpectedSettingsJsonSchema() *entity.JSONSchema {
	//packageAdditional := reflect.TypeOf(additional.InnerSettings{}).PkgPath()
	packageBase := reflect.TypeOf(Settings{}).PkgPath()

	floatValueSchema := entity.NewNumberSchema()

	schema := entity.NewJSONSchema().
		SetSchema(entity.Draft202012).
		SetID(fmt.Sprintf("https://%s/%s", packageBase, "Settings"))

	schema.AddProperty("floatValue", floatValueSchema)

	return schema

}
