package host

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/consts"
)

var (
	errFailedToGetImportsHandler = errors.New("failed to get imports handler")
)

func newHost(i types.Instance) (*host, error) {
	imports := GetImportsHandler(i)
	if imports == nil {
		return nil, errFailedToGetImportsHandler
	}

	return &host{
		instance: i,
		imports:  imports,
	}, nil
}

func HostFunctions(i types.Instance) (map[string]interface{}, error) {
	fns := make(map[string]interface{})
	h, err := newHost(i)
	if err != nil {
		return nil, err
	}

	fns["abort"] = h.Abort
	fns["trace"] = h.Trace
	fns["seed"] = h.Seed
	fns["ws_log"] = h.Log
	fns["ws_get_data"] = h.GetResourceData
	fns["ws_set_data"] = h.SetResourceData
	fns["ws_get_event_type"] = h.GetEventType
	fns["ws_get_db"] = h.GetKVData
	fns["ws_set_db"] = h.SetKVData
	fns["ws_set_sql_db"] = h.ExecSQL
	fns["ws_get_sql_db"] = h.QuerySQL
	fns["ws_send_tx"] = h.SendTX
	fns["ws_send_tx_with_operator"] = h.SendTXWithOperator
	fns["ws_call_contract"] = h.CallContract
	fns["ws_get_env"] = h.Env
	fns["ws_send_mqtt_msg"] = h.PubMQTT
	fns["ws_submit_metrics"] = h.SubmitMetrics
	fns["ws_api_call"] = h.AsyncAPICall

	return fns, nil
}

type host struct {
	instance types.Instance
	imports  types.ImportsHandler
}

func (h *host) Log(level, msgaddr, msgsize int32) int32 {
	msg, err := h.instance.GetMemory(msgaddr, msgsize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}

	h.imports.Log(consts.LogLevel(level), string(msg))
	return consts.RESULT_OK.Int32()
}

