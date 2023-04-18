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
	ProjectName string     `in:"path" name:"projectName"`
	StrategyID  types.SFID `in:"path" name:"strategyID"`
}

func (r *GetStrategy) Path() string {
	return "/:projectName/:strategyID"
}

func (r *GetStrategy) Output(ctx context.Context) (interface{}, error) {
	a := middleware.MustCurrentAccountFromContext(ctx)
	if _, err := a.ValidateProjectPermByPrjName(ctx, r.ProjectName); err != nil {
		return nil, err
	}

	return strategy.GetStrategyByStrategyID(ctx, r.StrategyID)
}

type ListStrategy struct {
	httpx.MethodGet
	ProjectName string `in:"path" name:"projectName"`
	strategy.ListStrategyReq
}

func (r *ListStrategy) Path() string {
	return "/:projectName"
}

func (r *ListStrategy) Output(ctx context.Context) (interface{}, error) {
	a := middleware.MustCurrentAccountFromContext(ctx)
	if m, err := a.ValidateProjectPermByPrjName(ctx, r.ProjectName); err != nil {
		return nil, err
	} else {
		r.SetCurrentProjectID(m.ProjectID)
		return strategy.ListStrategy(ctx, &r.ListStrategyReq)
	}
}
