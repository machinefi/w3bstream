package middleware

import (
	"context"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/account"
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

func (v *CurrentAccount) WithProjectContextByName(ctx context.Context, prjName string) (context.Context, error) {
	a := CurrentAccountFromContext(ctx)
	d := types.MustDBExecutorFromContext(ctx)
	m := &models.Project{ProjectInfo: models.ProjectInfo{Name: prjName}}
	if err := m.FetchByName(d); err != nil {
		return ctx, status.CheckDatabaseError(err, "GetProjectByName")
	}
	if a.AccountID != m.AccountID {
		return ctx, status.NoProjectPermission
	}
	ctx = types.WithProject(ctx, m)
	return ctx, nil
}

func (v *CurrentAccount) WithProjectContextByID(ctx context.Context, prjID types.SFID) (context.Context, error) {
	a := CurrentAccountFromContext(ctx)
	d := types.MustDBExecutorFromContext(ctx)
	m := &models.Project{RelProject: models.RelProject{ProjectID: prjID}}
	if err := m.FetchByProjectID(d); err != nil {
		return ctx, status.CheckDatabaseError(err, "GetProjectByProjectID")
	}
	if a.AccountID != m.AccountID {
		return ctx, status.NoProjectPermission
	}
	ctx = types.WithProject(ctx, m)
	return ctx, nil
}

func (v *CurrentAccount) WithAppletContext(ctx context.Context, appletID types.SFID) (context.Context, error) {
	d := types.MustDBExecutorFromContext(ctx)

	app := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	if err := app.FetchByAppletID(d); err != nil {
		return ctx, status.CheckDatabaseError(err, "GetAppletByAppletID")
	}

	_ctx, err := v.WithProjectContextByID(ctx, app.ProjectID)
	if err != nil {
		return ctx, err
	}
	ctx = types.WithApplet(_ctx, app)
	return ctx, nil
}

func (v *CurrentAccount) WithInstanceContext(ctx context.Context, instanceID types.SFID) (context.Context, error) {
	d := types.MustDBExecutorFromContext(ctx)

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
	d := types.MustDBExecutorFromContext(ctx)
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
func (v *CurrentAccount) ValidateProjectPermByPrjName(ctx context.Context, projectName string) (*models.Project, error) {
	d := types.MustDBExecutorFromContext(ctx)
	a := CurrentAccountFromContext(ctx)
	m := &models.Project{ProjectInfo: models.ProjectInfo{Name: projectName}}

	if err := m.FetchByName(d); err != nil {
		return nil, status.CheckDatabaseError(err, "GetProjectByProjectID")
	}
	if a.AccountID != m.AccountID {
		return nil, status.NoProjectPermission
	}
	return m, nil
}
