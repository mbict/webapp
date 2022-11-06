package webappv2

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// Response wraps a http.ResponseWriter and implements its interface to be used
// by an HTTP handler to construct an HTTP response.
// See: https://golang.org/pkg/net/http/#ResponseWriter

type Response interface {
	http.ResponseWriter

	SetStatusCode(code int) error
	StatusCode() int
	TotalBytesSent() int64
	HeaderSend() bool
}

type response struct {
	http.ResponseWriter
	status      int
	size        int64
	headersSend bool
}

func (r *response) StatusCode() int {
	return r.status
}

func (r *response) SetStatusCode(code int) error {
	if r.headersSend == true {
		return errors.New("headers already send")
	}
	r.status = code
	return nil
}

func (r *response) TotalBytesSent() int64 {
	return r.size
}

func (r *response) HeaderSend() bool {
	return r.headersSend
}

// newResponse creates a new instance of response.
func newResponse(w http.ResponseWriter) *response {
	return &response{
		ResponseWriter: w,
		size:           0,
		status:         http.StatusOK,
		headersSend:    false,
	}
}

// Header returns the header map for the writer that will be sent by
// WriteHeader. Changing the header after a call to WriteHeader (or Write) has
// no effect unless the modified headers were declared as trailers by setting
// the "Trailer" header before the call to WriteHeader (see example)
// To suppress implicit response headers, set their value to nil.
// Example: https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
func (r *response) Header() http.Header {
	return r.ResponseWriter.Header()
}

// WriteHeader sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus, explicit calls to WriteHeader are mainly
// used to send error codes.
func (r *response) WriteHeader(code int) {
	if r.headersSend {
		//r.webapp.Logger().Warn("headers already send")
		return
	}
	r.status = code
	r.ResponseWriter.WriteHeader(r.status)
	r.headersSend = true
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *response) Write(b []byte) (n int, err error) {
	if !r.headersSend {
		if r.status == 0 {
			r.status = http.StatusOK
		}
		r.WriteHeader(r.status)
	}
	n, err = r.ResponseWriter.Write(b)
	r.size += int64(n)
	return
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *response) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (r *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

func (r *response) reset(w http.ResponseWriter) {
	r.ResponseWriter = w
	r.size = 0
	r.status = http.StatusOK
	r.headersSend = false
}
