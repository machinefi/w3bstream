package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Account w3bstream account
// @def primary                    ID
// @def unique_index UI_account_id AccountID
//
//go:generate toolkit gen model Account --database DB
type Account struct {
	datatypes.PrimaryID
	RelAccount
	AccountInfo
	datatypes.OperationTimesWithDeleted
}

type RelAccount struct {
	// @rel Account.AccountID
	AccountID types.SFID `db:"f_account_id" json:"accountID"` // AccountID  account id
}

type AccountInfo struct {
	Role               enums.AccountRole  `db:"f_role,default=2"    json:"role"`
	State              enums.AccountState `db:"f_state,default='1'" json:"state"`
	Avatar             string             `db:"f_avatar,default=''" json:"avatar,omitempty"`
	Meta               Meta               `db:"f_meta,default='{}'" json:"meta,omitempty"`
	OperatorPrivateKey string             `db:"f_prvkey,default=''" json:"-"`
}
