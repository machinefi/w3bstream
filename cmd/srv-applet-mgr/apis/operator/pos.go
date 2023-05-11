package operator

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
)

type CreateOperator struct {
	httpx.MethodPost
	operator.CreateReq `in:"body"`
}

func (r *CreateOperator) Path() string { return "/" }

func (r *CreateOperator) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	r.CreateReq.AccountID = ca.AccountID

	return operator.Create(ctx, &r.CreateReq)
}
