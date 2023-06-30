package webapp

import (
	"sort"
)

type RouteInfo interface {
	// Name of the route
	Name() string

	// Method this route is registered for
	Method() string

	// Path returns the path template
	Path() string

	// Params returns all the registered param names in the url
	Params() []string

	// Reverse generates a URL from route and provided parameters.
	Reverse(params ...interface{}) string
}

type Routes interface {
	// URI generates a URI from handler.
	URI(handler HandlerFunc, params ...interface{}) string

	// Reverse generates a URL from route name and provided parameters.
	Reverse(name string, params ...interface{}) string

	// Get returns the registered routes by name
	Get(name string) RouteInfo

	// Get returns all the registered routes
	All() []RouteInfo
}

type routes map[string]RouteInfo

func (r routes) URI(handler HandlerFunc, params ...interface{}) string {
	//TODO implement me
	panic("implement me")
}

func (r routes) Reverse(name string, params ...interface{}) string {
	if routeInfo, ok := r[name]; ok {
		return routeInfo.Reverse(params...)
	}
	return ""
}

func (r routes) Get(name string) RouteInfo {
	if routeInfo, ok := r[name]; ok {
		return routeInfo
	}
	return nil
}

func (r routes) All() []RouteInfo {
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ri := make([]RouteInfo, len(r))
	i := 0
	for _, key := range keys {
		ri[i] = r[key]
		i++
	}
	return ri
}
