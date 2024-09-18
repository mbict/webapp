package binder

import (
	"github.com/mbict/webapp"
)

func queryGetterFunc(c webapp.Context) getter {
	return mapGetter(c.Request().URL.Query())
}
