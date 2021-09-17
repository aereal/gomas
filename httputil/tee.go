package httputil

import (
	"io"
	"net/http"
)

func NewTeeResponseWriter(rw http.ResponseWriter, buf io.Writer) *TeeResponseWriter {
	return &TeeResponseWriter{rw: rw, mw: io.MultiWriter(rw, buf)}
}

// TeeResponseWriter is an http.ResponseWriter implementation and writes response body into original http.ResponseWriter and given io.Writer.
//
// You can use TeeResponseWriter for capturing sent response body in http.Handler.
type TeeResponseWriter struct {
	rw            http.ResponseWriter
	mw            io.Writer
	statusCode    int
	writtenHeader bool
}

var _ http.ResponseWriter = &TeeResponseWriter{}

func (w *TeeResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *TeeResponseWriter) WriteHeader(statusCode int) {
	w.rw.WriteHeader(statusCode)
	w.writtenHeader = true
}

func (w *TeeResponseWriter) Write(b []byte) (int, error) {
	if !w.writtenHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.mw.Write(b)
}

// StatusCode returns captured response status.
//
// The return value may be zero if no response sent.
func (w *TeeResponseWriter) StatusCode() int {
	return w.statusCode
}
