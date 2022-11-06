package router

import (
	"net/http"
	"webappv2"
)

type group struct {
	prefix     string
	middleware []webappv2.MiddlewareFunc
	router     *Router
}

func (g *group) Use(middleware ...webappv2.MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *group) CONNECT(path string, h webappv2.HandlerFunc, m ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	return g.Add(http.MethodConnect, path, h, m...)
}

func (g *group) DELETE(path string, h webappv2.HandlerFunc, m ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	return g.Add(http.MethodDelete, path, h, m...)
}

func (g *group) GET(path string, h webappv2.HandlerFunc, m ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	return g.Add(http.MethodGet, path, h, m...)
}

func (g *group) HEAD(path string, h webappv2.HandlerFunc, m ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	return g.Add(http.MethodHead, path, h, m...)
}

func (g *group) OPTIONS(path string, h webappv2.HandlerFunc, m ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	return g.Add(http.MethodOptions, path, h, m...)
}

func (g *group) PATCH(path string, h webappv2.HandlerFunc, m ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	return g.Add(http.MethodPatch, path, h, m...)
}

func (g *group) POST(path string, h webappv2.HandlerFunc, m ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	return g.Add(http.MethodPost, path, h, m...)
}

func (g *group) PUT(path string, h webappv2.HandlerFunc, m ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	return g.Add(http.MethodPut, path, h, m...)
}

func (g *group) TRACE(path string, h webappv2.HandlerFunc, m ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	return g.Add(http.MethodTrace, path, h, m...)
}

func (g *group) Any(path string, handler webappv2.HandlerFunc, middleware ...webappv2.MiddlewareFunc) []webappv2.RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (g *group) Match(methods []string, path string, handler webappv2.HandlerFunc, middleware ...webappv2.MiddlewareFunc) []webappv2.RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (g *group) Group(prefix string, middleware ...webappv2.MiddlewareFunc) webappv2.RouteGroup {
	m := make([]webappv2.MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)

	return &group{
		prefix:     g.prefix + prefix,
		middleware: m,
		router:     g.router,
	}
}

func (g *group) Add(method, path string, handler webappv2.HandlerFunc, middleware ...webappv2.MiddlewareFunc) webappv2.RouteInfo {
	m := make([]webappv2.MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)

	return g.router.add(method, g.prefix+path, handler, m...)
}
