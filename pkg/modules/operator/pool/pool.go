package pool

import (
	"context"
	"fmt"
	"sync"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	optypes "github.com/machinefi/w3bstream/pkg/modules/operator/pool/types"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Pool struct {
	db        sqlx.DBExecutor
	mux       sync.RWMutex
	operators map[string]*optypes.SyncOperator
}

func (p *Pool) getKey(accountID types.SFID, opName string) string {
	return fmt.Sprintf("%d-%s", accountID, opName)
}

func (p *Pool) Get(ctx context.Context, accountID types.SFID, opName string) (*optypes.SyncOperator, error) {
	key := p.getKey(accountID, opName)

	p.mux.RLock()
	op, ok := p.operators[key]
	p.mux.RUnlock()

	if ok {
		return op, nil
	}

	return p.setOperator(ctx, accountID, opName)
}

func (p *Pool) setOperator(ctx context.Context, accountID types.SFID, opName string) (*optypes.SyncOperator, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	key := p.getKey(accountID, opName)
	sop, ok := p.operators[key]
	if ok {
		return sop, nil
	}

	op, err := operator.GetByAccountAndName(types.WithMgrDBExecutor(ctx, p.db), accountID, opName)
	if err != nil {
		return nil, err
	}
	nsop := &optypes.SyncOperator{
		Op: op,
	}
	p.operators[key] = nsop
	return nsop, nil
}

// operator memory pool
// TODO support operator delete
func NewPool(mgrDB sqlx.DBExecutor) optypes.Pool {
	return &Pool{
		db:        mgrDB,
		operators: make(map[string]*optypes.SyncOperator),
	}
}
