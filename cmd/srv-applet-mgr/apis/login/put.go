package login

import (
	"context"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/account"
)

type LoginByUsername struct {
	httpx.MethodPut
	account.LoginByUsernameReq `in:"body"`
}

func (r *LoginByUsername) Path() string { return "/login" }

func (r *LoginByUsername) Output(ctx context.Context) (interface{}, error) {
	ac, err := account.ValidateLoginByUsername(ctx, &r.LoginByUsernameReq)
	if err != nil {
		return nil, err
	}
	return token(ctx, ac.AccountID)
}

type LoginByEthAddress struct {
	httpx.MethodPut
	account.LoginByEthAddressReq `in:"body"`
}

func (r *LoginByEthAddress) Path() string { return "/eth_login/" }

func (r *LoginByEthAddress) Output(ctx context.Context) (interface{}, error) {
	ac, err := account.ValidateLoginByEthAddress(ctx, &r.LoginByEthAddressReq)
	if err != nil {
		return nil, err
	}
	return token(ctx, ac.AccountID)
}

func token(ctx context.Context, accountID types.SFID) (*account.LoginRsp, error) {
	j := jwt.MustConfFromContext(ctx)

	tok, err := j.GenerateTokenByPayload(accountID)
	if err != nil {
		return nil, status.InternalServerError.StatusErr().WithDesc(err.Error())
	}

	return &account.LoginRsp{
		AccountID: accountID,
		Token:     tok,
		ExpireAt:  types.Timestamp{Time: time.Now().Add(j.ExpIn.Duration())},
		Issuer:    j.Issuer,
	}, nil
}
