// This is a generated source file. DO NOT EDIT
// Source: enums/flow_operator__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidFlowOperator = errors.New("invalid FlowOperator type")

func ParseFlowOperatorFromString(s string) (FlowOperator, error) {
	switch s {
	default:
		return FLOW_OPERATOR_UNKNOWN, InvalidFlowOperator
	case "":
		return FLOW_OPERATOR_UNKNOWN, nil
	case "FILTER":
		return FLOW_OPERATOR__FILTER, nil
	case "MAP":
		return FLOW_OPERATOR__MAP, nil
	case "FLATMAP":
		return FLOW_OPERATOR__FLATMAP, nil
	case "WINDOW":
		return FLOW_OPERATOR__WINDOW, nil
	case "GROUP":
		return FLOW_OPERATOR__GROUP, nil
	case "REDUCE":
		return FLOW_OPERATOR__REDUCE, nil
	}
}

func ParseFlowOperatorFromLabel(s string) (FlowOperator, error) {
	switch s {
	default:
		return FLOW_OPERATOR_UNKNOWN, InvalidFlowOperator
	case "":
		return FLOW_OPERATOR_UNKNOWN, nil
	case "FILTER":
		return FLOW_OPERATOR__FILTER, nil
	case "MAP":
		return FLOW_OPERATOR__MAP, nil
	case "FLATMAP":
		return FLOW_OPERATOR__FLATMAP, nil
	case "WINDOW":
		return FLOW_OPERATOR__WINDOW, nil
	case "GROUP":
		return FLOW_OPERATOR__GROUP, nil
	case "REDUCE":
		return FLOW_OPERATOR__REDUCE, nil
	}
}

func (v FlowOperator) Int() int {
	return int(v)
}

func (v FlowOperator) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case FLOW_OPERATOR_UNKNOWN:
		return ""
	case FLOW_OPERATOR__FILTER:
		return "FILTER"
	case FLOW_OPERATOR__MAP:
		return "MAP"
	case FLOW_OPERATOR__FLATMAP:
		return "FLATMAP"
	case FLOW_OPERATOR__WINDOW:
		return "WINDOW"
	case FLOW_OPERATOR__GROUP:
		return "GROUP"
	case FLOW_OPERATOR__REDUCE:
		return "REDUCE"
	}
}

func (v FlowOperator) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case FLOW_OPERATOR_UNKNOWN:
		return ""
	case FLOW_OPERATOR__FILTER:
		return "FILTER"
	case FLOW_OPERATOR__MAP:
		return "MAP"
	case FLOW_OPERATOR__FLATMAP:
		return "FLATMAP"
	case FLOW_OPERATOR__WINDOW:
		return "WINDOW"
	case FLOW_OPERATOR__GROUP:
		return "GROUP"
	case FLOW_OPERATOR__REDUCE:
		return "REDUCE"
	}
}

func (v FlowOperator) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.FlowOperator"
}

func (v FlowOperator) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{FLOW_OPERATOR__FILTER, FLOW_OPERATOR__MAP, FLOW_OPERATOR__FLATMAP, FLOW_OPERATOR__WINDOW, FLOW_OPERATOR__GROUP, FLOW_OPERATOR__REDUCE}
}

func (v FlowOperator) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidFlowOperator
	}
	return []byte(s), nil
}

func (v *FlowOperator) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseFlowOperatorFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *FlowOperator) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = FlowOperator(i)
	return nil
}

func (v FlowOperator) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
