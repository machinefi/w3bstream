package resource

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
)

type CondArgs struct {
	AccountID   types.SFID   `name:"-"`
	AppletIDs   []types.SFID `in:"query" name:"appletID,omitempty"`
	ResourceIDs []types.SFID `in:"query" name:"resourceID,omitempty"`
	FileNames   []string     `in:"query" name:"fileName,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m  = &models.ResourceMeta{}
		cs []builder.SqlCondition
	)

	if r.AccountID != 0 {
		cs = append(cs, m.ColAccountID().Eq(r.AccountID))
	}
	if len(r.AppletIDs) > 0 {
		cs = append(cs, m.ColAppletID().In(r.AppletIDs))
	}
	if len(r.ResourceIDs) > 0 {
		cs = append(cs, m.ColResourceID().In(r.ResourceIDs))
	}
	if len(r.FileNames) > 0 {
		cs = append(cs, m.ColFileName().In(r.FileNames))
	}
	cs = append(cs, m.ColDeletedAt().Neq(0))

	return builder.And(cs...)
}

type RemoveByAppletIDReq struct {
	AppletID types.SFID `json:"appletID"`
}

type ListReq struct {
	CondArgs
	datatypes.Pager
}

func (r *ListReq) Additions() builder.Additions {
	m := &models.ResourceMeta{}
	return builder.Additions{
		builder.OrderBy(
			builder.DescOrder(m.ColUpdatedAt()),
			builder.DescOrder(m.ColCreatedAt()),
		),
		r.Pager.Addition(),
	}
}

type ListRsp struct {
	Data  []models.ResourceMeta `json:"data"`
	Total int64                 `json:"total"`
}

type ListDetailRsp struct {
	Data  []*Detail `json:"data"`  // Data resource meta data list
	Total int64     `json:"total"` // Total resource meta count under current accountID
}

type Detail struct {
	AccountID   types.SFID `json:"accountID"   db:"f_acc_id"`
	ResourceID  types.SFID `json:"resourceID"  db:"f_res_id"`
	ProjectID   types.SFID `json:"projectID"   db:"f_prj_id"`
	ProjectName string     `json:"projectName" db:"f_prj_name"`
	AppletID    types.SFID `json:"appletID"    db:"f_app_id"`
	AppletName  string     `json:"appletName"  db:"f_app_name"`
	FileName    string     `json:"fileName"    db:"f_file_name"`
	datatypes.OperationTimes
}
