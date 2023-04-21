package event

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

type HandleEvent struct {
	httpx.MethodPost
	ProjectName          string `in:"path" name:"projectName"`
	event.HandleEventReq `in:"body"`
}

func (r *HandleEvent) Path() string { return "/:projectName" }

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	prj, err := project.GetProjectByProjectName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	ctx = types.WithProject(ctx, prj)

	return event.HandleEvents(ctx, r.ProjectName, &r.HandleEventReq), nil
}
