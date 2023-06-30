package router

import (
	"github.com/mbict/webapp"
	"net/http"
	"strings"
)

// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// Package httprouter is a trie based high performance HTTP request router.
//
// A trivial example is:
//
//  package main
//
//  import (
//      "fmt"
//      "github.com/mbict/httprouter"
//      "net/http"
//      "log"
//  )
//
//  func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//      fmt.Fprint(w, "Welcome!\n")
//  }
//
//  func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//      fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
//  }
//
//  func main() {
//      router := httprouter.New()
//      router.GET("/", Index)
//      router.GET("/hello/:name", Hello)
//
//      log.Fatal(http.ListenAndServe(":8080", router))
//  }
//
// The router matches incoming requests by the request method and the path.
// If a handle is registered for this path and method, the router delegates the
// request to that function.
// For the methods GET, POST, PUT, PATCH, DELETE and OPTIONS shortcut functions exist to
// register handles, for all other methods router.Handle can be used.
//
// The registered path, against which the router matches incoming requests, can
// contain two types of parameters:
//  Syntax    Type
//  :name     named parameter
//  *name     catch-all parameter
//
// Named parameters are dynamic path segments. They match anything until the
// next '/' or the path end:
//  Path: /blog/:category/:post
//
//  Requests:
//   /blog/go/request-routers            match: category="go", post="request-routers"
//   /blog/go/request-routers/           no match, but the router would redirect
//   /blog/go/                           no match
//   /blog/go/request-routers/comments   no match
//
// Catch-all parameters match anything until the path end, including the
// directory index (the '/' before the catch-all). Since they match anything
// until the end, catch-all parameters must always be the final path element.
//  Path: /files/*filepath
//
//  Requests:
//   /files/                             match: filepath="/"
//   /files/LICENSE                      match: filepath="/LICENSE"
//   /files/templates/article.html       match: filepath="/templates/article.html"
//   /files                              no match, but the router would redirect
//
// The value of parameters is saved as a slice of the Param struct, consisting
// each of a key and a value. The slice is passed to the Handle func as a third
// parameter.
// There are two ways to retrieve the value of a parameter:
//  // by the name of the parameter
//  user := ps.ByName("user") // defined by :user or *user
//
//  // by the index of the parameter. This way you can also get the name (key)
//  thirdKey   := ps[2].Key   // the name of the 3rd parameter
//  thirdValue := ps[2].Value // the value of the 3rd parameter

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (path variables).
type Handle = webapp.HandlerFunc

//// Param is a single URL parameter, consisting of a key and a value.
//type Param struct {
//	Key   string
//	Value string
//}
//
//// Params is a Param-slice, as returned by the router.
//// The slice is ordered, the first URL parameter is also the first slice value.
//// It is therefore safe to read values by the index.
//type Params []Param
//
//// ByName returns the value of the first Param which key matches the given name.
//// If no matching Param is found, an empty string is returned.
//func (ps Params) ByName(name string) string {
//	for _, p := range ps {
//		if p.Key == name {
//			return p.Value
//		}
//	}
//	return ""
//}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees map[string]*node

	// If enabled, adds the matched route path onto the http.Request context
	// before invoking the handler.
	// The matched route path is only added to handlers of routes that were
	// registered when this option was enabled.
	SaveMatchedRoutePath bool

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 308 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 308 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// An optional http.Handler that is called on automatic OPTIONS requests.
	// The handler is only called if HandleOPTIONS is true and no OPTIONS
	// handler for the specific path was set.
	// The "Allowed" header is set before calling the handler.
	GlobalOPTIONS http.Handler

	// Cached value of global (*) allowed methods
	globalAllowed string

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound webapp.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed webapp.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})

	// Implements the basic router routing methods
	*group
}

