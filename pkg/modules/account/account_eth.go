package account

import (
	"context"
	"strings"

	// "github.com/spruceid/siwe-go"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CreateAccountByEthAddressReq struct {
	Address types.EthAddress `json:"address"          validate:"@ethAddress"`
	Avatar  string           `json:"avatar,omitempty" validate:"@url"`
}

type CreateAccountByEthAddressRsp struct {
	*models.Account
	Nonce string `json:"nonce"`
}

func CreateAccountByEthAddress(ctx context.Context, r *CreateAccountByEthAddressReq) (*CreateAccountByEthAddressRsp, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	g := confid.MustSFIDGeneratorFromContext(ctx)

	var (
		rel = &models.RelAccount{AccountID: g.MustGenSFID()}
		acc *models.Account
		aci *models.AccountIdentity
	)

	err := sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			acc = &models.Account{
				RelAccount: *rel,
				AccountInfo: models.AccountInfo{
					State:  enums.ACCOUNT_STATE__ENABLED,
					Role:   enums.ACCOUNT_ROLE__DEVELOPER,
					Avatar: r.Avatar,
				},
			}
			if err := acc.Create(db); err != nil {
				return status.CheckDatabaseError(err, "CreateAccount")
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			nonce, err := GenerateNonce()
			if err != nil {
				return status.GenNonceFailed.StatusErr().WithMsg(err.Error())
			}
			aci = &models.AccountIdentity{
				RelAccount: *rel,
				AccountIdentityInfo: models.AccountIdentityInfo{
					Type:       enums.ACCOUNT_IDENTITY_TYPE__ETHADDRESS,
					IdentityID: r.Address.String(),
					Source:     enums.ACCOUNT_SOURCE__SUBMIT,
					Meta:       models.Meta{models.AccountIdentityMetaKey_EthAddress_Nonce: nonce},
				},
			}
			if err := aci.Create(db); err != nil {
				return status.CheckDatabaseError(err, "CreateAccountIdentity")
			}
			return nil
		},
	).Do()

	_, l := conflog.FromContext(ctx).Start(ctx, "CreateAccountByEthAddress")
	defer l.End()

	if err != nil {
		l.Error(err)
		return nil, err
	}
	return &CreateAccountByEthAddressRsp{
		Account: acc,
		Nonce:   aci.Meta[models.AccountIdentityMetaKey_EthAddress_Nonce],
	}, nil
}

func GetAccountByEthAddress(ctx context.Context, address types.EthAddress) (acc *models.Account, aci *models.AccountIdentity, err error) {
	d := types.MustMgrDBExecutorFromContext(ctx)

	err = sqlx.NewTasks(d).With(
		func(db sqlx.DBExecutor) error {
			aci = &models.AccountIdentity{
				AccountIdentityInfo: models.AccountIdentityInfo{
					Type:       enums.ACCOUNT_IDENTITY_TYPE__ETHADDRESS,
					IdentityID: address.String(),
				},
			}
			if err := aci.FetchByTypeAndIdentityID(db); err != nil {
				return status.CheckDatabaseError(err, "FetchAccountIdentity")
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			if _, ok := aci.Meta[models.AccountIdentityMetaKey_EthAddress_Nonce]; ok {
				return nil
			}
			nonce, err := GenerateNonce()
			if err != nil {
				return status.GenNonceFailed.StatusErr().WithMsg(err.Error())
			}
			if aci.Meta == nil {
				aci.Meta = models.Meta{}
			}
			aci.Meta[models.AccountIdentityMetaKey_EthAddress_Nonce] = nonce
			if err = aci.UpdateByTypeAndIdentityID(db); err != nil {
				return status.CheckDatabaseError(err, "UpdateAccountIdentity")
			}
			return nil
		},
		func(db sqlx.DBExecutor) error {
			acc = &models.Account{}
			acc.AccountID = aci.AccountID
			if err := acc.FetchByAccountID(db); err != nil {
				return status.CheckDatabaseError(err, "FetchAccount")
			}
			if acc.State != enums.ACCOUNT_STATE__ENABLED {
				return status.DisabledAccount
			}
			return nil
		},
	).Do()
	if err != nil {
		return nil, nil, err
	}
	return
}

type GetNonceRsp struct {
	*models.Account
	Nonce string `json:"nonce"`
}

func GetNonce(ctx context.Context, address types.EthAddress) (*GetNonceRsp, error) {
	acc, aci, err := GetAccountByEthAddress(ctx, address)

	_, l := conflog.FromContext(ctx).Start(ctx, "GetNonce")
	if err != nil {
		l.Error(err)
		return nil, err
	}
	return &GetNonceRsp{
		Account: acc,
		Nonce:   aci.Meta[models.AccountIdentityMetaKey_EthAddress_Nonce],
	}, nil
}

type LoginByEthAddressReq struct {
	Address   types.EthAddress `json:"address" validate:"@ethAddress"`
	Nonce     string           `json:"nonce"`
	Signature string           `json:"signature"` // Signature should have '0x' prefix
}

func ValidateLoginByEthAddress(ctx context.Context, r *LoginByEthAddressReq) (*models.Account, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)

	acc, aci, err := GetAccountByEthAddress(ctx, r.Address)

	_, l := conflog.FromContext(ctx).Start(ctx, "LoginByEthAddress")
	defer l.End()

	if err != nil {
		l.Error(err)
		return nil, err
	}

	nonce := aci.Meta[models.AccountIdentityMetaKey_EthAddress_Nonce]
	if r.Nonce != nonce {
		l.Error(status.InvalidNonce)
		return nil, status.InvalidNonce
	}

	signature, err := hexutil.Decode(r.Signature)
	if err != nil {
		l.Error(err)
		return nil, status.InvalidSignature
	}
	signature[crypto.RecoveryIDOffset] -= 27

	message := accounts.TextHash([]byte(nonce))
	recovered, err := crypto.SigToPub(message, signature)
	if err != nil {
		return nil, status.InvalidNonceOrSignature.StatusErr().WithDesc(err.Error())
	}

	address := crypto.PubkeyToAddress(*recovered)
	if r.Address.String() != strings.ToLower(address.Hex()) {
		return nil, status.InvalidEthAddress
	}

	nonce, err = GenerateNonce()
	if err != nil {
		return nil, status.GenNonceFailed.StatusErr().WithDesc(err.Error())
	}
	aci.Meta[models.AccountIdentityMetaKey_EthAddress_Nonce] = nonce
	if err = aci.UpdateByID(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "UpdateNonce")
	}
	return acc, nil
}
