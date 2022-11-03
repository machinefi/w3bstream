// This is a generated source file. DO NOT EDIT
// Source: schema/constrain_type__generated.go

package schema

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidConstrainType = errors.New("invalid ConstrainType type")

func ParseConstrainTypeFromString(s string) (ConstrainType, error) {
	switch s {
	default:
		return CONSTRAIN_TYPE_UNKNOWN, InvalidConstrainType
	case "":
		return CONSTRAIN_TYPE_UNKNOWN, nil
	case "AUTOINCREMENT":
		return CONSTRAIN_TYPE__AUTOINCREMENT, nil
	case "NOT_NULL":
		return CONSTRAIN_TYPE__NOT_NULL, nil
	case "DEFAULT":
		return CONSTRAIN_TYPE__DEFAULT, nil
	}
}

func ParseConstrainTypeFromLabel(s string) (ConstrainType, error) {
	switch s {
	default:
		return CONSTRAIN_TYPE_UNKNOWN, InvalidConstrainType
	case "":
		return CONSTRAIN_TYPE_UNKNOWN, nil
	case "AUTOINCREMENT":
		return CONSTRAIN_TYPE__AUTOINCREMENT, nil
	case "NOT_NULL":
		return CONSTRAIN_TYPE__NOT_NULL, nil
	case "DEFAULT":
		return CONSTRAIN_TYPE__DEFAULT, nil
	}
}

func (v ConstrainType) Int() int {
	return int(v)
}

func (v ConstrainType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case CONSTRAIN_TYPE_UNKNOWN:
		return ""
	case CONSTRAIN_TYPE__AUTOINCREMENT:
		return "AUTOINCREMENT"
	case CONSTRAIN_TYPE__NOT_NULL:
		return "NOT_NULL"
	case CONSTRAIN_TYPE__DEFAULT:
		return "DEFAULT"
	}
}

func (v ConstrainType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case CONSTRAIN_TYPE_UNKNOWN:
		return ""
	case CONSTRAIN_TYPE__AUTOINCREMENT:
		return "AUTOINCREMENT"
	case CONSTRAIN_TYPE__NOT_NULL:
		return "NOT_NULL"
	case CONSTRAIN_TYPE__DEFAULT:
		return "DEFAULT"
	}
}

func (v ConstrainType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/schema.ConstrainType"
}

func (v ConstrainType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{CONSTRAIN_TYPE__AUTOINCREMENT, CONSTRAIN_TYPE__NOT_NULL, CONSTRAIN_TYPE__DEFAULT}
}

func (v ConstrainType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidConstrainType
	}
	return []byte(s), nil
}

func (v *ConstrainType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseConstrainTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *ConstrainType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = ConstrainType(i)
	return nil
}

func (v ConstrainType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
