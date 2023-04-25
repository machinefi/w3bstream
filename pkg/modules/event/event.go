package event

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

var _receiveEventMtc = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "w3b_receive_event_metrics",
	Help: "receive event counter metrics.",
}, []string{"project", "publisher"})

func init() {
	prometheus.MustRegister(_receiveEventMtc)
}

func OnEventReceived(ctx context.Context, projectName string, r *eventpb.Event) (ret *HandleEventResult, err error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "OnEventReceived")
	defer l.End()

	l = l.WithValues("project_name", projectName)

	ret = &HandleEventResult{
		ProjectName: projectName,
	}

	defer func() {
		if err != nil {
			ret.ErrMsg = err.Error()
		}
	}()

	eventType := enums.EVENTTYPEDEFAULT
	eventType, err = checkHeader(ctx, projectName, l, ret, r)
	if err != nil {
		return
	}
	l = l.WithValues("event_type", eventType)

	err = HandleEvent(ctx, projectName, eventType, ret, r.Payload)
	if err != nil {
		return
	}
	return ret, nil
}

func HandleEvent(ctx context.Context, projectName string, eventType string, ret *HandleEventResult, payload []byte) error {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "HandleEvent")
	defer l.End()

	l = l.WithValues("event_type", eventType)
	handlers, err := strategy.FindStrategyInstances(ctx, projectName, eventType)
	if err != nil {
		l.Error(err)
		return err
	}

	l.Info("matched strategies: %d", len(handlers))

	res := make(chan *wasm.EventHandleResult, len(handlers))

	wg := &sync.WaitGroup{}
	for _, v := range handlers {
		i := vm.GetConsumer(v.InstanceID)
		if i == nil {
			res <- &wasm.EventHandleResult{
				InstanceID: v.InstanceID.String(),
				Code:       -1,
				ErrMsg:     "instance not found",
			}
			continue
		}

		wg.Add(1)
		go func(v *strategy.InstanceHandler) {
			defer wg.Done()
			res <- i.HandleEvent(ctx, v.Handler, eventType, payload)
		}(v)
	}
	wg.Wait()
	close(res)

	for v := range res {
		if v == nil {
			continue
		}
		ret.WasmResults = append(ret.WasmResults, v)
	}
	return nil
}

func checkHeader(ctx context.Context, projectName string, l log.Logger, ret *HandleEventResult, r *eventpb.Event) (eventType string, err error) {
	publisherMtc := projectName
	eventType = enums.EVENTTYPEDEFAULT

	if r.Header != nil {
		if len(r.Header.Token) > 0 {
			var (
				pub *models.Publisher
				prj *models.Project
			)
			pub, err = publisherVerification(ctx, l, r)
			if err != nil {
				l.Error(err)
				return
			}
			prj, err = project.GetByName(ctx, projectName)
			if err != nil {
				return
			}
			publisherMtc = pub.Key
			pub, err = publisher.GetByProjectAndKey(ctx, prj.ProjectID, r.Header.GetPubId())
			if err != nil {
				l.Error(err)
				return
			}
			ret.PubID, ret.PubName = pub.PublisherID, pub.Name
			l.WithValues("pub_id", pub.PublisherID)
		}
		if len(r.Header.EventId) > 0 {
			ret.EventID = r.Header.EventId
		}
		if len(r.Header.EventType) > 0 {
			eventType = r.Header.EventType
		}
	}
	_receiveEventMtc.WithLabelValues(projectName, publisherMtc).Inc()
	return
}

func publisherVerification(ctx context.Context, l log.Logger, r *eventpb.Event) (*models.Publisher, error) {
	if r.Header == nil || len(r.Header.Token) == 0 {
		return nil, errors.New("message token is invalid")
	}

	d := types.MustMgrDBExecutorFromContext(ctx)
	publisherJwt := jwt.MustConfFromContext(ctx)

	claim, err := publisherJwt.ParseToken(r.Header.Token)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	v, ok := claim.Payload.(string)
	if !ok {
		l.Error(errors.New("claim of publisher convert string error"))
		return nil, status.InvalidAuthValue
	}
	publisherID := types.SFID(0)
	if err := publisherID.UnmarshalText([]byte(v)); err != nil {
		return nil, status.InvalidAuthPublisherID
	}

	m := &models.Publisher{RelPublisher: models.RelPublisher{PublisherID: publisherID}}
	err = m.FetchByPublisherID(d)
	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "FetchByPublisherID")
	}

	return m, nil
}

func HandleEvents(ctx context.Context, projectName string, r *HandleEventReq) []*HandleEventResult {
	results := make([]*HandleEventResult, 0, len(r.Events))
	for i := range r.Events {
		ret, _ := OnEventReceived(ctx, projectName, &r.Events[i])
		results = append(results, ret)
	}
	return results
}

func OnEvent(ctx context.Context, data []byte) (ret []*wasm.EventHandleResult) {
	l := types.MustLoggerFromContext(ctx)
	r := types.MustStrategyResultsFromContext(ctx)

	// TODO @zhiwei matrix
	results := make(chan *wasm.EventHandleResult, len(r))

	wg := &sync.WaitGroup{}
	for _, v := range r {
		l = l.WithValues(
			"acc", v.AccountID,
			"prj", v.ProjectName,
			"app", v.AppletName,
			"ins", v.InstanceID,
			"hdl", v.Handler,
			"tpe", v.EventType,
		)
		ins := vm.GetConsumer(v.InstanceID)
		if ins == nil {
			l.Warn(errors.New("instance not running"))
			results <- &wasm.EventHandleResult{
				InstanceID: v.InstanceID.String(),
				Code:       -1,
				ErrMsg:     "instance not found",
			}
			continue
		}

		wg.Add(1)
		go func(v *types.StrategyResult) {
			defer wg.Done()

			cost := timer.Start()
			select {
			case <-time.After(time.Second * 5):
			default:
				rv := ins.HandleEvent(ctx, v.Handler, v.EventType, data)
				results <- rv
				l.WithValues("cst", cost().Milliseconds()).Info("")
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
