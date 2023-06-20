package monitor

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

const monitorBatchUpdateLimit = 100

type ControlContractLog struct {
	httpx.MethodPut
	blockchain.BatchUpdateMonitorReq `in:"body"`
	Cmd                              enums.MonitorCmd `in:"path" name:"cmd"`
}

func (r *ControlContractLog) Path() string { return "/contract_log/:cmd" }

func (r *ControlContractLog) Output(ctx context.Context) (interface{}, error) {
	if err := r.validate(ctx); err != nil {
		return nil, err
	}

	b, _ := convCmd(r.Cmd)
	return nil, blockchain.BatchUpdateContractLogPausedBySFIDs(ctx, r.IDs, b)
}

func (r *ControlContractLog) validate(ctx context.Context) error {
	if l := len(r.IDs); l == 0 || l > monitorBatchUpdateLimit {
		return status.InvalidContractLogIDs
	}

	if _, err := convCmd(r.Cmd); err != nil {
		return err
	}
	ca := middleware.MustCurrentAccountFromContext(ctx)

	cs, err := blockchain.ListContractLogBySFIDs(ctx, r.IDs)
	if err != nil {
		return err
	}
	if len(cs) == 0 {
		return status.InvalidContractLogIDs
	}
	pm := map[string]bool{}
	for _, c := range cs {
		pm[c.ProjectName] = true
	}

	return vaildateProjects(ctx, ca.AccountID, pm)
}

type ControlChainTx struct {
	httpx.MethodPut
	blockchain.BatchUpdateMonitorReq `in:"body"`
	Cmd                              enums.MonitorCmd `in:"path" name:"cmd"`
}

func (r *ControlChainTx) Path() string { return "/chain_tx/:cmd" }

func (r *ControlChainTx) Output(ctx context.Context) (interface{}, error) {
	if err := r.validate(ctx); err != nil {
		return nil, err
	}

	b, _ := convCmd(r.Cmd)
	return nil, blockchain.BatchUpdateChainTxPausedBySFIDs(ctx, r.IDs, b)
}

func (r *ControlChainTx) validate(ctx context.Context) error {
	if l := len(r.IDs); l == 0 || l > monitorBatchUpdateLimit {
		return status.InvalidChainTxIDs
	}

	if _, err := convCmd(r.Cmd); err != nil {
		return err
	}
	ca := middleware.MustCurrentAccountFromContext(ctx)

	cs, err := blockchain.ListChainTxBySFIDs(ctx, r.IDs)
	if err != nil {
		return err
	}
	if len(cs) == 0 {
		return status.InvalidChainTxIDs
	}
	pm := map[string]bool{}
	for _, c := range cs {
		pm[c.ProjectName] = true
	}

	return vaildateProjects(ctx, ca.AccountID, pm)
}

type ControlChainHeight struct {
	httpx.MethodPut
	blockchain.BatchUpdateMonitorReq `in:"body"`
	Cmd                              enums.MonitorCmd `in:"path" name:"cmd"`
}

func (r *ControlChainHeight) Path() string { return "/chain_height/:cmd" }

func (r *ControlChainHeight) Output(ctx context.Context) (interface{}, error) {
	if err := r.validate(ctx); err != nil {
		return nil, err
	}

	b, _ := convCmd(r.Cmd)
	return nil, blockchain.BatchUpdateChainHeightPausedBySFIDs(ctx, r.IDs, b)
}

func (r *ControlChainHeight) validate(ctx context.Context) error {
	if l := len(r.IDs); l == 0 || l > monitorBatchUpdateLimit {
		return status.InvalidChainHeightIDs
	}

	if _, err := convCmd(r.Cmd); err != nil {
		return err
	}
	ca := middleware.MustCurrentAccountFromContext(ctx)

	cs, err := blockchain.ListChainHeightBySFIDs(ctx, r.IDs)
	if err != nil {
		return err
	}
	if len(cs) == 0 {
		return status.InvalidChainHeightIDs
	}
	pm := map[string]bool{}
	for _, c := range cs {
		pm[c.ProjectName] = true
	}

	return vaildateProjects(ctx, ca.AccountID, pm)
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

func vaildateProjects(ctx context.Context, accountID types.SFID, pNames map[string]bool) error {
	pns := []string{}
	for p := range pNames {
		pns = append(pns, p)
	}
	ps, err := project.ListByCond(ctx, &project.CondArgs{Names: pns})
	if err != nil {
		return err
	}
	for _, p := range ps {
		if p.AccountID != accountID {
			return status.NoProjectPermission
		}
	}
	return nil
}
