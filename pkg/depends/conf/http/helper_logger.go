package http

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
)

func TraceLogHandler(tr trace.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = b3.New().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			ctx, span := tr.Start(ctx, "Operator", trace.WithTimestamp(time.Now()))
			defer func() {
				span.End(trace.WithTimestamp(time.Now()))
			}()

			var (
				l   = logger.SpanLogger(span)
				lrw = NewLoggerResponseWriter(rw)
			)

			b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)).Inject(ctx, propagation.HeaderCarrier(lrw.Header()))
			meta := metax.ParseMeta(lrw.Header().Get("X-Meta"))
			meta["_id"] = []string{span.SpanContext().TraceID().String()}

			ctx = metax.ContextWithMeta(ctx, meta)
			ctx = logr.WithLogger(ctx, l)

			cost := timer.Start()
			next.ServeHTTP(lrw, r.WithContext(ctx))
			duration := strconv.FormatInt(cost().Microseconds(), 10) + "μs"

			operator := metax.ParseMeta(lrw.Header().Get("X-Meta")).Get("operator")
			if operator == "" {
				operator = lrw.Header().Get("X-Meta")
			}
			if operator != "" {
				span.SetName(operator)
			}

			kvs := []interface{}{
				"@tag", "access",
				"@rmt", httpx.ClientIP(r),
				"@cst", duration,
				"@mtd", r.Method,
				"@url", OmitAuthorization(r.URL),
				"@code", lrw.code,
			}

			if lrw.err != nil {
				if lrw.code >= http.StatusInternalServerError {
					l.WithValues(kvs...).Error(lrw.err)
				} else {
					l.WithValues(kvs...).Warn(lrw.err)
				}
			} else {
				l.WithValues(kvs...).Info("")
			}
		})
	}
}

func NewLoggerResponseWriter(rw http.ResponseWriter) *LoggerResponseWriter {
	lrw := &LoggerResponseWriter{ResponseWriter: rw}
	if v, ok := rw.(http.Hijacker); ok {
		lrw.Hijacker = v
	}
	if v, ok := rw.(http.Flusher); ok {
		lrw.Flusher = v
	}
	return lrw
}

type LoggerResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
	http.Flusher

	written bool
	code    int
	err     error
}

func (rw *LoggerResponseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *LoggerResponseWriter) WriteHeader(sc int) {
	if !rw.written {
		rw.ResponseWriter.WriteHeader(sc)
		rw.code = sc
		rw.written = true
	}
}

func (rw *LoggerResponseWriter) Write(data []byte) (int, error) {
	if rw.err == nil && rw.code >= http.StatusBadRequest {
		rw.err = errors.New(string(data))
	}
	return rw.ResponseWriter.Write(data)
}

func OmitAuthorization(u *url.URL) string {
	query := u.Query()
	query.Del("authorization")
	u.RawQuery = query.Encode()
	return u.String()
}
