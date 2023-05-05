package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	basetypes "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateInstance struct {
	httpx.MethodPost
	AppletID         basetypes.SFID `in:"path" name:"appletID"`
	deploy.CreateReq `in:"body"`
}

func (r *CreateInstance) Path() string {
	return "/applet/:appletID"
}

func (r *CreateInstance) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}

	if _, ok := types.InstanceFromContext(ctx); ok {
		return nil, status.MultiInstanceDeployed
	}
	return deploy.Upsert(ctx, &r.CreateReq, enums.INSTANCE_STATE__STARTED)
}
