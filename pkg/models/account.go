package models

import (
	"database/sql/driver"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Account w3bstream account
// @def primary                  AccountID DeletedAt
// @def unique_index ui_username Username
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
	Username     string                    `db:"f_username"                  json:"username"`
	IdentityType enums.AccountIdentityType `db:"f_identity_type,default='0'" json:"identityType"`
	State        enums.AccountState        `db:"f_state,default='0'"         json:"-"`
	Password     AccountPassword           `db:"f_password"                  json:"-"`
	Vendor       AccountVendor             `db:"f_vendor,default=''"         json:"-"`
	Meta         Meta                      `db:"f_meta,default=''"           json:"meta"`
}

type AccountPassword struct {
	Type     enums.PasswordType `json:"type,omitempty"`
	Password string             `json:"password"`
	Scope    string             `json:"scope,omitempty"`
	Desc     string             `json:"desc,omitempty"`
}

func (AccountPassword) DataType(drv string) string { return "text" }

func (v AccountPassword) Value() (driver.Value, error) {
	return datatypes.JSONValue(v)
}

func (v *AccountPassword) Scan(src interface{}) error {
	return datatypes.JSONScan(src, v)
}

// AccountVendor third part vendor
type AccountVendor struct {
	From     string `json:"from"`
	Identity string `json:"identity"`
}

func (AccountVendor) DataType(drv string) string { return "text" }

func (v AccountVendor) Value() (driver.Value, error) {
	return datatypes.JSONValue(v)
}

func (v *AccountVendor) Scan(src interface{}) error {
	return datatypes.JSONScan(src, v)
}
