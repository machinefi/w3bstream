package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// AccessKey api access key
// @def primary ID
// @def unique_index ui_name AccountID Name
// @def unique_index ui_rand Rand
//
//go:generate toolkit gen model AccessKey --database DB
type AccessKey struct {
	datatypes.PrimaryID
	RelAccount
	AccessKeyInfo
	datatypes.OperationTimesWithDeleted
}

type AccessKeyInfo struct {
	IdentityID   types.SFID                  `db:"f_identity_id"`
	IdentityType enums.AccessKeyIdentityType `db:"f_identity_type"`
	Name         string                      `db:"f_name"`
	Rand         string                      `db:"f_rand"`
	ExpiredAt    types.Timestamp             `db:"f_expired_at,default='0'"`
	LastUsed     types.Timestamp             `db:"f_last_used,default='0'"`
	Description  string                      `db:"f_desc,default=''"`
	_Privilege   interface{}                 `db:"-"` // TODO add privilege for account api key
}
