package applet

import (
	"context"
	"mime/multipart"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type CondArgs struct {
	ProjectID types.SFID   `name:"-"`
	AppletIDs []types.SFID `in:"query" name:"appletID,omitempty"`
	Names     []string     `in:"query" name:"names,omitempty"`
	NameLike  string       `in:"query" name:"name,omitempty"`
	LNameLike string       `in:"query" name:"lName,omitempty"`
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

type ListReq struct {
	CondArgs
	datatypes.Pager
}

type ListRsp struct {
	Data  []models.Applet `json:"data"`
	Total int64           `json:"total"`
}

type Detail struct {
	*models.Applet
	models.ResourceInfo
	*models.InstanceInfo
}

type ListDetailRsp struct {
	Data  []*Detail `json:"data"`
	Total int64     `json:"total"`
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

// BuildStrategies, must be built. if nil default strategy returned
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
	*models.Applet
	*models.Instance `json:"instance"`
	*models.Resource `json:"resource,omitempty"`
	Strategies       []models.Strategy `json:"strategies,omitempty"`
}

type UpdateReq struct {
	File *multipart.FileHeader `name:"file,omitempty"`
	Info `name:"info"`
}

// BuildStrategies try build, if invalid nil return
func (r *UpdateReq) BuildStrategies(ctx context.Context) []models.Strategy {
	if len(r.Strategies) == 0 {
		return nil
	}
	app := types.MustAppletFromContext(ctx)
	prj := types.MustProjectFromContext(ctx)
	sty := make([]models.Strategy, 0, len(r.Info.Strategies))
	ids := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFIDs(len(r.Info.Strategies))
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

func detail(app *models.Applet, ins *models.Instance, res *models.Resource) *Detail {
	ret := &Detail{Applet: app}

	if res != nil {
		ret.ResourceInfo = res.ResourceInfo
	}

	if ins != nil {
		ins.State, _ = vm.GetInstanceState(ins.InstanceID)
		ret.InstanceInfo = &ins.InstanceInfo
	}

	return ret
}
