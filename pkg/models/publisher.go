package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// Publisher database model demo
// @def primary                        ID
// @def unique_index UI_publisher_id   PublisherID
// @def unique_index UI_publisher_key  Key
//
//go:generate toolkit gen model Publisher --database DB
type Publisher struct {
	datatypes.PrimaryID
	RelProject
	RelPublisher
	PublisherInfo
	datatypes.OperationTimes
}

type RelPublisher struct {
	PublisherID types.SFID `db:"f_publisher_id" json:"publisherID"`
}

type PublisherInfo struct {
	Name  string `db:"f_name"             json:"name"`
	Key   string `db:"f_key"              json:"key"` // Key the unique identifier for publisher
	Token string `db:"f_token,default=''" json:"token"`
}
