package wasmtime

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bytecodealliance/wasmtime-go/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/reactivex/rxgo/v2"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
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
	ctx      context.Context
	id       types.SFID
	rt       *Runtime
	state    *atomic.Uint32
	res      *mapx.Map[uint32, []byte]
	evs      *mapx.Map[uint32, []byte]
	handlers map[string]*wasmtime.Func
	kvs      wasm.KVStore
	msgQueue chan *Task
	ch       chan rxgo.Item
}

func NewInstanceByCode(ctx context.Context, id types.SFID, code []byte, st enums.InstanceState) (i *Instance, err error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "NewInstanceByCode")
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
		ctx:      ctx,
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

	go ins.queueWorker()
	go func() {
		observable := ins.streamCompute(ins.ch)
		initSink(observable, ins.ctx, "db", "Customer")
	}()

	return ins, nil
}

var _ wasm.Instance = (*Instance)(nil)

func (i *Instance) ID() string { return i.id.String() }

func (i *Instance) Start(ctx context.Context) error {
	log.FromContext(ctx).WithValues("instance", i.ID()).Info("started")
	i.setState(enums.INSTANCE_STATE__STARTED)
	return nil
}

func (i *Instance) Stop(ctx context.Context) error {
	log.FromContext(ctx).WithValues("instance", i.ID()).Info("stopped")
	i.setState(enums.INSTANCE_STATE__STOPPED)
	return nil
}

func (i *Instance) setState(st wasm.InstanceState) {
	i.state.Store(uint32(st))
}

func (i *Instance) State() wasm.InstanceState { return wasm.InstanceState(i.state.Load()) }

func (i *Instance) HandleEvent(ctx context.Context, fn, eventType string, data []byte) *wasm.EventHandleResult {
	if i.State() != enums.INSTANCE_STATE__STARTED {
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			Code:       wasm.ResultStatusCode_Failed,
			ErrMsg:     "instance not running",
		}
	}

	select {
	case <-time.After(5 * time.Second):
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			Code:       wasm.ResultStatusCode_Failed,
			ErrMsg:     "fail to add the event to the VM",
		}
	case i.msgQueue <- newTask(ctx, fn, eventType, data):
		eventID := types.MustEventIDFromContext(ctx)
		log.FromContext(ctx).WithValues("eid", eventID).Debug("put task in queue.")

		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			Code:       wasm.ResultStatusCode_OK,
			ErrMsg:     "",
		}
	}
}

func (i *Instance) queueWorker() {
	for {
		task, more := <-i.msgQueue
		log.FromContext(task.ctx).WithValues("eid", task.EventID).Debug(
			fmt.Sprintf("queue len is %d and more is %t", len(i.msgQueue), more))
		if !more {
			return
		}

		log.FromContext(task.ctx).WithValues("eid", task.EventID).Debug("get task from queue.")

		if task.EventType == "OP_DEMO" {
			log.FromContext(task.ctx).WithValues("eid", task.EventID).Info("OP_DEMO start.")
			i.ch <- rxgo.Of(task)
			continue
		}

		res := i.handle(task.ctx, task)
		log.FromContext(task.ctx).WithValues("eid", task.EventID).Debug("event process completed.")
		if len(res.ErrMsg) > 0 {
			job.Dispatch(i.ctx, job.NewWasmLogTask(i.ctx, conflog.Level(log.ErrorLevel).String(), "vmTask", res.ErrMsg))
		} else {
			job.Dispatch(i.ctx, job.NewWasmLogTask(
				i.ctx,
				conflog.Level(log.InfoLevel).String(),
				"vmTask",
				fmt.Sprintf("the event, whose eventtype is %s, is successfully handled by %s, ", task.EventType, task.Handler),
			))
		}
	}
}

func (i *Instance) streamCompute(ch chan rxgo.Item) rxgo.Observable {
	return rxgo.FromChannel(ch).Filter(i.filterFunc).Map(i.mapFunc).
		GroupByDynamic(i.groupByKey, rxgo.WithBufferedChannel(10), rxgo.WithErrorStrategy(rxgo.ContinueOnError))
}

func initSink(observable rxgo.Observable, ctx context.Context, tye, schema string) {
	c := observable.Observe()
	for item := range c {

		switch item.V.(type) {
		case rxgo.GroupedObservable: // group operator
			go func() {
				obs := item.V.(rxgo.GroupedObservable)
				for i := range obs.Observe() {
					sink(ctx, i, tye, schema)
				}
			}()
		case rxgo.ObservableImpl: // window operator
			obs := item.V.(rxgo.ObservableImpl)
			for i := range obs.Observe() {
				sink(ctx, i, tye, schema)
			}
		default:
			sink(ctx, item, tye, schema)
		}
	}
}

