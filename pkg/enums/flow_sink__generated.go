// This is a generated source file. DO NOT EDIT
// Source: enums/flow_sink__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidFlowSink = errors.New("invalid FlowSink type")

func ParseFlowSinkFromString(s string) (FlowSink, error) {
	switch s {
	default:
		return FLOW_SINK_UNKNOWN, InvalidFlowSink
	case "":
		return FLOW_SINK_UNKNOWN, nil
	case "RMDB":
		return FLOW_SINK__RMDB, nil
	case "BLOCKCHAIN":
		return FLOW_SINK__BLOCKCHAIN, nil
	}
}

func ParseFlowSinkFromLabel(s string) (FlowSink, error) {
	switch s {
	default:
		return FLOW_SINK_UNKNOWN, InvalidFlowSink
	case "":
		return FLOW_SINK_UNKNOWN, nil
	case "RMDB":
		return FLOW_SINK__RMDB, nil
	case "BLOCKCHAIN":
		return FLOW_SINK__BLOCKCHAIN, nil
	}
}

func (v FlowSink) Int() int {
	return int(v)
}

func (v FlowSink) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case FLOW_SINK_UNKNOWN:
		return ""
	case FLOW_SINK__RMDB:
		return "RMDB"
	case FLOW_SINK__BLOCKCHAIN:
		return "BLOCKCHAIN"
	}
}

func (v FlowSink) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case FLOW_SINK_UNKNOWN:
		return ""
	case FLOW_SINK__RMDB:
		return "RMDB"
	case FLOW_SINK__BLOCKCHAIN:
		return "BLOCKCHAIN"
	}
}

func (v FlowSink) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.FlowSink"
}

func (v FlowSink) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{FLOW_SINK__RMDB, FLOW_SINK__BLOCKCHAIN}
}

func (v FlowSink) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidFlowSink
	}
	return []byte(s), nil
}

func (v *FlowSink) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseFlowSinkFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *FlowSink) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = FlowSink(i)
	return nil
}

func (v FlowSink) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
