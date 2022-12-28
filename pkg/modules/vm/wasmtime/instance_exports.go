package wasmtime

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"time"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

const (
	logTraceLevel uint32 = iota + 1
	logDebugLevel
	logInfoLevel
	logWarnLevel
	logErrorLevel
)

type (
	ExportFuncs struct {
		store  *wasmtime.Store
		res    *mapx.Map[uint32, []byte]
		db     wasm.KVStore
		pgDB   sqlx.DBExecutor
		logger conflog.Logger
		cl     *ChainClient
	}

	ChainClient struct {
		pvk   *ecdsa.PrivateKey
		chain *ethclient.Client
	}
)

func (ef *ExportFuncs) Log(c *wasmtime.Caller, logLevel, ptr, size int32) int32 {
	membuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	buf, err := read(membuf, ptr, size)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	switch uint32(logLevel) {
	case logTraceLevel:
		ef.logger.Trace(string(buf))
	case logDebugLevel:
		ef.logger.Debug(string(buf))
	case logInfoLevel:
		ef.logger.Info(string(buf))
	case logWarnLevel:
		ef.logger.Warn(errors.New(string(buf)))
	case logErrorLevel:
		ef.logger.Error(errors.New(string(buf)))
	default:
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetData(c *wasmtime.Caller, rid, vmAddrPtr, vmSizePtr int32) int32 {
	data, ok := ef.res.Load(uint32(rid))
	if !ok {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	if err := ef.copyDataIntoWasm(c, data, vmAddrPtr, vmSizePtr); err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) copyDataIntoWasm(c *wasmtime.Caller, data []byte, vmAddrPtr, vmSizePtr int32) error {
	allocFn := c.GetExport("alloc")
	if allocFn == nil {
		return errors.New("alloc is nil")
	}
	size := len(data)
	result, err := allocFn.Func().Call(ef.store, int32(size))
	if err != nil {
		return err
	}

	addr := result.(int32)

	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	if siz := copy(memBuf[addr:], data); siz != size {
		return errors.New("fail to copy data")
	}

	// fmt.Printf("host >> addr=%d\n", addr)
	// fmt.Printf("host >> size=%d\n", size)
	// fmt.Printf("host >> vmAddrPtr=%d\n", vmAddrPtr)
	// fmt.Printf("host >> vmSizePtr=%d\n", vmSizePtr)

	if err := putUint32Le(memBuf, vmAddrPtr, uint32(addr)); err != nil {
		return err
	}
	if err := putUint32Le(memBuf, vmSizePtr, uint32(size)); err != nil {
		return err
	}

	return nil
}

// TODO SetData if rid not exist, should be assigned by wasm?
func (ef *ExportFuncs) SetData(c *wasmtime.Caller, rid, addr, size int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	if addr > int32(len(memBuf)) || addr+size > int32(len(memBuf)) {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	buf, err := read(memBuf, addr, size)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	ef.res.Store(uint32(rid), buf)
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) SetDB(c *wasmtime.Caller, kAddr, kSize, vAddr, vSize int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	key, err := read(memBuf, kAddr, kSize)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}
	value, err := read(memBuf, vAddr, vSize)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	ef.logger.WithValues(
		"key", string(key),
		"val", string(value),
	).Info("host.SetDB")

	err = ef.db.Set(string(key), value)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_Failed)
	}
	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) SetSQLDB(c *wasmtime.Caller, addr, size int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	data, err := read(memBuf, addr, size)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	prestate, params, err := parseSQLQuery(data)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	if err := ef.execSQLQuery(prestate, params...); err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	return int32(wasm.ResultStatusCode_OK)
}

func parseSQLQuery(data []byte) (prestate string, params []interface{}, err error) {
	if !gjson.ValidBytes(data) {
		err = errors.New("query is invalid")
		return
	}

	res := gjson.ParseBytes(data)
	prestateRes := res.Get("statement")
	paramsRes := res.Get("params")
	if !prestateRes.Exists() || !paramsRes.Exists() {
		err = errors.New("query is invalid")
		return
	}
	prestate = prestateRes.String()

	params = make([]interface{}, 0)
	for _, para := range paramsRes.Array() {
		var res interface{}
		res, err = decodeSQLQueryParam(&para)
		if err != nil {
			return
		}
		params = append(params, res)
	}

	return
}

func decodeSQLQueryParam(in *gjson.Result) (ret interface{}, err error) {
	switch {
	case in.Get("int32").Exists():
		ret = int32(in.Get("int32").Int())
	case in.Get("int64").Exists():
		ret = int64(in.Get("int64").Int())
	case in.Get("float32").Exists():
		ret = float32(in.Get("float32").Float())
	case in.Get("float64").Exists():
		ret = float64(in.Get("float64").Float())
	case in.Get("string").Exists():
		ret = in.Get("string").String()
	case in.Get("bool").Exists():
		ret = in.Get("bool").Bool()
	case in.Get("bytes").Exists():
		ret, err = base64.StdEncoding.DecodeString(in.Get("bytes").String())
	case in.Get("time").Exists():
		ret, err = time.Parse(time.RFC3339, in.Get("time").String())
	default:
		err = errors.New("fail to decode the param")
	}
	return
}

func (ef *ExportFuncs) execSQLQuery(prestate string, params ...interface{}) error {
	_, err := ef.pgDB.ExecContext(context.Background(), prestate, params...)
	if err != nil {
		return err
	}
	return nil
}

func (ef *ExportFuncs) GetSQLDB(c *wasmtime.Caller, addr, size int32, vmAddrPtr, vmSizePtr int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	data, err := read(memBuf, addr, size)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	prestate, params, err := parseSQLQuery(data)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	rows, err := ef.querySQLQuery(prestate, params...)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	ret, err := jsonifyRows(rows)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}

	if err := ef.copyDataIntoWasm(c, ret, vmAddrPtr, vmSizePtr); err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) querySQLQuery(prestate string, params ...interface{}) (*sql.Rows, error) {
	return ef.pgDB.QueryContext(context.Background(), prestate, params...)
}

