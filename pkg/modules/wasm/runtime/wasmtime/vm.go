package wasmtime

import (
	"github.com/bytecodealliance/wasmtime-go/v8"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
)

// NewWasmtimeVM creates wasmtime vm
func NewWasmtimeVM(id string) types.VM {
	v := &VM{id: id}
	v.Init()
	return v
}

var _ types.VM = (*VM)(nil)

type VM struct {
	id     string
	engine *wasmtime.Engine
	store  *wasmtime.Store
}

func (vm *VM) ID() string { return vm.id }

func (vm *VM) Name() string { return "wasmtime" }

func (vm *VM) Init() {
	vm.engine = wasmtime.NewEngine()
	vm.store = wasmtime.NewStore(vm.engine)
}

func (vm *VM) NewModule(code []byte) (types.Module, error) {
	if len(code) == 0 {
		return nil, ErrInvalidWasmCode
	}

	mod, err := wasmtime.NewModule(vm.engine, code)
	if err != nil {
		return nil, errors.Wrap(ErrFailedToNewWasmModule, err.Error())
	}

	return NewWasmtimeModule(vm, mod, code)
}

func (vm *VM) Close() error {
	// TODO
	return nil
}
