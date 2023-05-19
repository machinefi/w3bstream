package ratelimit

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/ratelimit"
	"github.com/machinefi/w3bstream/pkg/types"
)

type UpdateTrafficRateLimit struct {
	httpx.MethodPut
	RateLimitID         types.SFID `in:"path" name:"rateLimitID"`
	ratelimit.UpdateReq `in:"body"`
}

func (r *UpdateTrafficRateLimit) Path() string {
	return "/:rateLimitID"
}

func (r *UpdateTrafficRateLimit) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	r.UpdateReq.RateLimitID = r.RateLimitID
	return ratelimit.Update(ctx, &r.UpdateReq)
}
