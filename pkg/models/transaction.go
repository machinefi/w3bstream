package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Transaction schema for blockchain transaction information
// @def primary                      ID
// @def unique_index UI_transaction_id  TransactionID
//
//go:generate toolkit gen model Transaction --database DB
type Transaction struct {
	datatypes.PrimaryID
	RelTransaction
	RelProject
	TransactionInfo
	datatypes.OperationTimesWithDeleted
}

type RelTransaction struct {
	TransactionID types.SFID `db:"f_transaction_id" json:"transactionID"`
}

type TransactionInfo struct {
	ChainName    enums.ChainName        `db:"f_chain_name"                    json:"chainName"`
	Nonce        uint64                 `db:"f_nonce,default='0'"             json:"nonce,omitempty"`
	Hash         string                 `db:"f_hash,default=''"               json:"hash,omitempty"`
	Sender       string                 `db:"f_sender,default=''"             json:"sender,omitempty"`
	Receiver     string                 `db:"f_receiver,default=''"           json:"receiver,omitempty"`
	Value        string                 `db:"f_value,default=''"              json:"value,omitempty"`
	Data         string                 `db:"f_data,default=''"               json:"data,omitempty"`
	OperatorName string                 `db:"f_operator_name,default=''"      json:"operatorName,omitempty"`
	State        enums.TransactionState `db:"f_state,default='0'"             json:"state,omitempty"`
	EventType    string                 `db:"f_event_type"                    json:"eventType"`
}
