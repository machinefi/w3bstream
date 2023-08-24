package logger

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
)

func Std() logr.Logger { return &std{lvl: logr.DebugLevel} }

func StdContext(ctx context.Context) (context.Context, logr.Logger) {
	l := Std()
	return logr.WithLogger(ctx, l), l
}

type std struct {
	lvl logr.Level
	kvs []interface{}
}

func (l *std) WithValues(kvs ...interface{}) logr.Logger {
	return &std{
		lvl: l.lvl,
		kvs: append(l.kvs, kvs...),
	}
}

func (l *std) Start(ctx context.Context, name string, kvs ...interface{}) (context.Context, logr.Logger) {
	return ctx, &std{
		lvl: l.lvl,
		kvs: append(l.kvs, kvs...),
	}
}

func (l *std) End() {}

func (l *std) Debug(format string, args ...interface{}) {
	gStdLogger.LogAttrs(context.Background(), slog.LevelDebug, fmt.Sprintf(format, args...), KVsToSlogAttr(l.kvs)...)
}

func (l *std) Info(format string, args ...interface{}) {
	gStdLogger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf(format, args...), KVsToSlogAttr(l.kvs)...)
}

func (l *std) Warn(err error) {
	gStdLogger.LogAttrs(context.Background(), slog.LevelWarn, err.Error(), KVsToSlogAttr(l.kvs)...)
}

func (l *std) Error(err error) {
	gStdLogger.LogAttrs(context.Background(), slog.LevelError, err.Error(), KVsToSlogAttr(l.kvs)...)
}
