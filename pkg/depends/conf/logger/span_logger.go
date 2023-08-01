package logger

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

func SpanLogger(span trace.Span) logr.Logger {
	return &spanLogger{span: span}
}

func NewSpanContext(ctx context.Context, name string) (context.Context, logr.Logger) {
	ctx, span := otel.Tracer(name).Start(ctx, name, trace.WithTimestamp(time.Now()))
	l := SpanLogger(span)
	return logr.WithLogger(ctx, l), l
}

type spanLogger struct {
	span trace.Span
	attr []attribute.KeyValue
}

var _ logr.Logger = (*spanLogger)(nil)

func (l *spanLogger) Start(ctx context.Context, name string, kvs ...any) (context.Context, logr.Logger) {
	span := trace.SpanFromContext(ctx)
	meta := metax.GetMetaFrom(ctx)

	if len(meta) > 0 {
		kvs = append(kvs, "meta", meta)
	}

	tp := span.TracerProvider()
	ctx, span = tp.Tracer(name).Start(ctx, name,
		trace.WithAttributes(KVsToAttr(kvs...)...),
		trace.WithTimestamp(time.Now()),
	)
	return ctx, &spanLogger{span: span}
}

func (l *spanLogger) End() {
	l.span.End(trace.WithTimestamp(time.Now()))
}

func (l *spanLogger) WithValues(kvs ...any) logr.Logger {
	return &spanLogger{
		span: l.span,
		attr: append(l.attr, KVsToAttr(kvs...)...),
	}
}

func (l *spanLogger) Debug(format string, args ...any) {
	l.info(slog.LevelDebug, format, args...)
}

func (l *spanLogger) Info(format string, args ...any) {
	l.info(slog.LevelInfo, format, args...)
}

func (l *spanLogger) Warn(err error) {
	l.error(slog.LevelWarn, err)
}

func (l *spanLogger) Error(err error) {
	l.error(slog.LevelError, err)
}

func (l *spanLogger) info(lv slog.Level, format string, args ...any) {
	msg := ""
	if len(args) == 0 {
		msg = format
	}
	msg = fmt.Sprintf(format, args...)

	options := []trace.EventOption{
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(attribute.String("@msg", msg)),
	}
	if len(l.attr) > 0 {
		options = append(options, trace.WithAttributes(l.attr...))
	}

	l.span.AddEvent("@"+lv.String(), options...)
}

func (l *spanLogger) error(lv slog.Level, err error) {
	if l.span == nil || err == nil || !l.span.IsRecording() {
		return
	}

	kvs := append(l.attr,
		attribute.String("@msg", err.Error()),
		attribute.String("@stack", fmt.Sprintf("%+v", err)),
	)

	l.span.SetStatus(codes.Error, err.Error())
	l.span.AddEvent(
		"@"+lv.String(),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(kvs...),
	)
}
