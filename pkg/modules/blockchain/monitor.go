package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
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

	e := &eventpb.Event{
		Header: &eventpb.Header{
			EventType: eventType,
		},
		Payload: data,
	}
	body, err := json.Marshal(e)
	if err != nil {
		logger.Error(err)
		return err
	}
	url := fmt.Sprintf("http://localhost:8888/srv-applet-mgr/v0/event/%s", projectName) // TODO move to config
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		logger.Error(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer resp.Body.Close()

	// TODO http code judge
	return nil
}
