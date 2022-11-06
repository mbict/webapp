package router

import (
	"fmt"
	"reflect"
	"runtime"
	"webappv2"
)

func newRouteInfo(method, path string, handler webappv2.HandlerFunc, middleware ...webappv2.MiddlewareFunc) *routeInfo {
	template, params := parsePath(path)

	//! TODO the params count is going to be part of the preallocation in webapp context

	return &routeInfo{
		name:     handlerName(handler),
		handler:  applyMiddleware(handler, middleware...),
		method:   method,
		path:     path,
		template: template,
		params:   params,
	}
}

type routeInfo struct {
	name     string
	handler  webappv2.HandlerFunc
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

func (r *routeInfo) Handler() webappv2.HandlerFunc {
	return r.handler
}

func parsePath(path string) (string, []string) {
	// TODO implement me getting the param names

	return "", make([]string, countParams(path))
}

func handlerName(h webappv2.HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}

func applyMiddleware(h webappv2.HandlerFunc, middleware ...webappv2.MiddlewareFunc) webappv2.HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
