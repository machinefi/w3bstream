package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
)

type ControlInstance struct {
	httpx.MethodPut
	InstanceID types.SFID      `in:"path" name:"instanceID"`
	Cmd        enums.DeployCmd `in:"path" name:"cmd"`
}

func (r *ControlInstance) Path() string { return "/:instanceID/:cmd" }

func (r *ControlInstance) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithInstanceContextBySFID(ctx, r.InstanceID)
	if err != nil {
		return nil, err
	}

	return nil, deploy.Deploy(ctx, r.Cmd)
}
