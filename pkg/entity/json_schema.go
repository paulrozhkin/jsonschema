package entity

// BaseSchema represents common fields for all schemas
type BaseSchema[T any] struct {
	Type        string        `json:"type"`                  // All DraftVersion
	Title       *string       `json:"title,omitempty"`       // All DraftVersion
	Description *string       `json:"description,omitempty"` // All DraftVersion
	Default     *T            `json:"default,omitempty"`     // DraftVersion-06 and later
	Examples    []T           `json:"examples,omitempty"`    // DraftVersion-06 and later
	Const       T             `json:"const,omitempty"`       // DraftVersion-06 and later
	Enum        []T           `json:"enum,omitempty"`        // All drafts
	AllOf       []*JSONSchema `json:"allOf,omitempty"`       // DraftVersion-04 and later
	AnyOf       []*JSONSchema `json:"anyOf,omitempty"`       // DraftVersion-04 and later
	OneOf       []*JSONSchema `json:"oneOf,omitempty"`       // DraftVersion-04 and later
	Not         *JSONSchema   `json:"not,omitempty"`         // DraftVersion-04 and later
	Comment     *string       `json:"$comment,omitempty"`    // DraftVersion-07 and later
	Deprecated  *bool         `json:"deprecated,omitempty"`  // DraftVersion-2019-09 and later
}

// JSONSchema represents the top-level structure of a JSON Schema
type JSONSchema struct {
	BaseSchema[any]
	DeprecatedID *string                `json:"id,omitempty"`          // DraftVersion-04 and later
	ID           *string                `json:"$id,omitempty"`         // DraftVersion-06 and later
	Schema       *DraftVersion          `json:"$schema,omitempty"`     // All DraftVersion
	Defs         map[string]*JSONSchema `json:"$defs,omitempty"`       // DraftVersion-06 and later
	Ref          *string                `json:"$ref,omitempty"`        // All DraftVersion
	DynamicRef   *string                `json:"$dynamicRef,omitempty"` // DraftVersion-2019-09 and later
	Anchor       *string                `json:"$anchor,omitempty"`     // DraftVersion-2019-09 and later
	Vocabulary   map[string]bool        `json:"$vocabulary,omitempty"` // DraftVersion-2019-09 and later
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

// ObjectSchema represents a schema for object values
type ObjectSchema struct {
	BaseSchema[map[string]any]
	Properties           map[string]*JSONSchema `json:"properties,omitempty"`           // All DraftVersion
	PatternProperties    map[string]*JSONSchema `json:"patternProperties,omitempty"`    // DraftVersion-04 and later
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

// NewJSONSchema Builder functions for JSONSchema and other schema types
func NewJSONSchema(schema DraftVersion) *JSONSchema {
	return &JSONSchema{
		BaseSchema: BaseSchema[any]{Type: "object"},
		Schema:     &schema,
		Defs:       make(map[string]*JSONSchema),
	}
}

func (s *JSONSchema) SetID(id string) *JSONSchema {
	s.ID = &id
	return s
}

func (s *JSONSchema) SetTitle(title string) *JSONSchema {
	s.Title = &title
	return s
}

func (s *JSONSchema) SetDescription(description string) *JSONSchema {
	s.Description = &description
	return s
}

func (s *JSONSchema) AddDefinition(name string, schema *JSONSchema) *JSONSchema {
	s.Defs[name] = schema
	return s
}

func NewNumberSchema() *NumberSchema {
	schema := new(NumberSchema)
	schema.Type = "number"
	return schema
}

func NewIntegerSchema() *IntegerSchema {
	schema := new(IntegerSchema)
	schema.Type = "integer"
	return schema
}

func (s *NumberSchema) SetMultipleOf(value float64) *NumberSchema {
	s.MultipleOf = &value
	return s
}

func (s *NumberSchema) SetMaximum(value float64) *NumberSchema {
	s.Maximum = &value
	return s
}

func (s *NumberSchema) SetMinimum(value float64) *NumberSchema {
	s.Minimum = &value
	return s
}

func NewStringSchema() *StringSchema {
	return &StringSchema{
		BaseSchema: BaseSchema[string]{Type: "string"},
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
		BaseSchema: BaseSchema[map[string]any]{Type: "object"},
		Properties: make(map[string]*JSONSchema),
	}
}

func (s *ObjectSchema) AddProperty(name string, schema *JSONSchema) *ObjectSchema {
	s.Properties[name] = schema
	return s
}

func (s *ObjectSchema) SetRequired(fields ...string) *ObjectSchema {
	s.Required = fields
	return s
}
