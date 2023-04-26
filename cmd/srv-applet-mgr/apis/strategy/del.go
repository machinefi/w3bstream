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

func (r *RemoveStrategy) Path() string { return "/data/:strategyID" }

func (r *RemoveStrategy) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithStrategyBySFID(ctx, r.StrategyID)
	if err != nil {
		return nil, err
	}

	return nil, strategy.RemoveBySFID(ctx, r.StrategyID)
}

type BatchRemoveStrategy struct {
	httpx.MethodDelete
	strategy.CondArgs
}

func (r *BatchRemoveStrategy) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	r.ProjectID = types.MustProjectFromContext(ctx).ProjectID
	return nil, strategy.Remove(ctx, &r.CondArgs)
}
