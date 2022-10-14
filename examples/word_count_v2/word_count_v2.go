package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

func main() {}

//go:wasm-module env
//export ws_log
func _ws_log(ptr uint32, size uint32)

//go:wasm-module env
//export ws_get_data
func _ws_get_data(rid uint32, ptr uint32, size uint32) int32

//go:wasm-module env
//export ws_get_db
func _ws_get_db(kaddr, ksize uint32) (v int32)

//go:wasm-module env
//export ws_set_db
func _ws_set_db(kaddr, ksize uint32, v int32)

var allocs = make(map[uintptr][]byte)

//export alloc
func alloc(size uintptr) unsafe.Pointer {
	buf := make([]byte, size)
	ptr := unsafe.Pointer(&buf[0])
	allocs[uintptr(ptr)] = buf
	return ptr
}

//export count
func _start(rid uint32) int32 {
	log(fmt.Sprintf("wasm >> received: %d", rid))
	str, err := getData(rid)
	if err != nil {
		log("wasm >> error" + err.Error())
		return -1
	}

	words := strings.Split(str, " ")
	counts := make(map[string]int32)
	for _, w := range words {
		if _, ok := counts[w]; !ok {
			kaddr, ksize := stringToPtr(w)
			counts[w] = _ws_get_db(kaddr, ksize) + 1
		} else {
			counts[w]++
		}
	}

	for k, cnt := range counts {
		kaddr, ksize := stringToPtr(k)
		_ws_set_db(kaddr, ksize, cnt)
		if _, ok := records[k]; !ok {
			records[k] = cnt
		} else {
			records[k] += cnt
		}
	}
	return 0
}

//export unique
func _unique(_ uint32) int32 {
	return int32(len(records))
}

var records = make(map[string]int32)

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

func ptrToBytes(ptr uint32, size uint32) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size),
		Cap:  uintptr(size),
	}))
}

func ptrToInt32(ptr uint32, size uint32) (v int32, err error) {
	if size != 4 {
		return 0, errors.New("invalid data size")
	}
	buf := ptrToBytes(ptr, size)
	err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &v)
	return
}

func stringToPtr(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}

func int32ToPtr(v int32) (uint32, uint32) {
	return uint32(uintptr(unsafe.Pointer(&v))), 4
}
