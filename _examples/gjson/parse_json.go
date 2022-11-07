package main

import (
	"fmt"

	"github.com/tidwall/gjson"

	common "github.com/machinefi/w3bstream/_examples/wasm_common_go"
)

func main() {}

//export start
func _start(rid uint32) int32 {
	common.Log(fmt.Sprintf("start received: %d", rid))
	message, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log("error: " + err.Error())
		return -1
	}
	res := string(message)
	common.Log("wasm received message: " + res)
	common.Log("wasm get name(json string) from json: " + gjson.Get(res, "name").String())
	common.Log("wasm get name.age(int) from json: " + gjson.Get(res, "name.age").String())
	common.Log("wasm get friends(array) from json: " + gjson.Get(res, "friends").String())
	common.Log("wasm get friends[0].nets(array) from json: " + gjson.Get(res, "friends.0.nets").String())
	return 0
}
