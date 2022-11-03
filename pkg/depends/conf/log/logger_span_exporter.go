package log

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

type (
	SpanMapper   func(trace.ReadOnlySpan) trace.ReadOnlySpan
	WithExporter func(trace.SpanExporter) trace.SpanExporter
)

type StdExporter struct{ formatter logrus.Formatter }

func (e *StdExporter) Shutdown(ctx context.Context) error { return nil }

func (e *StdExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	for _, data := range spans {
		for _, ev := range data.Events() {
			if ev.Name == "" || ev.Name[0] != '@' {
				continue
			}
			level, err := logrus.ParseLevel(ev.Name[1:])
			if err != nil {
				continue
			}
			ent := logrus.NewEntry(logrus.StandardLogger())
			ent.Level = level
			ent.Time = ev.Time
			ent.Data = logrus.Fields{}
			for _, kv := range ev.Attributes {
				k := string(kv.Key)
				switch k {
				case "message":
					ent.Message = kv.Value.AsString()
				default:
					ent.Data[k] = kv.Value.AsInterface()
				}
			}
			for _, kv := range data.Attributes() {
				k := string(kv.Key)
				if k != "meta" {
					ent.Data[k] = kv.Value.AsInterface()
					continue
				}
				meta := metax.ParseMeta(kv.Value.AsString())
				for mk := range meta {
					if mk == "_id" {
						continue
					}
					ent.Data[mk] = meta[mk]
				}
			}
			spanCtx := data.SpanContext()
			ent.Data["span"] = data.Name()
			ent.Data["traceID"] = spanCtx.TraceID()
			if spanCtx.HasSpanID() {
				ent.Data["spanID"] = spanCtx.SpanID()
			}
			if data.Parent().IsValid() {
				ent.Data["parentSpanID"] = data.Parent().SpanID()
			}

			if content, err := e.formatter.Format(ent); err != nil {
				fmt.Println(err)
			} else {
				_, _ = os.Stdout.Write(content)
			}
		}
	}
	return nil
}

type SpanMapExporter struct {
	mappers []SpanMapper
	trace.SpanExporter
}

func WithSpanMapExporter(mappers ...SpanMapper) WithExporter {
	return func(exporters trace.SpanExporter) trace.SpanExporter {
		return &SpanMapExporter{mappers, exporters}
	}
}

func (e *SpanMapExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	snapshots := make([]trace.ReadOnlySpan, 0)

	for i := range spans {
		data := spans[i]
		for _, m := range e.mappers {
			data = m(data)
		}
		if data != nil {
			snapshots = append(snapshots, data)
		}
	}

	if len(snapshots) == 0 {
		return nil
	}
	return e.SpanExporter.ExportSpans(ctx, snapshots)
}

type ErrIgnoreExporter struct{ trace.SpanExporter }

func WithErrIgnoreExporter() WithExporter {
	return func(exporter trace.SpanExporter) trace.SpanExporter {
		return &ErrIgnoreExporter{exporter}
	}
}

func (e *ErrIgnoreExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	_ = e.SpanExporter.ExportSpans(ctx, spans)
	return nil
}

func InstallPipeline(output LoggerOutputType, format LoggerFormatType, exporter trace.SpanExporter) error {
	var stdexp *StdExporter
	switch format {
	case LOGGER_FORMAT_TYPE__JSON:
		stdexp = &StdExporter{&logrus.JSONFormatter{}}
	default:
		stdexp = &StdExporter{&logrus.TextFormatter{}}
	}

	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithSampler(trace.AlwaysSample()),
			trace.WithSyncer(
				WithSpanMapExporter(
					OutputFilter(output),
				)(stdexp)),
			trace.WithBatcher(
				WithSpanMapExporter(
					OutputFilter(output),
					SpanOnlyFilter(),
				)(WithErrIgnoreExporter()(exporter)),
			),
		),
	)

	return nil
}

func OutputFilter(output LoggerOutputType) SpanMapper {
	return func(s trace.ReadOnlySpan) trace.ReadOnlySpan {
		if output == LOGGER_OUTPUT_TYPE__NEVER {
			return nil
		}
		if output == LOGGER_OUTPUT_TYPE__ON_FAILURE && s.Status().Code == codes.Ok {
			return nil
		}
		return s
	}
}

func SpanOnlyFilter() SpanMapper {
	return func(s trace.ReadOnlySpan) trace.ReadOnlySpan {
		if s == nil {
			return nil
		}
		return s
	}
}
