package publisher

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/pkg/errors"
)

type CreatePublisherReq struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func CreatePublisher(ctx context.Context, projectID types.SFID, r *CreatePublisherReq) (*models.Publisher, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	j := jwt.MustConfFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx)
	defer l.End()

	// TODO generate token, maybe use public key
	publisherID := idg.MustGenSFID()
	token, err := j.GenerateTokenByPayload(publisherID)
	if err != nil {
		l.Error(err)
		return nil, status.InternalServerError.StatusErr().WithDesc(err.Error())
	}

	m := &models.Publisher{
		RelProject:    models.RelProject{ProjectID: projectID},
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
	l := types.MustLoggerFromContext(ctx)
	m := &models.Publisher{
		PublisherInfo: models.PublisherInfo{Key: publisherKey},
	}

	_, l = l.Start(ctx)

	if err := m.FetchByKey(d); err != nil {
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

	_, l = l.Start(ctx)
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

	_, l = l.Start(ctx)
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

func UpdatePublisher(ctx context.Context, publisherID types.SFID, r *CreatePublisherReq) (err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	j := jwt.MustConfFromContext(ctx)
	m := models.Publisher{RelPublisher: models.RelPublisher{PublisherID: publisherID}}

	_, l = l.Start(ctx)
	defer l.End()

	// TODO generate token, maybe use public key
	token, err := j.GenerateTokenByPayload(publisherID)
	if err != nil {
		l.Error(err)
		return status.InternalServerError.StatusErr().WithDesc(err.Error())
	}

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			return m.FetchByPublisherID(d)
		},
		func(db sqlx.DBExecutor) error {
			m.PublisherInfo.Name = r.Name
			m.PublisherInfo.Key = r.Key
			m.PublisherInfo.Token = token
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
	d := types.MustMgrDBExecutorFromContext(ctx)

	_, l = l.Start(ctx)
	defer l.End()

	pub := &models.Publisher{PublisherInfo: models.PublisherInfo{Key: pubKey}}
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
