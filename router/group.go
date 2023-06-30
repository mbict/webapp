package router

import (
	"github.com/mbict/webapp"
	"net/http"
)

type group struct {
	prefix     string
	middleware []webapp.MiddlewareFunc
	router     *Router
}

func (g *group) Use(middleware ...webapp.MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *group) CONNECT(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	return g.Add(http.MethodConnect, path, h, m...)
}

func (g *group) DELETE(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	return g.Add(http.MethodDelete, path, h, m...)
}

func (g *group) GET(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	return g.Add(http.MethodGet, path, h, m...)
}

func (g *group) HEAD(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	return g.Add(http.MethodHead, path, h, m...)
}

func (g *group) OPTIONS(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	return g.Add(http.MethodOptions, path, h, m...)
}

func (g *group) PATCH(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	return g.Add(http.MethodPatch, path, h, m...)
}

func (g *group) POST(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	return g.Add(http.MethodPost, path, h, m...)
}

func (g *group) PUT(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	return g.Add(http.MethodPut, path, h, m...)
}

func (g *group) TRACE(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	return g.Add(http.MethodTrace, path, h, m...)
}

func (g *group) Any(path string, handler webapp.HandlerFunc, middleware ...webapp.MiddlewareFunc) []webapp.RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (g *group) Match(methods []string, path string, handler webapp.HandlerFunc, middleware ...webapp.MiddlewareFunc) []webapp.RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (g *group) Group(prefix string, middleware ...webapp.MiddlewareFunc) webapp.RouteGroup {
	m := make([]webapp.MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)

	return &group{
		prefix:     g.prefix + prefix,
		middleware: m,
		router:     g.router,
	}
}

func (g *group) Add(method, path string, handler webapp.HandlerFunc, middleware ...webapp.MiddlewareFunc) webapp.RouteInfo {
	m := make([]webapp.MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)

	return g.router.add(method, g.prefix+path, handler, m...)
}
