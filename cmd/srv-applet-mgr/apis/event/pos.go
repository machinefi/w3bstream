package event

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/modules/event"
)

type HandleEvent struct {
	httpx.MethodPost
	ProjectName   string `in:"path" name:"projectName"`
	eventpb.Event `in:"body"`
}

func (r *HandleEvent) Path() string { return "/:projectName" }

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	return event.OnEventReceived(ctx, r.ProjectName, &r.Event)
}
