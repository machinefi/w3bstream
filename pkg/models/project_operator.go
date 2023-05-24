package models

import "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"

// ProjectOperator schema for project operator relationship
// @def primary                    ID
// @def unique_index UI_project_id ProjectID
//
//go:generate toolkit gen model ProjectOperator --database DB
type ProjectOperator struct {
	datatypes.PrimaryID
	RelProject
	RelOperator
	datatypes.OperationTimesWithDeleted
}
