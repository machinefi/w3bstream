package publisher

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
)

// Create Publisher
type CreatePublisher struct {
	httpx.MethodPost
	publisher.CreateReq `in:"body"`
}

func (r *CreatePublisher) Output(ctx context.Context) (interface{}, error) {
	acc := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := acc.WithProjectContextByName(acc.WithAccount(ctx), middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	return publisher.Create(ctx, &r.CreateReq)
}

type CreateAnonymousPublisher struct {
	httpx.MethodPost
}

func (r *CreateAnonymousPublisher) Path() string { return "/anonymous" }

func (r *CreateAnonymousPublisher) Output(ctx context.Context) (interface{}, error) {
	acc := middleware.MustCurrentAccountFromContext(ctx)
	prjName := middleware.MustProjectName(ctx)
	ctx, err := acc.WithProjectContextByName(acc.WithAccount(ctx), prjName)
	if err != nil {
		return nil, err
	}

	return publisher.CreateAnonymousPublisher(ctx)
}

type UpsertPublisher struct {
	httpx.MethodPost
	publisher.CreateReq `in:"body"`
}

func (r *UpsertPublisher) Path() string { return "/upsert" }

func (r *UpsertPublisher) Output(ctx context.Context) (interface{}, error) {
	acc := middleware.MustCurrentAccountFromContext(ctx)
	ctx, err := acc.WithProjectContextByName(acc.WithAccount(ctx), middleware.MustProjectName(ctx))
	if err != nil {
		return nil, err
	}

	return publisher.Upsert(ctx, &r.CreateReq)
}
