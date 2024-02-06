package wasmtime

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"golang.org/x/text/encoding/unicode"

	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	optypes "github.com/machinefi/w3bstream/pkg/modules/operator/pool/types"
	wasmapi "github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/types"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
	"github.com/machinefi/w3bstream/pkg/types/wasm/sql_util"
)

type (
	Import func(module, name string, f interface{}) error

	ABILinker interface {
		LinkABI(Import) error
	}

	ExportFuncs struct {
		rt      *Runtime
		prj     *models.Project
		app     *models.Applet
		ins     *models.Instance
		ctxID   *atomic.Value
		logs    []*models.WasmLog
		res     *mapx.Map[uint32, []byte]
		evs     *mapx.Map[uint32, []byte]
		env     *wasm.Env
		kvs     wasm.KVStore
		db      *wasm.Database
		logger  conflog.Logger
		cl      *wasm.ChainClient
		cf      *types.ChainConfig
		ctx     context.Context
		mq      *confmqtt.Client
		metrics metrics.CustomMetrics
		srv     wasmapi.Server
		opPool  optypes.Pool
	}
)

func NewExportFuncs(ctx context.Context, rt *Runtime) (*ExportFuncs, error) {
	ef := &ExportFuncs{
		prj:     types.MustProjectFromContext(ctx),
		app:     types.MustAppletFromContext(ctx),
		ins:     types.MustInstanceFromContext(ctx),
		ctxID:   &atomic.Value{},
		res:     wasm.MustRuntimeResourceFromContext(ctx),
		evs:     wasm.MustRuntimeEventTypesFromContext(ctx),
		kvs:     wasm.MustKVStoreFromContext(ctx),
		logger:  wasm.MustLoggerFromContext(ctx),
		srv:     types.MustWasmApiServerFromContext(ctx),
		opPool:  types.MustOperatorPoolFromContext(ctx),
		cl:      wasm.MustChainClientFromContext(ctx),
		cf:      types.MustChainConfigFromContext(ctx),
		db:      wasm.MustSQLStoreFromContext(ctx),
		env:     wasm.MustEnvFromContext(ctx),
		mq:      wasm.MustMQTTClientFromContext(ctx),
		metrics: wasm.MustCustomMetricsFromContext(ctx),
		rt:      rt,
		ctx:     ctx,
	}
	ef.ctxID.Store("")

	return ef, nil
}

var (
	_     wasm.ABI = (*ExportFuncs)(nil)
	_rand          = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func (ef *ExportFuncs) LinkABI(impt Import) error {
	for name, ff := range map[string]interface{}{
		"abort":                    ef.Abort,
		"trace":                    ef.Trace,
		"seed":                     ef.Seed,
		"ws_log":                   ef.Log,
		"ws_get_data":              ef.GetData,
		"ws_set_data":              ef.SetData,
		"ws_get_db":                ef.GetDB,
		"ws_set_db":                ef.SetDB,
		"ws_send_tx":               ef.SendTX,
		"ws_send_tx_with_operator": ef.SendTXWithOperator,
		"ws_call_contract":         ef.CallContract,
		"ws_set_sql_db":            ef.SetSQLDB,
		"ws_get_sql_db":            ef.GetSQLDB,
		"ws_get_env":               ef.GetEnv,
		"ws_send_mqtt_msg":         ef.SendMqttMsg,
		"ws_api_call":              ef.ApiCall,
	} {
		if err := impt("env", name, ff); err != nil {
			return err
		}
	}

	for name, ff := range map[string]interface{}{
		"ws_submit_metrics": ef.StatSubmit,
	} {
		if err := impt("stat", name, ff); err != nil {
			return err
		}
	}

	return nil
}

func (ef *ExportFuncs) _log(level conflog.Level, src string, msg any) {
	l := ef.logger.WithValues("@src", src)
	switch level {
	case conflog.DebugLevel:
		l.Debug(msg.(string))
	case conflog.InfoLevel:
		l.Info(msg.(string))
	case conflog.WarnLevel:
		l.Warn(msg.(error))
	case conflog.ErrorLevel:
		l.Error(msg.(error))
	default:
		l.Trace(msg.(string))
	}
	ef.logs = append(ef.logs, &models.WasmLog{
		WasmLogInfo: models.WasmLogInfo{
			ProjectName: ef.prj.Name,
			AppletName:  ef.app.Name,
			InstanceID:  ef.ins.InstanceID,
			EventID:     ef.ContextID(),
			Src:         src,
			Level:       level.String(),
			LogTime:     time.Now().UnixNano(),
			Msg:         fmt.Sprint(msg),
		},
		OperationTimes: datatypes.OperationTimes{
			CreatedAt: base.Timestamp{Time: time.Now()},
			UpdatedAt: base.Timestamp{Time: time.Now()},
		},
	})
}

func (ef *ExportFuncs) WasmLog(lv conflog.Level, msg any) {
	ef._log(lv, "wasm", msg)
}

func (ef *ExportFuncs) HostLog(lv conflog.Level, msg any) {
	ef._log(lv, "host", msg)
}

func (ef *ExportFuncs) EntryContext(ctx context.Context, ctxID string, tpe, data []byte) (uint32, bool) {
	if ef.ctxID.CompareAndSwap("", ctxID) {
		rid := uuid.New().ID() % uint32(maxInt)
		ef.evs.Store(rid, tpe)
		ef.res.Store(rid, data)
		return rid, true
	}
	return 0, false
}

func (ef *ExportFuncs) LeaveContext(ctx context.Context, ctxID string, rid uint32) bool {
	if ef.ctxID.Load().(string) == ctxID {
		idg := confid.MustSFIDGeneratorFromContext(ctx)
		ids := idg.MustGenSFIDs(len(ef.logs))
		d := types.MustMgrDBExecutorFromContext(ctx)

		for i, l := range ef.logs {
			l.WasmLogID = ids[i]
		}
		err := models.BatchCreateWasmLogs(d, ef.logs...)
		if err != nil {
			ef.WasmLog(conflog.WarnLevel, err)
		}
		ef.logs = ef.logs[:0]
		ef.evs.Remove(rid)
		ef.res.Remove(rid)
		return ef.ctxID.CompareAndSwap(ctxID, "")
	}
	return false
}

func (ef *ExportFuncs) ContextID() string {
	return ef.ctxID.Load().(string)
}

func (ef *ExportFuncs) Log(level, ptr, size int32) int32 {
	buf, err := ef.rt.Read(ptr, size)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	ef.WasmLog(conflog.Level(level), string(buf))
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) ApiCall(kAddr, kSize, vmAddrPtr, vmSizePtr int32) int32 {
	buf, err := ef.rt.Read(kAddr, kSize)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_TransDataFromVMFailed)
	}

	resp := ef.srv.Call(ef.ctx, buf)

	respJson, err := json.Marshal(resp)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_HostInternal)
	}

	if err = ef.rt.Copy(respJson, vmAddrPtr, vmSizePtr); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

