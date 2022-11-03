package log

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/textx"
)

func Span(tracer string, s trace.Span) Logger {
	return &span{lvl: DebugLevel, tracer: tracer, span: s}
}

func SpanContext(ctx context.Context, tracer string) (context.Context, Logger) {
	ctx, sp := otel.Tracer(tracer).Start(ctx, tracer, trace.WithTimestamp(time.Now()))
	l := Span(tracer, sp)
	return WithLogger(ctx, l), l
}

type span struct {
	lvl    Level
	tracer string
	span   trace.Span
	attrs  []attribute.KeyValue
}

func (l *span) SetLevel(lvl Level) Logger {
	logger := ptrx.Ptr(*l)
	logger.lvl = lvl
	return logger
}

func (l *span) Start(ctx context.Context, name string, kvs ...interface{}) (context.Context, Logger) {
	sp := trace.SpanFromContext(ctx)
	meta := metax.GetMetaFrom(ctx)

	if len(meta) > 0 {
		kvs = append(kvs, "meta", meta)
	}
	ctx, sp = sp.TracerProvider().Tracer(l.tracer).Start(
		ctx, name,
		trace.WithAttributes(toAttributes(kvs...)...),
		trace.WithTimestamp(time.Now()),
	)
	return ctx, &span{span: sp, tracer: l.tracer, lvl: l.lvl}
}

func (l *span) End() { l.span.End(trace.WithTimestamp(time.Now())) }

func (l *span) WithValues(kvs ...interface{}) Logger {
	return &span{
		tracer: l.tracer,
		span:   l.span,
		attrs:  append(l.attrs, toAttributes(kvs...)...),
		lvl:    l.lvl,
	}
}

func (l *span) Trace(format string, args ...interface{}) {
	l.info(TraceLevel, &printer{format, args})
}

func (l *span) Debug(format string, args ...interface{}) {
	l.info(DebugLevel, &printer{format, args})
}

func (l *span) Info(format string, args ...interface{}) {
	l.info(InfoLevel, &printer{format, args})
}

func (l *span) Warn(err error) { l.error(WarnLevel, err) }

func (l *span) Error(err error) { l.error(ErrorLevel, err) }

func (l *span) Fatal(err error) { l.error(FatalLevel, err) }

func (l *span) Panic(err error) { l.error(PanicLevel, err); panic(err) }

func (l *span) info(lvl Level, msg *printer) {
	if lvl > l.lvl {
		return
	}

	l.span.AddEvent(
		"@"+lvl.String(),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(l.attrs...),
		trace.WithAttributes(attribute.Stringer("message", msg)),
	)
}
func (l *span) error(lvl Level, err error) {
	if lvl > l.lvl {
		return
	}

	if l.span == nil || err == nil || !l.span.IsRecording() {
		return
	}

	attrs := append(l.attrs, attribute.String("message", err.Error()))

	if lvl <= ErrorLevel {
		attrs = append(attrs, attribute.String("stack", fmt.Sprintf("%+v", err)))
	}

	l.span.SetStatus(codes.Error, "")
	l.span.AddEvent(
		"@"+lvl.String(),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(attrs...),
	)
}

func toAttributes(kvs ...interface{}) []attribute.KeyValue {
	n := len(kvs)
	if n > 0 && n%2 == 0 {
		attrs := make([]attribute.KeyValue, n/2)
		for i := range attrs {
			k, v := kvs[2*i], kvs[2*i+1]

			ks, ok := k.(string)
			if !ok {
				continue
			}
			vs, err := textx.MarshalText(v)
			if err != nil {
				continue
			}
			attrs[i] = attribute.String(ks, string(vs))
		}
		return attrs
	}
	return nil

}

type printer struct {
	format string
	args   []interface{}
}

func (p *printer) String() string { return fmt.Sprintf(p.format, p.args...) }
