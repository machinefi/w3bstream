package traffic_limit

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
)

type CreateTrafficLimit struct {
	httpx.MethodPost
	trafficlimit.CreateReq `in:"body"`
}

func (r *CreateTrafficLimit) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	return trafficlimit.Create(ctx, &r.CreateReq)
}