// Abort is reserved for imported func env.abort() which is auto-generated by assemblyScript
func (ef *ExportFuncs) Abort(msgPtr int32, fileNamePtr int32, line int32, col int32) {
	msg, err := ef.readString(msgPtr)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, errors.Wrap(err, "fail to decode arguments in env.abort"))
		return
	}
	fileName, err := ef.readString(fileNamePtr)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, errors.Wrap(err, "fail to decode arguments in env.abort"))
		return
	}
	ef.HostLog(conflog.ErrorLevel, errors.Errorf("abort: %s at %s:%d:%d", msg, fileName, line, col))
}

func (ef *ExportFuncs) readString(ptr int32) (string, error) {
	if ptr < 4 {
		return "", errors.Errorf("the pointer address %d is invalid", ptr)
	}

	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()

	lenData, err := ef.rt.Read(ptr-4, 4) // sizeof(uint32) is 4
	if err != nil {
		return "", err
	}
	len := binary.LittleEndian.Uint32(lenData)
	data, err := ef.rt.Read(ptr, int32(len))
	if err != nil {
		return "", err
	}
	utf8bytes, err := decoder.Bytes(data)
	if err != nil {
		return "", err
	}
	return string(utf8bytes), nil
}

// Trace is reserved for imported func env.trace() which is auto-generated by assemblyScript
func (ef *ExportFuncs) Trace(msgPtr int32, _ int32, arr ...float64) {
	msg, err := ef.readString(msgPtr)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, errors.Wrap(err, "fail to decode arguments in env.abort"))
		return
	}

	str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ", "), "[]")
	if len(str) > 0 {
		str = " " + str
	}
	ef.HostLog(conflog.InfoLevel, fmt.Sprintf("trace: %s%s", msg, str))
}

// Seed is reserved for imported func env.seed() which is auto-generated by assemblyScript
func (ef *ExportFuncs) Seed() float64 {
	return _rand.Float64() * float64(time.Now().UnixNano())
}

