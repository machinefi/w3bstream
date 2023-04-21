package strategy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
)

type RemoveStrategy struct {
	httpx.MethodDelete
	strategy.RemoveStrategyReq
}

func (r *RemoveStrategy) Path() string { return "/:projectName" }

func (r *RemoveStrategy) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	if _, err := a.ValidateProjectPermByPrjName(ctx, r.ProjectName); err != nil {
		return nil, err
	}

	return nil, strategy.RemoveStrategy(ctx, &r.RemoveStrategyReq)
}
