package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
)

type RemoveApplet struct {
	httpx.MethodDelete
	applet.RemoveAppletReq
}

func (r *RemoveApplet) Path() string { return "/:appletID" }

func (r *RemoveApplet) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	app, err := applet.GetAppletByAppletID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}
	if _, err := a.ValidateProjectPerm(ctx, app.ProjectID); err != nil {
		return nil, err
	}

	return nil, applet.RemoveApplet(ctx, &r.RemoveAppletReq)
}
