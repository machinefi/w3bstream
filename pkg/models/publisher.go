package models

import (
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"
)

// Publisher database model demo
// @def primary                        ID
// @def unique_index UI_publisher_id   PublisherID
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
	PublisherID string `db:"f_publisher_id" json:"publisherID"`
}

type PublisherInfo struct {
	Name string `db:"f_name"      json:"name"`
}
