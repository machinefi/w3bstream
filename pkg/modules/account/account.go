package account

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/util"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateAccountByUsernameReq struct {
	Username  string              `json:"username"`
	Role      enums.AccountRole   `json:"role"`
	AvatarURL string              `json:"avatarURL,omitempty" validate:"@url"`
	Password  string              `json:"-"`
	Source    enums.AccountSource `json:"-"`
}

type CreateAccountByUsernameRsp struct {
	*models.Account
	Password string `json:"password"`
}

func CreateAccountByUsername(ctx context.Context, r *CreateAccountByUsernameReq) (*CreateAccountByUsernameRsp, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	g := confid.MustSFIDGeneratorFromContext(ctx)

	rel := &models.RelAccount{AccountID: g.MustGenSFID()}
	if r.Source == 0 {
		r.Source = enums.ACCOUNT_SOURCE__SUBMIT
	}
	acc := (*models.Account)(nil)
	passwd := r.Password

	err := sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			acc = &models.Account{
				RelAccount: *rel,
				AccountInfo: models.AccountInfo{
					State:              enums.ACCOUNT_STATE__ENABLED,
					Role:               r.Role,
					Avatar:             r.AvatarURL,
					OperatorPrivateKey: generateRandomPrivateKey(),
				},
			}
			if err := acc.Create(db); err != nil {
				return status.CheckDatabaseError(err, "CreateAccount")
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			if err := (&models.AccountIdentity{
				RelAccount: *rel,
				AccountIdentityInfo: models.AccountIdentityInfo{
					Type:       enums.ACCOUNT_IDENTITY_TYPE__USERNAME,
					IdentityID: r.Username,
					Source:     r.Source,
				},
			}).Create(db); err != nil {
				return status.CheckDatabaseError(err, "CreateAccountIdentity")
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			if passwd == "" {
				passwd = string(util.GenRandomPassword(8, 3))
			}
			if err := (&models.AccountPassword{
				RelAccount:         *rel,
				RelAccountPassword: models.RelAccountPassword{PasswordID: g.MustGenSFID()},
				AccountPasswordData: models.AccountPasswordData{
					Type: enums.PASSWORD_TYPE__LOGIN,
					Password: util.HashOfAccountPassword(
						rel.AccountID.String(), passwd,
					),
				},
			}).Create(db); err != nil {
				return status.CheckDatabaseError(err, "Create")
			}
			return nil
		},
	).Do()

	_, l = conflog.FromContext(ctx).Start(ctx, "CreateAccountByUsername")
	defer l.End()

	if err != nil {
		l.Error(err)
		return nil, err
	}
	return &CreateAccountByUsernameRsp{
		Account:  acc,
		Password: passwd,
	}, nil
}

type UpdatePasswordReq struct {
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
}

