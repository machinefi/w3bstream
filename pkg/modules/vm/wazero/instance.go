package wazero

import (
	"context"

	"github.com/google/uuid"
	"github.com/iotexproject/Bumblebee/x/mapx"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/modules/vm/common"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

var defaultRuntimeConfig = wazero.NewRuntimeConfig().
	WithFeatureBulkMemoryOperations(true).
	WithFeatureNonTrappingFloatToIntConversion(true).
	WithFeatureSignExtensionOps(true).
	WithFeatureMultiValue(true)

func NewInstanceByCode(code []byte, opts ...common.InstanceOptionSetter) (wasm.Instance, error) {
	ctx := context.Background()
	opt := &common.InstanceOption{
		Logger: common.DefaultLogger,
		Tasks:  &common.TaskQueue{Ch: make(chan *common.Task)},
	}

	for _, set := range opts {
		set(opt)
	}

	i := &Instance{
		opt:      opt,
		state:    enums.INSTANCE_STATE__CREATED,
		rt:       wazero.NewRuntimeWithConfig(ctx, defaultRuntimeConfig),
		handlers: make(map[string]api.Function),
		db:       make(map[string]int32),
	}

	{
		_, err := i.rt.NewModuleBuilder("env").
			ExportFunction("ws_get_data", i.GetData).
			ExportFunction("ws_set_data", i.SetData).
			ExportFunction("ws_get_db", i.GetDB).
			ExportFunction("ws_set_db", i.SetDB).
			ExportFunction("ws_log", i.Log).
			Instantiate(ctx, i.rt)
		if err != nil {
			return nil, err
		}
	}

	_, err := wasi_snapshot_preview1.Instantiate(ctx, i.rt)
	if err != nil {
		return nil, err
	}

	i.mod, err = i.rt.InstantiateModuleFromBinary(ctx, code)
	if err != nil {
		return nil, err
	}
	i.alloc = i.mod.ExportedFunction("alloc")
	i.free = i.mod.ExportedFunction("free")

	i.res = mapx.New[uint32, []byte]()
	i.ctx, i.cancel = context.WithCancel(ctx)

	return i, nil
}

type Instance struct {
	opt      *common.InstanceOption
	ctx      context.Context
	cancel   context.CancelFunc
	state    wasm.InstanceState
	rt       wazero.Runtime
	mod      api.Module
	res      *mapx.Map[uint32, []byte]
	alloc    api.Function
	free     api.Function
	handlers map[string]api.Function
	db       map[string]int32
}

var _ wasm.Instance = (*Instance)(nil)

func (i *Instance) Start() error {
	for {
		select {
		case <-i.ctx.Done():
			return i.ctx.Err()
		case task := <-i.opt.Tasks.Wait():
			task.Res <- i.handleEvent(task)
		}
	}
}

func (i *Instance) Stop() {
	i.state = enums.INSTANCE_STATE__STOPPED
	i.cancel()
}

func (i *Instance) State() wasm.InstanceState { return i.state }

func (i *Instance) HandleEvent(fn string, data []byte) ([]byte, wasm.ResultStatusCode) {
	task := &common.Task{
		Handler: fn,
		Payload: data,
		Res:     make(chan *common.EventHandleResult),
	}
	i.opt.Tasks.Push(task)

	res := <-task.Res
	return res.Response, res.Code
}

func (i *Instance) handleEvent(t *common.Task) *common.EventHandleResult {
	rid := i.AddResource(t.Payload)
	defer i.RmvResource(rid)

	hdl, ok := i.handlers[t.Handler]
	if !ok {
		hdl = i.mod.ExportedFunction(t.Handler)
		if hdl == nil {
			return &common.EventHandleResult{nil, wasm.ResultStatusCode_UnexportedHandler}
		}
		i.handlers[t.Handler] = hdl
	}

	results, err := hdl.Call(i.ctx, uint64(rid))
	if err != nil {
		return &common.EventHandleResult{nil, wasm.ResultStatusCode_Failed}
	}

	return &common.EventHandleResult{nil, wasm.ResultStatusCode(results[0])}
}

func (i *Instance) AddResource(data []byte) uint32 {
	id := uuid.New().ID()
	i.res.Store(id, data)
	return id
}

func (i *Instance) GetResource(id uint32) ([]byte, bool) { return i.res.Load(id) }

func (i *Instance) RmvResource(id uint32) { i.res.Remove(id) }

func (i *Instance) Get(k string) int32 { return i.db[k] }
