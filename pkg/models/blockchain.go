package models

import "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"

// TODO will delete after caller migrate to http api

// Blockchain database model blockchain
// @def primary                   ID
// @def unique_index UI_chain_id ChainID
//
//go:generate toolkit gen model Blockchain --database MonitorDB
type Blockchain struct {
	datatypes.PrimaryID
	RelBlockchain
	BlockchainInfo
	datatypes.OperationTimes
}

type RelBlockchain struct {
	ChainID uint64 `db:"f_chain_id"             json:"chainID"`
}

type BlockchainInfo struct {
	Address string `db:"f_chain_address"         json:"chainAddress"`
}
