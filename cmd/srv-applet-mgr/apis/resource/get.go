package resource

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/types"
)

type GetResource struct {
	httpx.MethodGet
	ResourceID types.SFID `in:"path" name:"resourceID"`
}

func (r *GetResource) Path() string { return "/:resourceID" }

func (r *GetResource) Output(ctx context.Context) (interface{}, error) {
	return resource.GetBySFID(ctx, r.ResourceID)
}
