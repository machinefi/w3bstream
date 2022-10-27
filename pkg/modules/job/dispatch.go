package job

import (
	"context"

	"github.com/iotexproject/Bumblebee/kit/mq"
	"github.com/iotexproject/w3bstream/pkg/types"
)

func Dispatch(ctx context.Context, t mq.Task) {
	l := types.MustLoggerFromContext(ctx)
	tb := types.MustTaskBoardFromContext(ctx)
	tw := types.MustTaskWorkerFromContext(ctx)

	_, l = l.WithValues(
		"subject", t.Subject(),
		"task_id", t.ID(),
	).Start(ctx, "Dispatch")
	l.Info("")

	if err := tb.Dispatch(tw.Channel, t); err != nil {
		l.Error(err)
	}
}
