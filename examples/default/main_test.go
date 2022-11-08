package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"webappv2"
	"webappv2/router"
)

func BenchmarkSimpleRouteNoAllocations(b *testing.B) {

	app := webappv2.New(
		webappv2.WithRouter(router.New()), // custom router
	)

	app.GET("/intention/{id}:activate", func(c webappv2.Context) error {
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
