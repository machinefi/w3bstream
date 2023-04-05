package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Tag tag for other object
// @def primary                            ID
// @def unique_index UI_tag_id             TagID
//
//go:generate toolkit gen model Tag --database DB
type Tag struct {
	datatypes.PrimaryID
	RelProject
	RelTag
	TagInfo
	datatypes.OperationTimesWithDeleted
}

type RelTag struct {
	TagID types.SFID `db:"f_tag_id" json:"tagID"`
}

type TagInfo struct {
	ReferenceID   types.SFID             `db:"f_reference_id" json:"referenceID"`
	ReferenceType enums.TagReferenceType `db:"f_reference_type" json:"referenceType"`
	Info          string                 `db:"f_info" json:"info"`
}
