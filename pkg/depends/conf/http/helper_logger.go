package http

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
)

func NewLogRoundTripper(logger *logrus.Entry) func(http.RoundTripper) http.RoundTripper {
	return func(roundTripper http.RoundTripper) http.RoundTripper {
		return &LogRoundTripper{
			logger: logger,
			next:   roundTripper,
		}
	}
}

type LogRoundTripper struct {
	logger *logrus.Entry
	next   http.RoundTripper
}

func (rt *LogRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	// inject h3 form context
	b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)).
		Inject(ctx, propagation.HeaderCarrier(req.Header))

	ctx, logger := log.Start(ctx, "Request")
	defer logger.End()

	rsp, err := rt.next.RoundTrip(req.WithContext(ctx))

	level, _ := log.ParseLevel(strings.ToLower(req.Header.Get("x-log-level")))
	if level == log.PanicLevel {
		level = log.TraceLevel
	}
	cost := timer.Start()

	logger.WithValues(
		"@cst", cost().Milliseconds(),
		"@mtd", req.Method[0:3],
		"@url", OmitAuthorization(req.URL),
	)

	if err == nil {
		if level >= log.InfoLevel {
			logger.Info("success")
		}
	} else {
		if level >= log.WarnLevel {
			logger.Warn(errors.Wrap(err, "http request failed"))
		}
	}
	return rsp, err
}

func TraceLogHandler(name string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			var (
				cost      = timer.StartSpan()
				tracer    = otel.Tracer(name)
				ctx, span = tracer.Start(
					b3.New().Extract(
						req.Context(),
						propagation.HeaderCarrier(req.Header),
					),
					"UnknownOperation",
					trace.WithTimestamp(cost.StartedAt()),
				)
				logger   = log.Span(name, span)
				lrw      = NewLoggerResponseWriter(rw)
				meta     = metax.ParseMeta(lrw.Header().Get("X-Meta"))
				operator = meta.Get("operator")
				level    log.Level
			)

			defer func() {
				span.End(trace.WithTimestamp(time.Now()))
			}()

			// for async pick
			b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)).
				Inject(ctx, propagation.HeaderCarrier(lrw.Header()))

			meta.Add("_id", span.SpanContext().TraceID().String())
			ctx = metax.ContextWithMeta(ctx, meta)
			ctx = log.WithLogger(ctx, logger)

			next.ServeHTTP(lrw, req.WithContext(ctx))

			if operator != "" {
				span.SetName(operator)
			}
			if lvl := strings.ToLower(req.Header.Get("x-log-level")); lvl != "" {
				level, _ = log.ParseLevel(lvl)
				if level == log.PanicLevel {
					level = log.TraceLevel
				}
			}
			kvs := []interface{}{
				"@tag", "access",
				"@rmt", httpx.ClientIP(req),
				"@cst", cost.Cost().Milliseconds(),
				"@mtd", req.Method[0:3],
				"@url", OmitAuthorization(req.URL),
				"@agent", req.Header.Get(httpx.HeaderUserAgent),
				"@status", lrw.code,
			}
			if lrw.err != nil {
				if lrw.code >= http.StatusInternalServerError {
					if level >= log.ErrorLevel {
						logger.WithValues(kvs...).Error(lrw.err)
					} else {
						if level >= log.WarnLevel {
							logger.WithValues(kvs...).Warn(lrw.err)
						}
					}
				}
			} else {
				if level >= log.InfoLevel {
					logger.WithValues(kvs...).Info("")
				}
			}
		})
	}
}

func TraceLogHandlerWithLogger(logger *logrus.Entry, name string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			var (
				cost      = timer.StartSpan()
				tracer    = otel.Tracer(name)
				ctx, span = tracer.Start(
					b3.New().Extract(
						req.Context(),
						propagation.HeaderCarrier(req.Header),
					),
					"UnknownOperation",
					trace.WithTimestamp(cost.StartedAt()),
				)
				lrw      = NewLoggerResponseWriter(rw)
				meta     = metax.ParseMeta(lrw.Header().Get("X-Meta"))
				operator = meta.Get("operator")
				level    log.Level
			)

			defer func() {
				span.End(trace.WithTimestamp(time.Now()))
			}()

			// for async pick
			b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)).
				Inject(ctx, propagation.HeaderCarrier(lrw.Header()))

			meta.Add("_id", span.SpanContext().TraceID().String())
			ctx = metax.ContextWithMeta(ctx, meta)

			next.ServeHTTP(lrw, req.WithContext(ctx))

			if operator != "" {
				span.SetName(operator)
			}
			lvl := strings.ToLower(req.Header.Get("x-log-level"))
			if level, _ = log.ParseLevel(lvl); level == log.PanicLevel {
				level = log.TraceLevel
			}
			fields := logrus.Fields{
				"@tag":    "access",
				"@cst":    cost.Cost().Milliseconds(),
				"@rmt":    httpx.ClientIP(req),
				"@mtd":    req.Method[0:3],
				"@url":    OmitAuthorization(req.URL),
				"@agent":  req.Header.Get(httpx.HeaderUserAgent),
				"@status": lrw.code,
			}
			if lrw.err != nil {
				if lrw.code >= http.StatusInternalServerError {
					if level >= log.ErrorLevel {
						logger.WithFields(fields).Error(lrw.err)
					}
				} else {
					if level >= log.WarnLevel {
						logger.WithFields(fields).Warn(lrw.err)
					}
				}
			} else {
				if level >= log.InfoLevel {
					logger.WithFields(fields).Info("")
				}
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
