package models

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	_ "github.com/machinefi/w3bstream/pkg/depends/util/strfmt"
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
	Name string `db:"f_name" json:"name" validate:"@projectName"` // Name project name
}

type ProjectBase struct {
	Public      datatypes.Bool `db:"f_public,default='2'" json:"public,omitempty"`   // Public is true, project receive event from anonymous publisher
	Version     string         `db:"f_version,default=''" json:"version,omitempty"`  // Version project version
	Proto       enums.Protocol `db:"f_proto,default='0'"  json:"protocol,omitempty"` // Proto project protocol for event publisher
	Description string         `db:"f_description,default=''"    json:"description,omitempty"`
}

func (m Project) DatabaseName() string {
	return "w3b_" + m.ProjectID.String()
}

func (m Project) Privileges() (usename, passwd string) {
	usename = m.DatabaseName()
	passwd = hex.EncodeToString(sha256.New().Sum([]byte(m.ProjectID.String())))
	return
}
