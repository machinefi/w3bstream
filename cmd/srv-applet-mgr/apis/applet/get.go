package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ListApplet struct {
	httpx.MethodGet
	applet.ListAppletReq
}

func (r *ListApplet) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}

	return applet.ListApplets(ctx, &r.ListAppletReq)
}

type AppletDetail struct {
	models.Applet
	models.ResourceInfo
	*models.InstanceInfo
}

type GetApplet struct {
	httpx.MethodGet
	AppletID types.SFID `in:"path" name:"appletID"`
}

func (r *GetApplet) Path() string { return "/:appletID" }

func (r *GetApplet) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}

	app := types.MustAppletFromContext(ctx)
	res := types.MustResourceFromContext(ctx)
	ins, _ := types.InstanceFromContext(ctx)

	ret := &AppletDetail{
		Applet:       *app,
		ResourceInfo: res.ResourceInfo,
	}
	if ins != nil {
		ret.InstanceInfo = &ins.InstanceInfo
	}

	return ret, nil
}
