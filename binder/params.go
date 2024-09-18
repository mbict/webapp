package binder

import (
	"github.com/mbict/webapp"
)

func paramsGetterFunc(c webapp.Context) getter {
	return paramsGetter(c.Param)
}

type paramsGetter func(key string) string

func (ps paramsGetter) Get(key string) string {
	return ps(key)
}

func (ps paramsGetter) Values(key string) []string {
	return []string{ps(key)}
}
