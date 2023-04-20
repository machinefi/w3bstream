package project

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

type GetProject struct {
	httpx.MethodGet
}

func (r *GetProject) Path() string { return "/data" }

func (r *GetProject) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}
	return types.MustProjectFromContext(ctx), nil
}

type ListProject struct {
	httpx.MethodGet
	project.ListReq
}

func (r *ListProject) Path() string { return "/datalist" }

func (r *ListProject) Output(ctx context.Context) (interface{}, error) {
	ctx = middleware.MustCurrentAccountFromContext(ctx).WithAccount(ctx)
	return project.List(ctx, &r.ListReq)
}
