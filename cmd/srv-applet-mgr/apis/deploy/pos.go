package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	basetypes "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateAndStartInstance struct {
	httpx.MethodPost
	AppletID         basetypes.SFID `in:"path" name:"appletID"`
	deploy.CreateReq `in:"body"`
}

func (r *CreateAndStartInstance) Path() string {
	return "/applet/:appletID"
}

func (r *CreateAndStartInstance) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}
	ctx, err = ca.WithResourceContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}

	if ins, ok := types.InstanceFromContext(ctx); ok {
		if err := deploy.Deploy(ctx, enums.DEPLOY_CMD__START); err != nil {
			return nil, err
		}
		return deploy.GetBySFID(ctx, ins.InstanceID)
	}
	return deploy.Upsert(ctx, &r.CreateReq, enums.INSTANCE_STATE__STARTED)
}
