package http

import (
	"context"
	"net/http"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

func NewContextInjectorMw(ctx contextx.WithContext) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				if ctx != nil {
					req = req.WithContext(ctx(req.Context()))
				}
				next.ServeHTTP(rw, req)
			},
		)
	}
}

type Middleware = func(http.Handler) http.Handler

func WithContextCompose(contexts ...contextx.WithContext) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		for i := range contexts {
			ctx = contexts[i](ctx)
		}
		return ctx
	}
}
