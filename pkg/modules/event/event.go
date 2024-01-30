package event

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
)

// HandleEvent support other module call
// TODO the full project info is not in context so query and set here. this impl
// is for support other module, which is temporary.
// And it will be deprecated when rpc/http is ready
func HandleEvent(ctx context.Context, tpe string, data []byte) (interface{}, error) {
	prj := &models.Project{ProjectName: models.ProjectName{
		Name: types.MustProjectFromContext(ctx).Name,
	}}

	err := prj.FetchByName(types.MustMgrDBExecutorFromContext(ctx))
	if err != nil {
		return nil, err
	}

	eventID := uuid.NewString() + "_monitor"
	ctx = types.WithEventID(ctx, eventID)

	if err := trafficlimit.TrafficLimit(ctx, enums.TRAFFIC_LIMIT_TYPE__EVENT); err != nil {
		results := append([]*Result{}, &Result{
			AppletName:  "",
			InstanceID:  0,
			Handler:     "",
			ReturnValue: nil,
			ReturnCode:  -1,
			Error:       err.Error(),
		})
		return results, nil
	}

	strategies, err := strategy.FilterByProjectAndEvent(ctx, prj.ProjectID, tpe)
	if err != nil {
		return nil, err
	}

	ctx = types.WithStrategyResults(ctx, strategies)

	return OnEvent(ctx, data), nil
}

