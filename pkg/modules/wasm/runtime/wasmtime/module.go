package wasmtime

import (
	"github.com/bytecodealliance/wasmtime-go/v8"

	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
)

// NewWasmtimeModule
func NewWasmtimeModule(vm *VM, mod *wasmtime.Module, code []byte) (*Module, error) {
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
	// debug    *DwarfInfo
}

func (m *Module) Init() {
	m.abiNames = m.GetABINameList()

	// if debug := ParseDwarf(m.code); debug != nil {
	// 	m.debug = debug
	// }
	m.code = nil
}

func (m *Module) NewInstance() types.Instance {
	return NewWasmtimeInstance(m.vm, m)
}

func (m *Module) GetABINameList() []string {
	exps := m.mod.Exports()
	names := make([]string, 0, len(exps))

	for _, e := range exps {
		if t := e.Type().FuncType(); t != nil {
			// if strings.HasPrefix(e.Name(), "ws_") {
			names = append(names, e.Name())
			// }
		}
	}
	return names
}
