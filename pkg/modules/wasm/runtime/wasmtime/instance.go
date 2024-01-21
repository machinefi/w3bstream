package wasmtime

import (
	"encoding/binary"
	"sync"
	"sync/atomic"

	"github.com/bytecodealliance/wasmtime-go/v15"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/wasm"
)

func NewWasmtimeInstance(vm *VM, mod *Module) wasm.Instance {
	i := &Instance{
		vm:  vm,
		mod: mod,
	}
	i.stopCond = sync.NewCond(&i.locker)

	return nil
}

var _ wasm.Instance = (*Instance)(nil)

type Instance struct {
	vm   *VM
	mod  *Module
	ins  *wasmtime.Instance
	abis []abi.ABI

	debug    *DwarfInfo
	locker   sync.Mutex
	started  atomic.Bool
	refCount int
	stopCond *sync.Cond

	mem *wasmtime.Memory
	fns sync.Map

	data any
}

var _ = (*Instance)(nil)

func (i *Instance) ID() string {
	return i.vm.id
}

func (i *Instance) Start() error {
	_ = i.abis

	abiNames := i.GetModule().GetABINameList()
	for _, _ = range abiNames {
	}

	// TODO Instantiate
	i.started.Store(true)

	return nil
}

func (i *Instance) Stop() {
	// TODO Deinstantiate
}

func (i *Instance) RegisterImports(name string) error {
	return nil
}

func (i *Instance) Malloc(size int32) (uint64, error) {
	if !i.started.Load() {
		return 0, ErrInstanceNotStart
	}

	mallocFn, err := i.GetExportsFunc("malloc")
	if err != nil {
		return 0, err
	}

	addr, err := mallocFn.Call(size)
	if err != nil {
		i.HandleError(err)
		return 0, err
	}
	return uint64(addr.(uint32)), nil
}

func (i *Instance) GetExportsFunc(name string) (wasm.Function, error) {
	if !i.started.Load() {
		return nil, ErrInstanceNotStart
	}

	if v, ok := i.fns.Load(name); ok {
		return v.(*wasmtimeNativeFunction), nil
	}

	export := i.ins.GetExport(i.vm.store, name)
	if export == nil {
		return nil, ErrInvalidExport
	}

	f := export.Func()
	if f == nil {
		return nil, ErrInvalidFunction
	}
	nf := newWasmtimeNativeFunction(i.vm.store, f)

	i.fns.Store(name, nf)

	return nf, nil
}

func (i *Instance) GetExportsMem(name string) ([]byte, error) {
	if !i.started.Load() {
		return nil, ErrInstanceNotStart
	}

	if i.mem == nil {
		exp := i.ins.GetExport(i.vm.store, name)
		if exp == nil {
			return nil, ErrInvalidExport
		}
		m := exp.Memory()
		if m == nil {
			return nil, ErrInvalidMemory
		}
		i.mem = m
	}

	return i.mem.UnsafeData(i.vm.store), nil
}

func (i *Instance) GetMemory(addr, size uint64) ([]byte, error) {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return nil, err
	}

	if checkIfOverflow(addr, size, mem) {
		return nil, ErrMemAccessOverflow
	}

	return mem[addr : addr+size], nil
}

func (i *Instance) PutMemory(addr, size uint64, mem []byte) error {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if need := uint64(len(mem)); need > size {
		size = need
	}

	if checkIfOverflow(addr, size, mem) {
		return ErrMemAccessOverflow
	}

	copy(mem[addr:], mem[:size])
	return nil
}

func (i *Instance) GetByte(addr uint64) (byte, error) {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}

	if checkIfOverflow(addr, 0, mem) {
		return 0, ErrMemAccessOverflow
	}

	return mem[addr], nil
}

func (i *Instance) PutByte(addr uint64, v byte) error {
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

func (i *Instance) GetUint32(addr uint64) (uint32, error) {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}

	if checkIfOverflow(addr, 4, mem) {
		return 0, ErrMemAccessOverflow
	}

	return binary.LittleEndian.Uint32(mem[addr:]), nil
}

func (i *Instance) PutUint32(addr uint64, v uint32) error {
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

func (i *Instance) GetModule() wasm.Module { return i.mod }

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

	if !i.started.Load() {
		return false
	}

	i.refCount++
	return true
}

func (i *Instance) Release() {
	i.locker.Lock()
	i.refCount--

	if i.refCount <= 0 {
		i.stopCond.Broadcast()
	}
	i.locker.Unlock()
}

func (i *Instance) HandleError(err error) {
	var trapErr *wasmtime.Trap
	if !errors.As(err, &trapErr) {
		return
	}

	frames := trapErr.Frames()
	if frames == nil {
		return
	}

	for _, f := range frames {
		// TODO @sincos output below info for wasm code debugging
		_ = f.FuncIndex()    // funcIndex
		_ = f.FuncOffset()   // funcOffset
		_ = f.ModuleOffset() // moduleOffset PC
		_ = ""               // file name
		_ = 0                // line number
		if i.debug != nil {
			pc := uint64(f.ModuleOffset())
			ln := i.debug.SeekPC(pc)
			if ln != nil {
				_ = ln.File.Name
				_ = ln.Line
			}
		}
	}
}
