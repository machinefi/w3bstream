package logr

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

type Logger interface {
	// Start to start span for tracing
	//
	// 	ctx log = log.Start(ctx, "SpanName")
	// 	defer log.End()
	//
	Start(context.Context, string, ...any) (context.Context, Logger)
	// End to end span
	End()

	// WithValues key value pairs
	WithValues(keyAndValues ...interface{}) Logger

	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(err error)
	Error(err error)
}

type keyLogger struct{}

func WithLogger(ctx context.Context, l Logger) context.Context {
	return contextx.WithValue(ctx, keyLogger{}, l)
}

func WithLoggerContext(l Logger) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return WithLogger(ctx, l)
	}
}

func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(keyLogger{}).(Logger); ok {
		return v
	}
	return Discard()
}

func Start(ctx context.Context, name string, kvs ...any) (context.Context, Logger) {
	return FromContext(ctx).Start(ctx, name, kvs...)
}
