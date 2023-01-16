// This is a generated source file. DO NOT EDIT
// Source: definition_scanner/enum__generated.go

package definition_scanner

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidEnum = errors.New("invalid Enum type")

func ParseEnumFromString(s string) (Enum, error) {
	switch s {
	default:
		return ENUM_UNKNOWN, InvalidEnum
	case "":
		return ENUM_UNKNOWN, nil
	case "ONE":
		return ENUM__ONE, nil
	case "TWO":
		return ENUM__TWO, nil
	}
}

func ParseEnumFromLabel(s string) (Enum, error) {
	switch s {
	default:
		return ENUM_UNKNOWN, InvalidEnum
	case "":
		return ENUM_UNKNOWN, nil
	case "one":
		return ENUM__ONE, nil
	case "two":
		return ENUM__TWO, nil
	}
}

func (v Enum) Int() int {
	return int(v)
}

func (v Enum) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case ENUM_UNKNOWN:
		return ""
	case ENUM__ONE:
		return "ONE"
	case ENUM__TWO:
		return "TWO"
	}
}

func (v Enum) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case ENUM_UNKNOWN:
		return ""
	case ENUM__ONE:
		return "one"
	case ENUM__TWO:
		return "two"
	}
}

func (v Enum) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/kit/httpswaggergen/testdata/definition_scanner.Enum"
}

func (v Enum) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{ENUM__ONE, ENUM__TWO}
}

func (v Enum) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidEnum
	}
	return []byte(s), nil
}

func (v *Enum) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseEnumFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *Enum) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = Enum(i)
	return nil
}

func (v Enum) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
