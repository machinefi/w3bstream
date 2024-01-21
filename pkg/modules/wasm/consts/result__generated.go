// This is a generated source file. DO NOT EDIT
// Source: internal/result__generated.go

package consts

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidResult = errors.New("invalid Result type")

func ParseResultFromString(s string) (Result, error) {
	switch s {
	default:
		return RESULT_UNKNOWN, InvalidResult
	case "":
		return RESULT_UNKNOWN, nil
	case "INVALID_MEM_ACCESS":
		return RESULT__INVALID_MEM_ACCESS, nil
	case "ENV_NOT_FOUND":
		return RESULT__ENV_NOT_FOUND, nil
	case "RESOURCE_NOT_FOUND":
		return RESULT__RESOURCE_NOT_FOUND, nil
	case "IMPORT_HANDLE_FAILED":
		return RESULT__IMPORT_HANDLE_FAILED, nil
	case "HOST_INVOKE_FAILED":
		return RESULT__HOST_INVOKE_FAILED, nil
	}
}

func ParseResultFromLabel(s string) (Result, error) {
	switch s {
	default:
		return RESULT_UNKNOWN, InvalidResult
	case "":
		return RESULT_UNKNOWN, nil
	case "INVALID_MEM_ACCESS":
		return RESULT__INVALID_MEM_ACCESS, nil
	case "ENV_NOT_FOUND":
		return RESULT__ENV_NOT_FOUND, nil
	case "RESOURCE_NOT_FOUND":
		return RESULT__RESOURCE_NOT_FOUND, nil
	case "IMPORT_HANDLE_FAILED":
		return RESULT__IMPORT_HANDLE_FAILED, nil
	case "HOST_INVOKE_FAILED":
		return RESULT__HOST_INVOKE_FAILED, nil
	}
}

func (v Result) Int() int {
	return int(v)
}

func (v Result) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case RESULT_UNKNOWN:
		return ""
	case RESULT__INVALID_MEM_ACCESS:
		return "INVALID_MEM_ACCESS"
	case RESULT__ENV_NOT_FOUND:
		return "ENV_NOT_FOUND"
	case RESULT__RESOURCE_NOT_FOUND:
		return "RESOURCE_NOT_FOUND"
	case RESULT__IMPORT_HANDLE_FAILED:
		return "IMPORT_HANDLE_FAILED"
	case RESULT__HOST_INVOKE_FAILED:
		return "HOST_INVOKE_FAILED"
	}
}

func (v Result) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case RESULT_UNKNOWN:
		return ""
	case RESULT__INVALID_MEM_ACCESS:
		return "INVALID_MEM_ACCESS"
	case RESULT__ENV_NOT_FOUND:
		return "ENV_NOT_FOUND"
	case RESULT__RESOURCE_NOT_FOUND:
		return "RESOURCE_NOT_FOUND"
	case RESULT__IMPORT_HANDLE_FAILED:
		return "IMPORT_HANDLE_FAILED"
	case RESULT__HOST_INVOKE_FAILED:
		return "HOST_INVOKE_FAILED"
	}
}

func (v Result) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/modules/wasm/internal.Result"
}

func (v Result) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{RESULT__INVALID_MEM_ACCESS, RESULT__ENV_NOT_FOUND, RESULT__RESOURCE_NOT_FOUND, RESULT__IMPORT_HANDLE_FAILED, RESULT__HOST_INVOKE_FAILED}
}

func (v Result) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidResult
	}
	return []byte(s), nil
}

func (v *Result) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseResultFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *Result) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = Result(i)
	return nil
}

func (v Result) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
