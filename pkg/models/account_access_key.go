package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// AccountAPIKey account api access key
// @def primary ID
// @def unique_index ui_account_key_name AccountID Name
// @def unique_index ui_access_key       AccessKey
//
//go:generate toolkit gen model AccountAccessKey --database DB
type AccountAccessKey struct {
	datatypes.PrimaryID
	RelAccount
	AccountAccessKeyInfo
	datatypes.OperationTimesWithDeleted
}

type AccountAccessKeyInfo struct {
	Name        string          `db:"f_name"`
	AccessKey   string          `db:"f_access_key"`
	ExpiredAt   types.Timestamp `db:"f_expired_at,default='0'"`
	Description string          `db:"f_desc,default=''"`
	_Privilege  interface{}     `db:"-"` // TODO add privilege for account api key
}
