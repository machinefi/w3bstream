package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// ResourceMeta database model wasm_resource_meta
// @def primary                        ID
// @def unique_index UI_meta_id        MetaID
// @def unique_index UI_res_acc_app   ResourceID AccountID AppletID
//
//go:generate toolkit gen model ResourceMeta --database DB
type ResourceMeta struct {
	datatypes.PrimaryID
	RelMeta
	RelResource
	RelAccount
	RelApplet
	MetaInfo
	datatypes.OperationTimesWithDeleted
}

type RelMeta struct {
	MetaID types.SFID `db:"f_meta_id" json:"metaID"`
}

type MetaInfo struct {
	FileName   string          `db:"f_file_name"            json:"fileName"`
}
