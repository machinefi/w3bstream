package mem_mq_test

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq/mem_mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq/worker"
)

var (
	tm = mem_mq.New(10000)
	ch = "mm"
)

func NewTask(subject, id string, payload ...interface{}) *Task {
	t := &Task{}
	t.SetSubject(subject)

	if id != "" {
		t.SetID(id)
	}

	if len(payload) == 0 {
		return t
	}
	if pl, ok := payload[0].([]byte); ok {
		t.Write(pl)
	}

	return t
}

type Task struct {
	mq.TaskHeader
	bytes.Buffer
}

func (t *Task) Payload() []byte { return t.Bytes() }

func BenchmarkTaskManager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = tm.Push(ch, NewTask("", fmt.Sprintf("%d", i), nil))
		_, _ = tm.Pop(ch)
	}
}

func TestTaskManager(t *testing.T) {
	_ = tm.Clear(ch)

	for i := 0; i < 1000; i++ {
		_ = tm.Push(ch, NewTask("", fmt.Sprintf("%d", i), nil))
		_ = tm.Push(ch, NewTask("", fmt.Sprintf("%d", i), nil))
		_ = tm.Push(ch, NewTask("", fmt.Sprintf("%d", i), nil))
		_ = tm.Push(ch, NewTask("", fmt.Sprintf("%d", i), nil))
		_ = tm.Push(ch, NewTask("", fmt.Sprintf("%d", i), nil))
	}

	wg := sync.WaitGroup{}
	wg.Add(1000)

	w := worker.New(func(ctx context.Context) error {
		task, err := tm.Pop(ch)
		if err != nil {
			return err
		}
		if task == nil {
			return nil
		}
		wg.Add(-1)
		return nil
	}, 10)

	go w.Start(context.Background())
	wg.Wait()
}
