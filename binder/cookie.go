package binder

import (
	"github.com/mbict/webapp"
	"net/http"
)

func cookieGetterFunc(c webapp.Context) getter {
	return cookieGetter(c.Request().Cookies())
}

type cookieGetter []*http.Cookie

func (c cookieGetter) Get(key string) string {
	for i := range c {
		if c[i].Name == key {
			return c[i].Value
		}
	}
	return ""
}

func (c cookieGetter) Values(key string) []string {
	var res []string
	for i := range c {
		if c[i].Name == key {
			res = append(res, c[i].Value)
		}
	}
	return res
}
