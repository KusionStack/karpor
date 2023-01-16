package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

var (
	_ http.ResponseWriter = &Response{}
	_ http.Hijacker       = &Response{}
)

type Response struct {
	http.ResponseWriter
	// statusCode cached http response status code
	statusCode int
	// body cached http response body
	body *bytes.Buffer
	// if underlying ResponseWriter supports it
	hijacker http.Hijacker
}

func NewResponse(responseWriter http.ResponseWriter) *Response {
	hijacker, _ := responseWriter.(http.Hijacker)
	return &Response{ResponseWriter: responseWriter, hijacker: hijacker, body: &bytes.Buffer{}}
}

func (r *Response) Write(p []byte) (int, error) {
	r.body.Write(p)
	return r.ResponseWriter.Write(p)
}

func (r *Response) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if r.hijacker == nil {
		return nil, nil, errors.New("http.Hijacker not implemented by underlying http.ResponseWriter")
	}
	return r.hijacker.Hijack()
}

func WithLogs(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		newWriter := NewResponse(writer)
		userAgent := req.UserAgent()

		// invoke real handler and cached response code and body
		handler.ServeHTTP(newWriter, req)

		basicMsgKV := []interface{}{"userAgent", userAgent, "httpcode", newWriter.statusCode}
		// status code is not 2XX or 3XX
		if newWriter.statusCode < http.StatusOK || newWriter.statusCode >= http.StatusMultipleChoices {
			klog.ErrorS(fmt.Errorf("status code is not 2XX or 3XX"), "Request failed", append(basicMsgKV, "body", newWriter.body.String())...)
		}
	})
}
