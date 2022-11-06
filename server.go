package webappv2

import (
	stdContext "context"
	"net/http"
)

type Server interface {

	// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
	ServeHTTP(w http.ResponseWriter, r *http.Request)

	// Start starts an HTTP server.
	Start(address string) error

	// StartServer starts a custom http server.
	StartServer(s *http.Server) (err error)

	// Close immediately stops the server.
	// It internally calls `http.Server#Close()`.
	Close() error

	// Shutdown stops the server gracefully.
	// It internally calls `http.Server#Shutdown()`.
	Shutdown(ctx stdContext.Context) error
}
