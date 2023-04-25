package job

import (
	"context"
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
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
				Msg:         subStringWithLength(msg, consts.WasmLogMaxLength),
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

// subStringWithLength
// If the length is negative, an empty string is returned.
// If the length is greater than the length of the input string, the entire string is returned.
// Otherwise, a substring of the input string with the specified length is returned.
func subStringWithLength(str string, length int) string {
	if length < 0 {
		return ""
	}
	rs := []rune(str)
	strLen := len(rs)

	if length > strLen {
		return str
	}
	return string(rs[0:length])
}
