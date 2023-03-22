package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Account w3bstream account
// @def primary                  AccountID DeletedAt
//
//go:generate toolkit gen model Account --database DB
type Account struct {
	RelAccount
	AccountInfo
	datatypes.OperationTimesWithDeleted
}

type RelAccount struct {
	// @rel Account.AccountID
	AccountID types.SFID `db:"f_account_id" json:"accountID"` // AccountID  account id
}

type AccountInfo struct {
	Role               enums.AccountRole  `db:"f_role"              json:"role"`
	State              enums.AccountState `db:"f_state,default='1'" json:"state"`
	Avatar             string             `db:"f_avatar,default=''" json:"avatar,omitempty"`
	Meta               Meta               `db:"f_meta,default='{}'" json:"meta,omitempty"`
	OperatorPrivateKey string             `db:"f_prvkey"            json:"-"`
}
