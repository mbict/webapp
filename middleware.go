package webappv2

import "net/http"

// MiddlewareFunc defines a function to process middleware.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Middleware interface {
	Handle(next HandlerFunc) HandlerFunc
}

// WrapMiddleware wraps `func(http.Handler) http.Handler` into `webapp.MiddlewareFunc`
func WrapMiddleware(m func(http.Handler) http.Handler) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) (err error) {
			m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				c.SetResponse(newResponse(w))
				err = next(c)
			})).ServeHTTP(c.Response(), c.Request())
			return
		}
	}
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
