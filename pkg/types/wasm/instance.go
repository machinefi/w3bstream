package wasm

import "github.com/iotexproject/w3bstream/pkg/enums"

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
	Start() error
	Stop()
	State() enums.InstanceState
	AddResource([]byte) uint32
	GetResource(uint32) ([]byte, bool)
	RmvResource(uint32)
	HandleEvent(fn string, data []byte) ([]byte, ResultStatusCode)
	Get(k string) int32
}

type EventConsumer interface {
	HandleEvent(handler string, payload []byte) ([]byte, ResultStatusCode)
}

type KVStore interface {
	Get(string) int32
}
