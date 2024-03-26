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

func NewImportFunc(ns string, f interface{}) *ImportFuncInfo {
	return &ImportFuncInfo{
		Namespace: ns,
		Func:      f,
	}
}

type ImportFuncInfo struct {
	Namespace string
	Func      interface{}
}

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

func HostFunctions(i types.Instance) (map[string]*ImportFuncInfo, error) {
	fns := make(map[string]*ImportFuncInfo)
	h, err := newHost(i)
	if err != nil {
		return nil, err
	}

	fns["abort"] = NewImportFunc("env", h.Abort)
	fns["trace"] = NewImportFunc("env", h.Trace)
	fns["seed"] = NewImportFunc("env", h.Seed)
	fns["ws_log"] = NewImportFunc("env", h.Log)
	fns["ws_get_data"] = NewImportFunc("env", h.GetResourceData)
	fns["ws_set_data"] = NewImportFunc("env", h.SetResourceData)
	fns["ws_get_event_type"] = NewImportFunc("env", h.GetEventType)
	fns["ws_get_db"] = NewImportFunc("env", h.GetKVData)
	fns["ws_set_db"] = NewImportFunc("env", h.SetKVData)
	fns["ws_set_sql_db"] = NewImportFunc("env", h.ExecSQL)
	fns["ws_get_sql_db"] = NewImportFunc("env", h.QuerySQL)
	fns["ws_send_tx"] = NewImportFunc("env", h.SendTX)
	fns["ws_send_tx_with_operator"] = NewImportFunc("env", h.SendTXWithOperator)
	fns["ws_call_contract"] = NewImportFunc("env", h.CallContract)
	fns["ws_get_env"] = NewImportFunc("env", h.Env)
	fns["ws_send_mqtt_msg"] = NewImportFunc("env", h.PubMQTT)
	fns["ws_submit_metrics"] = NewImportFunc("stat", h.SubmitMetrics)
	fns["ws_api_call"] = NewImportFunc("env", h.AsyncAPICall)

	return fns, nil
}

type host struct {
	instance types.Instance
	imports  types.ImportsHandler
}

