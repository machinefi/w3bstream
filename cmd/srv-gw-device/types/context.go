package types

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

type (
	CtxDBExecutor struct{}
	CtxLogger     struct{}
	CtxCore       struct{}
)

func WithDBExecutor(ctx context.Context, v sqlx.DBExecutor) context.Context {
	return contextx.WithValue(ctx, CtxDBExecutor{}, v)
}

func WithDBExecutorContext(v sqlx.DBExecutor) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxDBExecutor{}, v)
	}
}

func DBExecutorFromContext(ctx context.Context) (sqlx.DBExecutor, bool) {
	v, ok := ctx.Value(CtxDBExecutor{}).(sqlx.DBExecutor)
	return v, ok
}

func MustDBExecutorFromContext(ctx context.Context) sqlx.DBExecutor {
	v, ok := DBExecutorFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithLogger(ctx context.Context, v log.Logger) context.Context {
	return contextx.WithValue(ctx, CtxLogger{}, v)
}

func WithLoggerContext(v log.Logger) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxLogger{}, v)
	}
}

func LoggerFromContext(ctx context.Context) (log.Logger, bool) {
	v, ok := ctx.Value(CtxLogger{}).(log.Logger)
	return v, ok
}

func MustLoggerFromContext(ctx context.Context) log.Logger {
	v, ok := LoggerFromContext(ctx)
	must.BeTrue(ok)
	return v
}

func WithCore(ctx context.Context, v *Core) context.Context {
	return contextx.WithValue(ctx, CtxCore{}, v)
}

func WithCoreContext(v *Core) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, CtxCore{}, v)
	}
}

func CoreFromContext(ctx context.Context) (*Core, bool) {
	v, ok := ctx.Value(CtxCore{}).(*Core)
	return v, ok
}

func MustCoreFromContext(ctx context.Context) *Core {
	v, ok := CoreFromContext(ctx)
	must.BeTrue(ok)
	return v
}
