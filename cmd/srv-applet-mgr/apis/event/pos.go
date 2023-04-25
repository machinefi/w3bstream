package event

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/event"
)

type HandleEvent struct {
	httpx.MethodPost
	Channel   string `in:"path"  name:"channel"`
	EventType string `in:"path"  name:"eventType"`
	EventID   string `in:"query" name:"eventID,omitempty"`
	Payload   []byte `in:"body"`
}

func (r *HandleEvent) Path() string {
	return "/:channel/:eventType"
}

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	var (
		err error
		pub = middleware.MustPublisher(ctx)
		rsp = &event.HandleEventResult{
			ProjectName: r.Channel,
			PubID:       pub.PublisherID,
			PubName:     pub.Name,
			EventID:     r.EventID,
		}
	)

	ctx, err = pub.WithStrategiesContextByChannelAndType(ctx, r.Channel, r.EventType)
	if err != nil {
		rsp.ErrMsg = err.Error()
		return rsp, nil
	}
	rsp.WasmResults = event.OnEvent(ctx, r.Payload)
	return rsp, nil
}
