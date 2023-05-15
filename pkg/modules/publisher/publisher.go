package publisher

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	confjwt "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

var _registerPublisherMtc = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "w3b_register_publisher_metrics",
	Help: "register publisher counter metrics.",
}, []string{"project"})

func init() {
	prometheus.MustRegister(_registerPublisherMtc)
}

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

func RemoveBySFID(ctx context.Context, id types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Publisher{}

	return sqlx.NewTasks(d).With(
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
	).Do()
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

func Remove(ctx context.Context, r *CondArgs) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Publisher{}

	expr := builder.Delete().From(d.T(m), builder.Where(r.Condition()))

	if _, err := d.Exec(expr); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func Create(ctx context.Context, r *CreateReq) (*models.Publisher, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	prj := types.MustProjectFromContext(ctx)

	id := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID()
	token, err := confjwt.MustConfFromContext(ctx).GenerateTokenWithoutExpByPayload(id)
	if err != nil {
		return nil, status.GenPublisherTokenFailed.StatusErr().WithDesc(err.Error())
	}

	pub := &models.Publisher{
		RelProject:   models.RelProject{ProjectID: prj.ProjectID},
		RelPublisher: models.RelPublisher{PublisherID: id},
		PublisherInfo: models.PublisherInfo{
			Name:  r.Name,
			Key:   r.Key,
			Token: token,
		},
	}

	// TODO move matrix to other place @zhiwei
	_registerPublisherMtc.WithLabelValues(prj.Name).Inc()

	if err = pub.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.PublisherConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
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
