package event

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
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

	if r.IsDataPush() {
		return handleDataPush(ctx, r.Channel, r.Payload.Bytes())
	}

	pub, exist := middleware.PublisherFromContext(ctx)
	if !exist {
		return nil, errors.New("the publisher of the token is not found")
	}

	var (
		err error
		rsp = &event.EventRsp{
			Channel:     r.Channel,
			PublisherID: pub.PublisherID,
			EventID:     r.EventID,
		}
	)

	if ctx, err = pub.WithProjectContext(ctx); err != nil {
		return nil, err
	}

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

	prj := types.MustProjectFromContext(ctx)
	metrics.EventMetricsInc(ctx, prj.AccountID.String(), prj.Name, pub.Key, r.EventType)

	ctx = types.WithEventID(ctx, r.EventID)
	rsp.Results = event.OnEvent(ctx, r.Payload.Bytes())
	return rsp, nil
}

type (
	DataPushReqs []DataPushReq

	DataPushReq struct {
		DeviceID  string `json:"device_id"`
		EventType string `json:"event_type,omitempty"`
		Payload   string `json:"payload"`
		Timestamp int64  `json:"timestamp,omitempty"`
	}

	DataPushRsps []*DataPushRsp
	DataPushRsp  struct {
		Index   int             `json:"index"`
		Results []*event.Result `json:"results"`
	}
)

func handleDataPush(ctx context.Context, ch string, payload []byte) (interface{}, error) {
	var err error
	ca, exist := middleware.CurrentAccountFromContext(ctx)
	if !exist {
		return nil, errors.New("the account of the token is not found")
	}
	ctx = ca.WithAccount(ctx)
	ctx, err = ca.WithProjectContextByName(ctx, ch)
	if err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)
	wrapErr := func(i int, err error) *DataPushRsp {
		return &DataPushRsp{
			Index: i,
			Results: []*event.Result{{
				AppletName:  "",
				InstanceID:  0,
				Handler:     "",
				ReturnValue: nil,
				ReturnCode:  -1,
				Error:       err.Error(),
			}},
		}
	}

	reqs := DataPushReqs{}
	if err := json.Unmarshal(payload, &reqs); err != nil {
		return nil, errors.Wrap(err, "incorrect payload format for batched event")
	}

	rsps := DataPushRsps{}
	for i, v := range reqs {
		pub, err := createPublisherIfNotExist(ctx, prj.ProjectID, v.DeviceID)
		if err != nil {
			rsps = append(rsps, wrapErr(i, err))
			continue
		}
		eventType, eventID := createParamsIfNotExist(v.EventType, "")
		eventResults, err := handleEvent(
			ctx,
			prj,
			pub,
			eventType,
			eventID,
			[]byte(v.Payload),
		)
		if err != nil {
			rsps = append(rsps, wrapErr(i, err))
			continue
		}
		rsps = append(rsps, &DataPushRsp{
			Index:   i,
			Results: eventResults,
		})
	}

	return rsps, nil
}

func createPublisherIfNotExist(ctx context.Context, projectID types.SFID, name string) (*models.Publisher, error) {
	pub, err := publisher.GetByProjectAndKey(ctx, projectID, name)
	if err != nil {
		if err == status.PublisherNotFound {
			pub, err = publisher.Create(ctx, &publisher.CreateReq{
				Name: name,
				Key:  name,
			})
			if err != nil {
				return nil, err
			}
			return pub, nil
		}
		return nil, err
	}
	return pub, nil
}

func createParamsIfNotExist(eventType, eventID string) (string, string) {
	if eventType == "" {
		eventType = enums.EVENTTYPEDEFAULT
	}
	if eventID == "" {
		eventID = uuid.NewString() + "_event2"
	}
	return eventType, eventID
}

func handleEvent(ctx context.Context,
	prj *models.Project,
	pub *models.Publisher,
	eventType string,
	eventID string,
	payload []byte) ([]*event.Result, error) {
	if err := trafficlimit.TrafficLimit(ctx, enums.TRAFFIC_LIMIT_TYPE__EVENT); err != nil {
		return nil, err
	}

	res, err := strategy.FilterByProjectAndEvent(ctx, prj.ProjectID, eventType)
	if err != nil {
		return nil, err
	}
	ctx = types.WithStrategyResults(ctx, res)

	metrics.EventMetricsInc(ctx, prj.AccountID.String(), prj.Name, pub.Key, eventType)

	ctx = types.WithEventID(ctx, eventID)
	return event.OnEvent(ctx, payload), nil
}
