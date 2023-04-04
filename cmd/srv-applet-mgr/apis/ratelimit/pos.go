package ratelimit

import (
	"context"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/ratelimit"
)

type CreateTrafficRateLimit struct {
	httpx.MethodPost
	ProjectName                         string `in:"path" name:"projectName"`
	ratelimit.CreateTrafficRateLimitReq `in:"body"`
}

func (r *CreateTrafficRateLimit) Path() string {
	return "/:projectName"
}

func (r *CreateTrafficRateLimit) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	ctx, err := a.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}

	return ratelimit.CreateRateLimit(ctx, &r.CreateTrafficRateLimitReq)
}
