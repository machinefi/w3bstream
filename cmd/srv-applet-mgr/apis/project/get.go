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
	ProjectName string `in:"path" name:"projectName"`
}

func (r *GetProject) Path() string { return "/:projectName" }

func (r *GetProject) Output(ctx context.Context) (interface{}, error) {
	var (
		err error
		ca  = middleware.CurrentAccountFromContext(ctx)
	)
	ctx, err = ca.WithProjectContextByName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	return types.MustProjectFromContext(ctx), nil
}

type ListProject struct {
	httpx.MethodGet
	project.ListProjectReq
}

func (r *ListProject) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	r.ListProjectReq.SetCurrentAccount(ca.AccountID)
	return project.ListProject(ctx, &r.ListProjectReq)
}
