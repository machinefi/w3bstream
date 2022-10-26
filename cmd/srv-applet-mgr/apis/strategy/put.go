package strategy

import (
	"context"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/strategy"
)

type UpdateStrategy struct {
	httpx.MethodPut
	ProjectName                string     `in:"path" name:"projectName"`
	StrategyID                 types.SFID `in:"path" name:"strategyID"`
	strategy.CreateStrategyReq `in:"body"`
}

func (r *UpdateStrategy) Path() string {
	return "/:projectName/:strategyID"
}

func (r *UpdateStrategy) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	if _, err := a.ValidateProjectPermByPrjName(ctx, r.ProjectName); err != nil {
		return nil, err
	}

	return nil, strategy.UpdateStrategy(ctx, r.StrategyID, &r.CreateStrategyReq)
}
