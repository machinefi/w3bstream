package ratelimit

import (
	"context"
	"github.com/machinefi/w3bstream/pkg/types"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/ratelimit"
)

type UpdateTrafficRateLimit struct {
	httpx.MethodPut
	ProjectName                         string     `in:"path" name:"projectName"`
	RateLimitID                         types.SFID `in:"path" name:"rateLimitID"`
	ratelimit.CreateTrafficRateLimitReq `in:"body"`
}

func (r *UpdateTrafficRateLimit) Path() string {
	return "/:projectName/:rateLimitID"
}

func (r *UpdateTrafficRateLimit) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	ctx, err := a.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}

	return ratelimit.UpdateRateLimit(ctx, r.RateLimitID, &r.CreateTrafficRateLimitReq)
}
