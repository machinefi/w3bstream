package wasm

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	wsTypes "github.com/machinefi/w3bstream/pkg/types"
)

type ChainClient struct {
	pvk       *ecdsa.PrivateKey
	endpoints map[uint32]string
	clientMap map[uint32]*ethclient.Client
}

func NewChainClient(ctx context.Context) *ChainClient {
	c := &ChainClient{
		clientMap: make(map[uint32]*ethclient.Client, 0),
		endpoints: make(map[uint32]string),
	}
	ethcli, ok := wsTypes.ETHClientConfigFromContext(ctx)
	if !ok || ethcli == nil {
		return c
	}
	acc := middleware.CurrentAccountFromContext(ctx)
	if len(acc.OperatorPrivateKey) > 0 {
		c.pvk = crypto.ToECDSAUnsafe(common.FromHex(acc.OperatorPrivateKey))
	}
	if len(ethcli.Endpoints) > 0 {
		c.endpoints = decodeEndpoints(ethcli.Endpoints)
	}
	return c
}

func decodeEndpoints(in string) (ret map[uint32]string) {
	ret = make(map[uint32]string)
	if !gjson.Valid(in) {
		return
	}
	for k, v := range gjson.Parse(in).Map() {
		chainID, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		url := v.String()
		ret[uint32(chainID)] = url
	}
	return
}

func (c *ChainClient) SendTX(chainID uint32, toStr, valueStr, dataStr string) (string, error) {
	if c == nil {
		return "", nil
	}
	if c.pvk == nil {
		return "", errors.New("private key is empty")
	}
	cli, err := c.getEthClient(chainID)
	if err != nil {
		return "", err
	}
	var (
		sender = crypto.PubkeyToAddress(c.pvk.PublicKey)
		to     = common.HexToAddress(toStr)
	)
	value, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		return "", errors.New("fail to read tx value")
	}
	data, err := hex.DecodeString(strings.TrimPrefix(dataStr, "0x"))
	if err != nil {
		return "", err
	}
	nonce, err := cli.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return "", err
	}

	gasPrice, err := cli.SuggestGasPrice(context.Background())
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
	gasLimit, err := cli.EstimateGas(context.Background(), msg)
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

	chainid, err := cli.ChainID(context.Background())
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainid), c.pvk)
	if err != nil {
		return "", err
	}
	err = cli.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil

}

func (c *ChainClient) getEthClient(chainID uint32) (*ethclient.Client, error) {
	if cli, exist := c.clientMap[chainID]; exist {
		return cli, nil
	}
	chainEndpoint, exist := c.endpoints[chainID]
	if !exist {
		return nil, errors.Errorf("the chain %d is not supported", chainID)
	}
	chain, err := ethclient.Dial(chainEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "fail to dial the endpoint of the chain")
	}
	c.clientMap[chainID] = chain
	return chain, nil
}

func (c *ChainClient) CallContract(chainID uint32, toStr, dataStr string) ([]byte, error) {
	var (
		to = common.HexToAddress(toStr)
	)
	data, err := hex.DecodeString(strings.TrimPrefix(dataStr, "0x"))
	if err != nil {
		return nil, err
	}
	cli, err := c.getEthClient(chainID)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}
	return cli.CallContract(context.Background(), msg, nil)
}
