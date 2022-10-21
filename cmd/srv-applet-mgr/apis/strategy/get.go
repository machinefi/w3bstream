package strategy

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/strategy"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type GetStrategy struct {
	httpx.MethodGet
	ProjectID  types.SFID `in:"path" name:"projectID"`
	StrategyID types.SFID `in:"path" name:"strategyID"`
}

func (r *GetStrategy) Path() string {
	return "/:projectID/:strategyID"
}

func (r *GetStrategy) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	if _, err := a.ValidateProjectPerm(ctx, r.ProjectID); err != nil {
		return nil, err
	}

	return strategy.GetStrategyByStrategyID(ctx, r.StrategyID)
}

type ListStrategy struct {
	httpx.MethodGet
	ProjectID types.SFID `in:"path" name:"projectID"`
	strategy.ListStrategyReq
}

func (r *ListStrategy) Path() string {
	return "/:projectID"
}

func (r *ListStrategy) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	if _, err := a.ValidateProjectPerm(ctx, r.ProjectID); err != nil {
		return nil, err
	}

	r.SetCurrentProjectID(r.ProjectID)
	return strategy.ListStrategy(ctx, &r.ListStrategyReq)
}
