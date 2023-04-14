package middleware

import (
	"context"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/account"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ContextAccountAuth struct {
	httpx.MethodGet
}

var contextAccountAuthKey = reflect.TypeOf(ContextAccountAuth{}).String()

func (r *ContextAccountAuth) ContextKey() string { return contextAccountAuthKey }

func (r *ContextAccountAuth) Output(ctx context.Context) (interface{}, error) {
	v, ok := jwt.AuthFromContext(ctx).(string)
	if !ok {
		return nil, status.InvalidAuthValue
	}
	accountID := types.SFID(0)
	if err := accountID.UnmarshalText([]byte(v)); err != nil {
		return nil, status.InvalidAuthAccountID
	}
	ca, err := account.GetAccountByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return &CurrentAccount{*ca}, nil
}

func CurrentAccountFromContext(ctx context.Context) *CurrentAccount {
	return ctx.Value(contextAccountAuthKey).(*CurrentAccount)
}

type CurrentAccount struct {
	models.Account
}

func (v *CurrentAccount) WithProjectContextByName(ctx context.Context, name string) (context.Context, error) {
	a := CurrentAccountFromContext(ctx)
	prj, err := project.GetByAccountAndName(ctx, a.AccountID, name)
	if err != nil {
		return nil, err
	}
	return types.WithProject(ctx, prj), nil
}

func (v *CurrentAccount) WithProjectContextByID(ctx context.Context, id types.SFID) (context.Context, error) {
	a := CurrentAccountFromContext(ctx)
	prj, err := project.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	if a.AccountID != prj.AccountID {
		return nil, status.NoProjectPermission
	}
	return types.WithProject(ctx, prj), nil
}

func (v *CurrentAccount) WithAppletContext(ctx context.Context, appletID types.SFID) (context.Context, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)

	app := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	if err := app.FetchByAppletID(d); err != nil {
		return ctx, status.CheckDatabaseError(err, "GetAppletByAppletID")
	}

	_ctx, err := v.WithProjectContextByID(ctx, app.ProjectID)
	if err != nil {
		return ctx, err
	}

	_ctx, err = v.WithResourceContext(ctx, app.ResourceID)
	if err != nil {
		return ctx, err
	}

	ctx = types.WithApplet(_ctx, app)
	return ctx, nil
}

func (v *CurrentAccount) WithResourceContext(ctx context.Context, resID types.SFID) (context.Context, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)

	res := &models.Resource{RelResource: models.RelResource{ResourceID: resID}}
	if err := res.FetchByResourceID(d); err != nil {
		return ctx, status.CheckDatabaseError(err, "GetResourceByResourceID")
	}

	ctx = types.WithResource(ctx, res)
	return ctx, nil
}

func (v *CurrentAccount) WithInstanceContext(ctx context.Context, instanceID types.SFID) (context.Context, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)

	ins := &models.Instance{RelInstance: models.RelInstance{InstanceID: instanceID}}
	if err := ins.FetchByInstanceID(d); err != nil {
		return ctx, status.CheckDatabaseError(err, "GetInstanceByInstanceID")
	}
	_ctx, err := v.WithAppletContext(ctx, ins.AppletID)
	if err != nil {
		return ctx, err
	}
	ctx = types.WithInstance(_ctx, ins)
	return ctx, nil
}

// ValidateProjectPerm
// Deprecated: Use WithProjectContextByID instead
func (v *CurrentAccount) ValidateProjectPerm(ctx context.Context, prjID types.SFID) (*models.Project, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	a := CurrentAccountFromContext(ctx)
	m := &models.Project{RelProject: models.RelProject{ProjectID: prjID}}

	if err := m.FetchByProjectID(d); err != nil {
		return nil, status.CheckDatabaseError(err, "GetProjectByProjectID")
	}
	if a.AccountID != m.AccountID {
		return nil, status.NoProjectPermission
	}
	return m, nil
}

// ValidateProjectPermByPrjName
// Deprecated: Use WithProjectContextByName instead
func (v *CurrentAccount) ValidateProjectPermByPrjName(ctx context.Context, name string) (*models.Project, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	a := CurrentAccountFromContext(ctx)
	m := &models.Project{
		RelAccount:  models.RelAccount{AccountID: a.AccountID},
		ProjectName: models.ProjectName{Name: name},
	}

	if err := m.FetchByAccountIDAndName(d); err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, status.ProjectNotFound
		} else {
			return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
		}
	}
	if a.AccountID != m.AccountID {
		return nil, status.NoProjectPermission
	}
	return m, nil
}
