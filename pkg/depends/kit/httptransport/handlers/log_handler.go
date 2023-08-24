package handlers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
)

func LogHandler() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &loggerHandler{
			next: handler,
		}
	}
}

type loggerHandler struct {
	next http.Handler
}

type LoggerResponseWriter struct {
	rw         http.ResponseWriter
	written    bool
	statusCode int
	err        error
}

func (rw *LoggerResponseWriter) Header() http.Header { return rw.rw.Header() }

func (rw *LoggerResponseWriter) WriteErr(err error) { rw.err = err }

func (rw *LoggerResponseWriter) WriteHeader(sc int) {
	if rw.written {
		return
	}
	rw.rw.WriteHeader(sc)
	rw.statusCode = sc
	rw.written = true
}

func (rw *LoggerResponseWriter) Write(data []byte) (int, error) {
	if rw.err != nil && rw.statusCode >= http.StatusBadRequest {
		rw.err = errors.New(string(data))
	}
	return rw.rw.Write(data)
}

func (h *loggerHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	cost := timer.Start()
	reqID := req.Header.Get(httpx.HeaderRequestID)
	if reqID == "" {
		reqID = uuid.New().String()
	}

	var (
		w = &LoggerResponseWriter{rw: rw}
		l = logr.FromContext(req.Context())
	)

	defer func() {
		header := req.Header
		duration := strconv.FormatInt(cost().Microseconds(), 10) + "μs"
		fields := []interface{}{
			"@cst", duration,
			"@rmt", httpx.ClientIP(req),
			"@mtd", req.Method[0:3],
			"@url", req.URL.String(),
			"@agent", header.Get(httpx.HeaderUserAgent),
			"@status", w.statusCode,
		}
		if w.err != nil {
			if w.statusCode >= http.StatusInternalServerError {
				l.WithValues(fields).Error(w.err)
			} else {
				l.WithValues(fields).Warn(w.err)
			}
		} else {
			l.WithValues(fields).Info("")
		}
	}()

	h.next.ServeHTTP(
		w,
		req.WithContext(
			metax.ContextWithMeta(req.Context(), metax.ParseMeta(reqID)),
		),
	)
}
