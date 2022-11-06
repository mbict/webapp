package webappv2

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_webapp_registered_routes(T *testing.T) {

	webapp := New()

	webapp.GET("/test", func(c Context) error { return nil })
	webapp.POST("/test", func(c Context) error { return nil })
	g := webapp.Group("/test")
	{
		g.GET("/foo", func(c Context) error { return nil })
		b := g.Group("/bar")
		{
			b.GET("/blup", func(c Context) error { return nil })
		}
	}
	routes := webapp.Routes().All()

	require.Len(T, routes, 4)

	assert.Equal(T, "webappv2.Test_webapp_registered_routes.func1", routes[0].Name())
	assert.Equal(T, "/test", routes[0].Path())
	assert.Equal(T, "GET", routes[0].Method())
	assert.Equal(T, []string{}, routes[0].Params())

	assert.Equal(T, "webappv2.Test_webapp_registered_routes.func2", routes[1].Name())
	assert.Equal(T, "/test", routes[1].Path())
	assert.Equal(T, "POST", routes[1].Method())
	assert.Equal(T, []string{}, routes[1].Params())

	assert.Equal(T, "webappv2.Test_webapp_registered_routes.func3", routes[2].Name())
	assert.Equal(T, "/test/foo", routes[2].Path())
	assert.Equal(T, "GET", routes[2].Method())
	assert.Equal(T, []string{}, routes[2].Params())

	assert.Equal(T, "webappv2.Test_webapp_registered_routes.func4", routes[3].Name())
	assert.Equal(T, "/test/bar/blup", routes[3].Path())
	assert.Equal(T, "GET", routes[3].Method())
	assert.Equal(T, []string{}, routes[3].Params())
}

func Test_webapp_pre_runs_on_found_route(T *testing.T) {
	webapp := New()

	preDidRun := false
	webapp.Pre(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			preDidRun = true
			return next(c)
		}
	})
	webapp.GET("/test", func(c Context) error { return nil })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	rw := httptest.NewRecorder()
	webapp.ServeHTTP(rw, req)

	assert.True(T, preDidRun)
	assert.Equal(T, 200, rw.Code)
}

func Test_webapp_pre_runs_on_not_found_route(T *testing.T) {
	webapp := New()

	preDidRun := false
	webapp.Pre(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			preDidRun = true
			return next(c)
		}
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	rw := httptest.NewRecorder()
	webapp.ServeHTTP(rw, req)

	assert.True(T, preDidRun)
	assert.Equal(T, 404, rw.Code)
}
