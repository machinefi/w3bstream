package job

import (
	"context"
	"fmt"
	"time"

	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func NewWasmLogTask(ctx context.Context, logLevel, msg string) *WasmLogTask {
	return &WasmLogTask{
		wasmLog: &models.WasmLog{
			RelWasmLog: models.RelWasmLog{WasmLogID: confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID()},
			WasmLogInfo: models.WasmLogInfo{
				ProjectName: types.MustProjectFromContext(ctx).ProjectName.Name,
				AppletName:  types.MustAppletFromContext(ctx).Name,
				InstanceID:  types.MustInstanceFromContext(ctx).InstanceID,
				Level:       logLevel,
				LogTime:     base.AsTimestamp(time.Now()),
				Msg:         msg,
			},
		},
	}
}

type WasmLogTask struct {
	wasmLog *models.WasmLog
	mq.TaskState
}

var _ mq.Task = (*WasmLogTask)(nil)

func (t *WasmLogTask) Subject() string {
	return "DbLogStoring"
}

func (t *WasmLogTask) Arg() interface{} {
	return t.wasmLog
}

func (t *WasmLogTask) ID() string {
	return fmt.Sprintf("%s::%s", t.Subject(), t.wasmLog.WasmLogID)
}
