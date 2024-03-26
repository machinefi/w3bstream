package event

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
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
func HandleEvent(ctx context.Context, tpe string, data []byte) (*EventRsp, error) {
	prj := &models.Project{ProjectName: models.ProjectName{
		Name: types.MustProjectFromContext(ctx).Name,
	}}

	err := prj.FetchByName(types.MustMgrDBExecutorFromContext(ctx))
	if err != nil {
		return nil, err
	}
	ctx = types.WithProject(ctx, prj)

	ctx = types.WithPublisher(ctx, &models.Publisher{
		PrimaryID:    datatypes.PrimaryID{ID: 0},
		RelProject:   models.RelProject{ProjectID: prj.ProjectID},
		RelPublisher: models.RelPublisher{PublisherID: 0},
		PublisherInfo: models.PublisherInfo{
			Key:  "w3b_monitor",
			Name: "w3b_monitor",
		},
	})

	return Create(ctx, &EventReq{
		From:      enums.EVENT_SOURCE__MONITOR,
		Channel:   prj.Name,
		EventType: tpe,
		EventID:   uuid.NewString() + "_monitor",
		Timestamp: time.Now().UTC().UnixMilli(),
		Payload:   *bytes.NewBuffer(data),
	})
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
				// metrics.GeoCollect(ctx, data)
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
	_, l := logr.Start(ctx, "event.Create")
	defer l.End()

	prj := types.MustProjectFromContext(ctx)
	pub := types.MustPublisherFromContext(ctx)

	if err := trafficlimit.TrafficLimit(ctx, prj.ProjectID, enums.TRAFFIC_LIMIT_TYPE__EVENT); err != nil {
		return nil, err
	}

	strategies, err := strategy.FilterByProjectAndEvent(ctx, prj.ProjectID, r.EventType)
	if err != nil {
		return nil, err
	}

	events := make([]*models.Event, 0, len(strategies))
	results := make([]*Result, 0, len(strategies))
	for i, s := range strategies {
		if state, _ := vm.GetInstanceState(s.InstanceID); state != enums.INSTANCE_STATE__STARTED {
			continue
		}
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
				AutoCollect:  s.AutoCollect,
			},
		})
		results = append(results, &Result{
			AppletName: s.AppletName,
			InstanceID: s.InstanceID,
			Handler:    s.Handler,
		})
	}
	if err = models.BatchCreateEvents(types.MustMgrDBExecutorFromContext(ctx), events...); err != nil {
		return nil, status.DatabaseError.StatusErr().
			WithDesc(fmt.Sprintf("batch create event failed: %v", err))
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
	_, l := logr.Start(ctx, "event.BatchCreate")
	defer l.End()

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
			if se, ok := statusx.IsStatusErr(err); ok && se.Key == status.TrafficLimitExceededFailed.Key() {
				break
			}
			return nil, err
		}
		ret = append(ret, &DataPushRsp{
			Index:   i,
			Results: res.Results,
		})
	}
	return ret, nil
}

func handle(ctx context.Context, batch int64, prj types.SFID) int {
	ctx, l := logger.NewSpanContext(ctx, "event.handle")
	defer l.End()

	d := types.MustMgrDBExecutorFromContext(ctx)

	evs, err := models.BatchFetchLastUnhandledEvents(ctx, d, batch, prj)
	if err != nil {
		l.Error(err)
		return 0
	}
	if len(evs) == 0 {
		return 0
	}
	l.WithValues("batch", len(evs)).Info("")

	for _, v := range evs {
		v.HandledAt = time.Now().UnixMilli()
		v.Stage = enums.EVENT_STAGE__HANDLED

		if err = v.UpdateByIDWithFVs(d, builder.FieldValues{
			v.FieldStage():     v.Stage,
			v.FieldHandledAt(): v.HandledAt,
		}); err != nil {
			l.WithValues("evt", v.EventID).Error(err)
			continue
		}

		ins := vm.GetConsumer(v.InstanceID)
		if ins == nil {
			v.CompletedAt = time.Now().UTC().UnixMilli()
			v.ResultCode = -1
			v.Error = status.InstanceNotRunning.Key() + "_nil"
		} else {
			res := ins.HandleEvent(types.WithEventID(ctx, v.EventID), v.Handler, v.EventType, v.Input)
			v.CompletedAt = time.Now().UTC().UnixMilli()
			v.ResultCode = int32(res.Code)
			v.Error = res.ErrMsg
		}
		v.Stage = enums.EVENT_STAGE__COMPLETED
		l := l.WithValues(
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
		go func(ctx context.Context, v *models.Event) {
			metrics.EventMetricsInc(ctx, v)
			if v.AutoCollect == datatypes.TRUE {
				metrics.GeoCollect(ctx, v)
			}
		}(ctx, v)
	}
	return len(evs)
}

func NewDefaultEventHandleScheduler(prj types.SFID) *EventHandleScheduler {
	return NewEventHandleScheduler(time.Second*10, 100, prj)
}

func NewEventHandleScheduler(d time.Duration, batch int64, prj types.SFID) *EventHandleScheduler {
	return &EventHandleScheduler{
		prj:   prj,
		batch: batch,
		du:    d,
	}
}

type EventHandleScheduler struct {
	prj   types.SFID    // prj project sfid
	batch int64         // batch fetch
	du    time.Duration // du interval
}

func (s *EventHandleScheduler) Run(ctx context.Context) {
	for {
		if handled := handle(ctx, s.batch, s.prj); handled == 0 {
			time.Sleep(s.du)
		}
	}
}

func NewDefaultEventCleanupScheduler() *EventCleanupScheduler {
	return NewEventCleanupScheduler(time.Hour, 3*time.Hour*24)
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
