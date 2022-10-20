package common

import (
	"fmt"
	"reflect"
	"unsafe"
)

// PointerToString ptr
func PointerToString(ptr uint32, size uint32) string {
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size),
		Cap:  uintptr(size),
	}))
}

func PointerToBytes(ptr uint32, size uint32) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size),
		Cap:  uintptr(size),
	}))
}

func StringToPointer(v string) (addr, size uint32) {
	return BytesToPointer([]byte(v))
}

func BytesToPointer(v []byte) (addr, size uint32) {
	ptr := &v[0]
	pptr := uintptr(unsafe.Pointer(ptr))
	return uint32(pptr), uint32(len(v))
}

func GetDataByRID(rid uint32) ([]byte, error) {
	addr := uintptr(unsafe.Pointer(new(uint32)))
	size := uintptr(unsafe.Pointer(new(uint32)))

	code := _ws_get_data(rid, uint32(addr), uint32(size))
	if code != 0 {
		return nil, fmt.Errorf("get data failed: [rid:%d] [code:%d]", rid, code)
	}

	vaddr := *(*uint32)(unsafe.Pointer(addr))
	m := allocations.GetByAddr(vaddr)
	if m == nil {
		return nil, fmt.Errorf("get data by addr failed: [rid:%d] [addr:%d]", rid, vaddr)
	}

	allocations.AddResourceWithMem(rid, m)
	return m.data, nil
}

func GetDB(key string) int32 {
	addr, size := StringToPointer(key)
	return _ws_get_db(addr, size)
}

func SetDB(key string, v int32) {
	addr, size := StringToPointer(key)
	_ws_set_db(addr, size, v)
}

func SendTx(key string) int32 {
	addr, size := StringToPointer(key)
	return _ws_send_tx(addr, size)
}

func Log(message string) {
	ptr, size := StringToPointer(message)
	_ws_log(ptr, size)
}
