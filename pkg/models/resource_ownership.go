package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// ResourceOwnership database model resource ownership
// @def primary                          ID
// @def unique_index UI_resource_account ResourceID AccountID
//
//go:generate toolkit gen model ResourceOwnership --database DB
type ResourceOwnership struct {
	datatypes.PrimaryID
	RelResource
	RelAccount
	ResourceOwnerInfo
	datatypes.OperationTimes
}
type ResourceOwnerInfo struct {
	UploadedAt types.Timestamp `db:"f_uploaded_at"           json:"uploadedAt"`
	ExpireAt   types.Timestamp `db:"f_expire_at,default='0'" json:"expireAt"`
	Filename   string          `db:"f_filename,default=''"   json:"filename"`
	Comment    string          `db:"f_comment,default=''"    json:"comment"`
}
