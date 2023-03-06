package job

import (
	"context"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Dispatcher = func(mq.Task)

func NewDispatcher(ctx context.Context) Dispatcher {
	l, ok := types.LoggerFromContext(ctx)
	if !ok {
		l = conflog.Std()
	}

	tb := types.MustTaskBoardFromContext(ctx)
	tw := types.MustTaskWorkerFromContext(ctx)

	return func(t mq.Task) {
		l = l.WithValues("subject", t.Subject(), "task_id", t.ID())
		if err := tb.Dispatch(tw.Channel, t); err != nil {
			l.Error(err)
			return
		}
		l.Info("")
	}
}
