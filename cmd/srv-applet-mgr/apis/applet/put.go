package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
)

type UpdateApplet struct {
	httpx.MethodPut
	AppletID               types.SFID `in:"path" name:"appletID"`
	applet.UpdateAppletReq `in:"body" mime:"multipart"`
}

func (r *UpdateApplet) Path() string { return "/:appletID" }

func (r *UpdateApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)

	ctx, err := ca.WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}

	return nil, applet.UpdateApplet(ctx, r.AppletID, ca.AccountID, &r.UpdateAppletReq)
}

type UpdateAndDeploy struct {
	httpx.MethodPut
	AppletID                  types.SFID `in:"path" name:"appletID"`
	InstanceID                types.SFID `in:"path" name:"instanceID"`
	applet.UpdateAndDeployReq `in:"body" mime:"multipart"`
}

func (r *UpdateAndDeploy) Path() string {
	return "/:appletID/:instanceID"
}

func (r *UpdateAndDeploy) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}
	ctx, err = ca.WithInstanceContextBySFID(ctx, r.InstanceID)
	if err != nil {
		return nil, err
	}
	return nil, applet.UpdateAndDeploy(ctx, ca.AccountID, &r.UpdateAndDeployReq)
}
