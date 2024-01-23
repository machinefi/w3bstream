package wasmtime

import (
	"github.com/bytecodealliance/wasmtime-go/v15"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/wasm"
)

// NewWasmtimeVM creates wasmtime vm
func NewWasmtimeVM(id string) wasm.VM {
	v := &VM{id: id}
	v.Init()
	return v
}

var _ wasm.VM = (*VM)(nil)

type VM struct {
	id     string
	engine *wasmtime.Engine
	store  *wasmtime.Store
}

func (vm *VM) ID() string {
	return vm.id
}

func (vm *VM) Name() string {
	return "wasmtime"
}

func (vm *VM) Init() {
	vm.engine = wasmtime.NewEngine()
	vm.store = wasmtime.NewStore(vm.engine)
}

func (vm *VM) NewModule(code []byte) (wasm.Module, error) {
	if len(code) == 0 {
		return nil, ErrInvalidWasmCode
	}

	mod, err := wasmtime.NewModule(vm.engine, code)
	if err != nil {
		return nil, errors.Wrap(ErrNewWasmModule, err.Error())
	}

	return NewWasmtimeModule(vm, mod, code)
}

func (vm *VM) Close() error {
	return nil
}
