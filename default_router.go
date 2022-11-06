package webappv2

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

func NewDefaultRouter() Router {
	r := &defaultRouter{}
	r.defaultRouterGroup = &defaultRouterGroup{
		router: r,
	}
	return r
}

type defaultRouter struct {
	*defaultRouterGroup
}

func (r *defaultRouter) RouteNotFound(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (r *defaultRouter) Handle(c Context) error {
	return c.NotFound()

	//internal lookup
	//retreive route info
	//retreive the path params extracted from

}

func (r *defaultRouter) Lookup(method, path string) (HandlerFunc, RouteInfo) {
	//TODO implement me
	panic("implement me")
}

type defaultRouterGroup struct {
	prefix     string
	middleware []MiddlewareFunc
	router     *defaultRouter
}

func (g *defaultRouterGroup) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *defaultRouterGroup) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodConnect, path, h, m...)
}

func (g *defaultRouterGroup) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodDelete, path, h, m...)
}

func (g *defaultRouterGroup) GET(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodGet, path, h, m...)
}

func (g *defaultRouterGroup) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodHead, path, h, m...)
}

func (g *defaultRouterGroup) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodOptions, path, h, m...)
}

func (g *defaultRouterGroup) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodPatch, path, h, m...)
}

func (g *defaultRouterGroup) POST(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodPost, path, h, m...)
}

func (g *defaultRouterGroup) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodPut, path, h, m...)
}

func (g *defaultRouterGroup) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodTrace, path, h, m...)
}

func (g *defaultRouterGroup) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) []RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (g *defaultRouterGroup) Match(methods []string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) []RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (g *defaultRouterGroup) Group(prefix string, middleware ...MiddlewareFunc) RouteGroup {
	m := make([]MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)

	return &defaultRouterGroup{
		prefix:     g.prefix + prefix,
		middleware: m,
		router:     g.router,
	}
}

func (g *defaultRouterGroup) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) RouteInfo {
	m := make([]MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)

	return newRouteInfo(method, g.prefix+path, handler)
}

func newRouteInfo(method, path string, handler HandlerFunc) RouteInfo {
	template, params := parsePath(path)
	return &routeInfo{
		name:     handlerName(handler),
		method:   method,
		path:     path,
		template: template,
		params:   params,
	}
}

type routeInfo struct {
	name     string
	method   string
	path     string
	template string
	params   []string
}

func (r *routeInfo) Name() string {
	return r.name
}

func (r *routeInfo) Method() string {
	return r.method
}

func (r *routeInfo) Path() string {
	return r.path
}

func (r *routeInfo) Params() []string {
	return r.params
}

func (r *routeInfo) Reverse(params ...interface{}) string {
	return fmt.Sprintf(r.template, params...)
}

func parsePath(path string) (string, []string) {
	// TODO implement me
	return "", []string{}
}

func handlerName(h HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}
