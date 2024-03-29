package account_access_key

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

// ListAccountAccessKey get api access key list under current account
type ListAccountAccessKey struct {
	httpx.MethodGet
	access_key.ListReq
}

func (r *ListAccountAccessKey) Path() string { return "/datalist" }

func (r *ListAccountAccessKey) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	r.AccountID = ca.AccountID
	return access_key.List(ctx, &r.ListReq)
}

type ListAccessGroupMetas struct {
	httpx.MethodGet `summary:"List operator group metas"`
}

func (r *ListAccessGroupMetas) Path() string { return "/operator_group_metas" }

func (r *ListAccessGroupMetas) Output(_ context.Context) (interface{}, error) {
	return access_key.OperatorGroupMetaList(), nil
}

type GetAccessKeyByName struct {
	httpx.MethodGet
	Name string `in:"path" name:"name"`
}

func (r *GetAccessKeyByName) Path() string { return "/data/:name" }

func (r *GetAccessKeyByName) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.MustCurrentAccountFromContext(ctx)
	return access_key.GetByName(ca.WithAccount(ctx), r.Name)
}
