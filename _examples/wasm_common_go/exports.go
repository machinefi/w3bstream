package common

//go:wasm-module env
//export ws_log
func _ws_log(logLevel, ptr, size uint32) int32

//go:wasm-module env
//export ws_get_data
func _ws_get_data(rid, ptr, size uint32) int32

//go:wasm-module env
//export ws_set_data
func _ws_set_data(rid, ptr, size uint32) int32

//go:wasm-module env
//export ws_get_db
func _ws_get_db(kaddr, ksize, ptr, size uint32) int32

//go:wasm-module env
//export ws_set_db
func _ws_set_db(kaddr, ksize, vaddr, vsize uint32) int32

//go:wasm-module env
//export ws_send_tx
func _ws_send_tx(kaddr, ksize uint32) (v int32)

//go:wasm-module env
//export ws_get_redis_db
func _ws_get_redis_db(kaddr, ksize, ptr, size uint32) int32

//go:wasm-module env
//export ws_set_redis_db
func _ws_set_redis_db(kaddr, ksize, vaddr, vsize uint32) int32
