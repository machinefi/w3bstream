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
	PrivateKey    string `json:"privateKey"`
	ChainEndpoint string `chainEndpoint:"chainEndpoint"`

	pvk   *ecdsa.PrivateKey
	chain *ethclient.Client
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
	if len(c.ChainEndpoint) == 0 {
		return errors.New("no chain client is established due to empty chain endpoint")
	}

	chain, err := ethclient.Dial(c.ChainEndpoint)
	if err != nil {
		return errors.Wrap(err, "fail to dial the endpoint of the chain")
	}
	c.chain = chain

	if len(c.PrivateKey) > 0 {
		c.pvk = crypto.ToECDSAUnsafe(common.FromHex(c.PrivateKey))
	}

	return nil
}

func (c *ChainClient) SendTX(toStr, valueStr, dataStr string) (string, error) {
	if c == nil {
		return "", nil
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
	nonce, err := c.chain.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return "", err
	}

	gasPrice, err := c.chain.SuggestGasPrice(context.Background())
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
	gasLimit, err := c.chain.EstimateGas(context.Background(), msg)
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

	chainid, err := c.chain.ChainID(context.Background())
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainid), c.pvk)
	if err != nil {
		return "", err
	}
	err = c.chain.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil

}

func (c *ChainClient) CallContract(toStr, dataStr string) ([]byte, error) {
	var (
		to      = common.HexToAddress(toStr)
		data, _ = hex.DecodeString(dataStr)
	)

	msg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}

	return c.chain.CallContract(context.Background(), msg, nil)

}
