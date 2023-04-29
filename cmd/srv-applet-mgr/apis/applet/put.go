package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
)

type UpdateApplet struct {
	httpx.MethodPut
	AppletID         types.SFID `in:"path" name:"appletID"`
	applet.UpdateReq `in:"body" mime:"multipart"`
}

func (r *UpdateApplet) Path() string { return "/:appletID" }

func (r *UpdateApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithAppletContextBySFID(ca.WithAccount(ctx), r.AppletID)
	if err != nil {
		return nil, err
	}

	return applet.Update(ctx, &r.UpdateReq)
}
