package strategy

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CondArgs struct {
	ProjectID   types.SFID   `name:"-"`
	AppletIDs   []types.SFID `in:"query" name:"appletID,omitempty"`
	StrategyIDs []types.SFID `in:"query" name:"strategyID,omitempty"`
	EventTypes  []string     `in:"query" name:"eventType,omitempty"`
	Handlers    []string     `in:"query" name:"handler,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m  = &models.Strategy{}
		cs []builder.SqlCondition
	)

	if r.ProjectID != 0 {
		cs = append(cs, m.ColProjectID().Eq(r.ProjectID))
	}
	if len(r.AppletIDs) > 0 {
		cs = append(cs, m.ColAppletID().In(r.AppletIDs))
	}
	if len(r.StrategyIDs) > 0 {
		cs = append(cs, m.ColStrategyID().In(r.StrategyIDs))
	}
	if len(r.EventTypes) > 0 {
		cs = append(cs, m.ColEventType().In(r.EventTypes))
	}
	if len(r.Handlers) > 0 {
		cs = append(cs, m.ColHandler().In(r.Handlers))
	}
	cs = append(cs, m.ColDeletedAt().Eq(0))

	return builder.And(cs...)
}

type ListReq struct {
	CondArgs
	datatypes.Pager
}

func (r *ListReq) Additions() builder.Additions {
	m := &models.Strategy{}
	return builder.Additions{
		builder.OrderBy(
			builder.DescOrder(m.ColUpdatedAt()),
			builder.DescOrder(m.ColCreatedAt()),
		),
		r.Pager.Addition(),
	}
}

type ListRsp struct {
	Data  []models.Strategy `json:"data"`
	Total int64             `json:"total"`
}

type ListDetailRsp struct {
	Data  []*Detail `json:"data"`  // Data strategy data list
	Total int64     `json:"total"` // Total strategy count under current projectID
}

type Detail struct {
	types.StrategyResult
	datatypes.OperationTimes
}

type CreateReq struct {
	models.RelApplet
	models.StrategyInfo
}

type BatchCreateReq struct {
	Strategies []CreateReq `json:"strategies"`
}

type UpdateReq = CreateReq
