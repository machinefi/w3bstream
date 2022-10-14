package models

import (
	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"

	"github.com/iotexproject/w3bstream/pkg/enums"
)

// Instance database model instance
// @def primary                     ID
// @def unique_index UI_instance_id InstanceID
// @def index        I_applet_id    AppletID
// @def index        I_path         Path
//
//go:generate toolkit gen model Instance --database DB
type Instance struct {
	datatypes.PrimaryID
	RelInstance
	RelApplet
	InstanceInfo
	datatypes.OperationTimesWithDeleted
}

type RelInstance struct {
	InstanceID types.SFID `db:"f_instance_id" json:"instanceID"`
}

type InstanceInfo struct {
	Path  string              `db:"f_path"  json:"-"`
	State enums.InstanceState `db:"f_state" json:"state"`
}
