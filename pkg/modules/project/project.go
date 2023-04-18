// project management

package project

import (
	"context"
	"fmt"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/modules/config"
	"github.com/machinefi/w3bstream/pkg/modules/mq"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func GetProjectByProjectName(ctx context.Context, prjName string) (*models.Project, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Project{ProjectName: models.ProjectName{Name: prjName}}

	_, l = l.Start(ctx, "GetProjectByProjectName")
	defer l.End()

	if err := m.FetchByName(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetProjectByProjectName")
	}

	return m, nil
}

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
	rsp, err := applet.ListDetail(ctx, prj.ProjectID, &applet.ListReq{})
	if err != nil {
		return nil, err
	}

	return &Detail{
		ProjectID:   prj.ProjectID,
		ProjectName: prj.Name,
		Applets:     rsp.Data,
	}, nil
}

func List(ctx context.Context, acc types.SFID, r *ListReq) (*ListRsp, error) {
	var (
		err  error
		d    = types.MustMgrDBExecutorFromContext(ctx)
		prj  = types.MustProjectFromContext(ctx)
		ret  = &ListRsp{}
		cond = r.Condition(acc)
		adds = r.Additions()
	)

	ret.Data, err = prj.List(d, cond, adds...)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	ret.Total, err = prj.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	return ret, nil
}

func ListDetail(ctx context.Context, acc types.SFID, r *ListReq) (*ListDetailRsp, error) {
	ret := &ListDetailRsp{}

	lst, err := List(ctx, acc, r)
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

func Create(ctx context.Context, acc types.SFID, r *CreateReq) (*models.Project, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	prj := &models.Project{
		RelProject:  models.RelProject{ProjectID: idg.MustGenSFID()},
		RelAccount:  models.RelAccount{AccountID: acc},
		ProjectName: models.ProjectName{Name: acc.String() + "_" + r.Name},
		ProjectBase: models.ProjectBase{
			Version:     r.Version,
			Proto:       r.Proto,
			Description: r.Description,
			// TODO Issuer:      "",
			// TODO ExpIn:       0,
			// TODO SignKey:     "",
		},
	}

	if err := prj.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.ProjectNameConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return prj, nil
}

func CreateWithConfig(ctx context.Context, acc types.SFID, r *CreateReq) (*CreateRsp, error) {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		prj *models.Project
		err error
		cfg map[enums.ConfigType]wasm.Configuration
	)

	err = sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			prj, err = Create(ctx, acc, r)
			return err
		},
		func(d sqlx.DBExecutor) error {
			configs := r.Configs(prj.ProjectID.String())
			for _, c := range configs {
				if err = config.Create(ctx, prj.ProjectID, c); err != nil {
					return err
				}
				if cfg == nil {
					cfg = make(map[enums.ConfigType]wasm.Configuration)
				}
				cfg[c.ConfigType()] = c
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, err
	}

	for t, c := range cfg {
		if canBeInit, ok := c.(wasm.ConfigurationWithInit); ok {
			if err = canBeInit.Init(ctx); err != nil {
				return nil, status.ConfigInitFailed.StatusErr().WithDesc(
					fmt.Sprintf("project: %v type: %v err: %v", prj.ProjectID, t, err),
				)
			}
		}
	}
	return &CreateRsp{
		Project: prj,
		Configs: cfg,
	}, nil
}

func Remove(ctx context.Context) error {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		p = types.MustProjectFromContext(ctx)
	)

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			if err := p.DeleteByProjectID(d); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			return config.BatchRemove(ctx, &config.DataListParam{
				RelIDs: []types.SFID{p.ProjectID}},
			)
		},
		func(d sqlx.DBExecutor) error {
			return applet.BatchRemove(ctx, &applet.DataListParam{
				ProjectID: p.ProjectID,
			})
		},
	).Do()
}

func Init(ctx context.Context) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	_, l := conflog.FromContext(ctx).Start(ctx, "project.Init")
	defer l.End()

	data, err := (&models.Project{}).List(d, nil)
	if err != nil {
		l.Panic(err)
	}
	for i := range data {
		v := &data[i]
		l = l.WithValues("prj", v.Name)
		if err := mq.CreateChannel(ctx, v.Name, nil); err != nil {
			l.Panic(err)
		}
		l.Info("start subscribe")
	}
}
