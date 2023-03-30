package middleware

import (
	"context"
	"github.com/pkg/errors"

	confrate "github.com/machinefi/w3bstream/pkg/depends/conf/rate_limit"
)

type ReqRateLimit struct{}

//func (r *ReqRateLimit) ContextKey() string {
//	return ""
//}

func (r *ReqRateLimit) Output(ctx context.Context) (interface{}, error) {
	rl, ok := confrate.RateLimitKeyFromContext(ctx)
	if !ok {
		return nil, errors.New("")
	}

	rl.Limiter.Take()
	return nil, nil
}
