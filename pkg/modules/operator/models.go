package operator

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/models"
)

type CreateReq struct {
	AccountID  types.SFID `json:"-"`
	Name       string     `json:"name"`
	PrivateKey string     `json:"privateKey"`
}

type CondArgs struct {
	AccountID types.SFID `name:"-"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m  = &models.Operator{}
		cs []builder.SqlCondition
	)

	if r.AccountID != 0 {
		cs = append(cs, m.ColAccountID().Eq(r.AccountID))
	}
	cs = append(cs, m.ColDeletedAt().Eq(0))
	return builder.And(cs...)
}

type ListReq struct {
	CondArgs
}

type ListRsp struct {
	Data  []models.Operator `json:"data"`
	Total int64             `json:"total"`
}
