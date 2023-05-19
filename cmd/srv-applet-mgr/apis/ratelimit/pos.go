package ratelimit

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/ratelimit"
)

type CreateTrafficRateLimit struct {
	httpx.MethodPost
	ratelimit.CreateReq `in:"body"`
}

func (r *CreateTrafficRateLimit) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	return ratelimit.Create(ctx, &r.CreateReq)
}
