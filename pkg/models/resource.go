package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// Resource database model wasm_resource
// @def primary                            ID
// @def unique_index UI_resource_id        ResourceID
// @def unique_index UI_md5                Md5
//
//go:generate toolkit gen model Resource --database DB
type Resource struct {
	datatypes.PrimaryID
	RelResource
	ResourceInfo
	datatypes.OperationTimes
}

type RelResource struct {
	ResourceID types.SFID `db:"f_resource_id" json:"resourceID"`
}

type ResourceInfo struct {
	Path string `db:"f_path,default=''" json:"-"`
	Md5  string `db:"f_md5"  json:"-"`
}
