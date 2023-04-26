package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// ResourceMeta database model wasm_resource_meta
// @def primary                        ID
// @def unique_index UI_meta_id        MetaID
// @def unique_index UI_res_acc_name   ResourceID AccountID ResName
//
//go:generate toolkit gen model ResourceMeta --database DB
type ResourceMeta struct {
	datatypes.PrimaryID
	RelMeta
	RelResource
	RelAccount
	MetaInfo
	datatypes.OperationTimes
}

type RelMeta struct {
	MetaID types.SFID `db:"f_meta_id" json:"metaID"`
}

type MetaInfo struct {
	ResName string `db:"f_res_name"              json:"resName"`
	RefCnt  int    `db:"f_ref_cnt,default=0"     json:"refCnt"`
}
