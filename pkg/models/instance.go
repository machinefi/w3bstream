package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Instance database model instance
// @def primary                     ID
// @def unique_index UI_instance_id InstanceID
// @def index        I_applet_id    AppletID
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
	State enums.InstanceState `db:"f_state" json:"state"`
}
