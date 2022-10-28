package main

import (
	"fmt"

	common "github.com/iotexproject/w3bstream/examples/wasm_common_go"
)

func main() {}

//export start
func _start(rid uint32) int32 {
	common.Log(fmt.Sprintf("start received: %d", rid))
	res, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log("error: " + err.Error())
		return -1
	}
	common.Log(fmt.Sprintf("receive message: %s", res))
	count := common.GetDB("clicks") + 1
	if count%5 == 0 {
		/*
			common.SendTx(fmt.Sprintf(
				`{
					"to": "%s",
					"value": "0",
					"data": "40c10f19000000000000000000000000%s0000000000000000000000000000000000000000000000000de0b6b3a7640000"
				}`,
				"0x1ED83F5AD999262eC06Ed8f3B801e108024b3e9c",
				"97186a21fa8e7955c0f154f960d588c3aca44f14",
			))
		*/
		common.Log("send tx")
	}
	common.SetDB("clicks", count)
	return 0
}
