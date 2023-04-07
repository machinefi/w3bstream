package event

import (
	"context"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
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

func OnEventReceived(ctx context.Context, projectName string, r *eventpb.Event) (ret *HandleEventResult, err error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "OnEventReceived")
	defer l.End()

	l = l.WithValues("project_name", projectName)

	eventType := enums.EVENTTYPEDEFAULT
	if r.Header != nil && len(r.Header.EventType) > 0 {
		eventType = r.Header.EventType
	}
	l = l.WithValues("event_type", eventType)

	ret = &HandleEventResult{
		ProjectName: projectName,
		EventID:     r.Header.EventId,
	}

	defer func() {
		if err != nil {
			ret.ErrMsg = err.Error()
		}
	}()

	var (
		pub      *models.Publisher
		handlers []*strategy.InstanceHandler
	)

	pulisherMtc := projectName
	if r.Header != nil && len(r.Header.PubId) > 0 {
		pulisherMtc = r.Header.PubId
		pub, err = publisher.GetPublisherByPubKeyAndProjectName(ctx, r.Header.PubId, projectName)
		if err != nil {
			return
		}
		ret.PubID, ret.PubName = pub.PublisherID, pub.Name
		l.WithValues("pub_id", pub.PublisherID)
	}
	_receiveEventMtc.WithLabelValues(projectName, pulisherMtc).Inc()

	handlers, err = strategy.FindStrategyInstances(ctx, projectName, eventType)
	if err != nil {
		l.Error(err)
		return
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
			res <- i.HandleEvent(ctx, v.Handler, []byte(r.Payload))
			wg.Done()
		}(v)
	}
	wg.Wait()
	close(res)

	for v := range res {
		ret.WasmResults = append(ret.WasmResults, *v)
	}
	return ret, nil
}

func HandleEvents(ctx context.Context, projectName string, r *HandleEventReq) []*HandleEventResult {
	results := make([]*HandleEventResult, 0, len(r.Events))
	for i := range r.Events {
		ret, _ := OnEventReceived(ctx, projectName, &r.Events[i])
		results = append(results, ret)
	}
	return results
}
