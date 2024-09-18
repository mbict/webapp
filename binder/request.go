package binder

import (
	"github.com/mbict/webapp"
)

func requestGetterFunc(c webapp.Context) getter {
	return &requestGetter{Context: c}
}

type requestGetter struct {
	webapp.Context
}

func (r *requestGetter) Get(key string) string {
	switch key {
	case `remote-addr`:
		return (*r).RealIP()
	case `host`:
		return (*r).Request().Host
	case `method`:
		return (*r).Method()
	case `url`:
		return (*r).Request().URL.String()
	case `url:host`:
		return (*r).Request().URL.Host
	case `url:query`:
		return (*r).Request().URL.RawQuery
	case `url:path`:
		return (*r).Path()
	case `url:scheme`:
		return (*r).Request().URL.Scheme
	}

	return ""
}

func (r *requestGetter) Values(key string) []string {
	if v := r.Get(key); v != "" {
		return []string{v}
	}
	return []string{}
}
