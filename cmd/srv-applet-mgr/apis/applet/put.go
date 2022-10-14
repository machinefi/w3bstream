package applet

import (
	"context"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/applet"
)

type UpdateApplet struct {
	httpx.MethodPut
	AppletID types.SFID `in:"path" name:"appletID"`
	applet.UpdateAppletReq
}

func (r *UpdateApplet) Path() string { return "/:appletID" }

func (r *UpdateApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)

	app, err := applet.GetAppletByAppletID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}
	if _, err = ca.ValidateProjectPerm(ctx, app.ProjectID); err != nil {
		return nil, err
	}

	return nil, applet.UpdateApplet(ctx, r.AppletID, &r.UpdateAppletReq)
}
