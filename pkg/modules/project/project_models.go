package project

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type Detail struct {
	ProjectID   types.SFID       `json:"projectID"`
	ProjectName string           `json:"projectName"`
	Applets     []*applet.Detail `json:"applets,omitempty"`
}

type ListReq struct {
	ProjectIDs   []types.SFID `in:"query" name:"sfid,omitempty"`
	ProjectNames []string     `in:"query" name:"name,omitempty"`
	NameLike     string       `in:"query" name:"nameLike,omitempty"`
	LNameLike    string       `in:"query" name:"lNameLike,omitempty"`
	datatypes.Pager
}

func (r *ListReq) Condition(acc types.SFID) builder.SqlCondition {
	var (
		m  = &models.Project{}
		cs []builder.SqlCondition
	)

	if acc != 0 {
		cs = append(cs, m.ColAccountID().Eq(acc))
	}
	if len(r.ProjectIDs) > 0 {
		cs = append(cs, m.ColProjectID().In(r.ProjectIDs))
	}
	if len(r.ProjectNames) > 0 {
		for i := range r.ProjectNames {
			r.ProjectNames[i] = acc.String() + "_" + r.ProjectNames[i]
		}
		cs = append(cs, m.ColName().In(r.ProjectNames))
	}
	if r.NameLike != "" {
		cs = append(cs, m.ColName().Like(r.NameLike))
	}
	if r.LNameLike != "" {
		cs = append(cs, m.ColName().LLike(r.LNameLike))
	}
	return builder.And(cs...)
}

func (r *ListReq) Additions() builder.Additions {
	m := &models.Project{}
	return builder.Additions{
		builder.OrderBy(
			builder.DescOrder(m.ColUpdatedAt()),
			builder.DescOrder(m.ColCreatedAt()),
		),
		r.Pager.Addition(),
	}
}

type ListRsp struct {
	Data  []models.Project `json:"data"`
	Total int64            `json:"total"`
}

type ListDetailRsp struct {
	Data  []*Detail `json:"data"`
	Total int64     `json:"total"`
}

type CreateReq struct {
	models.ProjectName
	models.ProjectBase
	Envs   [][2]string  `json:"envs,omitempty"`
	Schema *wasm.Schema `json:"schema,omitempty"`
	// TODO if each project has its own mqtt broker should add *wasm.MqttClient
}

// func(d sqlx.DBExecutor) error {
// 	ctx = types.WithProject(ctx, m)
// 	if err := CreateOrUpdateProjectEnv(ctx, &wasm.Env{Env: r.Envs}); err != nil {
// 		return err
// 	}
// 	return nil
// },
// func(d sqlx.DBExecutor) error {
// 	if r.Schema == nil {
// 		sch := schema.NewSchema(r.Name)
// 		r.Schema = &wasm.Schema{Schema: *sch}
// 	}
// 	if err := CreateProjectSchema(ctx, r.Schema); err != nil {
// 		return err
// 	}
// 	return nil
// },

func (r *CreateReq) Configs(prefix string) []wasm.Configuration {
	sch := r.Schema
	if sch == nil {

	}
	return []wasm.Configuration{
		wasm.NewEnv(prefix),
		// TODO
	}
}

type CreateRsp struct {
	*models.Project `json:"project"`
	Configs         map[enums.ConfigType]wasm.Configuration `json:"configs,omitempty"`
}

type FullContext struct {
	*models.Project
	*models.Applet
	*models.Instance
	*models.Resource
}
