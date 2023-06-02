package middleware

import (
	"context"

	confrate "github.com/machinefi/w3bstream/pkg/depends/conf/rate_limit"
)

type EventReqRateLimit struct{}

func (r *EventReqRateLimit) Output(ctx context.Context) (interface{}, error) {
	rl := confrate.MustRateLimitKeyFromContext(ctx)

	rl.Limiter.Take()
	return nil, nil
}
