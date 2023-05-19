package middleware

import (
	"context"

	confrate "github.com/machinefi/w3bstream/pkg/depends/conf/rate_limit"
	"github.com/machinefi/w3bstream/pkg/errors/status"
)

type ReqRateLimit struct{}

func (r *ReqRateLimit) Output(ctx context.Context) (interface{}, error) {
	rl, ok := confrate.RateLimitKeyFromContext(ctx)
	if !ok {
		return nil, status.RateLimitKeyNotOk
	}

	rl.Limiter.Take()
	return nil, nil
}
