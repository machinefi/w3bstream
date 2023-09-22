package blockchain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type height struct {
	*monitor
	interval time.Duration
}

func (h *height) run(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for range ticker.C {
		h.do(ctx)
	}
}

func (h *height) do(ctx context.Context) {
	d := types.MustMonitorDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	chainConf := types.MustChainConfigFromContext(ctx)
	m := &models.ChainHeight{}

	_, l = l.Start(ctx, "height.run")
	defer l.End()

	cs, err := m.List(d, builder.And(m.ColFinished().Eq(datatypes.FALSE), m.ColPaused().Eq(datatypes.FALSE)))
	if err != nil {
		l.Error(errors.Wrap(err, "list chain height db failed"))
		return
	}
	for _, c := range cs {
		chain, ok := chainConf.GetChain(c.ChainID, c.ChainName)
		if !ok {
			l.WithValues("chainID", c.ChainID, "chainName", c.ChainName).Error(errors.New("blockchain not exist"))
			continue
		}
		res, err := h.checkHeightAndSendEvent(ctx, &c, chain)
		if err != nil {
			l.Error(errors.Wrap(err, "check chain height and send event failed"))
			continue
		}
		if res {
			c.Finished = datatypes.TRUE
			c.Uniq = c.ChainHeightID
			if err := c.UpdateByID(d); err != nil {
				l.Error(errors.Wrap(err, "update chain height db failed"))
			}
		}
	}
}

func (h *height) checkHeightAndSendEvent(ctx context.Context, c *models.ChainHeight, chain *types.Chain) (bool, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "height.checkHeightAndSendEvent")
	defer l.End()

	l = l.WithValues("type", "chain_height", "chain_height_id", c.ChainHeightID)

	headerNumber, err := h.getHeaderNumber(ctx, chain)
	if err != nil {
		l.Error(err)
		return false, err
	}
	if headerNumber < c.Height {
		l.WithValues("headerNumber", headerNumber, "chainHeight", c.Height).Debug("did not arrive")
		return false, nil
	}
	data, err := json.Marshal(struct {
		HeaderNumber uint64
	}{
		headerNumber,
	})
	if err != nil {
		l.Error(err)
		return false, err
	}
	if err := h.sendEvent(ctx, data, c.ProjectName, c.EventType); err != nil {
		l.Error(err)
		return false, err
	}
	return true, nil
}

func (h *height) getHeaderNumber(ctx context.Context, chain *types.Chain) (uint64, error) {
	switch {
	case chain.IsEth():
		client, err := ethclient.Dial(chain.Endpoint)
		if err != nil {
			return 0, err
		}
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return 0, err
		}
		return header.Number.Uint64(), nil

	case chain.IsSolana():
		cli := client.NewClient(chain.Endpoint)
		res, err := cli.RpcClient.GetBlockHeightWithConfig(context.Background(), rpc.GetBlockHeightConfig{Commitment: rpc.CommitmentFinalized})
		if err != nil {
			return 0, err
		}
		if res.Error != nil {
			return 0, res.Error
		}
		return res.Result, nil

	default:
		return 0, errors.New("unsupported chain")
	}
}
