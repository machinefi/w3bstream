package wasm

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/blocto/solana-go-sdk/client"
	solcommon "github.com/blocto/solana-go-sdk/common"
	soltypes "github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/types"
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

type PrivateKey struct {
	Operator *models.Operator
	Type     enums.OperatorKeyType
	Ecdsa    *ecdsa.PrivateKey
	Ed25519  ed25519.PrivateKey
}

type ChainClient struct {
	ProjectName string
	Operators   map[string]*PrivateKey
}

func (c *ChainClient) GlobalConfigType() ConfigType { return ConfigChains }

func (c *ChainClient) Init(parent context.Context) error {
	prj := wsTypes.MustProjectFromContext(parent)
	ops := wsTypes.MustOperatorsFromContext(parent)

	c.ProjectName = prj.Name
	if c.Operators == nil {
		c.Operators = make(map[string]*PrivateKey)
	}

	defaultOpID := base.SFID(0)
	if op, ok := wsTypes.ProjectOperatorFromContext(parent); ok {
		defaultOpID = op.OperatorID
	}

	for _, op := range ops {
		p := &PrivateKey{Type: op.Type, Operator: &op}
		b := common.FromHex(op.PrivateKey)

		if op.Type == enums.OPERATOR_KEY__ED25519 {
			pk := ed25519.PrivateKey(b)
			p.Ed25519 = pk
		} else {
			pk := crypto.ToECDSAUnsafe(b)
			p.Ecdsa = pk
		}

		c.Operators[op.Name] = p
		if defaultOpID == op.OperatorID {
			c.Operators[operator.DefaultOperatorName] = p
		}
	}

	return nil
}

func (c *ChainClient) WithContext(ctx context.Context) context.Context {
	return WithChainClient(ctx, c)
}

func (c *ChainClient) SendUserOpWithOperator(conf *types.ChainConfig, chainID uint64, chainName enums.ChainName, dataStr, operatorName string) (string, error) {
	pvk, ok := c.Operators[operatorName]
	if !ok {
		return "", errors.New("private key is empty")
	}

	chain, ok := conf.GetChain(chainID, chainName)
	if !ok {
		return "", errors.Errorf("the chain %d %s is not supported", chainID, chainName)
	}
	if !chain.IsUserOpSupported() {
		return "", errors.New("the chain not support user operation")
	}
	if pvk.Type != enums.OPERATOR_KEY__ECDSA {
		return "", errors.New("invalid operator key type, require ECDSA")
	}
	req, err := http.NewRequest("POST", conf.AAUserOpEndpoint, nil)
	if err != nil {
		return "", errors.Wrap(err, "call user operation service failed")
	}
	defer req.Body.Close()

	req.Header.Set("data", dataStr)
	req.Header.Set("chainEndpoint", chain.Endpoint)
	req.Header.Set("privateKey", pvk.Operator.PrivateKey)
	req.Header.Set("bundlerEndpoint", conf.AABundlerEndpoint)
	req.Header.Set("paymasterEndpoint", fmt.Sprintf("%s/%s", conf.AABundlerEndpoint, pvk.Operator.PaymasterKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "call user operation service failed")
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return "", errors.Wrap(err, "read user operation service response failed")
	}
	ss := strings.Split(string(body), "Transaction hash:")
	s := strings.Split(ss[1], "Done")
	return strings.TrimSpace(s[0]), nil
}

func (c *ChainClient) SendTXWithOperator(conf *types.ChainConfig, chainID uint64, chainName enums.ChainName, toStr, valueStr, dataStr, operatorName string) (string, error) {
	pvk, ok := c.Operators[operatorName]
	if !ok {
		return "", errors.New("private key is empty")
	}
	return c.sendTX(conf, chainID, chainName, toStr, valueStr, dataStr, pvk)
}

