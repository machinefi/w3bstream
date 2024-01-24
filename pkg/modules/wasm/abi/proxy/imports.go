package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflogger "github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	optypes "github.com/machinefi/w3bstream/pkg/modules/operator/pool/types"
	wasmapi "github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/types"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/consts"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
	"github.com/machinefi/w3bstream/pkg/types/wasm/sql_util"
)

var (
	errWithoutChainConfig       = errors.New("without chain config")
	errWithoutMQTTConfig        = errors.New("without mqtt config")
	errWithoutMetricsConfig     = errors.New("without metrics config")
	errWithoutDatabaseConfig    = errors.New("without database config")
	errWithoutKVStoreConfig     = errors.New("without kv store config")
	errWithoutAsyncServerConfig = errors.New("without async server config")
	errInvalidMetricsBody       = errors.New("invalid metrics body")
)

type chainctx struct {
	client    *wasm.ChainClient
	config    *types.ChainConfig
	operators optypes.Pool
}

func (v *chainctx) Valid() bool {
	return v != nil && v.client != nil && v.config != nil && v.operators != nil
}

func NewImports(ctx context.Context) *Imports {
	chain := &chainctx{}

	if v, ok := wasm.ChainClientFromContext(ctx); ok && v != nil {
		chain.client = v
	}
	if v, ok := types.ChainConfigFromContext(ctx); ok && v != nil {
		chain.config = v
	}
	if v, ok := types.OperatorPoolFromContext(ctx); ok && v != nil {
		chain.operators = v
	}

	i := &Imports{
		parent:    ctx,
		logger:    conflogger.Std(),
		logs:      make(chan *models.WasmLog, 1024),
		_rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		dbe:       types.MustMgrDBExecutorFromContext(ctx),
		idg:       confid.MustNewSFIDGenerator(),
		prj:       types.MustProjectFromContext(ctx),
		app:       types.MustAppletFromContext(ctx),
		ins:       types.MustInstanceFromContext(ctx),
		resources: mapx.New[uint32, []byte](),
		events:    mapx.New[uint32, string](),
		chain:     chain,
	}
	i.envs, _ = wasm.EnvFromContext(ctx)
	i.db, _ = wasm.SQLStoreFromContext(ctx)
	i.kv, _ = wasm.KVStoreFromContext(ctx)
	i.mqtt, _ = wasm.MQTTClientFromContext(ctx)
	i.metrics, _ = wasm.CustomMetricsFromContext(ctx)
	return i
}

func NewImportsDebugMode(ctx context.Context) *Imports {
	i := NewImports(ctx)
	i.debug = true
	return i
}

type Imports struct {
	parent    context.Context
	logger    logr.Logger
	debug     bool
	logs      chan *models.WasmLog
	_rand     *rand.Rand
	dbe       sqlx.DBExecutor
	idg       confid.SFIDGenerator
	prj       *models.Project
	app       *models.Applet
	ins       *models.Instance
	envs      *wasm.Env
	resources *mapx.Map[uint32, []byte]
	events    *mapx.Map[uint32, string]
	db        *wasm.Database
	kv        wasm.KVStore
	mqtt      *mqtt.Client
	metrics   metrics.CustomMetrics
	asyncsrv  wasmapi.Server
	chain     *chainctx
}

func (i *Imports) Log(lv consts.LogLevel, msg string) {
	if i.debug {
		l := i.logger.WithValues("@src", "wasm")
		switch lv {
		case consts.LOG_LEVEL__ERROR:
			l.Error(errors.New(msg))
		case consts.LOG_LEVEL__WARN:
			l.Warn(errors.New(msg))
		case consts.LOG_LEVEL__INFO:
			l.Info(msg)
		default:
			l.Debug(msg)
		}
	}
	i.logs <- &models.WasmLog{
		WasmLogInfo: models.WasmLogInfo{
			ProjectName: i.prj.Name,
			AppletName:  i.app.Name,
			InstanceID:  i.ins.InstanceID,
			Src:         "wasm",
			Level:       lv.String(),
			LogTime:     time.Now().UnixNano(),
			Msg:         subStringWithLength(msg, enums.WasmLogMaxLength),
		},
	}
}

func (i *Imports) LogInternal(lv consts.LogLevel, msg string, args ...any) {
	l := i.logger.WithValues(append(args, "@src", "host")...)

	switch lv {
	case consts.LOG_LEVEL__ERROR:
		l.Error(errors.New(msg))
	case consts.LOG_LEVEL__WARN:
		l.Warn(errors.New(msg))
	case consts.LOG_LEVEL__INFO:
		l.Info(msg)
	default:
		l.Debug(msg)
	}
	msg = fmt.Sprintf("%s %v", msg, args)
	i.logs <- &models.WasmLog{
		WasmLogInfo: models.WasmLogInfo{
			ProjectName: i.prj.Name,
			AppletName:  i.app.Name,
			InstanceID:  i.ins.InstanceID,
			Src:         "host",
			Level:       lv.String(),
			LogTime:     time.Now().UnixNano(),
			Msg:         subStringWithLength(msg, enums.WasmLogMaxLength),
		},
	}

}

func (i *Imports) FlushLog() {
	logs := make([]*models.WasmLog, 0, 128)
	// lock
	size := 128
	if len(i.logs) < size {
		size = len(i.logs)
	}
	for idx := 0; idx < size; idx++ {
		logs = append(logs, <-i.logs)
	}
	// unlock
	ids := i.idg.MustGenSFIDs(len(logs))
	for i, l := range logs {
		l.WasmLogID = ids[i]
	}
	err := models.BatchCreateWasmLogs(i.dbe, logs...)
	if err != nil {
		i.logger.WithValues("@src", "host").Error(err)
	}
}

