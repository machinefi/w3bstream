package publisher

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
)

type CreateReq struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type UpdateReq = CreateReq

type ListReq struct {
	PublisherIDs []types.SFID `in:"query" name:"publisherIDs"`
	Names        []string     `in:"query" name:"name"`
	Keys         []string     `in:"query" name:"key"`
	datatypes.Pager
}

func (r *ListReq) Condition(prj types.SFID) builder.SqlCondition {
	var (
		m  = &models.Publisher{}
		cs []builder.SqlCondition
	)

	if prj != 0 {
		cs = append(cs, m.ColProjectID().Eq(prj))
	}
	if len(r.PublisherIDs) > 0 {
		cs = append(cs, m.ColPublisherID().In(r.PublisherIDs))
	}
	if len(r.Names) > 0 {
		cs = append(cs, m.ColName().In(r.Names))
	}
	if len(r.Keys) > 0 {
		cs = append(cs, m.ColKey().In(r.Keys))
	}

	return builder.And(cs...)
}

func (r *ListReq) Additions() builder.Additions {
	m := &models.Publisher{}
	return builder.Additions{
		builder.OrderBy(
			builder.DescOrder(m.ColUpdatedAt()),
			builder.DescOrder(m.ColCreatedAt()),
		),
		r.Pager.Addition(),
	}
}

type ListRsp struct {
	Data  []models.Publisher `json:"data"`
	Total int64              `json:"total"`
}

type Detail struct {
	ProjectName string `json:"projectName" db:"f_project_name"`
	models.Publisher
	datatypes.OperationTimes
}

type ListDetailRsp struct {
	Total int64     `json:"total"`
	Data  []*Detail `json:"data"`
}
