package main

import (
	"fmt"
	"strings"

	common "github.com/machinefi/w3bstream/_examples/wasm_common_go"
)

func main() {}

//export start
func _start(rid uint32) int32 {
	common.Log(fmt.Sprintf("start rid: %d", rid))
	str, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log("error: " + err.Error())
		return -1
	}

	words := strings.Split(string(str), " ")
	records := make(map[string]int32)
	for _, w := range words {
		if _, ok := records[w]; !ok {
			records[w] = common.GetDB(w) + 1
		} else {
			records[w]++
		}
	}

	for k, cnt := range records {
		common.SetDB(k, cnt)
	}
	return 0
}
