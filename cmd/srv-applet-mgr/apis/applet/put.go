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
	ca := middleware.CurrentAccountFromContext(ctx)

	ctx, err := ca.WithAppletContext(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}

	return nil, applet.UpdateApplet(ctx, r.AppletID, &r.UpdateAppletReq)
}
