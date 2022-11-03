// This is a generated source file. DO NOT EDIT
// Source: enums/account_identity_type__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidAccountIdentityType = errors.New("invalid AccountIdentityType type")

func ParseAccountIdentityTypeFromString(s string) (AccountIdentityType, error) {
	switch s {
	default:
		return ACCOUNT_IDENTITY_TYPE_UNKNOWN, InvalidAccountIdentityType
	case "":
		return ACCOUNT_IDENTITY_TYPE_UNKNOWN, nil
	case "MOBILE":
		return ACCOUNT_IDENTITY_TYPE__MOBILE, nil
	case "EMAIL":
		return ACCOUNT_IDENTITY_TYPE__EMAIL, nil
	case "USERNAME":
		return ACCOUNT_IDENTITY_TYPE__USERNAME, nil
	case "BUILTIN":
		return ACCOUNT_IDENTITY_TYPE__BUILTIN, nil
	}
}

func ParseAccountIdentityTypeFromLabel(s string) (AccountIdentityType, error) {
	switch s {
	default:
		return ACCOUNT_IDENTITY_TYPE_UNKNOWN, InvalidAccountIdentityType
	case "":
		return ACCOUNT_IDENTITY_TYPE_UNKNOWN, nil
	case "MOBILE":
		return ACCOUNT_IDENTITY_TYPE__MOBILE, nil
	case "EMAIL":
		return ACCOUNT_IDENTITY_TYPE__EMAIL, nil
	case "USERNAME":
		return ACCOUNT_IDENTITY_TYPE__USERNAME, nil
	case "BUILTIN":
		return ACCOUNT_IDENTITY_TYPE__BUILTIN, nil
	}
}

func (v AccountIdentityType) Int() int {
	return int(v)
}

func (v AccountIdentityType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCOUNT_IDENTITY_TYPE_UNKNOWN:
		return ""
	case ACCOUNT_IDENTITY_TYPE__MOBILE:
		return "MOBILE"
	case ACCOUNT_IDENTITY_TYPE__EMAIL:
		return "EMAIL"
	case ACCOUNT_IDENTITY_TYPE__USERNAME:
		return "USERNAME"
	case ACCOUNT_IDENTITY_TYPE__BUILTIN:
		return "BUILTIN"
	}
}

func (v AccountIdentityType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCOUNT_IDENTITY_TYPE_UNKNOWN:
		return ""
	case ACCOUNT_IDENTITY_TYPE__MOBILE:
		return "MOBILE"
	case ACCOUNT_IDENTITY_TYPE__EMAIL:
		return "EMAIL"
	case ACCOUNT_IDENTITY_TYPE__USERNAME:
		return "USERNAME"
	case ACCOUNT_IDENTITY_TYPE__BUILTIN:
		return "BUILTIN"
	}
}

func (v AccountIdentityType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.AccountIdentityType"
}

func (v AccountIdentityType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{ACCOUNT_IDENTITY_TYPE__MOBILE, ACCOUNT_IDENTITY_TYPE__EMAIL, ACCOUNT_IDENTITY_TYPE__USERNAME, ACCOUNT_IDENTITY_TYPE__BUILTIN}
}

func (v AccountIdentityType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidAccountIdentityType
	}
	return []byte(s), nil
}

func (v *AccountIdentityType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseAccountIdentityTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *AccountIdentityType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = AccountIdentityType(i)
	return nil
}

func (v AccountIdentityType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