func jsonifyRows(rawRows *sql.Rows) ([]byte, error) {
	if rawRows == nil {
		return nil, errors.New("rows are empty")
	}

	columnTypes, err := rawRows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	rows := make([]interface{}, 0)
	for rawRows.Next() {
		scanArgs := make([]interface{}, len(columnTypes))
		for i := range columnTypes {
			switch columnTypes[i].DatabaseTypeName() {
			case "VARCHAR", "TEXT", "CHAR":
				scanArgs[i] = new(sql.NullString)
			case "TIMESTAMP", "TIME", "DATE":
				scanArgs[i] = new(sql.NullTime)
			case "BOOL", "BOOLEAN":
				scanArgs[i] = new(sql.NullBool)
			case "INT", "INTEGER", "SMALLINT", "BIGINT", "INT2", "INT4", "INT8":
				scanArgs[i] = new(sql.NullInt64)
			case "FLOAT", "FLOAT4", "FLOAT8", "DOUBLE":
				scanArgs[i] = new(sql.NullFloat64)
			default:
				// fmt.Println(columnTypes[i].DatabaseTypeName(), columnTypes[i].ScanType().Name())
				scanArgs[i] = new(sql.NullString)
			}
		}

		if err := rawRows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		entryMap := make(map[string]interface{}, len(columnTypes))
		for i := 0; i < len(columnTypes); i++ {
			colName := columnTypes[i].Name()
			switch v := scanArgs[i].(type) {
			case *sql.NullBool:
				if !v.Valid {
					entryMap[colName] = nil
					continue
				}
				entryMap[colName], err = v.Value()
			case *sql.NullString:
				if !v.Valid {
					entryMap[colName] = nil
					continue
				}
				entryMap[colName], err = v.Value()
			case *sql.NullFloat64:
				if !v.Valid {
					entryMap[colName] = nil
					continue
				}
				entryMap[colName], err = v.Value()
			case *sql.NullInt64:
				if !v.Valid {
					entryMap[colName] = nil
					continue
				}
				entryMap[colName], err = v.Value()
			// TODO: support time encodings
			// case *sql.NullTime:
			// 	if !v.Valid {
			// 		entryMap[colName] = nil
			// 		continue
			// 	}
			// 	entryMap[colName], err = v.Value()
			default:
				entryMap[colName] = scanArgs[i]
			}
			if err != nil {
				return nil, err
			}
		}
		rows = append(rows, entryMap)
	}
	if len(rows) == 0 {
		return []byte{}, nil
	}
	if len(rows) == 1 {
		return json.Marshal(rows[0])
	}
	return json.Marshal(rows)
}

