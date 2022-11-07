// This is a generated source file. DO NOT EDIT
// Source: example/sample__generated.go

package example

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidSample = errors.New("invalid Sample type")

func ParseSampleFromString(s string) (Sample, error) {
	switch s {
	default:
		return SAMPLE_UNKNOWN, InvalidSample
	case "":
		return SAMPLE_UNKNOWN, nil
	case "XXX":
		return SAMPLE__XXX, nil
	case "YYY":
		return SAMPLE__YYY, nil
	case "ZZZ":
		return SAMPLE__ZZZ, nil
	}
}

func ParseSampleFromLabel(s string) (Sample, error) {
	switch s {
	default:
		return SAMPLE_UNKNOWN, InvalidSample
	case "":
		return SAMPLE_UNKNOWN, nil
	case "样例XXX":
		return SAMPLE__XXX, nil
	case "样例YYY":
		return SAMPLE__YYY, nil
	case "样例ZZZ":
		return SAMPLE__ZZZ, nil
	}
}

func (v Sample) Int() int {
	return int(v)
}

func (v Sample) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case SAMPLE_UNKNOWN:
		return ""
	case SAMPLE__XXX:
		return "XXX"
	case SAMPLE__YYY:
		return "YYY"
	case SAMPLE__ZZZ:
		return "ZZZ"
	}
}

func (v Sample) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case SAMPLE_UNKNOWN:
		return ""
	case SAMPLE__XXX:
		return "样例XXX"
	case SAMPLE__YYY:
		return "样例YYY"
	case SAMPLE__ZZZ:
		return "样例ZZZ"
	}
}

func (v Sample) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/kit/enumgen/__examples__.Sample"
}

func (v Sample) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{SAMPLE__XXX, SAMPLE__YYY, SAMPLE__ZZZ}
}

func (v Sample) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidSample
	}
	return []byte(s), nil
}

func (v *Sample) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseSampleFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *Sample) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = Sample(i)
	return nil
}

func (v Sample) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
