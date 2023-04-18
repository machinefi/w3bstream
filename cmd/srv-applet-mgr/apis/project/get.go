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

func (r *GetProject) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return types.MustProjectFromContext(ctx), nil
}

type GetProjectDetail struct {
	httpx.MethodGet
}

func (r *GetProjectDetail) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.CurrentAccountFromContext(ctx).
		WithProjectContextByName(ctx, middleware.ProjectProviderFromContext(ctx))
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)

	return project.GetDetail(ctx, prj)
}

type ListProject struct {
	httpx.MethodGet
	project.ListReq
}

func (r *ListProject) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	return project.List(ctx, ca.AccountID, &r.ListReq)
}

type ListProjectDetail struct {
	httpx.MethodGet
	project.ListReq
}

func (r *ListProjectDetail) Path() string { return "/detail" }

func (r *ListProjectDetail) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	return project.ListDetail(ctx, ca.AccountID, &r.ListReq)
}
