package login

import (
	"context"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/account"
)

type LoginByUsername struct {
	httpx.MethodPut
	account.LoginByUsernameReq `in:"body"`
}

func (r *LoginByUsername) Path() string { return "" }

func (r *LoginByUsername) Output(ctx context.Context) (interface{}, error) {
	ac, err := account.ValidateLoginByUsername(ctx, &r.LoginByUsernameReq)
	if err != nil {
		return nil, err
	}
	return token(ctx, ac)
}

type LoginByEthAddress struct {
	httpx.MethodPut
	account.LoginByEthAddressReq `in:"body"`
}

func (r *LoginByEthAddress) Path() string { return "/wallet" }

func (r *LoginByEthAddress) Output(ctx context.Context) (interface{}, error) {
	ac, err := account.ValidateLoginByEthAddress(ctx, &r.LoginByEthAddressReq)
	if err != nil {
		return nil, err
	}
	return token(ctx, ac)
}

func token(ctx context.Context, a *models.Account) (*account.LoginRsp, error) {
	j := jwt.MustConfFromContext(ctx)

	tok, err := j.GenerateTokenByPayload(a.AccountID)
	if err != nil {
		return nil, status.InternalServerError.StatusErr().WithDesc(err.Error())
	}

	return &account.LoginRsp{
		AccountID:   a.AccountID,
		AccountRole: a.Role,
		Token:       tok,
		ExpireAt:    types.Timestamp{Time: time.Now().Add(j.ExpIn.Duration())},
		Issuer:      j.Issuer,
	}, nil
}
