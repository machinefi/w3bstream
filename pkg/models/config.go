package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Config database model config for configuration management
// @def primary                          ID
// @def unique_index UI_config_id        ConfigID
// @def unique_index UI_rel_type         RelID Type
//
//go:generate toolkit gen model Config --database DB
type Config struct {
	datatypes.PrimaryID
	RelConfig
	ConfigBase
	datatypes.OperationTimes
}

type RelConfig struct {
	ConfigID types.SFID `db:"f_config_id" json:"configID"`
}

type ConfigBase struct {
	RelID types.SFID       `db:"f_rel_id"           json:"relID"`
	Type  enums.ConfigType `db:"f_type"             json:"type"`
	Value []byte           `db:"f_value,default=''" json:"-"`
}
