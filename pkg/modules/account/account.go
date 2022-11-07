package account

import (
	"context"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/util"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateAccountByUsernameReq struct {
	Username string `json:"username"`
}

func CreateAccountByUsername(ctx context.Context, r *CreateAccountByUsernameReq) (*models.Account, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	g := confid.MustSFIDGeneratorFromContext(ctx)

	accountID := g.MustGenSFID()
	m := &models.Account{
		RelAccount: models.RelAccount{AccountID: accountID},
		AccountInfo: models.AccountInfo{
			Username:     r.Username,
			IdentityType: enums.ACCOUNT_IDENTITY_TYPE__USERNAME,
			State:        enums.ACCOUNT_STATE__ENABLED,
			Password: models.AccountPassword{
				Type: enums.PASSWORD_TYPE__LOGIN,
				Password: util.HashOfAccountPassword(
					accountID.String(),
					string(util.GenRandomPassword(8, 3)),
				),
			},
		},
	}

	_, l = l.Start(ctx, "CreateAccountByUsername")
	defer l.End()

	if err := m.Create(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "CreateAccount")
	}
	return m, nil
}

type UpdatePasswordReq struct {
	Password string `json:"password"`
}

func UpdateAccountPassword(ctx context.Context, accountID types.SFID, password string) error {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	m := &models.Account{RelAccount: models.RelAccount{AccountID: accountID}}

	_, l = l.Start(ctx, "UpdateAccountPassword")
	defer l.End()

	if err := m.FetchByAccountID(d); err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "FetchByAccountID")
	}

	// TODO should check account type and password type
	m.Password.Password = util.HashOfAccountPassword(accountID.String(), password)

	if err := m.UpdateByAccountID(d); err != nil {
		l.Error(err)
		return status.CheckDatabaseError(err, "UpdateByAccountID")
	}
	return nil
}

func ValidateAccountByLogin(ctx context.Context, username, password string) (*models.Account, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	m := &models.Account{}
	m.Username = username

	_, l = l.Start(ctx, "ValidateAccountByLogin")
	defer l.End()

	if err := m.FetchByUsername(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "FetchByUsername")
	}
	if m.Password.Password == util.HashOfAccountPassword(m.AccountID.String(), password) {
		return m, nil
	}
	l.Warn(errors.New("wrong password"))
	return nil, status.Unauthorized.StatusErr().WithDesc("invalid password")
}

func GetAccountByAccountID(ctx context.Context, accountID types.SFID) (*models.Account, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	m := &models.Account{RelAccount: models.RelAccount{AccountID: accountID}}
	_, l = l.Start(ctx, "GetAccountByAccountID")
	defer l.End()

	err := m.FetchByAccountID(d)
	if err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "FetchByAccountID")
	}
	return m, err
}

func CreateAdminIfNotExist(ctx context.Context) (string, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	accountID := idg.MustGenSFID()
	password := "iotex.W3B.admin"
	m := &models.Account{
		RelAccount: models.RelAccount{AccountID: accountID},
		AccountInfo: models.AccountInfo{
			Username:     "admin",
			IdentityType: enums.ACCOUNT_IDENTITY_TYPE__BUILTIN,
			State:        enums.ACCOUNT_STATE__ENABLED,
			Password: models.AccountPassword{
				Type: enums.PASSWORD_TYPE__LOGIN,
				Password: util.HashOfAccountPassword(
					accountID.String(),
					password,
				),
				Scope: "admin",
				Desc:  "builtin password",
			},
		},
	}

	_, l = l.Start(ctx, "CreateAdminIfNotExist")
	defer l.End()

	results := make([]models.Account, 0)
	err := d.QueryAndScan(builder.Select(nil).
		From(
			d.T(m),
			builder.Where(
				builder.And(
					m.ColUsername().Eq("admin"),
					m.ColIdentityType().Eq(enums.ACCOUNT_IDENTITY_TYPE__BUILTIN),
				),
			),
		), &results)
	if err != nil {
		l.Error(err)
		return "", status.CheckDatabaseError(err, "FetchAdminAccount")
	}
	if len(results) > 0 {
		l.Info("admin already exists, default password: `%s`", "iotex.W3B.admin")
		return results[0].Password.Password, nil
	}
	if err = m.Create(d); err != nil {
		l.Error(err)
		return "", err
	}
	l.Info("admin created")
	return password, nil
}
