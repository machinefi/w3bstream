package strategy

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/strategy"
)

type CreateStrategy struct {
	httpx.MethodPost
	ProjectName                     string `in:"path" name:"projectName"`
	strategy.CreateStrategyBatchReq `in:"body"`
}

func (r *CreateStrategy) Path() string {
	return "/:projectName"
}

func (r *CreateStrategy) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	if m, err := a.ValidateProjectPermByPrjName(ctx, r.ProjectName); err != nil {
		return nil, err
	} else {
		return nil, strategy.CreateStrategy(ctx, m.ProjectID, &r.CreateStrategyBatchReq)
	}
}
