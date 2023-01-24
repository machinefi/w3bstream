package wasm

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	wsTypes "github.com/machinefi/w3bstream/pkg/types"
)

type ChainClient struct {
	pvk       *ecdsa.PrivateKey
	clientMap map[uint32]*ethclient.Client
}

var _web3Endpoint = map[uint32]string{
	4689: "https://babel-api.mainnet.iotex.io",
	4690: "https://babel-api.testnet.iotex.io",
}

func NewChainClient(ctx context.Context) *ChainClient {
	c := &ChainClient{
		clientMap: make(map[uint32]*ethclient.Client, 0),
	}
	ethPvk, ok := wsTypes.ETHPvkConfigFromContext(ctx)
	if ok && len(ethPvk.PrivateKey) > 0 {
		c.pvk = crypto.ToECDSAUnsafe(common.FromHex(ethPvk.PrivateKey))
	}
	return c
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
	chainEndpoint, exist := _web3Endpoint[chainID]
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
