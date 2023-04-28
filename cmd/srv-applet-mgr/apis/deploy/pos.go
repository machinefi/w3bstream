package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
)

type CreateInstance struct {
	httpx.MethodPost
	AppletID         types.SFID `in:"path" name:"appletID"`
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
	return deploy.Create(ctx, &r.CreateReq)
}
