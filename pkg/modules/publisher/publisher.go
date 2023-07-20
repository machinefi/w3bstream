package publisher

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/types"
)

func GetBySFID(ctx context.Context, id types.SFID) (*models.Publisher, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Publisher{RelPublisher: models.RelPublisher{PublisherID: id}}

	if err := m.FetchByPublisherID(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.PublisherNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func GetByProjectAndKey(ctx context.Context, prj types.SFID, key string) (*models.Publisher, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Publisher{
		RelProject:    models.RelProject{ProjectID: prj},
		PublisherInfo: models.PublisherInfo{Key: key},
	}

	if err := m.FetchByProjectIDAndKey(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.PublisherNotFound
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return m, nil
}

func ListByCond(ctx context.Context, r *CondArgs) (data []models.Publisher, err error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Publisher{}
	)
	data, err = m.List(d, r.Condition())
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return data, nil
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m = &models.Publisher{}

		ret = &ListRsp{}
		err error

		cond = r.Condition()
		adds = r.Additions()
	)

	ret.Data, err = m.List(d, cond, adds...)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = m.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func ListDetail(ctx context.Context, r *ListReq) (*ListDetailRsp, error) {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)

		pub = &models.Publisher{}
		prj = types.MustProjectFromContext(ctx)
		ret = &ListDetailRsp{}
		err error

		cond = r.Condition()
		adds = r.Additions()
	)

	expr := builder.Select(builder.MultiWith(",",
		builder.Alias(prj.ColName(), "f_project_name"),
		pub.ColProjectID(),
		pub.ColPublisherID(),
		pub.ColName(),
		pub.ColKey(),
		pub.ColCreatedAt(),
		pub.ColUpdatedAt(),
	)).From(
		d.T(pub),
		append([]builder.Addition{
			builder.LeftJoin(d.T(prj)).On(pub.ColProjectID().Eq(prj.ColProjectID())),
			builder.Where(builder.And(cond, prj.ColDeletedAt().Neq(0))),
		}, adds...)...,
	)
	err = d.QueryAndScan(expr, ret.Data)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = pub.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func RemoveBySFID(ctx context.Context, acc *models.Account, prj *models.Project, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Publisher{}

	if err := sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			var err error
			m, err = GetBySFID(ctx, id)
			return err
		},
		func(d sqlx.DBExecutor) error {
			if err := m.DeleteByPublisherID(d); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			ctx := contextx.WithContextCompose(
				types.WithMgrDBExecutorContext(d),
				types.WithAccountContext(acc),
			)(ctx)
			return access_key.DeleteByName(ctx, "pub_"+id.String())
		},
	).Do(); err != nil {
		return err
	}
	metrics.PublisherMetricsDec(ctx, acc.AccountID.String(), prj.Name)
	return nil
}

func RemoveByProjectAndKey(ctx context.Context, prj types.SFID, key string) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Publisher{}

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			var err error
			m, err = GetByProjectAndKey(ctx, prj, key)
			return err
		},
		func(d sqlx.DBExecutor) error {
			if err := m.DeleteByProjectIDAndKey(d); err != nil {
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
}

func Remove(ctx context.Context, acc *models.Account, r *CondArgs) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Publisher{}
	prj := types.MustProjectFromContext(ctx)

	expr := builder.Delete().From(d.T(m), builder.Where(r.Condition()))

	res, err := d.Exec(expr)
	if err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	numDeleted, err := res.RowsAffected()
	if err != nil {
		return err
	}
	for i := 0; i < int(numDeleted); i++ {
		metrics.PublisherMetricsDec(ctx, acc.AccountID.String(), prj.Name)
	}
	return nil
}

func Create(ctx context.Context, r *CreateReq) (*models.Publisher, error) {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		prj = types.MustProjectFromContext(ctx)
		acc = types.MustAccountFromContext(ctx)
		idg = confid.MustSFIDGeneratorFromContext(ctx)
		pub *models.Publisher
		tok *access_key.CreateRsp
	)

	err := sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) (err error) {
			id := idg.MustGenSFID()
			tok, err = access_key.Create(types.WithMgrDBExecutor(ctx, d), &access_key.CreateReq{
				IdentityID:   id,
				IdentityType: enums.ACCESS_KEY_IDENTITY_TYPE__PUBLISHER,
				CreateReqBase: access_key.CreateReqBase{
					Name: "pub_" + id.String(),
					Desc: "pub_" + id.String(),
					Privileges: access_key.GroupAccessPrivileges{{
						Name: enums.ApiGroupEvent,
						Perm: enums.ACCESS_PERMISSION__READ_WRITE,
					}},
				},
			})
			if err != nil {
				return err
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			pub = &models.Publisher{
				RelProject:   models.RelProject{ProjectID: prj.ProjectID},
				RelPublisher: models.RelPublisher{PublisherID: tok.IdentityID},
				PublisherInfo: models.PublisherInfo{
					Name:  r.Name,
					Key:   r.Key,
					Token: tok.AccessKey,
				},
			}
			if err := pub.Create(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.PublisherConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()

	if err != nil {
		return nil, err
	}
	metrics.PublisherMetricsInc(ctx, acc.AccountID.String(), prj.Name)
	return pub, nil
}

func Update(ctx context.Context, r *UpdateReq) error {
	var (
		d = types.MustMgrDBExecutorFromContext(ctx)
		m *models.Publisher
	)

	// TODO gen publisher token m.Token = "", or not ?

	return sqlx.NewTasks(d).With(
		func(d sqlx.DBExecutor) error {
			ctx := types.WithMgrDBExecutor(ctx, d)
			var err error
			m, err = GetBySFID(ctx, r.PublisherID)
			return err
		},
		func(d sqlx.DBExecutor) error {
			m.Key = r.Key
			m.Name = r.Name
			if err := m.UpdateByPublisherID(d); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.PublisherConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
	).Do()
}
