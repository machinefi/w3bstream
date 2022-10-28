package publisher

import (
	"context"

	confid "github.com/iotexproject/Bumblebee/conf/id"
	"github.com/iotexproject/Bumblebee/conf/jwt"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"

	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreatePublisherReq struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func CreatePublisher(ctx context.Context, projectID types.SFID, r *CreatePublisherReq) (*models.Publisher, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	j := jwt.MustConfFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreatePublisher")
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
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Publisher{
		PublisherInfo: models.PublisherInfo{Key: publisherKey},
	}

	_, l = l.Start(ctx, "GetPublisherByPublisherKey")

	if err := m.FetchByKey(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetPublisherByPublisherKey")
	}

	return m, nil
}

type ListPublisherReq struct {
	projectID types.SFID
}

func (r *ListPublisherReq) SetCurrentProject(prjID types.SFID) { r.projectID = prjID }

func (r *ListPublisherReq) Condition() builder.SqlCondition {
	m := &models.Publisher{}
	return m.ColProjectID().Eq(r.projectID)
}

func (r *ListPublisherReq) Additions() builder.Additions { return nil }

type ListPublisherRsp struct {
	Total int64              `json:"total"`
	Data  []models.Publisher `json:"data"`
}

func ListPublisher(ctx context.Context, r *ListPublisherReq) (ret *ListPublisherRsp, err error) {
	l := types.MustLoggerFromContext(ctx)
	d := types.MustDBExecutorFromContext(ctx)

	m := &models.Publisher{}

	_, l = l.Start(ctx, "ListPublisher")
	defer l.End()

	ret = &ListPublisherRsp{}

	ret.Data, err = m.List(d, r.Condition(), r.Additions()...)
	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "ListPublisher")
	}

	ret.Total, err = m.Count(d, r.Condition())
	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "ListPublisherCount")
	}
	return ret, nil
}

type RemovePublisherReq struct {
	ProjectName  string       `in:"path" name:"projectName"`
	PublisherIDs []types.SFID `in:"query" name:"publisherID"`
}

func RemovePublisher(ctx context.Context, r *RemovePublisherReq) error {
	var (
		d          = types.MustDBExecutorFromContext(ctx)
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

func UpdatePublisher(ctx context.Context, publisherID types.SFID, r *CreatePublisherReq) (err error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	j := jwt.MustConfFromContext(ctx)
	m := models.Publisher{RelPublisher: models.RelPublisher{PublisherID: publisherID}}

	_, l = l.Start(ctx, "UpdatePublisher")
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
