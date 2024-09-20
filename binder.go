package webapp

import (
	"errors"
)

var ErrUnsupportedType = errors.New("decoder: unsupported type")

var DefaultBinder Binder

type Binder interface {
	Bind(c Context, i interface{}) error
	BindBody(c Context, i interface{}) error
	BindQueryParams(c Context, i interface{}) error
}

func NewBindError(err error) BindError {
	return BindError{err: err}
}

type BindError struct {
	err error
}

func (b BindError) Error() string {
	return b.err.Error()
}

func (b BindError) Unwrap() error {
	return b.err
}

func IsBindError(err error) bool {
	switch err.(type) {
	case BindError:
		return true
	}

	if err = errors.Unwrap(err); err != nil {
		return IsBindError(err)
	}
	return false
}
