package webapp

// JSONEncoding is the interface that encodes and decodes JSON to and from interfaces.
type JSONEncoding interface {
	Encode(c Context, i interface{}, indent string) error
	Decode(c Context, i interface{}) error
}

// XMLEncoding is the interface that encodes and decodes XML to and from interfaces.
type XMLEncoding interface {
	Encode(c Context, i interface{}, indent string) error
	Decode(c Context, i interface{}) error
}
