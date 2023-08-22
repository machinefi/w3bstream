package logger

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/exp/slog"

	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

func ErrIgnoreExporter() trace.SpanExporter {
	return &errIgnoreExporter{}
}

type errIgnoreExporter struct {
	trace.SpanExporter
}

func (e *errIgnoreExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	_ = e.SpanExporter.ExportSpans(ctx, spans)
	return nil
}

func (e *errIgnoreExporter) Shutdown(_ context.Context) error { return nil }

func StdoutSpanExporter() trace.SpanExporter {
	return &stdoutSpanExporter{}
}

type stdoutSpanExporter struct{}

var _ trace.SpanExporter = (*stdoutSpanExporter)(nil)

func (e *stdoutSpanExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	for i := range spans {
		data := spans[i]

		for _, event := range data.Events() {
			if event.Name == "" || event.Name[0] != '@' {
				continue
			}
			lv := slog.Level(0)
			if err := lv.UnmarshalText([]byte(strings.TrimPrefix(event.Name, "@"))); err != nil {
				continue
			}
			attr := make([]slog.Attr, 0)

			spanID := data.Name()
			if data.SpanContext().HasSpanID() {
				spanID = data.SpanContext().SpanID().String()
			}
			attr = append(attr, slog.String("span_id", spanID))
			attr = append(attr, slog.String("trace_id", data.SpanContext().TraceID().String()))
			if p := data.Parent(); p.IsValid() && p.HasSpanID() {
				attr = append(attr, slog.String("parent", p.SpanID().String()))
			}

			for _, kv := range event.Attributes {
				switch k := string(kv.Key); k {
				case "@msg":
					attr = append(attr, slog.String("@msg", kv.Value.AsString()))
				case "@stack":
					continue // too long
				case "":
					continue
				default:
					attr = append(attr, slog.Any(k, kv.Value.AsInterface()))
				}
			}
			for _, kv := range data.Attributes() {
				k := string(kv.Key)
				if k == "meta" {
					meta := metax.ParseMeta(kv.Value.AsString())
					for key := range meta {
						if key == "_id" {
							continue
						}
						attr = append(attr, slog.Any(k, meta[k]))
					}
					continue
				}
				if kv.Valid() {
					attr = append(attr, slog.Any(k, kv.Value.AsInterface()))
				}
			}

			gStdLogger.LogAttrs(ctx, lv, "", attr...)
		}
	}
	return nil
}

func (*stdoutSpanExporter) Shutdown(ctx context.Context) error { return nil }

func SpanMapExporter() trace.SpanExporter {
	return &spanMapExporter{}
}

type SpanMapper func(trace.ReadOnlySpan) trace.ReadOnlySpan

type spanMapExporter struct {
	mappers []SpanMapper
	trace.SpanExporter
}

var _ trace.SpanExporter = (*spanMapExporter)(nil)

func (e *spanMapExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	snapshots := make([]trace.ReadOnlySpan, 0)
	for i := range spans {
		data := spans[i]
		for _, m := range e.mappers {
			data = m(data)
		}
		if data == nil {
			continue
		}
		snapshots = append(snapshots, data)
	}
	if len(snapshots) == 0 {
		return nil
	}
	return e.SpanExporter.ExportSpans(ctx, snapshots)
}

func OutputFilter() SpanMapper {
	return func(span trace.ReadOnlySpan) trace.ReadOnlySpan {
		if gOutput == OUTPUT_TYPE__NEVER {
			return nil
		}
		code := span.Status().Code
		if gOutput == OUTPUT_TYPE__ON_FAILURE && code != codes.Error {
			return nil
		}
		return span
	}
}

func WithSpanMapExporter(mappers ...SpanMapper) func(trace.SpanExporter) trace.SpanExporter {
	return func(exporter trace.SpanExporter) trace.SpanExporter {
		return &spanMapExporter{
			mappers:      mappers,
			SpanExporter: exporter,
		}
	}
}

func WithErrIgnoreExporter() func(trace.SpanExporter) trace.SpanExporter {
	return func(exporter trace.SpanExporter) trace.SpanExporter {
		return &errIgnoreExporter{
			SpanExporter: exporter,
		}
	}
}
