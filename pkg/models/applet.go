package models

import (
	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"
)

// Applet database model applet
// @def primary                   ID
// @def unique_index UI_applet_id AppletID
// @def unique_index UI_name      Name
// @def unique_index UI_md5       Md5
//
//go:generate toolkit gen model Applet --database DB
type Applet struct {
	datatypes.PrimaryID
	RelProject
	RelApplet
	AppletInfo
	datatypes.OperationTimes
}

type RelApplet struct {
	AppletID types.SFID `db:"f_applet_id" json:"appletID"`
}

type AppletInfo struct {
	Name string `db:"f_name" json:"name"`
	Path string `db:"f_path" json:"-"`
	Md5  string `db:"f_md5"  json:"-"`
}
