package wasmtime

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"math/big"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	conflog "github.com/iotexproject/Bumblebee/conf/log"
	"github.com/iotexproject/Bumblebee/x/mapx"
	"github.com/pkg/errors"

	"github.com/tidwall/gjson"

	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

type (
	ExportFuncs struct {
		store  *wasmtime.Store
		res    *mapx.Map[uint32, []byte]
		db     map[string]int32
		logger conflog.Logger
		cl     *ChainClient
	}

	ChainClient struct {
		pvk   *ecdsa.PrivateKey
		chain *ethclient.Client
	}
)

func (ef *ExportFuncs) Log(c *wasmtime.Caller, ptr, size int32) {
	membuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	buf, err := read(membuf, ptr, size)
	if err != nil {
		panic(err)
		// return wasm.ResultStatusCode_Failed
	}
	ef.logger.Info(string(buf))
	// return int32(wasm.ResultStatusCode_OK)
}

func (ef *ExportFuncs) GetData(c *wasmtime.Caller, rid, vmAddrPtr, vmSizePtr int32) int32 {
	allocFn := c.GetExport("alloc")
	if allocFn == nil {
		return int32(wasm.ResultStatusCode_ImportNotFound)
	}
	data, ok := ef.res.Load(uint32(rid))
	if !ok {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}
	size := len(data)
	result, err := allocFn.Func().Call(ef.store, int32(size))
	if err != nil {
		return int32(wasm.ResultStatusCode_ImportCallFailed)
	}
	addr := result.(int32)

	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	if siz := copy(memBuf[addr:], data); siz != size {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}

	// fmt.Printf("host >> addr=%d\n", addr)
	// fmt.Printf("host >> size=%d\n", size)
	// fmt.Printf("host >> vmAddrPtr=%d\n", vmAddrPtr)
	// fmt.Printf("host >> vmSizePtr=%d\n", vmSizePtr)

	if err := putUint32Le(memBuf, vmAddrPtr, uint32(addr)); err != nil {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	if err := putUint32Le(memBuf, vmSizePtr, uint32(size)); err != nil {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	// fmt.Println("host >> get_data returned")
	return int32(wasm.ResultStatusCode_OK)
}

// TODO SetData if rid not exist, should be assigned by wasm?
func (ef *ExportFuncs) SetData(c *wasmtime.Caller, rid, addr, size int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	if addr > int32(len(memBuf)) || addr+size > int32(len(memBuf)) {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	buf, err := read(memBuf, addr, size)
	if err != nil {
		return int32(wasm.ResultStatusCode_TransDataToVMFailed)
	}
	ef.res.Store(uint32(rid), buf)
	return int32(wasm.ResultStatusCode_OK)
}

// TODO SetDB value should have type
func (ef *ExportFuncs) SetDB(c *wasmtime.Caller, kAddr, kSize, val int32) {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	key, _ := read(memBuf, kAddr, kSize)

	ef.logger.WithValues(
		"key", string(key),
		"val", val,
	).Info("host.SetDB")

	ef.db[string(key)] = val
}

func (ef *ExportFuncs) GetDB(c *wasmtime.Caller, kAddr, kSize int32) int32 {
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	key, err := read(memBuf, kAddr, kSize)
	if err != nil {
		return int32(wasm.ResultStatusCode_ResourceNotFound)
	}

	val := ef.db[string(key)]

	ef.logger.WithValues(
		"key", string(key),
		"val", val,
	).Info("host.GetDB")

	return val
}

// TODO: add chainID in sendtx abi
// TODO: make sendTX async, and add callback if possible
func (ef *ExportFuncs) SendTX(c *wasmtime.Caller, offset, size int32) int32 {
	if ef.cl == nil {
		return wasm.ResultStatusCode_Failed
	}
	memBuf := c.GetExport("memory").Memory().UnsafeData(ef.store)
	buf, err := read(memBuf, offset, size)
	if err != nil {
		return wasm.ResultStatusCode_Failed
	}
	ret := gjson.Parse(string(buf))
	// fmt.Println(ret)
	txHash, err := sentETHTx(ef.cl, ret.Get("to").String(), ret.Get("value").String(), ret.Get("data").String())
	if err != nil {
		return wasm.ResultStatusCode_Failed
	}
	ef.logger.Info("tx hash: %s", txHash)
	return int32(wasm.ResultStatusCode_OK)
}

func sentETHTx(cl *ChainClient, toStr string, valueStr string, dataStr string) (string, error) {
	var (
		sender = crypto.PubkeyToAddress(cl.pvk.PublicKey)
	)
	var (
		to       = common.HexToAddress(toStr)
		value, _ = new(big.Int).SetString(valueStr, 10)
		data, _  = hex.DecodeString(dataStr)
	)

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
	signedTx, _ := types.SignTx(tx, types.NewLondonSigner(chainid), cl.pvk)
	err = cl.chain.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
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
