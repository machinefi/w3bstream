package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
)

type ListApplet struct {
	httpx.MethodGet
	ProjectName string `in:"path" name:"projectName"`
	applet.ListAppletReq
}

func (r *ListApplet) Path() string { return "/:projectName" }

func (r *ListApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}

	return applet.ListApplets(ctx, &r.ListAppletReq)
}

type GetApplet struct {
	httpx.MethodGet
	applet.GetAppletReq
}

func (r *GetApplet) Path() string { return "/:projectName/:appletID" }

func (r *GetApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	ctx, err := ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}

	return applet.GetAppletByAppletID(ctx, r.AppletID)
}
