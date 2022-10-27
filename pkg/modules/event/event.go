package event

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/modules/project"
	"github.com/iotexproject/w3bstream/pkg/modules/publisher"
	"github.com/iotexproject/w3bstream/pkg/modules/strategy"
	"github.com/iotexproject/w3bstream/pkg/modules/vm"
	"github.com/iotexproject/w3bstream/pkg/types"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

type HandleEventRsp struct {
	Results []wasm.EventHandleResult `json:"results"`
}

func OnEventReceived(ctx context.Context, projectName string, r *eventpb.Event) (*HandleEventRsp, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "OnEventReceived")
	defer l.End()

	l = l.WithValues("project_name", projectName)

	eventType := types.EVENTTYPEDEFAULT
	if r.Header != nil && len(r.Header.EventType) > 0 {
		eventType = r.Header.EventType
	}
	l = l.WithValues("event_type", eventType)

	if r.Header != nil && len(r.Header.PubId) > 0 {
		puber, err := publisher.GetPublisherByPublisherKey(ctx, r.Header.PubId)
		if err != nil {
			l.Error(err)
			return nil, err
		}
		l = l.WithValues("publisher", puber.PublisherID)

		prj, err := project.GetProjectByProjectName(ctx, projectName)
		if err != nil {
			l.Error(err)
			return nil, err
		}
		l = l.WithValues("project_id", prj.ProjectID)

		if puber.ProjectID != prj.ProjectID {
			l.Error(errors.New("no project permission"))
			return nil, status.Forbidden.StatusErr().WithDesc("no project permission")
		}
	}

	instances, err := strategy.FindStrategyInstances(ctx, projectName, eventType)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	l.Info("matched strategies: %d", len(instances))

	res := make(chan *wasm.EventHandleResult, len(instances))
	wg := &sync.WaitGroup{}

	for _, v := range instances {
		i := vm.GetConsumer(v.InstanceID)
		if i == nil {
			continue
		}

		wg.Add(1)
		go func() {
			res <- i.HandleEvent(ctx, v.Handler, []byte(r.Payload))
			wg.Done()
		}()
	}

	wg.Wait()
	close(res)

	ret := &HandleEventRsp{}
	for v := range res {
		ret.Results = append(ret.Results, *v)
	}
	return ret, nil
}
