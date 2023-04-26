package event

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
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

// 	_receiveEventMtc.WithLabelValues(projectName, publisherMtc).Inc()

var Handler = func(ctx context.Context, ch string, ev *eventpb.Event) (interface{}, error) {
	return OnEventReceived(ctx, ch, ev)
}

func OnEventReceived(ctx context.Context, projectName string, r *eventpb.Event) (ret *HandleEventResult, err error) {
	return nil, nil
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
