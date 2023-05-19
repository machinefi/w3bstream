package ratelimit

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/models"
)

type CreateReq struct {
	models.RateLimitInfo
}

type UpdateReq = CreateReq

type CondArgs struct {
	ProjectName types.SFID   `name:"-"`
	AppletIDs   []types.SFID `in:"query" name:"appletID,omitempty"`
	Names       []string     `in:"query" name:"names,omitempty"`
	NameLike    string       `in:"query" name:"name,omitempty"`
	LNameLike   string       `in:"query" name:"lName,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m = &models.Applet{}
		c []builder.SqlCondition
	)
	if r.ProjectID != 0 {
		c = append(c, m.ColProjectID().Eq(r.ProjectID))
	}
	if len(r.AppletIDs) > 0 {
		c = append(c, m.ColAppletID().In(r.AppletIDs))
	}
	if len(r.Names) > 0 {
		c = append(c, m.ColName().In(r.Names))
	}
	if r.NameLike != "" {
		c = append(c, m.ColName().Like(r.NameLike))
	}
	if r.NameLike != "" {
		c = append(c, m.ColName().LLike(r.LNameLike))
	}
	return builder.And(c...)
}
