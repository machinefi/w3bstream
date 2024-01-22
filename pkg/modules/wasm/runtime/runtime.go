package runtime

import (
	"context"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/proxy"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/runtime/wasmtime"
)

var ErrUnknownRuntimeName = errors.New("unknown runtime name")

func NewRuntime(ctx context.Context, engine string, id string, code []byte) (types.Instance, error) {
	switch engine {
	case "wasmtime":
		vm := wasmtime.NewWasmtimeVM(id)
		mod, err := vm.NewModule(code)
		if err != nil {
			return nil, err
		}
		instance := wasmtime.NewWasmtimeInstance(vm.(*wasmtime.VM), mod.(*wasmtime.Module))
		imports := proxy.NewImports(ctx)
		instance.SetUserdata(&proxy.ABIContext{
			Imports:  imports,
			Instance: instance,
		})
		return instance, nil
	default:
		return nil, ErrUnknownRuntimeName
	}
}
