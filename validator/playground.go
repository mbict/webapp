package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/mbict/webapp"
)

type playgroundValidator struct {
	validator *validator.Validate
}

type PlaygroundValidatorOption func(validate *validator.Validate)

func WithCustomValidator(name string, fn validator.Func, callValidationEvenIfNull ...bool) PlaygroundValidatorOption {
	return func(v *validator.Validate) {
		v.RegisterValidation(name, fn, callValidationEvenIfNull...)
	}
}

func (p *playgroundValidator) Validate(i interface{}) error {
	return p.validator.Struct(i)
}

func PlaygroundValidator(options ...PlaygroundValidatorOption) webapp.Validator {
	v := validator.New()

	for _, option := range options {
		option(v)
	}

	return &playgroundValidator{
		validator: v,
	}
}
