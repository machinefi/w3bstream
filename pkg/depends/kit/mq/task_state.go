package mq

//go:generate toolkit gen enum TaskState
type TaskState uint8

const (
	TASK_STATE_UNKNOWN TaskState = iota
	TASK_STATE__SUCCEEDED
	TASK_STATE__FAILED
)

var TASK_STATE__PENDING = TASK_STATE_UNKNOWN

func (v TaskState) State() TaskState { return v }

func (v *TaskState) SetState(s TaskState) { *v = s }
