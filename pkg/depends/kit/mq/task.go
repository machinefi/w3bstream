package mq

import (
	"context"

	"github.com/google/uuid"
)

type Task interface {
	Subject() string
	ID() string
	State() TaskState
	SetState(TaskState)
}

type TaskUUID string

func (tid *TaskUUID) ID() string {
	if *tid == "" {
		*tid = TaskUUID(uuid.New().String())
	}
	return string(*tid)
}

func (tid *TaskUUID) SetID(id string) { *tid = TaskUUID(id) }

type WithArg interface {
	Arg() interface{}
}

type SetArg interface {
	SetArg(v interface{}) error
}

type TaskHeader struct {
	TaskUUID
	TaskState
	subject string
}

var _ Task = (*TaskHeader)(nil)

func (th *TaskHeader) Subject() string { return th.subject }

func (th *TaskHeader) SetSubject(s string) { th.subject = s }

func NewTaskBoard(tm TaskManager) *TaskBoard { return &TaskBoard{tm} }

type TaskBoard struct {
	tm TaskManager
}

func (b *TaskBoard) Dispatch(ch string, t Task) error {
	if t == nil {
		return nil
	}
	return b.tm.Push(ch, t)
}

type TaskOperator interface {
	Output(ctx context.Context) (interface{}, error)
	SetArg
}
