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

func (r *GetStrategy) Path() string { return "/data/:strategyID" }

func (r *GetStrategy) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithStrategyBySFID(ctx, r.StrategyID)
	if err != nil {
		return nil, err
	}

	return types.MustStrategyFromContext(ctx), nil
}

type ListStrategy struct {
	httpx.MethodGet
	strategy.ListReq
}

func (r *ListStrategy) Path() string { return "/data_list" }

func (r *ListStrategy) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return strategy.List(ctx, &r.ListReq)
}

type ListStrategyDetail struct {
	httpx.MethodGet
	strategy.ListReq
}

func (r *ListStrategyDetail) Path() string { return "/detail_list" }

func (r *ListStrategyDetail) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return strategy.ListDetail(ctx, &r.ListReq)
}
