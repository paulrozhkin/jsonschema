package base

import (
	"github.com/paulrozhkin/jsonschema/tests/parser/additional"
)

type Settings struct {
	ValInnerSettings additional.InnerSettings  `json:"valInnerSettings"`
	RefInnerSettings *additional.InnerSettings `json:"refInnerSettings,omitempty"`
	FloatValue       float32                   `json:"floatValue"`
}
