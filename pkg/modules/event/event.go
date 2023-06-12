package event

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
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

	strategies, err := strategy.FilterByProjectAndEvent(ctx, prj.ProjectID, t)
	if err != nil {
		return nil, err
	}

	ctx = types.WithStrategyResults(ctx, strategies)

	eventID := uuid.NewString() + "_monitor"
	ctx = types.WithEventID(ctx, eventID)

	if err := TrafficLimitEvent(ctx); err != nil {
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

	return OnEvent(ctx, data), nil
}

func TrafficLimitEvent(ctx context.Context) error {
	var (
		l   = types.MustLoggerFromContext(ctx)
		rDB = kvdb.MustRedisDBKeyFromContext(ctx)
		prj = types.MustProjectFromContext(ctx)

		valByte []byte
	)

	m, err := trafficlimit.GetByProjectAndType(ctx, prj.ProjectID, enums.TRAFFIC_LIMIT_TYPE__EVENT)
	if err != nil {
		se, ok := statusx.IsStatusErr(err)
		if !ok || !se.Is(status.TrafficLimitNotFound) {
			return err
		}
		l.Warn(err)
	}
	if m != nil {
		if valByte, err = rDB.IncrBy(fmt.Sprintf("%s::%s", prj.Name, m.ApiType.String()), []byte(strconv.Itoa(-1))); err != nil {
			l.Error(err)
			return status.DatabaseError.StatusErr().WithDesc(err.Error())
		}
		val, _ := strconv.Atoi(string(valByte))
		if val < 0 {
			return status.TrafficLimitExceededFailed
		}
	}
	return nil
}

func OnEvent(ctx context.Context, data []byte) (ret []*Result) {
	var (
		l       = types.MustLoggerFromContext(ctx)
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
