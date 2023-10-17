package clients

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/machinefi/w3bstream/pkg/enums"
	optypes "github.com/machinefi/w3bstream/pkg/modules/operator/pool/types"
	"github.com/pkg/errors"
)

type EthClient struct {
	endpoint string
}

func NewEthClient(endpoint string) *EthClient {
	return &EthClient{
		endpoint: endpoint,
	}
}

func (c *EthClient) TransactionByHash(ctx context.Context, hash string) (any, error) {
	client, err := ethclient.Dial(c.endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "dial chain address failed")
	}
	tx, _, err := client.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, errors.Wrap(err, "query transaction failed")
	}
	return tx, nil
}

func (c *EthClient) TransactionState(ctx context.Context, hash string) (enums.TransactionState, error) {
	client, err := ethclient.Dial(c.endpoint)
	if err != nil {
		return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "dial chain failed")
	}
	nh := common.HexToHash(hash)

	_, p, err := client.TransactionByHash(ctx, nh)
	if err != nil {
		if err == ethereum.NotFound {
			return enums.TRANSACTION_STATE__FAILED, nil
		} else {
			return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "get transaction by hash failed")
		}
	} else {
		if p {
			return enums.TRANSACTION_STATE__PENDING, nil
		}
	}

	receipt, err := client.TransactionReceipt(ctx, nh)
	if err != nil {
		if err == ethereum.NotFound {
			return enums.TRANSACTION_STATE__IN_BLOCK, nil
		}
		return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "get transaction receipt failed")
	}
	if receipt.Status == 0 {
		return enums.TRANSACTION_STATE__FAILED, nil
	}
	return enums.TRANSACTION_STATE__CONFIRMED, nil
}

func (c *EthClient) SendTransaction(ctx context.Context, toStr, valueStr, dataStr string, op *optypes.SyncOperator) (*ethtypes.Transaction, error) {
	cli, err := ethclient.Dial(c.endpoint)
	if err != nil {
		return nil, err
	}

	b := common.FromHex(op.Op.PrivateKey)
	pk := crypto.ToECDSAUnsafe(b)
	sender := crypto.PubkeyToAddress(pk.PublicKey)
	to := common.HexToAddress(toStr)

	value, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		return nil, errors.New("fail to read tx value")
	}
	data, err := hex.DecodeString(strings.TrimPrefix(dataStr, "0x"))
	if err != nil {
		return nil, err
	}

	gasPrice, err := cli.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		From:     sender,
		To:       &to,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}
	gasLimit, err := cli.EstimateGas(ctx, msg)
	if err != nil {
		return nil, err
	}

	chainid, err := cli.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	nonce, err := cli.PendingNonceAt(ctx, sender)
	if err != nil {
		return nil, err
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

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewLondonSigner(chainid), pk)
	if err != nil {
		return nil, err
	}

	err = cli.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
