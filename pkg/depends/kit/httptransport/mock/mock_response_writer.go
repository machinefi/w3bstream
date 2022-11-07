package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

func NewMockResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		header: http.Header{},
	}
}

type ResponseWriter struct {
	header     http.Header
	StatusCode int
	bytes.Buffer
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

func (w *ResponseWriter) Header() http.Header {
	if w.StatusCode == 0 {
		return w.header
	}

	header := http.Header{}

	for k, v := range w.header {
		header[k] = v
	}

	return header
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func (w *ResponseWriter) Response() *http.Response {
	resp := &http.Response{}
	resp.Header = w.header
	resp.StatusCode = w.StatusCode
	resp.Body = ioutil.NopCloser(&w.Buffer)
	return resp
}

func (w *ResponseWriter) MustDumpResponse() []byte {
	data, err := httputil.DumpResponse(w.Response(), true)
	if err != nil {
		panic(err)
	}
	return bytes.Replace(data, []byte("\r\n"), []byte("\n"), -1)
}
