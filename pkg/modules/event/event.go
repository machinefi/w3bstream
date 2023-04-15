package event

import (
	"context"
	"sync"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
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

type HandleEventResult struct {
	ProjectName string                   `json:"projectName"`
	PubID       types.SFID               `json:"pubID,omitempty"`
	PubName     string                   `json:"pubName,omitempty"`
	EventID     string                   `json:"eventID"`
	ErrMsg      string                   `json:"errMsg,omitempty"`
	WasmResults []wasm.EventHandleResult `json:"wasmResults"`
}

type HandleEventReq struct {
	Events []eventpb.Event `json:"events"`
}

func OnEventReceived(ctx context.Context, pl []byte) (ret []*wasm.EventHandleResult) {
	l := types.MustLoggerFromContext(ctx)
	r := types.MustStrategyResultsFromContext(ctx)

	results := make(chan *wasm.EventHandleResult, len(r))

	wg := &sync.WaitGroup{}
	for _, v := range r {
		l = l.WithValues(
			"acc_id", v.AccountID,
			"prj_name", v.ProjectName,
			"app_name", v.AppletName,
			"ins_id", v.InstanceID,
			"handler", v.Handler,
			"event_type", v.EventType,
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
				rv := ins.HandleEvent(ctx, v.Handler, pl)
				results <- rv
				l.WithValues("cost_ms", cost().Milliseconds()).Info("")
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

	// _, l = l.Start(ctx, "OnEventReceived")
	// defer l.End()

	// l = l.WithValues("project_name", projectName)

	// ret = &HandleEventResult{
	// 	ProjectName: projectName,
	// }

	// defer func() {
	// 	if err != nil {
	// 		ret.ErrMsg = err.Error()
	// 	}
	// }()

	// eventType := enums.EVENTTYPEDEFAULT
	// publisherMtc := projectName
	// if r.Header != nil {
	// 	if len(r.Header.EventId) > 0 {
	// 		ret.EventID = r.Header.EventId
	// 	}
	// 	if len(r.Header.PubId) > 0 {
	// 		publisherMtc = r.Header.PubId
	// 		var pub *models.Publisher
	// 		pub, err = publisher.GetPublisherByPubKeyAndProjectName(ctx, r.Header.PubId, projectName)
	// 		if err != nil {
	// 			return
	// 		}
	// 		ret.PubID, ret.PubName = pub.PublisherID, pub.Name
	// 		l.WithValues("pub_id", pub.PublisherID)
	// 	}
	// 	if len(r.Header.EventType) > 0 {
	// 		eventType = r.Header.EventType
	// 	}
	// 	if len(r.Header.Token) > 0 {
	// 		if err = publisherVerification(ctx, r, l); err != nil {
	// 			l.Error(err)
	// 			return
	// 		}
	// 	}
	// }
	// _receiveEventMtc.WithLabelValues(projectName, publisherMtc).Inc()

	// l = l.WithValues("event_type", eventType)
	// var handlers []*strategy.InstanceHandler
	// l = l.WithValues("event_type", eventType)
	// handlers, err = strategy.FindStrategyInstances(ctx, projectName, eventType)
	// if err != nil {
	// 	l.Error(err)
	// 	return
	// }

	// l.Info("matched strategies: %d", len(handlers))

	// res := make(chan *wasm.EventHandleResult, len(handlers))

	// wg := &sync.WaitGroup{}
	// for _, v := range handlers {
	// 	i := vm.GetConsumer(v.InstanceID)
	// 	if i == nil {
	// 		res <- &wasm.EventHandleResult{
	// 			InstanceID: v.InstanceID.String(),
	// 			Code:       -1,
	// 			ErrMsg:     "instance not found",
	// 		}
	// 		continue
	// 	}

	// 	wg.Add(1)
	// 	go func(v *strategy.InstanceHandler) {
	// 		defer wg.Done()
	// 		res <- i.HandleEvent(ctx, v.Handler, []byte(r.Payload))
	// 	}(v)
	// }
	// wg.Wait()
	// close(res)

	// for v := range res {
	// 	if v == nil {
	// 		continue
	// 	}
	// 	ret.WasmResults = append(ret.WasmResults, *v)
	// }
	// return ret, nil
}

func publisherVerification(ctx context.Context, r *eventpb.Event, l log.Logger) error {
	if r.Header == nil || len(r.Header.Token) == 0 {
		return errors.New("message token is invalid")
	}

	d := types.MustMgrDBExecutorFromContext(ctx)
	p := types.MustProjectFromContext(ctx)

	publisherJwt := &jwt.Jwt{
		Issuer:  p.ProjectBase.Issuer,
		ExpIn:   p.ProjectBase.ExpIn,
		SignKey: p.ProjectBase.SignKey,
	}
	claim, err := publisherJwt.ParseToken(r.Header.Token)
	if err != nil {
		l.Error(err)
		return err
	}

	v, ok := claim.Payload.(string)
	if !ok {
		l.Error(errors.New("claim of publisher convert string error"))
		return status.InvalidAuthValue
	}
	publisherID := types.SFID(0)
	if err := publisherID.UnmarshalText([]byte(v)); err != nil {
		return status.InvalidAuthPublisherID
	}

	m := &models.Publisher{RelPublisher: models.RelPublisher{PublisherID: publisherID}}
	err = m.FetchByPublisherID(d)
	if err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "FetchByPublisherID")
	}

	if m.ProjectID == p.ProjectID {
		return nil
	} else {
		return status.NoProjectPermission
	}
}
