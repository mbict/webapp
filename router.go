package webapp

import (
	"net/http"
)

type Router interface {
	// RouteNotFound registers a special-case route which is executed when no other route is found (i.e. HTTP 404 cases)
	// for current request URL.
	// Path supports static and named/any parameters just like other http method is defined. Generally path is ended with
	// wildcard/match-any character (`/*`, `/download/*` etc).
	//
	// Example: `e.RouteNotFound("/*", func(c webapp.Context) error { return c.NoContent(http.StatusNotFound) })`
	RouteNotFound(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	Lookup(method, path string) (HandlerFunc, RouteInfo)

	Handler

	RouteGroup
}

type RouteGroup interface {
	// Use adds middleware to the chain.
	// Important!: only routes registered after the use command will have the middleware, previous registered
	// routes are untouched
	Use(middleware ...MiddlewareFunc)

	// CONNECT registers a new CONNECT route for a path with matching handler in the
	// router with optional route-level middleware.
	CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	// DELETE registers a new DELETE route for a path with matching handler in the router
	// with optional route-level middleware.
	DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	// GET registers a new GET route for a path with matching handler in the router
	// with optional route-level middleware.
	GET(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	// HEAD registers a new HEAD route for a path with matching handler in the
	// router with optional route-level middleware.
	HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	// OPTIONS registers a new OPTIONS route for a path with matching handler in the
	// router with optional route-level middleware.
	OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	// PATCH registers a new PATCH route for a path with matching handler in the
	// router with optional route-level middleware.
	PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	// POST registers a new POST route for a path with matching handler in the
	// router with optional route-level middleware.
	POST(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	// PUT registers a new PUT route for a path with matching handler in the
	// router with optional route-level middleware.
	PUT(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	// TRACE registers a new TRACE route for a path with matching handler in the
	// router with optional route-level middleware.
	TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo

	// Any registers a new route for all HTTP methods and path with matching handler
	// in the router with optional route-level middleware.
	Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) []RouteInfo

	// Match registers a new route for multiple HTTP methods and path with matching
	// handler in the router with optional route-level middleware.
	Match(methods []string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) []RouteInfo

	// Group creates a new router group with prefix and optional group-level middleware.
	Group(prefix string, m ...MiddlewareFunc) RouteGroup

	// Add registers a new route for an HTTP method and path with matching handler
	// in the router with optional route-level middleware.
	Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) RouteInfo
}

// routeInfoGroup is a wrapper for the router route group to manage the returned routes
type routeInfoGroup struct {
	routes routes
	RouteGroup
}

func (g *routeInfoGroup) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodConnect, path, h, m...)
}

func (g *routeInfoGroup) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodDelete, path, h, m...)
}

func (g *routeInfoGroup) GET(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodGet, path, h, m...)
}

func (g *routeInfoGroup) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodHead, path, h, m...)
}

func (g *routeInfoGroup) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodOptions, path, h, m...)
}

func (g *routeInfoGroup) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodPatch, path, h, m...)
}

func (g *routeInfoGroup) POST(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodPost, path, h, m...)
}

func (g *routeInfoGroup) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodPut, path, h, m...)
}

func (g *routeInfoGroup) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	return g.Add(http.MethodTrace, path, h, m...)
}

func (g *routeInfoGroup) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) []RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (g *routeInfoGroup) Match(methods []string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) []RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (g *routeInfoGroup) Group(prefix string, middleware ...MiddlewareFunc) RouteGroup {
	rg := g.RouteGroup.Group(prefix, middleware...)
	return &routeInfoGroup{
		routes:     g.routes,
		RouteGroup: rg,
	}
}

func (g *routeInfoGroup) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) RouteInfo {
	routeInfo := g.RouteGroup.Add(method, path, handler, middleware...)
	g.routes[routeInfo.Name()] = routeInfo

	////keep track of the amount of params is the biggest route, used for context optimisation
	//numParams := len(routeInfo.Params())
	//if numParams > g.maxParams {
	//	g.maxParams = numParams
	//}

	return routeInfo
}
