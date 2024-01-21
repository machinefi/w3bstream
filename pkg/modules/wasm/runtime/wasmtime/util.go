package wasmtime

import "github.com/bytecodealliance/wasmtime-go/v15"

func checkIfOverflow(addr, size uint64, mem []byte) bool {
	return int(addr) > len(mem) || int(addr+size) > len(mem)
}

func newWasmtimeNativeFunction(store *wasmtime.Store, f *wasmtime.Func) *wasmtimeNativeFunction {
	return &wasmtimeNativeFunction{store: store, f: f}
}

type wasmtimeNativeFunction struct {
	store *wasmtime.Store
	f     *wasmtime.Func
}

func (wf *wasmtimeNativeFunction) Call(args ...any) (any, error) {
	return wf.f.Call(wf.store, args...)
}
