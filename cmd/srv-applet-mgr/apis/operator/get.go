package operator

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
)

type ListOperator struct {
	httpx.MethodGet
	operator.ListReq
}

func (r *ListOperator) Path() string { return "/datalist" }

func (r *ListOperator) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	r.ListReq.AccountID = ca.AccountID

	return operator.ListDetail(ctx, &r.ListReq)
}
