package types

// VM defines interface of the wasm virtual machine(engine)
type VM interface {
	ID() string
	// Name the name of the wasm vm
	Name() string
	// Init was invoked when creating a new wasm vm
	Init()
	// NewModule compile user code into a wasm module
	NewModule(code []byte) (Module, error)
	// Close to prevent from memory leaking
	Close() error
}

// Module defines interface of the wasm module
type Module interface {
	// Init was invoked when creating a wasm module
	Init()
	// NewInstance instantiates a wasm vm instance
	NewInstance() Instance
	// GetABINameList returns the abi name list exported from wasm module
	GetABINameList() []string
}

// Instance defines the wasm instance
type Instance interface {
	ID() string
	// RegisterImports
	RegisterImports(name string) error

	// Start starts wasm instance
	Start() error
	// Stop stops wasm instance
	Stop()
	// Started returns if instance is started
	Started() bool

	// GetModule get current wasm module
	GetModule() Module
	// GetExportsFunc returns the exported func of the wasm instance by name
	GetExportsFunc(name string) (Function, error)
	// GetExportsMem returns the exported memory of the wasm instance
	GetExportsMem(name string) ([]byte, error)

	// GetByte returns the wasm byte from the specified addr
	GetByte(addr uint64) (byte, error)
	// PutByte set a wasm byte to specified addr
	PutByte(addr uint64, v byte) error
	// GetUint32 returns the u32 value from the specified addr
	GetUint32(addr uint64) (uint32, error)
	// PutUint32 sets u32 value to the specified addr
	PutUint32(addr uint64, v uint32) error
	// GetMemory returns wasm memory
	GetMemory(addr uint64, size uint64) ([]byte, error)
	// PutMemory
	PutMemory(addr uint64, size uint64, data []byte) error
	// Malloc allocates size of memory from wasm default memory
	Malloc(size int32) (uint64, error)

	// GetUserdata get user-defined data
	GetUserdata() any
	// SetUserdata set user-define data into wasm instance
	SetUserdata(any)
	// Lock gets the exclusive ownership of the wasm instance and sets userdata
	Lock(any)
	// Unlock releases the exclusive ownership of the wasm instance and unref userdata
	Unlock()

	// Acquire increases the refer count of wasm instance
	Acquire() bool
	// Release decreases the refer count of wasm instance
	Release()

	Call(string, ...interface{}) (interface{}, error)

	HandleError(error)
}

type Function interface {
	Call(args ...any) (any, error)
}
