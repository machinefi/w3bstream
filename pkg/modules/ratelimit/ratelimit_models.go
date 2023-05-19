package ratelimit

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
)

type CreateReq struct {
	models.RateLimitInfo
}

type UpdateReq struct {
	RateLimitID types.SFID `json:"-"`
	models.RateLimitInfo
}

type CondArgs struct {
	ProjectID types.SFID             `in:"query" name:"projectID,omitempty"`
	ApiType   enums.RateLimitApiType `in:"query" name:"apiType,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m = &models.TrafficRateLimit{}
		c []builder.SqlCondition
	)
	if r.ProjectID != 0 {
		c = append(c, m.ColProjectID().Eq(r.ProjectID))
	}
	if r.ApiType != 0 {
		c = append(c, m.ColApiType().In(r.ApiType))
	}
	return builder.And(c...)
}

type ListReq struct {
	CondArgs
	datatypes.Pager
}

func (r *ListReq) Additions() builder.Additions {
	m := &models.TrafficRateLimit{}
	return builder.Additions{
		builder.OrderBy(
			builder.DescOrder(m.ColUpdatedAt()),
			builder.DescOrder(m.ColCreatedAt()),
		),
		r.Pager.Addition(),
	}
}

type ListRsp struct {
	Data  []models.TrafficRateLimit `json:"data"`
	Total int64                     `json:"total"`
}

type Detail struct {
	ProjectName string `json:"projectName" db:"f_project_name"`
	models.TrafficRateLimit
	datatypes.OperationTimes
}

type ListDetailRsp struct {
	Total int64     `json:"total"`
	Data  []*Detail `json:"data"`
}
