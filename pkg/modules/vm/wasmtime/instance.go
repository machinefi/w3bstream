package wasmtime

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bytecodealliance/wasmtime-go/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/reactivex/rxgo/v2"
	"github.com/tidwall/gjson"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/job"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

const (
	maxUint = ^uint32(0)
	maxInt  = int(maxUint >> 1)
	// TODO: add into config
	maxMsgPerInstance = 5000
)

type Instance struct {
	ctx         context.Context
	id          types.SFID
	rt          *Runtime
	state       *atomic.Uint32
	res         *mapx.Map[uint32, []byte]
	evs         *mapx.Map[uint32, []byte]
	handlers    map[string]*wasmtime.Func
	kvs         wasm.KVStore
	msgQueue    chan *Task
	ch          chan rxgo.Item
	source      []string
	operators   []wasm.Operator
	simpleOpMap map[string]string
	windOps     []wasm.Operator
	windOpMap   map[string]string
	sink        wasm.Sink
}

func NewInstanceByCode(ctx context.Context, id types.SFID, code []byte, st enums.InstanceState) (i *Instance, err error) {
	ctx, l := logr.Start(ctx, "modules.vm.wasmtime.NewInstanceByCode")
	defer l.End()

	res := mapx.New[uint32, []byte]()
	evs := mapx.New[uint32, []byte]()
	rt := NewRuntime()
	lk, err := NewExportFuncs(contextx.WithContextCompose(
		wasm.WithRuntimeResourceContext(res),
		wasm.WithRuntimeEventTypesContext(evs),
	)(ctx), rt)
	if err != nil {
		return nil, err
	}
	if err := rt.Link(lk, code); err != nil {
		return nil, err
	}
	state := &atomic.Uint32{}
	state.Store(uint32(st))

	ins := &Instance{
		rt:       rt,
		id:       id,
		state:    state,
		res:      res,
		evs:      evs,
		handlers: make(map[string]*wasmtime.Func),
		kvs:      wasm.MustKVStoreFromContext(ctx),
		msgQueue: make(chan *Task, maxMsgPerInstance),
		ch:       make(chan rxgo.Item),
	}

	flow, ok := wasm.FlowFromContext(ctx)
	if ok {
		ins.source = flow.Source.Strategies
		ins.operators = flow.Operators
		ins.simpleOpMap = make(map[string]string)
		ins.windOpMap = make(map[string]string)
		ins.windOps = make([]wasm.Operator, 0)
		ins.sink = flow.Sink
		go func() {
			observable := ins.streamCompute(ins.ch)
			ins.initSink(ins.ctx, observable)
		}()
	}

	return ins, nil
}

var _ wasm.Instance = (*Instance)(nil)

func (i *Instance) ID() string { return i.id.String() }

func (i *Instance) Start(ctx context.Context) error {
	ctx, l := logr.Start(ctx, "modules.vm.Instance.Start", "instance_id", i.ID())
	defer l.End()

	i.state.Store(uint32(enums.INSTANCE_STATE__STARTED))
	return nil
}

func (i *Instance) Stop(ctx context.Context) error {
	ctx, l := logr.Start(ctx, "modules.vm.Instance.Stop", "instance_id", i.ID())
	defer l.End()

	i.state.Store(uint32(enums.INSTANCE_STATE__STOPPED))
	return nil
}

func (i *Instance) State() wasm.InstanceState { return wasm.InstanceState(i.state.Load()) }

func (i *Instance) HandleEvent(ctx context.Context, fn, eventType string, data []byte) *wasm.EventHandleResult {
	ctx, l := logr.Start(ctx, "modules.vm.wasmtime.Instance.HandleEvent")
	defer l.End()

	if i.State() != enums.INSTANCE_STATE__STARTED {
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			Code:       wasm.ResultStatusCode_Failed,
			ErrMsg:     "instance not running",
		}
	}

	task := &Task{
		EventID:   types.MustEventIDFromContext(ctx),
		EventType: eventType,
		Handler:   fn,
		Payload:   data,
		TaskState: mq.TASK_STATE__PENDING,
		vm:        i,
		retrieve:  make(chan *wasm.EventHandleResult),
	}

	job.Dispatch(ctx, task)
	return task.Wait()
}

