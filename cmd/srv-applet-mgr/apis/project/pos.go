package project

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
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
	ca := middleware.MustCurrentAccountFromContext(ctx)

	prefix, err := middleware.ProjectNameModifier(ctx)
	if err != nil {
		return nil, err
	}
	r.Name = prefix + r.Name

	rsp, err := project.CreateProject(
		ctx, ca.AccountID, &r.CreateProjectReq,
		func(ctx context.Context, channel string, data *eventpb.Event) (interface{}, error) {
			return event.OnEventReceived(ctx, channel, data)
		},
	)
	if err != nil {
		return nil, err
	}
	rsp.Name, _ = middleware.ProjectNameForDisplay(rsp.Name)
	return rsp, nil
}
