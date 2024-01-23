package internal

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/wasm"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/consts"
)

type (
	Instance       = wasm.Instance
	ImportsHandler = wasm.ImportsHandler
)

func newHost(i Instance) (*host, error) {
	imports := GetImportsHandler(i)
	if imports == nil {
		return nil, nil
	}

	return &host{
		instance: i,
		imports:  imports,
		_rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func HostFunctions(i Instance) map[string]interface{} {
	funcs := make(map[string]interface{})
	h, err := newHost(i)
	if err != nil {
		return nil
	}

	funcs["abort"] = h.Abort
	funcs["trace"] = h.Trace
	funcs["seed"] = h.Seed
	funcs["ws_log"] = h.Log
	funcs["ws_get_data"] = h.GetResourceData
	funcs["ws_set_data"] = h.SetResourceData
	funcs["ws_get_db"] = h.GetKVData
	funcs["ws_set_db"] = h.SetKVData
	funcs["ws_send_tx"] = h.SendTX
	funcs["ws_send_tx_with_operator"] = h.SendTXWithOperator
	funcs["ws_call_contract"] = h.CallContract
	funcs["ws_set_sql_db"] = h.ExecSQL
	funcs["ws_get_sql_db"] = h.QuerySQL
	funcs["ws_get_env"] = h.Env
	funcs["ws_send_mqtt_msg"] = h.PubMQTT
	funcs["ws_api_call"] = h.AsyncAPICall

	return funcs
}

type host struct {
	instance Instance
	imports  ImportsHandler
	_rand    *rand.Rand
}

func (h *host) Log(level consts.LogLevel, msgaddr, msgsize int32) consts.Result {
	msg, err := h.instance.GetMemory(uint64(msgaddr), uint64(msgsize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}

	h.imports.Log(level, string(msg))
	return consts.RESULT_OK
}

func (h *host) Env(keyaddr, keysize, varaddrptr, varsizeptr int32) consts.Result {
	key, err := h.instance.GetMemory(uint64(keyaddr), uint64(keysize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	val, ok := h.imports.Env(string(key))
	if !ok {
		h.error(errors.Wrap(consts.RESULT__ENV_NOT_FOUND, string(key)))
		return consts.RESULT__ENV_NOT_FOUND
	}
	if err = CopyDataToInstance(h.instance, []byte(val), varaddrptr, varsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	return consts.RESULT_OK
}

func (h *host) GetResourceData(rid uint32, retaddrptr, retsizeptr int32) consts.Result {
	data, ok := h.imports.GetResourceData(rid)
	if !ok {
		h.error(errors.Wrap(consts.RESULT__RESOURCE_NOT_FOUND, fmt.Sprintf("GetResourceData:%d", rid)))
		return consts.RESULT__RESOURCE_NOT_FOUND
	}

	if err := CopyDataToInstance(h.instance, data, retaddrptr, retsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	return consts.RESULT_OK
}

func (h *host) SetResourceData(rid uint32, dataaddr, datasize int32) consts.Result {
	data, err := h.instance.GetMemory(uint64(dataaddr), uint64(datasize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	if err = h.imports.SetResourceData(rid, data); err != nil {
		h.error(errors.Wrap(err, fmt.Sprintf("SetResourceData:%d", rid)))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	return consts.RESULT_OK
}

func (h *host) GetKVData(keyaddr, keysize, retaddrptr, retsizeptr int32) consts.Result {
	key, err := h.instance.GetMemory(uint64(keyaddr), uint64(keysize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	val, err := h.imports.GetKVData(string(key))
	if err != nil {
		h.error(err)
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	if err = CopyDataToInstance(h.instance, val, retaddrptr, retsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	return consts.RESULT_OK
}

func (h *host) SetKVData(keyaddr, keysize, dataaddr, datasize int32) consts.Result {
	key, err := h.instance.GetMemory(uint64(keyaddr), uint64(keysize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}

	data, err := h.instance.GetMemory(uint64(dataaddr), uint64(datasize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}

	if err = h.imports.SetKVData(string(key), data); err != nil {
		h.error(errors.Wrap(err, "SetKVData"))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	return consts.RESULT_OK
}

func (h *host) ExecSQL(queryaddr, querysize int32) consts.Result {
	query, err := h.instance.GetMemory(uint64(queryaddr), uint64(querysize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}

	if err = h.imports.ExecSQL(string(query)); err != nil {
		h.error(errors.Wrap(err, "ExecSQL"))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	return consts.RESULT_OK
}

func (h *host) QuerySQL(queryaddr, querysize, retaddrptr, retsizeptr int32) consts.Result {
	query, err := h.instance.GetMemory(uint64(queryaddr), uint64(querysize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	res, err := h.imports.QuerySQL(string(query))
	if err != nil {
		h.error(errors.Wrap(err, "QuerySQL"))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	if err = CopyDataToInstance(h.instance, res, retaddrptr, retsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	return consts.RESULT_OK
}

func (h *host) SendTX(chainid int32, dataaddr, datasize, hashaddrptr, hashsizeptr int32) consts.Result {
	data, err := h.instance.GetMemory(uint64(dataaddr), uint64(datasize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	hash, err := h.imports.SendTX(chainid, data)
	if err != nil {
		h.error(errors.Wrap(err, "SendTX"))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	if err = CopyDataToInstance(h.instance, []byte(hash), hashaddrptr, hashsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	return consts.RESULT_OK
}

func (h *host) SendTXWithOperator(chainid int32, dataaddr, datasize, hashaddrptr, hashsizeptr int32) consts.Result {
	data, err := h.instance.GetMemory(uint64(dataaddr), uint64(datasize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	hash, err := h.imports.SendTXWithOperator(chainid, data)
	if err != nil {
		h.error(errors.Wrap(err, "SendTXWithOperator"))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	if err = CopyDataToInstance(h.instance, []byte(hash), hashaddrptr, hashsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	return consts.RESULT_OK
}

func (h *host) CallContract(chainid int32, dataaddr, datasize, resultaddrptr, resultsizeptr int32) consts.Result {
	data, err := h.instance.GetMemory(uint64(dataaddr), uint64(datasize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	res, err := h.imports.CallContract(chainid, data)
	if err != nil {
		h.error(errors.Wrap(err, "CallContract"))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	if err = CopyDataToInstance(h.instance, res, resultaddrptr, resultsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	return consts.RESULT_OK
}

func (h *host) PubMQTT(topicaddr, topicsize, msgaddr, msgsize int32) consts.Result {
	topic, err := h.instance.GetMemory(uint64(topicaddr), uint64(topicsize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	msg, err := h.instance.GetMemory(uint64(msgaddr), uint64(msgsize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}

	if err = h.imports.PubMQTT(string(topic), msg); err != nil {
		h.error(errors.Wrap(err, "PubMQTT"))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	return consts.RESULT_OK
}

func (h *host) Metrics(dataaddr, datasize int32) consts.Result {
	data, err := h.instance.GetMemory(uint64(dataaddr), uint64(datasize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	if err = h.imports.SubmitMetrics(data); err != nil {
		h.error(errors.Wrap(err, "SubmitMetrics"))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	return consts.RESULT_OK
}

func (h *host) AsyncAPICall(reqaddr, reqsize, rspaddrptr, rspsizeptr int32) consts.Result {
	req, err := h.instance.GetMemory(uint64(reqaddr), uint64(reqsize))
	if err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	rsp, err := h.imports.AsyncAPICall(req)
	if err != nil {
		h.error(errors.Wrap(err, "AsyncAPICall"))
		return consts.RESULT__IMPORT_HANDLE_FAILED
	}
	if err = CopyDataToInstance(h.instance, rsp, rspaddrptr, rspsizeptr); err != nil {
		h.error(err)
		return consts.RESULT__INVALID_MEM_ACCESS
	}
	return consts.RESULT_OK
}

func (h *host) Abort(msgaddr, filenameaddr, line, col int32) {
	msg, err := ReadStringFromAddr(h.instance, msgaddr)
	if err != nil {
		h.error(errors.Wrap(err, "message"))
		return
	}
	filename, err := ReadStringFromAddr(h.instance, filenameaddr)
	if err != nil {
		h.error(errors.Wrap(err, "filename"))
		return
	}
	h.error(errors.Errorf("%s at %s:%d:%d", msg, filename, line, col))
}

func (h *host) Trace(msgaddr, _ int32, arr ...float64) {
	msg, err := ReadStringFromAddr(h.instance, msgaddr)
	if err != nil {
		h.error(errors.Wrap(err, "message"))
		return
	}

	str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]")
	if len(str) > 0 {
		str = " " + str
	}
	h.debug(fmt.Sprintf("%s%s", msg, str))
}

func (h *host) Seed() float64 {
	return h._rand.Float64() * float64(time.Now().UnixNano())
}

func (h *host) debug(msg string) {
	var pcs [1]uintptr
	runtime.Callers(1, pcs[:])
	fs := runtime.CallersFrames([]uintptr{pcs[0]})
	f, _ := fs.Next()
	fn := f.Function
	h.imports.LogInternal(consts.LOG_LEVEL__DEBUG, fmt.Sprintf("[host:%s]:%s", fn, msg))
}

func (h *host) error(err error) {
	var pcs [1]uintptr
	runtime.Callers(1, pcs[:])
	fs := runtime.CallersFrames([]uintptr{pcs[0]})
	f, _ := fs.Next()
	fn := f.Function
	h.imports.LogInternal(consts.LOG_LEVEL__ERROR, fmt.Sprintf("[host:%s]:%s", fn, err.Error()))
}
