package strategy

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ListReq struct {
	CondArgs
	datatypes.Pager
}

type ListRsp struct {
	Data  []models.Strategy `json:"data"`
	Total int64             `json:"total"`
}

type CondArgs struct {
	StrategyIDs []types.SFID `in:"query" name:"strategyID"`
	AppletIDs   []types.SFID `in:"query" name:"appletID"`
	EventTypes  []string     `in:"query" name:"eventType"`
	Handlers    []string     `in:"query" name:"handler"`
}

func (arg *CondArgs) Condition(prj types.SFID) builder.SqlCondition {
	if arg == nil {
		return nil
	}

	m := &models.Strategy{}
	c := make([]builder.SqlCondition, 0)

	if prj != 0 {
		c = append(c, m.ColProjectID().Eq(prj))
	}
	if len(arg.StrategyIDs) > 0 {
		c = append(c, m.ColStrategyID().In(arg.StrategyIDs))
	}
	if len(arg.AppletIDs) > 0 {
		c = append(c, m.ColAppletID().In(arg.AppletIDs))
	}
	if len(arg.EventTypes) > 0 {
		c = append(c, m.ColEventType().In(arg.EventTypes))
	}
	if len(arg.Handlers) > 0 {
		c = append(c, m.ColHandler().In(arg.Handlers))
	}
	c = append(c, m.ColDeletedAt().Eq(0))
	return builder.And(c...)
}

type CreateData struct {
	models.RelApplet
	models.StrategyInfo
}

type CreateReq struct {
	Data []CreateData `json:"data"`
}

type CreateRsp struct {
	Data []*models.Strategy `json:"data"`
}

type UpdateReq = CreateData
