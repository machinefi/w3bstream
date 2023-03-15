package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// AccountPassword account password
// @def primary                          PasswordID
// @def unique_index ui_account_password AccountID Type DeletedAt
//
//go:generate toolkit gen model AccountPassword --database DB
type AccountPassword struct {
	RelAccount
	RelAccountPassword
	AccountPasswordData
	datatypes.OperationTimesWithDeleted
}

type RelAccountPassword struct {
	PasswordID types.SFID `db:"f_password_id" json:"passwordID"`
}

type AccountPasswordData struct {
	Type     enums.PasswordType `db:"f_type"              json:"type,omitempty"`   // Type password type
	Password string             `db:"f_password,size=32"  json:"password"`         // Password md5(md5(${account_id}-${password}))
	Scope    string             `db:"f_scope,default=''"  json:"scope,omitempty"`  // Scope comma separated
	Remark   string             `db:"f_remark,default=''" json:"remark,omitempty"` // Remark
}
