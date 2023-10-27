package clients

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/zksync-sdk/zksync2-go/clients"

	"github.com/machinefi/w3bstream/pkg/enums"
)

type (
	// ZKSyncClient is a client for zksync chain
	ZKSyncClient struct {
		*EthClient
	}
)

// NewZKSyncClient creates a new ZKSyncClient
func NewZKSyncClient(endpoint string) *ZKSyncClient {
	return &ZKSyncClient{
		EthClient: NewEthClient(endpoint),
	}
}

// TransactionByHash returns transaction by hash
func (c *ZKSyncClient) TransactionByHash(ctx context.Context, hash string) (any, error) {
	client, err := clients.Dial(c.endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "dial chain address failed")
	}

	tx, _, err := client.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, errors.Wrap(err, "query transaction failed")
	}
	return tx, nil
}

// TransactionState returns transaction state
func (c *ZKSyncClient) TransactionState(ctx context.Context, hash string) (enums.TransactionState, error) {
	client, err := clients.Dial(c.endpoint)
	if err != nil {
		return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "dial chain failed")
	}
	nh := common.HexToHash(hash)

	_, p, err := client.TransactionByHash(ctx, nh)
	if err != nil {
		if err == ethereum.NotFound {
			return enums.TRANSACTION_STATE__FAILED, nil
		}
		return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "get transaction by hash failed")
	}
	if p {
		return enums.TRANSACTION_STATE__PENDING, nil
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