func (i *Instance) queueWorker(ctx context.Context) {
	for {
		res := &wasm.EventHandleResult{}
		task, more := <-i.msgQueue

		ctx, l := logger.NewSpanContext(ctx, "modules.vm.wasmtime.Instance.queueWorker")
		l.WithValues("queue", len(i.msgQueue), "more", more).Debug("")

		if !more {
			return
		}
		if task == nil {
			return
		}
		l = l.WithValues("event_id", task.EventID)

		for _, typ := range i.source {
			if task.EventType == typ {
				l.Info("Flow_op start.")
				i.ch <- rxgo.Of(task)
				goto Next
			}
		}

		res = i.handle(ctx, task)
		if len(res.ErrMsg) > 0 {
			job.Dispatch(ctx, job.NewWasmLogTask(ctx, log.ErrorLevel.String(), "vmTask", res.ErrMsg))
		} else {
			job.Dispatch(ctx, job.NewWasmLogTask(ctx, log.InfoLevel.String(), "vmTask", fmt.Sprintf("the event, whose eventtype is %s, is successfully handled by %s, ", task.EventType, task.Handler)))
		}
	Next:
		l.End()
	}
}

func (i *Instance) streamCompute(ch chan rxgo.Item) rxgo.Observable {
	obs := rxgo.FromChannel(ch)
	for index, op := range i.operators {
		switch {
		case op.OpType == enums.FLOW_OPERATOR__FILTER:
			filterNum := index
			i.simpleOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__FILTER, filterNum)] = op.WasmFunc

			obs = obs.Filter(func(inter interface{}) bool {
				start := time.Now()
				res := false
				task := inter.(*Task)
				task.Handler = i.simpleOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__FILTER, filterNum)]

				rb, ok := i.runOp(task)
				if !ok {
					conflog.Std().Error(errors.New(fmt.Sprintf("%s result not found", op.WasmFunc)))
					return res
				}

				result := strings.ToLower(string(rb))
				if result == "true" {
					res = true
				} else if result == "false" {
					res = false
				} else {
					conflog.Std().Warn(errors.New("the value does not support"))
				}
				duration := time.Since(start)
				conflog.Std().Info(fmt.Sprintf("%s template cost %s", task.Handler, duration.String()))
				return res
			})
		case op.OpType == enums.FLOW_OPERATOR__MAP:
			mapNum := index
			i.simpleOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__MAP, mapNum)] = op.WasmFunc

			obs = obs.Map(func(ctx context.Context, inter interface{}) (interface{}, error) {
				start := time.Now()
				task := inter.(*Task)
				task.Handler = i.simpleOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__MAP, mapNum)]

				rb, ok := i.runOp(task)
				if !ok {
					conflog.Std().Error(errors.New(fmt.Sprintf("%s result not found", op.WasmFunc)))
					return nil, errors.New(fmt.Sprintf("%s result not found", op.WasmFunc))
				}

				task.Payload = rb
				duration := time.Since(start)
				conflog.Std().Info(fmt.Sprintf("%s template cost %s", task.Handler, duration.String()))
				return task, nil
			})
		case op.OpType == enums.FLOW_OPERATOR__WINDOW:
			obs = obs.WindowWithTime(rxgo.WithDuration(60 * time.Second))
		case op.OpType > enums.FLOW_OPERATOR__WINDOW:
			i.windOps = append(i.windOps, op)
		}
	}

	return obs
}

