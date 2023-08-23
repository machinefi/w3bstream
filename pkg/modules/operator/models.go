package operator

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
)

type CreateReq = models.OperatorInfo

type CondArgs struct {
	AccountID types.SFID `name:"-"`
}

type Detail struct {
	models.Operator
	Address string `json:"address"`
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
	datatypes.Pager
}

func (r *ListReq) Additions() builder.Additions {
	m := &models.Operator{}
	return builder.Additions{
		builder.OrderBy(
			builder.DescOrder(m.ColUpdatedAt()),
			builder.DescOrder(m.ColCreatedAt()),
		),
		r.Pager.Addition(),
	}
}

type ListRsp struct {
	Data  []models.Operator `json:"data"`
	Total int64             `json:"total"`
}

type ListDetailRsp struct {
	Data  []Detail `json:"data"`
	Total int64    `json:"total"`
}
