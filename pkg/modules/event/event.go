package event

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
)

// HandleEvent support other module call
// TODO the full project info is not in context so query and set here. this impl
// is for support other module, which is temporary.
// And it will be deprecated when rpc/http is ready
func HandleEvent(ctx context.Context, t string, data []byte) (interface{}, error) {
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

	strategies, err := strategy.FilterByProjectAndEvent(ctx, prj.ProjectID, t)
	if err != nil {
		return nil, err
	}

	ctx = types.WithStrategyResults(ctx, strategies)

	return OnEvent(ctx, data), nil
}

func OnEvent(ctx context.Context, data []byte) (ret []*Result) {
	ctx, l := logr.Start(ctx, "modules.event.OnEvent")
	defer l.End()

	var (
		r       = types.MustStrategyResultsFromContext(ctx)
		eventID = types.MustEventIDFromContext(ctx)

		results = make(chan *Result, len(r))
	)

	wg := &sync.WaitGroup{}
	for _, v := range r {
		l = l.WithValues(
			"prj", v.ProjectName,
			"app", v.AppletName,
			"ins", v.InstanceID,
			"hdl", v.Handler,
			"tpe", v.EventType,
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
			l.WithValues("eid", eventID).Debug("instance start to process.")
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
