package models

import (
	"database/sql/driver"

	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"
	"github.com/iotexproject/w3bstream/pkg/enums"
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
	Protocol enums.Protocol `db:"f_protocol"  json:"protocol"`
	Data     PublisherData  `db:"f_data"      json:"data"`
}

type PublisherData struct {
	MQTT *PublisherMQTT `json:"mqtt"`
}

type PublisherMQTT struct {
	Broker   string `json:"broker"`
	Port     uint64 `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topic    string `json:"topic"`
}

func (PublisherData) DataType(engine string) string {
	return "TEXT"
}

func (p PublisherData) Value() (driver.Value, error) {
	return datatypes.JSONValue(p)
}

func (p *PublisherData) Scan(src interface{}) error {
	return datatypes.JSONScan(src, p)
}
