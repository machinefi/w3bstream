package strategy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveStrategy struct {
	httpx.MethodDelete
	StrategyID types.SFID `in:"path" name:"strategyID"`
}

func (r *RemoveStrategy) Path() string { return "/:strategyID" }

func (r *RemoveStrategy) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithStrategyBySFID(ctx, r.StrategyID)
	if err != nil {
		return nil, err
	}

	return nil, strategy.Remove(ctx)
}

type BatchRemoveStrategy struct {
	httpx.MethodDelete
	strategy.ListReq
}

func (r *BatchRemoveStrategy) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}

	return nil, strategy.BatchRemove(ctx, &strategy.DataListParam{
		ProjectID: types.MustProjectFromContext(ctx).ProjectID,
		ListReq:   r.ListReq,
	})
}
