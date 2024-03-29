package webapp

import "net/http"

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(c Context) error

func (h HandlerFunc) Handle(c Context) error {
	return h(c)
}

// Handler defines the interface to implement to serve HTTP requests.
type Handler interface {
	Handle(c Context) error
}

// WrapHandler wraps `http.Handler` into `webapp.HandlerFunc`.
func WrapHandler(h http.Handler) HandlerFunc {
	return func(c Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
