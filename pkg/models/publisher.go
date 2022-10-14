package models

import (
	"github.com/iotexproject/Bumblebee/base/types"
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
	PublisherID types.SFID `db:"f_publisher_id" json:"publisherID"`
}

type PublisherInfo struct {
	Name  string `db:"f_name"             json:"name"`
	Key   string `db:"f_key"              json:"key"`
	Token string `db:"f_token,default=''" json:"token"`
}
