package wasmtime

import (
	"context"
	"encoding/binary"
	"log/slog"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/bytecodealliance/wasmtime-go/v8"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/proxy"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/host"
)

func NewWasmtimeInstance(vm *VM, mod *Module) types.Instance {
	i := &Instance{
		vm:  vm,
		mod: mod,
	}
	i.stopCond = sync.NewCond(&i.locker)

	return i
}

type Instance struct {
	vm  *VM
	mod *Module
	ins *wasmtime.Instance
	lnk *wasmtime.Linker

	externs  []wasmtime.AsExtern
	debug    *DwarfInfo
	locker   sync.Mutex
	started  atomic.Bool
	refCount int
	stopCond *sync.Cond

	mem *wasmtime.Memory
	fns sync.Map

	data any
}

var _ types.Instance = (*Instance)(nil)

func (i *Instance) ID() string {
	return i.vm.id
}

func (i *Instance) register(namespace, fnName string, fn interface{}) error {
	if namespace == "" || fnName == "" {
		return ErrInvalidImportFunc
	}

	if fn == nil || reflect.ValueOf(fn).IsNil() || reflect.TypeOf(fn).Kind() != reflect.Func {
		return ErrInvalidImportFunc
	}

	return i.lnk.FuncWrap(namespace, fnName, fn)

	// fnType := reflect.TypeOf(fn)

	// argsNum := fnType.NumIn()
	// argKinds := make([]*wasmtime.ValType, argsNum)
	// for i := 0; i < argsNum; i++ {
	// 	argKinds[i] = convertFromGoType(fnType.In(i))
	// }

	// retsNum := fnType.NumOut()
	// retKinds := make([]*wasmtime.ValType, retsNum)
	// for i := 0; i < retsNum; i++ {
	// 	retKinds[i] = convertFromGoType(fnType.Out(i))
	// }

	// return wasmtime.NewFunc(
	// 	i.vm.store,
	// 	wasmtime.NewFuncType(argKinds, retKinds),
	// 	func(caller *wasmtime.Caller, args []wasmtime.Val) (rets []wasmtime.Val, trap *wasmtime.Trap) {
	// 		if len(args) != len(argKinds) {
	// 			return nil, wasmtime.NewTrap("wasmtime: unmatched input number of arguments")
	// 		}

	// 		for i := range args {
	// 			if args[i].Kind() != argKinds[i].Kind() {
	// 				return nil, wasmtime.NewTrap(fmt.Sprintf("wasmtime: unmatched input type of argument: %d", i))
	// 			}
	// 		}

	// 		_args := make([]reflect.Value, len(args))
	// 		for i := range args {
	// 			_args[i] = convertToGoTypes(args[i])
	// 		}

	// 		defer func() {
	// 			if r := recover(); r != nil {
	// 				trap = wasmtime.NewTrap(fmt.Sprintf("wasmtime: call %s paniced, r: %v stack: %v", fnName, r, string(debug.Stack())))
	// 				rets = nil
	// 			}
	// 		}()

	// 		_rets := reflect.ValueOf(fn).Call(_args)
	// 		rets = make([]wasmtime.Val, len(_rets))
	// 		for i := range _rets {
	// 			rets[i] = convertToWasmtimeVal(_rets[i])
	// 		}
	// 		return rets, nil

	// 		// fn := caller.GetExport(fnName).Func()
	// 		// result, err := fn.Call(i.vm.store, _args...)
	// 		// if err != nil {
	// 		// 	return nil, wasmtime.NewTrap(err.Error())
	// 		// }
	// 		// if result == nil {
	// 		// 	return nil, nil
	// 		// }
	// 		// if v, ok := result.([]wasmtime.Val); ok {
	// 		// 	return v, nil
	// 		// }
	// 		// return []wasmtime.Val{convertToWasmtimeVal(result)}, nil
	// 	},
	// ), nil
}

func (i *Instance) RegisterImports(name string) error {
	if name != proxy.ABIName {
		return errors.Wrap(ErrUnknownABIName, name)
	}

	hostFns, err := host.HostFunctions(i)
	if err != nil {
		return err
	}

	for fnName, fn := range hostFns {
		if err = i.register("env", fnName, fn); err != nil {
			return err
		}
	}
	return nil
}

func (i *Instance) Start() error {
	i.lnk = wasmtime.NewLinker(i.vm.engine)
	if err := i.lnk.DefineWasi(); err != nil {
		slog.Error(err.Error())
		return err
	}

	err := i.RegisterImports(proxy.ABIName)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	i.ins, err = i.lnk.Instantiate(i.vm.store, i.mod.mod)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	i.started.Store(true)
	return nil
}

func (i *Instance) Stop() {
	i.locker.Lock()
	defer i.locker.Unlock()
	for i.refCount > 0 {
		i.stopCond.Wait()
	}
	if i.started.CompareAndSwap(true, false) {
		// TODO destroy
	}
}

func (i *Instance) Started() bool {
	return i.started.Load()
}

func (i *Instance) Malloc(size int32) (int32, error) {
	if !i.Started() {
		return 0, ErrInstanceNotStarted
	}

	// alloc func was implemented in w3bstream-golang-sdk
	fn, err := i.GetExportsFunc("alloc")
	if err != nil {
		fn, err = i.GetExportsFunc("malloc")
		if err != nil {
			return 0, err
		}
	}

	addr, err := fn.Call(size)
	if err != nil {
		i.HandleError(err)
		return 0, err
	}
	return addr.(int32), nil
}

