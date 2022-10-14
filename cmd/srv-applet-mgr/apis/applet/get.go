package applet

import (
	"context"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/applet"
)

type ListApplet struct {
	httpx.MethodGet
	ProjectID types.SFID `in:"path" name:"projectID"`
	applet.ListAppletReq
}

func (r *ListApplet) Path() string { return "/:projectID" }

func (r *ListApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	if _, err := ca.ValidateProjectPerm(ctx, r.ProjectID); err != nil {
		return nil, err
	}

	return applet.ListApplets(ctx, &r.ListAppletReq)
}

type GetApplet struct {
	httpx.MethodGet
	applet.GetAppletReq
}

func (r *GetApplet) Path() string { return "/:projectID/:appletID" }

func (r *GetApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	if _, err := ca.ValidateProjectPerm(ctx, r.ProjectID); err != nil {
		return nil, err
	}

	return applet.GetAppletByAppletID(ctx, r.AppletID)
}
