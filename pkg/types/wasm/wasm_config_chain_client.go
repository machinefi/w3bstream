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

	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	wsTypes "github.com/machinefi/w3bstream/pkg/types"
)

type ChainClient struct {
	projectName string
	endpoints   map[uint32]string
	clientMap   map[uint32]*ethclient.Client
	operators   map[string]*ecdsa.PrivateKey
}

// TODO impl ChainClient.Init

func NewChainClient(ctx context.Context, prj *models.Project, ops []models.Operator, p *models.ProjectOperator) *ChainClient {
	c := &ChainClient{
		projectName: prj.Name,
		clientMap:   make(map[uint32]*ethclient.Client, 0),
	}
	ethcli := wsTypes.MustETHClientConfigFromContext(ctx)
	c.endpoints = ethcli.Clients

	c.operators = convOperators(ops, p)
	return c
}

func convOperators(ops []models.Operator, p *models.ProjectOperator) map[string]*ecdsa.PrivateKey {
	res := make(map[string]*ecdsa.PrivateKey, len(ops))
	for _, op := range ops {
		res[op.Name] = crypto.ToECDSAUnsafe(common.FromHex(op.PrivateKey))
	}

	if p != nil {
		for _, op := range ops {
			if op.OperatorID == p.OperatorID {
				res[operator.DefaultOperatorName] = crypto.ToECDSAUnsafe(common.FromHex(op.PrivateKey))
				break
			}
		}
	}

	return res
}

func (c *ChainClient) SendTXWithOperator(chainID uint32, toStr, valueStr, dataStr, operatorName string) (string, error) {
	pvk, ok := c.operators[operatorName]
	if !ok {
		return "", errors.New("private key is empty")
	}
	return c.sendTX(chainID, toStr, valueStr, dataStr, pvk)
}

func (c *ChainClient) SendTX(chainID uint32, toStr, valueStr, dataStr string) (string, error) {
	pvk, ok := c.operators[operator.DefaultOperatorName]
	if !ok {
		return "", errors.New("private key is empty")
	}
	return c.sendTX(chainID, toStr, valueStr, dataStr, pvk)
}

func (c *ChainClient) sendTX(chainID uint32, toStr, valueStr, dataStr string, pvk *ecdsa.PrivateKey) (string, error) {
	cli, err := c.getEthClient(chainID)
	if err != nil {
		return "", err
	}
	var (
		sender = crypto.PubkeyToAddress(pvk.PublicKey)
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
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainid), pvk)
	if err != nil {
		return "", err
	}

	metrics.BlockChainTxMtc.WithLabelValues(c.projectName, strconv.Itoa(int(chainID))).Inc()

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

	metrics.BlockChainTxMtc.WithLabelValues(c.projectName, strconv.Itoa(int(chainID))).Inc()

	return cli.CallContract(context.Background(), msg, nil)
}
