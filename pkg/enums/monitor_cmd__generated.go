// This is a generated source file. DO NOT EDIT
// Source: enums/monitor_cmd__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidMonitorCmd = errors.New("invalid MonitorCmd type")

func ParseMonitorCmdFromString(s string) (MonitorCmd, error) {
	switch s {
	default:
		return MONITOR_CMD_UNKNOWN, InvalidMonitorCmd
	case "":
		return MONITOR_CMD_UNKNOWN, nil
	case "START":
		return MONITOR_CMD__START, nil
	case "PAUSE":
		return MONITOR_CMD__PAUSE, nil
	}
}

func ParseMonitorCmdFromLabel(s string) (MonitorCmd, error) {
	switch s {
	default:
		return MONITOR_CMD_UNKNOWN, InvalidMonitorCmd
	case "":
		return MONITOR_CMD_UNKNOWN, nil
	case "START":
		return MONITOR_CMD__START, nil
	case "PAUSE":
		return MONITOR_CMD__PAUSE, nil
	}
}

func (v MonitorCmd) Int() int {
	return int(v)
}

func (v MonitorCmd) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case MONITOR_CMD_UNKNOWN:
		return ""
	case MONITOR_CMD__START:
		return "START"
	case MONITOR_CMD__PAUSE:
		return "PAUSE"
	}
}

func (v MonitorCmd) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case MONITOR_CMD_UNKNOWN:
		return ""
	case MONITOR_CMD__START:
		return "START"
	case MONITOR_CMD__PAUSE:
		return "PAUSE"
	}
}

func (v MonitorCmd) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.MonitorCmd"
}

func (v MonitorCmd) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{MONITOR_CMD__START, MONITOR_CMD__PAUSE}
}

func (v MonitorCmd) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidMonitorCmd
	}
	return []byte(s), nil
}

func (v *MonitorCmd) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseMonitorCmdFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *MonitorCmd) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = MonitorCmd(i)
	return nil
}

func (v MonitorCmd) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
