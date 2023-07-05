package event

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

// TODO: Unify the event interface
var (
	Root  = kit.NewRouter(httptransport.Group("/event"), &middleware.EventReqRateLimit{})
	Root2 = kit.NewRouter(httptransport.Group("/event2"), &middleware.EventReqRateLimit{})
)

func init() {
	Root.Register(kit.NewRouter(&HandleEvent{}))
	Root2.Register(kit.NewRouter(&HandleDataPush{}))
}
