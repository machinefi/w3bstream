package mq

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq/worker"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
)

type TaskWorkerOption func(*taskWorkerOption)

type taskWorkerOption struct {
	Channel     string
	WorkerCount int
	OnFinished  func(ctx context.Context, t Task)
}

func WithChannel(ch string) TaskWorkerOption {
	return func(o *taskWorkerOption) { o.Channel = ch }
}

func WithWorkerCount(cnt int) TaskWorkerOption {
	return func(o *taskWorkerOption) { o.WorkerCount = cnt }
}

func WithFinishFunc(fn func(ctx context.Context, t Task)) TaskWorkerOption {
	return func(o *taskWorkerOption) { o.OnFinished = fn }
}

func NewTaskWorker(tm TaskManager, options ...TaskWorkerOption) *TaskWorker {
	tw := &TaskWorker{mgr: tm, ops: mapx.New[string, any]()}
	for _, opt := range options {
		opt(&tw.taskWorkerOption)
	}
	return tw
}

type TaskWorker struct {
	taskWorkerOption
	mgr    TaskManager
	ops    *mapx.Map[string, any]
	worker *worker.Worker
	with   contextx.WithContext
}

func (w *TaskWorker) Context() context.Context {
	if w.with != nil {
		return w.with(context.Background())
	}
	return context.Background()
}

func (w *TaskWorker) WithContextInjector(with contextx.WithContext) *TaskWorker {
	return &TaskWorker{
		taskWorkerOption: w.taskWorkerOption,
		mgr:              w.mgr,
		ops:              mapx.New[string, any](),
		worker:           w.worker,
		with:             with,
	}
}

func (w *TaskWorker) Register(router *kit.Router) {
	for _, route := range router.Routes() {
		factories := route.OperatorFactories()
		if len(factories) != 1 {
			continue
		}
		f := factories[0]
		w.ops.Store(f.Type.Name(), f)
	}
}

func (w *TaskWorker) Serve(router *kit.Router) error {
	w.Register(router)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	w.worker = worker.New(w.proc, w.WorkerCount)
	go func() {
		w.worker.Start(w.Context())
	}()

	<-stopCh
	return errors.New("TaskWorker server closed")
}

func (w *TaskWorker) operatorFactory(ch string) (*kit.OperatorFactory, error) {
	op, ok := w.ops.Load(ch)
	if !ok || op == nil {
		return nil, errors.Errorf("missed operator %s", ch)
	}
	return op.(*kit.OperatorFactory), nil
}

func (w *TaskWorker) proc(ctx context.Context) (err error) {
	var (
		t  Task
		se error // shadowed
	)
	t, err = w.mgr.Pop(w.Channel)
	if err != nil {
		return err
	}
	if t == nil {
		return nil
	}

	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("panic: %v", e)
		}

		if err != nil {
			t.SetState(TASK_STATE__FAILED)
		} else {
			t.SetState(TASK_STATE__SUCCEEDED)
		}

		if w.OnFinished != nil {
			w.OnFinished(ctx, t)
		}
	}()

	opf, se := w.operatorFactory(t.Subject())
	if se != nil {
		err = se
		return
	}

	op := opf.New()
	if with, ok := t.(WithArg); ok {
		if setter, ok := op.(SetArg); ok {
			if se = setter.SetArg(with.Arg()); se != nil {
				err = se
				return
			}
		}
	}

	meta := metax.ParseMeta(t.ID())
	meta.Add("task", w.Channel+"#"+t.Subject())

	if _, se = op.Output(metax.ContextWithMeta(ctx, meta)); se != nil {
		err = se
		return
	}
	return
}
