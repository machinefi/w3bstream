package strategy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/types"
)

type GetStrategy struct {
	httpx.MethodGet
	StrategyID types.SFID `in:"path" name:"strategyID"`
}

func (r *GetStrategy) Path() string {
	return "/data/:strategyID"
}

func (r *GetStrategy) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithStrategyBySFID(ctx, r.StrategyID)
	if err != nil {
		return nil, err
	}

	return types.MustStrategyFromContext(ctx), nil
}

type ListStrategy struct {
	httpx.MethodGet
	strategy.ListReq
}

func (r *ListStrategy) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	return strategy.List(ctx, &r.ListReq)
}
