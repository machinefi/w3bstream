package project

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type CondArgs struct {
	AccountID  types.SFID   `name:"-"`
	ProjectIDs []types.SFID `in:"query" name:"projectID,omitempty"`
	Names      []string     `in:"query" name:"name,omitempty"`
	Versions   []string     `in:"query" name:"version,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	m := &models.Project{}
	c := make([]builder.SqlCondition, 0)

	if r.AccountID != 0 {
		c = append(c, m.ColAccountID().Eq(r.AccountID))
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

type Detail struct {
	ProjectID   types.SFID       `json:"projectID"`
	ProjectName string           `json:"projectName"`
	Applets     []*applet.Detail `json:"applets,omitempty"`
}

type ListDetailRsp struct {
	Data  []*Detail `json:"data"`
	Total int64     `json:"total"`
}

type CreateReq struct {
	models.ProjectName
	models.ProjectBase
	Env      *wasm.Env      `json:"envs,omitempty"`
	Database *wasm.Database `json:"database,omitempty"`
}

type CreateRsp struct {
	*models.Project
	Env          *wasm.Env      `json:"envs,omitempty"`
	Database     *wasm.Database `json:"database,omitempty"`
	ChannelState datatypes.Bool `json:"channelState"`
}
