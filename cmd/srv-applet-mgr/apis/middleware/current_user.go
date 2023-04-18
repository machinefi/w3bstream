package middleware

import (
	"context"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/account"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
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

// WithProjectContextByName With project context by project name(in database)
func (v *CurrentAccount) WithProjectContextByName(ctx context.Context, name string) (context.Context, error) {
	prj, err := project.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if v.AccountID != prj.AccountID {
		return nil, status.NoProjectPermission
	}
	return types.WithProject(ctx, prj), nil
}

// WithProjectContextBySFID With project context by project SFID
func (v *CurrentAccount) WithProjectContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	prj, err := project.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	if v.AccountID != prj.AccountID {
		return nil, status.NoProjectPermission
	}
	return types.WithProject(ctx, prj), nil
}

// WithAppletContextBySFID With full contexts by applet SFID
func (v *CurrentAccount) WithAppletContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	var (
		app *models.Applet
		ins *models.Instance
		res *models.Resource
		err error
	)
	if app, err = applet.GetBySFID(ctx, id); err != nil {
		return nil, err
	}
	ctx = types.WithApplet(ctx, app)

	if ctx, err = v.WithProjectContextBySFID(ctx, app.ProjectID); err != nil {
		return nil, err
	}

	if ins, err = deploy.GetByAppletSFID(ctx, app.AppletID); err != nil {
		se, ok := statusx.IsStatusErr(err)
		if !ok || !se.Is(status.InstanceNotFound) {
			return nil, err
		}
	} else {
		ctx = types.WithInstance(ctx, ins)
	}

	if res, err = resource.GetBySFID(ctx, app.ResourceID); err != nil {
		return nil, err
	}

	return types.WithResource(ctx, res), nil
}

// WithInstanceContextBySFID With full contexts by instance SFID
func (v *CurrentAccount) WithInstanceContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	var (
		ins *models.Instance
		app *models.Applet
		res *models.Resource
		err error
	)
	if ins, err = deploy.GetBySFID(ctx, id); err != nil {
		return nil, err
	}
	ctx = types.WithInstance(ctx, ins)

	if app, err = applet.GetBySFID(ctx, ins.AppletID); err != nil {
		return nil, err
	}
	ctx = types.WithApplet(ctx, app)

	if res, err = resource.GetBySFID(ctx, app.ResourceID); err != nil {
		return nil, err
	}
	ctx = types.WithResource(ctx, res)

	if ctx, err = v.WithProjectContextBySFID(ctx, app.ProjectID); err != nil {
		return nil, err
	}
	return types.WithInstance(ctx, ins), nil
}

func (v *CurrentAccount) WithStrategyBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	sty, err := strategy.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithStrategy(ctx, sty)
	return v.WithProjectContextBySFID(ctx, sty.ProjectID)
}

func (v *CurrentAccount) WithPublisherBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	pub, err := publisher.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithPublisher(ctx, pub)
	return v.WithProjectContextBySFID(ctx, pub.ProjectID)
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

	if err := m.FetchByName(d); err != nil {
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
