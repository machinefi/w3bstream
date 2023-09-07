package transaction

import (
	"context"
	"time"

	solclient "github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func RunTxSyncer(ctx context.Context) {
	s := &syncer{interval: 5 * time.Second}
	go s.run(ctx)
}

type syncer struct {
	interval time.Duration
}

func (s *syncer) run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		s.do(ctx)
	}
}

func (s *syncer) do(ctx context.Context) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	chainConf := types.MustChainConfigFromContext(ctx)

	m := &models.Transaction{}

	_, l = l.Start(ctx, "transaction.syncer")
	defer l.End()

	cs, err := m.List(d, builder.And(m.ColState().Neq(enums.TRANSACTION_STATE__CONFIRMED), m.ColState().Neq(enums.TRANSACTION_STATE__FAILED)))
	if err != nil {
		l.Error(errors.Wrap(err, "list transaction db failed"))
		return
	}
	for _, c := range cs {
		chain, ok := chainConf.Chains[c.ChainName]
		if !ok {
			l.WithValues("chainName", c.ChainName).Error(errors.New("blockchain not exist"))
			continue
		}
		state := c.State
		if chain.IsSolana() {
			state, err = s.getSolanaState(chain, c.Hash)
			if err != nil {
				l.WithValues("chainName", c.ChainName).Error(errors.Wrap(err, "get solana state failed"))
				continue
			}
		} else {
			state, err = s.getEthState(chain, c.Hash)
			if err != nil {
				l.WithValues("chainName", c.ChainName).Error(errors.Wrap(err, "get eth state failed"))
				continue
			}
		}

		if state != c.State {
			c.State = state
			if err := c.UpdateByID(d); err != nil {
				l.WithValues("transactionID", c.TransactionID).Error(errors.Wrap(err, "update db failed"))
				continue
			}
		}
	}
}

func (s *syncer) getEthState(chain *types.Chain, hash string) (enums.TransactionState, error) {
	client, err := ethclient.Dial(chain.Endpoint)
	if err != nil {
		return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "dial chain failed")
	}
	h := common.HexToHash(hash)

	_, p, err := client.TransactionByHash(context.Background(), h)
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

	receipt, err := client.TransactionReceipt(context.Background(), h)
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

func (s *syncer) getSolanaState(chain *types.Chain, hash string) (enums.TransactionState, error) {
	cli := solclient.NewClient(chain.Endpoint)
	status, err := cli.GetSignatureStatus(context.Background(), hash)
	if err != nil {
		return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "query solana transaction failed")
	}
	switch *status.ConfirmationStatus {
	case rpc.CommitmentProcessed:
		return enums.TRANSACTION_STATE__PENDING, nil
	case rpc.CommitmentConfirmed:
		return enums.TRANSACTION_STATE__IN_BLOCK, nil
	case rpc.CommitmentFinalized:
		return enums.TRANSACTION_STATE__CONFIRMED, nil
	}
	return enums.TRANSACTION_STATE_UNKNOWN, errors.New("get solana transaction status failed")
}
