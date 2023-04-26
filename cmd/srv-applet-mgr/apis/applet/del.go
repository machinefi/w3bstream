package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
)

// RemoveApplet remove applet by applet id
type RemoveApplet struct {
	httpx.MethodDelete
	AppletID types.SFID `in:"path" name:"appletID"`
}

func (r *RemoveApplet) Path() string { return "/data/:appletID" }

func (r *RemoveApplet) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}

	return nil, applet.RemoveBySFID(ctx, r.AppletID)
}

// BatchRemoveApplet remove applets with condition under project permission
type BatchRemoveApplet struct {
	httpx.MethodDelete
	applet.CondArgs
}

func (r *BatchRemoveApplet) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	r.ProjectID = types.MustProjectFromContext(ctx).ProjectID
	return nil, applet.Remove(ctx, &r.CondArgs)
}
