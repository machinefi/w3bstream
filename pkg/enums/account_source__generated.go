// This is a generated source file. DO NOT EDIT
// Source: enums/account_source__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidAccountSource = errors.New("invalid AccountSource type")

func ParseAccountSourceFromString(s string) (AccountSource, error) {
	switch s {
	default:
		return ACCOUNT_SOURCE_UNKNOWN, InvalidAccountSource
	case "":
		return ACCOUNT_SOURCE_UNKNOWN, nil
	case "INIT":
		return ACCOUNT_SOURCE__INIT, nil
	case "SUBMIT":
		return ACCOUNT_SOURCE__SUBMIT, nil
	}
}

func ParseAccountSourceFromLabel(s string) (AccountSource, error) {
	switch s {
	default:
		return ACCOUNT_SOURCE_UNKNOWN, InvalidAccountSource
	case "":
		return ACCOUNT_SOURCE_UNKNOWN, nil
	case "INIT":
		return ACCOUNT_SOURCE__INIT, nil
	case "SUBMIT":
		return ACCOUNT_SOURCE__SUBMIT, nil
	}
}

func (v AccountSource) Int() int {
	return int(v)
}

func (v AccountSource) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCOUNT_SOURCE_UNKNOWN:
		return ""
	case ACCOUNT_SOURCE__INIT:
		return "INIT"
	case ACCOUNT_SOURCE__SUBMIT:
		return "SUBMIT"
	}
}

func (v AccountSource) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCOUNT_SOURCE_UNKNOWN:
		return ""
	case ACCOUNT_SOURCE__INIT:
		return "INIT"
	case ACCOUNT_SOURCE__SUBMIT:
		return "SUBMIT"
	}
}

func (v AccountSource) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.AccountSource"
}

func (v AccountSource) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{ACCOUNT_SOURCE__INIT, ACCOUNT_SOURCE__SUBMIT}
}

func (v AccountSource) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidAccountSource
	}
	return []byte(s), nil
}

func (v *AccountSource) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseAccountSourceFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *AccountSource) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = AccountSource(i)
	return nil
}

func (v AccountSource) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
