package wasmtime

import (
	"context"
	"encoding/binary"

	"github.com/bytecodealliance/wasmtime-go/v17"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/types"
)

var (
	ErrNotLinked           = errors.New("not linked")
	ErrAlreadyInstantiated = errors.New("already instantiated")
	ErrNotInstantiated     = errors.New("not instantiated")
	ErrFuncNotImported     = errors.New("func not imported")
	ErrAlreadyLinked       = errors.New("already linked")
)

type (
	Runtime struct {
		module   *wasmtime.Module
		linker   *wasmtime.Linker
		store    *wasmtime.Store
		instance *wasmtime.Instance
		engine   *wasmtime.Engine
	}
)

func NewRuntime() *Runtime {
	config := wasmtime.NewConfig()
	config.SetConsumeFuel(true)
	engine := wasmtime.NewEngineWithConfig(config)
	return &Runtime{
		engine: engine,
	}
}

func (rt *Runtime) Link(lk ABILinker, code []byte) error {
	if rt.module != nil {
		return ErrAlreadyLinked
	}
	linker := wasmtime.NewLinker(rt.engine)
	if err := lk.LinkABI(func(module, name string, fn interface{}) error {
		return linker.FuncWrap(module, name, fn)
	}); err != nil {
		return err
	}
	if err := linker.DefineWasi(); err != nil {
		return err
	}
	rt.linker = linker
	module, err := wasmtime.NewModule(rt.engine, code)
	if err != nil {
		return err
	}
	rt.module = module
	return nil
}

func (rt *Runtime) Instantiate(ctx context.Context) error {
	ctx, l := logr.Start(ctx, "modules.vm.wasmtime.Runtime.Instantiate")
	defer l.End()

	if rt.module == nil {
		return ErrNotLinked
	}
	if rt.instance != nil {
		return ErrAlreadyInstantiated
	}
	store := wasmtime.NewStore(rt.engine)
	store.SetWasi(wasmtime.NewWasiConfig())
	if fuel, _ := types.MaxWasmConsumeFuelFromContext(ctx); fuel > 0 {
		if err := store.SetFuel(fuel); err != nil {
			return err
		}
	}

	instance, err := rt.linker.Instantiate(store, rt.module)
	if err != nil {
		return err
	}
	rt.instance = instance
	rt.store = store

	return nil
}

func (rt *Runtime) Deinstantiate(ctx context.Context) {
	ctx, l := logr.Start(ctx, "modules.vm.wasmtime.Runtime.Deinstantiate")
	defer l.End()

	rt.instance = nil
	rt.store = nil
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

func putUint32Le(buf []byte, vmAddr int32, val uint32) error {
	if int32(len(buf)) < vmAddr+4 {
		return errors.New("overflow")
	}
	binary.LittleEndian.PutUint32(buf[vmAddr:], val)
	return nil
}

func (rt *Runtime) Call(ctx context.Context, name string, args ...interface{}) (interface{}, error) {
	ctx, l := logr.Start(ctx, "modules.vm.wasmtime.Runtime.Call", "func", name)
	defer l.End()

	if rt.module == nil {
		return nil, ErrNotLinked
	}
	if rt.instance == nil {
		return nil, ErrNotInstantiated
	}
	fn := rt.instance.GetFunc(rt.store, name)
	if fn == nil {
		return nil, ErrFuncNotImported
	}
	return fn.Call(rt.store, args...)
}

func (rt *Runtime) Read(addr, size int32) ([]byte, error) {
	if rt.module == nil {
		return nil, ErrNotLinked
	}
	if rt.instance == nil {
		return nil, ErrNotInstantiated
	}
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
	if rt.module == nil {
		return ErrNotLinked
	}
	if rt.instance == nil {
		return ErrNotInstantiated
	}
	size := len(hostData)
	addr, mem, err := rt.alloc(int32(size))
	if err != nil {
		return err
	}
	if copied := copy(mem[addr:], hostData); copied != size {
		return errors.New("fail to copy data")
	}
	if err = putUint32Le(mem, vmAddrPtr, uint32(addr)); err != nil {
		return err
	}

	return putUint32Le(mem, vmSizePtr, uint32(size))
}
