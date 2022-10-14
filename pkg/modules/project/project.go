// project management

package project

import (
	"context"

	confid "github.com/iotexproject/Bumblebee/conf/id"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"

	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/errors/status"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/iotexproject/w3bstream/pkg/models"
	"github.com/iotexproject/w3bstream/pkg/types"
)

type CreateProjectReq = models.ProjectInfo

func CreateProject(ctx context.Context, r *CreateProjectReq) (*models.Project, error) {
	d := types.MustDBExecutorFromContext(ctx)
	a := middleware.CurrentAccountFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	m := &models.Project{
		RelProject:  models.RelProject{ProjectID: idg.MustGenSFID()},
		RelAccount:  models.RelAccount{AccountID: a.AccountID},
		ProjectInfo: *r,
	}

	if err := m.Create(d); err != nil {
		return nil, err
	}

	return m, nil
}

type ListProjectReq struct {
	accountID  types.SFID
	IDs        []uint64     `in:"query" name:"ids,omitempty"`
	ProjectIDs []types.SFID `in:"query" name:"projectIDs,omitempty"`
	Names      []string     `in:"query" name:"names,omitempty"`
	NameLike   string       `in:"query" name:"name,omitempty"`
	datatypes.Pager
}

func (r *ListProjectReq) SetCurrentAccount(accountID types.SFID) {
	r.accountID = accountID
}

func (r *ListProjectReq) Condition() builder.SqlCondition {
	var (
		m  = &models.Project{}
		cs []builder.SqlCondition
	)

	cs = append(cs, m.ColAccountID().Eq(r.accountID))
	if len(r.IDs) > 0 {
		cs = append(cs, m.ColID().In(r.IDs))
	}
	if len(r.ProjectIDs) > 0 {
		cs = append(cs, m.ColProjectID().In(r.ProjectIDs))
	}
	if len(r.Names) > 0 {
		cs = append(cs, m.ColName().In(r.Names))
	}
	if r.NameLike != "" {
		cs = append(cs, m.ColName().Like(r.NameLike))
	}

	return builder.And(cs...)
}

func (r *ListProjectReq) Additions() builder.Additions {
	m := &models.Project{}
	return builder.Additions{
		builder.OrderBy(builder.DescOrder(m.ColCreatedAt())),
		r.Pager.Addition(),
	}
}

type ListProjectRsp struct {
	Data  []Detail `json:"data"`  // Data project data list
	Total int64    `json:"total"` // Total project count under current user
}

type Detail struct {
	ProjectID   types.SFID     `json:"projectID"`
	ProjectName string         `json:"projectName"`
	Applets     []AppletDetail `json:"applets,omitempty"`
	datatypes.OperationTimes
}

type AppletDetail struct {
	AppletID      types.SFID          `json:"appletID"`
	AppletName    string              `json:"appletName"`
	InstanceID    types.SFID          `json:"instanceID,omitempty"`
	InstanceState enums.InstanceState `json:"instanceState,omitempty"`
}

type detail struct {
	ProjectID     types.SFID          `db:"f_project_id"`
	ProjectName   string              `db:"f_project_name"`
	AppletID      types.SFID          `db:"f_applet_id"`
	AppletName    string              `db:"f_applet_name"`
	InstanceID    types.SFID          `db:"f_instance_id"`
	InstanceState enums.InstanceState `db:"f_instance_state"`
	datatypes.OperationTimes
}

func ListProject(ctx context.Context, r *ListProjectReq) (*ListProjectRsp, error) {
	var (
		d    = types.MustDBExecutorFromContext(ctx)
		ret  = &ListProjectRsp{}
		err  error
		cond = r.Condition()

		mProject  = &models.Project{}
		mApplet   = &models.Applet{}
		mInstance = &models.Instance{}
	)
	ret.Total, err = mProject.Count(d, cond)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "CountProject")
	}

	details := make([]detail, 0)

	err = d.QueryAndScan(
		builder.Select(
			builder.MultiWith(
				",",
				builder.Alias(mProject.ColProjectID(), "f_project_id"),
				builder.Alias(mProject.ColName(), "f_project_name"),
				builder.Alias(mApplet.ColAppletID(), "f_applet_id"),
				builder.Alias(mApplet.ColName(), "f_applet_name"),
				builder.Alias(mInstance.ColInstanceID(), "f_instance_id"),
				builder.Alias(mInstance.ColState(), "f_instance_state"),
				builder.Alias(mProject.ColCreatedAt(), "f_created_at"),
				builder.Alias(mProject.ColUpdatedAt(), "f_updated_at"),
			),
		).From(
			d.T(mProject),
			builder.LeftJoin(d.T(mApplet)).
				On(mProject.ColProjectID().Eq(mApplet.ColProjectID())),
			builder.LeftJoin(d.T(mInstance)).
				On(mApplet.ColAppletID().Eq(mInstance.ColAppletID())),
			builder.Where(cond),
			builder.OrderBy(
				builder.DescOrder(mProject.ColCreatedAt()),
				builder.AscOrder(mApplet.ColName()),
			),
			r.Pager.Addition(),
		),
		&details,
	)
	if err != nil {
		return nil, status.CheckDatabaseError(err, "ListProject")
	}

	detailsMap := make(map[types.SFID][]*detail)
	for i := range details {
		prjID := details[i].ProjectID
		detailsMap[prjID] = append(detailsMap[prjID], &details[i])
	}

	for prjID, vmap := range detailsMap {
		appletDetails := make([]AppletDetail, 0, len(vmap))
		for _, v := range vmap {
			if v.AppletID == 0 {
				continue
			}
			appletDetails = append(appletDetails, AppletDetail{
				AppletID:      v.AppletID,
				AppletName:    v.AppletName,
				InstanceID:    v.InstanceID,
				InstanceState: v.InstanceState,
			})
		}
		if len(appletDetails) == 0 {
			appletDetails = nil
		}
		ret.Data = append(ret.Data, Detail{
			ProjectID:   prjID,
			ProjectName: vmap[0].ProjectName,
			Applets:     appletDetails,
			OperationTimes: datatypes.OperationTimes{
				CreatedAt: vmap[0].CreatedAt,
				UpdatedAt: vmap[0].UpdatedAt,
			},
		})
	}

	return ret, nil
}

func GetProjectByProjectID(ctx context.Context, prjID types.SFID) (*Detail, error) {
	d := types.MustDBExecutorFromContext(ctx)
	ca := middleware.CurrentAccountFromContext(ctx)

	_, err := ca.ValidateProjectPerm(ctx, prjID)
	if err != nil {
		return nil, err
	}
	m := &models.Project{RelProject: models.RelProject{ProjectID: prjID}}

	if err := m.FetchByProjectID(d); err != nil {
		return nil, status.CheckDatabaseError(err, "GetProjectByProjectID")
	}

	ret, err := ListProject(ctx, &ListProjectReq{
		accountID:  ca.AccountID,
		ProjectIDs: []types.SFID{prjID},
	})

	if err != nil {
		return nil, err
	}

	if len(ret.Data) == 0 {
		return nil, status.NotFound
	}

	return &ret.Data[0], nil
}

func DeleteProject(_ context.Context, _ string) error {
	// TODO
	return nil
}
