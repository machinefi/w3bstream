package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"

	"github.com/iotexproject/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

// TODO move to config
const (
	listInterval  = 3 * time.Second
	blockInterval = 1000
)

func InitChainDB(ctx context.Context) error {
	d := types.MustDBExecutorFromContext(ctx)

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

func (l *monitor) sendEvent(data []byte, projectName string, et enums.EventType) error {
	// TODO event type
	e := &eventpb.Event{
		Payload: string(data),
	}
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://localhost:8888/srv-applet-mgr/v0/event/%s", projectName) // TODO move to config
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO http code judge
	return nil
}
