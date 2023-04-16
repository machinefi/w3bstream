package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveApplet struct {
	httpx.MethodDelete
	AppletID types.SFID `in:"path" name:"appletID"`
}

func (r *RemoveApplet) Path() string { return "/:appletID" }

func (r *RemoveApplet) Output(ctx context.Context) (interface{}, error) {
	var err error
	if ctx, err = middleware.CurrentAccountFromContext(ctx).
		WithAppletContextBySFID(ctx, r.AppletID); err != nil {
		return nil, err
	}
	app := types.MustAppletFromContext(ctx)
	ins, _ := types.InstanceFromContext(ctx) // if not deployed

	return nil, applet.RemoveAppletAndInstanceBySFID(ctx, app, ins)
}
