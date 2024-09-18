package webapp

import (
	"errors"
	"net/http"
	"strings"
)

// Validator is the interface that wraps the Validate function.
type Validator interface {
	//Validate could return an error of the type ValidationError or ValidationErrors on validation errors. If so you
	//return the error from your handlers and let the error handler handle the error as a bad request
	Validate(i interface{}) error
}

type ValidationError struct {
	Field           string   `json:"field,omitempty"`
	Message         string   `json:"message"`
	Validator       string   `json:"validator"`
	ConditionParams []string `json:"condition_param,omitempty"`
}

func (e ValidationError) Error() string {
	return e.Message
}

type ValidationErrors map[string][]ValidationError

func (e ValidationErrors) StatusCode() int {
	return http.StatusBadRequest
}

func (e ValidationErrors) Add(field, message, validator string, params ...string) ValidationErrors {

	//leave out the empty params
	if len(params) == 1 && params[0] == "" {
		params = []string{}
	}

	e[field] = append(e[field], NewValidationError(message, validator, params...))
	return e
}

func (e ValidationErrors) Error() string {
	messages := []string{}
	for _, errs := range e {
		for _, err := range errs {
			messages = append(messages, err.Message)
		}
	}
	return "validation errors: " + strings.Join(messages, ", ")
}

func NewValidationErrors() ValidationErrors {
	return make(ValidationErrors, 1)
}

func NewFieldValidationError(field, message, validator string, params ...string) ValidationError {
	return ValidationError{
		Field:           field,
		Message:         message,
		Validator:       validator,
		ConditionParams: params,
	}
}

func NewValidationError(message, validator string, params ...string) ValidationError {
	return ValidationError{
		Message:         message,
		Validator:       validator,
		ConditionParams: params,
	}
}

func IsValidationError(err error) bool {
	switch err.(type) {
	case ValidationError:
		return true
	case ValidationErrors:
		return true
	}

	if err = errors.Unwrap(err); err != nil {
		return IsValidationError(err)
	}
	return false
}
