package slog

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/exp/slog"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
)

func Logger(l *slog.Logger) logr.Logger {
	return &logger{slog: l}
}

type logger struct {
	ctx   context.Context
	slog  *slog.Logger
	spans []string
}

var _ logr.Logger = (*logger)(nil)

func (d *logger) WithValues(kvs ...any) logr.Logger {
	return &logger{
		spans: d.spans,
		slog:  d.slog.With(kvs...),
	}
}

func (d *logger) Start(ctx context.Context, name string, kvs ...any) (context.Context, logr.Logger) {
	spans := append(d.spans, name)

	if len(kvs) == 0 {
		return ctx, &logger{
			ctx:   ctx,
			spans: spans,
			slog:  d.slog.WithGroup(strings.Join(spans, "/")),
		}
	}

	return ctx, &logger{
		spans: spans,
		slog:  d.slog.WithGroup(strings.Join(spans, "/")).With(kvs...),
	}
}

func (d *logger) End() {
	if len(d.spans) != 0 {
		d.spans = d.spans[0 : len(d.spans)-1]
	}
}

func (d *logger) Debug(format string, args ...any) {
	if !d.slog.Enabled(d.ctx, slog.LevelDebug) {
		return
	}
	d.slog.Log(d.ctx, slog.LevelDebug, fmt.Sprintf(format, args...))
}

func (d *logger) Info(format string, args ...any) {
	if !d.slog.Enabled(d.ctx, slog.LevelInfo) {
		return
	}
	d.slog.Log(d.ctx, slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (d *logger) Warn(err error) {
	if !d.slog.Enabled(d.ctx, slog.LevelWarn) {
		return
	}
	d.slog.Log(d.ctx, slog.LevelWarn, err.Error(), slog.Any("err", err))
}

func (d *logger) Error(err error) {
	d.slog.Error(err.Error())
}
