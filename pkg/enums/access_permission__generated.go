// This is a generated source file. DO NOT EDIT
// Source: enums/access_permission__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidAccessPermission = errors.New("invalid AccessPermission type")

func ParseAccessPermissionFromString(s string) (AccessPermission, error) {
	switch s {
	default:
		return ACCESS_PERMISSION_UNKNOWN, InvalidAccessPermission
	case "":
		return ACCESS_PERMISSION_UNKNOWN, nil
	case "NO_ACCESS":
		return ACCESS_PERMISSION__NO_ACCESS, nil
	case "READONLY":
		return ACCESS_PERMISSION__READONLY, nil
	case "READ_WRITE":
		return ACCESS_PERMISSION__READ_WRITE, nil
	}
}

func ParseAccessPermissionFromLabel(s string) (AccessPermission, error) {
	switch s {
	default:
		return ACCESS_PERMISSION_UNKNOWN, InvalidAccessPermission
	case "":
		return ACCESS_PERMISSION_UNKNOWN, nil
	case "NO_ACCESS":
		return ACCESS_PERMISSION__NO_ACCESS, nil
	case "READONLY":
		return ACCESS_PERMISSION__READONLY, nil
	case "READ_WRITE":
		return ACCESS_PERMISSION__READ_WRITE, nil
	}
}

func (v AccessPermission) Int() int {
	return int(v)
}

func (v AccessPermission) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCESS_PERMISSION_UNKNOWN:
		return ""
	case ACCESS_PERMISSION__NO_ACCESS:
		return "NO_ACCESS"
	case ACCESS_PERMISSION__READONLY:
		return "READONLY"
	case ACCESS_PERMISSION__READ_WRITE:
		return "READ_WRITE"
	}
}

func (v AccessPermission) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case ACCESS_PERMISSION_UNKNOWN:
		return ""
	case ACCESS_PERMISSION__NO_ACCESS:
		return "NO_ACCESS"
	case ACCESS_PERMISSION__READONLY:
		return "READONLY"
	case ACCESS_PERMISSION__READ_WRITE:
		return "READ_WRITE"
	}
}

func (v AccessPermission) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.AccessPermission"
}

func (v AccessPermission) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{ACCESS_PERMISSION__NO_ACCESS, ACCESS_PERMISSION__READONLY, ACCESS_PERMISSION__READ_WRITE}
}

func (v AccessPermission) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidAccessPermission
	}
	return []byte(s), nil
}

func (v *AccessPermission) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseAccessPermissionFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *AccessPermission) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = AccessPermission(i)
	return nil
}

func (v AccessPermission) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
