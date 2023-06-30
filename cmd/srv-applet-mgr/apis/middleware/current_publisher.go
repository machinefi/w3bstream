package middleware

import (
	"context"
	"fmt"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
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
	vRaw := jwt.AuthFromContext(ctx)
	if _, ok := vRaw.(string); !ok {
		var content []byte
		switch v1 := vRaw.(type) {
		case []byte:
			content = v1
		case string:
			content = []byte(v1)
		case fmt.Stringer:
			content = []byte(v1.String())
		default:
			return nil, status.InvalidAuthValue
		}

		accountID := types.SFID(0)
		if err := accountID.UnmarshalText(content); err != nil {
			return nil, status.InvalidAuthAccountID
		}
		ca, err := account.GetAccountByAccountID(ctx, accountID)
		if err != nil {
			return nil, err
		}
		return &CurrentAccount{*ca}, nil
	}
	v := vRaw.(string)
	id := types.SFID(0)
	if err := id.UnmarshalText([]byte(v)); err != nil {
		return nil, status.InvalidAuthPublisherID
	}
	cp, err := publisher.GetBySFID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &CurrentPublisher{cp}, nil
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

type CurrentPublisher struct {
	*models.Publisher
}

func (v *CurrentPublisher) WithProjectContext(ctx context.Context) (context.Context, error) {
	p := MustPublisher(ctx)

	prj, err := project.GetBySFID(ctx, p.ProjectID)
	if err != nil {
		return nil, err
	}
	return types.WithProject(ctx, prj), nil
}

func (v *CurrentPublisher) WithAccountContext(ctx context.Context) (context.Context, error) {
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
