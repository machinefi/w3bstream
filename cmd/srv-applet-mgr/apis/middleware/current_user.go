package middleware

import (
	"context"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
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

func MustCurrentAccountFromContext(ctx context.Context) *CurrentAccount {
	ca, ok := ctx.Value(contextAccountAuthKey).(*CurrentAccount)
	must.BeTrue(ok)
	return ca
}

func CurrentAccountFromContext(ctx context.Context) (*CurrentAccount, bool) {
	ca, ok := ctx.Value(contextAccountAuthKey).(*CurrentAccount)
	return ca, ok
}

type CurrentAccount struct {
	models.Account
}

func (v *CurrentAccount) WithAccount(ctx context.Context) context.Context {
	return types.WithAccount(ctx, &v.Account)
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
	app, err := applet.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithApplet(ctx, app)

	ctx, err = v.WithProjectContextBySFID(ctx, app.ProjectID)
	if err != nil {
		return nil, err
	}

	ins, err := deploy.GetByAppletSFID(ctx, app.AppletID)
	if err != nil {
		se, ok := statusx.IsStatusErr(err)
		if !ok || !se.Is(status.InstanceNotFound) {
			return nil, err
		}
	} else {
		ctx = types.WithInstance(ctx, ins)
	}

	res, err := resource.GetBySFID(ctx, app.ResourceID)
	if err != nil {
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

	if ctx, err = v.WithProjectContextBySFID(ctx, app.ProjectID); err != nil {
		return nil, err
	}

	if res, err = resource.GetBySFID(ctx, app.ResourceID); err != nil {
		return nil, err
	}
	return types.WithResource(ctx, res), nil
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

func (v *CurrentAccount) WithResourceOwnerContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	_, err := resource.GetOwnerByAccountAndSFID(ctx, v.AccountID, id)
	if err != nil {
		return nil, err
	}
	// TODO if needed add ownership context
	return types.WithAccount(ctx, &v.Account), nil
}