func (h *host) Log(level, msgaddr, msgsize int32) int32 {
	msg, err := h.instance.GetMemory(msgaddr, msgsize)
	if err != nil {
		h.error(errors.Wrap(err, "Log::GetMemory"), "addr", msgaddr, "size", msgsize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("Log::GetMemory", "message", string(msg))

	h.imports.Log(consts.LogLevel(level), string(msg))
	return consts.RESULT_OK.Int32()
}

func (h *host) Env(keyaddr, keysize, varaddrptr, varsizeptr int32) int32 {
	key, err := h.instance.GetMemory(keyaddr, keysize)
	if err != nil {
		h.error(errors.Wrap(err, "Env::GetMemory"), "addr", keyaddr, "size", keysize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("Env::GetMemory", "key", string(key))

	val, ok := h.imports.Env(string(key))
	if !ok {
		h.error(consts.RESULT__ENV_NOT_FOUND, "key", string(key))
		return consts.RESULT__ENV_NOT_FOUND.Int32()
	}

	h.info("Env::GetMemory", "val", val)
	if err = CopyHostDataToWasm(h.instance, []byte(val), varaddrptr, varsizeptr); err != nil {
		h.error(errors.Wrap(err, "Env::CopyHostDataToWasm"), "addr", varaddrptr, "size", varsizeptr)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed", "key", key, "val", val)
	return consts.RESULT_OK.Int32()
}

func (h *host) Abort(msgaddr, filenameaddr, line, col int32) {
	msg, err := ReadStringFromAddr(h.instance, msgaddr)
	if err != nil {
		h.error(errors.Wrap(err, "Abort::ReadStringFromAddr1"), "addr", msgaddr)
		return
	}
	h.info("Abort::ReadStringFromAddr", "message", msg)

	filename, err := ReadStringFromAddr(h.instance, filenameaddr)
	if err != nil {
		h.error(errors.Wrap(err, "Abort::ReadStringFromAddr2"), "addr", filenameaddr)
		return
	}
	h.info("Abort::ReadStringFromAddr", "filename", filename)

	h.imports.Abort(msg, filename, line, col)
}

func (h *host) Trace(msgaddr, _ int32, arr ...float64) {
	msg, err := ReadStringFromAddr(h.instance, msgaddr)
	if err != nil {
		h.error(errors.Wrap(err, "Trace::ReadStringFromAddr"), "addr", msgaddr)
		return
	}
	h.info("Trace::ReadStringFromAddr", "message", msg)

	str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]")
	if len(str) > 0 {
		str = " " + str
	}
	h.imports.Trace(msg, str)
}

func (h *host) Seed() float64 {
	val := h.imports.Seed()
	h.info("succeed", "rand_value", val)
	return val
}

func (h *host) GetResourceData(rid, retaddrptr, retsizeptr int32) int32 {
	data, ok := h.imports.GetResourceData(uint32(rid))
	if !ok {
		h.error(consts.RESULT__RESOURCE_NOT_FOUND, "rid", rid)
		return consts.RESULT__RESOURCE_NOT_FOUND.Int32()
	}

	h.info("GetResourceData", "data", string(data))
	if err := CopyHostDataToWasm(h.instance, data, retaddrptr, retsizeptr); err != nil {
		h.error(errors.Wrap(err, "GetResourceData::CopyHostDataToWasm"), "addr", retaddrptr, "size", retsizeptr)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed", "rid", rid, "data", string(data))
	return consts.RESULT_OK.Int32()
}

func (h *host) SetResourceData(rid, dataaddr, datasize int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(errors.Wrap(err, "SetResourceData::GetMemory"), "addr", dataaddr, "size", datasize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	if err = h.imports.SetResourceData(uint32(rid), data); err != nil {
		h.error(err, "rid", rid)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed", "rid", rid, "data", string(data))
	return consts.RESULT_OK.Int32()
}

func (h *host) GetEventType(rid, retaddrptr, retsizeptr int32) int32 {
	data, ok := h.imports.GetEventType(uint32(rid))
	if !ok {
		h.error(consts.RESULT__RESOURCE_EVENT_NOT_FOUND, "rid", rid)
		return consts.RESULT__RESOURCE_EVENT_NOT_FOUND.Int32()
	}

	h.info("GetEventType", "event_type", data)
	if err := CopyHostDataToWasm(h.instance, []byte(data), retaddrptr, retsizeptr); err != nil {
		h.error(errors.Wrap(err, "GetEventType::CopyHostDataToWasm"), "addr", retaddrptr, "size", retsizeptr)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed", "rid", rid, "event_type", data)
	return consts.RESULT_OK.Int32()
}

func (h *host) GetKVData(keyaddr, keysize, retaddrptr, retsizeptr int32) int32 {
	key, err := h.instance.GetMemory(keyaddr, keysize)
	if err != nil {
		h.error(errors.Wrap(err, "GetKVData::GetMemory"), "addr", keyaddr, "size", keysize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("GetKVData::GetMemory", "key", string(key))

	val, err := h.imports.GetKVData(string(key))
	if err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	if val == nil {
		return consts.RESULT__KV_DATA_NOT_FOUND.Int32()
	}

	h.info("GetKVData", "val", string(val))
	if err = CopyHostDataToWasm(h.instance, val, retaddrptr, retsizeptr); err != nil {
		h.error(errors.Wrap(err, "GetKVData::CopyHostDataToWasm"), "addr", retaddrptr, "size", retsizeptr)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed", "key", string(key), "val", string(val))
	return consts.RESULT_OK.Int32()
}

func (h *host) SetKVData(keyaddr, keysize, dataaddr, datasize int32) int32 {
	key, err := h.instance.GetMemory(keyaddr, keysize)
	if err != nil {
		h.error(errors.Wrap(err, "SetKVData::GetMemory1"), "addr", keyaddr, "size", keysize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("SetKVData1", "key", string(key))

	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(errors.Wrap(err, "SetKVData::GetMemory2"), "addr", dataaddr, "size", datasize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("SetKVData1", "val", string(data))

	if err = h.imports.SetKVData(string(key), data); err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed", "key", string(key), "val", string(data))
	return consts.RESULT_OK.Int32()
}

func (h *host) ExecSQL(queryaddr, querysize int32) int32 {
	query, err := h.instance.GetMemory(queryaddr, querysize)
	if err != nil {
		h.error(errors.Wrap(err, "ExecSQL::GetMemory"), "addr", queryaddr, "size", querysize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("ExecSQL", "query", string(query))

	if err = h.imports.ExecSQL(string(query)); err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) QuerySQL(queryaddr, querysize, retaddrptr, retsizeptr int32) int32 {
	query, err := h.instance.GetMemory(queryaddr, querysize)
	if err != nil {
		h.error(errors.Wrap(err, "QuerySQL::GetMemory"), "addr", queryaddr, "size", querysize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("QuerySQL", "query", string(query))

	res, err := h.imports.QuerySQL(string(query))
	if err != nil {
		h.error(err, "query", string(query))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("QuerySQL", "result", string(res))

	if err = CopyHostDataToWasm(h.instance, res, retaddrptr, retsizeptr); err != nil {
		h.error(errors.Wrap(err, "QuerySQL::CopyHostDataToWasm"), "addr", retaddrptr, "size", retsizeptr)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) SendTX(chainid, dataaddr, datasize, hashaddrptr, hashsizeptr int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(errors.Wrap(err, "SendTX::GetMemory"), "addr", dataaddr, "size", datasize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("SendTX", "data", string(data))

	hash, err := h.imports.SendTX(uint32(chainid), data)
	if err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("SendTX", "hash", hash)

	if err = CopyHostDataToWasm(h.instance, []byte(hash), hashaddrptr, hashsizeptr); err != nil {
		h.error(errors.Wrap(err, "SendTX::CopyHostDataToWasm"), "addr", hashaddrptr, "size", hashsizeptr)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) SendTXWithOperator(chainid, dataaddr, datasize, hashaddrptr, hashsizeptr int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(errors.Wrap(err, "SendTXWithOperator::GetMemory"), "addr", dataaddr, "size", datasize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("SendTXWithOperator", "data", string(data))

	hash, err := h.imports.SendTXWithOperator(uint32(chainid), data)
	if err != nil {
		h.error(err, "chain", chainid, "data", string(data))
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}

	h.info("SendTXWithOperator", "hash", hash)
	if err = CopyHostDataToWasm(h.instance, []byte(hash), hashaddrptr, hashsizeptr); err != nil {
		h.error(errors.Wrap(err, "SendTXWithOperator::CopyHostDataToWasm"), "addr", hashaddrptr, "size", hashsizeptr)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) CallContract(chainid, dataaddr, datasize, resultaddrptr, resultsizeptr int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(errors.Wrap(err, "CallContract::GetMemory"), "addr", dataaddr, "size", datasize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("CallContract", "data", string(data))

	res, err := h.imports.CallContract(uint32(chainid), data)
	if err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("CallContract", "result", res)

	if err = CopyHostDataToWasm(h.instance, res, resultaddrptr, resultsizeptr); err != nil {
		h.error(errors.Wrap(err, "CallContract::CopyHostDataToWasm"), "addr", resultaddrptr, "size", resultsizeptr)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) PubMQTT(topicaddr, topicsize, msgaddr, msgsize int32) int32 {
	topic, err := h.instance.GetMemory(topicaddr, topicsize)
	if err != nil {
		h.error(errors.Wrap(err, "PutMQTT::GetMemory1"), "addr", topicaddr, "size", topicsize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("PubMQTT", "topic", string(topic))

	msg, err := h.instance.GetMemory(msgaddr, msgsize)
	if err != nil {
		h.error(errors.Wrap(err, "PutMQTT::GetMemory2"), "addr", msgaddr, "size", msgsize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("PubMQTT", "message", string(msg))

	if err = h.imports.PubMQTT(string(topic), msg); err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) SubmitMetrics(dataaddr, datasize int32) int32 {
	data, err := h.instance.GetMemory(dataaddr, datasize)
	if err != nil {
		h.error(errors.Wrap(err, "SubmitMetrics::GetMemory"), "dataaddr", dataaddr, "datasize", datasize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("SubmitMetrics", "data", string(data))

	if err = h.imports.SubmitMetrics(data); err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("succeed")
	return consts.RESULT_OK.Int32()
}

func (h *host) AsyncAPICall(reqaddr, reqsize, rspaddrptr, rspsizeptr int32) int32 {
	req, err := h.instance.GetMemory(reqaddr, reqsize)
	if err != nil {
		h.error(errors.Wrap(err, "AsyncAPICall::GetMemory"), "reqaddr", reqaddr, "reqsize", reqsize)
		return consts.RESULT__INVALID_MEM_ACCESS.Int32()
	}
	h.info("AsyncAPICall", "req", string(req))

	rsp, err := h.imports.AsyncAPICall(req)
	if err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED.Int32()
	}
	h.info("AsyncAPICall", "rsp", string(rsp))

	if err = CopyHostDataToWasm(h.instance, rsp, rspaddrptr, rspsizeptr); err != nil {
		h.error(errors.Wrap(err, "AsyncAPICall::CopyHostDataToWasm"), "addr", rspaddrptr, "size", rspsizeptr)
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

func (h *host) error(err error, args ...any) {
	h._log(consts.LOG_LEVEL__ERROR, err.Error(), args...)
}
