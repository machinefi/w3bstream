package applet

import (
	"context"
	"mime/multipart"
	"time"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type DataListParam struct {
	ProjectID types.SFID
	ListReq
}

func (r *DataListParam) Condition() builder.SqlCondition {
	return r.ListReq.Condition(r.ProjectID)
}

type ListReq struct {
	AppletIDs []types.SFID `in:"query" name:"appletID,omitempty"`
	Names     []string     `in:"query" name:"names,omitempty"`
	NameLike  string       `in:"query" name:"name,omitempty"`
	LNameLike string       `in:"query" name:"lName,omitempty"`
	datatypes.Pager
}

func (r *ListReq) Condition(prj types.SFID) builder.SqlCondition {
	var (
		m  = &models.Applet{}
		cs []builder.SqlCondition
	)
	if prj != 0 {
		cs = append(cs, m.ColProjectID().Eq(prj))
	}
	if len(r.AppletIDs) > 0 {
		cs = append(cs, m.ColAppletID().In(r.AppletIDs))
	}
	if len(r.Names) > 0 {
		cs = append(cs, m.ColName().In(r.Names))
	}
	if r.NameLike != "" {
		cs = append(cs, m.ColName().Like(r.NameLike))
	}
	if r.NameLike != "" {
		cs = append(cs, m.ColName().LLike(r.LNameLike))
	}
	return builder.And(cs...)
}

func (r *ListReq) Additions() builder.Additions {
	m := &models.Applet{}
	return builder.Additions{
		builder.OrderBy(
			builder.DescOrder(m.ColUpdatedAt()),
			builder.DescOrder(m.ColCreatedAt()),
		),
		r.Pager.Addition(),
	}
}

type ListRsp struct {
	Data  []models.Applet `json:"data"`
	Hints int64           `json:"hints"`
}

type Detail struct {
	models.Applet
	models.ResourceInfo
	*models.InstanceInfo
}

type ListDetailRsp struct {
	Data  []*Detail `json:"data"`
	Hints int64     `json:"hints"`
}

type Info struct {
	AppletName string                `json:"appletName"`
	WasmName   string                `json:"wasmName,omitempty"`
	WasmMd5    string                `json:"wasmMd5,omitempty"`
	WasmCache  *wasm.Cache           `json:"wasmCache,omitempty"`
	Strategies []models.StrategyInfo `json:"strategies,omitempty"`
}

type CreateReq struct {
	File *multipart.FileHeader `name:"file"`
	Info `name:"info"`
}

func (r *CreateReq) BuildStrategies(ctx context.Context) []models.Strategy {
	ids := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFIDs(len(r.Strategies) + 1)
	app := types.MustAppletFromContext(ctx)
	prj := types.MustProjectFromContext(ctx)
	sty := make([]models.Strategy, 0, len(r.Strategies))
	for i := range r.Strategies {
		sty = append(sty, models.Strategy{
			RelStrategy:  models.RelStrategy{StrategyID: ids[i]},
			RelProject:   models.RelProject{ProjectID: prj.ProjectID},
			RelApplet:    models.RelApplet{AppletID: app.AppletID},
			StrategyInfo: r.Strategies[i],
		})
	}
	if len(sty) == 0 {
		sty = append(sty, models.Strategy{
			RelStrategy:  models.RelStrategy{StrategyID: ids[0]},
			RelProject:   models.RelProject{ProjectID: prj.ProjectID},
			RelApplet:    models.RelApplet{AppletID: app.AppletID},
			StrategyInfo: models.DefaultStrategyInfo,
		})
	}
	return sty
}

type CreateRsp struct {
	models.RelApplet
	models.AppletInfo
	models.RelInstance
	models.InstanceInfo
}

type InfoForUpdate struct {
	AppletName string                `json:"appletName,omitempty"`
	WasmName   string                `json:"wasmName,omitempty"`
	WasmMd5    string                `json:"wasmMd5,omitempty"`
	WasmCache  *wasm.Cache           `json:"wasmCache,omitempty"`
	Strategies []models.StrategyInfo `json:"strategies,omitempty"`
}

type UpdateReq struct {
	File *multipart.FileHeader `name:"file,omitempty"`
	Info *InfoForUpdate        `name:"info"`
}

func (r *UpdateReq) Assignments() builder.Assignments {
	m := &models.Applet{}
	if r.Info == nil {
		return nil
	}
	ret := make(builder.Assignments, 2)
	if v := r.Info.AppletName; v != "" {
		ret = append(ret, m.ColName().ValueBy(v))
	}
	if v := r.Info.WasmName; v != "" {
		ret = append(ret, m.ColWasmName().ValueBy(v))
	}
	if len(ret) == 0 {
		return nil
	}
	return append(ret,
		m.ColUpdatedAt().ValueBy(types.Timestamp{Time: time.Now()}),
	)
}

func (r *UpdateReq) BuildStrategies(ctx context.Context) []models.Strategy {
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	app := types.MustAppletFromContext(ctx)
	prj := types.MustProjectFromContext(ctx)
	if r.Info == nil || len(r.Info.Strategies) == 0 {
		return []models.Strategy{{
			RelStrategy:  models.RelStrategy{StrategyID: idg.MustGenSFID()},
			RelProject:   models.RelProject{ProjectID: prj.ProjectID},
			RelApplet:    models.RelApplet{AppletID: app.AppletID},
			StrategyInfo: models.DefaultStrategyInfo,
		}}
	}

	sty := make([]models.Strategy, 0, len(r.Info.Strategies))
	ids := idg.MustGenSFIDs(len(r.Info.Strategies))
	for i := range r.Info.Strategies {
		sty = append(sty, models.Strategy{
			RelStrategy:  models.RelStrategy{StrategyID: ids[i]},
			RelProject:   models.RelProject{ProjectID: prj.ProjectID},
			RelApplet:    models.RelApplet{AppletID: app.AppletID},
			StrategyInfo: r.Info.Strategies[i],
		})
	}
	return sty
}

type UpdateRsp = CreateRsp
