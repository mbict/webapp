package webappv2

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
	"webappv2/container"
)

type Map map[string]interface{}

type Context interface {
	// Request returns `*http.Request`.
	Request() *http.Request

	// SetRequest sets `*http.Request`.
	SetRequest(r *http.Request)

	// Response returns `http.ResponseWriter`.
	Response() Response

	// SetResponse sets the `Response`
	SetResponse(http.ResponseWriter)

	// Routes return the route's helper for generating url
	Routes() Routes

	// CurrentRoute return the current matched route by the router
	CurrentRoute() RouteInfo

	// Method returns the http method used for the request
	Method() string

	// RealIP returns the client's network address based on `X-Forwarded-For`
	// or `X-Real-IP` request header.
	// The behavior can be configured using `Echo#IPExtractor`.
	RealIP() string

	// Path returns the registered path for the handler.
	Path() string

	// Param returns path parameter by name.
	Param(name string) string

	// ParamNames returns path parameter names.
	ParamNames() []string

	// SetParamNames sets path parameter names.
	SetParamNames(names ...string)

	// ParamValues returns path parameter values.
	ParamValues() []string

	// SetParamValues sets path parameter values.
	SetParamValues(values ...string)

	// QueryParam returns the query param for the provided name.
	QueryParam(name string) string

	// QueryParams returns the query parameters as `url.Values`.
	QueryParams() url.Values

	// QueryString returns the URL query string.
	QueryString() string

	// FormValue returns the form field value for the provided name.
	FormValue(name string) string

	// FormParams returns the form parameters as `url.Values`.
	FormParams() (url.Values, error)

	// FormFile returns the multipart form file for the provided name.
	FormFile(name string) (*multipart.FileHeader, error)

	// MultipartForm returns the multipart form.
	MultipartForm() (*multipart.Form, error)

	// Cookie returns the named cookie provided in the request.
	Cookie(name string) (*http.Cookie, error)

	// SetCookie adds a `Set-Cookie` header in HTTP response.
	SetCookie(cookie *http.Cookie)

	// Cookies returns the HTTP cookies sent with the request.
	Cookies() []*http.Cookie

	// Get retrieves data from the context.
	Get(key string) interface{}

	// Set saves data in the context.
	Set(key string, val interface{})

	// Bind binds the request params and body into provided type `i`. The default binder
	// does it based on Content-Type header.
	Bind(i interface{}) error

	// BindBody binds only the request body into provided type `i`. The default binder
	// does it based on Content-Type header.
	BindBody(i interface{}) error

	// BindQueryParams binds the all request params except the body into provided type `i`.
	BindQueryParams(i interface{}) error

	// Validate validates provided `i`. It is usually called after `Context#Bind()`.
	// Validator must be registered using `Echo#Validator`.
	Validate(i interface{}) error

	// Render renders a template with data and sends a text/html response with status
	// code. Renderer must be registered using `Echo.Renderer`.
	Render(code int, name string, data interface{}) error

	// HTML sends an HTTP response with status code.
	HTML(code int, html string) error

	// HTMLBlob sends an HTTP blob response with status code.
	HTMLBlob(code int, b []byte) error

	// String sends a string response with status code.
	String(code int, s string) error

	// JSON sends a JSON response with status code.
	JSON(code int, i interface{}) error

	// JSONPretty sends a pretty-print JSON with status code.
	JSONPretty(code int, i interface{}, indent string) error

	// JSONBlob sends a JSON blob response with status code.
	JSONBlob(code int, b []byte) error

	// JSONP sends a JSONP response with status code. It uses `callback` to construct
	// the JSONP payload.
	JSONP(code int, callback string, i interface{}) error

	// JSONPBlob sends a JSONP blob response with status code. It uses `callback`
	// to construct the JSONP payload.
	JSONPBlob(code int, callback string, b []byte) error

	// XML sends an XML response with status code.
	XML(code int, i interface{}) error

	// XMLPretty sends a pretty-print XML with status code.
	XMLPretty(code int, i interface{}, indent string) error

	// XMLBlob sends an XML blob response with status code.
	XMLBlob(code int, b []byte) error

	// Blob sends a blob response with status code and content type.
	Blob(code int, contentType string, b []byte) error

	// Stream sends a streaming response with status code and content type.
	Stream(code int, contentType string, r io.Reader) error

	// File sends a response with the content of the file.
	File(file string) error

	// Attachment sends a response as attachment, prompting client to save the
	// file.
	Attachment(file string, name string) error

	// Inline sends a response as inline, opening the file in the browser.
	Inline(file string, name string) error

	// NoContent sends a response without a body and a status code.
	NoContent() error

	// NotFound sends a response without body and a 404 code.
	NotFound() error

	// Created sends a response without a body and the url to the location of the created resource in the headers
	Created(url string) error

	// Redirect redirects the request to a provided URL with status code.
	Redirect(code int, url string) error

	// Error invokes the registered HTTP error handler. Generally used by middleware.
	Error(err error)

	// Container returns the DI Container instance.
	// You should rather use the container during initialization of your application and not during handling requests
	Container() container.Container
}

type ParamsContext interface {
	SetParamNames(names ...string)
	SetParamValues(values ...string)
	SetParamValue(name, values string)
	ParamValuesPtr() *[]string
}

type context struct {
	request  *http.Request
	response *response

	currentRoute RouteInfo
	paramNames   []string
	paramValues  []string
	storeLock    sync.RWMutex
	store        map[string]interface{}

	webapp *webapp
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) SetRequest(r *http.Request) {
	c.request = r
}

func (c *context) Response() Response {
	return c.response
}

func (c *context) SetResponse(r http.ResponseWriter) {
	if resp, ok := r.(*response); ok {
		c.response = resp
	}
	c.response.reset(r)
}

