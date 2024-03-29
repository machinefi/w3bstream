package middleware

import (
	"context"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/account"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/modules/cronjob"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/modules/resource"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ContextAccountAuth struct {
	httpx.MethodGet
}

var contextAccountAuthKey = reflect.TypeOf(ContextAccountAuth{}).String()

func (r *ContextAccountAuth) ContextKey() string { return contextAccountAuthKey }

func (r *ContextAccountAuth) Output(ctx context.Context) (interface{}, error) {
	ctx, l := logr.Start(ctx, "middleware.ContextAccountAuth.Output")
	defer l.End()

	pl, err := ParseJwtAuthContentFromContext(ctx)
	if err != nil {
		return nil, status.InvalidAuthAccountID.StatusErr().WithDesc(err.Error())
	}

	switch pl.IdentityType {
	// must unknown or account
	case enums.ACCESS_KEY_IDENTITY_TYPE__ACCOUNT, enums.ACCESS_KEY_IDENTITY_TYPE_UNKNOWN:
		ca, err := account.GetAccountByAccountID(ctx, pl.IdentityID)
		if err != nil {
			return nil, err
		}
		return &CurrentAccount{*ca}, nil
	default:
		return nil, status.InvalidAuthValue
	}
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

func (v *CurrentAccount) CheckRole(role enums.AccountRole) (*CurrentAccount, bool) {
	if v.Role == role {
		return v, true
	}
	return nil, false
}

func (v *CurrentAccount) WithAccount(ctx context.Context) context.Context {
	return types.WithAccount(ctx, &v.Account)
}

// WithProjectContextByName With project context by project name(in database)
func (v *CurrentAccount) WithProjectContextByName(ctx context.Context, name string) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithProjectContextByName")
	defer l.End()

	prj, err := project.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if v.Role == enums.ACCOUNT_ROLE__ADMIN || v.AccountID == prj.AccountID {
		return types.WithProject(ctx, prj), nil
	}

	return nil, status.NoProjectPermission
}

// WithProjectContextBySFID With project context by project SFID
func (v *CurrentAccount) WithProjectContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithProjectContextBySFID")
	defer l.End()

	prj, err := project.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	if v.Role == enums.ACCOUNT_ROLE__ADMIN || v.AccountID == prj.AccountID {
		return types.WithProject(ctx, prj), nil
	}

	return nil, status.NoProjectPermission
}

// WithAppletContextBySFID With applet contexts by applet SFID
func (v *CurrentAccount) WithAppletContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithAppletContextBySFID")
	defer l.End()

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

	return ctx, nil
}

func (v *CurrentAccount) WithResourceContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithResourceContextBySFID")
	defer l.End()

	res, err := resource.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}

	return types.WithResource(ctx, res), nil
}

// WithInstanceContextBySFID With full contexts by instance SFID
func (v *CurrentAccount) WithInstanceContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithInstanceContextBySFID")
	defer l.End()

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
	ctx, l := logr.Start(ctx, "CurrentAccount.WithStrategyBySFID")
	defer l.End()

	sty, err := strategy.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithStrategy(ctx, sty)
	return v.WithProjectContextBySFID(ctx, sty.ProjectID)
}

func (v *CurrentAccount) WithPublisherBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithPublisherBySFID")
	defer l.End()

	pub, err := publisher.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithPublisher(ctx, pub)
	return v.WithProjectContextBySFID(ctx, pub.ProjectID)
}

func (v *CurrentAccount) WithResourceOwnerContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithResourceOwnerContextBySFID")
	defer l.End()

	ship, err := resource.GetOwnerByAccountAndSFID(ctx, v.AccountID, id)
	if err != nil {
		return nil, err
	}
	if v.Role == enums.ACCOUNT_ROLE__ADMIN || v.AccountID == ship.AccountID {
		return types.WithResourceOwnership(ctx, ship), nil
	}

	return nil, status.NoResourcePermission
}

func (v *CurrentAccount) WithCronJobBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithCronJobBySFID")
	defer l.End()

	cronJob, err := cronjob.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithCronJob(ctx, cronJob)
	return v.WithProjectContextBySFID(ctx, cronJob.ProjectID)
}

func (v *CurrentAccount) WithOperatorBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithOperatorBySFID")
	defer l.End()

	op, err := operator.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	if v.AccountID != op.AccountID {
		return nil, status.NoOperatorPermission
	}
	return types.WithOperator(ctx, op), nil
}

func (v *CurrentAccount) WithContractLogBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithContractLogBySFID")
	defer l.End()

	cl, err := blockchain.GetContractLogBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithContractLog(ctx, cl)
	return v.WithProjectContextByName(ctx, cl.ProjectName)
}

func (v *CurrentAccount) WithChainHeightBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithChainHeightBySFID")
	defer l.End()

	h, err := blockchain.GetChainHeightBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithChainHeight(ctx, h)
	return v.WithProjectContextByName(ctx, h.ProjectName)
}

func (v *CurrentAccount) WithChainTxBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithChainTxBySFID")
	defer l.End()

	t, err := blockchain.GetChainTxBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithChainTx(ctx, t)
	return v.WithProjectContextByName(ctx, t.ProjectName)
}

func (v *CurrentAccount) WithTrafficLimitContextBySFID(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithTrafficLimitContextBySFID")
	defer l.End()

	traffic, err := trafficlimit.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = types.WithTrafficLimit(ctx, traffic)
	return v.WithProjectContextBySFID(ctx, traffic.ProjectID)
}

func (v *CurrentAccount) WithTrafficLimitContextBySFIDAndProjectName(ctx context.Context, id types.SFID) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentAccount.WithTrafficLimitContextBySFIDAndProjectName")
	defer l.End()

	traffic, err := trafficlimit.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	project := types.MustProjectFromContext(ctx)
	if v.Role == enums.ACCOUNT_ROLE__ADMIN || traffic.ProjectID == project.ProjectID {
		return types.WithTrafficLimit(ctx, traffic), nil
	}

	return nil, status.NoProjectPermission
}
