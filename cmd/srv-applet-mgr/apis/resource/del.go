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
	ResourceID types.SFID `in:"path" name:"resourceID"`
}

func (r *RemoveResource) Path() string { return "/:resourceID" }

func (r *RemoveResource) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithResourceOwnerContextBySFID(ctx, r.ResourceID)
	if err != nil {
		return nil, err
	}
	return nil, resource.RemoveBySFID(ctx, r.ResourceID)
}
