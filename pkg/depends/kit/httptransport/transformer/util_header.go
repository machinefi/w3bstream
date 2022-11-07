package transformer

import (
	"io"
	"net/http"
	"net/textproto"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
)

func MIMEHeader(headers ...textproto.MIMEHeader) textproto.MIMEHeader {
	header := textproto.MIMEHeader{}
	for _, h := range headers {
		for k, values := range h {
			for _, v := range values {
				header.Add(k, v)
			}
		}
	}
	return header
}

type HeaderWriter interface {
	httpx.WithHeader
	io.Writer
}

func WriterWithHeader(w io.Writer, header http.Header) HeaderWriter {
	return &headerWriter{Writer: w, header: header}
}

func (f *headerWriter) Header() http.Header { return f.header }

type headerWriter struct {
	io.Writer
	header http.Header
}
