package wasm

import (
	"context"

	"github.com/iotexproject/w3bstream/pkg/enums"
)

type VM interface {
	Name() string
	Init()
	NewModule(code []byte) Module
}

type Module interface {
	Init()
	NewInstance() Instance
	GetABI() []string
}

type Instance interface {
	ID() string
	Start(context.Context) error
	Stop(context.Context) error
	State() enums.InstanceState
	AddResource(context.Context, []byte) uint32
	RmvResource(context.Context, uint32)
	GetResource(uint32) ([]byte, bool)
	Get(k string) int32

	EventConsumer
}

type EventHandleResult struct {
	InstanceID string           `json:"instanceID"`
	Rsp        []byte           `json:"-"`
	Code       ResultStatusCode `json:"code"`
	ErrMsg     string           `json:"errMsg"`
}

type EventConsumer interface {
	HandleEvent(ctx context.Context, handler string, payload []byte) *EventHandleResult
}

type KVStore interface {
	Get(string) int32
}

type ContextHandler interface {
	Name() string
	GetImports() ImportsHandler
	SetImports(ImportsHandler)
	GetExports() ExportsHandler
	GetInstance() Instance
	SetInstance(Instance)
}

type ABI struct {
	Imports  ImportsHandler
	Instance Instance
}

func (a *ABI) Name() string { return NameVersion }

func (a *ABI) GetExports() ExportsHandler { return a }

func (a *ABI) GetImports() ImportsHandler { return a.Imports }

func (a *ABI) SetImports(i ImportsHandler) { a.Imports = i }

func (a *ABI) GetInstance() Instance { return a.Instance }

func (a *ABI) SetInstance(i Instance) { a.Instance = i }

func (a *ABI) Start() {}

func (a *ABI) Alloc() {}

func (a *ABI) Free() {}

type Memory interface {
	Read(context.Context, uint32, uint32) ([]byte, error)
	Write(context.Context, []byte)
}

type ImportsHandler interface {
	GetDB(keyAddr, keySize, valAddr, valSize uint32) (code int32)
	SetDB()
	GetData()
	SetData()
	Log(level uint32)
}

type Handler interface {
	Name() string
	Call(context.Context, ...interface{})
}

type ExportsHandler interface {
	Start()
	Alloc()
	Free()
}