func (i *Instance) initSink(ctx context.Context, observable rxgo.Observable) {
	c := observable.Observe()
	for item := range c {

		switch item.V.(type) {
		case rxgo.GroupedObservable: // group operator
			go func() {
				obs := item.V.(rxgo.GroupedObservable)
				// add other op like reduce
				for it := range obs.Observe() {
					i.sinkData(ctx, it)
				}
			}()
		case *rxgo.ObservableImpl: // window operator
			var (
				obs   = item.V
				index = 0
				op    = wasm.Operator{}
			)

			for index, op = range i.windOps {
				switch op.OpType {
				// last op
				case enums.FLOW_OPERATOR__REDUCE:
					reduceNum := index
					i.windOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__REDUCE, reduceNum)] = op.WasmFunc

					obs = obs.(*rxgo.ObservableImpl).Reduce(func(ctx context.Context, inter1 interface{}, inter2 interface{}) (interface{}, error) {
						start := time.Now()
						var task1, task2 *Task
						task2 = inter2.(*Task)
						task2.Handler = i.windOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__REDUCE, reduceNum)]

						tasks := make([]*Task, 0)
						if inter1 != nil {
							task1 = inter1.(*Task)
						}
						tasks = append(tasks, task1)
						tasks = append(tasks, task2)

						rb, ok := i.runOp(tasks...)
						if !ok {
							conflog.Std().Error(errors.New(fmt.Sprintf("%s result not found", op.WasmFunc)))
							return nil, errors.New(fmt.Sprintf("%s result not found", op.WasmFunc))
						}

						task2.Payload = rb
						duration := time.Since(start)
						conflog.Std().Info(fmt.Sprintf("%s template cost %s", task2.Handler, duration.String()))
						return task2, nil
					})
				case enums.FLOW_OPERATOR__GROUP:
					groupNum := index
					i.windOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__GROUP, groupNum)] = op.WasmFunc

					obs = obs.(*rxgo.ObservableImpl).GroupByDynamic(func(item rxgo.Item) string {
						start := time.Now()
						task := item.V.(*Task)
						task.Handler = i.windOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__GROUP, groupNum)]

						rb, ok := i.runOp(task)
						if !ok {
							conflog.Std().Error(errors.New(fmt.Sprintf("%s result not found", op.WasmFunc)))
							return "error"
						}

						groupKey := string(rb)
						duration := time.Since(start)
						conflog.Std().Info(fmt.Sprintf("%s template cost %s", task.Handler, duration.String()))
						return groupKey
					}, rxgo.WithBufferedChannel(2), rxgo.WithErrorStrategy(rxgo.ContinueOnError))
					goto skip
				}
			}

		skip:
			switch obs.(type) {
			case rxgo.OptionalSingle:
				for it := range obs.(rxgo.OptionalSingle).Observe() {
					i.sinkData(ctx, it)
				}
			case *rxgo.ObservableImpl:
				for it := range obs.(*rxgo.ObservableImpl).Observe() {
					// check group or common
					switch it.V.(type) {
					case rxgo.GroupedObservable:
						go func() {
							grpObs := it.V
							op := wasm.Operator{}
							// add other op like reduce
							// there are other ops after group op, should add here
							if index < len(i.windOps)-1 {
								for j := index; j < len(i.windOps); j++ {
									op = i.windOps[j]
									switch op.OpType {
									case enums.FLOW_OPERATOR__REDUCE:
										reduceNum := j
										i.windOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__REDUCE, reduceNum)] = op.WasmFunc

										grpObs = grpObs.(rxgo.GroupedObservable).Reduce(func(ctx context.Context, inter1 interface{}, inter2 interface{}) (interface{}, error) {
											start := time.Now()
											var task1, task2 *Task
											task2 = inter2.(*Task)
											task2.Handler = i.windOpMap[fmt.Sprintf("%s_%d", enums.FLOW_OPERATOR__REDUCE, reduceNum)]

											tasks := make([]*Task, 0)
											if inter1 != nil {
												task1 = inter1.(*Task)
											}
											tasks = append(tasks, task1)
											tasks = append(tasks, task2)

											rb, ok := i.runOp(tasks...)
											if !ok {
												conflog.Std().Error(errors.New(fmt.Sprintf("%s result not found", op.WasmFunc)))
												return nil, errors.New(fmt.Sprintf("%s result not found", op.WasmFunc))
											}

											task2.Payload = rb
											duration := time.Since(start)
											conflog.Std().Info(fmt.Sprintf("%s template cost %s", task2.Handler, duration.String()))
											return task2, nil
										})
									}
								}
							}
							switch grpObs.(type) {
							case rxgo.OptionalSingle:
								for it := range grpObs.(rxgo.OptionalSingle).Observe() {
									i.sinkData(ctx, it)
								}
							case *rxgo.ObservableImpl:
								for it := range grpObs.(*rxgo.ObservableImpl).Observe() {
									i.sinkData(ctx, it)
								}
							default:
								i.sinkData(ctx, it)
							}
						}()
					default:
						i.sinkData(ctx, it)
					}
				}
			}
		default:
			i.sinkData(ctx, item)
		}
	}
}

