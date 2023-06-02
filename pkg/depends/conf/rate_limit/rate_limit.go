package rate_limit

import (
	"context"

	"go.uber.org/ratelimit"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

type RateLimit struct {
	Count    int            `env:""`
	Duration types.Duration `env:""`

	Limiter ratelimit.Limiter `env:"-"`
}

func (r *RateLimit) SetDefault() {}

func (r *RateLimit) Init() {
	if r.Count > 0 && r.Duration > 0 {
		// TODO check Duration
		r.Limiter = ratelimit.New(r.Count, ratelimit.Per(r.Duration.Duration()))
	} else {
		r.Limiter = ratelimit.NewUnlimited()
	}
}

type rateLimitKey struct{}

func WithRateLimitKeyContext(rateLimit *RateLimit) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, rateLimitKey{}, rateLimit)
	}
}

func RateLimitKeyFromContext(ctx context.Context) (*RateLimit, bool) {
	j, ok := ctx.Value(rateLimitKey{}).(*RateLimit)
	return j, ok
}

func MustRateLimitKeyFromContext(ctx context.Context) *RateLimit {
	j, ok := ctx.Value(rateLimitKey{}).(*RateLimit)
	must.BeTrue(ok)
	return j
}