func UpdateAccountPassword(ctx context.Context, accountID types.SFID, r *UpdatePasswordReq) error {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	var (
		rel = models.RelAccount{AccountID: accountID}
		acc *models.Account
		aci *models.AccountIdentity
		ap  *models.AccountPassword
	)

	err := sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			acc = &models.Account{RelAccount: rel}
			if err := acc.FetchByAccountID(db); err != nil {
				return status.CheckDatabaseError(err, "FetchAccount")
			}
			if acc.State != enums.ACCOUNT_STATE__ENABLED {
				return status.DisabledAccount
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			aci = &models.AccountIdentity{
				RelAccount: rel,
				AccountIdentityInfo: models.AccountIdentityInfo{
					Type: enums.ACCOUNT_IDENTITY_TYPE__USERNAME,
				},
			}
			if err := aci.FetchByAccountIDAndType(db); err != nil {
				return status.CheckDatabaseError(err, "FetchAccountIdentity")
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			ap = &models.AccountPassword{
				RelAccount: rel,
				AccountPasswordData: models.AccountPasswordData{
					Type: enums.PASSWORD_TYPE__LOGIN,
				},
			}
			if err := ap.FetchByAccountIDAndType(db); err != nil {
				return status.CheckDatabaseError(err, "FetchAccountPassword")
			}
			if ap.Password != util.HashOfAccountPassword(accountID.String(), r.OldPassword) {
				return status.InvalidOldPassword
			}
			if r.OldPassword == r.Password {
				return status.InvalidNewPassword
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			ap.Password = util.HashOfAccountPassword(accountID.String(), r.Password)
			if err := ap.UpdateByAccountIDAndType(db); err != nil {
				return status.CheckDatabaseError(err, "UpdatePassword")
			}
			return nil
		},
	).Do()

	_, l = l.Start(ctx, "UpdateAccountPassword")
	defer l.End()

	if err != nil {
		l.Error(err)
		return err
	}
	return nil
}

type LoginByUsernameReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRsp struct {
	AccountID   types.SFID        `json:"accountID"`
	AccountRole enums.AccountRole `json:"accountRole"`
	Token       string            `json:"token"`
	ExpireAt    types.Timestamp   `json:"expireAt"`
	Issuer      string            `json:"issuer"`
}

func ValidateLoginByUsername(ctx context.Context, r *LoginByUsernameReq) (*models.Account, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	var (
		rel models.RelAccount
		aci *models.AccountIdentity
		acc *models.Account
	)

	err := sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			aci = &models.AccountIdentity{
				RelAccount: models.RelAccount{},
				AccountIdentityInfo: models.AccountIdentityInfo{
					Type:       enums.ACCOUNT_IDENTITY_TYPE__USERNAME,
					IdentityID: r.Username,
				},
			}
			if err := aci.FetchByTypeAndIdentityID(db); err != nil {
				return status.CheckDatabaseError(err, "FetchAccountIdentity")
			}
			rel.AccountID = aci.AccountID
			return nil
		},
		func(db sqlx.DBExecutor) error {
			acc = &models.Account{RelAccount: rel}
			if err := acc.FetchByAccountID(db); err != nil {
				return status.CheckDatabaseError(err, "FetchAccount")
			}
			if acc.State != enums.ACCOUNT_STATE__ENABLED {
				return status.DisabledAccount
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			ap := &models.AccountPassword{
				RelAccount: rel,
				AccountPasswordData: models.AccountPasswordData{
					Type: enums.PASSWORD_TYPE__LOGIN,
				},
			}
			if err := ap.FetchByAccountIDAndType(db); err != nil {
				return status.CheckDatabaseError(err, "FetchAccountPassword")
			}
			if util.HashOfAccountPassword(acc.AccountID.String(), r.Password) != ap.Password {
				return status.InvalidPassword
			}
			return nil
		},
	).Do()

	_, l = l.Start(ctx, "ValidateAccountByLogin")
	defer l.End()

	if err != nil {
		l.Error(err)
		return nil, err
	}
	return acc, nil
}

func GetAccountByAccountID(ctx context.Context, accountID types.SFID) (*models.Account, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
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
	_, l := conflog.FromContext(ctx).Start(ctx, "CreateAdminIfNotExist")
	defer l.End()

	ret, err := CreateAccountByUsername(ctx, &CreateAccountByUsernameReq{
		Username: "admin",
		Role:     enums.ACCOUNT_ROLE__ADMIN,
		Password: "iotex.W3B.admin",
		Source:   enums.ACCOUNT_SOURCE__INIT,
	})
	if err != nil {
		if sqlx.DBErr(err).IsConflict() {
			l.Info("admin already exists, default password: `%s`", "iotex.W3B.admin")
			return "", nil
		} else if se, ok := err.(*statusx.StatusErr); ok && se.Code == status.Conflict.Code() {
			l.Info("admin already exists, default password: `%s`", "iotex.W3B.admin")
			return "", nil
		} else {
			l.Error(err)
			return "", err
		}
	}
	return ret.Password, nil
}

func generateRandomPrivateKey() string {
	priKey, err := crypto.GenerateKey()
	if err != nil {
		return ""
	}
	return hex.EncodeToString(crypto.FromECDSA(priKey))
}
