package proxy

import (
	"math"

	"github.com/google/uuid"

	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
)

func init() {
	abi.RegisterABI(ABIName, func(i types.Instance) types.Context {
		return &ABIContext{
			Instance: i,
		}
	})
}

const ABIName = "w3bstream-wasm-proxy"

type ABIContext struct {
	Imports  types.ImportsHandler
	Instance types.Instance
}

func (a *ABIContext) Name() string {
	return ABIName
}

func (a *ABIContext) GetExports() types.Exports {
	return a
}

func (a *ABIContext) GetImports() types.ImportsHandler {
	return a.Imports
}

func (a *ABIContext) SetImports(imports types.ImportsHandler) {
	a.Imports = imports
}

func (a *ABIContext) GetInstance() types.Instance {
	return a.Instance
}

func (a *ABIContext) SetInstance(instance types.Instance) {
	a.Instance = instance
}

func (a *ABIContext) OnEventReceived(entry string, typ string, data []byte) (interface{}, uint64, error) {
	rid := int32(uuid.New().ID() % math.MaxInt32)

	if err := a.Imports.SetResourceData(uint32(rid), data); err != nil {
		return nil, 0, err
	}
	defer a.Imports.RemoveResourceData(uint32(rid))
	if err := a.Imports.SetEventType(uint32(rid), typ); err != nil {
		return nil, 0, err
	}
	defer a.Imports.RemoveEventType(uint32(rid))
	defer a.Imports.FlushLog()

	return a.Instance.Call(entry, int(rid))
}

func (a *ABIContext) OnCreated(i types.Instance) {
	if err := i.RegisterImports(a.Name()); err != nil {
		panic(err)
	}
}

func (a *ABIContext) OnDestroy(i types.Instance) {}
