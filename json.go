package webappv2

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DefaultJSONEncoder implements JSON encoding using encoding/json.
type DefaultJSONEncoder struct{}

// Encode converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d DefaultJSONEncoder) Encode(c Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Decode reads a JSON from a request body and converts it into an interface.
func (d DefaultJSONEncoder) Decode(c Context, i interface{}) error {
	err := json.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return NewHTTPErrorWithInternal(http.StatusBadRequest, err, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset))
	} else if se, ok := err.(*json.SyntaxError); ok {
		return NewHTTPErrorWithInternal(http.StatusBadRequest, err, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()))
	}
	return err
}
