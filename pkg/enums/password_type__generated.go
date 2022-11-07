// This is a generated source file. DO NOT EDIT
// Source: enums/password_type__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidPasswordType = errors.New("invalid PasswordType type")

func ParsePasswordTypeFromString(s string) (PasswordType, error) {
	switch s {
	default:
		return PASSWORD_TYPE_UNKNOWN, InvalidPasswordType
	case "":
		return PASSWORD_TYPE_UNKNOWN, nil
	case "LOGIN":
		return PASSWORD_TYPE__LOGIN, nil
	case "PERSONAL_TOKEN":
		return PASSWORD_TYPE__PERSONAL_TOKEN, nil
	}
}

func ParsePasswordTypeFromLabel(s string) (PasswordType, error) {
	switch s {
	default:
		return PASSWORD_TYPE_UNKNOWN, InvalidPasswordType
	case "":
		return PASSWORD_TYPE_UNKNOWN, nil
	case "LOGIN":
		return PASSWORD_TYPE__LOGIN, nil
	case "PERSONAL_TOKEN":
		return PASSWORD_TYPE__PERSONAL_TOKEN, nil
	}
}

func (v PasswordType) Int() int {
	return int(v)
}

func (v PasswordType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case PASSWORD_TYPE_UNKNOWN:
		return ""
	case PASSWORD_TYPE__LOGIN:
		return "LOGIN"
	case PASSWORD_TYPE__PERSONAL_TOKEN:
		return "PERSONAL_TOKEN"
	}
}

func (v PasswordType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case PASSWORD_TYPE_UNKNOWN:
		return ""
	case PASSWORD_TYPE__LOGIN:
		return "LOGIN"
	case PASSWORD_TYPE__PERSONAL_TOKEN:
		return "PERSONAL_TOKEN"
	}
}

func (v PasswordType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.PasswordType"
}

func (v PasswordType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{PASSWORD_TYPE__LOGIN, PASSWORD_TYPE__PERSONAL_TOKEN}
}

func (v PasswordType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidPasswordType
	}
	return []byte(s), nil
}

func (v *PasswordType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParsePasswordTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *PasswordType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = PasswordType(i)
	return nil
}

func (v PasswordType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
