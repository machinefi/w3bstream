package wasmtime

import (
	"github.com/bytecodealliance/wasmtime-go/v15"

	"github.com/machinefi/w3bstream/pkg/modules/wasm"
)

// NewWasmtimeModule
func NewWasmtimeModule(vm *VM, mod *wasmtime.Module, code []byte) (wasm.Module, error) {
	m := &Module{
		vm:   vm,
		mod:  mod,
		code: code,
	}
	m.Init()
	return m, nil
}

type Module struct {
	vm       *VM
	mod      *wasmtime.Module
	abiNames []string
	code     []byte
	debug    *DwarfInfo
}

func (m *Module) Init() {
	m.abiNames = m.GetABINameList()

	if debug := ParseDwarf(m.code); debug != nil {
		m.debug = debug
	}
	m.code = nil
}

func (m *Module) NewInstance() wasm.Instance {
	return nil
}

func (m *Module) GetABINameList() []string {
	exps := m.mod.Exports()
	names := make([]string, 0, len(exps))

	for _, e := range exps {
		if t := e.Type().FuncType(); t != nil {
			names = append(names, e.Name())
		}
	}
	return names
}
