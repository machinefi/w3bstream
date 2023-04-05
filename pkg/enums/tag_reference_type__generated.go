// This is a generated source file. DO NOT EDIT
// Source: enums/tag_reference_type__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidTagReferenceType = errors.New("invalid TagReferenceType type")

func ParseTagReferenceTypeFromString(s string) (TagReferenceType, error) {
	switch s {
	default:
		return TAG_REFERENCE_TYPE_UNKNOWN, InvalidTagReferenceType
	case "":
		return TAG_REFERENCE_TYPE_UNKNOWN, nil
	case "PROJECT":
		return TAG_REFERENCE_TYPE__PROJECT, nil
	}
}

func ParseTagReferenceTypeFromLabel(s string) (TagReferenceType, error) {
	switch s {
	default:
		return TAG_REFERENCE_TYPE_UNKNOWN, InvalidTagReferenceType
	case "":
		return TAG_REFERENCE_TYPE_UNKNOWN, nil
	case "PROJECT":
		return TAG_REFERENCE_TYPE__PROJECT, nil
	}
}

func (v TagReferenceType) Int() int {
	return int(v)
}

func (v TagReferenceType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case TAG_REFERENCE_TYPE_UNKNOWN:
		return ""
	case TAG_REFERENCE_TYPE__PROJECT:
		return "PROJECT"
	}
}

func (v TagReferenceType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case TAG_REFERENCE_TYPE_UNKNOWN:
		return ""
	case TAG_REFERENCE_TYPE__PROJECT:
		return "PROJECT"
	}
}

func (v TagReferenceType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.TagReferenceType"
}

func (v TagReferenceType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{TAG_REFERENCE_TYPE__PROJECT}
}

func (v TagReferenceType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidTagReferenceType
	}
	return []byte(s), nil
}

func (v *TagReferenceType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseTagReferenceTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *TagReferenceType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = TagReferenceType(i)
	return nil
}

func (v TagReferenceType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