func (i *Instance) GetExportsFunc(name string) (types.Function, error) {
	if !i.Started() {
		return nil, ErrInstanceNotStarted
	}

	if v, ok := i.fns.Load(name); ok {
		return v.(*wasmtimeNativeFunction), nil
	}

	export := i.ins.GetExport(i.vm.store, name)
	if export == nil {
		return nil, errors.Wrap(ErrInvalidExportFunc, name)
	}

	f := export.Func()
	if f == nil {
		return nil, errors.Wrap(ErrInvalidExportFunc, name)
	}
	nf := newWasmtimeNativeFunction(i.vm.store, f)

	i.fns.Store(name, nf)

	return nf, nil
}

func (i *Instance) GetExportsMem(name string) ([]byte, error) {
	if !i.Started() {
		return nil, ErrInstanceNotStarted
	}

	if i.mem == nil {
		exp := i.ins.GetExport(i.vm.store, name)
		if exp == nil {
			return nil, errors.Wrap(ErrInvalidExportMem, name)
		}
		m := exp.Memory()
		if m == nil {
			return nil, errors.Wrap(ErrInvalidExportMem, name)
		}
		i.mem = m
	}

	return i.mem.UnsafeData(i.vm.store), nil
}

func (i *Instance) GetMemory(addr, size int32) ([]byte, error) {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return nil, err
	}

	if checkIfOverflow(addr, size, mem) {
		return nil, ErrMemAccessOverflow
	}

	return mem[addr : addr+size], nil
}

func (i *Instance) PutMemory(addr, size int32, data []byte) error {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if need := int32(len(data)); need > size {
		size = need
	}

	if checkIfOverflow(addr, size, mem) {
		return ErrMemAccessOverflow
	}

	copy(mem[addr:], data[:size])
	return nil
}

func (i *Instance) GetByte(addr int32) (byte, error) {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}

	if checkIfOverflow(addr, 0, mem) {
		return 0, ErrMemAccessOverflow
	}

	return mem[addr], nil
}

func (i *Instance) PutByte(addr int32, v byte) error {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if checkIfOverflow(addr, 0, mem) {
		return ErrMemAccessOverflow
	}

	mem[addr] = v
	return nil
}

func (i *Instance) GetUint32(addr int32) (uint32, error) {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}

	if checkIfOverflow(addr, 4, mem) {
		return 0, ErrMemAccessOverflow
	}

	return binary.LittleEndian.Uint32(mem[addr:]), nil
}

func (i *Instance) PutUint32(addr int32, v uint32) error {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if checkIfOverflow(addr, 4, mem) {
		return ErrMemAccessOverflow
	}

	binary.LittleEndian.PutUint32(mem[addr:], v)
	return nil
}

func (i *Instance) GetModule() types.Module {
	return i.mod
}

func (i *Instance) GetUserdata() any {
	return i.data
}

func (i *Instance) SetUserdata(data any) {
	i.data = data
}

func (i *Instance) Lock(data any) {
	i.locker.Lock()
	i.data = data
}

func (i *Instance) Unlock() {
	i.locker.Unlock()
	i.data = nil
}

func (i *Instance) Acquire() bool {
	i.locker.Lock()
	defer i.locker.Unlock()

	if !i.Started() {
		return false
	}

	i.refCount++
	return true
}

func (i *Instance) Release() {
	i.locker.Lock()
	defer i.locker.Unlock()
	i.refCount--

	if i.refCount <= 0 {
		i.stopCond.Broadcast()
	}
	i.vm.store.GC()
}

func (i *Instance) Call(name string, args ...interface{}) (interface{}, error) {
	if !i.Started() {
		return nil, ErrInstanceNotStarted
	}

	// if v, ok := i.fns.Load(name); ok {
	// 	return v.(*wasmtimeNativeFunction), nil
	// }

	// export := i.ins.GetExport(i.vm.store, name)
	// if export == nil {
	// 	return nil, errors.Wrap(ErrInvalidExportFunc, name)
	// }

	// f := export.Func()
	// if f == nil {
	// 	return nil, errors.Wrap(ErrInvalidExportFunc, name)
	// }

	f := i.ins.GetFunc(i.vm.store, name)
	if f == nil {
		return nil, errors.Wrap(ErrInvalidExportFunc, name)
	}

	ret, err := f.Call(i.vm.store, args...)
	if err != nil {
		i.HandleError(err)
		return nil, err
	}
	return ret, nil
}

func (i *Instance) HandleError(err error) {
	// if i.debug == nil {
	// 	return
	// }
	var trapErr *wasmtime.Trap
	if !errors.As(err, &trapErr) {
		return
	}

	frames := trapErr.Frames()
	if frames == nil {
		return
	}

	for _, f := range frames {
		args := []any{
			"func_index", f.FuncIndex(),
			"func_offset", f.FuncOffset(),
			"instance_id", i.vm.id,
		}
		pc := uint64(f.ModuleOffset())
		if i.debug != nil {
			if l := i.debug.SeekPC(pc); l != nil {
				args = append(args,
					"filename", l.File.Name,
					"line", l.Line,
				)
			}
		}
		slog.Log(context.Background(), slog.LevelError, err.Error(), args...)
	}
}
