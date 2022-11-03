package httpx

import (
	"context"
	"net/http"
	"net/url"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

type WithContentType interface {
	ContentType() string
}

type WithStatusCode interface {
	StatusCode() int
}

type WithCookies interface {
	Cookies() []*http.Cookie
}

type RedirectDescriber interface {
	WithStatusCode
	Location() *url.URL
}

type WithHeader interface {
	Header() http.Header
}

type ckStatusCode struct{}

func ContextWithStatusCode(ctx context.Context, code int) context.Context {
	return contextx.WithValue(ctx, ckStatusCode{}, code)
}

func StatusCodeFromContext(ctx context.Context) int {
	if code, ok := ctx.Value(ckStatusCode{}).(int); ok {
		return code
	}
	return http.StatusOK
}
