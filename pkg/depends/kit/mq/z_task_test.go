package mq_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq/mem_mq"
)

func NewTask(subject, id string, arg ...interface{}) *Task {
	t := &Task{}
	t.SetSubject(subject)

	if id != "" {
		t.SetID(id)
	}

	if len(arg) > 0 {
		t.SetArg(arg[0])
	}
	return t
}

type Task struct {
	mq.TaskHeader
	argv string
}

func (t *Task) Arg() interface{} { return t.argv }

func (t *Task) SetArg(v interface{}) {
	switch argv := v.(type) {
	case string:
		t.argv = argv
	case []byte:
		t.argv = string(argv)
	}
}

var (
	managers = []mq.TaskManager{mem_mq.New(10000) /*redis_mq.New()*/}
	channel  = "cc"
	router   = kit.NewRouter()
)

func init() {
	for _, m := range managers {
		_ = m.Clear(channel)
	}
}

type OpA struct{}

func (a *OpA) Output(ctx context.Context) (interface{}, error) {
	fmt.Println(metax.GetMetaFrom(ctx))
	return nil, nil
}

type OpB struct {
	bytes.Buffer
}

func (b *OpB) SetArg(v interface{}) error {
	var data []byte

	switch arg := v.(type) {
	case []byte:
		data = arg
	case string:
		data = []byte(arg)
	default:
		return nil
	}
	_, err := b.Write(data)
	return err
}

func (b *OpB) Output(ctx context.Context) (interface{}, error) {
	fmt.Println(metax.GetMetaFrom(ctx), b.String())
	return nil, nil
}

func TestTaskWorker(t *testing.T) {
	for _, tm := range managers {
		tb := mq.NewTaskBoard(tm)
		n := 100

		for i := 0; i < n; i++ {
			for j := 0; j < 5; j++ {
				_ = tb.Dispatch(
					channel,
					NewTask("OpA", fmt.Sprintf("OpA%d", i), []byte("A")),
				)
				_ = tb.Dispatch(
					channel,
					NewTask("OpB", fmt.Sprintf("OpB%d", i), []byte("B")),
				)
			}
		}
		router.Register(kit.NewRouter(&OpA{}))
		router.Register(kit.NewRouter(&OpB{}))

		tw := mq.NewTaskWorker(tm,
			mq.WithChannel(channel),
			mq.WithWorkerCount(2),
			mq.WithFinishFunc(func(ctx context.Context, task mq.Task) {
				NewWithT(t).Expect(task.State()).To(Equal(mq.TASK_STATE__SUCCEEDED))
			}),
		)

		go func() {
			fmt.Println(tw.Serve(router))
		}()

		time.Sleep(800 * time.Microsecond)

		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(os.Interrupt)

		time.Sleep(time.Second)
	}
}
