package wasm

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
	"github.com/machinefi/w3bstream/pkg/modules/blockchain/clients"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	optypes "github.com/machinefi/w3bstream/pkg/modules/operator/pool/types"
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
	Type    enums.OperatorKeyType
	Ecdsa   []byte
	Ed25519 ed25519.PrivateKey
}

type ChainClient struct {
	ProjectName string
	Operators   map[string]*PrivateKey
}

type SendTxResp struct {
	ChainName enums.ChainName
	Nonce     uint64
	Hash      string
	Sender    string
	Receiver  string
	Data      string
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
		p := &PrivateKey{Type: op.Type}
		b := common.FromHex(op.PrivateKey)

		if op.Type == enums.OPERATOR_KEY__ED25519 {
			pk := ed25519.PrivateKey(b)
			p.Ed25519 = pk
		} else {
			p.Ecdsa = b
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

func (c *ChainClient) SendTXWithOperator(conf *wsTypes.ChainConfig, chainID uint64, chainName enums.ChainName, toStr, valueStr, dataStr, operatorName string, opPool optypes.Pool, prj *models.Project) (*SendTxResp, error) {
	op, err := opPool.Get(prj.AccountID, operatorName)
	if err != nil {
		return nil, err
	}
	return c.sendTX(conf, chainID, chainName, toStr, valueStr, dataStr, op)
}

func (c *ChainClient) SendTX(conf *wsTypes.ChainConfig, chainID uint64, chainName enums.ChainName, toStr, valueStr, dataStr string, opPool optypes.Pool, prj *models.Project) (string, error) {
	op, err := opPool.Get(prj.AccountID, operator.DefaultOperatorName)
	if err != nil {
		return "", err
	}
	resp, err := c.sendTX(conf, chainID, chainName, toStr, valueStr, dataStr, op)
	if err != nil {
		return "", err
	}
	return resp.Hash, err
}

func (c *ChainClient) sendTX(conf *wsTypes.ChainConfig, chainID uint64, chainName enums.ChainName, toStr, valueStr, dataStr string, op *optypes.SyncOperator) (*SendTxResp, error) {
	chain, ok := conf.GetChain(chainID, chainName)
	if !ok {
		return nil, errors.Errorf("the chain %d %s is not supported", chainID, chainName)
	}
	if op.Op.PaymasterKey != "" {
		if !chain.IsAASupported() {
			return nil, errors.New("account abstraction not supported at the chain")
		}
		return c.sendUserOp(conf, chain, toStr, valueStr, dataStr, op)
	}
	if chain.IsSolana() {
		if op.Op.Type != enums.OPERATOR_KEY__ED25519 {
			return nil, errors.New("invalid operator key type, require ED25519")
		}
		return c.sendSolanaTX(chain, dataStr, op)
	}

	if op.Op.Type != enums.OPERATOR_KEY__ECDSA {
		return nil, errors.New("invalid operator key type, require ECDSA")
	}
	return c.sendEthTX(chain, toStr, valueStr, dataStr, op)
}

func (c *ChainClient) sendUserOp(conf *wsTypes.ChainConfig, chain *wsTypes.Chain, toStr, valueStr, dataStr string, op *optypes.SyncOperator) (*SendTxResp, error) {
	if toStr == "" || valueStr == "" {
		return nil, errors.New("missing to or value string")
	}

	op.Mux.Lock()
	defer op.Mux.Unlock()

	params, err := json.Marshal(struct {
		PrivateKey            string `json:"privateKey,omitempty"`
		To                    string `json:"to,omitempty"`
		Value                 string `json:"value,omitempty"`
		Data                  string `json:"data,omitempty"`
		ChainRPC              string `json:"chainRPC,omitempty"`
		BundlerRPC            string `json:"bundlerRPC,omitempty"`
		PaymasterRPC          string `json:"paymasterRPC,omitempty"`
		EntryPointAddress     string `json:"entryPointAddress,omitempty"`
		AccountFactoryAddress string `json:"accountFactoryAddress,omitempty"`
	}{
		PrivateKey:            op.Op.PrivateKey,
		To:                    toStr,
		Value:                 valueStr,
		Data:                  dataStr,
		ChainRPC:              chain.Endpoint,
		BundlerRPC:            chain.AABundlerEndpoint,
		PaymasterRPC:          fmt.Sprintf("%s/%s", chain.AAPaymasterEndpoint, op.Op.PaymasterKey),
		EntryPointAddress:     chain.AAEntryPointContractAddress,
		AccountFactoryAddress: chain.AAAccountFactoryContractAddress,
	})
	if err != nil {
		return nil, errors.Wrap(err, "build aa service params failed")
	}

	req, err := http.NewRequest("POST", conf.AAUserOpEndpoint, bytes.NewReader(params))
	if err != nil {
		return nil, errors.Wrap(err, "build aa service http request failed")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "call aa service failed")
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, errors.Wrap(err, "read aa service response failed")
	}
	jsonResp := struct {
		TxHash string `json:"txHash,omitempty"`
	}{}
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return nil, errors.Wrap(err, "unmarshal aa service response failed")
	}

	b := common.FromHex(op.Op.PrivateKey)
	pk := crypto.ToECDSAUnsafe(b)
	sender := crypto.PubkeyToAddress(pk.PublicKey)

	return &SendTxResp{
		ChainName: chain.Name,
		Sender:    sender.String(),
		Hash:      jsonResp.TxHash,
		Receiver:  toStr,
		Data:      dataStr,
	}, nil
}

