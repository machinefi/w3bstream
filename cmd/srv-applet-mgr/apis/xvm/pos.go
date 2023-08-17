package xvm

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/xvm"
)

type CreateRisc0VM struct {
	httpx.MethodPost
	xvm.CreateRisc0VmReq `in:"body" mime:"multipart"`
}

func (r *CreateRisc0VM) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ca.WithAccount(ctx), middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	return xvm.CreateRisc0Vm(ctx, &r.CreateRisc0VmReq)
}
