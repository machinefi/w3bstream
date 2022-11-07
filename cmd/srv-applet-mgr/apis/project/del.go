package project

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/project"
)

type RemoveProject struct {
	httpx.MethodDelete
	ProjectName string `in:"path" name:"projectName"`
}

func (r *RemoveProject) Path() string { return "/:projectName" }

func (r *RemoveProject) Output(ctx context.Context) (interface{}, error) {
	a := middleware.CurrentAccountFromContext(ctx)
	if m, err := a.ValidateProjectPermByPrjName(ctx, r.ProjectName); err != nil {
		return nil, err
	} else {
		return nil, project.RemoveProjectByProjectID(ctx, m.ProjectID)
	}
}
