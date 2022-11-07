package worker_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/mq/worker"
)

func TestWorker(t *testing.T) {
	count := int64(0)

	w := worker.New(func(ctx context.Context) error {
		c := atomic.LoadInt64(&count)
		atomic.StoreInt64(&count, c+1)
		time.Sleep(100 * time.Millisecond)
		return nil
	}, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	w.Start(ctx)
	t.Log(count)
}
