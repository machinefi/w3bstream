package main

import (
	"fmt"

	common "github.com/iotexproject/w3bstream/examples/wasm_common_go"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

//export start
func _start(rid uint32) int32 {
	common.Log(fmt.Sprintf("start rid: %d", rid))

	message, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log(err.Error())
		return -1
	}

	defer func() {
		if common.FreeResource(rid) {
			common.Log(fmt.Sprintf("resource %v released", rid))
		}
	}()

	common.Log(fmt.Sprintf("get resource %v: `%s`", rid, string(message)))
	return 0
}
