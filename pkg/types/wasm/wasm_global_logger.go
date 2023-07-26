package wasm

import (
	"context"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Logger struct {
	logger conflog.Logger
}

func (l *Logger) GlobalConfigType() ConfigType { return ConfigLogger }

func (l *Logger) Init(parent context.Context) error {
	log := types.MustLoggerFromContext(parent)
	prj := types.MustProjectFromContext(parent)
	app := types.MustAppletFromContext(parent)
	ins := types.MustInstanceFromContext(parent)

	l.logger = log.WithValues("@src", "wasm", "prj", prj.Name, "app", app.Name, "ins", ins.InstanceID)
	return nil
}

func (l *Logger) WithContext(ctx context.Context) context.Context {
	return WithLogger(ctx, l.logger)
}
