package wasm

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/pkg/errors"
)

type ChainClient struct {
	PrivateKey string `json:"privateKey"`

	pvk       *ecdsa.PrivateKey
	clientMap map[uint32]*ethclient.Client
}

var _web3Endpoint = map[uint32]string{
	4689: "https://babel-api.mainnet.iotex.io",
	4690: "https://babel-api.testnet.iotex.io",
}

func (c *ChainClient) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__CHAIN_CLIENT
}

func (c *ChainClient) WithContext(ctx context.Context) context.Context {
	if err := c.Build(); err != nil {
		return ctx
	}
	return WithChainClient(ctx, c)
}

func (c *ChainClient) Build() error {
	if len(c.PrivateKey) > 0 {
		c.pvk = crypto.ToECDSAUnsafe(common.FromHex(c.PrivateKey))
	}

	return nil
}

func (c *ChainClient) SendTX(chainID uint32, toStr, valueStr, dataStr string) (string, error) {
	if c == nil {
		return "", nil
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
	data, err := hex.DecodeString(dataStr)
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
		to      = common.HexToAddress(toStr)
		data, _ = hex.DecodeString(dataStr)
	)

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
