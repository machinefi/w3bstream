package main

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/tidwall/gjson"
)

func main() {}

//go:wasm-module env
//export ws_log
func _ws_log(ptr uint32, size uint32)

//go:wasm-module env
//export ws_get_data
func _ws_get_data(rid uint32, ptr uint32, size uint32) int32

var allocs = make(map[uintptr][]byte)

//export alloc
func alloc(size uintptr) unsafe.Pointer {
	buf := make([]byte, size)
	ptr := unsafe.Pointer(&buf[0])
	allocs[uintptr(ptr)] = buf
	return ptr
}

//export start
func _start(rid uint32) int32 {
	log(fmt.Sprintf("wasm >> received: %d", rid))
	message, err := getData(rid)
	if err != nil {
		log("wasm >> error" + err.Error())
		return -1
	}
	log("wasm >> " + message)
	log("wasm get IoTeX from json >> " + gjson.Get(message, "IoTeX").String())
	return 0
}

// log a message to the console using _log.
func log(message string) {
	ptr, size := stringToPtr(message)
	_ws_log(ptr, size)
}

func getData(rid uint32) (string, error) {
	addr := uintptr(unsafe.Pointer(new(uint32)))
	size := uintptr(unsafe.Pointer(new(uint32)))

	code := _ws_get_data(rid, uint32(addr), uint32(size))
	if code != 0 {
		return "", fmt.Errorf("get data failed: [rid:%d] [code:%d]", rid, code)
	}

	log(fmt.Sprintf("wasm >> addr=%d", addr))
	log(fmt.Sprintf("wasm >> size=%d", size))

	vAddr := (*uint32)(unsafe.Pointer(addr))
	vSize := (*uint32)(unsafe.Pointer(size))

	log(fmt.Sprintf("wasm >> *vaddr=%d", *vAddr))
	log(fmt.Sprintf("wasm >> *vsize=%d", *vSize))

	return ptrToString(*vAddr, *vSize), nil
}

func ptrToString(ptr uint32, size uint32) string {
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size),
		Cap:  uintptr(size),
	}))
}

func stringToPtr(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}
