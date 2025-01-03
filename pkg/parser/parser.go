package parser

import "github.com/paulrozhkin/jsonschema/pkg/entity"

type Parser interface {
	Parse() (*entity.TypeMetadata, error)
}
