package resource

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
)

type ListResources struct {
	httpx.MethodGet
	resource.ListReq
}

func (r *ListResources) Path() string { return "/datalist" }

func (r *ListResources) Output(ctx context.Context) (interface{}, error) {
	return resource.ListResourceMeta(ctx, &r.ListReq)
}

type ListResourcesDetail struct {
	httpx.MethodGet
	resource.ListReq
}

func (r *ListResourcesDetail) Path() string { return "/details" }

func (r *ListResourcesDetail) Output(ctx context.Context) (interface{}, error) {
	return resource.ListResourceMetaDetail(ctx, &r.ListReq)
}
