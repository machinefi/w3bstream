package main

import (
	"fmt"

	"github.com/mailru/easyjson"

	common "github.com/iotexproject/w3bstream/examples/wasm_common_go"

	"github.com/iotexproject/w3bstream/examples/easyjson/model"
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
	common.Log("wasm received: " + string(message))
	student := model.Student{}
	easyjson.Unmarshal(message, &student)
	common.Log("wasm get struct.name from json:" + student.Name)
	common.Log("wasm change student name to Jane ")
	student.Name = "Jane"
	common.Log("wasm get new name from the struct:" + student.Name)
	msg, err := easyjson.Marshal(student)
	common.Log("wasm get json from struct: " + string(msg))
	return 0
}
