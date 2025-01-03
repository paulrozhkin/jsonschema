package additional

type InnerSettings struct {
	StringValue string `json:"stringValue"`
	IntValue    int    `json:"intValue" jsonschema:"minimum=0,maximum=10"`
	BoolValue   bool   `json:"boolValue"`
}
