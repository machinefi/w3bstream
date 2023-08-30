package types

import (
	"context"
	"sync"

	basetypes "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/models"
)

type SyncOperator struct {
	Mux sync.Mutex
	Op  *models.Operator
}

type Pool interface {
	Get(ctx context.Context, accountID basetypes.SFID, opName string) (*SyncOperator, error)
}
