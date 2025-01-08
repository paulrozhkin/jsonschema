package entity

import "reflect"

type DataType interface {
	// IsType return true, if schema contains type
	IsType(dataType JSONSchemaDataType) bool
	// IsTypes return true, if schema corresponds to the types
	IsTypes(dataTypes []JSONSchemaDataType) bool
}

// BaseSchema represents common fields for all schemas
type BaseSchema[T any] struct {
	Type        JSONSchemaType `json:"type,omitempty"`        // All DraftVersion
	Title       *string        `json:"title,omitempty"`       // All DraftVersion
	Description *string        `json:"description,omitempty"` // All DraftVersion
	Default     *T             `json:"default,omitempty"`     // DraftVersion-06 and later
	Examples    []*T           `json:"examples,omitempty"`    // DraftVersion-06 and later
	Const       *T             `json:"const,omitempty"`       // DraftVersion-06 and later
	Enum        []*T           `json:"enum,omitempty"`        // All drafts
	Comment     *string        `json:"$comment,omitempty"`    // DraftVersion-07 and later
	Deprecated  *bool          `json:"deprecated,omitempty"`  // DraftVersion-2019-09 and later
}

func (s *BaseSchema[T]) IsType(dataType JSONSchemaDataType) bool {
	for i := range s.Type {
		if s.Type[i] == dataType {
			return true
		}
	}
	return false
}

func (s *BaseSchema[T]) IsTypes(dataTypes []JSONSchemaDataType) bool {
	return reflect.DeepEqual(s.Type, dataTypes)
}

// ObjectSchema represents a schema for object values
type ObjectSchema struct {
	BaseSchema[map[string]any]
	Properties           map[string]DataType    `json:"properties,omitempty"`           // All DraftVersion
	PatternProperties    map[string]DataType    `json:"patternProperties,omitempty"`    // DraftVersion-04 and later
	AdditionalProperties *AdditionalProperties  `json:"additionalProperties,omitempty"` // All DraftVersion
	MaxProperties        *int                   `json:"maxProperties,omitempty"`        // All DraftVersion
	MinProperties        *int                   `json:"minProperties,omitempty"`        // All DraftVersion
	Required             []string               `json:"required,omitempty"`             // All DraftVersion
	Dependencies         map[string]*Dependency `json:"dependencies,omitempty"`         // DraftVersion-04, DraftVersion-06, DraftVersion-07
	DependentSchemas     map[string]*JSONSchema `json:"dependentSchemas,omitempty"`     // DraftVersion-2019-09 and later
	DependentRequired    map[string][]string    `json:"dependentRequired,omitempty"`    // DraftVersion-2019-09 and later
	PropertyNames        *JSONSchema            `json:"propertyNames,omitempty"`        // DraftVersion-06 and later

	// Conditional Validation
	If   *JSONSchema `json:"if,omitempty"`   // DraftVersion-07 and later
	Then *JSONSchema `json:"then,omitempty"` // DraftVersion-07 and later
	Else *JSONSchema `json:"else,omitempty"` // DraftVersion-07 and later

	AllOf []*JSONSchema `json:"allOf,omitempty"` // DraftVersion-04 and later
	AnyOf []*JSONSchema `json:"anyOf,omitempty"` // DraftVersion-04 and later
	OneOf []*JSONSchema `json:"oneOf,omitempty"` // DraftVersion-04 and later
	Not   *JSONSchema   `json:"not,omitempty"`   // DraftVersion-04 and later
}

// JSONSchema represents the top-level structure of a JSON Schema
type JSONSchema struct {
	DeprecatedID *string             `json:"id,omitempty"`          // DraftVersion-04 and later
	ID           *string             `json:"$id,omitempty"`         // DraftVersion-06 and later
	Schema       *DraftVersion       `json:"$schema,omitempty"`     // All DraftVersion
	Defs         map[string]DataType `json:"$defs,omitempty"`       // DraftVersion-06 and later
	Ref          *string             `json:"$ref,omitempty"`        // All DraftVersion
	DynamicRef   *string             `json:"$dynamicRef,omitempty"` // DraftVersion-2019-09 and later
	Anchor       *string             `json:"$anchor,omitempty"`     // DraftVersion-2019-09 and later
	Vocabulary   map[string]string   `json:"$vocabulary,omitempty"` // DraftVersion-2019-09 and later

	ObjectSchema
}

