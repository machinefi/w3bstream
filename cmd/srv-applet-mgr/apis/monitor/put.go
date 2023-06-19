package monitor

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
)

type ControlContractLog struct {
	httpx.MethodPut
	blockchain.UpdateMonitorReq `in:"body"`
	Cmd                         enums.MonitorCmd `in:"path" name:"cmd"`
}

func (r *ControlContractLog) Path() string { return "/contract_log/:cmd" }

func (r *ControlContractLog) Output(ctx context.Context) (interface{}, error) {
	if len(r.IDs) == 0 {
		return nil, status.BadRequest
	}
	ca := middleware.MustCurrentAccountFromContext(ctx)

	for _, id := range r.IDs {
		if _, err := ca.WithContractLogBySFID(ctx, id); err != nil {
			return nil, err
		}
	}

	b, err := convCmd(r.Cmd)
	if err != nil {
		return nil, err
	}
	return nil, blockchain.UpdateContractLogPausedBySFIDs(ctx, r.IDs, b)
}

type ControlChainTx struct {
	httpx.MethodPut
	blockchain.UpdateMonitorReq `in:"body"`
	Cmd                         enums.MonitorCmd `in:"path" name:"cmd"`
}

func (r *ControlChainTx) Path() string { return "/chain_tx/:cmd" }

func (r *ControlChainTx) Output(ctx context.Context) (interface{}, error) {
	if len(r.IDs) == 0 {
		return nil, status.BadRequest
	}
	ca := middleware.MustCurrentAccountFromContext(ctx)

	for _, id := range r.IDs {
		if _, err := ca.WithChainTxBySFID(ctx, id); err != nil {
			return nil, err
		}
	}

	b, err := convCmd(r.Cmd)
	if err != nil {
		return nil, err
	}
	return nil, blockchain.UpdateChainTxPausedBySFIDs(ctx, r.IDs, b)
}

type ControlChainHeight struct {
	httpx.MethodPut
	blockchain.UpdateMonitorReq `in:"body"`
	Cmd                         enums.MonitorCmd `in:"path" name:"cmd"`
}

func (r *ControlChainHeight) Path() string { return "/chain_height/:cmd" }

func (r *ControlChainHeight) Output(ctx context.Context) (interface{}, error) {
	if len(r.IDs) == 0 {
		return nil, status.BadRequest
	}
	ca := middleware.MustCurrentAccountFromContext(ctx)

	for _, id := range r.IDs {
		if _, err := ca.WithChainHeightBySFID(ctx, id); err != nil {
			return nil, err
		}
	}

	b, err := convCmd(r.Cmd)
	if err != nil {
		return nil, err
	}
	return nil, blockchain.UpdateChainHeightPausedBySFIDs(ctx, r.IDs, b)
}

func convCmd(c enums.MonitorCmd) (datatypes.Bool, error) {
	switch c {
	case enums.MONITOR_CMD__START:
		return datatypes.FALSE, nil
	case enums.MONITOR_CMD__PAUSE:
		return datatypes.TRUE, nil
	}
	return datatypes.FALSE, status.UnknownMonitorCommand
}