func (h *host) Env(keyaddr, keysize, varaddrptr, varsizeptr int32) int32 {
	key, err := h.instance.GetMemory(keyaddr, keysize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	val, ok := h.imports.Env(string(key))
	if !ok {
		h.error(consts.RESULT__ENV_NOT_FOUND, "key", string(key))
		return consts.RESULT__ENV_NOT_FOUND.Int32()
	}
	if err = CopyHostDataToWasm(h.instance, []byte(val), varaddrptr, varsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) Abort(msgaddr, filenameaddr, line, col int32) {
	msg, err := ReadStringFromAddr(h.instance, msgaddr)
	if err != nil {
		h.error(err)
		return
	}
	filename, err := ReadStringFromAddr(h.instance, filenameaddr)
	if err != nil {
		h.error(err)
		return
	}
	h.imports.Abort(msg, filename, line, col)
}

func (h *host) Trace(msgaddr, _ int32, arr ...float64) {
	msg, err := ReadStringFromAddr(h.instance, msgaddr)
	if err != nil {
		h.error(err)
		return
	}

	str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]")
	if len(str) > 0 {
		str = " " + str
	}
	h.imports.Trace(msg, str)
}

func (h *host) Seed() float64 {
	h.info("succeed")
	return h.imports.Seed()
}

func (h *host) GetResourceData(rid, retaddrptr, retsizeptr int32) int32 {
	data, ok := h.imports.GetResourceData(uint32(rid))
	if !ok {
		h.error(consts.RESULT__RESOURCE_NOT_FOUND, "rid", rid)
		return consts.RESULT__RESOURCE_NOT_FOUND.Int32()
	}

	if err := CopyHostDataToWasm(h.instance, data, retaddrptr, retsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) SetResourceData(rid, dataaddr, datasize int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	if err = h.imports.SetResourceData(uint32(rid), data); err != nil {
		h.error(err, "rid", rid)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) GetEventType(rid, retaddrptr, retsizeptr int32) int32 {
	data, ok := h.imports.GetEventType(uint32(rid))
	if !ok {
		h.error(consts.RESULT__RESOURCE_EVENT_NOT_FOUND, "rid", rid)
		return consts.RESULT__RESOURCE_EVENT_NOT_FOUND.Int32()
	}

	if err := CopyHostDataToWasm(h.instance, []byte(data), retaddrptr, retsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) GetKVData(keyaddr, keysize, retaddrptr, retsizeptr int32) int32 {
	key, err := h.instance.GetMemory(keyaddr, keysize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	val, err := h.imports.GetKVData(string(key))
	if err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	if val == nil {
		return consts.RESULT__KV_DATA_NOT_FOUND.Int32()
	}
	if err = CopyHostDataToWasm(h.instance, val, retaddrptr, retsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) SetKVData(keyaddr, keysize, dataaddr, datasize int32) int32 {
	key, err := h.instance.GetMemory(keyaddr, keysize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}

	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}

	if err = h.imports.SetKVData(string(key), data); err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) ExecSQL(queryaddr, querysize int32) int32 {
	query, err := h.instance.GetMemory(queryaddr, querysize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}

	if err = h.imports.ExecSQL(string(query)); err != nil {
		h.error(err, "query", string(query))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) QuerySQL(queryaddr, querysize, retaddrptr, retsizeptr int32) int32 {
	query, err := h.instance.GetMemory(queryaddr, querysize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	res, err := h.imports.QuerySQL(string(query))
	if err != nil {
		h.error(err, "query", string(query))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	if err = CopyHostDataToWasm(h.instance, res, retaddrptr, retsizeptr); err != nil {
		h.error(err, "result", string(res))
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) SendTX(chainid, dataaddr, datasize, hashaddrptr, hashsizeptr int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	hash, err := h.imports.SendTX(uint32(chainid), data)
	if err != nil {
		h.error(err, "chain", chainid, "data", string(data))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	if err = CopyHostDataToWasm(h.instance, []byte(hash), hashaddrptr, hashsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) SendTXWithOperator(chainid, dataaddr, datasize, hashaddrptr, hashsizeptr int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	hash, err := h.imports.SendTXWithOperator(uint32(chainid), data)
	if err != nil {
		h.error(err, "chain", chainid, "data", string(data))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	if err = CopyHostDataToWasm(h.instance, []byte(hash), hashaddrptr, hashsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) CallContract(chainid, dataaddr, datasize, resultaddrptr, resultsizeptr int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	res, err := h.imports.CallContract(uint32(chainid), data)
	if err != nil {
		h.error(err, "chain", chainid, "data", string(data))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	if err = CopyHostDataToWasm(h.instance, res, resultaddrptr, resultsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) PubMQTT(topicaddr, topicsize, msgaddr, msgsize int32) int32 {
	topic, err := h.instance.GetMemory(topicaddr, topicsize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	msg, err := h.instance.GetMemory(msgaddr, msgsize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}

	if err = h.imports.PubMQTT(string(topic), msg); err != nil {
		h.error(err, "topic", string(topic), "message", string(msg))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) SubmitMetrics(dataaddr, datasize int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	if err = h.imports.SubmitMetrics(data); err != nil {
		h.error(err, "data", string(data))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) AsyncAPICall(reqaddr, reqsize, rspaddrptr, rspsizeptr int32) int32 {
	req, err := h.instance.GetMemory(reqaddr, reqsize)
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	rsp, err := h.imports.AsyncAPICall(req)
	if err != nil {
		h.error(err, "req", string(req))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	if err = CopyHostDataToWasm(h.instance, rsp, rspaddrptr, rspsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) _log(lv consts.LogLevel, msg string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	fs := runtime.CallersFrames([]uintptr{pcs[0]})
	f, _ := fs.Next()
	parts := strings.Split(f.Function, ".")
	fn := parts[len(parts)-1]
	h.imports.LogInternal(lv, msg, append(args, "host_func", fn, "instance_id", h.instance.ID())...)
}

func (h *host) info(msg string, args ...any) {
	h._log(consts.LOG_LEVEL__INFO, msg, args...)
}

func (h *host) Info(msg string, args ...any) {
	h._log(consts.LOG_LEVEL__INFO, msg, args...)
}

func (h *host) error(err error, args ...any) {
	h._log(consts.LOG_LEVEL__INFO, err.Error(), args...)
}
