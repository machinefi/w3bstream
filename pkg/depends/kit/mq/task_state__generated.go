// This is a generated source file. DO NOT EDIT
// Source: mq/task_state__generated.go

package mq

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidTaskState = errors.New("invalid TaskState type")

func ParseTaskStateFromString(s string) (TaskState, error) {
	switch s {
	default:
		return TASK_STATE_UNKNOWN, InvalidTaskState
	case "":
		return TASK_STATE_UNKNOWN, nil
	case "SUCCEEDED":
		return TASK_STATE__SUCCEEDED, nil
	case "FAILED":
		return TASK_STATE__FAILED, nil
	}
}

func ParseTaskStateFromLabel(s string) (TaskState, error) {
	switch s {
	default:
		return TASK_STATE_UNKNOWN, InvalidTaskState
	case "":
		return TASK_STATE_UNKNOWN, nil
	case "SUCCEEDED":
		return TASK_STATE__SUCCEEDED, nil
	case "FAILED":
		return TASK_STATE__FAILED, nil
	}
}

func (v TaskState) Int() int {
	return int(v)
}

func (v TaskState) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case TASK_STATE_UNKNOWN:
		return ""
	case TASK_STATE__SUCCEEDED:
		return "SUCCEEDED"
	case TASK_STATE__FAILED:
		return "FAILED"
	}
}

func (v TaskState) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case TASK_STATE_UNKNOWN:
		return ""
	case TASK_STATE__SUCCEEDED:
		return "SUCCEEDED"
	case TASK_STATE__FAILED:
		return "FAILED"
	}
}

func (v TaskState) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/depends/kit/mq.TaskState"
}

func (v TaskState) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{TASK_STATE__SUCCEEDED, TASK_STATE__FAILED}
}

func (v TaskState) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidTaskState
	}
	return []byte(s), nil
}

func (v *TaskState) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseTaskStateFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *TaskState) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = TaskState(i)
	return nil
}

func (v TaskState) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
