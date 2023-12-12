package middleware

import (
	"context"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/account"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/modules/strategy"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ContextPublisherAuth struct {
	httpx.MethodGet
}

var ctxPublisherAuthKey = reflect.TypeOf(ContextAccountAuth{}).String()

func (r *ContextPublisherAuth) ContextKey() string { return ctxPublisherAuthKey }

func (r *ContextPublisherAuth) Output(ctx context.Context) (interface{}, error) {
	ctx, l := logr.Start(ctx, "middleware.ContextPublisherAuth.Output")
	defer l.End()

	pl, err := ParseJwtAuthContentFromContext(ctx)
	if err != nil {
		return nil, status.InvalidAuthPublisherID.StatusErr().WithDesc(err.Error())
	}

	switch pl.IdentityType {
	case enums.ACCESS_KEY_IDENTITY_TYPE__ACCOUNT:
		ca, err := account.GetAccountByAccountID(ctx, pl.IdentityID)
		if err != nil {
			return nil, err
		}
		return &CurrentAccount{*ca}, nil
	default: // unknown or publisher
		cp, err := publisher.GetBySFID(ctx, pl.IdentityID)
		if err != nil {
			return nil, err
		}
		return &CurrentPublisher{cp}, nil
	}
}

func PublisherFromContext(ctx context.Context) (*CurrentPublisher, bool) {
	p, ok := ctx.Value(ctxPublisherAuthKey).(*CurrentPublisher)
	return p, ok
}

func MustPublisher(ctx context.Context) *CurrentPublisher {
	p, ok := ctx.Value(ctxPublisherAuthKey).(*CurrentPublisher)
	must.BeTrue(ok)
	return p
}

func MaybePublisher(ctx context.Context) (*CurrentPublisher, bool) {
	v := ctx.Value(ctxPublisherAuthKey)
	p, ok := v.(*CurrentPublisher)
	return p, ok
}

type CurrentPublisher struct {
	*models.Publisher
}

func (v *CurrentPublisher) WithProjectContext(ctx context.Context) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentPublisher.WithProjectContext")
	defer l.End()

	p := MustPublisher(ctx)

	prj, err := project.GetBySFID(ctx, p.ProjectID)
	if err != nil {
		return nil, err
	}
	return types.WithProject(ctx, prj), nil
}

func (v *CurrentPublisher) WithAccountContext(ctx context.Context) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentPublisher.WithAccountContext")
	defer l.End()

	var (
		err error
		acc *models.Account
	)
	if ctx, err = v.WithProjectContext(ctx); err != nil {
		return nil, err
	}
	prj := types.MustProjectFromContext(ctx)

	if acc, err = account.GetAccountByAccountID(ctx, prj.AccountID); err != nil {
		return nil, err
	}
	return types.WithAccount(ctx, acc), nil
}

func (v *CurrentPublisher) WithStrategiesByChanAndType(ctx context.Context, ch, tpe string) (context.Context, error) {
	ctx, l := logr.Start(ctx, "CurrentPublisher.WithStrategiesByChanAndType")
	defer l.End()
	var (
		err error
		res []*types.StrategyResult
	)
	prj, ok := types.ProjectFromContext(ctx)
	if !ok {
		if ctx, err = v.WithProjectContext(ctx); err != nil {
			return nil, err
		}
		prj = types.MustProjectFromContext(ctx)
	}

	if prj.Name != ch {
		return nil, status.InvalidEventChannel
	}

	res, err = strategy.FilterByProjectAndEvent(ctx, prj.ProjectID, tpe)
	if err != nil {
		return nil, err
	}
	return types.WithStrategyResults(ctx, res), nil
}