func (c *ChainClient) SendTX(conf *types.ChainConfig, chainID uint64, chainName enums.ChainName, toStr, valueStr, dataStr string) (string, error) {
	pvk, ok := c.Operators[operator.DefaultOperatorName]
	if !ok {
		return "", errors.New("private key is empty")
	}
	return c.sendTX(conf, chainID, chainName, toStr, valueStr, dataStr, pvk)
}

func (c *ChainClient) sendTX(conf *types.ChainConfig, chainID uint64, chainName enums.ChainName, toStr, valueStr, dataStr string, pvk *PrivateKey) (string, error) {
	chain, ok := conf.GetChain(chainID, chainName)
	if !ok {
		return "", errors.Errorf("the chain %d %s is not supported", chainID, chainName)
	}
	if chain.IsSolana() {
		if pvk.Type != enums.OPERATOR_KEY__ED25519 {
			return "", errors.New("invalid operator key type, require ED25519")
		}
		return c.sendSolanaTX(chain, dataStr, pvk.Ed25519)
	}

	if pvk.Type != enums.OPERATOR_KEY__ECDSA {
		return "", errors.New("invalid operator key type, require ECDSA")
	}
	return c.sendEthTX(chain, toStr, valueStr, dataStr, pvk.Ecdsa)
}

func (c *ChainClient) sendSolanaTX(chain *types.Chain, dataStr string, pvk ed25519.PrivateKey) (string, error) {
	cli := client.NewClient(chain.Endpoint)
	account := soltypes.Account{
		PublicKey:  solcommon.PublicKeyFromBytes(pvk.Public().(ed25519.PublicKey)),
		PrivateKey: pvk,
	}
	ins := []soltypes.Instruction{}
	if err := json.Unmarshal([]byte(dataStr), &ins); err != nil {
		return "", errors.Wrap(err, "invalid data format")
	}
	if len(ins) == 0 {
		return "", errors.New("missing instruction data")
	}

	resp, err := cli.GetLatestBlockhash(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "failed to get solana latest block hash")
	}
	tx, err := soltypes.NewTransaction(soltypes.NewTransactionParam{
		Message: soltypes.NewMessage(soltypes.NewMessageParam{
			FeePayer:        account.PublicKey,
			RecentBlockhash: resp.Blockhash,
			Instructions:    ins,
		}),
		Signers: []soltypes.Account{account},
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to build solana raw tx")
	}
	hash, err := cli.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", errors.Wrap(err, "failed to send solana tx")
	}
	return hash, nil
}

func (c *ChainClient) sendEthTX(chain *types.Chain, toStr, valueStr, dataStr string, pvk *ecdsa.PrivateKey) (string, error) {
	if toStr == "" || valueStr == "" {
		return "", errors.New("missing to or value string")
	}
	cli, err := ethclient.Dial(chain.Endpoint)
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
	tx := ethtypes.NewTx(
		&ethtypes.LegacyTx{
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
	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewLondonSigner(chainid), pvk)
	if err != nil {
		return "", err
	}

	metrics.BlockChainTxMtc.WithLabelValues(c.ProjectName, strconv.Itoa(int(chain.ChainID))).Inc()

	err = cli.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}

func (c *ChainClient) getEthClient(conf *types.ChainConfig, chainID uint64, chainName enums.ChainName) (*ethclient.Client, error) {
	chain, ok := conf.GetChain(chainID, chainName)
	if !ok {
		return nil, errors.Errorf("the chain %d %s is not supported", chainID, chainName)
	}

	return ethclient.Dial(chain.Endpoint)
}

func (c *ChainClient) CallContract(conf *types.ChainConfig, chainID uint64, chainName enums.ChainName, toStr, dataStr string) ([]byte, error) {
	var (
		to = common.HexToAddress(toStr)
	)
	data, err := hex.DecodeString(strings.TrimPrefix(dataStr, "0x"))
	if err != nil {
		return nil, err
	}
	cli, err := c.getEthClient(conf, chainID, chainName)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}

	metrics.BlockChainTxMtc.WithLabelValues(c.ProjectName, strconv.Itoa(int(chainID))).Inc()

	return cli.CallContract(context.Background(), msg, nil)
}