func sink(ctx context.Context, item rxgo.Item, tye, schema string) {
	//customer := item.V.(models.Customer)
	conflog.Std().Info(fmt.Sprintf("customer: %v", string(item.V.(*Task).Payload)))

	//switch tye {
	//case "db":
	//	d, _ := wasm.SQLStoreFromContext(ctx)
	//	if err := customer.Create(d); err != nil {
	//		conflog.Std().Error(err)
	//	}
	//case "blockchain":
	//
	//default:
	//
	//}
}

func (i *Instance) filterFunc(inter interface{}) bool {
	res := false

	task := inter.(*Task)
	task.Handler = "filterAge"

	rid := i.AddResource(task.ctx, []byte(task.EventType), task.Payload)
	defer i.RmvResource(task.ctx, rid)

	code := i.handleByRid(task.ctx, task.Handler, rid).Code
	conflog.Std().Info(fmt.Sprintf("%s wasm code %d", task.Handler, code))

	if code < 0 {
		return res
	}

	rb, ok := i.GetResource(uint32(code))
	if !ok {
		conflog.Std().Error(errors.New("not found"))
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

	return res
}

func (i *Instance) mapFunc(c context.Context, inter interface{}) (interface{}, error) {
	task := inter.(*Task)
	task.Handler = "mapTax"

	rid := i.AddResource(task.ctx, []byte(task.EventType), task.Payload)
	defer i.RmvResource(task.ctx, rid)

	code := i.handleByRid(task.ctx, task.Handler, rid).Code
	conflog.Std().Info(fmt.Sprintf("mapTax wasm code %d", code))

	if code < 0 {
		conflog.Std().Error(errors.New(fmt.Sprintf("%s %s error.", string(inter.(*Task).Payload), "mapTax")))
		return nil, errors.New(fmt.Sprintf("%s %s error.", string(inter.(*Task).Payload), "mapTax"))
	}

	rb, ok := i.GetResource(uint32(code))
	if !ok {
		conflog.Std().Error(errors.New("mapTax result not found"))
		return nil, errors.New("mapTax result not found")
	}

	task.Payload = rb
	return task, nil
}

func (i *Instance) groupByKey(item rxgo.Item) string {
	task := item.V.(*Task)
	task.Handler = "groupByAge"

	rid := i.AddResource(task.ctx, []byte(task.EventType), task.Payload)
	defer i.RmvResource(task.ctx, rid)

	code := i.handleByRid(task.ctx, task.Handler, rid).Code
	conflog.Std().Info(fmt.Sprintf("groupByAge wasm code %d", code))

	if code < 0 {
		conflog.Std().Error(errors.New(fmt.Sprintf("%v %s error.", string(item.V.(*Task).Payload), "groupByAge")))
		return "error"
	}

	rb, ok := i.GetResource(uint32(code))
	if !ok {
		conflog.Std().Error(errors.New("groupByAge result not found"))
		return "error"
	}

	groupKey := string(rb)
	return groupKey
}

func (i *Instance) handleByRid(ctx context.Context, handlerName string, rid uint32) *wasm.EventHandleResult {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "instance.handleByRid")
	defer l.End()

	if err := i.rt.Instantiate(); err != nil {
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			ErrMsg:     err.Error(),
			Code:       wasm.ResultStatusCode_Failed,
		}
	}
	defer i.rt.Deinstantiate()

	result, err := i.rt.Call(handlerName, int32(rid))
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
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "instance.Handle")
	defer l.End()

	rid := i.AddResource(ctx, []byte(task.EventType), task.Payload)
	defer i.RmvResource(ctx, rid)

	if err := i.rt.Instantiate(); err != nil {
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			ErrMsg:     err.Error(),
			Code:       wasm.ResultStatusCode_Failed,
		}
	}
	defer i.rt.Deinstantiate()

	l.WithValues("eid", task.EventID).Debug("call wasm runtime.")

	// TODO support wasm return data(not only code) for HTTP responding
	result, err := i.rt.Call(task.Handler, int32(rid))
	l.WithValues("eid", task.EventID).Debug("call wasm runtime completed.")
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

func (i *Instance) AddResource(ctx context.Context, eventType, data []byte) uint32 {
	var id = int32(uuid.New().ID() % uint32(maxInt))
	i.res.Store(uint32(id), data)
	i.evs.Store(uint32(id), eventType)
	return uint32(id)
}

func (i *Instance) GetResource(id uint32) ([]byte, bool) {
	return i.res.Load(id)
}

func (i *Instance) RmvResource(ctx context.Context, id uint32) {
	i.res.Remove(id)
	i.evs.Remove(id)
}
