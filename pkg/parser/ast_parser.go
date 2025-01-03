package parser

import "github.com/paulrozhkin/jsonschema/pkg/entity"

type AstParser struct {
	parsingObj any
}

func NewAstParser(parsingObj any) *AstParser {
	return &AstParser{parsingObj: parsingObj}
}

func (p *AstParser) Parse() (*entity.TypeMetadata, error) {
	panic("implement me")
}
