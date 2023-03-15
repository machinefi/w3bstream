// This is a generated source file. DO NOT EDIT
// Source: enums/account_role__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidAccountRole = errors.New("invalid AccountRole type")

func ParseAccountRoleFromString(s string) (AccountRole, error) {
	switch s {
	default:
		return ACCOUNT_ROLE_UNKNOWN, InvalidAccountRole
	case "":
		return ACCOUNT_ROLE_UNKNOWN, nil
	case "ADMIN":
		return ACCOUNT_ROLE__ADMIN, nil
	case "DEVELOPER":
		return ACCOUNT_ROLE__DEVELOPER, nil
	}
}

func ParseAccountRoleFromLabel(s string) (AccountRole, error) {
	switch s {
	default:
		return ACCOUNT_ROLE_UNKNOWN, InvalidAccountRole
	case "":
		return ACCOUNT_ROLE_UNKNOWN, nil
	case "ADMIN":
		return ACCOUNT_ROLE__ADMIN, nil
	case "DEVELOPER":
		return ACCOUNT_ROLE__DEVELOPER, nil
	}
}

func (v AccountRole) Int() int {
	return int(v)
}

func (v AccountRole) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCOUNT_ROLE_UNKNOWN:
		return ""
	case ACCOUNT_ROLE__ADMIN:
		return "ADMIN"
	case ACCOUNT_ROLE__DEVELOPER:
		return "DEVELOPER"
	}
}

func (v AccountRole) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCOUNT_ROLE_UNKNOWN:
		return ""
	case ACCOUNT_ROLE__ADMIN:
		return "ADMIN"
	case ACCOUNT_ROLE__DEVELOPER:
		return "DEVELOPER"
	}
}

func (v AccountRole) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.AccountRole"
}

func (v AccountRole) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{ACCOUNT_ROLE__ADMIN, ACCOUNT_ROLE__DEVELOPER}
}

func (v AccountRole) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidAccountRole
	}
	return []byte(s), nil
}

func (v *AccountRole) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseAccountRoleFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *AccountRole) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = AccountRole(i)
	return nil
}

func (v AccountRole) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
