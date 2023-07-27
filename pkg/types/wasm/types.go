package wasm

import (
	"context"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/types"
)

// NewLogContext create a wasm log context. this context is from wasm runtime
func NewLogContext(ctx context.Context, t LogType, lv conflog.Level, msg string) *LogContext {
	tr := MustTraceInfoFromContext(ctx)
	return &LogContext{
		LogID:       confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID(),
		Type:        t,
		Level:       lv,
		ProjectName: tr.ProjectName,
		AppletName:  tr.AppletName,
		InstanceID:  tr.InstanceID,
		Message:     msg,
	}
}

type LogContext struct {
	LogID       types.SFID
	Type        LogType
	Level       conflog.Level
	ProjectName string
	AppletName  string
	InstanceID  types.SFID
	Message     string
}
