package event

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
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

	ctx, err = pub.WithStrategiesByChanAndType(ctx, r.Channel, r.EventType)
	if err != nil {
		rsp.Error = statusx.FromErr(err).Key
		return rsp, nil
	}

	prj := types.MustProjectFromContext(ctx)
	_eventMtc.WithLabelValues(prj.AccountID.String(), prj.Name, pub.Key, r.EventType).Inc()

	ctx = types.WithEventID(ctx, r.EventID)
	rsp.Results = event.OnEvent(ctx, r.Payload.Bytes())
	return rsp, nil
}
