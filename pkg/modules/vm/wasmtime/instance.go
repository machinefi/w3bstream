package wasmtime

import (
	"bytes"
	"context"
	"encoding/binary"
	"time"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/google/uuid"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/job"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func NewInstanceByCode(ctx context.Context, id types.SFID, code []byte) (i *Instance, err error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "NewInstanceByCode")
	defer l.End()

	res := mapx.New[uint32, []byte]()
	rt := NewRuntime()
	lk, err := NewExportFuncs(wasm.WithRuntimeResource(ctx, res), rt)
	if err != nil {
		return nil, err
	}
	if err := rt.Initiate(lk, code); err != nil {
		return nil, err
	}

	return &Instance{
		rt:       rt,
		id:       id,
		state:    enums.INSTANCE_STATE__CREATED,
		res:      res,
		handlers: make(map[string]*wasmtime.Func),
		kvs:      wasm.MustKVStoreFromContext(ctx),
	}, nil
}

type Instance struct {
	id       types.SFID
	rt       *Runtime
	state    wasm.InstanceState
	res      *mapx.Map[uint32, []byte]
	handlers map[string]*wasmtime.Func
	kvs      wasm.KVStore
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

func (i *Instance) HandleEvent(ctx context.Context, fn string, data []byte) *wasm.EventHandleResult {
	if i.state != enums.INSTANCE_STATE__STARTED {
		return &wasm.EventHandleResult{
			InstanceID: i.id.String(),
			Code:       wasm.ResultStatusCode_Failed,
			ErrMsg:     "instance not running",
		}
	}

	t := NewTask(i, fn, data)
	job.Dispatch(ctx, t)
	return t.Wait(time.Second * 5)
}

func (i *Instance) Handle(ctx context.Context, t *Task) *wasm.EventHandleResult {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "instance.Handle")
	defer l.End()

	rid := i.AddResource(ctx, t.Payload)
	defer i.RmvResource(ctx, rid)

	result, err := i.rt.Call(t.Handler, int32(rid))
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

const MaxUint = ^uint32(0)
const MaxInt = int(MaxUint >> 1)

func (i *Instance) AddResource(ctx context.Context, data []byte) uint32 {
	var id = int32(uuid.New().ID() % uint32(MaxInt))
	i.res.Store(uint32(id), data)
	return uint32(id)
}

func (i *Instance) GetResource(id uint32) ([]byte, bool) {
	return i.res.Load(id)
}

func (i *Instance) RmvResource(ctx context.Context, id uint32) {
	i.res.Remove(id)
}

func (i *Instance) Get(k string) int32 {
	data, _ := i.kvs.Get(k)
	var ret int32
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return ret
}
