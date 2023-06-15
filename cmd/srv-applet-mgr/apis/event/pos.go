package event

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
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

	if ctx, err = pub.WithProjectContext(ctx); err != nil {
		return nil, err
	}

	prj := types.MustProjectFromContext(ctx)
	metrics.EventMtc.WithLabelValues(prj.AccountID.String(), prj.Name, pub.Key, r.EventType).Inc()

	if err := trafficlimit.TrafficLimit(ctx, enums.TRAFFIC_LIMIT_TYPE__EVENT); err != nil {
		rsp.Results = append([]*event.Result{}, &event.Result{
			AppletName:  "",
			InstanceID:  0,
			Handler:     "",
			ReturnValue: nil,
			ReturnCode:  -1,
			Error:       err.Error(),
		})
		return rsp, nil
	}

	ctx, err = pub.WithStrategiesByChanAndType(ctx, r.Channel, r.EventType)
	if err != nil {
		rsp.Error = statusx.FromErr(err).Key
		return rsp, nil
	}

	ctx = types.WithEventID(ctx, r.EventID)
	rsp.Results = event.OnEvent(ctx, r.Payload.Bytes())
	return rsp, nil
}
