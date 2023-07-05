package deploy

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type CondArgs struct {
	ProjectID   types.SFID            `name:"-"`
	InstanceIDs []types.SFID          `in:"query" name:"instanceID,omitempty"`
	AppletIDs   []types.SFID          `in:"query" name:"appletID,omitempty"`
	States      []enums.InstanceState `in:"query" name:"state,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m = &models.Instance{}
		c []builder.SqlCondition
	)

	if r.ProjectID != 0 {
		c = append(c, (&models.Applet{}).ColProjectID().Eq(r.ProjectID))
	}
	if len(r.InstanceIDs) > 0 {
		if len(r.InstanceIDs) == 1 {
			c = append(c, m.ColInstanceID().Eq(r.InstanceIDs[0]))
		} else {
			c = append(c, m.ColInstanceID().In(r.InstanceIDs))
		}
	}
	if len(r.AppletIDs) > 0 {
		if len(r.AppletIDs) == 1 {
			c = append(c, m.ColAppletID().Eq(r.AppletIDs[0]))
		} else {
			c = append(c, m.ColAppletID().In(r.AppletIDs))
		}
	}
	if len(r.States) > 0 {
		if len(r.States) == 1 {
			c = append(c, m.ColState().Eq(r.States[0]))
		} else {
			c = append(c, m.ColState().In(r.States))
		}
	}
	return builder.And(c...)
}

type ListReq struct {
	CondArgs
	datatypes.Pager
}

type ListRsp struct {
	Data  []models.Instance `json:"data"`
	Total int64             `json:"total"`
}

type CreateReq struct {
	Cache *wasm.Cache `json:"cache,omitempty"`
}
