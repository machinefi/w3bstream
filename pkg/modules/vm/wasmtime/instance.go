package wasmtime

import (
	"context"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/google/uuid"
	"github.com/iotexproject/Bumblebee/x/mapx"

	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/modules/vm/common"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

func NewInstanceByCode(code []byte, opts ...common.InstanceOptionSetter) (wasm.Instance, error) {
	ctx := context.Background()

	opt := &common.InstanceOption{
		Logger: common.DefaultLogger,
		Tasks:  &common.TaskQueue{Ch: make(chan *common.Task)},
	}

	for _, set := range opts {
		set(opt)
	}

	res := mapx.New[uint32, []byte]()
	db := make(map[string]int32)

	vmEngine := wasmtime.NewEngineWithConfig(wasmtime.NewConfig())
	vmStore := wasmtime.NewStore(vmEngine)
	linker := wasmtime.NewLinker(vmEngine)

	ef := ExportFuncs{vmStore, res, db, opt.Logger}
	_ = linker.FuncWrap("env", "ws_get_data", ef.GetData)
	_ = linker.FuncWrap("env", "ws_set_data", ef.SetData)
	_ = linker.FuncWrap("env", "ws_get_db", ef.GetDB)
	_ = linker.FuncWrap("env", "ws_set_db", ef.SetDB)
	_ = linker.FuncWrap("env", "ws_log", ef.Log)

	_ = linker.DefineWasi()

	wasiConfig := wasmtime.NewWasiConfig()
	vmStore.SetWasi(wasiConfig)

	vmModule, err := wasmtime.NewModule(vmEngine, code)
	if err != nil {
		return nil, err
	}
	vmInstance, err := linker.Instantiate(vmStore, vmModule)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &Instance{
		tasks:      opt.Tasks,
		ctx:        ctx,
		cancel:     cancel,
		state:      enums.INSTANCE_STATE__CREATED,
		vmEngine:   vmEngine,
		vmStore:    vmStore,
		vmModule:   vmModule,
		vmInstance: vmInstance,
		res:        res,
		handlers:   make(map[string]*wasmtime.Func),
		db:         db,
	}, nil
}

type Instance struct {
	tasks      *common.TaskQueue
	ctx        context.Context
	cancel     context.CancelFunc
	state      wasm.InstanceState
	vmEngine   *wasmtime.Engine
	vmStore    *wasmtime.Store
	vmModule   *wasmtime.Module
	vmInstance *wasmtime.Instance
	res        *mapx.Map[uint32, []byte]
	handlers   map[string]*wasmtime.Func
	db         map[string]int32
}

var _ wasm.Instance = (*Instance)(nil)

func (i *Instance) Start() error {
	for {
		select {
		case <-i.ctx.Done():
			return i.ctx.Err()
		case task := <-i.tasks.Wait():
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
	i.tasks.Push(task)

	res := <-task.Res
	return res.Response, res.Code
}

func (i *Instance) handleEvent(t *common.Task) *common.EventHandleResult {
	rid := i.AddResource(t.Payload)
	defer i.RmvResource(rid)

	hdl, ok := i.handlers[t.Handler]
	if !ok {
		hdl = i.vmInstance.GetFunc(i.vmStore, t.Handler)
		if hdl == nil {
			return &common.EventHandleResult{nil, wasm.ResultStatusCode_UnexportedHandler}
		}
		i.handlers[t.Handler] = hdl
	}

	result, err := hdl.Call(i.vmStore, int32(rid))
	if err != nil {
		return &common.EventHandleResult{nil, wasm.ResultStatusCode_Failed}
	}

	return &common.EventHandleResult{nil, wasm.ResultStatusCode(result.(int32))}
}

func (i *Instance) AddResource(data []byte) uint32 {
	id := uuid.New().ID()
	i.res.Store(id, data)
	return id
}

func (i *Instance) GetResource(id uint32) ([]byte, bool) { return i.res.Load(id) }

func (i *Instance) RmvResource(id uint32) { i.res.Remove(id) }

func (i *Instance) Get(k string) int32 { return i.db[k] }
