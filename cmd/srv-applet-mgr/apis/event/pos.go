package event

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/modules/event"
)

type HandleEvent struct {
	httpx.MethodPost
	event.EventReq
}

func (r *HandleEvent) Path() string {
	return "/:channel"
}

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	r.EventReq.SetDefault()

	var (
		err error
		pub = middleware.MustPublisher(ctx)
		rsp = &event.EventRsp{
			Channel:     r.Channel,
			PublisherID: pub.PublisherID,
			EventID:     r.EventID,
		}
	)

	// TODO @zhiwei add event matrix to proxy client transport
	_receiveEventMtc.WithLabelValues(r.Channel, pub.Key).Inc()

	ctx, err = pub.WithStrategiesByChanAndType(ctx, r.Channel, r.EventType)
	if err != nil {
		rsp.Error = statusx.FromErr(err).Key
		return rsp, nil
	}
	rsp.Results = event.OnEvent(ctx, r.Payload)
	return rsp, nil
}
