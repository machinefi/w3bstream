package publisher

import (
	"context"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

var _registerPublisherMtc = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "w3b_register_publisher_metrics",
	Help: "register publisher counter metrics.",
}, []string{"project"})

func init() {
	prometheus.MustRegister(_registerPublisherMtc)
}

func GetPublisherByPubKeyAndProjectName(ctx context.Context, pubKey, prjName string) (*models.Publisher, error) {
	l := types.MustLoggerFromContext(ctx)
	d := types.MustMgrDBExecutorFromContext(ctx)

	_, l = l.Start(ctx, "GetPublisherByPubKeyAndProjectID")
	defer l.End()

	pub := &models.Publisher{PublisherInfo: models.PublisherInfo{Key: pubKey}}
	// TODO change prjName to projectID, then use FetchByProjectIDAndKey
	if err := pub.FetchByKey(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetPublisherByKey")
	}

	l = l.WithValues("pub_id", pub.PublisherID)
	prj, err := project.GetProjectByProjectName(ctx, prjName)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	l = l.WithValues("project_id", prj.ProjectID)

	if pub.ProjectID != prj.ProjectID {
		l.Error(errors.New("no project permission"))
		return nil, status.NoProjectPermission
	}
	return pub, nil
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

func Remove(ctx context.Context) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := types.MustPublisherFromContext(ctx)

	if err := m.DeleteByPublisherID(d); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func BatchRemoveBySFIDs(ctx context.Context, ids []types.SFID) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := &models.Publisher{}

	expr := builder.Delete().From(d.T(m), builder.Where(m.ColPublisherID().In(ids)))

	if _, err := d.Exec(expr); err != nil {
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}

func Create(ctx context.Context, r *CreateReq) (*models.Publisher, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	prj := types.MustProjectFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	pub := &models.Publisher{
		RelProject:   models.RelProject{ProjectID: prj.ProjectID},
		RelPublisher: models.RelPublisher{PublisherID: idg.MustGenSFID()},
		PublisherInfo: models.PublisherInfo{
			Name: r.Name,
			Key:  r.Key,
			// TODO gen publsiher token: Token: "",
		},
	}

	// TODO matrix
	_registerPublisherMtc.WithLabelValues(prj.Name).Inc()

	if err := pub.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.PublisherKeyConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return pub, nil
}

func List(ctx context.Context, r *ListReq) (*ListRsp, error) {
	var (
		d    = types.MustMgrDBExecutorFromContext(ctx)
		prj  = types.MustProjectFromContext(ctx)
		pub  = &models.Publisher{}
		ret  = &ListRsp{}
		err  error
		cond = r.Condition(prj.ProjectID)
		adds = r.Additions()
	)

	ret.Data, err = pub.List(d, cond, adds...)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	ret.Total, err = pub.Count(d, cond)
	if err != nil {
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}

func ListDetail(ctx context.Context, r *ListReq) (*ListDetailRsp, error) {
	var (
		d   = types.MustMgrDBExecutorFromContext(ctx)
		prj = types.MustProjectFromContext(ctx)
		pub = &models.Publisher{}
		ret = &ListDetailRsp{}
		err error

		cond = r.Condition(prj.ProjectID)
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

func Update(ctx context.Context, r *UpdateReq) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	m := types.MustPublisherFromContext(ctx)

	m.Key = r.Key
	m.Name = r.Name
	// TODO gen publisher token m.Token = "", or not ?

	if err := m.UpdateByPublisherID(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return status.PublisherKeyConflict
		}
		return status.DatabaseError.StatusErr().WithDesc(err.Error())
	}
	return nil
}
