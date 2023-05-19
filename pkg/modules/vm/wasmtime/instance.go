package wasmtime

import (
	"context"
	"fmt"
	"time"

	"github.com/bytecodealliance/wasmtime-go/v8"
	"github.com/google/uuid"

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
	state    wasm.InstanceState
	res      *mapx.Map[uint32, []byte]
	evs      *mapx.Map[uint32, []byte]
	handlers map[string]*wasmtime.Func
	kvs      wasm.KVStore
	msgQueue chan *Task
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

	ins := &Instance{
		ctx:      ctx,
		rt:       rt,
		id:       id,
		state:    st,
		res:      res,
		evs:      evs,
		handlers: make(map[string]*wasmtime.Func),
		kvs:      wasm.MustKVStoreFromContext(ctx),
		msgQueue: make(chan *Task, maxMsgPerInstance),
	}

	go ins.queueWorker()

	return ins, nil
}

var _ wasm.Instance = (*Instance)(nil)

func (i *Instance) ID() string { return i.id.String() }

func (i *Instance) Start(ctx context.Context) error {
	log.FromContext(ctx).WithValues("instance", i.ID()).Info("started")
	i.state = enums.INSTANCE_STATE__STARTED
	return nil
}

func (i *Instance) Stop(ctx context.Context) error {
	log.FromContext(ctx).WithValues("instance", i.ID()).Info("stopped")
	i.state = enums.INSTANCE_STATE__STOPPED
	return nil
}

func (i *Instance) State() wasm.InstanceState { return i.state }

func (i *Instance) HandleEvent(ctx context.Context, fn, eventType string, data []byte) *wasm.EventHandleResult {
	if i.state != enums.INSTANCE_STATE__STARTED {
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
		if !more {
			return
		}
		res := i.handle(task.ctx, task)
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

	// TODO support wasm return data(not only code) for HTTP responding
	result, err := i.rt.Call(task.Handler, int32(rid))
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
