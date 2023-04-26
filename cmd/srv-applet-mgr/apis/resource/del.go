package resource

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
)

type RemoveResource struct {
	httpx.MethodDelete
	ResourceID                   types.SFID `in:"path" name:"resourceID"`
	resource.RemoveByAppletIDReq `in:"body"`
}

func (r *RemoveResource) Path() string { return "/:resourceID" }

func (r *RemoveResource) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	return nil, resource.RemoveByResIDAndAccIDAndAppID(ctx, r.ResourceID, ca.AccountID, r.RemoveByAppletIDReq.AppletID)
}

type BatchRemoveResource struct {
	httpx.MethodDelete
	resource.CondArgs
}

func (r *BatchRemoveResource) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)

	r.AccountID = ca.AccountID
	return nil, resource.Remove(ctx, &r.CondArgs)
}
