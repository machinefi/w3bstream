package custommetrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	Metrics interface {
		Counter(string) Counter
		Gauge(string) Gauge
	}

	Counter interface {
		Inc()
		Add(float64)
	}

	Gauge interface {
		Set(float64)
		// Add(float64)
		// Sub(float64)
	}
)

var (
	_labels               = []string{"account", "project", "customlabel"}
	_customMetricsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wasm_custom_counter_metrics",
			Help: "custom counter metrics emitted from wasm vm",
		},
		_labels,
	)
	_customMetricsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wasm_custom_gague_metrics",
			Help: "custom gague metrics emitted from wasm vm",
		},
		_labels,
	)
)

func init() {
	prometheus.MustRegister(_customMetricsCounter, _customMetricsGauge)
}

type (
	metrics struct {
		preDefLabels []string
		counters     sync.Map
		gagues       sync.Map
	}
)

func NewCustomMetric(account string, project string) Metrics {
	return &metrics{
		preDefLabels: []string{account, project},
		counters:     sync.Map{},
		gagues:       sync.Map{},
	}
}

func (m *metrics) Counter(customLabel string) Counter {
	value, exist := m.counters.Load(customLabel)
	if !exist {
		value = &counter{
			labels:  append(m.preDefLabels[:], customLabel),
			counter: _customMetricsCounter,
		}
		m.counters.Store(customLabel, value)
	}
	return value.(Counter)
}

func (m *metrics) Gauge(customLabel string) Gauge {
	value, exist := m.gagues.Load(customLabel)
	if !exist {
		value = &gauge{
			labels: append(m.preDefLabels[:], customLabel),
			gauge:  _customMetricsGauge,
		}
		m.gagues.Store(customLabel, value)
	}
	return value.(Gauge)
}

type (
	counter struct {
		labels  []string
		counter *prometheus.CounterVec
	}
)

func (c *counter) Inc() {
	c.counter.WithLabelValues(c.labels...).Inc()
}

func (c *counter) Add(val float64) {
	c.counter.WithLabelValues(c.labels...).Add(val)
}

type (
	gauge struct {
		labels []string
		gauge  *prometheus.GaugeVec
	}
)

func (g *gauge) Set(val float64) {
	g.gauge.WithLabelValues(g.labels...).Set(val)
}

// func (g *gauge) Add(val float64) {
// 	g.gauge.WithLabelValues(g.labels...).Add(val)
// }

// func (g *gauge) Sub(val float64) {
// 	g.gauge.WithLabelValues(g.labels...).Sub(val)
// }
