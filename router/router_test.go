package router

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"webappv2"
	"webappv2/container"
)

type testContext struct {
	params []string
}

func (t *testContext) SetParamValue(name, values string) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) ParamValuesPtr() *[]string {
	return &t.params
}

func (t *testContext) Request() *http.Request {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) SetRequest(r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Response() webappv2.Response {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) SetResponse(writer http.ResponseWriter) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Routes() webappv2.Routes {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) CurrentRoute() webappv2.RouteInfo {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Method() string {
	return http.MethodGet
}

func (t *testContext) RealIP() string {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Path() string {
	return "/test"
}

func (t *testContext) Param(name string) string {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) ParamNames() []string {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) SetParamNames(names ...string) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) ParamValues() []string {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) SetParamValues(values ...string) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) QueryParam(name string) string {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) QueryParams() url.Values {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) QueryString() string {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) FormValue(name string) string {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) FormParams() (url.Values, error) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) FormFile(name string) (*multipart.FileHeader, error) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) MultipartForm() (*multipart.Form, error) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Cookie(name string) (*http.Cookie, error) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) SetCookie(cookie *http.Cookie) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Cookies() []*http.Cookie {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Get(key string) interface{} {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Set(key string, val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Bind(i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Validate(i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Render(code int, name string, data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) HTML(code int, html string) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) HTMLBlob(code int, b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) String(code int, s string) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) JSON(code int, i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) JSONPretty(code int, i interface{}, indent string) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) JSONBlob(code int, b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) JSONP(code int, callback string, i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) JSONPBlob(code int, callback string, b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) XML(code int, i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) XMLPretty(code int, i interface{}, indent string) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) XMLBlob(code int, b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Blob(code int, contentType string, b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Stream(code int, contentType string, r io.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) File(file string) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Attachment(file string, name string) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Inline(file string, name string) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) NoContent() error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) NotFound() error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Created(url string) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Redirect(code int, url string) error {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Error(err error) {
	//TODO implement me
	panic("implement me")
}

func (t *testContext) Container() container.Container {
	//TODO implement me
	panic("implement me")
}

func BenchmarkRouterHandleNoAllocations(b *testing.B) {

	r := New()

	r.GET("/intention/{id}:activate", func(c webappv2.Context) error { return nil })

	//req, _ := http.NewRequest(http.MethodGet, "/test/intention/1234:activate", nil)
	//rw := httptest.NewRecorder()
	c := &testContext{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			r.Handle(c)
		}
	}
}
