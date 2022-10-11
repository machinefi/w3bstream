package models

import (
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"
)

// Blockchain database model blockchain
// @def primary                   ID
//
//go:generate toolkit gen model Blockchain --database DB
type Blockchain struct {
	datatypes.PrimaryID
	RelBlockchain
	BlockchainInfo
	datatypes.OperationTimes
}

type RelBlockchain struct {
	BlockchainID string `db:"f_blockchain_id" json:"blockchainID"`
}

type BlockchainInfo struct {
	BlockchainAddress string `db:"f_blockchainAddress"            json:"blockchainAddress"`
	ContractAddress   string `db:"f_contractAddress"              json:"contractAddress"`
	BlockStart        uint64 `db:"f_blockStart"                   json:"blockStart"`
	BlockCurrent      uint64 `db:"f_blockCurrent"                 json:"blockCurrent"`
	ProjectID         string `db:"f_project_id"                   json:"projectID"`
	AppletID          string `db:"f_applet_id"                    json:"appletID"`
	Handler           string `db:"f_handler"                      json:"handler"`
}
