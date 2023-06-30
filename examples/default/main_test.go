package main

import (
	"github.com/mbict/webapp"
	"github.com/mbict/webapp/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkSimpleRouteNoAllocations(b *testing.B) {

	app := webapp.New(
		webapp.WithRouter(router.New()), // custom router
	)

	app.GET("/intention/{id}:activate", func(c webapp.Context) error {
		return c.String(200, "hi")
	})

	req, _ := http.NewRequest(http.MethodGet, "/intention/1234:activate", nil)
	rw := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			app.ServeHTTP(rw, req)
		}
	}
}
