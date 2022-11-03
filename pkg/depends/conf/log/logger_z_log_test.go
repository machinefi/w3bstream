package log_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/trace"

	. "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

func ExampleLogger() {
	ctx := WithLogger(context.Background(), Std())

	l := FromContext(ctx).WithValues("k", "k")

	l.Debug("test %d", 1)
	l.Trace("test %d", 1)
	l.Info("test %d", 1)
	// Output:
}

func ExampleLogger_Start() {
	ctx := WithLogger(context.Background(), Std())

	_, l := Start(ctx, "span", "k", "k")
	defer l.End()

	l.Debug("test %d", 1)
	l.Trace("test %d", 1)
	l.Info("test %d", 1)
	// Output:
}

var (
	exporter, _ = zipkin.New("http://localhost:9411/api/v2/spans")
)

func TestLog(t *testing.T) {
	t.Run("OutputOnNever", func(t *testing.T) {
		var logger = Log{
			Level:    DebugLevel,
			Output:   LOGGER_OUTPUT_TYPE__NEVER,
			Format:   LOGGER_FORMAT_TYPE__TEXT,
			Exporter: exporter,
		}

		logger.SetDefault()
		logger.Init()

		ctx := metax.ContextWithMeta(context.Background(), metax.Meta{"_id": {"from context"}, "operator": {"GetByID"}})
		doLog(ctx)
	})

	t.Run("OutputOnFailure", func(t *testing.T) {
		var logger = Log{
			Level:    DebugLevel,
			Output:   LOGGER_OUTPUT_TYPE__ON_FAILURE,
			Format:   LOGGER_FORMAT_TYPE__TEXT,
			Exporter: exporter,
		}

		logger.SetDefault()
		logger.Init()

		ctx := metax.ContextWithMeta(context.Background(), metax.Meta{"_id": {"from context"}, "operator": {"GetByID"}})
		doLog(ctx)
	})

	t.Run("OutputAlways", func(t *testing.T) {
		var logger = Log{
			Output:   LOGGER_OUTPUT_TYPE__ALWAYS,
			Format:   LOGGER_FORMAT_TYPE__TEXT,
			Level:    DebugLevel,
			Exporter: exporter,
		}

		logger.SetDefault()
		logger.Init()

		ctx := metax.ContextWithMeta(context.Background(), metax.Meta{"_id": {"from context"}, "operator": {"GetByID"}})
		doLog(ctx)
	})

	t.Run("OutputAlwaysJSON", func(t *testing.T) {
		var logger = Log{
			Output:   LOGGER_OUTPUT_TYPE__ALWAYS,
			Format:   LOGGER_FORMAT_TYPE__JSON,
			Level:    DebugLevel,
			Exporter: exporter,
		}

		logger.SetDefault()
		logger.Init()

		ctx := metax.ContextWithMeta(context.Background(), metax.Meta{"_id": {"from context"}, "operator": {"GetByID"}})
		doLog(ctx)
	})
}

func doLog(ctx context.Context) {
	tracerName := "xxx"
	tracer := otel.Tracer(tracerName)

	ctx, span := tracer.Start(ctx, "op", trace.WithTimestamp(time.Now()))
	defer func() {
		span.End(trace.WithTimestamp(time.Now()))
	}()

	ctx = WithLogger(ctx, Span(tracerName, span))

	someActionWithSpan(ctx)

	otherActions(ctx)
}

func someActionWithSpan(ctx context.Context) {
	_, log := Start(ctx, "SomeActionWithSpan")
	defer log.End()

	log.Info("info xxx")
	log.Debug("debug xxx")
	log.Warn(errors.New("warn xxx"))
}

func otherActions(ctx context.Context) {
	log := FromContext(ctx)

	log.WithValues("test_key", 2).Info("test")
	log.Error(errors.New(""))
}
