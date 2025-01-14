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

func ExpectedSettingsMetadata() *entity.JsonSchemaMetadata {
	jsonSchemaMetadata := entity.NewJsonSchemaMetadata()
	packageAdditional := reflect.TypeOf(additional.InnerSettings{}).PkgPath()
	packageBase := reflect.TypeOf(Settings{}).PkgPath()
	innerSettingsNodeMetadata := []*entity.DataTypeMetadata{
		{FieldName: "StringValue", Package: "", TypeName: "string", TypeKind: "string", Nodes: nil, Tags: map[string][]string{"json": {"stringValue"}}, IsPointer: false},
		{FieldName: "IntValue", Package: "", TypeName: "int", TypeKind: "int", Nodes: nil, Tags: map[string][]string{"json": {"intValue"}, "jsonschema": {"minimum=0", "maximum=10"}}, IsPointer: false},
		{FieldName: "BoolValue", Package: "", TypeName: "bool", TypeKind: "bool", Nodes: nil, Tags: map[string][]string{"json": {"boolValue"}}, IsPointer: false},
	}
	innerSettingsMetadata := &entity.DataTypeMetadata{
		Package:   packageAdditional,
		TypeName:  "InnerSettings",
		TypeKind:  "struct",
		Nodes:     innerSettingsNodeMetadata,
		FieldName: "",
		Tags:      nil,
		IsPointer: false,
	}
	jsonSchemaMetadata.Types[fmt.Sprintf("%s#%s", packageAdditional, "InnerSettings")] = innerSettingsMetadata

	rootSettingsMetadata := &entity.DataTypeMetadata{
		TypeName: "Settings",
		Package:  packageBase,
		TypeKind: "struct",
		Nodes: []*entity.DataTypeMetadata{
			{FieldName: "ValInnerSettings", Ref: innerSettingsMetadata, Tags: map[string][]string{"json": {"valInnerSettings"}}, IsPointer: false},
			{FieldName: "RefInnerSettings", Ref: innerSettingsMetadata, Tags: map[string][]string{"json": {"refInnerSettings", "omitempty"}}, IsPointer: true},
			{FieldName: "FloatValue", Package: "", TypeName: "float32", TypeKind: "float32", Nodes: nil, Tags: map[string][]string{"json": {"floatValue"}}, IsPointer: false},
		},
		Tags:      nil,
		IsPointer: false,
	}
	jsonSchemaMetadata.Types[fmt.Sprintf("%s#%s", packageBase, "Settings")] = rootSettingsMetadata
	jsonSchemaMetadata.Root = rootSettingsMetadata
	return jsonSchemaMetadata
}

func ExpectedSettingsJsonSchema() *entity.JSONSchema {
	packageBase := reflect.TypeOf(Settings{}).PkgPath()

	floatValueSchema := entity.NewNumberSchema()
	innerSettingsSchema := entity.NewObjectSchema().
		AddProperty("stringValue", entity.NewStringSchema()).
		AddProperty("intValue", entity.NewIntegerSchema().SetMinimum(0).SetMaximum(10)).
		AddProperty("boolValue", entity.NewBooleanSchema()).
		AddRequired("stringValue", "intValue", "boolValue")
	innerSettingsId := fmt.Sprintf("#/$defs/%s", "InnerSettings")

	schema := entity.NewJSONSchema()
	schema.
		SetSchema(entity.Draft202012).
		SetID(fmt.Sprintf("https://%s/%s", packageBase, "Settings")).
		AddDefinition("InnerSettings", innerSettingsSchema).
		AddProperty("valInnerSettings", entity.NewJSONEmptySchema().SetRef(innerSettingsId)).
		AddProperty("refInnerSettings", entity.NewJSONEmptySchema().SetRef(innerSettingsId)).
		AddProperty("floatValue", floatValueSchema).
		AddRequired("valInnerSettings", "floatValue")
	return schema

}
