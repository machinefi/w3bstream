package account

import (
	"context"
	"strings"

	// "github.com/spruceid/siwe-go"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/spruceid/siwe-go"
)

func FetchOrCreateAccountByEthAddress(ctx context.Context, address types.EthAddress) (*models.Account, *models.AccountIdentity, error) {
	d := types.MustMgrDBExecutorFromContext(ctx)
	g := confid.MustSFIDGeneratorFromContext(ctx)

	var (
		rel    = models.RelAccount{AccountID: g.MustGenSFID()}
		acc    *models.Account
		aci    *models.AccountIdentity
		exists bool
	)

	err := sqlx.NewTasks(d).With(
		// fetch AccountIdentity
		func(db sqlx.DBExecutor) error {
			aci = &models.AccountIdentity{
				AccountIdentityInfo: models.AccountIdentityInfo{
					Type:       enums.ACCOUNT_IDENTITY_TYPE__ETHADDRESS,
					IdentityID: address.String(),
				},
			}
			err := aci.FetchByTypeAndIdentityID(db)
			if err != nil {
				if sqlx.DBErr(err).IsNotFound() {
					exists = false
					return nil
				} else {
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
			} else {
				exists = true
				rel.AccountID = aci.AccountID
				return nil
			}
		},
		// create or fetch Account
		func(db sqlx.DBExecutor) error {
			acc = &models.Account{RelAccount: rel}
			if exists {
				if err := acc.FetchByAccountID(db); err != nil {
					if sqlx.DBErr(err).IsNotFound() {
						return status.AccountNotFound
					}
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
				if acc.State != enums.ACCOUNT_STATE__ENABLED {
					return status.DisabledAccount
				}
				return nil
			} else {
				acc.Role = enums.ACCOUNT_ROLE__DEVELOPER
				acc.State = enums.ACCOUNT_STATE__ENABLED
				if err := acc.Create(db); err != nil {
					return status.DatabaseError.StatusErr().WithDesc(err.Error())
				}
				return nil
			}
		},
		// create AccountIdentity
		func(db sqlx.DBExecutor) error {
			if exists {
				return nil
			}
			aci.RelAccount = rel
			aci.Source = enums.ACCOUNT_SOURCE__SUBMIT
			if err := aci.Create(db); err != nil {
				if sqlx.DBErr(err).IsConflict() {
					return status.AccountIdentityConflict
				}
				return status.DatabaseError.StatusErr().WithDesc(err.Error())
			}
			return nil
		},
		func(d sqlx.DBExecutor) error {
			if exists {
				return nil
			}
			req := operator.CreateReq{
				AccountID:  rel.AccountID,
				Name:       operator.DefaultOperatorName,
				PrivateKey: generateRandomPrivateKey(),
			}
			ctx := types.WithMgrDBExecutor(ctx, d)
			_, err := operator.Create(ctx, &req)
			return err
		},
	).Do()

	_, l := conflog.FromContext(ctx).Start(ctx, "FetchOrCreateAccountByEthAddress")
	defer l.End()

	if err != nil {
		l.Error(err)
		return nil, nil, err
	}
	return acc, aci, nil
}

type LoginByEthAddressReq struct {
	Message   string `json:"message"`   // Message siwe serialized message
	Signature string `json:"signature"` // Signature should have '0x' prefix
}

func ValidateLoginByEthAddress(ctx context.Context, r *LoginByEthAddressReq) (*models.Account, error) {
	_, l := conflog.FromContext(ctx).Start(ctx, "LoginByEthAddress")
	defer l.End()

	msg, err := siwe.ParseMessage(r.Message)
	if err != nil {
		l.Error(err)
		return nil, status.InvalidEthLoginMessage.StatusErr().WithDesc(err.Error())
	}

	if _, err = msg.Verify(r.Signature, nil, nil, nil); err != nil {
		l.Error(err)
		return nil, status.InvalidEthLoginSignature.StatusErr().WithDesc(err.Error())
	}

	address := strings.ToLower(msg.GetAddress().String())

	if lst, ok := types.WhiteListFromContext(ctx); ok {
		if !lst.Validate(address) {
			return nil, status.WhiteListForbidden
		}
	}

	acc, _, err := FetchOrCreateAccountByEthAddress(ctx, types.EthAddress(address))

	if err != nil {
		l.Error(err)
		return nil, err
	}
	return acc, nil
}
