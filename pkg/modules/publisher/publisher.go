package publisher

import (
	"context"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
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

type CreatePublisherReq struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func CreatePublisher(ctx context.Context, project *models.Project, r *CreatePublisherReq) (*models.Publisher, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)
	publisherJwt := jwt.MustConfFromContext(ctx)

	_, l = l.Start(ctx, "CreatePublisher")
	defer l.End()

	publisherID := idg.MustGenSFID()
	token, err := publisherJwt.GenerateTokenWithoutExpByPayload(publisherID)
	if err != nil {
		l.Error(err)
		return nil, status.InternalServerError.StatusErr().WithDesc(err.Error())
	}

	_registerPublisherMtc.WithLabelValues(project.Name).Inc()

	m := &models.Publisher{
		RelProject:    models.RelProject{ProjectID: project.ProjectID},
		RelPublisher:  models.RelPublisher{PublisherID: publisherID},
		PublisherInfo: models.PublisherInfo{Name: r.Name, Key: r.Key, Token: token},
	}
	if err = m.Create(d); err != nil {
		l.Error(err)
		return nil, err
	}

	return m, nil
}

func GetPublisherByPublisherKey(ctx context.Context, publisherKey string) (*models.Publisher, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	p := types.MustProjectFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Publisher{
		RelProject:    models.RelProject{ProjectID: p.ProjectID},
		PublisherInfo: models.PublisherInfo{Key: publisherKey},
	}

	_, l = l.Start(ctx, "GetPublisherByPublisherKey")

	if err := m.FetchByProjectIDAndKey(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetPublisherByPublisherKey")
	}

	return m, nil
}

type ListPublisherReq struct {
	projectID types.SFID
	datatypes.Pager
}

func (r *ListPublisherReq) SetCurrentProject(prjID types.SFID) { r.projectID = prjID }

func (r *ListPublisherReq) Condition() builder.SqlCondition {
	m := &models.Publisher{}
	return m.ColProjectID().Eq(r.projectID)
}

func (r *ListPublisherReq) Additions() builder.Additions { return nil }

type InfoPublisher struct {
	models.Publisher
	ProjectName string `db:"f_project_name"`
	datatypes.OperationTimes
}

type ListPublisherRsp struct {
	Total int64           `json:"total"`
	Data  []InfoPublisher `json:"data"`
}

func ListPublisher(ctx context.Context, r *ListPublisherReq) (*ListPublisherRsp, error) {
	var (
		l = types.MustLoggerFromContext(ctx)
		d = types.MustMgrDBExecutorFromContext(ctx)

		ret        = &ListPublisherRsp{}
		err        error
		cond       = r.Condition()
		mPublisher = &models.Publisher{}
		mProject   = &models.Project{}
	)

	_, l = l.Start(ctx, "ListPublisher")
	defer l.End()

	ret.Total, err = mPublisher.Count(d, cond)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "CountPublisher")
	}

	details := make([]InfoPublisher, 0)
	err = d.QueryAndScan(
		builder.Select(
			builder.MultiWith(
				",",
				builder.Alias(mPublisher.ColProjectID(), "f_project_id"),
				builder.Alias(mProject.ColName(), "f_project_name"),
				builder.Alias(mPublisher.ColPublisherID(), "f_publisher_id"),
				builder.Alias(mPublisher.ColName(), "f_name"),
				builder.Alias(mPublisher.ColKey(), "f_key"),
				builder.Alias(mPublisher.ColToken(), "f_token"),
				builder.Alias(mPublisher.ColCreatedAt(), "f_created_at"),
				builder.Alias(mPublisher.ColUpdatedAt(), "f_updated_at"),
			),
		).From(
			d.T(mPublisher),
			builder.LeftJoin(d.T(mProject)).
				On(mPublisher.ColProjectID().Eq(mProject.ColProjectID())),
			builder.Where(cond),
			builder.OrderBy(
				builder.DescOrder(mPublisher.ColCreatedAt()),
				builder.AscOrder(mPublisher.ColName()),
			),
			r.Pager.Addition(),
		),
		&details,
	)
	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "ListPublisher")
	}

	ret.Data = details
	return ret, err
}

type RemovePublisherReq struct {
	ProjectName  string       `in:"path" name:"projectName"`
	PublisherIDs []types.SFID `in:"query" name:"publisherID"`
}

func RemovePublisher(ctx context.Context, r *RemovePublisherReq) error {
	var (
		d          = types.MustMgrDBExecutorFromContext(ctx)
		l          = types.MustLoggerFromContext(ctx)
		mPublisher = &models.Publisher{}
		err        error
	)

	_, l = l.Start(ctx, "RemovePublisher")
	defer l.End()

	return sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			for _, id := range r.PublisherIDs {
				mPublisher.PublisherID = id
				if err = mPublisher.DeleteByPublisherID(d); err != nil {
					l.Error(err)
					return status.CheckDatabaseError(err, "DeleteByPublisherID")
				}
			}
			return nil
		},
	).Do()
}

func UpdatePublisher(ctx context.Context, project *models.Project, publisherID types.SFID, r *CreatePublisherReq) (err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := models.Publisher{RelPublisher: models.RelPublisher{PublisherID: publisherID}}

	_, l = l.Start(ctx, "UpdatePublisher")
	defer l.End()

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			return m.FetchByPublisherID(d)
		},
		func(db sqlx.DBExecutor) error {
			m.PublisherInfo.Name = r.Name
			m.PublisherInfo.Key = r.Key
			return m.UpdateByPublisherID(d)
		},
	).Do()

	if err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "UpdatePublisher")
	}

	return
}

func GetPublisherByPubKeyAndProjectName(ctx context.Context, pubKey, prjName string) (*models.Publisher, error) {
	l := types.MustLoggerFromContext(ctx)
	prj, err := project.GetProjectByProjectName(ctx, prjName)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	l = l.WithValues("project_id", prj.ProjectID)

	d := types.MustMgrDBExecutorFromContext(ctx)
	_, l = l.Start(ctx, "GetPublisherByPubKeyAndProjectID")
	defer l.End()

	pub := &models.Publisher{
		RelProject:    models.RelProject{ProjectID: prj.ProjectID},
		PublisherInfo: models.PublisherInfo{Key: pubKey},
	}
	// TODO change prjName to projectID, then use FetchByProjectIDAndKey
	if err := pub.FetchByProjectIDAndKey(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetPublisherByProjectIDAndKey")
	}

	l = l.WithValues("pub_id", pub.PublisherID)

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