func (ef *ExportFuncs) GetDB(c *wasmtime.Caller,
	kAddr, kSize int32, vmAddrPtr, vmSizePtr int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	key, err := read(memBuf, kAddr, kSize)
	if err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	val, exist := ef.db.Get(string(key))
	if exist != nil || val == nil {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	ef.logger.WithValues(
		"key", string(key),
		"val", string(val),
	).Info("host.GetDB")

	if err := ef.copyDataIntoWasm(c, val, vmAddrPtr, vmSizePtr); err != nil {
		ef.logger.Error(err)
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	return int32(wasm.ResultStatusCode_OK)
}

// TODO: add chainID in sendtx abi
// TODO: make sendTX async, and add callback if possible
func (ef *ExportFuncs) SendTX(c *wasmtime.Caller, offset, size int32) int32 {
	if ef.cl == nil {
		ef.logger.Error(errors.New("eth client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}
	if ef.cl.pvk == nil {
		ef.logger.Error(errors.New("private key is empty"))
		return wasm.ResultStatusCode_Failed
	}
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	buf, err := read(memBuf, offset, size)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	// fmt.Println(ret)
	txHash, err := sendETHTx(ef.cl, ret.Get("to").String(), ret.Get("value").String(), ret.Get("data").String())
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	ef.logger.Info("tx hash: %s", txHash)
	return int32(wasm.ResultStatusCode_OK)
}

func sendETHTx(cl *ChainClient, toStr string, valueStr string, dataStr string) (string, error) {
	var (
		sender = crypto.PubkeyToAddress(cl.pvk.PublicKey)
		to     = common.HexToAddress(toStr)
	)
	value, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		return "", errors.New("fail to read tx value")
	}
	data, err := hex.DecodeString(dataStr)
	if err != nil {
		return "", err

	}
	nonce, err := cl.chain.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return "", err
	}

	gasPrice, err := cl.chain.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		From:     sender,
		To:       &to,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}
	gasLimit, err := cl.chain.EstimateGas(context.Background(), msg)
	if err != nil {
		return "", err
	}

	// Create a new transaction
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       &to,
			Value:    value,
			Data:     data,
		})

	chainid, err := cl.chain.ChainID(context.Background())
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainid), cl.pvk)
	if err != nil {
		return "", err
	}
	err = cl.chain.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}

func (ef *ExportFuncs) CallContract(c *wasmtime.Caller,
	offset, size int32, vmAddrPtr, vmSizePtr int32) int32 {
	if ef.cl == nil {
		ef.logger.Error(errors.New("eth client doesn't exist"))
		return wasm.ResultStatusCode_Failed
	}
	if ef.cl.pvk == nil {
		ef.logger.Error(errors.New("private key is empty"))
		return wasm.ResultStatusCode_Failed
	}
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	buf, err := read(memBuf, offset, size)
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	// fmt.Println(ret)
	data, err := callContract(ef.cl.chain, ret.Get("to").String(), ret.Get("data").String())
	if err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	if err := ef.copyDataIntoWasm(c, data, vmAddrPtr, vmSizePtr); err != nil {
		ef.logger.Error(err)
		return wasm.ResultStatusCode_Failed
	}
	return int32(wasm.ResultStatusCode_OK)
}

func callContract(cl *ethclient.Client, toStr string, dataStr string) ([]byte, error) {
	var (
		to      = common.HexToAddress(toStr)
		data, _ = hex.DecodeString(dataStr)
	)

	msg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}

	return cl.CallContract(context.Background(), msg, nil)
}

func putUint32Le(buf []byte, addr int32, num uint32) error {
	if int32(len(buf)) < addr+4 {
		return errors.New("overflow")
	}
	binary.LittleEndian.PutUint32(buf[addr:], num)
	return nil
}

func read(memBuf []byte, addr int32, size int32) ([]byte, error) {
	if addr > int32(len(memBuf)) || addr+size > int32(len(memBuf)) {
		return nil, errors.New("overflow")
	}
	buf := make([]byte, size)
	if siz := copy(buf, memBuf[addr:addr+size]); int32(siz) != size {
		return nil, errors.New("overflow")
	}
	return buf, nil
}