// NumericSchema represents a base schema for numeric values
type NumericSchema[T any] struct {
	BaseSchema[T]
	MultipleOf       *T `json:"multipleOf,omitempty"`       // DraftVersion-06 and later
	Maximum          *T `json:"maximum,omitempty"`          // All DraftVersion
	ExclusiveMaximum *T `json:"exclusiveMaximum,omitempty"` // DraftVersion-06 and later
	Minimum          *T `json:"minimum,omitempty"`          // All DraftVersion
	ExclusiveMinimum *T `json:"exclusiveMinimum,omitempty"` // DraftVersion-06 and later
}

// NumberSchema represents a schema for number values
type NumberSchema struct {
	NumericSchema[float64]
}

// IntegerSchema represents a schema for integer values
type IntegerSchema struct {
	NumericSchema[int]
}

// StringSchema represents a schema for string values
type StringSchema struct {
	BaseSchema[string]
	MaxLength        *int        `json:"maxLength,omitempty"`        // All DraftVersion
	MinLength        *int        `json:"minLength,omitempty"`        // All DraftVersion
	Pattern          *string     `json:"pattern,omitempty"`          // All DraftVersion
	Format           *string     `json:"format,omitempty"`           // DraftVersion-04 and later
	ContentMediaType *string     `json:"contentMediaType,omitempty"` // DraftVersion-07 and later
	ContentEncoding  *string     `json:"contentEncoding,omitempty"`  // DraftVersion-07 and later
	ContentSchema    *JSONSchema `json:"contentSchema,omitempty"`    // DraftVersion-07 and later
}

// BooleanSchema represents a schema for boolean values
type BooleanSchema struct {
	BaseSchema[bool]
}

// NullSchema represents a schema for null values
type NullSchema struct {
	BaseSchema[any]
}

// ArraySchema represents a schema for array values
type ArraySchema struct {
	BaseSchema[[]any]
	Items            *JSONSchema   `json:"items,omitempty"`            // All DraftVersion
	PrefixItems      []*JSONSchema `json:"prefixItems,omitempty"`      // DraftVersion-2020-12 and later
	Contains         *JSONSchema   `json:"contains,omitempty"`         // DraftVersion-06 and later
	MaxItems         *int          `json:"maxItems,omitempty"`         // All DraftVersion
	MinItems         *int          `json:"minItems,omitempty"`         // All DraftVersion
	UniqueItems      *bool         `json:"uniqueItems,omitempty"`      // DraftVersion-04 and later
	MinContains      *int          `json:"minContains,omitempty"`      // DraftVersion-2019-09 and later
	MaxContains      *int          `json:"maxContains,omitempty"`      // DraftVersion-2019-09 and later
	UnevaluatedItems *JSONSchema   `json:"unevaluatedItems,omitempty"` // DraftVersion-2019-09 and later
}

// AdditionalProperties represents the additionalProperties keyword
// Can be either a boolean or a JSONSchema
type AdditionalProperties struct {
	Bool   *bool       `json:"-"` // All DraftVersion
	Schema *JSONSchema `json:"-"` // All DraftVersion
}

// Dependency represents the dependencies keyword
// Can be either an array of strings or a JSONSchema
type Dependency struct {
	PropertyDependencies []string    `json:"-"` // DraftVersion-04, DraftVersion-06, DraftVersion-07
	SchemaDependency     *JSONSchema `json:"-"` // DraftVersion-04, DraftVersion-06, DraftVersion-07
}

