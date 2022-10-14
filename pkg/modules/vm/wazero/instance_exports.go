package wazero

import (
	"fmt"

	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

func (i *Instance) Log(offset, size uint32) {
	buf, ok := i.mod.Memory().Read(i.ctx, offset, size)
	if !ok {
		panic(fmt.Sprintf("Memory.Read(%d,%d) out of range)", offset, size))
	}
	i.opt.Logger.Info(string(buf))
}

func (i *Instance) GetData(rid uint32, vmAddrPtr, vmSizePtr uint32) wasm.ResultStatusCode {
	if i.alloc == nil {
		return wasm.ResultStatusCode_ImportNotFound
	}

	data, ok := i.res.Load(rid)
	if !ok {
		return wasm.ResultStatusCode_ResourceNotFound
	}
	size := len(data)
	results, err := i.alloc.Call(i.ctx, uint64(size))
	if err != nil {
		return wasm.ResultStatusCode_ImportCallFailed
	}
	addr := results[0]

	if !i.mod.Memory().Write(i.ctx, uint32(addr), data) {
		return wasm.ResultStatusCode_TransDataToVMFailed
	}

	fmt.Printf("host >> addr=%d\n", addr)
	fmt.Printf("host >> size=%d\n", size)
	fmt.Printf("host >> vmAddrPtr=%d\n", vmAddrPtr)
	fmt.Printf("host >> vmSizePtr=%d\n", vmSizePtr)

	if !i.mod.Memory().WriteUint32Le(i.ctx, vmAddrPtr, uint32(addr)) {
		return wasm.ResultStatusCode_TransDataToVMFailed
	}
	if !i.mod.Memory().WriteUint32Le(i.ctx, vmSizePtr, uint32(size)) {
		return wasm.ResultStatusCode_TransDataToVMFailed
	}

	fmt.Println("host >> get_data returned")
	return wasm.ResultStatusCode_OK
}

// TODO SetData if rid not exist, should be assigned by wasm?

func (i *Instance) SetData(rid uint32, addr, size uint32) wasm.ResultStatusCode {
	buf, ok := i.mod.Memory().Read(i.ctx, addr, size)
	if !ok {
		return wasm.ResultStatusCode_TransDataFromVMFailed
	}
	i.res.Store(rid, buf)

	return 0
}

// TODO SetDB value should have type

func (i *Instance) SetDB(kAddr, kSize uint32, val int32) {
	key, ok := i.mod.Memory().Read(i.ctx, kAddr, kSize)
	if !ok {
		return
	}

	i.opt.Logger.WithValues(
		"key", string(key),
		"val", val,
	).Info("host.SetDB")

	i.db[string(key)] = val
}

func (i *Instance) GetDB(kAddr, kSize uint32) int32 {
	key, ok := i.mod.Memory().Read(i.ctx, kAddr, kSize)
	if !ok {
		return 0
	}

	val := i.db[string(key)]

	i.opt.Logger.WithValues(
		"key", string(key),
		"val", val,
	).Info("host.GetDB")

	return val
}
