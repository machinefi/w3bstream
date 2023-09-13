package models

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Proof database model proof
// @def primary                          ID
// @def unique_index UI_proof_id        ProofID
//
//go:generate toolkit gen model Proof --database DB
type Proof struct {
	datatypes.PrimaryID
	RelProject
	RelProof
	ProofInfo
	datatypes.OperationTimes
}

type RelProof struct {
	ProofID types.SFID `db:"f_proof_id" json:"proofID"`
}

type ProofInfo struct {
	Name         string            `db:"f_name" json:"name"`
	TemplateName string            `db:"f_template_name" json:"templateName"`
	ImageID      string            `db:"f_image_id" json:"imageID"`
	InputData    string            `db:"f_input_data" json:"inputData"`
	Receipt      string            `db:"f_receipt,default='',size=102400" json:"receipt"`
	Status       enums.ProofStatus `db:"f_status" json:"status"`
}
