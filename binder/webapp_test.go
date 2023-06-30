package binder

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"webappv2"
)

/*func Test_webapp_handler_gets_called(T *testing.T) {
	app := webappv2.New(webappv2.WithRouter(router.New()), webappv2.WithBinder(Binder()))

	handlerCalled := false
	app.GET("/test", func(c webappv2.Context) error { handlerCalled = true; return nil })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.True(T, handlerCalled)
	assert.Equal(T, 200, rw.Code)
}

func Test_webapp_404_http_error_on_no_matching_route(T *testing.T) {
	app := webappv2.New(webappv2.WithRouter(New()))

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Equal(T, 404, rw.Code)
}

func Test_webapp_group_route(T *testing.T) {
	app := webappv2.New(webappv2.WithRouter(New()))

	handlerCalled := false
	app.Group("/foo").Group("/bar").GET("/test", func(c webappv2.Context) error { handlerCalled = true; return nil })

	req, _ := http.NewRequest(http.MethodGet, "/foo/bar/test", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.True(T, handlerCalled)
	assert.Equal(T, 200, rw.Code)
}

func Test_webapp_middleware_gets_called(T *testing.T) {
	var callstack []string
	mockedMiddleware := func(name string) webappv2.MiddlewareFunc {
		return func(next webappv2.HandlerFunc) webappv2.HandlerFunc {
			return func(c webappv2.Context) error {
				callstack = append(callstack, name)

				return next(c)
			}

		}
	}

	app := webappv2.New(webappv2.WithRouter(New()))
	app.Use(mockedMiddleware("main"))
	app.Group("/foo", mockedMiddleware("foo")).
		Group("/bar", mockedMiddleware("bar")).
		GET("/test", func(c webappv2.Context) error { callstack = append(callstack, "handler"); return nil }, mockedMiddleware("test"))

	req, _ := http.NewRequest(http.MethodGet, "/foo/bar/test", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Equal(T, []string{"main", "foo", "bar", "test", "handler"}, callstack)
	assert.Equal(T, 200, rw.Code)
}

func Test_webapp_no_middleware_gets_called_on_no_found_route(T *testing.T) {
	var callstack []string
	mockedMiddleware := func(name string) webappv2.MiddlewareFunc {
		return func(next webappv2.HandlerFunc) webappv2.HandlerFunc {
			return func(c webappv2.Context) error {
				callstack = append(callstack, name)

				return next(c)
			}
		}
	}

	app := webappv2.New(webappv2.WithRouter(New()))
	app.Use(mockedMiddleware("main"))
	app.Group("/foo", mockedMiddleware("foo")).
		Group("/bar", mockedMiddleware("bar")).
		GET("/test", func(c webappv2.Context) error { callstack = append(callstack, "handler"); return nil }, mockedMiddleware("test"))

	req, _ := http.NewRequest(http.MethodGet, "/foo/bar/baz", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Len(T, callstack, 0)
	assert.Equal(T, 404, rw.Code)
}

func Test_webapp_path_params(T *testing.T) {
	app := webappv2.New(webappv2.WithRouter(New()))
	app.GET("/test/@foo/@bar", func(c webappv2.Context) error {

		//TODO: make the names filled in
		assert.Equal(T, []string{"", ""}, c.ParamNames())
		assert.Equal(T, []string{"bar", "baz"}, c.ParamValues())

		return nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/test/bar/baz", nil)
	rw := httptest.NewRecorder()
	app.ServeHTTP(rw, req)

	assert.Equal(T, 200, rw.Code)
}*/

func BenchmarkSimpleRouteNoAllocations(b *testing.B) {

	app := webappv2.New(webappv2.WithRouter(DefaultBinder()))
	app.GET("/test/@foo/@bar", func(c webappv2.Context) error {
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