func (c *context) Routes() Routes {
	return c.webapp.Routes()
}

func (c *context) CurrentRoute() RouteInfo {
	return c.currentRoute
}

func (c *context) Method() string {
	return c.request.Method
}

func (c *context) RealIP() string {
	//TODO implement me
	panic("implement me")
}

func (c *context) Path() string {
	return c.request.URL.Path
}

func (c *context) Param(name string) string {
	for i, n := range c.paramNames {
		if i < len(c.paramValues) {
			if n == name {
				return c.paramValues[i]
			}
		}
	}
	return ""
}

func (c *context) ParamNames() []string {
	return c.paramNames
}

func (c *context) SetParamNames(names ...string) {
	c.paramNames = names

	//we grow the param values if we try to set more params than allocated
	l := len(c.paramNames)
	if l > len(c.paramValues) {
		//copy values
		newParamValues := make([]string, l)
		copy(newParamValues, c.paramValues)
		c.paramValues = newParamValues
	}
}

func (c *context) ParamValues() []string {
	return c.paramValues[:len(c.paramNames)]
}

func (c *context) SetParamValues(values ...string) {
	//grow param values if more values are provided than allocated,
	//the new values will not be available until the setParamsNames matches the same length
	l := len(values)
	if l > len(c.paramValues) {
		c.paramValues = make([]string, l)
	}

	for i := 0; i < l; i++ {
		c.paramValues[i] = values[i]
	}
}

func (c *context) SetParamValue(name, value string) {
	for i, n := range c.paramNames {
		if n == name {
			c.paramValues[i] = value
			return
		}
	}
}

func (c *context) ParamValuesPtr() *[]string {
	return &c.paramValues
}

func (c *context) QueryParam(name string) string {
	//TODO implement me
	panic("implement me")
}

func (c *context) QueryParams() url.Values {
	//TODO implement me
	panic("implement me")
}

func (c *context) QueryString() string {
	//TODO implement me
	panic("implement me")
}

func (c *context) FormValue(name string) string {
	//TODO implement me
	panic("implement me")
}

func (c *context) FormParams() (url.Values, error) {
	//TODO implement me
	panic("implement me")
}

func (c *context) FormFile(name string) (*multipart.FileHeader, error) {
	//TODO implement me
	panic("implement me")
}

func (c *context) MultipartForm() (*multipart.Form, error) {
	//TODO implement me
	panic("implement me")
}

func (c *context) Cookie(name string) (*http.Cookie, error) {
	//TODO implement me
	panic("implement me")
}

func (c *context) SetCookie(cookie *http.Cookie) {
	//TODO implement me
	panic("implement me")
}

func (c *context) Cookies() []*http.Cookie {
	//TODO implement me
	panic("implement me")
}

func (c *context) Get(key string) interface{} {
	c.storeLock.RLock()
	defer c.storeLock.RUnlock()

	return c.store[key]
}

func (c *context) Set(key string, val interface{}) {
	c.storeLock.Lock()
	defer c.storeLock.Unlock()

	if c.store == nil {
		c.store = make(Map)
	}
	c.store[key] = val
}

func (c *context) Bind(i interface{}) error {
	return c.webapp.binder.Bind(c, i)
}

func (c *context) BindBody(i interface{}) error {
	return c.webapp.binder.BindBody(c, i)
}

func (c *context) BindQueryParams(i interface{}) error {
	return c.webapp.binder.BindQueryParams(c, i)
}

func (c *context) Validate(i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) Render(code int, name string, data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) HTML(code int, html string) error {
	return c.HTMLBlob(code, []byte(html))
}

func (c *context) HTMLBlob(code int, b []byte) error {
	return c.Blob(code, MIMETextHTMLCharsetUTF8, b)
}

func (c *context) String(code int, s string) error {
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(s))
}

func (c *context) JSON(code int, i interface{}) error {
	c.Response().Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.response.SetStatusCode(code)
	return c.webapp.JsonEncoder().Encode(c, i, "")
}

func (c *context) JSONPretty(code int, i interface{}, indent string) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) JSONBlob(code int, b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) JSONP(code int, callback string, i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) JSONPBlob(code int, callback string, b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) XML(code int, i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) XMLPretty(code int, i interface{}, indent string) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) XMLBlob(code int, b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) Blob(code int, contentType string, b []byte) error {
	c.Response().Header().Set(HeaderContentType, contentType)
	c.response.WriteHeader(code)
	_, err := c.response.Write(b)
	return err
}

func (c *context) Stream(code int, contentType string, r io.Reader) error {
	c.Response().Header().Set(HeaderContentType, contentType)
	c.response.WriteHeader(code)
	_, err := io.Copy(c.response, r)
	return err
}

func (c *context) File(file string) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) Attachment(file string, name string) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) Inline(file string, name string) error {
	//TODO implement me
	panic("implement me")
}

func (c *context) NoContent() error {
	c.response.WriteHeader(http.StatusNoContent)
	return nil
}

func (c *context) NotFound() error {
	c.response.WriteHeader(http.StatusNotFound)
	return nil
}

func (c *context) Created(url string) error {
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(http.StatusCreated)
	return nil
}

func (c *context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return ErrInvalidRedirectCode
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}

func (c *context) Error(err error) {
	//TODO implement me
	panic("implement me")
}

func (c *context) Container() container.Container {
	//TODO implement me
	panic("implement me")
}

// reset will clean the context and initialize with the new values
func (c *context) reset(request *http.Request, response http.ResponseWriter) {
	c.request = request
	c.response.reset(response)
	c.store = nil
	c.paramNames = nil
	c.paramValues = c.paramValues[0:c.webapp.maxParams]
	for i := 0; i < c.webapp.maxParams; i++ {
		c.paramValues[i] = ""
	}
}
