package webappv2

import (
	stdContext "context"
	"net"
	"net/http"
	"sync"
)

type WebApp interface {
	Service
	Server

	// Pre adds middleware to the chain which is run before router.
	Pre(middleware ...MiddlewareFunc)

	RouteGroup
}

func New(options ...Option) WebApp {
	app := &webapp{}
	app.contextPool.New = func() any {
		return app.newContext()
	}
	app.router = NewDefaultRouter()
	app.routes = make(routes)
	app.jsonEncoder = DefaultJSONEncoder{}
	app.errorHandler = DefaultErrorHandler

	// Apply options that overwrite the default behaviour of the webapp
	for _, option := range options {
		option(app)
	}

	// Init after options
	app.routeInfoGroup = &routeInfoGroup{
		routes:     app.routes,
		RouteGroup: app.router.Group(""),
		webapp:     app,
	}
	app.handler = app.router.Handle

	return app
}

type webapp struct {
	preMiddleware []MiddlewareFunc
	errorHandler  ErrorHandler

	contextPool sync.Pool
	maxParams   int

	router  Router
	routes  routes
	handler HandlerFunc

	jsonEncoder JSONEncoding

	*routeInfoGroup
}

func (a *webapp) Handle(c Context) error {
	return a.handler(c)
}

func (a *webapp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Acquire/Allocate a context from the pool
	c := a.contextPool.Get().(*context)
	c.reset(r, w)

	// Execute chain
	if err := a.handler(c); err != nil {
		a.errorHandler(c, err)
	}

	// Release context back to the pool
	a.contextPool.Put(c)
}

func (a *webapp) Start(address string) error {
	server := new(http.Server)
	server.Handler = a
	server.Addr = address

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	//TODO add close and listener to webapp for gracfull closing
	//TODO examine the tcp handling

	return server.Serve(l)
}

func (a *webapp) StartServer(s *http.Server) (err error) {
	//TODO implement me
	panic("implement me")
}

func (a *webapp) Close() error {
	//TODO implement me
	panic("implement me")
}

func (a *webapp) Shutdown(ctx stdContext.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a *webapp) Binder() Binder {
	//TODO implement me
	panic("implement me")
}

func (a *webapp) Validator() Validator {
	//TODO implement me
	panic("implement me")
}

func (a *webapp) Router() Router {
	return a.router
}

func (a *webapp) Routes() Routes {
	return a.routes
}

func (a *webapp) Renderer() Renderer {
	//TODO implement me
	panic("implement me")
}

func (a *webapp) Logger() Logger {
	//TODO implement me
	panic("implement me")
}

func (a *webapp) JsonEncoder() JSONEncoding {
	return a.jsonEncoder
}

func (a *webapp) XmlEncoder() XMLEncoding {
	//TODO implement me
	panic("implement me")
}

func (a *webapp) Pre(middleware ...MiddlewareFunc) {
	a.preMiddleware = append(a.preMiddleware, middleware...)
	a.handler = applyMiddleware(a.router.Handle, a.preMiddleware...)
}

func (a *webapp) RouteNotFound(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo {
	//TODO implement me
	panic("implement me")
}

// newContext is the factory method for creating a new one
func (w *webapp) newContext() Context {
	return &context{
		request:     nil,
		response:    newResponse(nil),
		paramNames:  nil,
		paramValues: make([]string, w.maxParams),
		webapp:      w,
	}
}
