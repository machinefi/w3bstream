package types

import (
	"context"
	"net/http"
)

type HttpRequest struct {
	Method string
	Url    string
	Header http.Header
	Body   []byte
}

type HttpResponse struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	Header     http.Header
	Body       []byte
}

type Server interface {
	Call(ctx context.Context, data []byte) *HttpResponse
}