func (i *Instance) sinkData(ctx context.Context, item rxgo.Item) {
	rowByte := item.V.(*Task).Payload

	switch i.sink.SinkType {
	case enums.FLOW_SINK__RMDB:
		db, err := sql.Open(i.sink.SinkInfo.DBInfo.DBType, i.sink.SinkInfo.DBInfo.Endpoint)
		if err != nil {
			conflog.Std().Error(err)
		}
		err = db.Ping()
		if err != nil {
			conflog.Std().Error(err)
		}

		sqlStringPrefix := fmt.Sprintf("INSERT INTO %s (", i.sink.SinkInfo.DBInfo.Table)
		sqlStringSuffix := fmt.Sprintf(") VALUES (")
		params := make([]interface{}, 0)
		for index, c := range i.sink.SinkInfo.DBInfo.Columns {
			params = append(params, gjson.GetBytes(rowByte, c).String())
			sqlStringPrefix = sqlStringPrefix + c + ","
			sqlStringSuffix = sqlStringSuffix + "$" + strconv.Itoa(index+1) + ","
		}
		sqlString := fmt.Sprintf("%s%s);", sqlStringPrefix[:len(sqlStringPrefix)-1], sqlStringSuffix[:len(sqlStringSuffix)-1])

		_, err = db.ExecContext(context.Background(), sqlString, params...)
		if err != nil {
			conflog.Std().Error(err)
		}
	case enums.FLOW_SINK__BLOCKCHAIN:

	default:

	}
}

func (i *Instance) runOp(task ...*Task) ([]byte, bool) {
	var (
		ctx     context.Context
		handler string

		rids = make([]interface{}, 0)
	)

	for _, t := range task {
		// if task is nil,  set rid is 0
		var rid uint32 = 0
		if t != nil {
			rid = i.AddResource([]byte(t.EventType), t.Payload)
			// ctx = t.ctx
			handler = t.Handler
		}

		rids = append(rids, int32(rid))
	}
	defer func() {
		for _, rid := range rids {
			i.RmvResource(uint32(rid.(int32)))
		}
	}()

	start := time.Now()
	code := i.handleByRid(ctx, handler, rids...).Code
	duration := time.Since(start)
	conflog.Std().Info(fmt.Sprintf("%s wasm cost %s", handler, duration.String()))

	conflog.Std().Info(fmt.Sprintf("%s wasm code %d", handler, code))

	if code < 0 {
		conflog.Std().Error(errors.New(fmt.Sprintf("%s wasm code run error", handler)))
		return nil, false
	}

	return i.GetResource(uint32(code))
}

func (i *Instance) handleByRid(ctx context.Context, handlerName string, rids ...interface{}) *wasm.EventHandleResult {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "instance.handleByRid")
	defer l.End()

	if err := i.rt.Instantiate(ctx); err != nil {
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			ErrMsg:     err.Error(),
			Code:       wasm.ResultStatusCode_Failed,
		}
	}
	defer i.rt.Deinstantiate(ctx)

	result, err := i.rt.Call(ctx, handlerName, rids...)
	if err != nil {
		l.Error(err)
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			ErrMsg:     err.Error(),
			Code:       wasm.ResultStatusCode_Failed,
		}
	}

	return &wasm.EventHandleResult{
		InstanceID: i.id.String(),
		Code:       wasm.ResultStatusCode(result.(int32)),
	}
}

func (i *Instance) handle(ctx context.Context, task *Task) *wasm.EventHandleResult {
	ctx, l := logr.Start(ctx, "modules.vm.wasmtime.Instance.handle",
		"event_id", task.EventID,
		"instance_id", i.id,
	)
	defer l.End()

	l.Info("start processing task")
	rid := i.AddResource([]byte(task.EventType), task.Payload)
	defer i.RmvResource(rid)

	if err := i.rt.Instantiate(ctx); err != nil {
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			ErrMsg:     err.Error(),
			Code:       wasm.ResultStatusCode_Failed,
		}
	}
	defer i.rt.Deinstantiate(ctx)

	// TODO support wasm return data(not only code) for HTTP responding
	result, err := i.rt.Call(ctx, task.Handler, int32(rid))
	l.Debug("call wasm runtime completed.")
	if err != nil {
		l.Error(err)
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			ErrMsg:     err.Error(),
			Code:       wasm.ResultStatusCode_Failed,
		}
	}

	return &wasm.EventHandleResult{
		InstanceID: i.id.String(),
		Code:       wasm.ResultStatusCode(result.(int32)),
	}
}

func (i *Instance) AddResource(eventType, data []byte) uint32 {
	var id = int32(uuid.New().ID() % uint32(maxInt))
	i.res.Store(uint32(id), data)
	i.evs.Store(uint32(id), eventType)
	return uint32(id)
}

func (i *Instance) GetResource(id uint32) ([]byte, bool) {
	return i.res.Load(id)
}

func (i *Instance) RmvResource(id uint32) {
	i.res.Remove(id)
	i.evs.Remove(id)
}
