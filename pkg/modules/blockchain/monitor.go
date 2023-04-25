package blockchain

import (
	"context"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
)

// TODO move to config
const (
	listInterval  = 3 * time.Second
	blockInterval = 1000
)

func InitChainDB(ctx context.Context) error {
	d := types.MustMonitorDBExecutorFromContext(ctx)

	m := &models.Blockchain{
		RelBlockchain:  models.RelBlockchain{ChainID: 4690},
		BlockchainInfo: models.BlockchainInfo{Address: "https://babel-api.testnet.iotex.io"},
	}

	results := make([]models.Account, 0)
	err := d.QueryAndScan(builder.Select(nil).
		From(
			d.T(m),
			builder.Where(
				builder.And(
					m.ColChainID().Eq(4690),
				),
			),
		), &results)
	if err != nil {
		return status.CheckDatabaseError(err, "FetchChain")
	}
	if len(results) > 0 {
		return nil
	}
	return m.Create(d)
}

func Monitor(ctx context.Context) {
	m := &monitor{}
	c := &contract{
		monitor:       m,
		listInterval:  listInterval,
		blockInterval: blockInterval,
	}
	h := &height{
		monitor:  m,
		interval: listInterval,
	}
	t := &tx{
		monitor:  m,
		interval: listInterval,
	}
	go c.run(ctx)
	go h.run(ctx)
	go t.run(ctx)
}

type monitor struct{}

func (l *monitor) sendEvent(ctx context.Context, data []byte, projectName string, eventType string) error {
	logger := types.MustLoggerFromContext(ctx)

	_, logger = logger.Start(ctx, "monitor.sendEvent")
	defer logger.End()

	ret := &event.HandleEventResult{
		ProjectName: projectName,
	}
	return event.HandleEvent(ctx, projectName, eventType, ret, data)
}
