package binder

import (
	"github.com/mbict/webapp"
)

var defaultDefaultsGetter = defaultsGetter{}

// DefaultGetterFactory will return a value getter that always returns the key as the value
// int Property `default:"123"` will set the value 123
func defaultsGetterFunc(_ webapp.Context) getter {
	return defaultDefaultsGetter
}

type defaultsGetter struct{}

func (_ defaultsGetter) Get(key string) string {
	return key
}

func (_ defaultsGetter) Values(key string) []string {
	return []string{key}
}
