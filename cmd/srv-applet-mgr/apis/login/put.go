package login

import (
	"context"
	"time"

	base "github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/conf/jwt"
	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/pkg/types"

	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/modules/account"
)

type Login struct {
	httpx.MethodPut
	InBody `in:"body"`
}

type InBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	AccountID types.SFID     `json:"accountID"`
	Token     string         `json:"token"`
	ExpireAt  base.Timestamp `json:"expireAt"`
	Issuer    string         `json:"issuer"`
}

func (r *Login) Path() string { return "/login" }

func (r *Login) Output(ctx context.Context) (interface{}, error) {
	ac, err := account.ValidateAccountByLogin(ctx, r.Username, r.Password)
	if err != nil {
		return nil, err
	}
	j := jwt.MustConfFromContext(ctx)

	token, err := j.GenerateTokenByPayload(ac.AccountID)
	if err != nil {
		return nil, status.InternalServerError.StatusErr().WithDesc(err.Error())
	}

	return &Response{
		AccountID: ac.AccountID,
		Token:     token,
		ExpireAt:  base.Timestamp{Time: time.Now().Add(j.ExpIn.Duration())},
		Issuer:    j.Issuer,
	}, nil
}
