package container

import "errors"

var ErrCannotResolveInstance = errors.New("cannot resolve instance")

type Container interface {
	Invoke(function interface{})
	InvokeE(function interface{}) error
}
