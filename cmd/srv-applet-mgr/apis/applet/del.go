package applet

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
)

type RemoveApplet struct {
	httpx.MethodDelete
	AppletID types.SFID `in:"path"  name:"appletID"`
}

func (r *RemoveApplet) Path() string { return "/:appletID" }

func (r *RemoveApplet) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)

	_, err := ca.WithAppletContext(ctx, r.AppletID)
	if err != nil {
		return nil, err
	}
	return nil, applet.RemoveApplet(ctx, r.AppletID)
}
