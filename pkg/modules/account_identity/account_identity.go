package account_identity

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

// GetBySFIDAndType get account identity from model by account id and identity type
func GetBySFIDAndType(ctx context.Context, id types.SFID, t enums.AccountIdentityType) (*models.AccountIdentity, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.AccountIdentity{
		RelAccount:          models.RelAccount{AccountID: id},
		AccountIdentityInfo: models.AccountIdentityInfo{Type: t},
	}

	if err := m.FetchByAccountIDAndType(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.AccountIdentityNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}
