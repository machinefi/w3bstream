package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ControlInstance struct {
	httpx.MethodPut
	InstanceID types.SFID      `in:"path" name:"instanceID"`
	Cmd        enums.DeployCmd `in:"path" name:"cmd"`
}

func (r *ControlInstance) Path() string { return "/:instanceID/:cmd" }

func (r *ControlInstance) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)

	ctx, err := ca.WithInstanceContextBySFID(ctx, r.InstanceID)
	if err != nil {
		return nil, err
	}
	state, err := deploy.Deploy(ctx, r.Cmd)
	if err != nil {
		return nil, err
	}
	ins := types.MustInstanceFromContext(ctx)
	ins.State = state
	return ins, nil
}
