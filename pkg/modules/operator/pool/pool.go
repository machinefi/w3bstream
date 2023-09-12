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

func (p *Pool) Get(accountID types.SFID, opName string) (*optypes.SyncOperator, error) {
	key := p.getKey(accountID, opName)

	p.mux.RLock()
	op, ok := p.operators[key]
	p.mux.RUnlock()

	if ok {
		return op, nil
	}

	return p.setOperator(accountID, opName)
}

func (p *Pool) setOperator(accountID types.SFID, opName string) (*optypes.SyncOperator, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	key := p.getKey(accountID, opName)
	sop, ok := p.operators[key]
	if ok {
		return sop, nil
	}

	op, err := operator.GetByAccountAndName(types.WithMgrDBExecutor(context.Background(), p.db), accountID, opName)
	if err != nil {
		return nil, err
	}
	nsop := &optypes.SyncOperator{
		Op: op,
	}
	p.operators[key] = nsop
	return nsop, nil
}

func (p *Pool) Delete(id types.SFID) error {
	ctx := types.WithMgrDBExecutor(context.Background(), p.db)
	op, err := operator.GetBySFID(ctx, id)
	if err != nil {
		return err
	}
	if err := operator.RemoveBySFID(ctx, id); err != nil {
		return err
	}

	key := p.getKey(op.AccountID, op.Name)

	p.mux.Lock()
	defer p.mux.Unlock()
	delete(p.operators, key)

	return nil
}

// operator memory pool
func NewPool(mgrDB sqlx.DBExecutor) optypes.Pool {
	return &Pool{
		db:        mgrDB,
		operators: make(map[string]*optypes.SyncOperator),
	}
}
