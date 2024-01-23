package job

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	confmq "github.com/machinefi/w3bstream/pkg/depends/conf/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier"
	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier/lark"
)

var channels sync.Map

func AddChannel(ch string) {
	channels.Store(ch, time.NewTicker(time.Minute*5))
}

func Ticker(ch string) *time.Ticker {
	t, ok := channels.Load(ch)
	if !ok {
		return nil
	}
	return t.(*time.Ticker)
}

func Dispatch(ctx context.Context, t mq.Task) {
	ctx, l := logr.Start(ctx, "modules.job.Dispatch",
		"subject", t.Subject(),
		"task_id", t.ID(),
	)
	defer l.End()

	tasks := confmq.MustMqFromContext(ctx)
	ch := tasks.TaskWorker.Channel

	if err := tasks.TaskBoard.Dispatch(tasks.TaskWorker.Channel, t); err != nil {
		tik := Ticker(ch)
		if tik == nil {
			AddChannel(ch)
		}
		select {
		case <-Ticker(ch).C:
			break
		default:
			if body, _err := lark.Build(ctx, "job dispatching", "WARNING", err.Error()); _err != nil {
				l.Warn(errors.Wrap(_err, "build lark message"))
			} else {
				if _err = robot_notifier.Push(ctx, body, nil); _err != nil {
					l.Warn(errors.Wrap(_err, "notifier push message"))
				}
			}

		}
		l.Error(err)
	}
}
