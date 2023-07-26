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

	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	wsTypes "github.com/machinefi/w3bstream/pkg/types"
)

func NewChainClient(ctx context.Context, prj *models.Project, ops []models.Operator, op *models.ProjectOperator) *ChainClient {
	ctx = contextx.WithContextCompose(
		wsTypes.WithProjectContext(prj),
		wsTypes.WithOperatorsContext(ops),
		wsTypes.WithProjectOperatorContext(op),
	)(ctx)

	cli := &ChainClient{}
	_ = cli.Init(ctx)
	return cli
}

type ChainClient struct {
	projectName string
	endpoints   map[uint32]string
	clientMap   map[uint32]*ethclient.Client
	operators   map[string]*ecdsa.PrivateKey
}

func (c *ChainClient) GlobalConfigType() ConfigType { return ConfigChains }

func (c *ChainClient) Init(parent context.Context) error {
	prj := wsTypes.MustProjectFromContext(parent)
	ops := wsTypes.MustOperatorsFromContext(parent)

	c.projectName = prj.Name
	if c.clientMap == nil {
		c.clientMap = make(map[uint32]*ethclient.Client)
	}
	if c.operators == nil {
		c.operators = make(map[string]*ecdsa.PrivateKey)
	}

	defaultOpID := base.SFID(0)
	if op, ok := wsTypes.ProjectOperatorFromContext(parent); ok {
		defaultOpID = op.OperatorID
	}

	for _, op := range ops {
		pk := crypto.ToECDSAUnsafe(common.FromHex(op.PrivateKey))
		c.operators[op.Name] = pk
		if defaultOpID == op.OperatorID {
			c.operators[operator.DefaultOperatorName] = pk
		}
	}

	ethcli := wsTypes.MustETHClientConfigFromContext(parent)
	c.endpoints = ethcli.Clients

	return nil
}

func (c *ChainClient) WithContext(ctx context.Context) context.Context {
	return WithChainClient(ctx, c)
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
