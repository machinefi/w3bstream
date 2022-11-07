// This is a generated source file. DO NOT EDIT
// Source: enums/account_state__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidAccountState = errors.New("invalid AccountState type")

func ParseAccountStateFromString(s string) (AccountState, error) {
	switch s {
	default:
		return ACCOUNT_STATE_UNKNOWN, InvalidAccountState
	case "":
		return ACCOUNT_STATE_UNKNOWN, nil
	case "ENABLED":
		return ACCOUNT_STATE__ENABLED, nil
	case "DISABLED":
		return ACCOUNT_STATE__DISABLED, nil
	}
}

func ParseAccountStateFromLabel(s string) (AccountState, error) {
	switch s {
	default:
		return ACCOUNT_STATE_UNKNOWN, InvalidAccountState
	case "":
		return ACCOUNT_STATE_UNKNOWN, nil
	case "ENABLED":
		return ACCOUNT_STATE__ENABLED, nil
	case "DISABLED":
		return ACCOUNT_STATE__DISABLED, nil
	}
}

func (v AccountState) Int() int {
	return int(v)
}

func (v AccountState) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCOUNT_STATE_UNKNOWN:
		return ""
	case ACCOUNT_STATE__ENABLED:
		return "ENABLED"
	case ACCOUNT_STATE__DISABLED:
		return "DISABLED"
	}
}

func (v AccountState) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCOUNT_STATE_UNKNOWN:
		return ""
	case ACCOUNT_STATE__ENABLED:
		return "ENABLED"
	case ACCOUNT_STATE__DISABLED:
		return "DISABLED"
	}
}

func (v AccountState) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.AccountState"
}

func (v AccountState) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{ACCOUNT_STATE__ENABLED, ACCOUNT_STATE__DISABLED}
}

func (v AccountState) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidAccountState
	}
	return []byte(s), nil
}

func (v *AccountState) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseAccountStateFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *AccountState) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = AccountState(i)
	return nil
}

func (v AccountState) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
