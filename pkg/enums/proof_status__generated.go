// This is a generated source file. DO NOT EDIT
// Source: enums/proof_status__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidProofStatus = errors.New("invalid ProofStatus type")

func ParseProofStatusFromString(s string) (ProofStatus, error) {
	switch s {
	default:
		return PROOF_STATUS_UNKNOWN, InvalidProofStatus
	case "":
		return PROOF_STATUS_UNKNOWN, nil
	case "GENERATING":
		return PROOF_STATUS__GENERATING, nil
	case "SUCCEEDED":
		return PROOF_STATUS__SUCCEEDED, nil
	case "FAILED":
		return PROOF_STATUS__FAILED, nil
	}
}

func ParseProofStatusFromLabel(s string) (ProofStatus, error) {
	switch s {
	default:
		return PROOF_STATUS_UNKNOWN, InvalidProofStatus
	case "":
		return PROOF_STATUS_UNKNOWN, nil
	case "GENERATING":
		return PROOF_STATUS__GENERATING, nil
	case "SUCCEEDED":
		return PROOF_STATUS__SUCCEEDED, nil
	case "FAILED":
		return PROOF_STATUS__FAILED, nil
	}
}

func (v ProofStatus) Int() int {
	return int(v)
}

func (v ProofStatus) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case PROOF_STATUS_UNKNOWN:
		return ""
	case PROOF_STATUS__GENERATING:
		return "GENERATING"
	case PROOF_STATUS__SUCCEEDED:
		return "SUCCEEDED"
	case PROOF_STATUS__FAILED:
		return "FAILED"
	}
}

func (v ProofStatus) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case PROOF_STATUS_UNKNOWN:
		return ""
	case PROOF_STATUS__GENERATING:
		return "GENERATING"
	case PROOF_STATUS__SUCCEEDED:
		return "SUCCEEDED"
	case PROOF_STATUS__FAILED:
		return "FAILED"
	}
}

func (v ProofStatus) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.ProofStatus"
}

func (v ProofStatus) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{PROOF_STATUS__GENERATING, PROOF_STATUS__SUCCEEDED, PROOF_STATUS__FAILED}
}

func (v ProofStatus) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidProofStatus
	}
	return []byte(s), nil
}

func (v *ProofStatus) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseProofStatusFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *ProofStatus) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = ProofStatus(i)
	return nil
}

func (v ProofStatus) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
