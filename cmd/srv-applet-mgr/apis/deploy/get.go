package deploy

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type GetInstanceByInstanceID struct {
	httpx.MethodGet
	InstanceID types.SFID `in:"path" name:"instanceID"`
}

func (r *GetInstanceByInstanceID) Path() string {
	return "/instance/:instanceID"
}

func (r *GetInstanceByInstanceID) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithInstanceContextBySFID(ctx, r.InstanceID)
	if err != nil {
		return nil, err
	}

	return types.MustInstanceFromContext(ctx), nil
}

type GetInstanceByAppletID struct {
	httpx.MethodGet
	AppletID types.SFID `in:"path" name:"appletID"`
}

func (r *GetInstanceByAppletID) Path() string {
	return "/applet/:appletID"
}

func (r *GetInstanceByAppletID) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithAppletContextBySFID(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}
	return types.MustInstanceFromContext(ctx), nil
}
