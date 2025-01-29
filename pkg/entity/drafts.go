package entity

// DraftVersion defines supported draft versions for JSON Schema
type DraftVersion string

const (
	Draft04     DraftVersion = "http://json-schema.org/draft-04/schema#"
	Draft06     DraftVersion = "http://json-schema.org/draft-06/schema#"
	Draft07     DraftVersion = "http://json-schema.org/draft-07/schema#"
	Draft201909 DraftVersion = "https://json-schema.org/draft/2019-09/schema"
	Draft202012 DraftVersion = "https://json-schema.org/draft/2020-12/schema"
)
