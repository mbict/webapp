package router

import (
	"github.com/mbict/webapp"
)

type Option func(router *Router)

func WithNotFoundHandler(h webapp.Handler) Option {
	return func(r *Router) {
		r.NotFound = h
	}
}
