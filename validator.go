package webappv2

// Validator is the interface that wraps the Validate function.
type Validator interface {
	Validate(i interface{}) error
}
