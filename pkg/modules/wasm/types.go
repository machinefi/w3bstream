package runtime

type WasmVM interface {
	ID() string
	Name() string
	Init()
	NewModule(code []byte) (WasmModule, error)
	Close() error
}

type WasmModule interface {
	Init()
	NewInstance() WasmInstance
	GetABINameList() []string
}

type WasmInstance interface {
	RegisterImports(name string) error

	Start() error
	Stop()

	GetExportsFunc(name string) (WasmFunction, error)
	GetExportsMem(name string) ([]byte, error)

	GetMemory(addr uint64, size uint64) ([]byte, error)
	PutMemory(addr uint64, size uint64, data []byte) error
	Malloc(size int32) (uint64, error)

	HandleError(error)
}

type WasmFunction interface {
	Call(args ...any) (any, error)
}

type ABI interface {
	Name() string
	GetABIImports() any
	SetABIImports(imps any)
	GetABIExports() any
	// ABIHandler
}