func (c *ChainClient) sendSolanaTX(chain *wsTypes.Chain, dataStr string, op *optypes.SyncOperator) (*SendTxResp, error) {
	cli := client.NewClient(chain.Endpoint)
	b := common.FromHex(op.Op.PrivateKey)
	pk := ed25519.PrivateKey(b)
	account := soltypes.Account{
		PublicKey:  solcommon.PublicKeyFromBytes(pk.Public().(ed25519.PublicKey)),
		PrivateKey: pk,
	}
	ins := []soltypes.Instruction{}
	if err := json.Unmarshal([]byte(dataStr), &ins); err != nil {
		return nil, errors.Wrap(err, "invalid data format")
	}
	if len(ins) == 0 {
		return nil, errors.New("missing instruction data")
	}

	resp, err := cli.GetLatestBlockhash(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get solana latest block hash")
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
		return nil, errors.Wrap(err, "failed to build solana raw tx")
	}

	op.Mux.Lock()
	defer op.Mux.Unlock()

	hash, err := cli.SendTransaction(context.Background(), tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send solana tx")
	}
	return &SendTxResp{
		ChainName: chain.Name,
		Hash:      hash,
		Sender:    account.PublicKey.String(),
		Data:      dataStr,
	}, nil
}

func (c *ChainClient) sendEthTX(chain *wsTypes.Chain, toStr, valueStr, dataStr string, op *optypes.SyncOperator) (*SendTxResp, error) {
	if toStr == "" || valueStr == "" {
		return nil, errors.New("missing to or value string")
	}

	op.Mux.Lock()
	defer op.Mux.Unlock()

	b := common.FromHex(op.Op.PrivateKey)
	pk := crypto.ToECDSAUnsafe(b)
	sender := crypto.PubkeyToAddress(pk.PublicKey)
	client := NewEthClient(chain)
	tx, err := client.SendTransaction(context.Background(), toStr, valueStr, dataStr, op)
	if err != nil {
		return nil, err
	}
	return &SendTxResp{
		ChainName: chain.Name,
		Nonce:     tx.Nonce(),
		Hash:      tx.Hash().Hex(),
		Sender:    sender.String(),
		Receiver:  toStr,
		Data:      dataStr,
	}, nil
}

func (c *ChainClient) getEthClient(conf *wsTypes.ChainConfig, chainID uint64, chainName enums.ChainName) (*ethclient.Client, error) {
	chain, ok := conf.GetChain(chainID, chainName)
	if !ok {
		return nil, errors.Errorf("the chain %d %s is not supported", chainID, chainName)
	}

	return ethclient.Dial(chain.Endpoint)
}

func (c *ChainClient) CallContract(conf *wsTypes.ChainConfig, chainID uint64, chainName enums.ChainName, toStr, dataStr string) ([]byte, error) {
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

// TODO: move to a more appropriate place
type EthClient interface {
	TransactionByHash(ctx context.Context, hash string) (any, error)
	TransactionState(ctx context.Context, hash string) (enums.TransactionState, error)
	SendTransaction(ctx context.Context, toStr, valueStr, dataStr string, op *optypes.SyncOperator) (*ethtypes.Transaction, error)
}

// NewEthClient creates a new EthClient according to the chain type
func NewEthClient(chain *types.Chain) EthClient {
	if chain.IsZKSync() {
		return clients.NewZKSyncClient(chain.Endpoint)
	}
	return clients.NewEthClient(chain.Endpoint)
}
