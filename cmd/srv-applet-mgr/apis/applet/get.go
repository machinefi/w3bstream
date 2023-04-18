package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ListApplet struct {
	httpx.MethodGet
	applet.ListAppletReq
}

func (r *ListApplet) Path() string { return "/datalist" }

func (r *ListApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	return applet.ListApplets(ctx, &r.ListAppletReq)
}

type GetApplet struct {
	httpx.MethodGet
	AppletID types.SFID
}

func (r *GetApplet) Path() string { return "/data/:appletID" }

func (r *GetApplet) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}

	return types.MustAppletFromContext(ctx), nil
}
