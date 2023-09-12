package operator

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveOperator struct {
	httpx.MethodDelete
	OperatorID types.SFID `in:"path" name:"operatorID"`
}

func (r *RemoveOperator) Path() string { return "/data/:operatorID" }

func (r *RemoveOperator) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithOperatorBySFID(ctx, r.OperatorID)
	if err != nil {
		return nil, err
	}
	pool := types.MustOperatorPoolFromContext(ctx)
	return nil, pool.Delete(r.OperatorID)
}