func OnEvent(ctx context.Context, data []byte) (ret []*Result) {
	ctx, l := logr.Start(ctx, "event.OnEvent", "event_id", types.MustEventIDFromContext(ctx))
	defer l.End()

	var (
		r       = types.MustStrategyResultsFromContext(ctx)
		results = make(chan *Result, len(r))
	)

	wg := &sync.WaitGroup{}
	for _, v := range r {
		l = l.WithValues(
			"prj", v.ProjectName,
			"app", v.AppletName,
		)
		ins := vm.GetConsumer(v.InstanceID)
		if ins == nil {
			l.Warn(errors.New("instance not running"))
			results <- &Result{
				AppletName:  v.AppletName,
				InstanceID:  v.InstanceID,
				Handler:     v.Handler,
				ReturnValue: nil,
				ReturnCode:  -1,
				Error:       status.InstanceNotRunning.Key(),
			}
			continue
		}

		wg.Add(1)
		go func(v *types.StrategyResult) {
			defer wg.Done()
			l.WithValues("ins", v.InstanceID, "hdl", v.Handler, "tpe", v.EventType).Info("handled")
			rv := ins.HandleEvent(ctx, v.Handler, v.EventType, data)
			results <- &Result{
				AppletName:  v.AppletName,
				InstanceID:  v.InstanceID,
				Handler:     v.Handler,
				ReturnValue: nil,
				ReturnCode:  int(rv.Code),
				Error:       rv.ErrMsg,
			}
		}(v)

		go func(v *types.StrategyResult) {
			if v.AutoCollect == datatypes.BooleanValue(true) {
				metrics.GeoCollect(ctx, data)
			}
		}(v)
	}
	wg.Wait()
	close(results)

	for v := range results {
		if v == nil {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

func Create(ctx context.Context, r *EventReq) (*EventRsp, error) {
	prj := types.MustProjectFromContext(ctx)
	pub := types.MustPublisherFromContext(ctx)

	strategies, err := strategy.FilterByProjectAndEvent(ctx, prj.ProjectID, r.EventType)
	if err != nil {
		return nil, err
	}

	events := make([]*models.Event, 0, len(strategies))
	results := make([]*Result, 0, len(strategies))
	for i, s := range strategies {
		events = append(events, &models.Event{
			EventContext: models.EventContext{
				Stage:        enums.EVENT_STAGE__RECEIVED,
				From:         r.From,
				AccountID:    prj.AccountID,
				ProjectID:    prj.ProjectID,
				ProjectName:  prj.Name,
				PublisherID:  pub.PublisherID,
				PublisherKey: pub.Key,
				EventID:      r.EventID,
				Index:        i,
				EventType:    r.EventType,
				InstanceID:   s.InstanceID,
				Handler:      s.Handler,
				Input:        r.Payload.Bytes(),
				Total:        len(strategies),
				PublishedAt:  r.Timestamp,
				ReceivedAt:   time.Now().UTC().UnixMilli(),
			},
		})
		results = append(results, &Result{
			AppletName: s.AppletName,
			InstanceID: s.InstanceID,
			Handler:    s.Handler,
		})
	}
	if err = models.BatchCreateEvents(types.MustMgrDBExecutorFromContext(ctx), events...); err != nil {
		return nil, err
	}
	return &EventRsp{
		Channel:      r.Channel,
		PublisherID:  pub.PublisherID,
		PublisherKey: pub.Key,
		EventID:      r.EventID,
		Timestamp:    time.Now().UTC().UnixMilli(),
		Results:      results,
	}, nil
}

func BatchCreate(ctx context.Context, reqs DataPushReqs) (DataPushRsps, error) {
	ret := make(DataPushRsps, 0, len(reqs))
	prj := types.MustProjectFromContext(ctx)
	for i, v := range reqs {
		pub, err := publisher.CreateIfNotExist(ctx, &publisher.CreateReq{
			Name: v.DeviceID,
			Key:  v.DeviceID,
		})
		if err != nil {
			return nil, err
		}
		ctx = types.WithPublisher(ctx, pub)
		r := &EventReq{
			Channel:   prj.Name,
			EventType: v.EventType,
			EventID:   uuid.NewString() + "_event2",
			Timestamp: v.Timestamp,
			Payload:   *bytes.NewBuffer([]byte(v.Payload)),
		}
		res, err := Create(ctx, r)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &DataPushRsp{
			Index:   i,
			Results: res.Results,
		})
	}
	return ret, nil
}

func handle(ctx context.Context) int {
	ctx, l := logger.NewSpanContext(ctx, "event.handle")
	defer l.End()

	d := types.MustMgrDBExecutorFromContext(ctx)

	evs, err := models.BatchFetchLast100UnhandledEvents(d)
	if err != nil {
		l.Error(err)
		return 0
	}
	l.WithValues("batch", len(evs)).Info("")

	for _, v := range evs {
		v.HandledAt = time.Now().UnixMilli()

		v.Stage = enums.EVENT_STAGE__HANDLED
		ins := vm.GetConsumer(v.InstanceID)
		if ins == nil {
			v.CompletedAt = time.Now().UTC().UnixMilli()
			v.ResultCode = -1
			v.Error = status.InstanceNotRunning.Key()
		} else {
			res := ins.HandleEvent(ctx, v.Handler, v.EventType, v.Input)
			v.CompletedAt = time.Now().UTC().UnixMilli()
			v.ResultCode = int32(res.Code)
			v.Error = res.ErrMsg
		}
		v.Stage = enums.EVENT_STAGE__COMPLETED
		l = l.WithValues(
			"evt", v.EventID,
			"ins", v.InstanceID,
			"hdl", v.Handler,
			"res_code", v.ResultCode,
		)
		if v.Error != "" {
			l.Error(errors.New(v.Error))
		}

		if err = v.UpdateByIDWithFVs(d, builder.FieldValues{
			v.FieldStage():       v.Stage,
			v.FieldHandledAt():   v.HandledAt,
			v.FieldCompletedAt(): v.CompletedAt,
			v.FieldError():       v.Error,
			v.FieldResultCode():  v.ResultCode,
		}); err != nil {
			l.Error(err)
		}
		go metrics.EventMetricsInc(ctx, v.AccountID.String(), v.ProjectName, v.PublisherKey, v.EventType)
	}
	return len(evs)
}

func NewEventHandleScheduler(d time.Duration) *EventHandleScheduler {
	return &EventHandleScheduler{d: d}
}

type EventHandleScheduler struct {
	d time.Duration
}

func (s *EventHandleScheduler) Run(ctx context.Context) {
	for {
		if handled := handle(ctx); handled == 0 {
			time.Sleep(s.d)
		}
	}
}

func NewEventCleanupScheduler(d time.Duration, keep time.Duration) *EventCleanupScheduler {
	return &EventCleanupScheduler{
		d:    d,
		keep: keep,
	}
}

type EventCleanupScheduler struct {
	d    time.Duration
	keep time.Duration
}

func (s *EventCleanupScheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.d)
	m := &models.Event{}
	d := types.MustMgrDBExecutorFromContext(ctx)
	t := d.T(m)
	for {
		_, l := logger.NewSpanContext(ctx, "event.cleanup")

		ts := time.Now().UTC().UnixMilli() - s.keep.Milliseconds()
		_, err := d.Exec(builder.Delete().From(t, builder.Where(m.ColReceivedAt().Lt(ts))))
		if err != nil {
			l.Error(errors.Wrap(err, "event cleanup"))
			l.End()
		} else {
			l.Info("event cleanup")
		}
		l.End()
		<-ticker.C
	}
}
