package resource

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/types"
)

// DownloadResource download resource by resource id
type DownloadResource struct {
	httpx.MethodGet
	ResourceID types.SFID `in:"path" name:"resourceID"`
}

func (r *DownloadResource) Path() string { return "/data/:resourceID" }

func (r *DownloadResource) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithResourceOwnerContextBySFID(ctx, r.ResourceID)
	if err != nil {
		return nil, err
	}
	ship := types.MustResourceOwnershipFromContext(ctx)

	_, data, err := resource.GetContentBySFID(ctx, r.ResourceID)
	if err != nil {
		return nil, err
	}
	file := httpx.NewAttachment(ship.Filename, "text/plain")
	file.Write(data)

	return file, nil
}

type GetDownloadResourceUrl struct {
	httpx.MethodGet
	ResourceID types.SFID `in:"path" name:"resourceID"`
}

func (r *GetDownloadResourceUrl) Path() string { return "/url/:resourceID" }

func (r *GetDownloadResourceUrl) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithResourceOwnerContextBySFID(ctx, r.ResourceID)
	if err != nil {
		return nil, err
	}

	return resource.GetDownloadUrlBySFID(ctx, r.ResourceID)
}

type ListResources struct {
	httpx.MethodGet
	resource.ListReq
}

func (r *ListResources) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)

	r.AccountID = ca.AccountID
	return resource.List(ctx, &r.ListReq)
}
