package strategy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
)

type UpdateStrategy struct {
	httpx.MethodPut
	StrategyID         types.SFID `in:"path" name:"strategyID"`
	strategy.UpdateReq `in:"body"`
}

func (r *UpdateStrategy) Path() string {
	return "/:strategyID"
}

func (r *UpdateStrategy) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithStrategyBySFID(ctx, r.StrategyID)
	if err != nil {
		return nil, err
	}
	return nil, strategy.Update(ctx, r.StrategyID, &r.UpdateReq)
}