// NewAdditionalPropertiesBool creates a new AdditionalProperties instance
func NewAdditionalPropertiesBool(value bool) *AdditionalProperties {
	return &AdditionalProperties{Bool: &value}
}

func NewAdditionalPropertiesSchema(schema *JSONSchema) *AdditionalProperties {
	return &AdditionalProperties{Schema: schema}
}

// NewDependencyProperties creates a new Dependency instance
func NewDependencyProperties(properties []string) *Dependency {
	return &Dependency{PropertyDependencies: properties}
}

func NewDependencySchema(schema *JSONSchema) *Dependency {
	return &Dependency{SchemaDependency: schema}
}

// NewJSONEmptySchema Builder functions for JSONSchema without Type
func NewJSONEmptySchema() *JSONSchema {
	return &JSONSchema{}
}

// NewJSONSchema Builder functions for JSONSchema and other schema types
func NewJSONSchema() *JSONSchema {
	return &JSONSchema{
		ObjectSchema: *NewObjectSchema(),
	}
}

func (s *JSONSchema) SetTitle(title string) *JSONSchema {
	s.Title = &title
	return s
}

func (s *JSONSchema) SetID(id string) *JSONSchema {
	s.ID = &id
	return s
}

func (s *JSONSchema) SetSchemaVersion(version DraftVersion) *JSONSchema {
	s.Schema = &version
	return s
}

func (s *JSONSchema) SetDescription(description string) *JSONSchema {
	s.Description = &description
	return s
}

func (s *JSONSchema) SetSchema(draft DraftVersion) *JSONSchema {
	s.Schema = &draft
	return s
}

func (s *JSONSchema) AddDefinition(name string, schema DataType) *JSONSchema {
	if s.Defs == nil {
		s.Defs = make(map[string]DataType)
	}
	s.Defs[name] = schema
	return s
}

func (s *JSONSchema) SetRef(ref string) *JSONSchema {
	s.Ref = &ref
	return s
}

func NewNumberSchema() *NumberSchema {
	schema := new(NumberSchema)
	schema.Type = JSONSchemaType{JSONSchemaNumber}
	return schema
}

func NewIntegerSchema() *IntegerSchema {
	schema := new(IntegerSchema)
	schema.Type = JSONSchemaType{JSONSchemaInteger}
	return schema
}

func (s *IntegerSchema) SetMultipleOf(value int) *IntegerSchema {
	s.MultipleOf = &value
	return s
}

func (s *IntegerSchema) SetMaximum(value int) *IntegerSchema {
	s.Maximum = &value
	return s
}

func (s *IntegerSchema) SetMinimum(value int) *IntegerSchema {
	s.Minimum = &value
	return s
}

func NewStringSchema() *StringSchema {
	return &StringSchema{
		BaseSchema: BaseSchema[string]{Type: JSONSchemaType{JSONSchemaString}},
	}
}

func (s *StringSchema) SetMaxLength(length int) *StringSchema {
	s.MaxLength = &length
	return s
}

func (s *StringSchema) SetPattern(pattern string) *StringSchema {
	s.Pattern = &pattern
	return s
}

func NewObjectSchema() *ObjectSchema {
	return &ObjectSchema{
		BaseSchema: BaseSchema[map[string]any]{Type: JSONSchemaType{JSONSchemaObject}},
	}
}

func (s *ObjectSchema) AddProperty(name string, schema DataType) *ObjectSchema {
	if s.Properties == nil {
		s.Properties = make(map[string]DataType)
	}
	s.Properties[name] = schema
	return s
}

func (s *ObjectSchema) AddRequired(requiredProperties ...string) *ObjectSchema {
	s.Required = append(s.Required, requiredProperties...)
	return s
}

func (s *ObjectSchema) SetRequired(fields ...string) *ObjectSchema {
	s.Required = fields
	return s
}

func NewBooleanSchema() *BooleanSchema {
	schema := new(BooleanSchema)
	schema.Type = JSONSchemaType{JSONSchemaBoolean}
	return schema
}
