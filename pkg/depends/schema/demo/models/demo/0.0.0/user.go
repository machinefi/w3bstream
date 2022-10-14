// This is a generated source file. DO NOT EDIT
// Source: demo/user.go

package demo

import (
	"github.com/iotexproject/Bumblebee/base/types"
)

// User User Table
// @def primary                  ID
// @def unique_index ui_username Username
type User struct {
	ID        int64           `db:"f_id,autoincrement"`   // user id
	Username  string          `db:"f_username"`           // username
	Gender    uint8           `db:"f_gender,default='1'"` // gender 1 male 2 female
	CreatedAt types.Timestamp `db:"f_created_at,default='0'"`
	UpdatedAt types.Timestamp `db:"f_updated_at,default='0'"`
	DeletedAt types.Timestamp `db:"f_deleted_at,default='0'"`
}
