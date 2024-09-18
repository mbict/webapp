package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/mbict/webapp"
	"strings"
)

type PlaygroundValidatorOption func(validator *playgroundValidator)

func WithCustomValidator(name string, fn validator.Func, callValidationEvenIfNull ...bool) PlaygroundValidatorOption {
	return func(validator *playgroundValidator) {
		validator.validate.RegisterValidation(name, fn, callValidationEvenIfNull...)
	}
}

func WithValidator(validate *validator.Validate) PlaygroundValidatorOption {
	return func(validator *playgroundValidator) {
		validator.validate = validate
	}
}

func WithTranslator(ut ut.Translator) PlaygroundValidatorOption {
	return func(validator *playgroundValidator) {
		validator.translator = ut
	}
}

func NewPlaygroundValidator(options ...PlaygroundValidatorOption) webapp.Validator {
	translator, _ := ut.New(en.New()).GetTranslator("en")
	validate := validator.New()
	v := &playgroundValidator{
		validate:   validate,
		translator: translator,
	}

	for _, option := range options {
		option(v)
	}

	return v
}

type playgroundValidator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func (cv *playgroundValidator) Validate(i interface{}) error {
	return WrapValidationErrors(cv.validate.Struct(i), cv.translator)
}

func IsPlaygroundValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

func WrapValidationErrors(err error, translator ut.Translator) error {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		errs := webapp.NewValidationErrors()
		for _, err := range fieldErrors {
			errs = errs.Add(
				err.Field(),
				err.Translate(translator),
				err.Tag(),
				strings.Split(err.Param(), ",")...,
			)
		}
		return errs
	}
	return err
}
