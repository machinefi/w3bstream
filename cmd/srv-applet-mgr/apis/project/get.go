package project

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/project"
)

type GetProject struct {
	httpx.MethodGet
	ProjectName string `in:"path" name:"projectName" validate:"@projectName"`
}

func (r *GetProject) Path() string { return "/:projectName" }

func (r *GetProject) Output(ctx context.Context) (interface{}, error) {
	return project.GetProjectByProjectName(ctx, r.ProjectName)
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
