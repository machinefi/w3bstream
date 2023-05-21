package event

import (
	"context"
	"fmt"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
	"strconv"
	"sync"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

var Handler = func(ctx context.Context, data []byte) []*Result {
	return OnEvent(ctx, data)
}

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
	return OnEvent(ctx, data), nil
}

func OnEvent(ctx context.Context, data []byte) (ret []*Result) {
	var (
		l   = types.MustLoggerFromContext(ctx)
		r   = types.MustStrategyResultsFromContext(ctx)
		rDB = kvdb.MustRedisDBKeyFromContext(ctx)
		prj = types.MustProjectFromContext(ctx)

		results = make(chan *Result, len(r))

		valByte []byte
	)

	m, err := trafficlimit.GetByProjectAndType(ctx, prj.ProjectID, enums.TRAFFIC_LIMIT_TYPE__EVENT)
	if err != nil {
		se, ok := statusx.IsStatusErr(err)
		if !ok || !se.Is(status.TrafficLimitNotFound) {
			ret = append(ret, &Result{
				AppletName:  "",
				InstanceID:  0,
				Handler:     "",
				ReturnValue: nil,
				ReturnCode:  -1,
				Error:       err.Error(),
			})
			return
		}
		l.Warn(err)
	}
	if m != nil {
		if valByte, err = rDB.IncrBy(fmt.Sprintf("%s::%s", prj.Name, m.ApiType.String()), []byte(strconv.Itoa(-1))); err != nil {
			l.Error(err)
			ret = append(ret, &Result{
				AppletName:  "",
				InstanceID:  0,
				Handler:     "",
				ReturnValue: nil,
				ReturnCode:  -1,
				Error:       status.DatabaseError.StatusErr().WithDesc(err.Error()).Key,
			})
			return
		}
		val, _ := strconv.Atoi(string(valByte))
		if val < 0 {
			ret = append(ret, &Result{
				AppletName:  "",
				InstanceID:  0,
				Handler:     "",
				ReturnValue: nil,
				ReturnCode:  -1,
				Error:       status.TrafficLimitExceeded.Key(),
			})
			return
		}
	}

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