func (r *Router) RouteNotFound(path string, h webapp.HandlerFunc, m ...webapp.MiddlewareFunc) webapp.RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (r *Router) Lookup(method, path string) (webapp.HandlerFunc, webapp.RouteInfo) {
	/*if root := r.trees[method]; root != nil {
		handler, ps, tsr := root.getValue(path, r.getParams)
		if handler == nil {
			r.putParams(ps)
			return nil, nil, tsr
		}
		if ps == nil {
			return handler, nil, tsr
		}
		return handler, *ps, tsr
	}*/
	return nil, nil
}

func (r *Router) Handle(c webapp.Context) error {

	path := c.Path()
	if root := r.trees[c.Method()]; root != nil {

		//get params buffer
		pc := c.(webapp.ParamsContext)
		ps := pc.ParamValuesPtr()
		*ps = (*ps)[0:0] // reset slice

		if routeInfo, _ := root.getValue(path, ps); routeInfo != nil {
			pc.SetParamNames(routeInfo.Params()...)

			return routeInfo.handler(c)

		} /*else if c.Method() != http.MethodConnect && path != "/" {
			// Moved Permanently, request with GET method
			code := http.StatusMovedPermanently
			if c.Method() != http.MethodGet {
				// Permanent Redirect, request with same method
				code = http.StatusPermanentRedirect
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = fixedPath
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}*/
	}

	if c.Method() == http.MethodOptions && r.HandleOPTIONS {
		// Handle OPTIONS requests
		if allow := r.allowed(path, http.MethodOptions); allow != "" {
			c.Response().Header().Set("Allow", allow)

			//todo determine what todo here
			if r.GlobalOPTIONS != nil {
				//r.GlobalOPTIONS.ServeHTTP(w, req)
			}
			return nil
		}
	} else { //Handle method not allowed
		if allow := r.allowed(path, c.Method()); allow != "" {
			c.Response().Header().Set("Allow", allow)
			return webapp.ErrMethodNotAllowed
		}
	}

	// Handle 404
	return r.NotFound.Handle(c)
}

func (r *Router) add(method, path string, handler webapp.HandlerFunc, middleware ...webapp.MiddlewareFunc) webapp.RouteInfo {

	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if handler == nil {
		panic("handler must not be nil")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root

		r.globalAllowed = r.allowed("*", "")
	}

	routeInfo := newRouteInfo(method, path, handler, middleware...)
	root.addRoute(path, routeInfo)

	// Lazy-init paramsPool alloc func
	/* todo move this into the webapp while registering routes, and modify the context creation flow
	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)
			return &ps
		}
	}*/

	return routeInfo
}

// New returns a new initialized Router.
func New(options ...Option) *Router {
	r := &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		NotFound: webapp.HandlerFunc(func(c webapp.Context) error {
			return c.NotFound()
		}),
	}

	r.group = &group{
		router: r,
	}

	for _, option := range options {
		option(r)
	}

	return r
}

//// Handler is an adapter which allows the usage of an http.Handler as a
//// request handle.
//// The Params are available in the request context under ParamsKey.
//func (r *Router) Handler(method, path string, handler http.Handler) {
//	r.Handle(method, path,
//		func(w http.ResponseWriter, req *http.Request, p Params) {
//			if len(p) > 0 {
//				ctx := req.Context()
//				ctx = context.WithValue(ctx, ParamsKey, p)
//				req = req.WithContext(ctx)
//			}
//			handler.ServeHTTP(w, req)
//		},
//	)
//}

//// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
//// request handle.
//func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
//	r.Handler(method, path, handler)
//}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//
//	router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
/*func (r *Router) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)

	r.GET(path, func(w http.ResponseWriter, req *http.Request, ps Params) {
		req.URL.Path = ps.ByName("filepath")
		fileServer.ServeHTTP(w, req)
	})
}*/

func (r *Router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" { // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range r.trees {
				if method == http.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else { // specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == http.MethodOptions {
				continue
			}

			handle, _ := r.trees[method].getValue(path, nil)
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Add request method to list of allowed methods
		allowed = append(allowed, http.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}

	return allow
}
