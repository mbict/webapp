package binder

import "github.com/mbict/webapp"

func headerGetterFunc(c webapp.Context) getter {
	return mapGetter(c.Request().Header)
}
