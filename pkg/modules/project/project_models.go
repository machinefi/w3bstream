package project

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CondArgs struct {
	ProjectIDs []types.SFID `in:"query" name:"projectID,omitempty"`
	Names      []string     `in:"query" name:"name,omitempty"`
	Versions   []string     `in:"query" name:"version,omitempty"`
}

func (r *CondArgs) Condition(acc types.SFID) builder.SqlCondition {
	m := &models.Project{}
	c := make([]builder.SqlCondition, 0)

	if acc != 0 {
		c = append(c, m.ColAccountID().Eq(acc))
	}
	if len(r.ProjectIDs) > 0 {
		c = append(c, m.ColProjectID().In(r.ProjectIDs))
	}
	if len(r.Names) > 0 {
		// TODO how to fix name's prefix?
		c = append(c, m.ColName().In(r.Names))
	}
	if len(r.Versions) > 0 {
		c = append(c, m.ColVersion().In(r.Versions))
	}
	c = append(c, m.ColDeletedAt().Eq(0))
	return builder.And(c...)
}

type ListReq struct {
	CondArgs
	datatypes.Pager
}

type ListRsp struct {
	Data  []models.Project `json:"data"`
	Total int64            `json:"total"`
}
