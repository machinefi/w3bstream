package wasmtime

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"time"

	"github.com/bytecodealliance/wasmtime-go"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/job"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func NewInstanceByCode(ctx context.Context, id types.SFID, code []byte) (*Instance, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "NewInstanceByCode")
	defer l.End()

	vmEngine := wasmtime.NewEngineWithConfig(wasmtime.NewConfig())
	vmStore := wasmtime.NewStore(vmEngine)
	linker := wasmtime.NewLinker(vmEngine)
	res := mapx.New[uint32, []byte]()
	db := make(map[string][]byte)

	cl, err := buildChainClient(l, ctx)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	ef := ExportFuncs{
		store: vmStore,
		res:   res,
		db:    db,
		cl:    cl,
		logger: types.MustLoggerFromContext(ctx).WithValues(
			"@src", "wasm",
			"@namespace", types.MustProjectFromContext(ctx).Name,
			"@applet", types.MustAppletFromContext(ctx).Name,
		),
	}
	_ = linker.FuncWrap("env", "ws_get_data", ef.GetData)
	_ = linker.FuncWrap("env", "ws_set_data", ef.SetData)
	_ = linker.FuncWrap("env", "ws_get_db", ef.GetDB)
	_ = linker.FuncWrap("env", "ws_set_db", ef.SetDB)
	_ = linker.FuncWrap("env", "ws_log", ef.Log)
	_ = linker.FuncWrap("env", "ws_send_tx", ef.SendTX)
	_ = linker.FuncWrap("env", "ws_call_contract", ef.CallContract)

	_ = linker.DefineWasi()

	wasiConfig := wasmtime.NewWasiConfig()
	vmStore.SetWasi(wasiConfig)

	vmModule, err := wasmtime.NewModule(vmEngine, code)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	vmInstance, err := linker.Instantiate(vmStore, vmModule)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	return &Instance{
		id:         id,
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

func buildChainClient(l log.Logger, ctx context.Context) (*ChainClient, error) {
	ethConf, ok := types.ETHClientConfigFromContext(ctx)
	if !ok {
		return nil, errors.New("fail to read eth client conf")
	}
	if len(ethConf.ChainEndpoint) == 0 {
		l.Warn(errors.New("no chain client is established due to empty chain endpoint"))
		return nil, nil
	}
	chain, err := ethclient.Dial(ethConf.ChainEndpoint)
	if err != nil {
		l.Error(errors.New("fail to dial the endpoint of the chain"))
		return nil, err
	}
	var pvk *ecdsa.PrivateKey
	if len(ethConf.PrivateKey) > 0 {
		pvk = crypto.ToECDSAUnsafe(gethCommon.FromHex(ethConf.PrivateKey))
	}
	return &ChainClient{
		pvk:   pvk,
		chain: chain,
	}, nil
}

type Instance struct {
	id         types.SFID
	state      wasm.InstanceState
	vmEngine   *wasmtime.Engine
	vmStore    *wasmtime.Store
	vmModule   *wasmtime.Module
	vmInstance *wasmtime.Instance
	res        *mapx.Map[uint32, []byte]
	handlers   map[string]*wasmtime.Func
	db         map[string][]byte
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

	hdl, ok := i.handlers[t.Handler]
	if !ok {
		hdl = i.vmInstance.GetFunc(i.vmStore, t.Handler)
		if hdl == nil {
			return &wasm.EventHandleResult{
				InstanceID: i.id.String(),
				ErrMsg:     "handler not exists",
				Code:       wasm.ResultStatusCode_UnexportedHandler,
			}
		}
		i.handlers[t.Handler] = hdl
	}

	l.Info("call handler:%s", t.Handler)

	result, err := hdl.Call(i.vmStore, int32(rid))
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
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "instance.AddResource")
	defer l.End()

	var id = int32(uuid.New().ID() % uint32(MaxInt))
	i.res.Store(uint32(id), data)

	l.WithValues("res_id", id, "payload", string(data)).Info("added")

	return uint32(id)
}

func (i *Instance) GetResource(id uint32) ([]byte, bool) { return i.res.Load(id) }

func (i *Instance) RmvResource(ctx context.Context, id uint32) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "instance.RmvResource")
	defer l.End()

	i.res.Remove(id)
	l.WithValues("res_id", id).Info("removed")

}

func (i *Instance) Get(k string) int32 {
	data := i.db[k]
	var ret int32
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return ret
}
