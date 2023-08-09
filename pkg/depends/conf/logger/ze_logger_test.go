package logger_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/conf/tracer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

func ExampleSpanLogger() {
	{
		var c = logger.Config{
			Level:  logr.DebugLevel,
			Output: logger.OUTPUT_TYPE__NEVER,
			Format: logger.FORMAT_TYPE__TEXT,
		}

		c.SetDefault()
		if c.Init() != nil {
			return
		}

		ctx := metax.ContextWithMeta(context.Background(), metax.Meta{"_id": {"from context"}, "operator": {"GetByID"}})
		doLog(ctx, "OutputNever")
	}

	{
		var c = logger.Config{
			Level:  logr.InfoLevel,
			Output: logger.OUTPUT_TYPE__ON_FAILURE,
			Format: logger.FORMAT_TYPE__JSON,
		}

		c.SetDefault()
		if c.Init() != nil {
			return
		}

		ctx := metax.ContextWithMeta(context.Background(), metax.Meta{"_id": {"from context"}, "operator": {"GetByID"}})
		doLog(ctx, "OutputOnFailureText")
	}

	{
		var c = logger.Config{
			Output: logger.OUTPUT_TYPE__ALWAYS,
			Format: logger.FORMAT_TYPE__TEXT,
			Level:  logr.DebugLevel,
		}

		c.SetDefault()
		if c.Init() != nil {
			return
		}

		ctx := metax.ContextWithMeta(context.Background(), metax.Meta{"_id": {"from context"}, "operator": {"GetByID"}})
		doLog(ctx, "OutputAlwaysText")
	}

	{
		var c = logger.Config{
			Output: logger.OUTPUT_TYPE__ALWAYS,
			Format: logger.FORMAT_TYPE__JSON,
			Level:  logr.DebugLevel,
		}

		c.SetDefault()
		if c.Init() != nil {
			return
		}

		ctx := metax.ContextWithMeta(context.Background(), metax.Meta{"_id": {"from context"}, "operator": {"GetByID"}})
		doLog(ctx, "OutputAlwaysJson")
	}

	_ = tr.Shutdown(context.Background())

	// Output:
}

var tr *tracer.Config

func init() {
	tr = &tracer.Config{
		GrpcEndpoint: types.Endpoint{
			Scheme:   "http",
			Hostname: "localhost",
			Port:     4317,
		},
		ServiceName:    "test_log",
		ServiceVersion: "1.0.0",
		InstanceID:     uuid.NewString(),
		DebugMode:      datatypes.TRUE,
	}
	tr.SetDefault()
	if err := tr.Init(); err != nil {
		panic(err)
	}
}

func doLog(ctx context.Context, name string) {
	println(name)
	defer println()

	t := otel.Tracer(name)

	ctx, span := t.Start(ctx, "op", trace.WithTimestamp(time.Now()))
	defer func() {
		span.End(trace.WithTimestamp(time.Now()))
	}()

	ctx = logr.WithLogger(ctx, logger.SpanLogger(span))

	someActionWithSpanAndLog(ctx)

	otherActionsLogOnly(ctx)
}

func someActionWithSpanAndLog(ctx context.Context) {
	_, l := logr.Start(ctx, "SomeActionWithSpan")
	defer l.End()

	l.Info("msg")
	l.Debug("msg")
	l.Warn(errors.New("err"))
	l.Error(errors.New("err"))
}

func otherActionsLogOnly(ctx context.Context) {
	l := logr.FromContext(ctx)

	l.WithValues("test_key", 2).Info("test")
	l.Error(errors.New("any"))
}
