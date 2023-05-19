package event

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/event"), &middleware.ReqRateLimit{})

func init() {
	Root.Register(kit.NewRouter(&HandleEvent{}))
}

var _eventMtc = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "inbound_events_metrics",
		Help: "received events metrics.",
	},
	[]string{"account", "project", "publisher", "eventtype"},
)

func init() {
	prometheus.MustRegister(_eventMtc)
}
