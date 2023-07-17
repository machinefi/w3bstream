package types

import (
	"context"
	"net/http"
)

const W3bstreamSystemProjectID = "W3bstreamSystemProjectID"

type Server interface {
	Call(ctx context.Context, data []byte) *http.Response
}
