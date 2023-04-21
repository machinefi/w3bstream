package project

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/project"
)

type CreateProject struct {
	httpx.MethodPost
	project.CreateProjectReq `in:"body"`
}

func (r *CreateProject) Output(ctx context.Context) (interface{}, error) {
	return project.CreateProject(
		ctx, &r.CreateProjectReq,
		func(ctx context.Context, channel string, data *eventpb.Event) (interface{}, error) {
			return event.OnEventReceived(ctx, channel, data)
		},
	)
}
