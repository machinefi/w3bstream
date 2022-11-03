package datatypes

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

type OperationTimes struct {
	CreatedAt types.Timestamp `db:"f_created_at,default='0'" json:"createdAt"`
	UpdatedAt types.Timestamp `db:"f_updated_at,default='0'" json:"updatedAt"`
}

func (o *OperationTimes) Additions() builder.Additions {
	return builder.Additions{
		builder.OrderBy(builder.DescOrder(builder.Col("f_created_at"))),
		builder.OrderBy(builder.DescOrder(builder.Col("f_updated_at"))),
	}
}

type OperationTimesWithDeleted struct {
	OperationTimes
	DeletedAt types.Timestamp `db:"f_deleted_at,default='0'" json:"-"`
}

func (o *OperationTimesWithDeleted) Condition() builder.SqlCondition {
	return builder.Col("f_deleted_at").Eq(0)
}
