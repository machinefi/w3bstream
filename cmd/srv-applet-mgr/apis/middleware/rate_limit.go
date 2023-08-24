package middleware

import (
	"context"

	confrate "github.com/machinefi/w3bstream/pkg/depends/conf/rate_limit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
)

type EventReqRateLimit struct{}

func (r *EventReqRateLimit) Output(ctx context.Context) (interface{}, error) {
	ctx, l := logr.Start(ctx, "middleware.EventReqRateLimit.Output")
	defer l.End()

	rl := confrate.MustRateLimitKeyFromContext(ctx)

	rl.Limiter.Take()
	return nil, nil
}
