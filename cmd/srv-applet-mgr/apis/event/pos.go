package event

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/project"
)

type HandleEvent struct {
	httpx.MethodPost
	event.EventCore `in:"body"`
}

func (r *HandleEvent) Path() string { return "/" }

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	prj, err := project.GetProjectInfoByProjectID(ctx, r.ProjectID)
	if err != nil {
		return nil, err
	}
	e := &eventpb.Event{
		Header: &eventpb.Header{
			EventType: r.EventType,
		},
		Payload: r.Payload,
	}

	return event.OnEventReceived(ctx, prj.ProjectName.Name, e)
}
