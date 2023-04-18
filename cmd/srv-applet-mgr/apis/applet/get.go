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
	applet.ListReq
}

func (r *ListApplet) Path() string { return "/data_list" }

func (r *ListApplet) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)

	return applet.List(ctx, prj.ProjectID, &r.ListReq)
}

type ListAppletDetail struct {
	httpx.MethodGet
	applet.ListReq
}

func (r *ListAppletDetail) Path() string { return "/detail_list" }

func (r *ListAppletDetail) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)

	return applet.ListDetail(ctx, prj.ProjectID, &r.ListReq)
}

type GetApplet struct {
	httpx.MethodGet
	AppletID types.SFID `in:"path" name:"appletID"`
}

func (r *GetApplet) Path() string { return "/data/:appletID" }

func (r *GetApplet) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}
	app := types.MustAppletFromContext(ctx)
	ins, _ := types.InstanceFromContext(ctx)
	res := types.MustResourceFromContext(ctx)

	return applet.GetDetail(ctx, app, ins, res), nil
}

type GetAppletDetail struct {
	httpx.MethodGet
	AppletID types.SFID `in:"path" name:"appletID"`
}

func (r *GetAppletDetail) Path() string { return "/detail/:appletID" }

func (r *GetAppletDetail) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}
	app := types.MustAppletFromContext(ctx)
	ins, _ := types.InstanceFromContext(ctx)
	res := types.MustResourceFromContext(ctx)

	return applet.GetDetail(ctx, app, ins, res), nil
}
