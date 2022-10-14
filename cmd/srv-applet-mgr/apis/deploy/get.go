package deploy

import (
	"context"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/modules/applet"
	"github.com/iotexproject/w3bstream/pkg/modules/deploy"
)

type GetInstanceByInstanceID struct {
	httpx.MethodGet
	InstanceID types.SFID `in:"path" name:"instanceID"`
}

func (r *GetInstanceByInstanceID) Path() string {
	return "/instance/:instanceID"
}

func (r *GetInstanceByInstanceID) Output(ctx context.Context) (interface{}, error) {
	_, err := validateByInstance(ctx, r.InstanceID)
	if err != nil {
		return nil, err
	}

	return deploy.GetInstanceByInstanceID(ctx, r.InstanceID)
}

type GetInstanceByAppletID struct {
	httpx.MethodGet
	AppletID types.SFID `in:"path" name:"appletID"`
}

func (r *GetInstanceByAppletID) Path() string {
	return "/applet/:appletID"
}

func (r *GetInstanceByAppletID) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)

	app, err := applet.GetAppletByAppletID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}

	if _, err = ca.ValidateProjectPerm(ctx, app.ProjectID); err != nil {
		return nil, err
	}

	return deploy.GetInstanceByAppletID(ctx, r.AppletID)
}
