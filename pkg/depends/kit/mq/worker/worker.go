package worker

import (
	"context"
	"sync"
)

func New(proc func(ctx context.Context) error, num int) *Worker {
	return &Worker{num: num, proc: proc}
}

type Worker struct {
	num  int
	proc func(ctx context.Context) error
	wg   sync.WaitGroup
}

func (w *Worker) Start(ctx context.Context) {
	w.wg.Add(w.num)

	for i := 0; i < w.num; i++ {
		go func() {
			defer w.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					w.proc(ctx)
				}
			}
		}()
	}

	w.wg.Wait()
}
