package main

import (
	"fmt"
	"strings"

	common "github.com/machinefi/w3bstream/_examples/wasm_common_go"
)

func main() {}

//export start
func _start(rid uint32) int32 {
	common.Log(fmt.Sprintf("start received: %d", rid))
	str, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log("error:" + err.Error())
		return -1
	}

	words := strings.Split(string(str), " ")
	counts := make(map[string]int32)
	for _, w := range words {
		if _, ok := counts[w]; !ok {
			counts[w] = common.GetDB(w) + 1
		} else {
			counts[w]++
		}
	}

	for k, cnt := range counts {
		common.SetDB(k, cnt)
		if _, ok := records[k]; !ok {
			records[k] = cnt
		} else {
			records[k] += cnt
		}
	}
	return 0
}

//export word_count
func _unique(_ uint32) int32 {
	return int32(len(records))
}

var records = make(map[string]int32)