func (i *Imports) Env(key string) (string, bool) {
	if i.envs == nil {
		return os.LookupEnv(key)
	}
	return i.envs.Get(key)
}

func (i *Imports) Abort(msg, filename string, line, col int32) {
	// var pcs [1]uintptr
	// runtime.Callers(1, pcs[:])
	// fs := runtime.CallersFrames([]uintptr{pcs[0]})
	// f, _ := fs.Next()
	// fn := f.Function
	i.LogInternal(consts.LOG_LEVEL__ERROR, msg, "filename", filename, "line", line, "col", col)
}

func (i *Imports) Trace(msg, trace string) {
	// var pcs [1]uintptr
	// runtime.Callers(1, pcs[:])
	// fs := runtime.CallersFrames([]uintptr{pcs[0]})
	// f, _ := fs.Next()
	// fn := f.Function
	i.LogInternal(consts.LOG_LEVEL__DEBUG, msg, "trace", trace)
}

func (i *Imports) Seed() float64 {
	return i._rand.Float64() * float64(time.Now().UnixNano())
}

func (i *Imports) GetResourceData(rid uint32) ([]byte, bool) {
	return i.resources.Load(rid)
}

func (i *Imports) SetResourceData(rid uint32, data []byte) error {
	i.resources.Store(rid, data)
	return nil
}

func (i *Imports) RemoveResourceData(rid uint32) {
	i.resources.Remove(rid)
}

func (i *Imports) GetEventType(rid uint32) (string, bool) {
	return i.events.Load(rid)
}

func (i *Imports) SetEventType(rid uint32, typ string) error {
	i.events.Store(rid, typ)
	return nil
}

func (i *Imports) RemoveEventType(rid uint32) {
	i.events.Remove(rid)
}

func (i *Imports) GetKVData(key string) ([]byte, error) {
	if i.kv == nil {
		return nil, errWithoutKVStoreConfig
	}
	return i.kv.Get(key)
}

func (i *Imports) SetKVData(key string, data []byte) error {
	if i.kv == nil {
		return errWithoutKVStoreConfig
	}
	return i.kv.Set(key, data)
}

func (i *Imports) ExecSQL(q string) error {
	if i.db == nil {
		return errWithoutDatabaseConfig
	}
	prestate, params, err := sql_util.ParseQuery([]byte(q))
	if err != nil {
		return err
	}
	db, err := i.db.WithDefaultSchema()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(context.Background(), prestate, params...)
	return err
}

func (i *Imports) QuerySQL(q string) ([]byte, error) {
	if i.db == nil {
		return nil, errWithoutDatabaseConfig
	}
	prestate, params, err := sql_util.ParseQuery([]byte(q))
	if err != nil {
		return nil, err
	}
	db, err := i.db.WithDefaultSchema()
	if err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(context.Background(), prestate, params...)
	if err != nil {
		return nil, err
	}
	return sql_util.JsonifyRows(rows)
}

func (i *Imports) SendTX(chainID uint32, data []byte) (string, error) {
	if !i.chain.Valid() {
		return "", errWithoutChainConfig
	}
	res := gjson.Parse(string(data))
	return i.chain.client.SendTX(
		i.chain.config,
		uint64(chainID),
		"",
		res.Get("to").String(),
		res.Get("value").String(),
		res.Get("data").String(),
		i.chain.operators,
		i.prj,
	)
}

func (i *Imports) SendTXWithOperator(chainID uint32, data []byte) (string, error) {
	if !i.chain.Valid() {
		return "", errWithoutChainConfig
	}
	res := gjson.Parse(string(data))
	tx, err := i.chain.client.SendTXWithOperator(
		i.chain.config,
		uint64(chainID),
		"",
		res.Get("to").String(),
		res.Get("value").String(),
		res.Get("data").String(),
		res.Get("operatorName").String(),
		i.chain.operators,
		i.prj,
	)
	if err != nil {
		return "", err
	}
	return tx.Hash, err
}

func (i *Imports) CallContract(chainID uint32, data []byte) ([]byte, error) {
	if !i.chain.Valid() {
		return nil, errWithoutChainConfig
	}
	res := gjson.Parse(string(data))

	return i.chain.client.CallContract(
		i.chain.config,
		uint64(chainID),
		"",
		res.Get("to").String(),
		res.Get("data").String(),
	)
}

func (i *Imports) PubMQTT(topic string, message []byte) error {
	if i.mqtt == nil {
		return errWithoutMQTTConfig
	}
	return i.mqtt.WithTopic(topic).Publish(message)
}

func (i *Imports) SubmitMetrics(data []byte) error {
	if i.metrics == nil {
		return errWithoutMetricsConfig
	}

	datastr := string(data)
	if !gjson.Valid(datastr) {
		return errInvalidMetricsBody
	}
	object := gjson.Parse(datastr)
	if object.IsArray() {
		return errInvalidMetricsBody
	}
	return i.metrics.Submit(object)
}

func (i *Imports) AsyncAPICall(req []byte) ([]byte, error) {
	if i.asyncsrv == nil {
		return nil, errWithoutAsyncServerConfig
	}

	rsp := i.asyncsrv.Call(i.parent, req)
	return json.Marshal(rsp)
}
