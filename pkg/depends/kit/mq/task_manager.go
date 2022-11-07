package mq

type TaskManager interface {
	Push(ch string, t Task) error
	Pop(ch string) (Task, error)
	Remove(ch string, id string) error
	Clear(ch string) error
}
