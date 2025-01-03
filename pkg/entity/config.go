package entity

type Config struct {
	// SchemaVersion determinate draft that will be generate
	SchemaVersion Draft

	// FieldNameTag will change the tag used to get field names. json tags are used by default.
	FieldNameTag []string

	// Namer allows customizing of type names. The default is to use the type's name
	// provided by the reflect package.
	//Namer func(reflect.Type) string

	// KeyNamer allows customizing of key names.
	// The default is to use the key's name as is, or the json tag if present.
	// If a json tag is present, KeyNamer will receive the tag's name as an argument, not the original key name.
	KeyNamer func(string) string
}
