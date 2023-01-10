package wasmtime

import (
	"encoding/binary"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/pkg/errors"
)

type (
	Runtime struct {
		store    *wasmtime.Store
		instance *wasmtime.Instance
	}
)

func NewRuntime() *Runtime {
	return &Runtime{}
}

func (rt *Runtime) Initiate(lk ABILinker, code []byte) error {
	engine := wasmtime.NewEngineWithConfig(wasmtime.NewConfig())
	module, err := wasmtime.NewModule(engine, code)
	if err != nil {
		return err
	}
	linker := wasmtime.NewLinker(engine)
	if err = lk.LinkABI(func(module, name string, fn interface{}) error {
		return linker.FuncWrap(module, name, fn)
	}); err != nil {
		return err
	}

	_ = linker.DefineWasi()
	rt.store = wasmtime.NewStore(engine)
	rt.store.SetWasi(wasmtime.NewWasiConfig())
	rt.instance, err = linker.Instantiate(rt.store, module)

	return err
}

func (rt *Runtime) newMemory() []byte {
	return rt.instance.GetExport(rt.store, "memory").Memory().UnsafeData(rt.store)
}

func (rt *Runtime) alloc(size int32) (int32, []byte, error) {
	fn := rt.instance.GetExport(rt.store, "alloc")
	if fn == nil {
		return 0, nil, errors.New("alloc is nil")
	}
	result, err := fn.Func().Call(rt.store, size)
	if err != nil {
		return 0, nil, err
	}
	return result.(int32), rt.newMemory(), nil
}

func (rt *Runtime) Call(name string, args ...interface{}) (interface{}, error) {
	fn := rt.instance.GetFunc(rt.store, name)
	if fn == nil {
		return nil, errors.Errorf("runtime: %s fn is not imported", name)
	}
	return fn.Call(rt.store, args...)
}

func (rt *Runtime) Read(addr, size int32) ([]byte, error) {
	mem := rt.newMemory()
	if addr > int32(len(mem)) || addr+size > int32(len(mem)) {
		return nil, errors.New("overflow")
	}
	buf := make([]byte, size)
	if copied := copy(buf, mem[addr:addr+size]); int32(copied) != size {
		return nil, errors.New("overflow")
	}
	return buf, nil
}

func (rt *Runtime) Copy(hostData []byte, vmAddrPtr, vmSizePtr int32) error {
	size := len(hostData)
	addr, mem, err := rt.alloc(int32(size))
	if err != nil {
		return err
	}
	if copied := copy(mem[addr:], hostData); copied != size {
		return errors.New("fail to copy data")
	}
	if err = rt.PutUint32Le(mem, vmAddrPtr, uint32(addr)); err != nil {
		return err
	}
	if err = rt.PutUint32Le(mem, vmSizePtr, uint32(size)); err != nil {
		return err
	}

	return nil
}

func (rt *Runtime) PutUint32Le(buf []byte, vmAddr int32, val uint32) error {
	if int32(len(buf)) < vmAddr+4 {
		return errors.New("overflow")
	}
	binary.LittleEndian.PutUint32(buf[vmAddr:], val)
	return nil
}
