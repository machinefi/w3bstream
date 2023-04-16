package wasm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
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
	// AddResource(context.Context, []byte) uint32
	// RmvResource(context.Context, uint32)
	// GetResource(uint32) ([]byte, bool)
	// Get(k string) int32

	EventConsumer
}

type EventHandleResult struct {
	InstanceID string           `json:"instanceID"`
	Rsp        []byte           `json:"-"`
	Code       ResultStatusCode `json:"code"`
	ErrMsg     string           `json:"errMsg"`
}

type EventConsumer interface {
	HandleEvent(ctx context.Context, handler, eventType string, payload []byte) *EventHandleResult
}

type KVStore interface {
	Get(string) ([]byte, error)
	Set(key string, value []byte) error
}

type SQLStore interface {
	sqlx.SqlExecutor
}

type ContextHandler interface {
	Name() string
	GetImports() ImportsHandler
	SetImports(ImportsHandler)
	GetExports() ExportsHandler
	GetInstance() Instance
	SetInstance(Instance)
}

type ABI interface {
	Log(loglevel, ptr, size int32) int32
	GetData(rid, vmAddrPtr, vmSizePtr int32) int32
	SetData(rid, addr, size int32) int32
	GetDB(kAddr, kSize, vmAddrPtr, vmSizePtr int32) int32
	SetDB(kAddr, kSize, vAddr, vSize int32) int32
	SendTX(chainid int32, offset, size, vmAddrPtr, vmSizePtr int32) int32
	CallContract(chainid int32, offset, size, vmAddrPtr, vmSizePtr int32) int32
	SetSQLDB(addr, size int32) int32
	GetSQLDB(addr, size, vmAddrPtr, vmSizePtr int32) int32
	GetEnv(kAddr, kSize, vmAddrPtr, vmSizePtr int32) int32
	SendMQTT(topicAddr, topicSize, msgAddr, msgSize int32) int32
}

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
