package job

import (
	"fmt"
	"time"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func NewWasmLogTask(v *wasm.LogContext) *WasmLogTask {
	l := conflog.Std()
	l.Debug(fmt.Sprintf("new log task with %s-%s", v.Type, v.Message))
	task := &WasmLogTask{
		wasmLog: &models.WasmLog{
			RelWasmLog: models.RelWasmLog{WasmLogID: v.LogID},
			WasmLogInfo: models.WasmLogInfo{
				ProjectName: v.ProjectName,
				AppletName:  v.AppletName,
				InstanceID:  v.InstanceID,
				Src:         string(v.Type),
				Level:       v.Level.String(),
				LogTime:     time.Now().UnixNano(),
				Msg:         subStringWithLength(v.Message, enums.WasmLogMaxLength),
			},
		},
	}
	l.Debug(fmt.Sprintf("log record is %v", task.wasmLog))
	return task
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
