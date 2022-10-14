package wasmtime

import (
	"encoding/binary"

	"github.com/bytecodealliance/wasmtime-go"
	conflog "github.com/iotexproject/Bumblebee/conf/log"
	"github.com/iotexproject/Bumblebee/x/mapx"
	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

type WasmtimeExportFunc interface{}

type ExportFuncs struct {
	store  *wasmtime.Store
	res    *mapx.Map[uint32, []byte]
	db     map[string]int32
	logger conflog.Logger
}

func (ef *ExportFuncs) Log(c *wasmtime.Caller, ptr, size int32) {
	membuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	buf, err := read(membuf, ptr, size)
	if err != nil {
		return
	}
	ef.logger.Info(string(buf))
}

func (ef *ExportFuncs) GetData(c *wasmtime.Caller, rid, vmAddrPtr, vmSizePtr int32) int32 {
	allocFn := c.GetExport("alloc")
	if allocFn == nil {
		return int32(wasm.ResultStatusCode_ImportNotFound)
	}
	data, ok := ef.res.Load(uint32(rid))
	if !ok {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}
	size := len(data)
	result, err := allocFn.Func().Call(ef.store, int32(size))
	if err != nil {
		return int32(wasm.ResultStatusCode_ImportCallFailed)
	}
	addr := result.(int32)

	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	if siz := copy(memBuf[addr:], data); siz != size {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	// fmt.Printf("host >> addr=%d\n", addr)
	// fmt.Printf("host >> size=%d\n", size)
	// fmt.Printf("host >> vmAddrPtr=%d\n", vmAddrPtr)
	// fmt.Printf("host >> vmSizePtr=%d\n", vmSizePtr)

	if err := putUint32Le(memBuf, vmAddrPtr, uint32(addr)); err != nil {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	if err := putUint32Le(memBuf, vmSizePtr, uint32(size)); err != nil {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	// fmt.Println("host >> get_data returned")
	return int32(wasm.ResultStatusCode_OK)
}

// TODO SetData if rid not exist, should be assigned by wasm?
func (ef *ExportFuncs) SetData(c *wasmtime.Caller, rid, addr, size int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	if addr > int32(len(memBuf)) || addr+size > int32(len(memBuf)) {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	buf, err := read(memBuf, addr, size)
	if err != nil {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	ef.res.Store(uint32(rid), buf)
	return int32(wasm.ResultStatusCode_OK)
}

// TODO SetDB value should have type
func (ef *ExportFuncs) SetDB(c *wasmtime.Caller, kAddr, kSize, val int32) {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	key, _ := read(memBuf, kAddr, kSize)

	ef.logger.WithValues(
		"key", string(key),
		"val", val,
	).Info("host.SetDB")

	ef.db[string(key)] = val
}

func (ef *ExportFuncs) GetDB(c *wasmtime.Caller, kAddr, kSize int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	key, err := read(memBuf, kAddr, kSize)
	if err != nil {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	val := ef.db[string(key)]

	ef.logger.WithValues(
		"key", string(key),
		"val", val,
	).Info("host.GetDB")

	return val
}

func putUint32Le(buf []byte, addr int32, num uint32) error {
	if int32(len(buf)) < addr+4 {
		return errors.New("overflow")
	}
	binary.LittleEndian.PutUint32(buf[addr:], num)
	return nil
}

func read(memBuf []byte, addr int32, size int32) ([]byte, error) {
	if addr > int32(len(memBuf)) || addr+size > int32(len(memBuf)) {
		return nil, errors.New("overflow")
	}
	buf := make([]byte, size)
	if siz := copy(buf, memBuf[addr:addr+size]); int32(siz) != size {
		return nil, errors.New("overflow")
	}
	return buf, nil
}
