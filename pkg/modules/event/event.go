package event

import (
	"context"
	"time"

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

type HandleEventResult struct {
	InstanceID types.SFID            `json:"instanceID"`
	Error      string                `json:"error"`
	ResultCode wasm.ResultStatusCode `json:"resultCode"`
}

type HandleEventRsp []HandleEventResult

func OnEventReceived(ctx context.Context, projectName string, r *eventpb.Event) (HandleEventRsp, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "OnEventReceived")
	defer l.End()

	l = l.WithValues("project_name", projectName)

	eventType := types.EVENTTYPEDEFAULT
	if r.Header != nil {
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

	ret := make(HandleEventRsp, 0, len(instances))

	for _, v := range instances {
		consumer := vm.GetConsumer(v.InstanceID)
		if consumer == nil {
			continue
		}
		cctx, _ := context.WithTimeout(ctx, 3*time.Second)
		_, code, err := consumer.HandleEvent(cctx, v.Handler, []byte(r.Payload))

		if err != nil {
			ret = append(ret, HandleEventResult{
				InstanceID: v.InstanceID,
				Error:      err.Error(),
				ResultCode: code,
			})
		} else {
			ret = append(ret, HandleEventResult{
				InstanceID: v.InstanceID,
				ResultCode: code,
			})
		}

	}
	return ret, nil
}
