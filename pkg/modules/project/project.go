// project management

package project

import (
	"context"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/transporter/mqtt"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func GetBySFID(ctx context.Context, prj types.SFID) (*models.Project, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)

	m := &models.Project{
		RelProject: models.RelProject{ProjectID: prj},
	}
	if err := m.FetchByProjectID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ProjectNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func GetByName(ctx context.Context, name string) (*models.Project, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Project{
		ProjectName: models.ProjectName{Name: name},
	}
	if err := m.FetchByName(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ProjectNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func GetDetail(ctx context.Context, prj *models.Project) (*Detail, error) {
	rsp, err := applet.ListDetail(ctx, &applet.ListReq{
		CondArgs: applet.CondArgs{ProjectID: prj.ProjectID},
	})
	if err != nil {
		return nil, err
	}

	return &Detail{
		ProjectID:   prj.ProjectID,
		ProjectName: prj.Name,
		Applets:     rsp.Data,
	}, nil
}

func ListByCond(ctx context.Context, r *CondArgs) ([]models.Project, error) {
	var (
		d    = types.MustMgrDBExecutorFromContext(ctx)
		prj  = &models.Project{}
		cond = r.Condition()
	)

	data, err := prj.List(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		err  error
		d    = types.MustMgrDBExecutorFromContext(ctx)
		prj  = &models.Project{}
		ret  = &ListRsp{}
		cond = r.Condition()
	)

	ret.Data, err = prj.List(d, cond, r.Pager.Addition())
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = prj.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	return ret, nil
}

func ListDetail(ctx context.Context, r *ListReq) (*ListDetailRsp, error) {
	ret := &ListDetailRsp{}

	lst, err := List(ctx, r)
	if err != nil {
		return nil, err
	}
	ret.Total = lst.Total

	for i := range lst.Data {
		detail, err := GetDetail(ctx, &lst.Data[i])
		if err != nil {
			return nil, err
		}
		ret.Data = append(ret.Data, detail)
	}
	return ret, nil
}

func Create(ctx context.Context, r *CreateReq) (*CreateRsp, error) {
	ctx, l := logr.Start(ctx, "modules.project.Create")
	defer l.End()

	d := types.MustMgrDBExecutorFromContext(ctx)
	acc := types.MustAccountFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	prj := &models.Project{
		RelProject:  models.RelProject{ProjectID: idg.MustGenSFID()},
		RelAccount:  models.RelAccount{AccountID: acc.AccountID},
		ProjectName: models.ProjectName{Name: r.Name},
		ProjectBase: models.ProjectBase{
			Version:     r.Version,
			Proto:       r.Proto,
			Description: r.Description,
		},
	}

	rsp := &CreateRsp{
		Project: prj,
	}

	err := sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			if err := prj.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.ProjectNameConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			ctx = types.WithProject(ctx, prj)
			return nil
		},
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			if r.Env == nil {
				r.Env = &wasm.Env{}
			}
			_, err := config.Create(ctx, prj.ProjectID, r.Env)
			if err != nil {
				return err
			}
			rsp.Env = r.Env
			if r.Database == nil {
				r.Database = &wasm.Database{}
			}
			_, err = config.Create(ctx, prj.ProjectID, r.Database)
			if err != nil {
				return err
			}
			rsp.Database = r.Database
			if r.Flow == nil {
				r.Flow = &wasm.Flow{}
			}
			_, err = config.Create(ctx, prj.ProjectID, r.Flow)
			if err != nil {
				return err
			}
			rsp.Flow = r.Flow
			return nil
		},
	).Do()
	if err != nil {
		l.Error(err)
		return nil, err
	}

	if err = mqtt.Subscribe(ctx, prj.Name); err != nil {
		l.WithValues("prj", prj.Name).Warn(errors.Wrap(err, "channel create failed"))
	}
	rsp.ChannelState = datatypes.BooleanValue(err == nil)

	return rsp, nil
}

func RemoveBySFID(ctx context.Context, id types.SFID) (err error) {
	ctx, l := logr.Start(ctx, "modules.project.RemoveBySFID", "porject_id", id)
	defer l.End()

	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		p *models.Project
	)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			if p, err = GetBySFID(ctx, id); err != nil {
				return err
			}
			mqtt.Stop(ctx, p.Name)
			l = l.WithValues("project_name", p.Name)
			l.Info("stop subscribing")
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if err := p.DeleteByProjectID(d); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			return config.Remove(ctx, &config.CondArgs{RelIDs: []types.SFID{p.ProjectID}})
		},
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			return applet.Remove(ctx, &applet.CondArgs{ProjectID: p.ProjectID})
		},
	).Do()
}

func Init(ctx context.Context) error {
	ctx, l := logr.Start(ctx, "modules.project.Init")
	defer l.End()

	d := types.MustMgrDBExecutorFromContext(ctx)

	data, err := (&models.Project{}).List(d, nil)
	if err != nil {
		return err
	}
	for i := range data {
		v := &data[i]
		l = l.WithValues("prj", v.Name)
		ctx = types.WithProject(ctx, v)
		if err = mqtt.Subscribe(ctx, v.Name); err != nil {
			l.Warn(errors.Wrap(err, "channel create failed"))
		}
		l.Info("start subscribe")
	}
	return nil
}
