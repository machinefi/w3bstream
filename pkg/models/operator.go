package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

// Operator schema for blockchain operate information
// @def primary                      ID
// @def unique_index UI_operator_id  OperatorID
// @def unique_index UI_name         AccountID Name
//
//go:generate toolkit gen model Operator --database DB
type Operator struct {
	datatypes.PrimaryID
	RelAccount
	RelOperator
	OperatorInfo
	datatypes.OperationTimesWithDeleted
}

type RelOperator struct {
	OperatorID types.SFID `db:"f_operator_id" json:"operatorID"`
}

type OperatorInfo struct {
	PrivateKey string `db:"f_private_key" json:"-"`
	Name       string `db:"f_name" json:"name"`
}
