package router

import (
	"github.com/mbict/webapp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func Test_webapp_registered_routes(t *testing.T) {
//
//	app := webapp.New(webapp.WithRouter(New()))
//
//	app.GET("/test", func(c webapp.Context) error { return nil })
//	app.POST("/test", func(c webapp.Context) error { return nil })
//	g := app.Group("/test")
//	{
//		g.GET("/foo", func(c webapp.Context) error { return nil })
//		b := g.Group("/bar")
//		{
//			b.GET("/blup", func(c webapp.Context) error { return nil })
//		}
//	}
//	routes := app.Routes().All()
//
//	require.Len(t, routes, 4)
//
//	assert.Equal(t, "webapp.Test_webapp_registered_routes.func1", routes[0].Name())
//	assert.Equal(t, "/test", routes[0].Path())
//	assert.Equal(t, "GET", routes[0].Method())
//	assert.Equal(t, []string{}, routes[0].Params())
//
//	assert.Equal(t, "webapp.Test_webapp_registered_routes.func2", routes[1].Name())
//	assert.Equal(t, "/test", routes[1].Path())
//	assert.Equal(t, "POST", routes[1].Method())
//	assert.Equal(t, []string{}, routes[1].Params())
//
//	assert.Equal(t, "webapp.Test_webapp_registered_routes.func3", routes[2].Name())
//	assert.Equal(t, "/test/foo", routes[2].Path())
//	assert.Equal(t, "GET", routes[2].Method())
//	assert.Equal(t, []string{}, routes[2].Params())
//
//	assert.Equal(t, "webapp.Test_webapp_registered_routes.func4", routes[3].Name())
//	assert.Equal(t, "/test/bar/blup", routes[3].Path())
//	assert.Equal(t, "GET", routes[3].Method())
//	assert.Equal(t, []string{}, routes[3].Params())
//}

func Test_webapp_handler_gets_called(t *testing.T) {
	app := webapp.New(webapp.WithRouter(New()))

	handlerCalled := false
	app.GET("/test", func(c webapp.Context) error { handlerCalled = true; return nil })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.True(t, handlerCalled)
	assert.Equal(t, 200, rw.Code)
}

func Test_webapp_404_http_error_on_no_matching_route(t *testing.T) {
	app := webapp.New(webapp.WithRouter(New()))

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Equal(t, 404, rw.Code)
}

func Test_webapp_group_route(t *testing.T) {
	app := webapp.New(webapp.WithRouter(New()))

	handlerCalled := false
	app.Group("/foo").Group("/bar").GET("/test", func(c webapp.Context) error { handlerCalled = true; return nil })

	req, _ := http.NewRequest(http.MethodGet, "/foo/bar/test", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.True(t, handlerCalled)
	assert.Equal(t, 200, rw.Code)
}

func Test_webapp_middleware_gets_called(t *testing.T) {
	var callstack []string
	mockedMiddleware := func(name string) webapp.MiddlewareFunc {
		return func(next webapp.HandlerFunc) webapp.HandlerFunc {
			return func(c webapp.Context) error {
				callstack = append(callstack, name)

				return next(c)
			}

		}
	}

	app := webapp.New(webapp.WithRouter(New()))
	app.Use(mockedMiddleware("main"))
	app.Group("/foo", mockedMiddleware("foo")).
		Group("/bar", mockedMiddleware("bar")).
		GET("/test", func(c webapp.Context) error { callstack = append(callstack, "handler"); return nil }, mockedMiddleware("test"))

	req, _ := http.NewRequest(http.MethodGet, "/foo/bar/test", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Equal(t, []string{"main", "foo", "bar", "test", "handler"}, callstack)
	assert.Equal(t, 200, rw.Code)
}

func Test_webapp_no_middleware_gets_called_on_no_found_route(t *testing.T) {
	var callstack []string
	mockedMiddleware := func(name string) webapp.MiddlewareFunc {
		return func(next webapp.HandlerFunc) webapp.HandlerFunc {
			return func(c webapp.Context) error {
				callstack = append(callstack, name)

				return next(c)
			}
		}
	}

	app := webapp.New(webapp.WithRouter(New()))
	app.Use(mockedMiddleware("main"))
	app.Group("/foo", mockedMiddleware("foo")).
		Group("/bar", mockedMiddleware("bar")).
		GET("/test", func(c webapp.Context) error { callstack = append(callstack, "handler"); return nil }, mockedMiddleware("test"))

	req, _ := http.NewRequest(http.MethodGet, "/foo/bar/baz", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Len(t, callstack, 0)
	assert.Equal(t, 404, rw.Code)
}

func Test_webapp_path_params(t *testing.T) {
	app := webapp.New(webapp.WithRouter(New()))
	app.GET("/test/{foo}/{bar}", func(c webapp.Context) error {

		assert.Equal(t, []string{"foo", "bar"}, c.ParamNames())
		assert.Equal(t, []string{"bar", "baz"}, c.ParamValues())

		return nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/test/bar/baz", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Equal(t, 200, rw.Code)
}

func Test_webapp_path_params_with_action_keyword(t *testing.T) {
	app := webapp.New(webapp.WithRouter(New()))
	app.GET("/test/{foo}:action", func(c webapp.Context) error {

		//TODO: make the names filled in
		assert.Equal(t, []string{"foo"}, c.ParamNames())
		assert.Equal(t, []string{"bar"}, c.ParamValues())

		return nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/test/bar:action", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Equal(t, 200, rw.Code)
}

func Test_webapp_path_params_with_action_keyword_no_override(t *testing.T) {
	app := webapp.New(webapp.WithRouter(New()))
	app.GET("/test/{foo}", func(c webapp.Context) error {

		//TODO: make the names filled in
		assert.Equal(t, []string{"foo"}, c.ParamNames())
		assert.Equal(t, []string{"aabc-id-here"}, c.ParamValues())

		return nil
	})

	app.GET("/test/{foo}:action", func(c webapp.Context) error {

		assert.Fail(t, "should not be called")

		return nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/test/aabc-id-here", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Equal(t, 200, rw.Code)
}

func Test_webapp_path_params_with_main_no_override_action_keyword(t *testing.T) {
	app := webapp.New(webapp.WithRouter(New()))
	app.GET("/test/{foo}", func(c webapp.Context) error {
		assert.Fail(t, "should not be called")
		return nil
	})

	app.GET("/test/{foo}:action", func(c webapp.Context) error {
		//TODO: make the names filled in
		assert.Equal(t, []string{"foo"}, c.ParamNames())
		assert.Equal(t, []string{"aabc-id-here"}, c.ParamValues())

		return nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/test/aabc-id-here:action", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Equal(t, 200, rw.Code)
}

func BenchmarkSimpleRouteNoAllocations(b *testing.B) {

	app := webapp.New(webapp.WithRouter(New()))
	app.GET("/test/@foo/@bar", func(c webapp.Context) error {
		return nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/test/bar/baz", nil)
	rw := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			app.ServeHTTP(rw, req)
		}
	}
}
