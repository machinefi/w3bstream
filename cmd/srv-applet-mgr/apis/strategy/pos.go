package strategy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
)

type CreateStrategy struct {
	httpx.MethodPost
	strategy.CreateReq `in:"body"`
}

func (r *CreateStrategy) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	return strategy.Create(ctx, &r.CreateReq)
}
