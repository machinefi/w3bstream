package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Project schema for project information
// @def primary                    ID
// @def unique_index UI_project_id ProjectID
// @def unique_index UI_name       Name
//
//go:generate toolkit gen model Project --database DB
type Project struct {
	datatypes.PrimaryID
	RelProject
	RelAccount
	ProjectName
	ProjectBase
	datatypes.OperationTimesWithDeleted
}

type RelProject struct {
	ProjectID types.SFID `db:"f_project_id" json:"projectID"`
}

type ProjectName struct {
	Name string `db:"f_name" json:"name"` // Name project name
}

type ProjectBase struct {
	Version     string         `db:"f_version,default=''" json:"version,omitempty"`  // Version project version
	Proto       enums.Protocol `db:"f_proto,default='0'"  json:"protocol,omitempty"` // Proto project protocol for event publisher
	Description string         `db:"f_description,default=''"    json:"description,omitempty"`
	Issuer      string         `db:"f_issuer,default='web3stream'" json:"issuer,omitempty"`
	ExpIn       types.Duration `db:"f_exp_in,default='0'" json:"expIn,omitempty"`
	SignKey     string         `db:"f_sign_key,default='web3stream'" json:"signKey,omitempty"`
}
