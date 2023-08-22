package job

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier"
	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier/lark"
	"github.com/machinefi/w3bstream/pkg/types"
)

func Dispatch(ctx context.Context, t mq.Task) {
	ctx, l := logr.Start(ctx, "modules.job.Dispatch",
		"subject", t.Subject(),
		"task_id", t.ID(),
	)
	defer l.End()

	tb := types.MustTaskBoardFromContext(ctx)
	tw := types.MustTaskWorkerFromContext(ctx)

	if err := tb.Dispatch(tw.Channel, t); err != nil {
		if body, err := lark.Build(ctx, "job dispatching", "WARNING", err.Error()); err != nil {
			_ = robot_notifier.Push(ctx, body, nil)
		}
		l.Error(err)
	}
}
