package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// AccountIdentity account identity
// @def primary ID
// @def unique_index ui_account_identity AccountID Type
// @def unique_index ui_identity_id      Type IdentityID
// @def index        i_identity_id       IdentityID
// @def index        i_source            Source
//
//go:generate toolkit gen model AccountIdentity --database DB
type AccountIdentity struct {
	datatypes.PrimaryID
	RelAccount
	AccountIdentityInfo
	datatypes.OperationTimesWithDeleted
}

type AccountIdentityInfo struct {
	Type       enums.AccountIdentityType `db:"f_type"            json:"type"`
	IdentityID string                    `db:"f_identity_id"     json:"identityID"`
	Source     enums.AccountSource       `db:"f_source"          json:"source"`
	Meta       Meta                      `db:"f_meta,default=''" json:"meta"`
}

const (
	AccountIdentityMetaKey_EthAddress_Nonce string = "nonce"
)