func (ef *ExportFuncs) GetData(rid, vmAddrPtr, vmSizePtr int32) int32 {
	data, ok := ef.res.Load(uint32(rid))
	if !ok {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	if err := ef.rt.Copy(data, vmAddrPtr, vmSizePtr); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

// TODO SetData if rid not exist, should be assigned by wasm?
func (ef *ExportFuncs) SetData(rid, addr, size int32) int32 {
	buf, err := ef.rt.Read(addr, size)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	ef.res.Store(uint32(rid), buf)
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetDB(kAddr, kSize int32, vmAddrPtr, vmSizePtr int32) int32 {
	key, err := ef.rt.Read(kAddr, kSize)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	val, exist := ef.kvs.Get(string(key))
	if exist != nil || val == nil {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	ef.HostLog(conflog.InfoLevel, fmt.Sprintf("host.GetDB %s:%s", string(key), strconv.Quote(string(val))))

	if err = ef.rt.Copy(val, vmAddrPtr, vmSizePtr); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) SetDB(kAddr, kSize, vAddr, vSize int32) int32 {
	key, err := ef.rt.Read(kAddr, kSize)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}
	val, err := ef.rt.Read(vAddr, vSize)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	ef.HostLog(conflog.InfoLevel, fmt.Sprintf("host.SetDB %s:%s", string(key), strconv.Quote(string(val))))

	err = ef.kvs.Set(string(key), val)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_Failed)
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) SetSQLDB(addr, size int32) int32 {
	if ef.db == nil {
		return int32(wasm.ResultStatusCode_NoDBContext)
	}
	data, err := ef.rt.Read(addr, size)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	prestate, params, err := sql_util.ParseQuery(data)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}

	db, err := ef.db.WithDefaultSchema()
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	_, err = db.ExecContext(context.Background(), prestate, params...)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}

	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetSQLDB(addr, size int32, vmAddrPtr, vmSizePtr int32) int32 {
	if ef.db == nil {
		return int32(wasm.ResultStatusCode_NoDBContext)
	}
	data, err := ef.rt.Read(addr, size)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	prestate, params, err := sql_util.ParseQuery(data)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}

	db, err := ef.db.WithDefaultSchema()
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	rows, err := db.QueryContext(context.Background(), prestate, params...)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}

	ret, err := sql_util.JsonifyRows(rows)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}

	if err = ef.rt.Copy(ret, vmAddrPtr, vmSizePtr); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

// TODO: make sendTX async, and add callback if possible
func (ef *ExportFuncs) SendTX(chainID int32, offset, size, vmAddrPtr, vmSizePtr int32) int32 {
	ef.HostLog(conflog.InfoLevel, fmt.Sprintf("offset %d size %d vmAddrPtr %d vmSizePtr %d", offset, size, vmAddrPtr, vmSizePtr))
	if ef.cl == nil {
		ef.HostLog(conflog.ErrorLevel, errors.New("eth client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}
	buf, err := ef.rt.Read(offset, size)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	txHash, err := ef.cl.SendTX(ef.cf, uint64(chainID), "", ret.Get("to").String(), ret.Get("value").String(), ret.Get("data").String(), ef.opPool, ef.prj)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	if err := ef.rt.Copy([]byte(txHash), vmAddrPtr, vmSizePtr); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) SendTXWithOperator(chainID int32, offset, size, vmAddrPtr, vmSizePtr int32) int32 {
	if ef.cl == nil {
		ef.HostLog(conflog.ErrorLevel, errors.New("eth client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}
	buf, err := ef.rt.Read(offset, size)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	txResp, err := ef.cl.SendTXWithOperator(ef.cf, uint64(chainID), "", ret.Get("to").String(), ret.Get("value").String(), ret.Get("data").String(), ret.Get("operatorName").String(), ef.opPool, types.MustProjectFromContext(ef.ctx))
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	if err := ef.rt.Copy([]byte(txResp.Hash), vmAddrPtr, vmSizePtr); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) SendMqttMsg(topicAddr, topicSize, msgAddr, msgSize int32) int32 {
	if ef.mq == nil {
		ef.HostLog(conflog.ErrorLevel, errors.New("mq client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}

	var (
		topicBuf []byte
		msgBuf   []byte
		err      error
	)

	topicBuf, err = ef.rt.Read(topicAddr, topicSize)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	msgBuf, err = ef.rt.Read(msgAddr, msgSize)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	err = ef.mq.WithTopic(string(topicBuf)).Publish(string(msgBuf))
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) CallContract(chainID int32, offset, size int32, vmAddrPtr, vmSizePtr int32) int32 {
	if ef.cl == nil {
		ef.HostLog(conflog.ErrorLevel, errors.New("eth client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}
	buf, err := ef.rt.Read(offset, size)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	data, err := ef.cl.CallContract(ef.cf, uint64(chainID), "", ret.Get("to").String(), ret.Get("data").String())
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	if err = ef.rt.Copy(data, vmAddrPtr, vmSizePtr); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetEnv(kAddr, kSize int32, vmAddrPtr, vmSizePtr int32) int32 {
	if ef.env == nil {
		return int32(wasm.ResultStatusCode_EnvKeyNotFound)
	}
	key, err := ef.rt.Read(kAddr, kSize)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	val, ok := ef.env.Get(string(key))
	if !ok {
		return int32(wasm.ResultStatusCode_EnvKeyNotFound)
	}

	if err = ef.rt.Copy([]byte(val), vmAddrPtr, vmSizePtr); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetEventType(rid, vmAddrPtr, vmSizePtr int32) int32 {
	data, ok := ef.res.Load(uint32(rid))
	if !ok {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	if err := ef.rt.Copy(data, vmAddrPtr, vmSizePtr); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) StatSubmit(vmAddrPtr, vmSizePtr int32) int32 {
	buf, err := ef.rt.Read(vmAddrPtr, vmSizePtr)
	if err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	str := string(buf)
	if !gjson.Valid(str) {
		err = errors.New("invalid json")
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	object := gjson.Parse(str)
	if object.IsArray() {
		err = errors.New("json object should not be an array")
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}

	if err = ef.metrics.Submit(object); err != nil {
		ef.HostLog(conflog.ErrorLevel, err)
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}
