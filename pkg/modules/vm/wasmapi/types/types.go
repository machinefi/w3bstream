package types

import (
	"context"
	"net/http"
)

type Server interface {
	Call(ctx context.Context, data []byte) *http.Response
}
