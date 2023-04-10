package mws

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func MetricsHandler() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &metricsHandler{next: handler}
	}
}

type metricsHandler struct{ next http.Handler }

func (h *metricsHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet && req.URL.Path == "/metrics" {
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
		).(http.HandlerFunc)(rw, req)
		return
	}

	h.next.ServeHTTP(rw, req)
}
