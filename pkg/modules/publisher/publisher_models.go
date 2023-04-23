package publisher

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
)

type CreateReq struct {
	ProjectID types.SFID `json:"-"`
	Name      string     `json:"name"`
	Key       string     `json:"key"`
}

type UpdateReq struct {
	PublisherID types.SFID `json:"-"`
	Name        string     `json:"name"`
	Key         string     `json:"key"`
}

type CondArgs struct {
	ProjectIDs   []types.SFID `name:"-"`
	PublisherIDs []types.SFID `in:"query" name:"publisherIDs"`
	Names        []string     `in:"query" name:"name"`
	Keys         []string     `in:"query" name:"key"`
	NameLike     string       `in:"query" name:"name"`
	LNameLike    string       `in:"query" name:"lname"`
	RNameLike    string       `in:"query" name:"rname"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m = &models.Publisher{}
		c []builder.SqlCondition
	)

	if len(r.ProjectIDs) > 0 {
		if len(r.ProjectIDs) == 1 {
			c = append(c, m.ColProjectID().Eq(r.ProjectIDs[0]))
		} else {
			c = append(c, m.ColProjectID().In(r.ProjectIDs))
		}
	}
	if len(r.PublisherIDs) > 0 {
		if len(r.PublisherIDs) == 1 {
			c = append(c, m.ColPublisherID().Eq(r.PublisherIDs[0]))
		} else {
			c = append(c, m.ColPublisherID().In(r.PublisherIDs))
		}
	}
	if len(r.Names) > 0 {
		if len(r.PublisherIDs) == 1 {
			c = append(c, m.ColName().Eq(r.Names[0]))
		} else {
			c = append(c, m.ColName().In(r.Names))
		}
	}
	if len(r.Keys) > 0 {
		if len(r.PublisherIDs) == 1 {
			c = append(c, m.ColKey().Eq(r.Keys[0]))
		} else {
			c = append(c, m.ColKey().In(r.Keys))
		}
	}
	if r.NameLike != "" {
		c = append(c, m.ColName().Like(r.NameLike))
	}
	if r.LNameLike != "" {
		c = append(c, m.ColName().LLike(r.LNameLike))
	}
	if r.RNameLike != "" {
		c = append(c, m.ColName().RLike(r.RNameLike))
	}

	return builder.And(c...)
}

type ListReq struct {
	CondArgs
	datatypes.Pager
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
