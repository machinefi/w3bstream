package tag

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateTagReq struct {
	models.TagInfo
}

func CreateTag(ctx context.Context, project *models.Project, r *CreateTagReq) (*models.Tag, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "CreateTag")
	defer l.End()

	m := &models.Tag{
		RelProject: models.RelProject{ProjectID: project.ProjectID},
		RelTag:     models.RelTag{TagID: idg.MustGenSFID()},
		TagInfo:    models.TagInfo{ReferenceID: r.ReferenceID, ReferenceType: r.ReferenceType, Info: r.Info},
	}

	if err := m.Create(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "CreateTag")
	}
	return m, nil
}

type RemoveTagReq struct {
	TagIDs []types.SFID `in:"query" name:"tagID,omitempty"`
}

func RemoveTag(ctx context.Context, project *models.Project, r *RemoveTagReq) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Tag{}

	_, l = l.Start(ctx, "RemoveTag")
	defer l.End()

	return sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			for _, id := range r.TagIDs {
				if _, err := db.Exec(
					builder.Delete().From(
						db.T(m),
						builder.Where(
							builder.And(
								m.ColProjectID().Eq(project.ProjectID),
								m.ColTagID().Eq(id),
							),
						),
					),
				); err != nil {
					l.Error(err)
					return status.CheckDatabaseError(err, "DeleteByTagID")
				}
			}
			return nil
		},
	).Do()
}

type ListTagReq struct {
	TagIDs         []uint64                 `in:"query" name:"tagID,omitempty"`
	ReferenceIDs   []types.SFID             `in:"query" name:"referenceID,omitempty"`
	ReferenceTypes []enums.TagReferenceType `in:"query" name:"referenceType,omitempty"`
	datatypes.Pager
}

func (r *ListTagReq) Condition(projectID types.SFID) builder.SqlCondition {
	var (
		m  = &models.Tag{}
		cs []builder.SqlCondition
	)
	if len(r.TagIDs) > 0 {
		cs = append(cs, m.ColTagID().In(r.TagIDs))
	}
	if len(r.ReferenceIDs) > 0 {
		cs = append(cs, m.ColReferenceID().In(r.ReferenceIDs))
	}
	if len(r.ReferenceTypes) > 0 {
		cs = append(cs, m.ColReferenceType().In(r.ReferenceTypes))
	}
	cs = append(cs, m.ColProjectID().Eq(projectID))

	return builder.And(cs...)
}

func (r *ListTagReq) Additions() builder.Additions {
	m := &models.Tag{}
	return builder.Additions{
		builder.OrderBy(builder.DescOrder(m.ColCreatedAt())),
		r.Pager.Addition(),
	}
}

type ListTagRsp struct {
	Data  []models.Tag `json:"data"`
	Hints int64        `json:"hints"`
}

func ListTags(ctx context.Context, project *models.Project, r *ListTagReq) (*ListTagRsp, error) {
	tag := &models.Tag{}

	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "ListTags")
	defer l.End()

	tags, err := tag.List(d, r.Condition(project.ProjectID), r.Additions()...)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	hints, err := tag.Count(d, r.Condition(project.ProjectID))
	if err != nil {
		l.Error(err)
		return nil, err
	}
	return &ListTagRsp{tags, hints}, nil
}
